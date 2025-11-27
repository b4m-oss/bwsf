package core

import (
	"fmt"

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
// ここではテストで呼び出しを検証できるよう、最低限のメソッドだけ定義します。
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

// PushEnvCore / PullEnvCore / ListDotenvsCore / SetupBitwardenCore は
// テスト駆動で実装していくコアロジック関数です。
// 現時点では Red フェーズのため、未実装エラーを返します。

func PushEnvCore(
	fromDir, projectName string,
	fs FileSystem,
	bw BwClient,
	cfg *config.Config,
	promptPassword func() (string, error),
	logger Logger,
) error {
	return fmt.Errorf("PushEnvCore not implemented yet")
}

func PullEnvCore(
	outputDir, projectName string,
	fs FileSystem,
	bw BwClient,
	cfg *config.Config,
	promptPassword func() (string, error),
	confirmOverwrite func(path string) (bool, error),
	logger Logger,
) error {
	return fmt.Errorf("PullEnvCore not implemented yet")
}

func ListDotenvsCore(
	bw BwClient,
	cfg *config.Config,
	promptPassword func() (string, error),
	logger Logger,
) ([]Item, error) {
	return nil, fmt.Errorf("ListDotenvsCore not implemented yet")
}

func SetupBitwardenCore(
	fs FileSystem,
	bw BwClient,
	logger Logger,
	selectHostType func() (string, error),
	inputURL func() (string, error),
	inputEmail func() (string, error),
	inputPassword func() (string, error),
) error {
	return fmt.Errorf("SetupBitwardenCore not implemented yet")
}

// WithUnlockRetry は Bitwarden がロックされている場合に Unlock/Login を挟んでリトライする共通処理です。
// TDD の Red フェーズのため、現時点では未実装のスタブとして panic を置いています。
func WithUnlockRetry(
	bw BwClient,
	cfg *config.Config,
	promptPassword func() (string, error),
	logger Logger,
	fn func() error,
) error {
	// Red フェーズ用の暫定実装: まずは単純に fn() を一度だけ呼び出す。
	// Unlock/Login の挙動は今後テストを増やしながら実装していく。
	if err := fn(); err != nil {
		return err
	}
	return nil
}


