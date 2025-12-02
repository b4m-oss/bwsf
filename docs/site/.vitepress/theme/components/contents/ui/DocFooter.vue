<script setup lang="ts">
import { computed } from 'vue'
import { useData, useRoute, withBase } from 'vitepress'

const { page, theme, site } = useData()
const route = useRoute()

// 前後のページ情報を取得
// VitePressはpage.value.prev/nextを提供するが、カスタムテーマでは生成されない可能性がある
// その場合は、sidebarの設定から手動で計算する
const prev = computed(() => {
  const pageData = page.value as any
  // まず、VitePressが自動生成したprev/nextを確認
  if (pageData.prev) {
    return pageData.prev
  }
  
  // 自動生成されていない場合、sidebarから計算
  return calculatePrevNext().prev
})

const next = computed(() => {
  const pageData = page.value as any
  // まず、VitePressが自動生成したprev/nextを確認
  if (pageData.next) {
    return pageData.next
  }
  
  // 自動生成されていない場合、sidebarから計算
  return calculatePrevNext().next
})

// sidebarの設定から前後のページを計算する関数
function calculatePrevNext() {
  const currentPath = route.path
  const sidebar = (theme.value as any)?.sidebar || []
  
  // sidebarをフラット化して、すべてのリンクを取得
  const allLinks: Array<{ text: string; link: string }> = []
  
  function flattenSidebar(items: any[]) {
    items.forEach((item) => {
      if (item.link) {
        allLinks.push({ text: item.text, link: item.link })
      }
      if (item.items && Array.isArray(item.items)) {
        flattenSidebar(item.items)
      }
    })
  }
  
  sidebar.forEach((group: any) => {
    if (group.items && Array.isArray(group.items)) {
      flattenSidebar(group.items)
    }
  })
  
  // 現在のページのインデックスを探す
  const currentIndex = allLinks.findIndex((link) => {
    const linkPath = withBase(link.link)
    return linkPath === currentPath || currentPath.startsWith(linkPath)
  })
  
  let prevItem: { text: string; link: string } | null = null
  let nextItem: { text: string; link: string } | null = null
  
  if (currentIndex > 0) {
    prevItem = allLinks[currentIndex - 1]
  }
  
  if (currentIndex >= 0 && currentIndex < allLinks.length - 1) {
    nextItem = allLinks[currentIndex + 1]
  }
  
  return {
    prev: prevItem ? { text: prevItem.text, link: prevItem.link } : undefined,
    next: nextItem ? { text: nextItem.text, link: nextItem.link } : undefined
  }
}

// ラベルを取得（設定ファイルから）
const prevLabel = computed(() => {
  const themeConfig = theme.value as any
  return themeConfig?.docFooter?.prev || '前のページ'
})

const nextLabel = computed(() => {
  const themeConfig = theme.value as any
  return themeConfig?.docFooter?.next || '次のページ'
})

// 前後があるかチェック
const hasNavigation = computed(() => {
  return !!prev.value || !!next.value
})

// リンクを生成
const getLink = (item: any) => {
  if (!item) return ''
  return withBase(item.link)
}
</script>

<template>
  <nav v-if="hasNavigation" class="doc-footer">
    <div class="doc-footer-nav">
      <div v-if="prev" class="doc-footer-nav-item doc-footer-nav-prev">
        <a :href="getLink(prev)" class="doc-footer-nav-link">
          <span class="doc-footer-nav-label">
            <IconLoader name="icon-chevron" :width="12" :height="12" flip="horizontal" fill="transparent" :stroke-color="'var(--text-option)'" aria-hidden="true" />
            {{ prevLabel }}
          </span>
          {{ prev?.text || '' }}
        </a>
      </div>
      <div v-else class="doc-footer-nav-item doc-footer-nav-prev doc-footer-nav-item-empty" aria-hidden="true"></div>
      
      <div v-if="next" class="doc-footer-nav-item doc-footer-nav-next">
        <a :href="getLink(next)" class="doc-footer-nav-link">
          <span class="doc-footer-nav-label">
            {{ nextLabel }}
            <IconLoader name="icon-chevron" :width="12" :height="12" fill="transparent" :stroke-color="'var(--text-option)'" aria-hidden="true" />
          </span>
          {{ next?.text || '' }}
        </a>
      </div>
      <div v-else class="doc-footer-nav-item doc-footer-nav-next"></div>
    </div>
  </nav>
</template>

<style scoped>
.doc-footer {
  margin-top: 4rem;
  padding-top: 3rem;
  border-top: 1px solid var(--storoke-light);
}

.doc-footer-nav {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
}

.doc-footer-nav-item {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--stroke-option);
  border-radius: 0.25em;
  padding: .5em .75em;

  &:hover {
    background-color: var(--bg-option);
  }
}

.doc-footer-nav-prev {
  text-align: left;
}

.doc-footer-nav-next {
  text-align: right;
}

.doc-footer-nav-item-empty {
  opacity: 0;
}

.doc-footer-nav-label {
  font-size: 1.2rem;
  color: var(--text-option);
  margin-bottom: 0.5rem;
}

.doc-footer-nav-link {
  font-size: 1.6rem;
  font-weight: 600;
  color: var(--text-link);
  text-decoration: none;
  transition: color 0.2s ease;
  display: flex;
  flex-direction: column;
  /* align-items: center; */
  justify-content: center;
}

.doc-footer-nav-link:hover {
  color: var(--text-link-hover);
}

@media (max-width: 768px) {
  .doc-footer-nav {
    grid-template-columns: 1fr;
    gap: 2rem;
  }
  
  .doc-footer-nav-next {
    text-align: left;
  }
}
</style>

