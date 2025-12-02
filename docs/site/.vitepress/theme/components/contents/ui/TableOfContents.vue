<script setup lang="ts">
import { computed, onMounted, nextTick, ref } from 'vue'
import { useData } from 'vitepress'

const { page, theme } = useData()

// 見出しデータを取得（VitePressの標準的な方法）
const headers = computed(() => {
  // page.value.headersが存在する場合はそれを使用
  if (page.value.headers && Array.isArray(page.value.headers) && page.value.headers.length > 0) {
    return page.value.headers
  }
  return []
})

// DOMから見出しを取得するフォールバック
const domHeaders = ref<any[]>([])

// DOMから見出しを取得する関数
const extractHeadersFromDOM = () => {
  // main要素内の見出しを取得
  const mainElement = document.querySelector('main') || document.querySelector('#main')
  const container = mainElement?.querySelector('.content-container') || mainElement || document
  
  const headingElements = container.querySelectorAll('h1, h2, h3, h4, h5, h6')
  const extractedHeaders: any[] = []
  
  headingElements.forEach((el, index) => {
    // 既存のIDを取得、なければ生成
    let id = el.id
    if (!id) {
      // テキストからスラッグを生成
      const text = el.textContent?.trim() || ''
      id = text
        .toLowerCase()
        .replace(/[^\w\s-]/g, '')
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-')
        .replace(/^-|-$/g, '') || `heading-${index}`
      el.id = id
    }
    
    const level = parseInt(el.tagName.charAt(1))
    const title = el.textContent?.trim() || ''
    
    // パーマリンクのリンクテキストを除外
    const cleanTitle = title.replace(/\s*Permalink.*$/, '').trim()
    
    if (cleanTitle) {
      extractedHeaders.push({
        level,
        title: cleanTitle,
        slug: id
      })
    }
  })
  
  return extractedHeaders
}

onMounted(() => {
  // nextTickでDOMの更新を待つ
  nextTick(() => {
    // page.value.headersが空の場合、DOMから直接取得
    if (headers.value.length === 0) {
      // 少し遅延を入れて、VitePressのコンテンツが完全にレンダリングされるのを待つ
      setTimeout(() => {
        const extracted = extractHeadersFromDOM()
        if (extracted.length > 0) {
          domHeaders.value = extracted
        }
      }, 300)
    }
  })
  
  // 定期的にチェック（VitePressのコンテンツが遅れてレンダリングされる場合に備える）
  const checkInterval = setInterval(() => {
    if (finalHeaders.value.length === 0 && headers.value.length === 0) {
      const extracted = extractHeadersFromDOM()
      if (extracted.length > 0) {
        domHeaders.value = extracted
        clearInterval(checkInterval)
      }
    } else {
      clearInterval(checkInterval)
    }
  }, 500)
  
  // 10秒後にタイムアウト
  setTimeout(() => {
    clearInterval(checkInterval)
  }, 10000)
})

// 最終的な見出しデータ（page.value.headersまたはDOMから取得したもの）
const finalHeaders = computed(() => {
  if (headers.value && headers.value.length > 0) {
    return headers.value
  }
  return domHeaders.value || []
})

// 目次のラベルを取得（設定ファイルから）
const label = computed(() => {
  const themeConfig = theme.value as any
  return themeConfig?.outline?.label || '目次'
})

// 見出しが存在するかチェック
const hasHeaders = computed(() => {
  return finalHeaders.value && finalHeaders.value.length > 0
})

// 見出しの階層構造を処理する関数
const processHeaders = (headers: any[]) => {
  const result: any[] = []
  const stack: any[] = []
  
  headers.forEach((header) => {
    const level = header.level
    
    // スタックから現在のレベルより深いものを削除
    while (stack.length > 0 && stack[stack.length - 1].level >= level) {
      stack.pop()
    }
    
    const item = {
      ...header,
      children: []
    }
    
    if (stack.length === 0) {
      result.push(item)
    } else {
      stack[stack.length - 1].children.push(item)
    }
    
    stack.push(item)
  })
  
  return result
}

