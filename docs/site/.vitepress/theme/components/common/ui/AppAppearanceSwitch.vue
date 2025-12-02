<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, type Ref, withDefaults, defineProps } from 'vue'
import { useData, useRoute } from 'vitepress'
import AppDropdown from './AppDropdown.vue'
import IconLoader from './IconLoader.vue'

type Appearance = 'auto' | 'light' | 'dark'

const props = withDefaults(
  defineProps<{
    /**
     * 設定パネル内などで、プルダウンではなくリストとして表示したい場合に true
     */
    inline?: boolean
  }>(),
  {
    inline: false
  }
)

const route = useRoute()
const { isDark, theme, site } = useData()

const APPEARANCE_STORAGE_KEY = 'vitepress-theme-appearance'

const appearance = ref<Appearance>('auto')
const mediaQuery: Ref<MediaQueryList | null> = ref(null)

// 言語判定（AppLanguageSwitch と同様のロジック）
const normalizedBase = computed(() => {
  const base = site.value.base || '/'
  if (base === '/' || base === '') return ''
  return base.endsWith('/') ? base.slice(0, -1) : base
})

const languagePath = computed(() => {
  const base = normalizedBase.value
  const path = route.path

  if (base && path.startsWith(base + '/')) {
    return path.slice(base.length)
  }

  if (base && path === base) {
    return '/'
  }

  return path
})

