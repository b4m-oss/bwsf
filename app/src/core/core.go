package core

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"bwenv/src/config"
)

// BwClient は Bitwarden CLI とのやり取りを抽象化するインターフェースです。
type BwClient interface {
	GetDotenvsFolderID() (string, error)
	ListItemsInFolder(folderID string) ([]Item, error)
	GetItemByName(folderID, name string) (*FullItem, error)
	GetItemByID(id string) (*FullItem, error)
	CreateNoteItem(folderID, name, notes string) error
	UpdateNoteItem(id, notes string) error
	Login(email, password, serverURL string) error
	Unlock(masterPassword string) error
}

// FileSystem はファイルシステム操作を抽象化するインターフェースです。
type FileSystem interface {
	OpenEnvFile(path string) ([]byte, error)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm uint32) error
	Stat(path string) (FileInfo, error)
	MkdirAll(path string, perm uint32) error
}

// FileInfo は Stat の結果に必要な最小限の情報を表します。
type FileInfo interface {
	IsNotExist() bool
}

// Logger はログ出力を抽象化するインターフェースです。
type Logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
}

// Item は dotenvs フォルダ内に保存される Bitwarden アイテムを表します。
type Item struct {
	ID   string
	Name string
}

// FullItem は Bitwarden アイテムの完全な情報を表します。
type FullItem struct {
	ID    string
	Name  string
	Notes string
}

// EnvData は .env ファイルのデータを表します。
type EnvData struct {
	Lines []string `json:"lines"`
}

// IsLockedError はエラーがロック関連かどうかを判定します。
func IsLockedError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "Bitwarden CLI is locked") ||
		strings.Contains(errMsg, "Master password") ||
		strings.Contains(errMsg, "master password")
}

// WithUnlockRetry は Bitwarden がロックされている場合に Unlock/Login を挟んでリトライする共通処理です。
func WithUnlockRetry(
	bw BwClient,
	cfg *config.Config,
	promptPassword func() (string, error),
	logger Logger,
	fn func() error,
) error {
	// 最初に fn() を実行
	err := fn()
	if err == nil {
		return nil
	}

	// ロック関連エラーでなければそのまま返す
	if !IsLockedError(err) {
		return err
	}

	// ロック状態なのでアンロックを試みる
	logger.Info("Bitwarden CLI is locked. Please enter your master password to unlock.")

	// パスワード入力
	password, promptErr := promptPassword()
	if promptErr != nil {
		return fmt.Errorf("failed to get master password: %w", promptErr)
	}

	// Unlock を試行
	unlockErr := bw.Unlock(password)
	if unlockErr == nil {
		// Unlock 成功、fn() を再実行
		logger.Info("Bitwarden CLI unlocked successfully")
		return fn()
	}

	// Unlock 失敗、cfg があれば Login を試みる
	if cfg != nil && cfg.Email != "" {
		logger.Info("Unlock failed, trying login then unlock...")
		loginErr := bw.Login(cfg.Email, password, cfg.SelfhostedURL)
		if loginErr != nil {
			return fmt.Errorf("failed to login Bitwarden CLI: %w", loginErr)
		}
		logger.Info("Bitwarden CLI logged in successfully")

		// Login 成功後、再度 Unlock を試行
		unlockErr = bw.Unlock(password)
		if unlockErr != nil {
			return fmt.Errorf("failed to unlock Bitwarden CLI after login: %w", unlockErr)
		}
		logger.Info("Bitwarden CLI unlocked successfully")

		// fn() を再実行
		return fn()
	}

	return fmt.Errorf("failed to unlock Bitwarden CLI: %w", unlockErr)
}

// PushEnvCore は .env ファイルを Bitwarden にプッシュするコアロジックです。
func PushEnvCore(
	fromDir, projectName string,
	fs FileSystem,
	bw BwClient,
	cfg *config.Config,
	promptPassword func() (string, error),
	logger Logger,
) error {
	// .env ファイルのパスを決定
	envPath := filepath.Join(fromDir, ".env")

	// .env ファイルを読み込む
	content, err := fs.OpenEnvFile(envPath)
	if err != nil {
		// fromDir が "." または ".." の場合、/project-root へフォールバック
		if fromDir == "." || fromDir == ".." {
			fallbackPath := filepath.Join("/project-root", ".env")
			content, err = fs.OpenEnvFile(fallbackPath)
			if err != nil {
				return fmt.Errorf("failed to open .env file: %w", err)
			}
		} else {
			return fmt.Errorf("failed to open .env file: %w", err)
		}
	}

	// .env をパースして JSON に変換
	envData := parseEnvContent(content)
	jsonData, err := envDataToJSON(envData)
	if err != nil {
		return fmt.Errorf("failed to convert to JSON: %w", err)
	}

	// dotenvs フォルダ ID を取得
	var folderID string
	err = WithUnlockRetry(bw, cfg, promptPassword, logger, func() error {
		var innerErr error
		folderID, innerErr = bw.GetDotenvsFolderID()
		return innerErr
	})
	if err != nil {
		return fmt.Errorf("failed to get dotenvs folder: %w", err)
	}

	// 既存アイテムを検索
	var existingItem *FullItem
	err = WithUnlockRetry(bw, cfg, promptPassword, logger, func() error {
		var innerErr error
		existingItem, innerErr = bw.GetItemByName(folderID, projectName)
		return innerErr
	})
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	// 既存アイテムがあれば更新、なければ新規作成
	if existingItem != nil {
		err = WithUnlockRetry(bw, cfg, promptPassword, logger, func() error {
			return bw.UpdateNoteItem(existingItem.ID, jsonData)
		})
		if err != nil {
			return fmt.Errorf("failed to update item: %w", err)
		}
	} else {
		err = WithUnlockRetry(bw, cfg, promptPassword, logger, func() error {
			return bw.CreateNoteItem(folderID, projectName, jsonData)
		})
		if err != nil {
			return fmt.Errorf("failed to create item: %w", err)
		}
	}

	return nil
}

