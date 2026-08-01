package infra

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"bwsf/src/config"
	"bwsf/src/core"

	"github.com/google/uuid"
)

// ErrAPINotImplemented is returned by API vault methods until Issue #53 Step 4.
var ErrAPINotImplemented = errors.New(
	"API vault operations are not implemented yet (Issue #53 Step 4). " +
		"Authenticate with `bwsf auth` first. " +
		"Use `bwsf backend --set bw` to switch back to the Bitwarden CLI backend",
)

// ErrAPINotAuthenticated means Identity login has not succeeded in this process.
var ErrAPINotAuthenticated = errors.New(
	"API backend is not authenticated. Run `bwsf auth` to store your Personal API Key and obtain a token",
)

// ErrAPIUnlockNotImplemented is returned until Issue #53 Step 3 (master password → crypto keys).
var ErrAPIUnlockNotImplemented = errors.New(
	"API vault unlock (master password / crypto keys) is not implemented yet (Issue #53 Step 3)",
)

// ApiBwClient is the Bitwarden API backend adapter.
// Step 2: Personal API Key + Identity token. Vault CRUD remains stubbed.
type ApiBwClient struct {
	cfg      *config.Config
	store    SecretStore
	identity *IdentityClient

	mu    sync.Mutex
	token *TokenSet
}

// NewApiBwClient creates an ApiBwClient with OS keyring and default HTTP client.
func NewApiBwClient(cfg *config.Config) *ApiBwClient {
	return NewApiBwClientWithDeps(cfg, NewKeyringStore(), NewIdentityClient())
}

// NewApiBwClientWithDeps creates an ApiBwClient with injectable dependencies (tests).
func NewApiBwClientWithDeps(cfg *config.Config, store SecretStore, identity *IdentityClient) *ApiBwClient {
	if cfg == nil {
		cfg = &config.Config{}
	}
	if store == nil {
		store = NewKeyringStore()
	}
	if identity == nil {
		identity = NewIdentityClient()
	}
	return &ApiBwClient{
		cfg:      cfg,
		store:    store,
		identity: identity,
	}
}

// EnsureDeviceIdentifier returns a stable device id, persisting one into config when missing.
func EnsureDeviceIdentifier(cfg *config.Config) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("config is nil")
	}
	if cfg.DeviceIdentifier != "" {
		return cfg.DeviceIdentifier, nil
	}
	id := uuid.NewString()
	cfg.DeviceIdentifier = id
	if err := config.SaveConfig(cfg); err != nil {
		return "", fmt.Errorf("failed to persist device_identifier: %w", err)
	}
	return id, nil
}

// Authenticate loads the Personal API Key from the secret store and obtains an Identity token.
func (c *ApiBwClient) Authenticate(ctx context.Context) error {
	creds, err := LoadAPICredentials(c.store)
	if err != nil {
		if errors.Is(err, ErrSecretNotFound) {
			return ErrAPINotAuthenticated
		}
		return err
	}
	return c.authenticateWithCredentials(ctx, creds)
}

// AuthenticateWithCredentials stores creds (optional persist) and obtains a token.
func (c *ApiBwClient) AuthenticateWithCredentials(ctx context.Context, creds APICredentials, persist bool) error {
	if persist {
		if err := SaveAPICredentials(c.store, creds); err != nil {
			return err
		}
	}
	return c.authenticateWithCredentials(ctx, creds)
}

func (c *ApiBwClient) authenticateWithCredentials(ctx context.Context, creds APICredentials) error {
	if creds.ClientID == "" || creds.ClientSecret == "" {
		return fmt.Errorf("client_id and client_secret are required")
	}

	identityBase, err := ResolveIdentityBase(c.cfg.HostType, c.cfg.SelfhostedURL)
	if err != nil {
		return err
	}

	deviceID, err := EnsureDeviceIdentifier(c.cfg)
	if err != nil {
		return err
	}
	device := DefaultDeviceInfo(deviceID)

	token, err := c.identity.FetchClientCredentialsToken(ctx, identityBase, creds.ClientID, creds.ClientSecret, device)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.token = token
	c.mu.Unlock()
	return nil
}

