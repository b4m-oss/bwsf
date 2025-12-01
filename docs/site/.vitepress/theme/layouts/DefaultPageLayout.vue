<script setup lang="ts">
// 一般的なページ用のレイアウトサンプルです。
// 例:
//   <DefaultPageLayout>
//     # Markdown 見出し
//     ページ本文...
//   </DefaultPageLayout>
//
// といった形で、docs 全体に共通のレイアウトを持たせることができます。
</script>

<template>
  <div class="default-page-layout">
    <header class="default-page-layout__header">
      <slot name="header" />
    </header>

    <main class="default-page-layout__content">
      <slot />
    </main>

    <footer class="default-page-layout__footer">
      <slot name="footer" />
    </footer>
  </div>
</template>

<style scoped>
.default-page-layout {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.default-page-layout__content {
  max-width: 768px;
  margin: 0 auto;
}
</style>


