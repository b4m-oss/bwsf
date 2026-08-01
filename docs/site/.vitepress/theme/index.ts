import DefaultTheme from 'vitepress/theme'
import type { Theme } from 'vitepress'
import { h } from 'vue'
import CookieConsent from './CookieConsent.vue'
import CookieSettings from './CookieSettings.vue'
import { useLanguageRedirect } from './useLanguageRedirect'
import { useGoogleAnalytics } from './useGoogleAnalytics'

export default {
  extends: DefaultTheme,
  Layout() {
    return h(DefaultTheme.Layout, null, {
      'layout-bottom': () => h(CookieConsent)
    })
  },
  enhanceApp({ app }) {
    app.component('CookieSettings', CookieSettings)
  },
  setup() {
    useLanguageRedirect()
    useGoogleAnalytics()
  }
} satisfies Theme
