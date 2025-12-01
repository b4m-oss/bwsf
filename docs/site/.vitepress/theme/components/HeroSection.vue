<script setup lang="ts">
// トップページなどで再利用できる Hero コンポーネントのサンプルです。
// 例:
//   <HeroSection
//     title="bwsf"
//     subtitle="Bitwarden CLI で .env を安全に管理"
//   />
//
// props はお好みで増やして構いません。

const props = defineProps<{
  title: string
  subtitle?: string
  label?: string
}>()
</script>

<template>
  <section class="hero">
    <p v-if="label" class="hero__label">
      {{ label }}
    </p>
    <h1 class="hero__title">
      {{ title }}
    </h1>
    <p v-if="subtitle" class="hero__subtitle">
      {{ subtitle }}
    </p>
    <div class="hero__actions">
      <slot name="actions" />
    </div>
  </section>
</template>

<style scoped>
.hero {
  text-align: center;
  padding: 3rem 1.5rem 2rem;
}

.hero__label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--vp-c-text-2);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.hero__title {
  font-size: clamp(2.5rem, 4vw, 3.25rem);
  line-height: 1.1;
  margin: 0.5rem 0 1rem;
}

.hero__subtitle {
  font-size: 1.05rem;
  color: var(--vp-c-text-2);
  max-width: 40rem;
  margin: 0 auto;
}

.hero__actions {
  margin-top: 1.75rem;
  display: flex;
  justify-content: center;
  gap: 1rem;
  flex-wrap: wrap;
}
</style>


