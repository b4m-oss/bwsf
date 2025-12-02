<script setup lang="ts">
// NOTE: 型定義上は存在しない仮想モジュールを扱うため、いくつか any キャストを行っています。
import { ref, shallowRef, watch, computed, createApp } from 'vue'
import MiniSearch, { type SearchResult } from 'minisearch'
// @ts-expect-error: VitePress が提供する仮想モジュール
import localSearchIndex from '@localSearchIndex'
import { dataSymbol, useData, useRouter } from 'vitepress'
// @ts-expect-error: 内部ユーティリティを直接参照
import { pathToFile } from 'vitepress/dist/client/app/utils'

interface LocalSearchDoc {
  id: string
  title: string
  titles: string[]
  text?: string
  snippet?: string
  snippetHtml?: string
}

const vpData = useData()
const { localeIndex } = vpData
const router = useRouter()

const searchIndex = shallowRef<MiniSearch<LocalSearchDoc> | null>(null)
const query = ref('')
const results = ref<(SearchResult & LocalSearchDoc)[]>([])
const isOpen = ref(false)
let searchToken = 0

// ローカルインデックスのロード（ロケールごとに切り替え）
watch(
  () => localeIndex.value,
  async (locale) => {
    const loader = (localSearchIndex as any)[locale]
    if (!loader) {
      searchIndex.value = null
      return
    }
    const mod = await loader()
    searchIndex.value = MiniSearch.loadJSON<LocalSearchDoc>(mod.default, {
      fields: ['title', 'titles', 'text'],
      storeFields: ['title', 'titles', 'text'],
      searchOptions: {
        fuzzy: 0.2,
        prefix: true,
        boost: { title: 4, text: 2, titles: 1 }
      }
    })
  },
  { immediate: true }
)

// クエリ変更時に検索 ＋ スニペット生成（非同期）
watch(query, async (q) => {
  const index = searchIndex.value
  const token = ++searchToken

  if (!index) {
    results.value = []
    return
  }

  const trimmed = q.trim()
  if (!trimmed) {
    results.value = []
    return
  }

  const searched = index.search(trimmed).slice(0, 16) as (SearchResult &
    LocalSearchDoc)[]

  const enriched: (SearchResult & LocalSearchDoc)[] = []

  for (const item of searched) {
    const fullText = await fetchPageText(item.id)
    if (token !== searchToken) {
      // 途中で別の検索が走った場合は結果を破棄
      return
    }
    const { plain, html } = createSnippet(fullText, trimmed)

    enriched.push({
      ...item,
      snippet: plain,
      snippetHtml: html
    })
  }

  if (token === searchToken) {
    results.value = enriched
  }
})

const placeholder = computed(() =>
  localeIndex.value === 'ja' ? 'ドキュメントを検索' : 'Search docs'
)

function onFocus() {
  isOpen.value = true
}

function onBlur(e: FocusEvent) {
  const next = e.relatedTarget as HTMLElement | null
  if (!next || !next.closest('.header-search-dropdown')) {
    isOpen.value = false
  }
}

function navigateTo(result: SearchResult & LocalSearchDoc) {
  isOpen.value = false
  query.value = ''
  router.go(result.id)
}

function onKeydownEsc() {
  isOpen.value = false
}

async function fetchPageText(id: string): Promise<string> {
  const [path] = id.split('#')
  const file = pathToFile(path)
  if (!file) return ''

  try {
    const mod: any = await import(/* @vite-ignore */ file)
    const comp = mod.default ?? mod
    if (!comp || (!comp.render && !comp.setup)) return ''

    const app = createApp(comp)
    // ドキュメント中の未登録コンポーネントの warning を抑制
    app.config.warnHandler = () => {}
    app.provide(dataSymbol, vpData)

    const div = document.createElement('div')
    app.mount(div)
    const text = div.textContent || ''
    app.unmount()
    return text
  } catch (e) {
    console.error('[SearchBox] failed to load page for excerpt', e)
    return ''
  }
}

