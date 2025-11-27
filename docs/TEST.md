# テスト仕様書

このドキュメントは、`bwenv` CLI のリファクタ後の構成と、それぞれの関数／メソッドの処理内容をテスト仕様書としてまとめたものです。

---

## 1. 全体構成

### main / CLI エントリポイント

- `main`
  - `cmd.Execute()` を呼び出し、Cobra ベースの CLI を起動する。
- `Execute`
  - `rootCmd.Execute()` を実行して CLI コマンド群を起動する。
  - `rootCmd.Execute()` がエラーを返した場合は `Logger` を通じてエラーメッセージを表示する。
  - エラー発生時には `os.Exit(1)` でプロセスを終了する。

### レイヤ構造の方針

- **CLI 層**
  - Cobra コマンド (`runPush`, `runPull`, `runList`, `runSetup`) からなる。
  - 引数／フラグのパースと、エラー時の Exit/メッセージ表示のみを担当する「薄いラッパー」とする。
  - 実際の処理はすべて「コアロジック層」の関数に委譲する。
- **コアロジック層**
  - ビジネスロジック（Push/Pull/List/Setup）を環境依存のない形で実装する。
  - I/O は `FileSystem` / `BwClient` / `Logger` インターフェースに抽象化する。
- **インフラ層**
  - `exec.Command` による `bw` CLI 呼び出しや、`os` ベースのファイル操作など、現行実装に近い具体クラスを持つ。

---

## 2. 抽象インターフェース

### BwClient インターフェース

- `type BwClient interface`
  - Bitwarden CLI (`bw`) とのやり取りを抽象化する。
  - すべてのメソッドは **ネットワーク／CLI 側のエラー** を `error` で返す。

#### GetDotenvsFolderID

- "dotenvs" フォルダの ID を取得する。
- フォルダが存在しない場合は `"dotenvs folder not found"` に相当するエラーを返す。
- Bitwarden CLI がロックされている場合は、`ErrBitwardenLocked` を含むエラーを返す。

#### ListItemsInFolder

- 指定フォルダ ID 内のアイテム一覧を `[]Item` として返す。
- フォルダが空の場合でも、空スライスと `nil` エラーを返す。
- CLI ロック時は `ErrBitwardenLocked` を含むエラーを返す。

#### GetItemByName

- 指定フォルダ ID 内から、名前が一致するアイテムを検索して `*FullItem` を返す。
- 見つからない場合は `nil, nil` を返す。
- ロック状態の場合は `ErrBitwardenLocked` を含むエラーを返す。

#### GetItemByID

- 指定 ID のアイテムを取得して `*FullItem` を返す。
- 存在しない ID の場合は適切なエラーを返す。
- ロック状態の場合は `ErrBitwardenLocked` を含むエラーを返す。

#### CreateNoteItem

- フォルダ ID／名前／ノート文字列を受け取り、Secure Note（タイプ 2）として新規アイテムを作成する。
- 成功時は `nil`、失敗時は CLI 出力を含んだエラーを返す。

#### UpdateNoteItem

- アイテム ID と新しいノート文字列を受け取り、既存アイテムの `notes` フィールドを更新する。
- 成功時は `nil` を返す。
- ロック状態の場合は `ErrBitwardenLocked` を含むエラーを返す。

#### Login

- メールアドレス・パスワード・サーバー URL を受け取り、`bw login` 相当の処理を行う。
- 成功時は `nil` を返す。
- 失敗時は CLI 出力を含んだエラーを返す。

#### Unlock

- マスターパスワードを受け取り、`bw unlock` 相当の処理を行う。
- セッションキーの設定／検証まで内部で行い、成功時は `nil` を返す。
- ロック解除ができなかった場合は、詳細なメッセージを含むエラーを返す。

### FileSystem インターフェース

- `type FileSystem interface`
  - 実ディスクへのアクセスを抽象化する。

#### OpenEnvFile

- `.env` ファイルのパスを受け取り、読み取り用ハンドル／リーダを返す。
- 存在しない場合は `"not found"` 相当のエラーを返す。

#### ReadFile

- 任意パスのファイル内容を文字列として読み込んで返す。

#### WriteFile

- 任意パスに文字列を保存する。
- パーミッションは呼び出し側（コアロジック）から指定できるようにする。

#### Stat

- 任意パスの存在確認／属性取得を行う。
- 存在しない場合は明示的に判定できるエラー（`os.IsNotExist` と同等の扱い）を返す。

#### MkdirAll

- ディレクトリを再帰的に作成する。
- すでに存在する場合は成功扱いとする。

### Logger インターフェース

- `type Logger interface`
  - 既存の `utils.Error/Success/...` をラップするための抽象化。

#### Error / Errorf / Errorln

- エラー系メッセージを標準エラー出力に送る。
- カラー出力を有効／無効にするかは具象実装に委ねる。

