package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// NoteItem / FullItem 構造体のテスト
// =============================================================================

// 正常系: NoteItem 構造体が正しく初期化できる
func TestNoteItem_Struct(t *testing.T) {
	item := NoteItem{
		Type:     2,
		Name:     "test-project",
		Notes:    `{"lines":["KEY=value"]}`,
		FolderID: "folder-123",
		SecureNote: SecureNote{
			Type: 0,
		},
	}

	assert.Equal(t, 2, item.Type)
	assert.Equal(t, "test-project", item.Name)
	assert.Equal(t, `{"lines":["KEY=value"]}`, item.Notes)
	assert.Equal(t, "folder-123", item.FolderID)
	assert.Equal(t, 0, item.SecureNote.Type)
}

// 正常系: FullItem 構造体が正しく初期化できる
func TestFullItem_Struct(t *testing.T) {
	item := FullItem{
		ID:       "item-456",
		Name:     "my-project",
		Type:     2,
		Notes:    `{"lines":["DB_HOST=localhost"]}`,
		FolderID: "folder-123",
		SecureNote: SecureNote{
			Type: 0,
		},
	}

	assert.Equal(t, "item-456", item.ID)
	assert.Equal(t, "my-project", item.Name)
	assert.Equal(t, 2, item.Type)
	assert.Equal(t, `{"lines":["DB_HOST=localhost"]}`, item.Notes)
	assert.Equal(t, "folder-123", item.FolderID)
}

// 正常系: SecureNote 構造体
func TestSecureNote_Struct(t *testing.T) {
	note := SecureNote{
		Type: 0, // Text type
	}

	assert.Equal(t, 0, note.Type)
}

// =============================================================================
// GetItemByName / GetItemByID / CreateNoteItem / UpdateNoteItem の構造テスト
// =============================================================================

// 注: これらの関数は実際の bw コマンドを呼び出すため、
// 単体テストではモック化が必要。ここでは関数シグネチャの確認のみ

// 正常系: GetItemByName の関数シグネチャ確認
func TestGetItemByName_Signature(t *testing.T) {
	var fn func(folderID, itemName string) (*FullItem, error) = GetItemByName
	assert.NotNil(t, fn)
}

// 正常系: GetItemByID の関数シグネチャ確認
func TestGetItemByID_Signature(t *testing.T) {
	var fn func(itemID string) (*FullItem, error) = GetItemByID
	assert.NotNil(t, fn)
}

// 正常系: CreateNoteItem の関数シグネチャ確認
func TestCreateNoteItem_Signature(t *testing.T) {
	var fn func(folderID, name, notes string) error = CreateNoteItem
	assert.NotNil(t, fn)
}

// 正常系: UpdateNoteItem の関数シグネチャ確認
func TestUpdateNoteItem_Signature(t *testing.T) {
	var fn func(itemID, notes string) error = UpdateNoteItem
	assert.NotNil(t, fn)
}


