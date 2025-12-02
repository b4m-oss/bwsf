<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

interface Props {
  messages: string[]
  typingSpeedMs?: number
  deletingSpeedMs?: number
  stayDurationMs?: number
  delayAfterTypeMs?: number
  delayAfterDeleteMs?: number
  loop?: boolean
  idleDelayMs?: number
}

const props = withDefaults(defineProps<Props>(), {
  typingSpeedMs: 90,
  deletingSpeedMs: 50,
  stayDurationMs: 1200,
  delayAfterTypeMs: 800,
  delayAfterDeleteMs: 600,
  loop: true,
  idleDelayMs: 20000,
})

const currentMessageIndex = ref(0)
const currentText = ref('')
const isDeleting = ref(false)
const isAnimating = ref(false)
const isWaitingAfterType = ref(false)
const isWaitingAfterDelete = ref(false)

let timer: ReturnType<typeof setTimeout> | null = null
let idleTimer: ReturnType<typeof setTimeout> | null = null

const hasMessages = computed(() => props.messages && props.messages.length > 0)
const hasAnimatedOnce = ref(false)

const isCaretVisible = computed(
  () => isAnimating.value && !isWaitingAfterType.value && !isWaitingAfterDelete.value,
)

const clearTimer = () => {
  if (timer !== null) {
    clearTimeout(timer)
    timer = null
  }
}

const clearIdleTimer = () => {
  if (idleTimer !== null) {
    clearTimeout(idleTimer)
    idleTimer = null
  }
}

const scheduleTick = (delay: number) => {
  clearTimer()
  timer = setTimeout(tick, delay)
}

const tick = () => {
  if (!hasMessages.value || !isAnimating.value) {
    return
  }

  const messages = props.messages
  const message = messages[currentMessageIndex.value] ?? ''

  // タイプ完了後の待機フェーズ → 削除フェーズへ遷移
  if (isWaitingAfterType.value) {
    isWaitingAfterType.value = false
    isDeleting.value = true
    scheduleTick(props.deletingSpeedMs)
    return
  }

  // 削除完了後の待機フェーズ → 次の文言のタイプフェーズへ遷移
  if (isWaitingAfterDelete.value) {
    isWaitingAfterDelete.value = false
    isDeleting.value = false
    scheduleTick(props.typingSpeedMs)
    return
  }

  if (!isDeleting.value) {
    // タイプイン中
    if (currentText.value.length < message.length) {
      currentText.value = message.slice(0, currentText.value.length + 1)
      scheduleTick(props.typingSpeedMs)
    } else {
      // 文章を全て表示し終えたら、一定時間待ってから削除フェーズへ
      isWaitingAfterType.value = true
      scheduleTick(props.delayAfterTypeMs)
    }
  } else {
    // 削除中（バックスペース風）
    if (currentText.value.length > 0) {
      currentText.value = message.slice(0, currentText.value.length - 1)
      scheduleTick(props.deletingSpeedMs)
    } else {
      // 全削除が完了したら次のメッセージへ
      const isLastMessage = currentMessageIndex.value === messages.length - 1
      if (!props.loop && isLastMessage) {
        clearTimer()
        return
      }

      currentMessageIndex.value = isLastMessage ? 0 : currentMessageIndex.value + 1
      // 次の文言のタイプイン前に待機を挟む
      isWaitingAfterDelete.value = true
      scheduleTick(props.delayAfterDeleteMs)
    }
  }
}

const resetIdleTimer = () => {
  clearIdleTimer()
  idleTimer = setTimeout(() => {
    startAnimation()
  }, props.idleDelayMs)
}

const handleUserActivity = () => {
  stopAnimation()
  resetIdleTimer()
}

const attachIdleListeners = () => {
  if (typeof window === 'undefined') return
  window.addEventListener('mousemove', handleUserActivity)
  window.addEventListener('scroll', handleUserActivity)
  window.addEventListener('keydown', handleUserActivity)
  window.addEventListener('touchstart', handleUserActivity)
}

const detachIdleListeners = () => {
  if (typeof window === 'undefined') return
  window.removeEventListener('mousemove', handleUserActivity)
  window.removeEventListener('scroll', handleUserActivity)
  window.removeEventListener('keydown', handleUserActivity)
  window.removeEventListener('touchstart', handleUserActivity)
}

const startAnimation = () => {
  if (!hasMessages.value || isAnimating.value) {
    return
  }

  const messages = props.messages
  const firstMessage = messages[0] ?? ''

  // 初回アニメーションは、1つ目のメッセージが既にフル表示されている前提で削除から開始
  if (!hasAnimatedOnce.value && currentMessageIndex.value === 0 && currentText.value === firstMessage) {
    isDeleting.value = true
  }

  isAnimating.value = true
  hasAnimatedOnce.value = true

  const initialDelay = isDeleting.value ? props.deletingSpeedMs : props.typingSpeedMs
  scheduleTick(initialDelay)
}

const stopAnimation = () => {
  if (!isAnimating.value) return

  // タイプイン途中で止まった場合は、その文言を最後まで表示する
  if (hasMessages.value && !isDeleting.value && !isWaitingAfterType.value && !isWaitingAfterDelete.value) {
    const messages = props.messages
    const message = messages[currentMessageIndex.value] ?? ''
    if (currentText.value.length < message.length) {
      currentText.value = message
    }
  }

  isAnimating.value = false
  clearTimer()
}

onMounted(() => {
  if (!hasMessages.value) return

  // 初期ロード時は 1つ目の文言をそのまま表示しておく
  const messages = props.messages
  currentText.value = messages[0] ?? ''
  currentMessageIndex.value = 0
  isDeleting.value = false
  isAnimating.value = false

  attachIdleListeners()
  resetIdleTimer()
})

onBeforeUnmount(() => {
  clearTimer()
  clearIdleTimer()
  detachIdleListeners()
})
</script>

<template>
  <span class="caption" aria-live="polite">
    <span class="caption-text">
      {{ currentText }}
    </span>
    <span
      v-if="isCaretVisible"
      class="caption-caret"
      aria-hidden="true"
    >
      |
    </span>
  </span>
</template>

<style scoped>
.caption {
  display: inline-flex;
  align-items: center;
  gap: 0.15rem;
}

.caption-text {
  white-space: nowrap;
}

.caption-caret {
  display: inline-block;
  width: 0.1em;
  animation: caret-blink 1s step-end infinite;
  opacity: 0.7;
}

@keyframes caret-blink {
  0%,
  50% {
    opacity: 0.7;
  }
  51%,
  100% {
    opacity: 0;
  }
}
</style>


