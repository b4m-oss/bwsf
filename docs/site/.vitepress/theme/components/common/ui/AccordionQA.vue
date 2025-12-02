<script setup lang="ts">
import { ref, watch, inject, onMounted, onUnmounted, withDefaults, defineProps } from 'vue'
import IconLoader from './IconLoader.vue'

const props = withDefaults(
  defineProps<{
    /** 初期状態で開いているかどうか */
    open?: boolean
    /** 質問のタイトル */
    title?: string
  }>(),
  {
    open: false,
    title: ''
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
  <div class="accordion-qa" :class="{ 'accordion-qa--open': isOpen }">
    <button
      type="button"
      class="accordion-qa__summary"
      :aria-expanded="isOpen"
      @click="toggle"
    >
      <div class="accordion-qa__title-wrapper">
        <div class="accordion-qa__icon-wrapper accordion-qa__icon-wrapper--question">
          <IconLoader name="icon-question" :width="24" :height="24" />
        </div>
        <div class="accordion-qa__title-content">
          {{ title }}
        </div>
      </div>
    </button>
    <Transition name="accordion">
      <div v-if="isOpen" class="accordion-qa__content">
        <div class="accordion-qa__content-wrapper">
          <div class="accordion-qa__icon-wrapper accordion-qa__icon-wrapper--answer">
            <IconLoader name="icon-answer" :width="24" :height="24" />
          </div>
          <div class="accordion-qa__content-text">
            <slot />
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.accordion-qa {
  border: 1px solid var(--stroke-light);
  border-radius: 0.5rem;
  margin-bottom: 1.5rem;
  background: var(--bg-main, transparent);
}

.accordion-qa__title-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  width: 100%;
  align-items: center;
}

.accordion-qa__icon-wrapper {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.4rem;
  height: 2.4rem;
  margin-top: 0.2rem;
}

.accordion-qa__icon-wrapper--question {
  color: var(--text-link, #20368d);
}

.accordion-qa__icon-wrapper--answer {
  color: var(--text-link, #20368d);
}

.accordion-qa__title-content {
  flex: 1;
  font-size: 1.4rem;
  font-weight: 600;
  color: var(--text-main, var(--text-bold));
  line-height: 1.5;
}

.accordion-qa__content-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  width: 100%;
}

.accordion-qa__summary {
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

.accordion-qa__summary:hover {
  color: var(--text-bold);
}

.accordion-qa--open .accordion-qa__summary {
  border-bottom: 1px solid var(--storoke-light, var(--stroke-light, rgba(148, 163, 184, 0.4)));
  margin-bottom: 1rem;
}

.accordion-qa__content {
  padding: 0 1.5rem 1.5rem 1.5rem;
  font-size: 1.3rem;
  line-height: 1.6;
  color: var(--text-main);
  overflow: hidden;
  padding-top: 1.5rem;
  padding-bottom: 2.5rem;
}

.accordion-qa__content-text {
  flex: 1;
  font-size: 1.3rem;
  line-height: 1.6;
  color: var(--text-main);
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