#### Info / Infoln

- 情報メッセージを標準出力に送る。

#### Success / Successln

- 成功メッセージを標準出力に送る。

#### Warning / Warningln

- 警告メッセージを標準出力に送る。

---

## 3. 共通ユーティリティ（ロック時リトライ）

### WithUnlockRetry

- シグネチャ（イメージ）:
  - `func WithUnlockRetry(bw BwClient, cfg *config.Config, promptPassword func() (string, error), logger Logger, fn func() error) error`
- 処理内容:
  - `fn()` を一度実行し、エラーが `ErrBitwardenLocked` またはロック関連エラーの場合のみ、アンロック／ログインを行う。
  - アンロックフロー:
    - `promptPassword` を使ってマスターパスワードを取得する。
    - `bw.Unlock` を実行する。
    - 失敗した場合、`cfg` が存在すれば `bw.Login(cfg.Email, password, cfg.SelfhostedURL)` を試みる。
    - ログイン成功後に再度 `bw.Unlock` を実行する。
  - アンロック／ログインに成功した場合、`fn()` を再実行し、その結果を返す。
  - 途中でのパスワード取得失敗やアンロック／ログイン失敗は、適切なエラーメッセージをラップして返す。

### ParseEnvFile / EnvDataToJSON / RestoreEnvFileFromJSON / FindEnvFile

- 既存の挙動を維持するが、テストでは `FileSystem` のモックを介して利用できるようにする（`os` 直呼び出しを内部に閉じ込める）。

---

## 4. コアロジック層の関数

### PushEnvCore

- シグネチャ（イメージ）:
  - `func PushEnvCore(fromDir, projectName string, fs FileSystem, bw BwClient, cfg *config.Config, promptPassword func() (string, error), logger Logger) error`
- 処理:
  - `fromDir` と `projectName` に基づいて対象 `.env` ファイルパスを決定する。
    - `fromDir` が `"."` または `".."` の場合、`/project-root` へのフォールバックを含む。
  - `fs.OpenEnvFile` を利用して `.env` を開き、`ParseEnvFile` で `EnvData` にパースする。
  - `EnvDataToJSON` で JSON 文字列に変換する。
  - `WithUnlockRetry` を使って `"dotenvs"` フォルダ ID 取得 (`bw.GetDotenvsFolderID`) を行う。
  - 同じく `WithUnlockRetry` を通じて、`bw.GetItemByName(folderID, projectName)` で既存アイテムを検索する。
  - 以下の分岐:
    - 既存アイテムが **ある** 場合:
      - 呼び出し側（CLI 層）で「上書き確認」済みである前提とし、`WithUnlockRetry` を通じて `bw.UpdateNoteItem(item.ID, jsonData)` を実行する。
    - 既存アイテムが **ない** 場合:
      - `WithUnlockRetry` 経由で `bw.CreateNoteItem(folderID, projectName, jsonData)` を実行し、新規作成する。
  - エラー発生時はそれぞれのエラーを呼び出し元に返す（Exit は CLI 層で行う）。

### PullEnvCore

- シグネチャ（イメージ）:
  - `func PullEnvCore(outputDir, projectName string, fs FileSystem, bw BwClient, cfg *config.Config, promptPassword func() (string, error), confirmOverwrite func(path string) (bool, error), logger Logger) error`
- 処理:
  - `outputDir` が `"."` または `".."` の場合、`/project-root` に正規化する。
  - `WithUnlockRetry` 経由で `bw.GetDotenvsFolderID()` を呼び出し、フォルダ ID を取得する。
  - `WithUnlockRetry` 経由で `bw.GetItemByName(folderID, projectName)` を呼び出し、対応するアイテムを取得する。
  - アイテムが存在しない場合は `"Item '%s' not found in dotenvs folder"` 相当のエラーを返す。
  - `RestoreEnvFileFromJSON(item.Notes)` で `.env` 内容文字列を復元する。
  - 出力先 `.env` パス（`outputDir/.env`）を構成し、`fs.Stat` で存在を確認する。
    - すでに存在する場合は `confirmOverwrite(envPath)` を呼び出し、`false` の場合は `"operation cancelled"` 相当のエラーを返す、または `nil` を返して何もしない。
  - 必要に応じて `fs.MkdirAll(outputDir)` を呼び出し、ディレクトリを作成する。
  - `fs.WriteFile(envPath, content, 0644)` で `.env` を書き出す。

### ListDotenvsCore

- シグネチャ（イメージ）:
  - `func ListDotenvsCore(bw BwClient, cfg *config.Config, promptPassword func() (string, error), logger Logger) ([]Item, error)`
