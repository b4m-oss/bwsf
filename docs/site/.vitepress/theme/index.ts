import type { Theme } from 'vitepress'
import CookieConsent from './components/common/CookieConsent.vue'
import CookieSettings from './components/common/CookieSettings.vue'
import IconLoader from './components/common/ui/IconLoader.vue'
import Accordion from './components/common/ui/Accordion.vue'
import AccordionQA from './components/common/ui/AccordionQA.vue'
import AccordionList from './components/common/ui/AccordionList.vue'
import ToggleButton from './components/common/ui/ToggleButton.vue'
import { useLanguageRedirect } from './useLanguageRedirect'
import { useGoogleAnalytics } from './useGoogleAnalytics'
import 'reset-css'

// Nuxt 風の pages/layouts/components 構成用に、
// レイアウト・コンポーネントをここで登録しておきます。
import HomeLayout from './layouts/HomeLayout.vue'
import Header from './components/contents/ui/Header.vue'
import Footer from './components/contents/ui/Footer.vue'
import SubPageLayout from './layouts/SubPageLayout.vue'
import AppLayout from './layouts/AppLayout.vue'
import HeroSection from './components/contents/HeroSection.vue'
import HeroFeatureCard from './components/contents/HeroFeatureCard.vue'
import SidebarMenu from './components/contents/ui/Sidebar.vue'
import TableOfContents from './components/contents/ui/TableOfContents.vue'
import DocFooter from './components/contents/ui/DocFooter.vue'

// サイト全体で使うカスタム CSS
import './styles/main.css'

export default {
  // 完全に自前実装のレイアウトを使う（DefaultTheme には依存しない）
  Layout: AppLayout,
  enhanceApp({ app }) {
    // グローバルコンポーネントとして登録しておくことで、
    // 任意の Markdown ページから <HomeLayout> や <HeroSection> を使えるようにします。
    app.component('CookieSettings', CookieSettings)
    app.component('Header', Header)
    app.component('Footer', Footer)
    app.component('SidebarMenu', SidebarMenu)
    app.component('TableOfContents', TableOfContents)
    app.component('DocFooter', DocFooter)

    // layouts
    app.component('HomeLayout', HomeLayout)
    app.component('SubPageLayout', SubPageLayout)

    // components
    app.component('HeroSection', HeroSection)
    app.component('HeroFeatureCard', HeroFeatureCard)
    app.component('CookieConsent', CookieConsent)
    app.component('IconLoader', IconLoader)
    app.component('Accordion', Accordion)
    app.component('AccordionQA', AccordionQA)
    app.component('AccordionList', AccordionList)
    app.component('ToggleButton', ToggleButton)
  },
  setup() {
    useLanguageRedirect()
    useGoogleAnalytics()
  }
} satisfies Theme

