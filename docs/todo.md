# Todo

## v0.2.0

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