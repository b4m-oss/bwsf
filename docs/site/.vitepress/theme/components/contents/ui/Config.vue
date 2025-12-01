<script setup lang="ts">
import AppDropdown from '../../common/ui/AppDropdown.vue';
import IconLoader from '../../common/ui/IconLoader.vue';
import AppAppearanceSwitch from '../../common/ui/AppAppearanceSwitch.vue';
import AppLanguageSwitch from '../../common/ui/AppLanguageSwitch.vue';
</script>

<template>
  <AppDropdown class="config">
    <template #trigger="{ isOpen, toggle }">
      <button
        type="button"
        :aria-expanded="isOpen"
        aria-haspopup="listbox"
        @click="toggle"
      >
        <IconLoader name="icon-cog" :width="23" :height="23" fill="currentColor" />
      </button>
    </template>
    
    <template #menu="{ close }">
      <div class="config__dropdown">
        <div class="config__menu">
          <b>言語設定</b>
          <AppLanguageSwitch inline />
        </div>
        <div class="config__menu">
          <b>テーマ設定</b>
          <AppAppearanceSwitch inline />
        </div>
      </div>
    </template>
  </AppDropdown>
</template>

<style scoped>
.config {
  position: relative;
}

.config__dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 0.25rem;
  min-width: 16rem;
  padding: 0.75rem 1rem;
  border-radius: 0.5rem;
  background: var(--bg-main);
  border: 1px solid var(--text-option);
  box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.1);
  z-index: 100;

  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.config__menu {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.25rem;

  b {
    font-size: 1.1rem;
    font-weight: 800;
    color: var(--text-option);
    margin-bottom: 0.15rem;
  }
}
</style>