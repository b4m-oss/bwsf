# インストール

## 動作環境

### 対応OS

| OS | 状態 |
|---|---|
| macOS | ✅ 対応 |
| Linux | ✅ 対応 |
| Windows | 🚧 計画中 |

### 依存関係

**Bitwarden CLI (`bw`)** が必要です。先にインストールしてください：

```bash
# macOS
brew install bitwarden-cli

# Linux (Snap)
sudo snap install bw

# npm (クロスプラットフォーム)
npm install -g @bitwarden/cli
```

その他のオプションについては、[公式 Bitwarden CLI ドキュメント](https://bitwarden.com/help/cli/#download-and-install)を参照してください。

## bwenv のインストール

### macOS

```bash
brew tap b4m-oss/tap && brew install bwenv
```

### Linux

::: tip
Linux では、先に [Homebrew on Linux](https://docs.brew.sh/Homebrew-on-Linux) をインストールする必要があります。
:::

```bash
brew tap b4m-oss/tap && brew install bwenv
```

## インストールの確認

```bash
bwenv -v
# bwenv version x.x.x
```

## 初期設定

インストール後、セットアップコマンドを実行して Bitwarden 接続を設定します：

```bash
bwenv setup
```

以下の入力を求められます：
1. Bitwarden サーバー URL（Bitwarden Cloud の場合は空欄）
2. Bitwarden のメールアドレス
3. マスターパスワード

## アンインストール

```bash
brew uninstall bwenv
```

## アップグレード

```bash
brew upgrade bwenv
```


