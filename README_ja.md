# bwenv

bwenvは、[Bitwarden](https://bitwarden.com/)を使用して.envファイルを管理するCLIツールです。

## 概要

bwenvコマンドは、Bitwardenで管理されているdotenvファイルをサポートします。

簡単な使用方法は以下の通りです：

| コマンド | |
|----|----|
| bwenv push | .envファイルをBitwardenホストにプッシュ |
| bwenv pull | Bitwardenホストから.envファイルをプル |
| bwenv list | Bitwardenホストに保存されている.envファイルの一覧を表示 |

## 動機

私たちは長い間、Bitwardenをパスワードマネージャーとして使用しています。
また、.envファイルをBitwardenホストに保存し、シェルスクリプトとして管理しています。
このプロジェクトは、手作りのシェルスクリプトをGoで書かれたモダンなCLIコマンドに移行するものです。

## 要件

**`bw`** コマンドがマシンにインストールされている必要があります。

[bwコマンドのインストール方法については、こちらのドキュメントをお読みください。](https://bitwarden.com/help/cli/#download-and-install)

### 対応OS

- [計画中] macOS
- [計画中] Linux
- [計画中] Windows

## インストール

[!Note]
これは計画中です。

- **macOS & Linux**: Homebrew
- **Windows**: Chocolaty

## 開発

### 要件

**Docker** が開発マシンにインストールされている必要があります。

### 開発環境の起動

```
git clone https://github.com/b4m-oss/bwenv.git
cd bwenv
make run
```

## ライセンス

[MIT License](./LICENSE)