// Command bwapi is a disposable spike CLI for Issue #53 Step 0.
// It is NOT part of the production bwsf binary.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bnema/bitwarden-go-sdk/bitwarden"
)

const (
	spikeFolderName = "bwsf-spike"
	spikeNoteName   = "bwsf-spike-note"

	envEmail        = "BWSF_SPIKE_EMAIL"
	envPassword     = "BWSF_SPIKE_PASSWORD"
	envServerURL    = "BWSF_SPIKE_SERVER_URL"
	envTOTP         = "BWSF_SPIKE_TOTP"
	envClientID     = "BWSF_SPIKE_CLIENT_ID"
	envClientSecret = "BWSF_SPIKE_CLIENT_SECRET"
	envIdentityURL  = "BWSF_SPIKE_IDENTITY_URL"

	cloudIdentityUS = "https://identity.bitwarden.com"

	// Vaultwarden rejects /connect/token when deviceName is blank.
	// Official Bitwarden cloud is more lenient; always send device metadata.
	spikeDeviceType       = "8" // LinuxDesktop (numeric; string names also OK on VW)
	spikeDeviceIdentifier = "bwsf-spike-device"
	spikeDeviceName       = "bwsf-spike"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "apikey", "--apikey":
			if err := runScenarioB(context.Background()); err != nil {
				fail("Scenario B", err)
			}
			return
		case "help", "-h", "--help":
			printUsage()
			return
		default:
			fmt.Fprintf(os.Stderr, "unknown argument %q\n\n", os.Args[1])
			printUsage()
			os.Exit(2)
		}
	}

	if err := runScenarioA(context.Background()); err != nil {
		fail("Scenario A", err)
	}
}

func fail(scenario string, err error) {
	fmt.Fprintf(os.Stderr, "FAIL (%s): %v\n", scenario, err)
	explainLoginFailure(err)
	os.Exit(1)
}

// explainLoginFailure prints SDK error fields and a short hint when the
// opaque "unknown status=400" pattern appears (SDK drops response bodies when
// Vaultwarden's OAuth error JSON has an empty "error" string).
func explainLoginFailure(err error) {
	var bwErr *bitwarden.Error
	if !errors.As(err, &bwErr) || bwErr == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "SDK error detail: kind=%s status=%d code=%q message=%q op=%q\n",
		string(bwErr.Kind), bwErr.StatusCode, bwErr.Code, bwErr.Message, bwErr.Op)
	if bwErr.Kind == bitwarden.ErrorKindUnknown && bwErr.StatusCode == 400 {
		fmt.Fprintf(os.Stderr, "hint: opaque HTTP 400 often means Vaultwarden rejected /identity/connect/token\n")
		fmt.Fprintf(os.Stderr, "      (e.g. missing deviceName) while SDK discarded the response body.\n")
		fmt.Fprintf(os.Stderr, "      Ensure LoginOptions includes DeviceName/DeviceType/DeviceIdentifier.\n")
		fmt.Fprintf(os.Stderr, "      See README.md § Vaultwarden troubleshooting.\n")
	}
}

func printUsage() {
	fmt.Fprintf(os.Stdout, `bwapi — Issue #53 Step 0 spike (NOT production bwsf)

Usage:
  bwapi              Run Scenario A (master-password login + vault CRUD)
  bwapi apikey       Run Scenario B (Personal API Key → access_token probe)
  bwapi --apikey     Same as apikey

Scenario A env:
  %s            (required)
  %s         (required)
  %s     (optional; self-hosted / Vaultwarden; must be https)
  %s             (optional; 2FA authenticator code; prompts if missing)

Scenario B env:
  %s       (required)
  %s   (required)
  %s     (optional; derives identity as <url>/identity)
  %s    (optional; overrides identity base)

See README.md for success criteria and Vaultwarden smoke steps.
`, envEmail, envPassword, envServerURL, envTOTP, envClientID, envClientSecret, envServerURL, envIdentityURL)
}

// ---------------------------------------------------------------------------
// Scenario A — main go/no-go
// ---------------------------------------------------------------------------

