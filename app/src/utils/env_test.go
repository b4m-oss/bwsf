package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// ParseEnvFile のテスト
// =============================================================================

// 正常系: .env ファイルを正しくパースできる
func TestParseEnvFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := `# Comment line
KEY1=value1
KEY2="quoted value"

KEY3=value3`
	os.WriteFile(envPath, []byte(content), 0644)

	data, err := ParseEnvFile(envPath)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data.Lines, 5)
	assert.Equal(t, "# Comment line", data.Lines[0])
	assert.Equal(t, "KEY1=value1", data.Lines[1])
	assert.Equal(t, `KEY2="quoted value"`, data.Lines[2])
	assert.Equal(t, "", data.Lines[3]) // 空行
	assert.Equal(t, "KEY3=value3", data.Lines[4])
}

// 正常系: 空の .env ファイル
func TestParseEnvFile_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	os.WriteFile(envPath, []byte(""), 0644)

	data, err := ParseEnvFile(envPath)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Empty(t, data.Lines)
}

// 異常系: ファイルが存在しない
func TestParseEnvFile_FileNotFound(t *testing.T) {
	data, err := ParseEnvFile("/nonexistent/.env")

	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "failed to open .env file")
}

// =============================================================================
// EnvDataToJSON のテスト
// =============================================================================

// 正常系: EnvData を JSON に変換
func TestEnvDataToJSON_Success(t *testing.T) {
	data := &EnvData{
		Lines: []string{"KEY1=value1", "KEY2=value2"},
	}

	jsonStr, err := EnvDataToJSON(data)

	assert.NoError(t, err)
	assert.Contains(t, jsonStr, `"lines"`)
	assert.Contains(t, jsonStr, `"KEY1=value1"`)
	assert.Contains(t, jsonStr, `"KEY2=value2"`)
}

// 正常系: 空の EnvData
func TestEnvDataToJSON_Empty(t *testing.T) {
	data := &EnvData{
		Lines: []string{},
	}

	jsonStr, err := EnvDataToJSON(data)

	assert.NoError(t, err)
	assert.Contains(t, jsonStr, `"lines": []`)
}

// =============================================================================
// RestoreEnvFileFromJSON のテスト
// =============================================================================

// 正常系: JSON から .env 内容を復元
func TestRestoreEnvFileFromJSON_Success(t *testing.T) {
	jsonStr := `{"lines":["KEY1=value1","KEY2=value2","# comment"]}`

	content, err := RestoreEnvFileFromJSON(jsonStr)

	assert.NoError(t, err)
	assert.Equal(t, "KEY1=value1\nKEY2=value2\n# comment", content)
}

// 正常系: 空の lines
func TestRestoreEnvFileFromJSON_EmptyLines(t *testing.T) {
	jsonStr := `{"lines":[]}`

	content, err := RestoreEnvFileFromJSON(jsonStr)

	assert.NoError(t, err)
	assert.Equal(t, "", content)
}

// 異常系: 壊れた JSON
func TestRestoreEnvFileFromJSON_InvalidJSON(t *testing.T) {
	jsonStr := `not valid json`

	content, err := RestoreEnvFileFromJSON(jsonStr)

	assert.Error(t, err)
	assert.Empty(t, content)
	assert.Contains(t, err.Error(), "failed to unmarshal JSON")
}

// =============================================================================
// FindEnvFile のテスト
// =============================================================================

// 正常系: .env ファイルが存在する
func TestFindEnvFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	os.WriteFile(envPath, []byte("KEY=value"), 0644)

	foundPath, err := FindEnvFile(tmpDir)

	assert.NoError(t, err)
	assert.Equal(t, envPath, foundPath)
}

// 異常系: .env ファイルが存在しない
func TestFindEnvFile_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	foundPath, err := FindEnvFile(tmpDir)

	assert.Error(t, err)
	assert.Empty(t, foundPath)
	assert.Contains(t, err.Error(), ".env file not found")
}

// =============================================================================
// 往復テスト（ParseEnvFile -> EnvDataToJSON -> RestoreEnvFileFromJSON）
// =============================================================================

// 正常系: パース → JSON 変換 → 復元で元の内容が維持される
func TestEnvFile_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	originalContent := `# Database settings
DB_HOST=localhost
DB_PORT=5432
DB_PASSWORD="secret password"

# API settings
API_KEY=abc123`
	os.WriteFile(envPath, []byte(originalContent), 0644)

	// パース
	data, err := ParseEnvFile(envPath)
	assert.NoError(t, err)

	// JSON 変換
	jsonStr, err := EnvDataToJSON(data)
	assert.NoError(t, err)

	// 復元
	restoredContent, err := RestoreEnvFileFromJSON(jsonStr)
	assert.NoError(t, err)

	// 元の内容と一致することを確認
	assert.Equal(t, originalContent, restoredContent)
}

