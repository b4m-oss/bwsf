package infra

import (
	"bwenv/src/core"
	"bwenv/src/utils"
)

// RealBwClient は core.BwClient インターフェースの実装で、
// utils パッケージの既存関数をラップします。
type RealBwClient struct{}

// NewBwClient は RealBwClient のインスタンスを作成します。
func NewBwClient() *RealBwClient {
	return &RealBwClient{}
}

// GetDotenvsFolderID は dotenvs フォルダの ID を取得します。
func (c *RealBwClient) GetDotenvsFolderID() (string, error) {
	return utils.GetDotenvsFolderID()
}

// ListItemsInFolder は指定フォルダ内のアイテム一覧を取得します。
func (c *RealBwClient) ListItemsInFolder(folderID string) ([]core.Item, error) {
	items, err := utils.ListItemsInFolder(folderID)
	if err != nil {
		return nil, err
	}

	// utils.Item から core.Item に変換
	result := make([]core.Item, len(items))
	for i, item := range items {
		result[i] = core.Item{
			ID:   item.ID,
			Name: item.Name,
		}
	}
	return result, nil
}

// GetItemByName は指定フォルダ内のアイテムを名前で検索します。
func (c *RealBwClient) GetItemByName(folderID, name string) (*core.FullItem, error) {
	item, err := utils.GetItemByName(folderID, name)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}

	// utils.FullItem から core.FullItem に変換
	return &core.FullItem{
		ID:    item.ID,
		Name:  item.Name,
		Notes: item.Notes,
	}, nil
}

// GetItemByID は指定 ID のアイテムを取得します。
func (c *RealBwClient) GetItemByID(id string) (*core.FullItem, error) {
	item, err := utils.GetItemByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}

	// utils.FullItem から core.FullItem に変換
	return &core.FullItem{
		ID:    item.ID,
		Name:  item.Name,
		Notes: item.Notes,
	}, nil
}

// CreateNoteItem は新しいノートアイテムを作成します。
func (c *RealBwClient) CreateNoteItem(folderID, name, notes string) error {
	return utils.CreateNoteItem(folderID, name, notes)
}

// UpdateNoteItem は既存のノートアイテムを更新します。
func (c *RealBwClient) UpdateNoteItem(id, notes string) error {
	return utils.UpdateNoteItem(id, notes)
}

// Login は Bitwarden CLI にログインします。
func (c *RealBwClient) Login(email, password, serverURL string) error {
	success, errorMsg := utils.BwLogin(email, password, serverURL)
	if !success {
		return &LoginError{Message: errorMsg}
	}
	return nil
}

// Unlock は Bitwarden CLI をアンロックします。
func (c *RealBwClient) Unlock(masterPassword string) error {
	success, errorMsg := utils.BwUnlock(masterPassword)
	if !success {
		return &UnlockError{Message: errorMsg}
	}
	return nil
}

// LoginError はログイン失敗時のエラーです。
type LoginError struct {
	Message string
}

func (e *LoginError) Error() string {
	return e.Message
}

// UnlockError はアンロック失敗時のエラーです。
type UnlockError struct {
	Message string
}

func (e *UnlockError) Error() string {
	return e.Message
}