func runScenarioA(ctx context.Context) error {
	email := os.Getenv(envEmail)
	password := os.Getenv(envPassword)
	if email == "" || password == "" {
		return fmt.Errorf("%s and %s are required", envEmail, envPassword)
	}

	opts := []bitwarden.Option{}
	var derived identityAPIURLs
	if serverURL := strings.TrimSpace(os.Getenv(envServerURL)); serverURL != "" {
		normalized, err := normalizeServerURL(serverURL)
		if err != nil {
			return err
		}
		if normalized != serverURL {
			fmt.Printf("NewClient: WithServerURL(%s)  (normalized from %q)\n", normalized, serverURL)
		} else {
			fmt.Printf("NewClient: WithServerURL(%s)\n", normalized)
		}
		derived = deriveSelfHostedURLs(normalized)
		fmt.Printf("Derived identity=%s\n", derived.Identity)
		fmt.Printf("Derived api=%s\n", derived.API)
		opts = append(opts, bitwarden.WithServerURL(normalized))
	} else {
		fmt.Println("NewClient: cloud default (US)")
		fmt.Printf("Identity (cloud US)=%s\n", cloudIdentityUS)
	}

	client, err := bitwarden.NewClient(opts...)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	fmt.Printf("BeginLogin… (deviceName=%s deviceType=%s)\n", spikeDeviceName, spikeDeviceType)
	result, err := client.BeginLogin(ctx, bitwarden.LoginOptions{
		Email:            email,
		Password:         password,
		DeviceType:       spikeDeviceType,
		DeviceIdentifier: spikeDeviceIdentifier,
		DeviceName:       spikeDeviceName,
	})
	if err != nil {
		if derived.Identity != "" {
			diagnoseConnectToken(ctx, derived.Identity)
		}
		return fmt.Errorf("BeginLogin: %w", err)
	}

	if result.Challenge != nil {
		providers := result.Challenge.Providers()
		fmt.Printf("2FA required; providers: %v\n", providers)
		code := strings.TrimSpace(os.Getenv(envTOTP))
		if code == "" {
			fmt.Print("Enter 2FA code (authenticator): ")
			line, readErr := bufio.NewReader(os.Stdin).ReadString('\n')
			if readErr != nil {
				result.Challenge.Close()
				return fmt.Errorf("read 2FA code: %w", readErr)
			}
			code = strings.TrimSpace(line)
		}
		if code == "" {
			result.Challenge.Close()
			return fmt.Errorf("2FA code empty (set %s or type at prompt)", envTOTP)
		}
		fmt.Println("CompleteLogin…")
		result, err = client.CompleteLogin(ctx, bitwarden.CompleteLoginOptions{
			Challenge: result.Challenge,
			Provider:  bitwarden.TwoFactorProviderAuthenticator,
			Code:      code,
			Remember:  false,
		})
		if err != nil {
			return fmt.Errorf("CompleteLogin: %w", err)
		}
	}

	fmt.Printf("OK login account=%s\n", result.AccountID)

	fmt.Println("Sync…")
	if err := client.Sync(ctx); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}
	fmt.Println("OK Sync")

	folderID, err := ensureSpikeFolder(ctx, client)
	if err != nil {
		return err
	}
	fmt.Printf("OK folder %q id=%s\n", spikeFolderName, folderID)

	if err := exerciseSecureNote(ctx, client, folderID); err != nil {
		return err
	}

	fmt.Println("SUCCESS (Scenario A): login, sync, folder, secure-note create/get/update all OK")
	return nil
}

func ensureSpikeFolder(ctx context.Context, client *bitwarden.Client) (string, error) {
	folders, err := client.Folders().List(ctx)
	if err != nil {
		return "", fmt.Errorf("Folders.List: %w", err)
	}
	for _, f := range folders {
		if f.Name == spikeFolderName {
			fmt.Printf("Folder %q already exists\n", spikeFolderName)
			return f.ID, nil
		}
	}
	fmt.Printf("Creating folder %q…\n", spikeFolderName)
	created, err := client.Folders().Create(ctx, spikeFolderName)
	if err != nil {
		return "", fmt.Errorf("Folders.Create: %w", err)
	}
	return created.ID, nil
}

func exerciseSecureNote(ctx context.Context, client *bitwarden.Client, folderID string) error {
	vault := client.Vault()

	fmt.Printf("Creating Secure Note %q…\n", spikeNoteName)
	created, err := vault.Create(ctx, bitwarden.Item{
		Type:     bitwarden.ItemTypeSecureNote,
		Name:     spikeNoteName,
		Notes:    "bwsf spike initial notes",
		FolderID: folderID,
	})
	if err != nil {
		return fmt.Errorf("Vault.Create(secure_note): %w", err)
	}
	fmt.Printf("OK create id=%s\n", created.ID)

	fmt.Printf("Get by name via Search(%q)…\n", spikeNoteName)
	found, err := vault.Search(ctx, spikeNoteName)
	if err != nil {
		return fmt.Errorf("Vault.Search: %w", err)
	}
	var got *bitwarden.Item
	for i := range found {
		if found[i].Name == spikeNoteName && found[i].Type == bitwarden.ItemTypeSecureNote {
			got = &found[i]
			break
		}
	}
	if got == nil {
		// Fallback: Get by ID from create result.
		item, getErr := vault.Get(ctx, created.ID)
		if getErr != nil {
			return fmt.Errorf("secure note not found by Search and Get failed: %w", getErr)
		}
		got = &item
		fmt.Println("WARN: Search did not return note; used Get by create ID")
	}
	fmt.Printf("OK get id=%s notes_len=%d\n", got.ID, len(got.Notes))

	updatedNotes := fmt.Sprintf("bwsf spike updated at %s", time.Now().UTC().Format(time.RFC3339))
	fmt.Println("Updating notes…")
	updated, err := vault.Update(ctx, got.ID, bitwarden.Item{
		Type:     bitwarden.ItemTypeSecureNote,
		Name:     spikeNoteName,
		Notes:    updatedNotes,
		FolderID: folderID,
	})
	if err != nil {
		return fmt.Errorf("Vault.Update: %w", err)
	}
	if updated.Notes != updatedNotes {
		return fmt.Errorf("Vault.Update: notes mismatch after update")
	}
	fmt.Printf("OK update id=%s\n", updated.ID)
	return nil
}

