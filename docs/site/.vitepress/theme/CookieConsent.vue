<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, withBase } from 'vitepress'

const route = useRoute()
const showBanner = ref(false)

const isJapanese = computed(() => route.path.startsWith('/ja/') || route.path === '/ja')

const texts = computed(() => {
  if (isJapanese.value) {
    return {
      message: '当サイトでは、サイト分析と最適なユーザー体験を提供するためにCookieを使用しています。',
      accept: '同意する',
      decline: '拒否する',
      learnMore: '詳細を見る',
      policyLink: '/ja/cookie-policy'
    }
  }
  return {
    message: 'We use cookies to analyze site traffic and provide you with the best experience.',
    accept: 'Accept',
    decline: 'Decline',
    learnMore: 'Learn more',
    policyLink: '/en/cookie-policy'
  }
})

const policyFullPath = computed(() => withBase(texts.value.policyLink))

onMounted(() => {
  const consent = localStorage.getItem('cookie-consent')
  if (!consent) {
    showBanner.value = true
  }
})

function acceptCookies() {
  localStorage.setItem('cookie-consent', 'accepted')
  showBanner.value = false
  // Enable Google Analytics if it was disabled
  if (typeof window !== 'undefined' && (window as any).gtag) {
    (window as any).gtag('consent', 'update', {
      analytics_storage: 'granted'
    })
  }
}

function declineCookies() {
  localStorage.setItem('cookie-consent', 'declined')
  showBanner.value = false
  // Disable Google Analytics
  if (typeof window !== 'undefined' && (window as any).gtag) {
    (window as any).gtag('consent', 'update', {
      analytics_storage: 'denied'
    })
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition name="slide-up">
      <div v-if="showBanner" class="cookie-consent">
        <div class="cookie-consent-content">
          <p class="cookie-consent-message">{{ texts.message }}</p>
          <div class="cookie-consent-actions">
            <a :href="policyFullPath" class="cookie-consent-link">
              {{ texts.learnMore }}
            </a>
            <button @click="declineCookies" class="cookie-consent-button decline">
              {{ texts.decline }}
            </button>
            <button @click="acceptCookies" class="cookie-consent-button accept">
              {{ texts.accept }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.cookie-consent {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: var(--vp-c-bg-soft);
  border-top: 1px solid var(--vp-c-border);
  box-shadow: 0 -4px 12px rgba(0, 0, 0, 0.1);
}

.cookie-consent-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  flex-wrap: wrap;
}

.cookie-consent-message {
  margin: 0;
  font-size: 14px;
  color: var(--vp-c-text-1);
  flex: 1;
  min-width: 200px;
}

.cookie-consent-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.cookie-consent-link {
  font-size: 14px;
  color: var(--vp-c-brand-1);
  text-decoration: none;
  white-space: nowrap;
}

.cookie-consent-link:hover {
  text-decoration: underline;
}

.cookie-consent-button {
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.cookie-consent-button.decline {
  background: transparent;
  border: 1px solid var(--vp-c-border);
  color: var(--vp-c-text-2);
}

.cookie-consent-button.decline:hover {
  border-color: var(--vp-c-text-2);
  color: var(--vp-c-text-1);
}

.cookie-consent-button.accept {
  background: var(--vp-c-brand-1);
  border: 1px solid var(--vp-c-brand-1);
  color: var(--vp-c-white);
}

.cookie-consent-button.accept:hover {
  background: var(--vp-c-brand-2);
  border-color: var(--vp-c-brand-2);
}

/* Slide up transition */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(100%);
  opacity: 0;
}

/* Mobile responsive */
@media (max-width: 640px) {
  .cookie-consent-content {
    flex-direction: column;
    align-items: stretch;
    text-align: center;
  }

  .cookie-consent-actions {
    justify-content: center;
  }
}
</style>
