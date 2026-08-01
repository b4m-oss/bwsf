import { onMounted } from 'vue'

const GA_MEASUREMENT_ID = 'G-CNCM19EPXM'

export function useGoogleAnalytics() {
  function isAnalyticsEnabled(): boolean {
    if (typeof localStorage === 'undefined') return false
    const consent = localStorage.getItem('cookie-consent')
    // Only enable if explicitly accepted or not yet decided (null)
    // If declined, do not enable
    return consent !== 'declined'
  }

  function loadGoogleAnalytics() {
    if (typeof window === 'undefined') return
    if (!isAnalyticsEnabled()) return
    
    // Check if already loaded
    if (document.querySelector(`script[src*="googletagmanager.com/gtag"]`)) return

    // Load gtag.js
    const script = document.createElement('script')
    script.async = true
    script.src = `https://www.googletagmanager.com/gtag/js?id=${GA_MEASUREMENT_ID}`
    document.head.appendChild(script)

    // Initialize gtag
    script.onload = () => {
      (window as any).dataLayer = (window as any).dataLayer || []
      function gtag(...args: any[]) {
        (window as any).dataLayer.push(arguments)
      }
      (window as any).gtag = gtag
      gtag('js', new Date())
      gtag('config', GA_MEASUREMENT_ID)
    }
  }

  onMounted(() => {
    loadGoogleAnalytics()
  })
}