// ---------------------------------------------------------------------------
// Scenario B — Personal API Key probe (SDK has no API-key login)
// ---------------------------------------------------------------------------

func runScenarioB(ctx context.Context) error {
	clientID := os.Getenv(envClientID)
	clientSecret := os.Getenv(envClientSecret)
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("%s and %s are required", envClientID, envClientSecret)
	}

	identityBase, err := resolveIdentityBase()
	if err != nil {
		return err
	}
	tokenURL := strings.TrimRight(identityBase, "/") + "/connect/token"
	fmt.Printf("SDK gap: bitwarden-go-sdk v0.4.0 has no Personal API Key / client_credentials login.\n")
	fmt.Printf("Probing raw HTTP POST %s\n", tokenURL)

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("scope", "api")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("POST connect/token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("connect/token HTTP %d: %s", resp.StatusCode, truncate(string(body), 400))
	}

	var parsed struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("decode token JSON: %w", err)
	}
	if parsed.AccessToken == "" {
		return fmt.Errorf("access_token missing in response")
	}

	fmt.Printf("OK access_token obtained (type=%s expires_in=%d len=%d)\n",
		parsed.TokenType, parsed.ExpiresIn, len(parsed.AccessToken))
	fmt.Println("SUCCESS (Scenario B): Personal API Key → access_token (vault CRUD not required)")
	return nil
}

func resolveIdentityBase() (string, error) {
	if v := strings.TrimSpace(os.Getenv(envIdentityURL)); v != "" {
		return strings.TrimRight(v, "/"), nil
	}
	if server := strings.TrimSpace(os.Getenv(envServerURL)); server != "" {
		normalized, err := normalizeServerURL(server)
		if err != nil {
			return "", err
		}
		return deriveSelfHostedURLs(normalized).Identity, nil
	}
	return cloudIdentityUS, nil
}

type identityAPIURLs struct {
	Identity string
	API      string
}

// normalizeServerURL mirrors SDK self-hosted rules enough for spike DX:
// absolute https URL, trailing slashes stripped from path.
func normalizeServerURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("parse %s: %w", envServerURL, err)
	}
	if !u.IsAbs() || u.Host == "" {
		return "", fmt.Errorf("%s must be absolute URL with host", envServerURL)
	}
	if u.Scheme != "https" {
		return "", fmt.Errorf("%s must use https (SDK rejects http)", envServerURL)
	}
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawQuery = ""
	u.Fragment = ""
	return u.Scheme + "://" + u.Host + u.Path, nil
}

func deriveSelfHostedURLs(normalizedBase string) identityAPIURLs {
	base := strings.TrimRight(normalizedBase, "/")
	return identityAPIURLs{
		Identity: base + "/identity",
		API:      base + "/api",
	}
}

// diagnoseConnectToken probes /connect/token without the user's password so
// PO can see Vaultwarden's JSON body (SDK often surfaces only "unknown status=400").
func diagnoseConnectToken(ctx context.Context, identityBase string) {
	tokenURL := strings.TrimRight(identityBase, "/") + "/connect/token"
	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("username", "diagnose@example.invalid")
	form.Set("password", "not-a-real-hash")
	form.Set("scope", "api offline_access")
	form.Set("client_id", "desktop")
	form.Set("deviceType", spikeDeviceType)
	form.Set("deviceIdentifier", spikeDeviceIdentifier)
	form.Set("deviceName", spikeDeviceName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "diagnose: build request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "diagnose: POST %s: %v\n", tokenURL, err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
	fmt.Fprintf(os.Stderr, "diagnose: POST %s → HTTP %d\n", tokenURL, resp.StatusCode)
	fmt.Fprintf(os.Stderr, "diagnose: body=%s\n", truncate(string(body), 500))

	// Also show what missing deviceName looks like (common VW 400).
	form.Del("deviceName")
	req2, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return
	}
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req2.Header.Set("Accept", "application/json")
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return
	}
	defer resp2.Body.Close()
	body2, _ := io.ReadAll(io.LimitReader(resp2.Body, 1<<16))
	fmt.Fprintf(os.Stderr, "diagnose (no deviceName): HTTP %d body=%s\n", resp2.StatusCode, truncate(string(body2), 300))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
