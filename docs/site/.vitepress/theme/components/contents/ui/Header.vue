<script setup lang="ts">
import { withBase } from 'vitepress'
// import AppAppearanceSwitch from '../../common/ui/AppAppearanceSwitch.vue'
// import AppLanguageSwitch from '../../common/ui/AppLanguageSwitch.vue'
import Config from './Config.vue';
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
        <span class="caption">- .env management no more effort</span>
      </a>
    </p>
    
    <nav class="header-nav">
      <ul class="header-nav-list">
        <li class="header-nav-link"><a :href="withBase('/ja/guide/getting-started')">Getting staretd</a></li>
        <li class="header-nav-link"><a :href="withBase('/ja/guide/commands')">コマンド一覧</a></li>
        <li class="header-nav-link"><a :href="withBase('/ja/guide/features')">機能</a></li>
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
  }
}
</style>