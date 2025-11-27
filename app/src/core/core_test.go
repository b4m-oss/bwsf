package core

import (
	"errors"
	"fmt"
	"testing"

	"bwenv/src/config"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// テスト用モック実装
// =============================================================================

// ErrBitwardenLocked はテスト用のロックエラー
var ErrBitwardenLocked = errors.New("Bitwarden CLI is locked")

// --- mockBwClient ---

type mockBwClient struct {
	calls []string

	// GetDotenvsFolderID の挙動制御
	folderID    string
	folderIDErr error

	// ListItemsInFolder の挙動制御
	items        []Item
	listItemsErr error

	// GetItemByName の挙動制御
	itemByName    *FullItem
	itemByNameErr error

	// GetItemByID の挙動制御
	itemByID    *FullItem
	itemByIDErr error

	// CreateNoteItem の挙動制御
	createErr error

	// UpdateNoteItem の挙動制御
	updateErr error

	// Login の挙動制御
	loginErr error

	// Unlock の挙動制御
	unlockErr error
}

func (m *mockBwClient) GetDotenvsFolderID() (string, error) {
	m.calls = append(m.calls, "GetDotenvsFolderID")
	if m.folderIDErr != nil {
		return "", m.folderIDErr
	}
	if m.folderID == "" {
		return "dummy-folder-id", nil
	}
	return m.folderID, nil
}

func (m *mockBwClient) ListItemsInFolder(folderID string) ([]Item, error) {
	m.calls = append(m.calls, fmt.Sprintf("ListItemsInFolder(%s)", folderID))
	if m.listItemsErr != nil {
		return nil, m.listItemsErr
	}
	return m.items, nil
}

func (m *mockBwClient) GetItemByName(folderID, name string) (*FullItem, error) {
	m.calls = append(m.calls, fmt.Sprintf("GetItemByName(%s,%s)", folderID, name))
	if m.itemByNameErr != nil {
		return nil, m.itemByNameErr
	}
	return m.itemByName, nil
}

func (m *mockBwClient) GetItemByID(id string) (*FullItem, error) {
	m.calls = append(m.calls, fmt.Sprintf("GetItemByID(%s)", id))
	if m.itemByIDErr != nil {
		return nil, m.itemByIDErr
	}
	return m.itemByID, nil
}

func (m *mockBwClient) CreateNoteItem(folderID, name, notes string) error {
	m.calls = append(m.calls, fmt.Sprintf("CreateNoteItem(%s,%s)", folderID, name))
	return m.createErr
}

func (m *mockBwClient) UpdateNoteItem(id, notes string) error {
	m.calls = append(m.calls, fmt.Sprintf("UpdateNoteItem(%s)", id))
	return m.updateErr
}

func (m *mockBwClient) Login(email, password, serverURL string) error {
	m.calls = append(m.calls, fmt.Sprintf("Login(%s,%s)", email, serverURL))
	return m.loginErr
}

func (m *mockBwClient) Unlock(masterPassword string) error {
	m.calls = append(m.calls, "Unlock")
	return m.unlockErr
}

// --- mockFileSystem ---

type mockFileInfo struct {
	notExist bool
}

func (f *mockFileInfo) IsNotExist() bool {
	return f.notExist
}

type mockFileSystem struct {
	calls []string

	// OpenEnvFile の挙動制御
	envFileContent []byte
	envFileErr     error

	// ReadFile の挙動制御
	readContent []byte
	readErr     error

	// WriteFile の挙動制御
	writtenPath string
	writtenData []byte
	writeErr    error

	// Stat の挙動制御
	statInfo FileInfo
	statErr  error

	// MkdirAll の挙動制御
	mkdirErr error
}

func (m *mockFileSystem) OpenEnvFile(path string) ([]byte, error) {
	m.calls = append(m.calls, fmt.Sprintf("OpenEnvFile(%s)", path))
	if m.envFileErr != nil {
		return nil, m.envFileErr
	}
	return m.envFileContent, nil
}

func (m *mockFileSystem) ReadFile(path string) ([]byte, error) {
	m.calls = append(m.calls, fmt.Sprintf("ReadFile(%s)", path))
	if m.readErr != nil {
		return nil, m.readErr
	}
	return m.readContent, nil
}

func (m *mockFileSystem) WriteFile(path string, data []byte, perm uint32) error {
	m.calls = append(m.calls, fmt.Sprintf("WriteFile(%s)", path))
	m.writtenPath = path
	m.writtenData = data
	return m.writeErr
}

func (m *mockFileSystem) Stat(path string) (FileInfo, error) {
	m.calls = append(m.calls, fmt.Sprintf("Stat(%s)", path))
	if m.statErr != nil {
		return nil, m.statErr
	}
	if m.statInfo != nil {
		return m.statInfo, nil
	}
	return &mockFileInfo{notExist: true}, nil
}

func (m *mockFileSystem) MkdirAll(path string, perm uint32) error {
	m.calls = append(m.calls, fmt.Sprintf("MkdirAll(%s)", path))
	return m.mkdirErr
}

// --- mockLogger ---

type mockLogger struct {
	errors   []string
	infos    []string
	warnings []string
	success  []string
}

func (l *mockLogger) Error(args ...interface{}) {
	l.errors = append(l.errors, fmt.Sprint(args...))
}

func (l *mockLogger) Info(args ...interface{}) {
	l.infos = append(l.infos, fmt.Sprint(args...))
}

func (l *mockLogger) Warning(args ...interface{}) {
	l.warnings = append(l.warnings, fmt.Sprint(args...))
}

func (l *mockLogger) Success(args ...interface{}) {
	l.success = append(l.success, fmt.Sprint(args...))
}

// =============================================================================
// WithUnlockRetry のテスト
// =============================================================================

// 正常系: fn が一度で成功する場合は Unlock/Login を行わず、そのまま成功すること
func TestWithUnlockRetry_SuccessWithoutRetry(t *testing.T) {
	bw := &mockBwClient{}
	logger := &mockLogger{}
	cfg := &config.Config{}

	called := 0
	fn := func() error {
		called++
		return nil
	}

	err := WithUnlockRetry(
		bw,
		cfg,
		func() (string, error) {
			return "", errors.New("prompt should not be called")
		},
		logger,
		fn,
	)

	assert.NoError(t, err)
	assert.Equal(t, 1, called, "fn should be called exactly once")
	assert.NotContains(t, bw.calls, "Unlock")
	assert.NotContains(t, bw.calls, "Login")
}

// 正常系: fn が 1 回目で ErrBitwardenLocked を返し、Unlock 成功後の 2 回目で成功するケース
func TestWithUnlockRetry_LockThenUnlockSuccess(t *testing.T) {
	bw := &mockBwClient{}
	logger := &mockLogger{}
	cfg := &config.Config{Email: "test@example.com"}

	callCount := 0
	fn := func() error {
		callCount++
		if callCount == 1 {
			return ErrBitwardenLocked
		}
		return nil
	}

	promptCalled := false
	err := WithUnlockRetry(
		bw,
		cfg,
		func() (string, error) {
			promptCalled = true
			return "password123", nil
		},
		logger,
		fn,
	)

	assert.NoError(t, err)
	assert.Equal(t, 2, callCount, "fn should be called twice")
	assert.True(t, promptCalled, "promptPassword should be called")
	assert.Contains(t, bw.calls, "Unlock")
}

// 異常系: fn がロック関連以外のエラーを返した場合は、Unlock/Login を試みずにエラーをそのまま返すこと
func TestWithUnlockRetry_NonLockErrorPropagates(t *testing.T) {
	bw := &mockBwClient{}
	logger := &mockLogger{}
	cfg := &config.Config{}

	expectedErr := errors.New("some error")
	fn := func() error {
		return expectedErr
	}

	err := WithUnlockRetry(
		bw,
		cfg,
		func() (string, error) {
			return "pwd", nil
		},
		logger,
		fn,
	)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.NotContains(t, bw.calls, "Unlock")
	assert.NotContains(t, bw.calls, "Login")
}

// 異常系: promptPassword がエラーを返した場合、Unlock/Login が呼ばれずにそのエラーが返る
func TestWithUnlockRetry_PromptPasswordError(t *testing.T) {
	bw := &mockBwClient{}
	logger := &mockLogger{}
	cfg := &config.Config{}

	callCount := 0
	fn := func() error {
		callCount++
		if callCount == 1 {
			return ErrBitwardenLocked
		}
		return nil
	}

	promptErr := errors.New("prompt cancelled")
	err := WithUnlockRetry(
		bw,
		cfg,
		func() (string, error) {
			return "", promptErr
		},
		logger,
		fn,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "prompt cancelled")
	assert.Equal(t, 1, callCount, "fn should be called only once before prompt error")
	assert.NotContains(t, bw.calls, "Unlock")
}

// 異常系: Unlock と Login の両方が失敗する場合
func TestWithUnlockRetry_UnlockAndLoginBothFail(t *testing.T) {
	bw := &mockBwClient{
		unlockErr: errors.New("unlock failed"),
		loginErr:  errors.New("login failed"),
	}
	logger := &mockLogger{}
	cfg := &config.Config{Email: "test@example.com", SelfhostedURL: "https://bw.example.com"}

	callCount := 0
	fn := func() error {
		callCount++
		return ErrBitwardenLocked
	}

	err := WithUnlockRetry(
		bw,
		cfg,
		func() (string, error) {
			return "password", nil
		},
		logger,
		fn,
	)

	assert.Error(t, err)
	assert.Equal(t, 1, callCount, "fn should not be retried after unlock/login failure")
	assert.Contains(t, bw.calls, "Unlock")
	assert.Contains(t, bw.calls, "Login(test@example.com,https://bw.example.com)")
}

// =============================================================================
// PushEnvCore のテスト
// =============================================================================

// 正常系: .env が存在し、既存アイテムがないケースで CreateNoteItem が呼ばれる
func TestPushEnvCore_CreateNewItem(t *testing.T) {
	bw := &mockBwClient{
		folderID:   "folder-123",
		itemByName: nil, // アイテムなし
	}
	fs := &mockFileSystem{
		envFileContent: []byte("KEY=value\n"),
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	err := PushEnvCore(
		".",
		"my-project",
		fs,
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		logger,
	)

	// Red フェーズ: 現在は not implemented エラーが返る
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

// 正常系: .env が存在し、既存アイテムがあるケースで UpdateNoteItem が呼ばれる
func TestPushEnvCore_UpdateExistingItem(t *testing.T) {
	bw := &mockBwClient{
		folderID:   "folder-123",
		itemByName: &FullItem{ID: "item-456", Name: "my-project", Notes: "{}"},
	}
	fs := &mockFileSystem{
		envFileContent: []byte("KEY=newvalue\n"),
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	err := PushEnvCore(
		".",
		"my-project",
		fs,
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

// 正常系: fromDir が "." で .env がなく、/project-root/.env へフォールバック
func TestPushEnvCore_FallbackToProjectRoot(t *testing.T) {
	bw := &mockBwClient{folderID: "folder-123"}
	fs := &mockFileSystem{
		envFileContent: []byte("FALLBACK=true\n"),
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	err := PushEnvCore(
		".",
		"my-project",
		fs,
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

// 異常系: .env ファイルが見つからない
func TestPushEnvCore_EnvFileNotFound(t *testing.T) {
	bw := &mockBwClient{folderID: "folder-123"}
	fs := &mockFileSystem{
		envFileErr: errors.New(".env file not found"),
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	err := PushEnvCore(
		".",
		"my-project",
		fs,
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
}

// 異常系: GetDotenvsFolderID がエラーを返す
func TestPushEnvCore_GetFolderIDError(t *testing.T) {
	bw := &mockBwClient{
		folderIDErr: errors.New("dotenvs folder not found"),
	}
	fs := &mockFileSystem{
		envFileContent: []byte("KEY=value\n"),
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	err := PushEnvCore(
		".",
		"my-project",
		fs,
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
}

// 異常系: CreateNoteItem がエラーを返す
func TestPushEnvCore_CreateItemError(t *testing.T) {
	bw := &mockBwClient{
		folderID:   "folder-123",
		itemByName: nil,
		createErr:  errors.New("failed to create item"),
	}
	fs := &mockFileSystem{
		envFileContent: []byte("KEY=value\n"),
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	err := PushEnvCore(
		".",
		"my-project",
		fs,
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
}

// =============================================================================
// PullEnvCore のテスト
// =============================================================================

// 正常系: アイテムが存在し、出力先に .env がない場合に新規作成
func TestPullEnvCore_CreateNewEnvFile(t *testing.T) {
	bw := &mockBwClient{
		folderID:   "folder-123",
		itemByName: &FullItem{ID: "item-456", Name: "my-project", Notes: `{"lines":["KEY=value"]}`},
	}
	fs := &mockFileSystem{
		statInfo: &mockFileInfo{notExist: true},
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	err := PullEnvCore(
		".",
		"my-project",
		fs,
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		func(path string) (bool, error) { return true, nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

// 正常系: 出力ディレクトリが存在しない場合に MkdirAll が呼ばれる
func TestPullEnvCore_CreateOutputDirectory(t *testing.T) {
	bw := &mockBwClient{
		folderID:   "folder-123",
		itemByName: &FullItem{ID: "item-456", Name: "my-project", Notes: `{"lines":["KEY=value"]}`},
	}
	fs := &mockFileSystem{
		statInfo: &mockFileInfo{notExist: true},
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	err := PullEnvCore(
		"/custom/output",
		"my-project",
		fs,
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		func(path string) (bool, error) { return true, nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
}

// 正常系: outputDir が "." の場合に /project-root へ正規化
func TestPullEnvCore_NormalizeOutputDir(t *testing.T) {
	bw := &mockBwClient{
		folderID:   "folder-123",
		itemByName: &FullItem{ID: "item-456", Name: "my-project", Notes: `{"lines":["KEY=value"]}`},
	}
	fs := &mockFileSystem{
		statInfo: &mockFileInfo{notExist: true},
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	err := PullEnvCore(
		".",
		"my-project",
		fs,
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		func(path string) (bool, error) { return true, nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
}

// 異常系: アイテムが見つからない
func TestPullEnvCore_ItemNotFound(t *testing.T) {
	bw := &mockBwClient{
		folderID:   "folder-123",
		itemByName: nil, // アイテムなし
	}
	fs := &mockFileSystem{}
	logger := &mockLogger{}
	cfg := &config.Config{}

	err := PullEnvCore(
		".",
		"my-project",
		fs,
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		func(path string) (bool, error) { return true, nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
}

// 異常系: confirmOverwrite が false を返す場合
func TestPullEnvCore_OverwriteCancelled(t *testing.T) {
	bw := &mockBwClient{
		folderID:   "folder-123",
		itemByName: &FullItem{ID: "item-456", Name: "my-project", Notes: `{"lines":["KEY=value"]}`},
	}
	fs := &mockFileSystem{
		statInfo: &mockFileInfo{notExist: false}, // ファイルが存在する
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	err := PullEnvCore(
		".",
		"my-project",
		fs,
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		func(path string) (bool, error) { return false, nil }, // キャンセル
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
}

// 異常系: RestoreEnvFileFromJSON が壊れた JSON でエラー
func TestPullEnvCore_InvalidJSON(t *testing.T) {
	bw := &mockBwClient{
		folderID:   "folder-123",
		itemByName: &FullItem{ID: "item-456", Name: "my-project", Notes: "not valid json"},
	}
	fs := &mockFileSystem{
		statInfo: &mockFileInfo{notExist: true},
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	err := PullEnvCore(
		".",
		"my-project",
		fs,
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		func(path string) (bool, error) { return true, nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
}

// =============================================================================
// ListDotenvsCore のテスト
// =============================================================================

// 正常系: 3 件のアイテムが返る
func TestListDotenvsCore_ReturnsItems(t *testing.T) {
	bw := &mockBwClient{
		folderID: "folder-123",
		items: []Item{
			{ID: "1", Name: "project-a"},
			{ID: "2", Name: "project-b"},
			{ID: "3", Name: "project-c"},
		},
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	items, err := ListDotenvsCore(
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
	assert.Nil(t, items)
}

// 正常系: 空スライスが返る
func TestListDotenvsCore_ReturnsEmptySlice(t *testing.T) {
	bw := &mockBwClient{
		folderID: "folder-123",
		items:    []Item{},
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	items, err := ListDotenvsCore(
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
	assert.Nil(t, items)
}

// 異常系: GetDotenvsFolderID がロック関連エラーを返し、リトライも失敗
func TestListDotenvsCore_FolderIDLockError(t *testing.T) {
	bw := &mockBwClient{
		folderIDErr: ErrBitwardenLocked,
		unlockErr:   errors.New("unlock failed"),
		loginErr:    errors.New("login failed"),
	}
	logger := &mockLogger{}
	cfg := &config.Config{Email: "test@example.com"}

	items, err := ListDotenvsCore(
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
	assert.Nil(t, items)
}

// 異常系: ListItemsInFolder がエラーを返す
func TestListDotenvsCore_ListItemsError(t *testing.T) {
	bw := &mockBwClient{
		folderID:     "folder-123",
		listItemsErr: errors.New("failed to list items"),
	}
	logger := &mockLogger{}
	cfg := &config.Config{}

	items, err := ListDotenvsCore(
		bw,
		cfg,
		func() (string, error) { return "pwd", nil },
		logger,
	)

	// Red フェーズ
	assert.Error(t, err)
	assert.Nil(t, items)
}

// =============================================================================
// SetupBitwardenCore のテスト
// =============================================================================

// 正常系: cloud が選択され、Login 成功、SaveConfig が呼ばれる
func TestSetupBitwardenCore_CloudSuccess(t *testing.T) {
	bw := &mockBwClient{}
	fs := &mockFileSystem{}
	logger := &mockLogger{}

	err := SetupBitwardenCore(
		fs,
		bw,
		logger,
		func() (string, error) { return "cloud", nil },
		func() (string, error) { return "", errors.New("should not be called") },
		func() (string, error) { return "test@example.com", nil },
		func() (string, error) { return "password123", nil },
	)

	// Red フェーズ
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

// 正常系: selfhosted が選択され、URL 入力、Login 成功
func TestSetupBitwardenCore_SelfhostedSuccess(t *testing.T) {
	bw := &mockBwClient{}
	fs := &mockFileSystem{}
	logger := &mockLogger{}

	err := SetupBitwardenCore(
		fs,
		bw,
		logger,
		func() (string, error) { return "selfhosted", nil },
		func() (string, error) { return "https://bw.example.com", nil },
		func() (string, error) { return "test@example.com", nil },
		func() (string, error) { return "password123", nil },
	)

	// Red フェーズ
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

// 異常系: selectHostType がエラーを返す
func TestSetupBitwardenCore_SelectHostTypeError(t *testing.T) {
	bw := &mockBwClient{}
	fs := &mockFileSystem{}
	logger := &mockLogger{}

	err := SetupBitwardenCore(
		fs,
		bw,
		logger,
		func() (string, error) { return "", errors.New("user cancelled") },
		func() (string, error) { return "", nil },
		func() (string, error) { return "", nil },
		func() (string, error) { return "", nil },
	)

	// Red フェーズ
	assert.Error(t, err)
}

// 異常系: Login がエラーを返す場合、設定保存が行われない
func TestSetupBitwardenCore_LoginError(t *testing.T) {
	bw := &mockBwClient{
		loginErr: errors.New("invalid credentials"),
	}
	fs := &mockFileSystem{}
	logger := &mockLogger{}

	err := SetupBitwardenCore(
		fs,
		bw,
		logger,
		func() (string, error) { return "cloud", nil },
		func() (string, error) { return "", nil },
		func() (string, error) { return "test@example.com", nil },
		func() (string, error) { return "wrongpassword", nil },
	)

	// Red フェーズ
	assert.Error(t, err)
}

// =============================================================================
// BwClient インターフェースのテスト（モック検証）
// =============================================================================

// GetDotenvsFolderID: 正常系
func TestBwClient_GetDotenvsFolderID_Success(t *testing.T) {
	bw := &mockBwClient{folderID: "test-folder-id"}

	id, err := bw.GetDotenvsFolderID()

	assert.NoError(t, err)
	assert.Equal(t, "test-folder-id", id)
	assert.Contains(t, bw.calls, "GetDotenvsFolderID")
}

// GetDotenvsFolderID: フォルダが見つからない
func TestBwClient_GetDotenvsFolderID_NotFound(t *testing.T) {
	bw := &mockBwClient{
		folderIDErr: errors.New("dotenvs folder not found"),
	}

	id, err := bw.GetDotenvsFolderID()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dotenvs folder not found")
	assert.Empty(t, id)
}

// GetDotenvsFolderID: ロック状態
func TestBwClient_GetDotenvsFolderID_Locked(t *testing.T) {
	bw := &mockBwClient{
		folderIDErr: ErrBitwardenLocked,
	}

	id, err := bw.GetDotenvsFolderID()

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBitwardenLocked)
	assert.Empty(t, id)
}

// ListItemsInFolder: 正常系 - 3 件
func TestBwClient_ListItemsInFolder_Success(t *testing.T) {
	bw := &mockBwClient{
		items: []Item{
			{ID: "1", Name: "a"},
			{ID: "2", Name: "b"},
			{ID: "3", Name: "c"},
		},
	}

	items, err := bw.ListItemsInFolder("folder-id")

	assert.NoError(t, err)
	assert.Len(t, items, 3)
}

// ListItemsInFolder: 正常系 - 空
func TestBwClient_ListItemsInFolder_Empty(t *testing.T) {
	bw := &mockBwClient{
		items: []Item{},
	}

	items, err := bw.ListItemsInFolder("folder-id")

	assert.NoError(t, err)
	assert.Empty(t, items)
}

// ListItemsInFolder: ロック状態
func TestBwClient_ListItemsInFolder_Locked(t *testing.T) {
	bw := &mockBwClient{
		listItemsErr: ErrBitwardenLocked,
	}

	items, err := bw.ListItemsInFolder("folder-id")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBitwardenLocked)
	assert.Nil(t, items)
}

// GetItemByName: 正常系 - アイテム発見
func TestBwClient_GetItemByName_Found(t *testing.T) {
	bw := &mockBwClient{
		itemByName: &FullItem{ID: "item-1", Name: "my-project", Notes: "{}"},
	}

	item, err := bw.GetItemByName("folder-id", "my-project")

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "my-project", item.Name)
}

// GetItemByName: 正常系 - アイテムなし
func TestBwClient_GetItemByName_NotFound(t *testing.T) {
	bw := &mockBwClient{
		itemByName: nil,
	}

	item, err := bw.GetItemByName("folder-id", "nonexistent")

	assert.NoError(t, err)
	assert.Nil(t, item)
}

// GetItemByName: ロック状態
func TestBwClient_GetItemByName_Locked(t *testing.T) {
	bw := &mockBwClient{
		itemByNameErr: ErrBitwardenLocked,
	}

	item, err := bw.GetItemByName("folder-id", "my-project")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBitwardenLocked)
	assert.Nil(t, item)
}

// GetItemByID: 正常系
func TestBwClient_GetItemByID_Success(t *testing.T) {
	bw := &mockBwClient{
		itemByID: &FullItem{ID: "item-123", Name: "test", Notes: "content"},
	}

	item, err := bw.GetItemByID("item-123")

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "item-123", item.ID)
}

// GetItemByID: 異常系 - 空出力
func TestBwClient_GetItemByID_Empty(t *testing.T) {
	bw := &mockBwClient{
		itemByIDErr: errors.New("no output from bw get item command"),
	}

	item, err := bw.GetItemByID("item-123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no output")
	assert.Nil(t, item)
}

// GetItemByID: ロック状態
func TestBwClient_GetItemByID_Locked(t *testing.T) {
	bw := &mockBwClient{
		itemByIDErr: ErrBitwardenLocked,
	}

	item, err := bw.GetItemByID("item-123")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBitwardenLocked)
	assert.Nil(t, item)
}

// CreateNoteItem: 正常系
func TestBwClient_CreateNoteItem_Success(t *testing.T) {
	bw := &mockBwClient{}

	err := bw.CreateNoteItem("folder-id", "new-item", "{}")

	assert.NoError(t, err)
	assert.Contains(t, bw.calls, "CreateNoteItem(folder-id,new-item)")
}

// CreateNoteItem: 異常系
func TestBwClient_CreateNoteItem_Error(t *testing.T) {
	bw := &mockBwClient{
		createErr: errors.New("failed to create item"),
	}

	err := bw.CreateNoteItem("folder-id", "new-item", "{}")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create item")
}

// UpdateNoteItem: 正常系
func TestBwClient_UpdateNoteItem_Success(t *testing.T) {
	bw := &mockBwClient{}

	err := bw.UpdateNoteItem("item-123", "new notes")

	assert.NoError(t, err)
	assert.Contains(t, bw.calls, "UpdateNoteItem(item-123)")
}

// UpdateNoteItem: 異常系
func TestBwClient_UpdateNoteItem_Error(t *testing.T) {
	bw := &mockBwClient{
		updateErr: errors.New("failed to update item"),
	}

	err := bw.UpdateNoteItem("item-123", "new notes")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update item")
}

// Login: 正常系
func TestBwClient_Login_Success(t *testing.T) {
	bw := &mockBwClient{}

	err := bw.Login("test@example.com", "password", "")

	assert.NoError(t, err)
	assert.Contains(t, bw.calls, "Login(test@example.com,)")
}

// Login: 異常系
func TestBwClient_Login_Error(t *testing.T) {
	bw := &mockBwClient{
		loginErr: errors.New("invalid credentials"),
	}

	err := bw.Login("test@example.com", "wrong", "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

// Unlock: 正常系
func TestBwClient_Unlock_Success(t *testing.T) {
	bw := &mockBwClient{}

	err := bw.Unlock("masterpassword")

	assert.NoError(t, err)
	assert.Contains(t, bw.calls, "Unlock")
}

// Unlock: 異常系
func TestBwClient_Unlock_Error(t *testing.T) {
	bw := &mockBwClient{
		unlockErr: errors.New("wrong master password"),
	}

	err := bw.Unlock("wrongpassword")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wrong master password")
}

// =============================================================================
// FileSystem インターフェースのテスト（モック検証）
// =============================================================================

// WriteFile 後に同じ内容が取得できる
func TestFileSystem_WriteAndRead(t *testing.T) {
	fs := &mockFileSystem{}

	err := fs.WriteFile("/test/.env", []byte("KEY=value"), 0644)
	assert.NoError(t, err)
	assert.Equal(t, "/test/.env", fs.writtenPath)
	assert.Equal(t, []byte("KEY=value"), fs.writtenData)
}

// MkdirAll: 正常系
func TestFileSystem_MkdirAll_Success(t *testing.T) {
	fs := &mockFileSystem{}

	err := fs.MkdirAll("/test/dir", 0755)

	assert.NoError(t, err)
	assert.Contains(t, fs.calls, "MkdirAll(/test/dir)")
}

// MkdirAll: 異常系
func TestFileSystem_MkdirAll_Error(t *testing.T) {
	fs := &mockFileSystem{
		mkdirErr: errors.New("permission denied"),
	}

	err := fs.MkdirAll("/readonly/dir", 0755)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
}

// Stat: ファイルが存在しない
func TestFileSystem_Stat_NotExist(t *testing.T) {
	fs := &mockFileSystem{
		statInfo: &mockFileInfo{notExist: true},
	}

	info, err := fs.Stat("/nonexistent")

	assert.NoError(t, err)
	assert.True(t, info.IsNotExist())
}

// Stat: ファイルが存在する
func TestFileSystem_Stat_Exists(t *testing.T) {
	fs := &mockFileSystem{
		statInfo: &mockFileInfo{notExist: false},
	}

	info, err := fs.Stat("/existing")

	assert.NoError(t, err)
	assert.False(t, info.IsNotExist())
}
