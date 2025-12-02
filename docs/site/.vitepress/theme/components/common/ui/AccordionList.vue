<script setup lang="ts">
import { ref, provide, computed } from 'vue'
import ToggleButton from './ToggleButton.vue'

// 子アコーディオンの開閉状態を管理
const accordionStates = ref<Map<string, boolean>>(new Map())
const accordionRefs = ref<Set<string>>(new Set())

// すべてのアコーディオンが開いているかどうか
const allOpen = computed(() => {
  if (accordionRefs.value.size === 0) return false
  return Array.from(accordionStates.value.values()).every(state => state === true)
})

// すべてのアコーディオンが閉じているかどうか
const allClosed = computed(() => {
  if (accordionRefs.value.size === 0) return true
  return Array.from(accordionStates.value.values()).every(state => state === false)
})

// アコーディオンを登録
const registerAccordion = (id: string) => {
  accordionRefs.value.add(id)
  accordionStates.value.set(id, false)
}

// アコーディオンの開閉状態を更新
const updateAccordionState = (id: string, isOpen: boolean) => {
  accordionStates.value.set(id, isOpen)
}

// すべて展開/折りたたみ
const toggleAll = (value: boolean) => {
  accordionRefs.value.forEach(id => {
    accordionStates.value.set(id, value)
  })
  // 子コンポーネントに状態変更を通知
  emitToggleAll(value)
}

// provideで子コンポーネントに機能を提供
provide('accordionList', {
  registerAccordion,
  updateAccordionState,
  getState: (id: string) => accordionStates.value.get(id) ?? false
})

// すべて展開/折りたたみのイベントを発火
const emitToggleAll = (shouldOpen: boolean) => {
  const event = new CustomEvent('accordion-toggle-all', {
    detail: { isOpen: shouldOpen }
  })
  window.dispatchEvent(event)
}
</script>

<template>
  <div class="accordion-list">
    <div class="accordion-list__header">
      <div class="accordion-list__toggle-wrapper">
        <span class="accordion-list__toggle-label">
          {{ allOpen ? 'すべて折りたたむ' : 'すべて展開する' }}
        </span>
        <ToggleButton
          :model-value="allOpen"
          @update:model-value="toggleAll"
          :active-label="''"
          :inactive-label="''"
          :show-label="false"
        />
      </div>
    </div>
    <div class="accordion-list__content">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.accordion-list {
  margin-bottom: 2rem;
}

.accordion-list__header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 1rem;
}

.accordion-list__toggle-wrapper {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.accordion-list__toggle-label {
  font-size: 1.2rem;
  font-weight: 500;
  color: var(--text-main, var(--text-bold));
  user-select: none;
}

.accordion-list__content {
  display: flex;
  flex-direction: column;
}
</style>

