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

// =============================================================================
// IsExampleFile のテスト
// =============================================================================

// 正常系: .env.example は example ファイル
func TestIsExampleFile_DotEnvExample(t *testing.T) {
	assert.True(t, IsExampleFile(".env.example"))
}

// 正常系: .env.staging.example は example ファイル
func TestIsExampleFile_DotEnvStagingExample(t *testing.T) {
	assert.True(t, IsExampleFile(".env.staging.example"))
}

// 正常系: .env.example.staging は example ファイル
func TestIsExampleFile_DotEnvExampleStaging(t *testing.T) {
	assert.True(t, IsExampleFile(".env.example.staging"))
}

// 正常系: .env は example ファイルではない
func TestIsExampleFile_DotEnv(t *testing.T) {
	assert.False(t, IsExampleFile(".env"))
}

// 正常系: .env.staging は example ファイルではない
func TestIsExampleFile_DotEnvStaging(t *testing.T) {
	assert.False(t, IsExampleFile(".env.staging"))
}

// 正常系: .env.local は example ファイルではない
func TestIsExampleFile_DotEnvLocal(t *testing.T) {
	assert.False(t, IsExampleFile(".env.local"))
}

// =============================================================================
// FindEnvFiles のテスト
// =============================================================================

// 正常系: 複数の .env* ファイルを検出
func TestFindEnvFiles_MultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// テスト用ファイルを作成
	os.WriteFile(filepath.Join(tmpDir, ".env"), []byte("KEY=value"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".env.local"), []byte("KEY=local"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".env.staging"), []byte("KEY=staging"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".env.production"), []byte("KEY=prod"), 0644)

	files, err := FindEnvFiles(tmpDir)

	assert.NoError(t, err)
	assert.Len(t, files, 4)

	// .env が最初に来ることを確認
	assert.Equal(t, ".env", filepath.Base(files[0]))
}

// 正常系: .example ファイルは除外される
func TestFindEnvFiles_ExcludesExampleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// テスト用ファイルを作成
	os.WriteFile(filepath.Join(tmpDir, ".env"), []byte("KEY=value"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".env.example"), []byte("KEY=example"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".env.staging"), []byte("KEY=staging"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".env.staging.example"), []byte("KEY=staging-example"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".env.example.staging"), []byte("KEY=example-staging"), 0644)

	files, err := FindEnvFiles(tmpDir)

	assert.NoError(t, err)
	assert.Len(t, files, 2) // .env と .env.staging のみ

	// ファイル名を確認
	var names []string
	for _, f := range files {
		names = append(names, filepath.Base(f))
	}
	assert.Contains(t, names, ".env")
	assert.Contains(t, names, ".env.staging")
	assert.NotContains(t, names, ".env.example")
	assert.NotContains(t, names, ".env.staging.example")
	assert.NotContains(t, names, ".env.example.staging")
}

// 正常系: .env ファイルがない場合は空スライス
func TestFindEnvFiles_NoEnvFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// .env 以外のファイルを作成
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# README"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "config.json"), []byte("{}"), 0644)

	files, err := FindEnvFiles(tmpDir)

	assert.NoError(t, err)
	assert.Empty(t, files)
}

// 正常系: .env のみの場合
func TestFindEnvFiles_OnlyDotEnv(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, ".env"), []byte("KEY=value"), 0644)

	files, err := FindEnvFiles(tmpDir)

	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, ".env", filepath.Base(files[0]))
}

// 正常系: ファイルはアルファベット順にソートされる（.env が先頭）
func TestFindEnvFiles_SortedAlphabetically(t *testing.T) {
	tmpDir := t.TempDir()

	// 順序がバラバラにファイルを作成
	os.WriteFile(filepath.Join(tmpDir, ".env.staging"), []byte("KEY=staging"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".env.production"), []byte("KEY=prod"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".env"), []byte("KEY=value"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".env.local"), []byte("KEY=local"), 0644)

	files, err := FindEnvFiles(tmpDir)

	assert.NoError(t, err)
	assert.Len(t, files, 4)

	// ソート順を確認
	assert.Equal(t, ".env", filepath.Base(files[0]))
	assert.Equal(t, ".env.local", filepath.Base(files[1]))
	assert.Equal(t, ".env.production", filepath.Base(files[2]))
	assert.Equal(t, ".env.staging", filepath.Base(files[3]))
}

// 異常系: 存在しないディレクトリ
func TestFindEnvFiles_NonExistentDir(t *testing.T) {
	files, err := FindEnvFiles("/nonexistent/directory")

	assert.Error(t, err)
	assert.Nil(t, files)
	assert.Contains(t, err.Error(), "failed to read directory")
}

// 正常系: ディレクトリは無視される
func TestFindEnvFiles_IgnoresDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// .env ファイルと .env ディレクトリを作成
	os.WriteFile(filepath.Join(tmpDir, ".env"), []byte("KEY=value"), 0644)
	os.MkdirAll(filepath.Join(tmpDir, ".env.d"), 0755)

	files, err := FindEnvFiles(tmpDir)

	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, ".env", filepath.Base(files[0]))
}

