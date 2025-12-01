package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// CheckBwCommand のテスト
// =============================================================================

// 注: これらのテストは実際の bw コマンドの存在に依存するため、
// CI 環境では条件付きでスキップするか、モック化が必要

// 正常系: bw コマンドが存在する場合（環境依存）
func TestCheckBwCommand_Exists(t *testing.T) {
	// 実際の環境で bw がインストールされているかどうかで結果が変わる
	installed, path := CheckBwCommand()

	if installed {
		assert.NotEmpty(t, path)
	} else {
		assert.Empty(t, path)
	}
}

// =============================================================================
// BwLogin の構造テスト
// =============================================================================

// 注: BwLogin は実際の bw コマンドを呼び出すため、
// 単体テストではモック化が必要。ここでは関数シグネチャの確認のみ

// 正常系: 関数が正しいシグネチャを持つ
func TestBwLogin_Signature(t *testing.T) {
	// BwLogin の型を確認（コンパイル時チェック）
	var fn func(email, password, serverURL string) (bool, string) = BwLogin
	assert.NotNil(t, fn)
}