const currentLang = computed(() => {
  const path = languagePath.value
  const match = path.match(/^\/([a-zA-Z-]+)\//)
  if (match) return match[1]
  return 'en'
})

const messages = {
  ja: {
    title:
      'テーマの表示モード（ライト / ダーク / デバイスに合わせる）を切り替えます',
    auto: {
      label: 'デバイス設定',
    },
    light: {
      label: 'ライト',
    },
    dark: {
      label: 'ダーク',
    }
  },
  en: {
    title: 'Switch appearance (light / dark / follow system)',
    auto: {
      label: 'device conf',
    },
    light: {
      label: 'Light',
    },
    dark: {
      label: 'Dark',
    }
  }
} as const

const currentMessages = computed(() => {
  const lang = currentLang.value as keyof typeof messages
  return messages[lang] ?? messages.en
})

const title = computed(() => {
  return theme.value.appearanceSwitchTitle || currentMessages.value.title
})

const currentLabel = computed(() => {
  switch (appearance.value) {
    case 'light':
      return currentMessages.value.light.label
    case 'dark':
      return currentMessages.value.dark.label
    case 'auto':
    default:
      return currentMessages.value.auto.label
  }
})

const appearanceOptions = computed<
  { value: Appearance; label: string; description?: string }[]
>(() => [
  {
    value: 'auto',
    label: currentMessages.value.auto.label,
  },
  {
    value: 'light',
    label: currentMessages.value.light.label,
  },
  {
    value: 'dark',
    label: currentMessages.value.dark.label,
  }
])

function ensureMediaQuery() {
  if (typeof window === 'undefined') return
  if (!mediaQuery.value) {
    mediaQuery.value = window.matchMedia('(prefers-color-scheme: dark)')
  }
}

function applyAppearance(mode: Appearance) {
  appearance.value = mode

  if (typeof window === 'undefined') return

  ensureMediaQuery()

  if (mode === 'auto') {
    const prefersDark = mediaQuery.value?.matches ?? false
    isDark.value = prefersDark
    window.localStorage.setItem(APPEARANCE_STORAGE_KEY, 'auto')
  } else {
    isDark.value = mode === 'dark'
    window.localStorage.setItem(APPEARANCE_STORAGE_KEY, mode)
  }
}

function handleMediaChange(ev: MediaQueryListEvent) {
  if (appearance.value !== 'auto') return
  isDark.value = ev.matches
}

onMounted(() => {
  if (typeof window === 'undefined') return

  ensureMediaQuery()

  const stored = window.localStorage.getItem(
    APPEARANCE_STORAGE_KEY
  ) as Appearance | null

  const initial: Appearance =
    stored === 'light' || stored === 'dark' || stored === 'auto'
      ? stored
      : 'auto'

  applyAppearance(initial)

  mediaQuery.value?.addEventListener('change', handleMediaChange)
})

onBeforeUnmount(() => {
  mediaQuery.value?.removeEventListener('change', handleMediaChange)
})

function onSelect(mode: Appearance, close: () => void) {
  applyAppearance(mode)
  close()
}
</script>

<template>
  <!-- 通常: トップバー用のドロップダウン -->
  <AppDropdown v-if="!props.inline" class="app-appearance-switch">
    <template #trigger="{ isOpen, toggle }">
      <button
        type="button"
        class="button button-ghost app-appearance-switch__button"
        :aria-expanded="isOpen"
        aria-haspopup="listbox"
        :title="title"
        @click="toggle"
      >
        <span class="text">
          {{ currentLabel }}
        </span>
        <IconLoader
          name="icon-chevron"
          :width="12"
          :height="12"
          fill="currentColor"
          aria-hidden="true"
          class="chevron"
        />
      </button>
    </template>

    <template #menu="{ close }">
      <ul class="app-appearance-switch__menu" role="listbox">
        <li
          v-for="option in appearanceOptions"
          :key="option.value"
          class="app-appearance-switch__item"
        >
          <button
            type="button"
            class="app-appearance-switch__option"
            role="option"
            :aria-selected="appearance === option.value"
            @click="onSelect(option.value, close)"
          >
            <IconLoader name="icon-check" :width="14" :height="14" stroke-color="var(--text-main)" aria-hidden="true" v-if="appearance === option.value" class="icon-selected" />
            <span class="label">
              {{ option.label }}
            </span>
          </button>
        </li>
      </ul>
    </template>
  </AppDropdown>

  <!-- インライン: 設定パネル内用のリスト表示 -->
  <div v-else class="app-appearance-switch app-appearance-switch--inline">
    <ul class="app-appearance-switch__menu app-appearance-switch__menu--inline" role="listbox">
      <li
        v-for="option in appearanceOptions"
        :key="option.value"
        class="app-appearance-switch__item"
      >
        <button
          type="button"
          class="app-appearance-switch__option"
          role="option"
          :aria-selected="appearance === option.value"
          @click="applyAppearance(option.value)"
        >
          <IconLoader
            name="icon-check"
            :width="14"
            :height="14"
            stroke-color="var(--text-main)"
            aria-hidden="true"
            v-if="appearance === option.value"
            class="icon-selected"
          />
          <span class="label">
            {{ option.label }}
          </span>
        </button>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.app-appearance-switch {
  position: relative;
  width: 100%;
  border: 1px solid var(--stroke-option);
  border-radius: .3rem;

  .app-appearance-switch__button {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    font-size: 1.1rem;
    padding: 0.4rem 1.1rem;

    .text {
      white-space: nowrap;
    }

    .chevron {
      transition: transform 0.2s ease;
    }
  }
}

.app-appearance-switch__menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 0.25rem;
  min-width: 14rem;
  border-radius: 0.5rem;
  background: var(--bg-main);
  border: 1px solid var(--text-option);
  box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.1);
  z-index: 100;
  padding: 0.25rem 0;

  .app-appearance-switch__item {
    list-style: none;
    position: relative;
    .icon-selected {
      position: absolute;
      left: 0.95rem;
      top: 50%;
      transform: translateY(-50%);
      width: 1.2rem;
      height: 1.2rem;
      fill: var(--text-lightest);
    }
  }

  .app-appearance-switch__option {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.15rem;
    width: 100%;
    padding: 0.6rem 1.25rem;
    border: none;
    background: transparent;
    text-align: left;
    cursor: pointer;
    color: var(--text-main);
    font-size: 1.3rem;
    font-weight: 500;
    padding-left: 3rem;

    &:hover {
      background: var(--bg-accent);
      /* color: var(--text-lightest); */
    }

    .label {
      font-weight: 600;
    }

  }
}

.app-appearance-switch__menu--inline {
  position: static;
  margin-top: 0;
  min-width: 0;
  border: none;
  box-shadow: none;
  padding: 0;
}
</style>