// EnsureAccessToken returns a valid access token, refreshing or re-authing when needed.
func (c *ApiBwClient) EnsureAccessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	token := c.token
	c.mu.Unlock()

	if token.Valid() {
		return token.AccessToken, nil
	}

	// Prefer refresh_token when present (client_credentials often has none).
	if token != nil && token.RefreshToken != "" {
		identityBase, err := ResolveIdentityBase(c.cfg.HostType, c.cfg.SelfhostedURL)
		if err != nil {
			return "", err
		}
		deviceID, err := EnsureDeviceIdentifier(c.cfg)
		if err != nil {
			return "", err
		}
		creds, _ := LoadAPICredentials(c.store)
		refreshed, err := c.identity.RefreshAccessToken(
			ctx, identityBase, token.RefreshToken, creds.ClientID, DefaultDeviceInfo(deviceID),
		)
		if err == nil {
			c.mu.Lock()
			c.token = refreshed
			c.mu.Unlock()
			return refreshed.AccessToken, nil
		}
	}

	if err := c.Authenticate(ctx); err != nil {
		return "", err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token == nil || c.token.AccessToken == "" {
		return "", ErrAPINotAuthenticated
	}
	return c.token.AccessToken, nil
}

// IsAuthenticated reports whether this process holds a non-expired access token.
func (c *ApiBwClient) IsAuthenticated() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.token.Valid()
}

// ClearSession drops the in-memory token (does not remove Keychain secrets).
func (c *ApiBwClient) ClearSession() {
	c.mu.Lock()
	c.token = nil
	c.mu.Unlock()
}

// TokenExpiresAt exposes token expiry for diagnostics/tests.
func (c *ApiBwClient) TokenExpiresAt() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token == nil {
		return time.Time{}
	}
	return c.token.ExpiresAt
}

func (c *ApiBwClient) GetDotenvsFolderID() (string, error) {
	return "", ErrAPINotImplemented
}

func (c *ApiBwClient) DotenvsFolderExists() (bool, error) {
	return false, ErrAPINotImplemented
}

func (c *ApiBwClient) CreateDotenvsFolder() error {
	return ErrAPINotImplemented
}

func (c *ApiBwClient) ListItemsInFolder(folderID string) ([]core.Item, error) {
	return nil, ErrAPINotImplemented
}

func (c *ApiBwClient) GetItemByName(folderID, name string) (*core.FullItem, error) {
	return nil, ErrAPINotImplemented
}

func (c *ApiBwClient) GetItemByID(id string) (*core.FullItem, error) {
	return nil, ErrAPINotImplemented
}

func (c *ApiBwClient) CreateNoteItem(folderID, name, notes string) error {
	return ErrAPINotImplemented
}

func (c *ApiBwClient) UpdateNoteItem(id, notes string) error {
	return ErrAPINotImplemented
}

// Login for the API backend authenticates with a stored Personal API Key.
// email/password are ignored (CLI-compatible signature); use `bwsf auth` to store the key.
func (c *ApiBwClient) Login(email, password, serverURL string) error {
	_ = email
	_ = password
	if serverURL != "" && c.cfg.SelfhostedURL == "" {
		c.cfg.SelfhostedURL = serverURL
		if c.cfg.HostType == "" {
			c.cfg.HostType = "selfhosted"
		}
	}
	return c.Authenticate(context.Background())
}

// Unlock is not available until Step 3 (master password → decryption keys).
func (c *ApiBwClient) Unlock(masterPassword string) error {
	_ = masterPassword
	if !c.IsAuthenticated() {
		if err := c.Authenticate(context.Background()); err != nil {
			return err
		}
	}
	return ErrAPIUnlockNotImplemented
}
