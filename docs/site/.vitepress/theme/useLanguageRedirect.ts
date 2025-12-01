import { onMounted } from 'vue'
import { useRoute, useData } from 'vitepress'

export function useLanguageRedirect() {
  const route = useRoute()
  const { site } = useData()

  function getPreferredLanguage(): 'ja' | 'en' {
    if (typeof navigator === 'undefined') return 'en'
    
    const browserLang = navigator.language || (navigator as any).userLanguage || ''
    return browserLang.toLowerCase().startsWith('ja') ? 'ja' : 'en'
  }

  function getBasePath(): string {
    const base = site.value.base || '/'
    return base.endsWith('/') ? base.slice(0, -1) : base
  }

  function isRootPage(): boolean {
    // Check if current path is the root page (e.g., /, /bwenv/, or similar)
    const path = route.path
    return path === '/' || path === ''
  }

  function hasAlreadyRedirected(): boolean {
    if (typeof sessionStorage === 'undefined') return false
    return sessionStorage.getItem('language-auto-redirected') === 'true'
  }

  function markAsRedirected() {
    if (typeof sessionStorage !== 'undefined') {
      sessionStorage.setItem('language-auto-redirected', 'true')
    }
  }

  function performRedirect() {
    if (typeof window === 'undefined') return
    
    // Only redirect on the root page
    if (!isRootPage()) return
    
    // Don't redirect if already redirected in this session
    if (hasAlreadyRedirected()) return

    const preferredLang = getPreferredLanguage()
    const basePath = getBasePath()
    const targetPath = `${basePath}/${preferredLang}/`

    markAsRedirected()
    window.location.replace(targetPath)
  }

  onMounted(() => {
    performRedirect()
  })
}
