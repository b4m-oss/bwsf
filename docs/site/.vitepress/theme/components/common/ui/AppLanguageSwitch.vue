<script setup lang="ts">
import { computed, defineProps, withDefaults } from 'vue'
import { useRoute, useData, withBase } from 'vitepress'
import IconLoader from './IconLoader.vue'
import AppDropdown from './AppDropdown.vue'

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
const { site } = useData()

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
</script>

<template>
  <!-- 通常: トップバー用のドロップダウン -->
  <AppDropdown v-if="!props.inline" class="app-language-switch">
    <template #trigger="{ isOpen, toggle }">
      <button
        type="button"
        class="button button-ghost"
        :aria-expanded="isOpen"
        aria-haspopup="listbox"
        @click="toggle"
      >
        <IconLoader
          name="icon-lang"
          :width="14"
          :height="14"
          fill="currentColor"
          aria-label="言語を切り替える"
        />
        <span class="text">
          {{ currentLabel }}
        </span>
      </button>
    </template>

    <template #menu="{ close }">
      <ul class="app-language-switch__menu" role="listbox">
        <li
          v-for="locale in localeOptions"
          :key="locale.lang"
          class="app-language-switch__item"
        >
          <a
            :href="locale.href"
            class="app-language-switch__link"
            @click="close"
          >
            <IconLoader
              name="icon-check"
              :width="14"
              :height="14"
              stroke-color="var(--text-main)"
              aria-hidden="true"
              v-if="currentLang === locale.lang"
              class="icon-selected"
            />
            {{ locale.label }}
          </a>
        </li>
      </ul>
    </template>
  </AppDropdown>

  <!-- インライン: 設定パネル内用のリスト表示 -->
  <div v-else class="app-language-switch app-language-switch--inline">
    <ul class="app-language-switch__menu app-language-switch__menu--inline" role="listbox">
      <li
        v-for="locale in localeOptions"
        :key="locale.lang"
        class="app-language-switch__item"
      >
        <a
          :href="locale.href"
          class="app-language-switch__link"
        >
          <IconLoader
            name="icon-check"
            :width="14"
            :height="14"
            fill="transparent"
            strokeColor="green"
            aria-hidden="true"
            v-if="currentLang === locale.lang"
            class="icon-selected"
          />
          {{ locale.label }}
        </a>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.app-language-switch {
  position: relative;

  width: 100%;
  border: 1px solid var(--stroke-option);
  border-radius: .4rem;

  .button {
    font-size: 1.1rem;
    padding: 0.5rem 1rem;
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 0.5rem;
  }
}

.app-language-switch__menu {
  position: absolute;
  top: 100%;
  left: 0;
  margin-top: 0.25rem;
  min-width: 7rem;
  border-radius: 0.5rem;
  background: var(--bg-main);
  border: 1px solid var(--text-option);
  box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.1);
  z-index: 100;

  .app-language-switch__item {
    list-style: none;
    position: relative;
    font-size: 1.3rem;
    font-weight: 500;
    color: var(--text-option);

    .icon-selected {
      position: absolute;
      left: 0.95rem;
      top: 50%;
      transform: translateY(-50%);
      width: 1.2rem;
      height: 1.2rem;
      /* fill: var(--text-lightest); */
    }

    .app-language-switch__link {
      display: block;
      width: 100%;
      height: 100%;
      padding: 1em 2em;
      padding-left: 3rem;
      text-decoration: none;
      /* color: var(--text-option); */

      &:hover {
        background: var(--bg-accent);
      }
    }
  }
}

.app-language-switch__menu--inline {
  position: static;
  margin-top: 0;
  min-width: 0;
  border: none;
  box-shadow: none;
}
</style>