function createSnippet(
  text: string,
  query: string
): { plain: string; html: string } {
  const normalized = text.replace(/\s+/g, ' ').trim()
  if (!normalized) return { plain: '', html: '' }

  const q = query.toLowerCase().split(/\s+/)[0] || ''
  const lower = normalized.toLowerCase()

  let index = q ? lower.indexOf(q) : -1
  if (index === -1) {
    index = 0
  }

  const radius = 60
  let start = Math.max(0, index - radius)
  let end = Math.min(normalized.length, index + q.length + radius)

  let snippet = normalized.slice(start, end)
  if (start > 0) snippet = '…' + snippet
  if (end < normalized.length) snippet += '…'

  const plain = snippet
  const terms = query
    .split(/\s+/)
    .map((s) => s.trim())
    .filter(Boolean)

  const html = highlightTerms(plain, terms)

  return { plain, html }
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function escapeRegExp(str: string): string {
  return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

function highlightTerms(text: string, terms: string[]): string {
  const safe = escapeHtml(text)
  if (!terms.length) return safe

  const pattern = terms.map(escapeRegExp).join('|')
  const re = new RegExp(`(${pattern})`, 'gi')
  return safe.replace(re, '<mark>$1</mark>')
}
</script>

<template>
  <div class="header-search" @keydown.esc.stop.prevent="onKeydownEsc">
    <input
      v-model="query"
      type="search"
      class="header-search-input"
      :placeholder="placeholder"
      @focus="onFocus"
      @blur="onBlur"
    >

    <div
      v-if="isOpen && results.length"
      class="header-search-dropdown"
    >
      <ul>
        <li
          v-for="item in results"
          :key="item.id"
        >
          <button
            type="button"
            class="header-search-item"
            @mousedown.prevent="navigateTo(item)"
          >
            <span class="title">
              {{ item.titles?.[item.titles.length - 1] || item.title }}
            </span>
            <span class="page-title">{{ item.title }}</span>
            <span
              v-if="item.snippetHtml"
              class="snippet"
              v-html="item.snippetHtml"
            />
            <span class="path">{{ item.id }}</span>
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.header-search {
  position: relative;
  display: flex;
  align-items: center;
}

.header-search-input {
  width: 18rem;
  max-width: 100%;
  padding: 0.45rem 0.9rem;
  border-radius: 999px;
  border: 1px solid var(--stroke-light);
  background-color: var(--bg-main);
  color: var(--text-main);
  font-size: 1.2rem;
}

.header-search-input::placeholder {
  color: var(--text-option);
}

.header-search-dropdown {
  position: absolute;
  top: 120%;
  right: 0;
  width: 34rem;
  max-width: 70vw;
  max-height: 24rem;
  padding: 0.75rem 0;
  border-radius: 0.9rem;
  border: 1px solid var(--stroke-light);
  background-color: var(--bg-main);
  box-shadow: 0 18px 45px rgba(15, 23, 42, 0.45);
  overflow-y: auto;
  z-index: 200;
}

.header-search-dropdown ul {
  list-style: none;
  margin: 0;
  padding: 0;
}

.header-search-item {
  width: 100%;
  padding: 0.65rem 1.1rem;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.2rem;
  background: none;
  border: none;
  text-align: left;
  cursor: pointer;
}

.header-search-item:hover {
  background-color: var(--bg-highlight);
}

.header-search-item .title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--text-bold);
}

.header-search-item .page-title {
  font-size: 1.1rem;
  font-weight: 500;
  color: var(--text-option);
}

.header-search-item .path {
  font-size: 1rem;
  color: var(--text-option);
}

.header-search-item .snippet {
  font-size: 1.05rem;
  color: var(--text-main);
}

.header-search-item mark {
  background-color: rgba(248, 250, 252, 0.12);
  color: var(--text-bold);
  padding: 0 0.1em;
  border-radius: 0.2em;
}

@media (max-width: 860px) {
  .header-search-input {
    width: 14rem;
    font-size: 1.1rem;
  }

  .header-search-dropdown {
    width: min(28rem, 90vw);
  }
}
</style>


