---
title: ホーム
description: Bitwarden CLIで.envファイルを安全に管理
titleTemplate: bwsf - .env管理をもっと手軽に
layout: home
---

<HomeLayout>

## クイックスタート

```bash
# Homebrew でインストール
brew tap b4m-oss/tap && brew install bwsf

# 初期設定
bwsf setup

# Bitwarden から .env をプル
cd /path/to/your_project
bwsf pull

# Bitwarden に .env をプッシュ
bwsf push
```

## 仕組み

bwsf は公式の Bitwarden CLI（`bw`）を使用して、`.env` ファイルを安全に保存・取得します。環境変数は Bitwarden ボールト内の専用 `dotenvs` フォルダに**ノートアイテム**として保存されます。

各プロジェクトの `.env` ファイルはディレクトリ名で識別されるため、複数のプロジェクトを簡単に整理・管理できます。

</HomeLayout>