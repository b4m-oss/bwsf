# bwsf

bwsf（Bitwarden Secured Files）は、[Bitwarden](https://bitwarden.com/)を使用して.envファイルを管理するCLIツールです。

[English version is here.](./README.md)

## 🚨🚨 破壊的変更 🚨🚨

### CLI名の変更

v0.11.0から、`bwenv`は`bwsf`に名前が変更されました。これは、既にbwenvコマンドが存在していたためです。混乱を避けるため、CLI名を変更することにしました。

#### 移行方法

設定ディレクトリの名前を変更してください。

```bash
mv ~/.config/bwenv ~/.config/bwsf
```

現在のバージョンをアンインストールし、最新バージョンを再インストールしてください。

```bash
brew uninstall bwenv
brew install bwsf
```

### 複数の`.env.environment`ファイル

v0.9.0から、bwsfは`.env | .env.staging | .env.production`のような複数の環境用.envファイルを保存できるようになりました。

これに伴い、BitwardenのNoteアイテムに保存されるデータ構造が変更されました。

v0.8.0以前に保存されたデータは、v0.9.0以降では互換性がありません。

移行システムは提供しません。

## 🚨🚨 破壊的変更 🚨🚨

v0.9.0から、bwenvは`.env | .env.staging | .env.production`のような複数環境の.envファイルを保存するようになりました。
これに伴い、BitwardenのNoteアイテムに保存されるデータ構造が変更されました。
v0.8.0以前に保存されたデータは、v0.9.0以降との互換性がありません。
マイグレーションシステムは提供されません。

## 概要

bwsfコマンドは、Bitwardenで管理されているdotenvファイルをサポートします。

簡単な使用方法は以下の通りです：

| コマンド | |
|----|----|
| bwsf push | .envファイルをBitwardenホストにプッシュ |
| bwsf pull | Bitwardenホストから.envファイルをプル |
| bwsf list | Bitwardenホストに保存されている.envファイルの一覧を表示 |

## 動機

私たちは長い間、Bitwardenをパスワードマネージャーとして使用しています。
また、.envファイルをBitwardenホストに保存し、シェルスクリプトとして管理しています。
このプロジェクトは、手作りのシェルスクリプトをGoで書かれたモダンなCLIコマンドに移行するものです。

## 要件

**`bw`** コマンドがマシンにインストールされている必要があります。

[bwコマンドのインストール方法については、こちらのドキュメントをお読みください。](https://bitwarden.com/help/cli/#download-and-install)

** Homebrew **: インストールに必要です。

### 対応OS

- macOS
- Linux
- [計画中] Windows

## インストール

| OS | コマンド |
|----|----|
| macOS / Linux| brew tap b4m-oss/tap && brew install bwenv |

## インストールの確認

```shell
bwenv -v
# bwenv version 0.5.5
```

## 使い方

### 初期セットアップ

```shell
bwenv setup
```

Bitwardenホストとアカウント情報を設定します。

### Bitwardenホストから.envファイルをプル

```shell
cd /path/to/your_project
bwenv pull
```

bwenvはカレントディレクトリの名前を使用して、Bitwardenホスト内の.envデータを検索します。
存在する場合、カレントディレクトリに.envファイルとしてデータをプルします。
カレントディレクトリに既に.envファイルがある場合、bwenvは上書きするかどうかを確認します。
データはBitwardenのNoteアイテムとして保存されます。

### Bitwardenホストに.envファイルをプッシュ

bwenvはカレントディレクトリの.envデータをBitwardenホストにプッシュします。
dotenvフォルダに同じ名前のBitwardenのNoteアイテムが存在する場合、bwenvは上書きするかどうかを確認します。

### Bitwardenホストの.envデータ一覧

```shell
bwenv list
```

Bitwardenホストから.envデータの一覧を取得します。
プロジェクト名のリストが標準出力に表示されます。

## アンインストール

```shell
brew uninstall bwenv
```

## FAQ

<details>
<summary>Q. Bitwardenアカウントを持っていません。</summary>

bwenvを使用するには、Bitwardenアカウントが必要です。

[Bitwarden Cloud](https://bitwarden.com/)にアクセスして、アカウントを登録できます。

無料で、クレジットカードも不要です。

</details>

<details>
<summary>Q. Bitwardenのセルフホストユーザーです。</summary>

もちろん、bwenvはBitwardenのセルフホストユーザーでも利用可能です。

初期セットアップ時にセルフホストのURLを入力できます。

</details>

<details>
<summary>Q. .envファイルはBitwardenホストにどのように保存されますか？</summary>

.envファイルはJSON形式に変換されます。bwenvはBitwardenのNoteアイテムを作成し、NoteセクションにそのJSONを保存します。

</details>

<details>
<summary>Q. Bitwardenのアカウント情報はどこに保存されますか？</summary>

bwenvは設定データを`~/.config/bwenv/`に保存します。

ただし、セキュリティ情報（パスワードなど）は一切保存されません。

</details>

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
