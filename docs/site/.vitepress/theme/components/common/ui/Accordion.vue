<script setup lang="ts">
import { ref, watch, inject, onMounted, onUnmounted, withDefaults, defineProps } from 'vue'

const props = withDefaults(
  defineProps<{
    /** 初期状態で開いているかどうか */
    open?: boolean
  }>(),
  {
    open: false
  }
)

const emit = defineEmits<{
  toggle: [isOpen: boolean]
}>()

const isOpen = ref(props.open)

// AccordionListからの制御を受け取る
const accordionList = inject<{
  registerAccordion: (id: string) => void
  updateAccordionState: (id: string, isOpen: boolean) => void
  getState: (id: string) => boolean
} | null>('accordionList', null)

// ユニークなIDを生成
const accordionId = `accordion-${Math.random().toString(36).substr(2, 9)}`

// AccordionListに登録
onMounted(() => {
  if (accordionList) {
    accordionList.registerAccordion(accordionId)
    // グローバルイベントをリッスン
    window.addEventListener('accordion-toggle-all', handleToggleAll)
  }
})

onUnmounted(() => {
  window.removeEventListener('accordion-toggle-all', handleToggleAll)
})

// 一斉展開/折りたたみのハンドラ
const handleToggleAll = (event: Event) => {
  if (event instanceof CustomEvent && accordionList) {
    const shouldOpen = event.detail.isOpen
    isOpen.value = shouldOpen
    accordionList.updateAccordionState(accordionId, shouldOpen)
  }
}

// props.openの変更を監視
watch(() => props.open, (newValue) => {
  isOpen.value = newValue
})

// 開閉を切り替え
const toggle = () => {
  isOpen.value = !isOpen.value
  emit('toggle', isOpen.value)
  if (accordionList) {
    accordionList.updateAccordionState(accordionId, isOpen.value)
  }
}
</script>

<template>
  <div class="accordion" :class="{ 'accordion--open': isOpen }">
    <button
      type="button"
      class="accordion__summary"
      :aria-expanded="isOpen"
      @click="toggle"
    >
      <slot name="title" />
    </button>
    <Transition name="accordion">
      <div v-if="isOpen" class="accordion__content">
        <slot name="content" />
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.accordion {
  border: 1px solid var(--storoke-light, var(--stroke-light, rgba(148, 163, 184, 0.4)));
  border-radius: 0.5rem;
  margin-bottom: 1rem;
  background: var(--bg-main, transparent);
}

.accordion__summary {
  width: 100%;
  display: flex;
  align-items: center;
  padding: 1rem 1.5rem;
  cursor: pointer;
  font-size: 1.4rem;
  font-weight: 600;
  color: var(--text-main, var(--text-bold));
  user-select: none;
  transition: color 0.2s ease;
  border: none;
  background: transparent;
  text-align: left;
}

.accordion__summary:hover {
  color: var(--text-bold);
}

.accordion--open .accordion__summary {
  border-bottom: 1px solid var(--storoke-light, var(--stroke-light, rgba(148, 163, 184, 0.4)));
  margin-bottom: 1rem;
}

.accordion__content {
  padding: 0 1.5rem 1.5rem 1.5rem;
  font-size: 1.3rem;
  line-height: 1.6;
  color: var(--text-main);
  overflow: hidden;
}

/* アニメーション */
.accordion-enter-active {
  transition: all 0.3s ease-out;
}

.accordion-leave-active {
  transition: all 0.2s linear;
}

.accordion-enter-from {
  max-height: 0;
  opacity: 0;
  padding-top: 0;
  padding-bottom: 0;
}

.accordion-enter-to {
  max-height: 2000px;
  opacity: 1;
}

.accordion-leave-from {
  max-height: 2000px;
  opacity: 1;
}

.accordion-leave-to {
  max-height: 0;
  opacity: 0;
  padding-top: 0;
  padding-bottom: 0;
}
</style>

