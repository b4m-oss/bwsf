<script setup lang="ts">
import IconLoader from '../common/ui/IconLoader.vue'

const props = defineProps<{
  icon?: string
  title: string
  description: string
  width?: number | string
  height?: number | string
  name?: string
  fill?: string
  strokeColor?: string
  strokeWidth?: number | string
  flip?: string | boolean
  rotate?: number | string
  posX?: number
  posY?: number
}>()
</script>

<template>
  <section class="feature-item">
    <div class="feature-item__icon">
      <IconLoader
        :width="props.width"
        :height="props.height"
        :name="props.name || props.icon"
        :fill="props.fill"
        :stroke-color="props.strokeColor"
        :stroke-width="props.strokeWidth"
        :flip="props.flip"
        :rotate="props.rotate"
        :style="props.posX !== undefined || props.posY !== undefined
          ? {
              position: 'absolute',
              left: `${props.posX ?? 0}px`,
              top: `${props.posY ?? 0}px`,
            }
          : {}"
      />
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
  position: relative;
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