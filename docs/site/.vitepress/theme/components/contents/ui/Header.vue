<script setup lang="ts">
import { withBase } from 'vitepress'
// import AppAppearanceSwitch from '../../common/ui/AppAppearanceSwitch.vue'
// import AppLanguageSwitch from '../../common/ui/AppLanguageSwitch.vue'
import Config from './Config.vue'
import SearchBox from './SearchBox.vue'
import TypewriterCaption from '../../common/ui/TypewriterCaption.vue'

const messages = [
  '- .env management no more effort',
  '- Bring Bitwarden and .env together',
  '- Focus on your code, not secrets sync',
]
</script>

<template>
  <header id="global-header">
    <p class="logo">
      <a :href="withBase('/')">
        <img
          :src="withBase('/logo.svg')"
          alt="bwsfのロゴマーク"
          class="logo-image"
        >
        <span class="logo-text">bwsf</span>
        <TypewriterCaption
          :messages="messages"
          :typing-speed-ms="90"
          :deleting-speed-ms="10"
          :delay-after-type-ms="7000"
          :delay-after-delete-ms="1000"
        />
      </a>
    </p>
    
    <nav class="header-nav">
      <ul class="header-nav-list">
        <li class="header-nav-link header-nav-search">
          <SearchBox :placeholder="localeIndex === 'ja' ? 'ドキュメントを検索' : 'Search docs'" :search-options="{ fuzzy: 0.2, prefix: true, boost: { title: 4, text: 2, titles: 1 } }" />
        </li>
        <li class="header-nav-link">
          <a :href="withBase('/ja/guide/getting-started')">Getting staretd</a>
        </li>
        <li class="header-nav-link">
          <a :href="withBase('/ja/guide/commands')">コマンド一覧</a>
        </li>
        <li class="header-nav-link">
          <a :href="withBase('/ja/guide/features')">機能</a>
        </li>
      </ul>
      <div class="header-configuration">
        <Config />
      </div>
    </nav>
  </header>
</template>

<style scoped>
#global-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  display: grid;
  grid-template-columns: 30% 70%;
  grid-template-rows: 1fr;
  align-items: center;

  padding-left: 2rem;
  padding-right: 2rem;
  padding-top: 1.5rem;
  padding-bottom: 1.5rem;

  border-bottom: 1px solid var(--storoke-light);

  .logo {
    a {
      display: flex;
      flex-direction: row;
      align-items: center;
      gap: 0.5rem;
    }

    .logo-image {
      width: 3.2rem;
      height: auto;
    }
    .logo-text {
      font-size: 1.5rem;
      font-weight: 800;
      color: var(--text-bold);
    }
    .caption {
      font-size: 1.1rem;
      font-weight: 400;
      color: var(--text-option);
    }
  }

  .header-nav {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    align-items: center;
    gap: 1.75rem;
  }


  .header-nav-list {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    gap: 4.25rem;

    .header-nav-link {
      font-size: 1.3rem;
      font-weight: 600;
      color: var(--text-option);
    }

    .header-nav-search {
      position: relative;
      display: flex;
      align-items: center;
      min-width: 16rem;
    }
  }
}

@media (max-width: 860px) {
  #global-header {
    .logo {
      .caption {
        display: none;
      }
    }
  }
}
</style>