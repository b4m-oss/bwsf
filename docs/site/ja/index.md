---
title: ホーム
description: Bitwarden CLIで.envファイルを安全に管理
titleTemplate: bwsf - .env管理をもっと手軽に
layout: home
---

<script setup lang="ts">
import { withBase } from 'vitepress'
</script>

<HomeLayout>

<template #hero>
  <hgroup class="title">
    <h1 class="heading-1">
    <IconLoader name="icon-chevron" :width="64" :height="64" fill="transparent" class="prompt-icon" />
    bwsf
    </h1>
    <div class="description">
      <p>.envファイルの管理はCLIに任せる</p>
      <p>プロジェクトも、環境も、メンバーも、すべて一元管理する</p>
      <p>オンプレミスでの利用も可能</p>
      <p>そんな夢のコマンドが、<span class="cmd">bwsf</span>です。</p>
    </div>
  </hgroup>
  
  <nav class="hero-nav">
    <a :href="withBase('/ja/guide/getting-started')" class="button button-super getting-started">今すぐスタート<IconLoader name="icon-arrow" :width="14" :height="14" :strokeColor="'#ffffff'" /></a>
    <a href="https://github.com/b4m-oss/bwsf" class="button github" taget="_blank" rel="noopener"><IconLoader name="icon-github" :width="17" :height="17" fill="#ffffff" />GitHub</a>
    <p class="caption">
      <span class="dev-by">開発：<a href="https://b4m.co.jp/" target="_blank" rel="noopener">合同会社 知的・自転車</a></span>
    </p>
    
  </nav>
</template>

<template #features>
  <HeroFeatureCard title="たった4文字の<br>コマンドライン" description="左手の指を4回動かすだけで打てる基本コマンド。オプションコマンドも、わかりやすくシンプルです。" />
  <HeroFeatureCard title="複数の開発環境に対応<br>独自名も可" description="開発（.local）、ステージング（.staging）、本番（.production）など、開発に利用されている複数の環境をまとめて管理できます。" />
  <HeroFeatureCard title="組織・メンバーで<br>活用しやすい" description="Bitwardenは、組織での利用も可能です。開発メンバーだけをdotenvsフォルダーに招待し、必要な権限を与えることができます。" />
  <HeroFeatureCard title="機密情報は外に出さない<br>自分たちで管理する" description="Bitwardenは、オープンソースです。オンプレミスで運用することもできます。" />
</template>

<style scoped>
.title {
  .heading-1 {
    position: relative;
    font-size: 7.2rem;
    font-weight: 800;
    color: var(--text-bold);
    margin-bottom: 3rem;
    padding-left: 3rem;
  }

  .prompt-icon {
    position: absolute;
    top: calc(50% + .2rem);
    transform: translateY(-50%);
    left: -4rem;
    opacity: 0;
    animation: promptBlink 1.8s ease-out 0.1s infinite;
  }

  .description {
    p {
      font-weight: 800;
      margin-bottom: 1.2em;
    }
    
    margin-bottom: 3rem;
  }

  
}

.hero-nav {
  display: grid;
  grid-template-columns: auto auto;
  grid-template-rows: 1fr 1fr;
  width: fit-content;

  column-gap: 2rem;

  .getting-started {
    grid-column: 1/2;
    grid-row: 1/2;
    display: flex;
    flex-flow: row;
    gap: 1rem;
    align-items: center;
  }

  .github {
    grid-column: 2/3;
    grid-row: 1/2;

    display: flex;
    flex-flow: row;
    gap: .5rem;
    align-items: center;
  }

  .caption {
    grid-column: 1/3;
    grid-row: 2/3;

    font-size: 1.1rem;
    padding-top: 1.2rem;

    a {
      color: var(--text-main);
    }
  }



}

</style>

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