// PullEnvCore は Bitwarden から .env ファイルをプルするコアロジックです。
func PullEnvCore(
	outputDir, projectName string,
	fs FileSystem,
	bw BwClient,
	cfg *config.Config,
	promptPassword func() (string, error),
	confirmOverwrite func(path string) (bool, error),
	logger Logger,
) error {
	// outputDir を正規化
	if outputDir == "." || outputDir == ".." {
		outputDir = "/project-root"
	}

	// dotenvs フォルダ ID を取得
	var folderID string
	err := WithUnlockRetry(bw, cfg, promptPassword, logger, func() error {
		var innerErr error
		folderID, innerErr = bw.GetDotenvsFolderID()
		return innerErr
	})
	if err != nil {
		return fmt.Errorf("failed to get dotenvs folder: %w", err)
	}

	// アイテムを取得
	var item *FullItem
	err = WithUnlockRetry(bw, cfg, promptPassword, logger, func() error {
		var innerErr error
		item, innerErr = bw.GetItemByName(folderID, projectName)
		return innerErr
	})
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	// アイテムが見つからない場合
	if item == nil {
		return fmt.Errorf("item '%s' not found in dotenvs folder", projectName)
	}

	// JSON から .env 内容を復元
	envContent, err := restoreEnvFileFromJSON(item.Notes)
	if err != nil {
		return fmt.Errorf("failed to restore .env from JSON: %w", err)
	}

	// 出力先パス
	envPath := filepath.Join(outputDir, ".env")

	// ファイルの存在確認
	info, err := fs.Stat(envPath)
	if err == nil && !info.IsNotExist() {
		// ファイルが存在する場合、上書き確認
		confirmed, confirmErr := confirmOverwrite(envPath)
		if confirmErr != nil {
			return fmt.Errorf("failed to confirm overwrite: %w", confirmErr)
		}
		if !confirmed {
			return nil // キャンセル
		}
	}

	// ディレクトリを作成（必要に応じて）
	if outputDir != "." && outputDir != "/project-root" {
		if err := fs.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// .env ファイルを書き出し
	if err := fs.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	return nil
}

// ListDotenvsCore は dotenvs フォルダ内のアイテム一覧を取得するコアロジックです。
func ListDotenvsCore(
	bw BwClient,
	cfg *config.Config,
	promptPassword func() (string, error),
	logger Logger,
) ([]Item, error) {
	// dotenvs フォルダ ID を取得
	var folderID string
	err := WithUnlockRetry(bw, cfg, promptPassword, logger, func() error {
		var innerErr error
		folderID, innerErr = bw.GetDotenvsFolderID()
		return innerErr
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get dotenvs folder: %w", err)
	}

	// アイテム一覧を取得
	var items []Item
	err = WithUnlockRetry(bw, cfg, promptPassword, logger, func() error {
		var innerErr error
		items, innerErr = bw.ListItemsInFolder(folderID)
		return innerErr
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	return items, nil
}

// SetupBitwardenCore は Bitwarden のセットアップを行うコアロジックです。
func SetupBitwardenCore(
	fs FileSystem,
	bw BwClient,
	logger Logger,
	selectHostType func() (string, error),
	inputURL func() (string, error),
	inputEmail func() (string, error),
	inputPassword func() (string, error),
) error {
	// 既存設定を読み込み
	existingConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load existing config: %w", err)
	}
	if existingConfig != nil {
		logger.Info("Existing configuration found. It will be overwritten.")
	}

	// ホストタイプを選択
	hostType, err := selectHostType()
	if err != nil {
		return fmt.Errorf("failed to select host type: %w", err)
	}

	// Self-hosted の場合は URL を入力
	var selfhostedURL string
	if hostType == "selfhosted" {
		selfhostedURL, err = inputURL()
		if err != nil {
			return fmt.Errorf("failed to get URL: %w", err)
		}
	}

	// メールアドレスを入力
	email, err := inputEmail()
	if err != nil {
		return fmt.Errorf("failed to get email: %w", err)
	}

	// パスワードを入力
	password, err := inputPassword()
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	// ログイン
	if err := bw.Login(email, password, selfhostedURL); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	// 設定を保存
	newConfig := &config.Config{
		HostType:      hostType,
		SelfhostedURL: selfhostedURL,
		Email:         email,
	}
	if err := config.SaveConfig(newConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	return nil
}

// parseEnvContent は .env ファイルの内容をパースします。
func parseEnvContent(content []byte) *EnvData {
	lines := strings.Split(string(content), "\n")
	// 末尾の空行を削除
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return &EnvData{Lines: lines}
}

// envDataToJSON は EnvData を JSON 文字列に変換します。
func envDataToJSON(data *EnvData) (string, error) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal env data to JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// restoreEnvFileFromJSON は JSON から .env ファイルの内容を復元します。
func restoreEnvFileFromJSON(jsonStr string) (string, error) {
	var data EnvData
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return strings.Join(data.Lines, "\n"), nil
}
