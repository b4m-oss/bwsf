<script setup lang="ts">
import { computed, withDefaults, defineProps } from 'vue'

const props = withDefaults(
  defineProps<{
    /** トグルの状態（v-model用） */
    modelValue: boolean
    /** アクティブ時のラベル */
    activeLabel?: string
    /** 非アクティブ時のラベル */
    inactiveLabel?: string
    /** ラベルを表示するかどうか */
    showLabel?: boolean
  }>(),
  {
    activeLabel: '',
    inactiveLabel: '',
    showLabel: true
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const toggle = () => {
  emit('update:modelValue', !props.modelValue)
}

const statusText = computed(() => {
  if (!props.showLabel) return ''
  return props.modelValue ? props.activeLabel : props.inactiveLabel
})
</script>

<template>
  <div class="toggle-wrapper">
    <button
      type="button"
      :class="['toggle-button', { active: modelValue }]"
      :aria-pressed="modelValue"
      @click="toggle"
    >
      <span class="toggle-slider"></span>
    </button>
    <span v-if="showLabel && statusText" class="toggle-status">
      {{ statusText }}
    </span>
  </div>
</template>

<style scoped>
.toggle-wrapper {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.toggle-button {
  position: relative;
  width: 48px;
  height: 26px;
  border-radius: 13px;
  border: none;
  background: var(--storoke-light, var(--stroke-light, rgba(148, 163, 184, 0.4)));
  cursor: pointer;
  transition: background 0.2s ease;
  flex-shrink: 0;
}

.toggle-button:hover {
  background: var(--bg-accent, rgba(148, 163, 184, 0.6));
}

.toggle-button.active {
  background: var(--text-link, #20368d);
}

.toggle-button.active:hover {
  background: var(--text-link-hover, #2468c6);
}

.toggle-slider {
  position: absolute;
  top: 3px;
  left: 3px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: white;
  transition: transform 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
}

.toggle-button.active .toggle-slider {
  transform: translateX(22px);
}

.toggle-status {
  font-size: 0.875rem;
  color: var(--text-main, var(--text-option));
  min-width: 50px;
  user-select: none;
}
</style>

