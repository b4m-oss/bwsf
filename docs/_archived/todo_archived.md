# Archived TODOs

## v0.2.0

### Feat: bwコマンドの存在をチェック

- 現在の端末に、bwコマンドがインストールされているかどうかを確認
- インストールされていれば「[INFO] ✅ bw command is installed!」と標準出力
- インストールされていなければ「[ERROR] ❌ bw command is not installed...」と標準出力

### Feat: Bitwardenセットアップ

- Bitwardenホストにログインする
- Bitwarden Cloudか、セルフホストかを確認する。対話で選択式。
- セルフホストの場合は、URLを入力させる
- Eメールアドレス、パスワードを入力させる。パスワードは入力を隠す。
- bwコマンドを使って、ログインを試みる
- 失敗したらエラーメッセージをそのまま表示
- ログインできたら「[INFO] ✅ Sign in to Bitwarden was successful!」と標準出力