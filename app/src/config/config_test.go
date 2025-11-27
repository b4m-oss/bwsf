package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// GetConfigPath のテスト
// =============================================================================

// 正常系: ホームディレクトリ配下の .config/bwenv/config.json パスが返る
func TestGetConfigPath_Success(t *testing.T) {
	path, err := GetConfigPath()

	assert.NoError(t, err)
	assert.Contains(t, path, ".config/bwenv/config.json")

	// ホームディレクトリで始まることを確認
	homeDir, _ := os.UserHomeDir()
	assert.True(t, filepath.HasPrefix(path, homeDir))
}

// =============================================================================
// LoadConfig のテスト
// =============================================================================

// 正常系: 設定ファイルが存在しない場合は nil, nil を返す
func TestLoadConfig_FileNotExist(t *testing.T) {
	// 一時的に HOME を変更してテスト用ディレクトリを使う
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.Nil(t, cfg)
}

// 正常系: 設定ファイルが存在する場合は正しくパースされる
func TestLoadConfig_Success(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// 設定ファイルを作成
	configDir := filepath.Join(tmpDir, ".config", "bwenv")
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config.json")
	configContent := `{"host_type":"selfhosted","selfhosted_url":"https://bw.example.com","email":"test@example.com"}`
	os.WriteFile(configPath, []byte(configContent), 0600)

	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "selfhosted", cfg.HostType)
	assert.Equal(t, "https://bw.example.com", cfg.SelfhostedURL)
	assert.Equal(t, "test@example.com", cfg.Email)
}

// 異常系: JSON が壊れている場合はエラー
func TestLoadConfig_InvalidJSON(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// 壊れた JSON を書き込む
	configDir := filepath.Join(tmpDir, ".config", "bwenv")
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config.json")
	os.WriteFile(configPath, []byte("not valid json"), 0600)

	cfg, err := LoadConfig()

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to parse config file")
}

// =============================================================================
// SaveConfig のテスト
// =============================================================================

// 正常系: 設定ファイルが正しく保存される
func TestSaveConfig_Success(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg := &Config{
		HostType:      "cloud",
		SelfhostedURL: "",
		Email:         "user@example.com",
	}

	err := SaveConfig(cfg)

	assert.NoError(t, err)

	// ファイルが作成されていることを確認
	configPath := filepath.Join(tmpDir, ".config", "bwenv", "config.json")
	_, statErr := os.Stat(configPath)
	assert.NoError(t, statErr)

	// 内容を確認
	content, _ := os.ReadFile(configPath)
	assert.Contains(t, string(content), `"host_type": "cloud"`)
	assert.Contains(t, string(content), `"email": "user@example.com"`)
}

// 正常系: ディレクトリが存在しない場合は作成される
func TestSaveConfig_CreatesDirectory(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg := &Config{
		HostType: "cloud",
		Email:    "test@example.com",
	}

	err := SaveConfig(cfg)

	assert.NoError(t, err)

	// ディレクトリが作成されていることを確認
	configDir := filepath.Join(tmpDir, ".config", "bwenv")
	info, statErr := os.Stat(configDir)
	assert.NoError(t, statErr)
	assert.True(t, info.IsDir())
}

// 正常系: 既存設定を上書き保存
func TestSaveConfig_Overwrite(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// 最初の設定
	cfg1 := &Config{
		HostType: "cloud",
		Email:    "old@example.com",
	}
	SaveConfig(cfg1)

	// 上書き
	cfg2 := &Config{
		HostType:      "selfhosted",
		SelfhostedURL: "https://new.example.com",
		Email:         "new@example.com",
	}
	err := SaveConfig(cfg2)

	assert.NoError(t, err)

	// 新しい内容になっていることを確認
	loaded, _ := LoadConfig()
	assert.Equal(t, "selfhosted", loaded.HostType)
	assert.Equal(t, "https://new.example.com", loaded.SelfhostedURL)
	assert.Equal(t, "new@example.com", loaded.Email)
}

