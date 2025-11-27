# Archived TODOs

## v0.6.0

- ✅ GPG署名
- ✅ ドキュメント整備
- ✅ Goバージョンの差異を統一

## v0.5.0

- ✅ brew tapでの、配布。まずはmacOSビルドのみ
- ✅ ローカルのmacOSで実行されるかを確認


## v0.4.0

- ✅ 単体テストの実装
- ✅ ソースコード内へのバージョン指定

## v0.3.0

- ✅ bw push
- ✅ bw pull

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

### Fix: .configを永続化

現状、.configをDocker内に作っているが、デバッグがしづらいので、永続化したい
docker containerも立ち上がったままにしたい

### Feat: bwenv listコマンド

- bwenv listと入力すると
- Bitwardenホスト内の「dotenvs」フォルダを探す
- dotenvsフォルダがなければエラー
- dotenvsフォルダが存在すれば、その中のアイテムの名称一覧を出力

### Enhance: 標準出力に色をつけたい

- 赤：ERROR
- 緑：何かを実行して、SUCCESSした時
- 黄：WARNING
- 水色：重要なINFO、または決定事項
- 薄紫：対話の質問

### Enhance: cloudかselfhostedかの質問

- 現状: 文字入力させている
- 改善: 選択肢を表示させ、矢印キーで選択、Enterで確定させたい