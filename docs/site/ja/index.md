---
layout: home

hero:
  name: bwenv
  text: 安全な .env 管理
  tagline: Bitwarden CLI で .env ファイルを管理
  actions:
    - theme: brand
      text: はじめる
      link: /ja/guide/getting-started
    - theme: alt
      text: GitHub で見る
      link: https://github.com/b4m-oss/bwenv

features:
  - icon: 🔐
    title: 安全なストレージ
    details: .env ファイルを Bitwarden のボールトに安全に保存。共有ドライブに平文で秘密情報を置く必要はありません。
  - icon: 🔄
    title: 簡単同期
    details: シンプルなコマンドで、ローカルマシンと Bitwarden 間で .env ファイルをプッシュ・プル。
  - icon: 📋
    title: マルチ環境
    details: 1つのプロジェクトで複数の環境ファイル（.env、.env.staging、.env.production）を管理。
  - icon: 🖥️
    title: クロスプラットフォーム
    details: macOS と Linux に対応。Windows サポートは計画中です。
---

## クイックスタート

```bash
# Homebrew でインストール
brew tap b4m-oss/tap && brew install bwenv

# 初期設定
bwenv setup

# Bitwarden から .env をプル
cd /path/to/your_project
bwenv pull

# Bitwarden に .env をプッシュ
bwenv push
```

## 仕組み

bwenv は公式の Bitwarden CLI（`bw`）を使用して、`.env` ファイルを安全に保存・取得します。環境変数は Bitwarden ボールト内の専用 `dotenvs` フォルダに**ノートアイテム**として保存されます。

各プロジェクトの `.env` ファイルはディレクトリ名で識別されるため、複数のプロジェクトを簡単に整理・管理できます。