- 処理:
  - `WithUnlockRetry` 経由で `bw.GetDotenvsFolderID()` を実行し、フォルダ ID を取得する。
  - `WithUnlockRetry` 経由で `bw.ListItemsInFolder(folderID)` を実行し、アイテム一覧を取得する。
  - 取得した `[]Item` をそのまま呼び出し元に返す。
  - アイテム 0 件は正常ケースとして空スライスを返す。

### SetupBitwardenCore

- シグネチャ（イメージ）:
  - `func SetupBitwardenCore(fs FileSystem, bw BwClient, logger Logger, selectHostType func() (string, error), inputURL func() (string, error), inputEmail func() (string, error), inputPassword func() (string, error)) error`
- 処理:
  - 既存設定を `LoadConfig` で読み込み、存在する場合は上書きされる旨を `logger` で通知する。
  - `selectHostType` で `"cloud"` または `"selfhosted"` を選択させる。
  - `"selfhosted"` の場合は `inputURL` で Self-hosted URL を取得する。
  - `inputEmail` でメールアドレスを取得する。
  - `inputPassword` でパスワードを非エコー入力させる。
  - `bw.Login(email, password, selfhostedURL)` を実行し、ログインを試みる。
  - ログイン成功時に `SaveConfig` で設定（HostType / URL / Email）を保存する。
  - 失敗時にはエラーを呼び出し元に返す（Exit は CLI 層が行う）。

---

## 5. CLI 層（Cobra コマンド）

### runPush

- Cobra から呼ばれるエントリポイント。
- 処理:
  - `--from` フラグをパースして `fromDir` を取得する。
  - `os.Getwd()` からカレントディレクトリ名を取り出し、`projectName` として利用する。
  - 具象 `BwClient` / `FileSystem` / `Logger` / `config.Config` を生成する。
  - 上書き確認のための UI 関数（`confirmOverwrite`）を定義しておく。
  - `PushEnvCore` を呼び出し、戻り値のエラーに応じて:
    - エラーあり: `Logger` でエラーメッセージを出力し `os.Exit(1)`。
    - 正常終了: 成功メッセージを表示し、終了コード 0 で返る。

### runPull

- 処理:
  - `--output` フラグをパースして `outputDir` を取得する。
  - `os.Getwd()` から `projectName` を決定する。
  - 具象 `BwClient` / `FileSystem` / `Logger` を生成する。
  - 上書き確認用 `confirmOverwrite` 関数を用意する。
  - `PullEnvCore` を呼び出し、エラー時はメッセージ表示＋`os.Exit(1)`、成功時は成功メッセージを表示する。

### runList

- 処理:
  - 具象 `BwClient` / `Logger` / `config.Config` を生成する。
  - `ListDotenvsCore` を呼び出し、戻り値の `[]Item` を標準出力に 1 行ずつ出力する。
  - アイテム 0 件の場合は `"No items found in dotenvs folder"` 相当のメッセージのみを出力する。

### runSetup

- 処理:
  - `BwClient` / `FileSystem` / `Logger` を生成する。
  - 入力 UI (`SelectHostType` / `InputURL` / `InputEmail` / `InputPassword`) を関数として `SetupBitwardenCore` に渡す。
  - `SetupBitwardenCore` の戻り値がエラーの場合はメッセージ表示＋`os.Exit(1)`。
  - 正常終了時にはサインイン成功メッセージを表示する。

---

## 6. 既存ユーティリティの仕様（インフラ実装側）

### Config 周り

- `GetConfigPath`
  - ユーザーのホームディレクトリ配下の `.config/bwenv/config.json` へのフルパスを返す。
  - ホームディレクトリ取得に失敗した場合はエラー。
- `LoadConfig`
  - 設定ファイルが存在しない場合は `nil, nil` を返す（エラー扱いしない）。
  - 存在する場合は JSON を `Config` にアンマーシャルして返す。
- `SaveConfig`
  - 必要に応じてディレクトリを作成し、`Config` を `0600` で書き出す。

### 入力／色付き出力ユーティリティ

- `SelectHostType` / `InputURL` / `InputEmail` / `InputPassword` / `ConfirmOverwrite`
  - 既存仕様を維持しつつ、テストではモック関数に差し替えられるようにする。
- `Error` / `Errorln` / `Success` / `Successln` / `Warning` / `Warningln` / `Info` / `Infoln` / `Question` / `Questionln`
  - `Logger` の具象実装として利用される。

### `.env` パースユーティリティ

- `ParseEnvFile` / `EnvDataToJSON` / `RestoreEnvFileFromJSON` / `FindEnvFile`
  - 既存の仕様を維持する（行順／コメント／クォートをそのまま保持）。

### ロック判定ユーティリティ

- `ErrBitwardenLocked`
  - Bitwarden CLI がロックされていることを表す代表的なエラー。
- `IsLockedError`
  - 渡されたエラーがロック関連であるかを判定し、`WithUnlockRetry` などから利用される。

