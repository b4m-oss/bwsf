<script setup lang="ts">
const props = defineProps<{
  icon?: string
  title: string
  description: string
}>()
</script>

<template>
  <section class="feature-item">
    <div class="feature-item__icon">
      <!-- {{ icon }} -->
        <p>🔐</p>
    </div>
    <h2 class="feature-item__title" v-html="title"></h2>
    <div class="feature-item__description">{{ description }}</div>
  </section>
</template>

<style scoped>
.feature-item {
  display: grid;
  grid-template-columns: 20% auto;
  grid-template-rows: auto auto;
  column-gap: 1rem;
  row-gap: 0; /* タイトル直下の余白を詰める */
  padding: 1.5rem;
  border-radius: 1rem;
  border: 1px solid var(--vp-c-divider);
  background: transparent;
}

.feature-item__title {
  margin-top: 0;
  border-top: 0;
  padding-top: 0;
  width: 100%;
  font-size: 1.8rem;
  line-height: 1.2;
  font-weight: 600;
  margin-bottom: 0 !important;
  color: var(--text-bold);
}

.feature-item__icon {
  grid-column: 1 / 2;
  grid-row: 1 / 2;
  text-align: center;
}

.feature-item__title {
  grid-column: 2 / 3;
  grid-row: 1 / 2;
}

.feature-item__description {
  grid-column: 2 / 3;
  grid-row: 2 / 3;
  font-size: 1.3rem;
  color: var(--main-text-option);
  line-height: 1.4;
}
</style>