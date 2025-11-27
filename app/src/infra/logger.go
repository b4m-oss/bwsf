package infra

import (
	"bwenv/src/utils"
)

// RealLogger は core.Logger インターフェースの実装で、
// utils パッケージのカラー出力関数をラップします。
type RealLogger struct{}

// NewLogger は RealLogger のインスタンスを作成します。
func NewLogger() *RealLogger {
	return &RealLogger{}
}

// Error はエラーメッセージを出力します。
func (l *RealLogger) Error(args ...interface{}) {
	utils.Errorln(args...)
}

// Info は情報メッセージを出力します。
func (l *RealLogger) Info(args ...interface{}) {
	utils.Infoln(args...)
}

// Success は成功メッセージを出力します。
func (l *RealLogger) Success(args ...interface{}) {
	utils.Successln(args...)
}

// Warning は警告メッセージを出力します。
func (l *RealLogger) Warning(args ...interface{}) {
	utils.Warningln(args...)
}

