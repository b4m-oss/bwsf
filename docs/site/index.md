---
title: bwsf
description: Manage .env files securely with Bitwarden CLI
titleTemplate: false
---

<script setup lang="ts">
import { withBase } from 'vitepress'

if (typeof window !== 'undefined' && typeof navigator !== 'undefined') {
  const rawLangs = navigator.languages && navigator.languages.length > 0
    ? navigator.languages
    : [navigator.language]

  const langs = rawLangs
    .filter((l): l is string => !!l)
    .map((l) => l.toLowerCase())

  const isJa = langs.some((l) => l.startsWith('ja'))
  const targetLocale = isJa ? 'ja' : 'en'
  const targetPath = targetLocale === 'ja' ? '/ja/' : '/en/'

  const targetUrl = withBase(targetPath)
  const target = new URL(targetUrl, window.location.origin)

  if (window.location.pathname !== target.pathname) {
    window.location.replace(target.href)
  }
}
</script>

<HomeLayout class="lang-chooser-root">

<template #hero>
  <hgroup class="title">
    <h1 class="heading-1">bwsf</h1>
    <div class="description">
      <p>Secure .env Management</p>
      <p>Choose your language / 言語を選択</p>
    </div>
  </hgroup>
  
  <nav class="hero-nav">
    <a :href="withBase('/en/')" class="button button-super">English</a>
    <a :href="withBase('/ja/')" class="button button-super">日本語</a>
  </nav>
</template>

<template #features>
</template>

</HomeLayout>

<style scoped>
.lang-chooser-root {
  opacity: 0;
  animation: lang-fade-in 0.4s ease-out 1s forwards;
}

@keyframes lang-fade-in {
  to {
    opacity: 1;
  }
}
</style>
