import DefaultTheme from 'vitepress/theme'
import type { Theme } from 'vitepress'
import { h } from 'vue'
import CookieConsent from './CookieConsent.vue'
import CookieSettings from './CookieSettings.vue'
import { useLanguageRedirect } from './useLanguageRedirect'
import { useGoogleAnalytics } from './useGoogleAnalytics'

// Nuxt 風の pages/layouts/components 構成用に、
// レイアウト・コンポーネントをここで登録しておきます。
import HomeLayout from './layouts/HomeLayout.vue'
import DefaultPageLayout from './layouts/DefaultPageLayout.vue'
import HeroSection from './components/HeroSection.vue'

// サイト全体で使うカスタム CSS
import './custom.css'

export default {
  extends: DefaultTheme,
  Layout() {
    // 既存どおり、DefaultTheme.Layout をベースに CookieConsent を差し込みます。
    return h(DefaultTheme.Layout, null, {
      'layout-bottom': () => h(CookieConsent)
    })
  },
  enhanceApp({ app }) {
    // グローバルコンポーネントとして登録しておくことで、
    // 任意の Markdown ページから <HomeLayout> や <HeroSection> を使えるようにします。
    app.component('CookieSettings', CookieSettings)

    // layouts
    app.component('HomeLayout', HomeLayout)
    app.component('DefaultPageLayout', DefaultPageLayout)

    // components
    app.component('HeroSection', HeroSection)
  },
  setup() {
    useLanguageRedirect()
    useGoogleAnalytics()
  }
} satisfies Theme

