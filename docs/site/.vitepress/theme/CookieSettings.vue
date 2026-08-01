<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vitepress'

const route = useRoute()
const analyticsEnabled = ref(true)
const saved = ref(false)

const isJapanese = computed(() => route.path.startsWith('/ja/') || route.path === '/ja')

const texts = computed(() => {
  if (isJapanese.value) {
    return {
      title: 'Cookie設定',
      analyticsLabel: 'Google Analytics（サイト分析）',
      analyticsDescription: 'サイトの利用状況を分析し、改善に役立てるために使用します。',
      enabled: '有効',
      disabled: '無効',
      save: '設定を保存',
      saved: '保存しました',
      note: '※ 設定を反映するにはページを再読み込みしてください。'
    }
  }
  return {
    title: 'Cookie Settings',
    analyticsLabel: 'Google Analytics (Site Analytics)',
    analyticsDescription: 'Used to analyze site usage and help us improve.',
    enabled: 'Enabled',
    disabled: 'Disabled',
    save: 'Save Settings',
    saved: 'Saved',
    note: '※ Please reload the page to apply the settings.'
  }
})

onMounted(() => {
  const consent = localStorage.getItem('cookie-consent')
  analyticsEnabled.value = consent !== 'declined'
})

function saveSettings() {
  localStorage.setItem('cookie-consent', analyticsEnabled.value ? 'accepted' : 'declined')
  saved.value = true
  setTimeout(() => {
    saved.value = false
  }, 2000)
}
</script>

<template>
  <div class="cookie-settings">
    <h2>{{ texts.title }}</h2>
    
    <div class="setting-item">
      <div class="setting-header">
        <label class="setting-label">{{ texts.analyticsLabel }}</label>
        <div class="toggle-wrapper">
          <button 
            @click="analyticsEnabled = !analyticsEnabled"
            :class="['toggle-button', { active: analyticsEnabled }]"
            :aria-pressed="analyticsEnabled"
          >
            <span class="toggle-slider"></span>
          </button>
          <span class="toggle-status">{{ analyticsEnabled ? texts.enabled : texts.disabled }}</span>
        </div>
      </div>
      <p class="setting-description">{{ texts.analyticsDescription }}</p>
    </div>

    <div class="setting-actions">
      <button @click="saveSettings" class="save-button">
        {{ saved ? texts.saved : texts.save }}
      </button>
    </div>
    
    <p class="setting-note">{{ texts.note }}</p>
  </div>
</template>

<style scoped>
.cookie-settings {
  margin: 2rem 0;
  padding: 1.5rem;
  border: 1px solid var(--vp-c-border);
  border-radius: 8px;
  background: var(--vp-c-bg-soft);
}

.cookie-settings h2 {
  margin: 0 0 1.5rem 0;
  padding: 0;
  border: none;
  font-size: 1.25rem;
}

.setting-item {
  padding: 1rem 0;
  border-bottom: 1px solid var(--vp-c-border);
}

.setting-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.setting-label {
  font-weight: 600;
  color: var(--vp-c-text-1);
}

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
  background: var(--vp-c-border);
  cursor: pointer;
  transition: background 0.2s ease;
}

.toggle-button.active {
  background: var(--vp-c-brand-1);
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
  color: var(--vp-c-text-2);
  min-width: 50px;
}

.setting-description {
  margin: 0.5rem 0 0 0;
  font-size: 0.875rem;
  color: var(--vp-c-text-2);
}

.setting-actions {
  margin-top: 1.5rem;
}

.save-button {
  padding: 0.625rem 1.25rem;
  border-radius: 6px;
  border: none;
  background: var(--vp-c-brand-1);
  color: white;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s ease;
}

.save-button:hover {
  background: var(--vp-c-brand-2);
}

.setting-note {
  margin: 1rem 0 0 0;
  font-size: 0.8rem;
  color: var(--vp-c-text-3);
}
</style>