const processedHeaders = computed(() => {
  if (!finalHeaders.value || finalHeaders.value.length === 0) {
    return []
  }
  return processHeaders(finalHeaders.value)
})

// アンカーリンクを生成
const getAnchor = (header: any) => {
  return `#${header.slug}`
}

// スムーズスクロールを実行する関数
const handleSmoothScroll = (e: Event, header: any) => {
  e.preventDefault()
  const targetId = header.slug
  const targetElement = document.getElementById(targetId)
  
  if (targetElement) {
    // ヘッダーの高さを考慮したオフセット（必要に応じて調整）
    const offset = 80
    const elementPosition = targetElement.getBoundingClientRect().top
    const offsetPosition = elementPosition + window.pageYOffset - offset
    
    window.scrollTo({
      top: offsetPosition,
      behavior: 'smooth'
    })
    
    // URLを更新（ブラウザの履歴に追加しない）
    const url = new URL(window.location.href)
    url.hash = `#${targetId}`
    window.history.replaceState({}, '', url.toString())
  }
}
</script>

<template>
  <nav class="table-of-contents">
    <h2 class="toc-title">{{ label }}</h2>
    <div v-if="!hasHeaders" style="font-size: 1.2rem; color: var(--text-option); padding: 1rem;">
      見出しが見つかりませんでした
    </div>
    <ul v-else class="toc-list">
      <li
        v-for="header in processedHeaders"
        :key="header.slug"
        class="toc-item"
        :class="`toc-level-${header.level}`"
      >
        <a :href="getAnchor(header)" class="toc-link" @click.prevent="handleSmoothScroll($event, header)">
          {{ header.title }}
        </a>
        <ul v-if="header.children && header.children.length > 0" class="toc-sublist">
          <li
            v-for="child in header.children"
            :key="child.slug"
            class="toc-item"
            :class="`toc-level-${child.level}`"
          >
            <a :href="getAnchor(child)" class="toc-link" @click.prevent="handleSmoothScroll($event, child)">
              {{ child.title }}
            </a>
            <ul v-if="child.children && child.children.length > 0" class="toc-sublist">
              <li
                v-for="grandchild in child.children"
                :key="grandchild.slug"
                class="toc-item"
                :class="`toc-level-${grandchild.level}`"
              >
                <a :href="getAnchor(grandchild)" class="toc-link" @click.prevent="handleSmoothScroll($event, grandchild)">
                  {{ grandchild.title }}
                </a>
              </li>
            </ul>
          </li>
        </ul>
      </li>
    </ul>
  </nav>
</template>

<style scoped>
.table-of-contents {
  position: sticky;
  top: 8rem;
  width: 100%;
  max-height: calc(100vh - 4rem);
  overflow-y: auto;
  padding: 0;
}

.toc-title {
  font-size: 1.4rem;
  font-weight: bold;
  color: var(--text-bold);
  margin-bottom: 1.2rem;
  line-height: 1.2;
}

.toc-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.toc-sublist {
  list-style: none;
  padding: 0;
  margin: 0.5rem 0 0 0;
  padding-left: 1.2rem;
  margin-top: .7rem;
  margin-bottom: 2rem;
}

.toc-item {
  margin-bottom: 1rem;
}

.toc-link {
  display: block;
  color: var(--text-main);
  text-decoration: none;
  font-size: 1.2rem;
  line-height: 1.4;
  transition: color 0.2s ease;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.toc-link:hover {
  color: var(--text-link);
}

.toc-level-2 .toc-link {
  font-weight: 600;
  font-size: 1.1rem;
}

.toc-level-3 .toc-link {
  font-weight: 500;
  font-size: 1rem;
}

.toc-level-4 .toc-link {
  font-weight: 400;
  font-size: 1rem;
}

/* スクロールバーのスタイル */
.table-of-contents::-webkit-scrollbar {
  width: 4px;
}

.table-of-contents::-webkit-scrollbar-track {
  background: transparent;
}

.table-of-contents::-webkit-scrollbar-thumb {
  background: var(--bg-accent-light);
  border-radius: 2px;
}

.table-of-contents::-webkit-scrollbar-thumb:hover {
  background: var(--bg-accent);
}
</style>

