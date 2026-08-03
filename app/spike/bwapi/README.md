# bwapi — Issue #53 Step 0 spike

本番の `bwsf` CLI とは **別モジュール・別バイナリ** の使い捨て実験です。  
`app/go.mod` や `app/src/` には SDK を入れず、コミュニティ Go SDK  
[`github.com/bnema/bitwarden-go-sdk`](https://github.com/bnema/bitwarden-go-sdk)（v0.4.x）で vault 操作が通るかだけを検証します。

> **注意:** SDK は alpha。公式サポートは Bitwarden cloud / 公式 self-hosted のみで、**Vaultwarden は非目標**と明記されています。HTTPS 必須（`http://` は SDK が拒否）。

---

## 必要な環境変数

### Scenario A（メイン go/no-go）

| 変数 | 必須 | 説明 |
|------|------|------|
| `BWSF_SPIKE_EMAIL` | ○ | ログイン用メール |
| `BWSF_SPIKE_PASSWORD` | ○ | マスターパスワード |
| `BWSF_SPIKE_SERVER_URL` | − | セルフホスト / Vaultwarden のベース URL（例: `https://vw.example.com`、末尾 `/` 可）。未設定時は Bitwarden US cloud。SDK は `{url}/identity` と `{url}/api` に導出 |
| `BWSF_SPIKE_TOTP` | − | 2FA が必要な場合の authenticator コード。未設定なら対話プロンプト |

### Scenario B（Personal API Key プローブ）

| 変数 | 必須 | 説明 |
|------|------|------|
| `BWSF_SPIKE_CLIENT_ID` | ○ | Personal API Key の `client_id`（例: `user.…`） |
| `BWSF_SPIKE_CLIENT_SECRET` | ○ | Personal API Key の `client_secret` |
| `BWSF_SPIKE_SERVER_URL` | − | 設定時は Identity を `<url>/identity` と導出 |
| `BWSF_SPIKE_IDENTITY_URL` | − | Identity ベース URL の上書き（未設定時: cloud は `https://identity.bitwarden.com`） |

シークレットをコードや git に書かないこと。

---

## ビルド・実行

SDK は **Go 1.26+** を要求します（`go.mod` の `go 1.26`）。

### goenv 注意（このマシン）

goenv は `go.mod` の `go 1.26` を見て **GOTOOLCHAIN より先に** バージョン解決します。  
`goenv versions` に 1.26 が無いと、次で失敗します:

```text
goenv: version '1.26' is not installed
```

インストール済みは例: `1.25.7`（global）。対処は次のどちらか。

#### 推奨: `./run.sh`（GOENV_VERSION + GOTOOLCHAIN=auto）

`run.sh` がインストール済みの 1.25.x を検出し、`GOENV_VERSION` と `GOTOOLCHAIN=auto` を付けて `go` を呼びます（toolchain が 1.26 を取得）。

```bash
cd app/spike/bwapi

# 依存整理
./run.sh mod tidy

# コンパイル確認（認証情報なしで OK）
./run.sh build -o /tmp/bwapi .
# または: make build

# Scenario A（ライブ）
export BWSF_SPIKE_EMAIL='you@example.com'
export BWSF_SPIKE_PASSWORD='…'
# 任意: export BWSF_SPIKE_SERVER_URL='https://your-server'
# 任意: export BWSF_SPIKE_TOTP='123456'
./run.sh
# または: make run / ./run.sh run .

# Scenario B（ライブ）
export BWSF_SPIKE_CLIENT_ID='user.…'
export BWSF_SPIKE_CLIENT_SECRET='…'
./run.sh apikey
# または: ./run.sh run . apikey
```

手動で同じことをする場合（このマシンで確認済み）:

```bash
cd app/spike/bwapi
GOENV_VERSION=1.25.7 GOTOOLCHAIN=auto go build -o /tmp/bwapi .
GOENV_VERSION=1.25.7 GOTOOLCHAIN=auto go run .
```

#### 代替: goenv に 1.26 を入れる

```bash
goenv install 1.26
cd app/spike/bwapi
go run .   # または GOTOOLCHAIN=auto go run .
```

本番 `make build` / Dockerfile / brew には **含まれません**。

---

## 成功条件

### Scenario A（go/no-go）

stdout に次が揃うこと:

1. `NewClient` 成功（cloud または `WithServerURL`）
2. `BeginLogin`（必要なら `CompleteLogin`）成功
3. `Sync` 成功
4. フォルダ `bwsf-spike` が既存 or 新規作成
5. Secure Note `bwsf-spike-note` の create → get/search → notes update 成功
6. 最後に `SUCCESS (Scenario A): …`

### Scenario B

- Identity `POST …/connect/token`（`grant_type=client_credentials`, `scope=api` + device メタ）で **`access_token` 取得**
- vault CRUD は不要
- SDK に API Key ログインが無いため、スパイクは raw HTTP でプローブする

### Scenario C（PO Vaultwarden スモーク）

1. Vaultwarden を **HTTPS** で用意（SDK は `http://` を拒否）
2. `export BWSF_SPIKE_SERVER_URL='https://<your-vaultwarden>'`
3. A と同じく email / password を設定し `go run .` を実行
4. 結果（成功 / 失敗メッセージ全文）を Issue #53 に報告

公式には VW 非対応のため、失敗してもスパイクの価値あり（制約の実証）。

---

## Vaultwarden troubleshooting

### 期待する URL 形

`BWSF_SPIKE_SERVER_URL=https://vw.example.com`（末尾スラッシュはスパイク／SDK が除去）のとき:

| 用途 | URL |
|------|-----|
| Identity | `https://vw.example.com/identity` |
| API | `https://vw.example.com/api` |
| Prelogin | `POST …/identity/accounts/prelogin/password` |
| Token | `POST …/identity/connect/token` |

末尾 `/` だけでは 400 にはなりません。SDK の `NewSelfHosted` が path の trailing slash を strip します。

### `BeginLogin: bitwarden: unknown status=400`

よくある原因（このスパイクで確認済み）:

1. **デバイスメタ未送信** — Vaultwarden は `/identity/connect/token` で空白を拒否する。  
   - Scenario A（password）: `device_name cannot be blank`  
   - Scenario B（`client_credentials`）: `device_identifier cannot be blank`（続けて name / type も必須）  
   公式 cloud より厳しい。スパイクは常に `deviceIdentifier` / `deviceName` / `deviceType` を付ける（A は SDK `LoginOptions`、B は raw form）。
2. **SDK がボディを捨てる** — VW の 400 JSON は `"error":""`（空文字）になりがちで、SDK は OAuth の `error` が空だと `KindUnknown` + `status=400` だけにする。本物の理由は `message` / `errorModel.message` 側。
3. **認証失敗も同じ opaque 400** — パスワード不正時も VW は同様の形のため、SDK 上は `unknown status=400` に見えることがある。

失敗時のスパイク出力:

- 導出済み `identity` / `api` URL
- SDK `Error` の kind / status / code / message（Scenario A）
- パスワードを使わない diagnose POST（`deviceName` あり／なし）でレスポンスボディ例を表示（Scenario A）
- Scenario B は raw HTTP のため、失敗時にレスポンスボディをそのまま表示

### VW API 不一致 vs 設定ミス

| 症状 | 解釈 |
|------|------|
| prelogin 200、token で `device_name` / `device_identifier cannot be blank` | **設定／クライアント側**（デバイスメタ不足）。修正可能 |
| prelogin / identity が 404 | ベース URL 誤り、またはリバースプロキシで `/identity` が届いていない |
| token が常に opaque 400 で diagnose も意味不明 | **VW と SDK の API 差**の可能性。Issue に diagnose 全文を貼る |

---

## Identity URL の導出（Scenario B）

| 条件 | Identity base | token endpoint |
|------|---------------|----------------|
| `BWSF_SPIKE_IDENTITY_URL` あり | その値 | `{base}/connect/token` |
| `BWSF_SPIKE_SERVER_URL` あり | `{server}/identity` | `{server}/identity/connect/token` |
| どちらも無し（cloud US） | `https://identity.bitwarden.com` | `…/connect/token` |

---

## 既知の SDK / 運用ギャップ（実装時メモ）

- **Personal API Key:** `bitwarden-go-sdk` v0.4.0 に client_credentials / API Key ログイン API なし → Scenario B は raw HTTP。
- **Vaultwarden:** README で non-goal。加えて **HTTPS 必須**（ローカル HTTP VW はそのままでは `NewClient(WithServerURL)` 不可）。
- **device メタ必須 (VW):** password / client_credentials とも `deviceIdentifier`・`deviceName`・`deviceType` が必要。SDK は `LoginOptions` が空だと form に載せない。スパイクは固定デバイスメタを送る。
- **opaque 400:** VW の token エラー JSON は `"error":""` が多く、SDK がボディを捨てて `unknown status=400` になる。失敗時は diagnose 出力を参照。
- **2FA:** `CompleteLogin` + `TwoFactorProviderAuthenticator`。他プロバイダは要コード変更。
- alpha API は破壊的変更があり得る。
