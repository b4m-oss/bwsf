<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useData, withBase } from 'vitepress'

const route = useRoute()
const { site } = useData()

const isOpen = ref(false)

// site.base（例: "/bwsf/"）を正規化して、先頭/末尾のスラッシュ重複を避ける
const normalizedBase = computed(() => {
  const base = site.value.base || '/'
  if (base === '/' || base === '') return ''
  return base.endsWith('/') ? base.slice(0, -1) : base
})

// ルーティング上の「言語付きパス」（例: "/ja/guide/..."）だけを取り出す
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

// locales 設定から候補一覧を作成
const localeOptions = computed(() => {
  const locales = site.value.locales ?? {}
  const entries = Object.entries(locales) as [string, any][]

  return entries.map(([key, locale]) => {
    const langCode = key // "en", "ja", ...
    const label = locale.label ?? langCode

    const path = languagePath.value
    let langPath: string

    const match = path.match(/^\/([a-zA-Z-]+)\//)
    if (match) {
      // 先頭の言語セグメントを置き換える
      langPath = path.replace(/^\/([a-zA-Z-]+)\//, `/${langCode}/`)
    } else {
      // 何もついていなければ、その言語のルートへ
      langPath = `/${langCode}/`
    }

    return {
      lang: langCode,
      label,
      href: withBase(langPath)
    }
  })
})

const currentLabel = computed(() => {
  const current = localeOptions.value.find((opt) => opt.lang === currentLang.value)
  return current?.label ?? currentLang.value
})

function toggle() {
  isOpen.value = !isOpen.value
}
</script>

<template>
  <div class="app-language-switch">
    <button
      type="button"
      class="app-language-switch__button"
      :aria-expanded="isOpen"
      aria-haspopup="listbox"
      @click="toggle"
    >
      <span class="icon" aria-hidden="true" />
      <span class="text">
        {{ currentLabel }}
      </span>
    </button>

    <ul
      v-if="isOpen"
      class="app-language-switch__menu"
      role="listbox"
    >
      <li
        v-for="locale in localeOptions"
        :key="locale.lang"
        class="app-language-switch__item"
      >
        <a
          :href="locale.href"
          class="app-language-switch__link"
        >
          {{ locale.label }}
        </a>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.app-language-switch {
  position: relative;
}

.app-language-switch__button {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.15rem 0.6rem;
  border-radius: 999px;
  border: 1px solid var(--vp-c-border);
  font-size: 0.8rem;
  color: var(--vp-c-text-2);
  background: var(--vp-c-bg-soft);
  cursor: pointer;
  transition: background-color 0.15s ease, border-color 0.15s ease,
    color 0.15s ease;
}

.app-language-switch__button:hover {
  background: var(--vp-c-bg-soft-up);
  border-color: var(--vp-c-brand-1);
  color: var(--vp-c-text-1);
}

.icon {
  width: 14px;
  height: 14px;
  mask: var(--vp-icon-languages) no-repeat center / contain;
  background-color: currentColor;
}

.text {
  white-space: nowrap;
}

.app-language-switch__menu {
  position: absolute;
  right: 0;
  margin-top: 0.25rem;
  padding: 0.25rem;
  min-width: 7rem;
  border-radius: 0.5rem;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  z-index: 20;
}

.app-language-switch__item {
  list-style: none;
}

.app-language-switch__link {
  display: block;
  padding: 0.25rem 0.5rem;
  font-size: 0.8rem;
  text-decoration: none;
  border-radius: 0.35rem;
}

.app-language-switch__link:hover {
  background: var(--bg-option);
}
</style>


