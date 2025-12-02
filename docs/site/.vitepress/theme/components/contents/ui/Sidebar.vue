<script setup lang="ts">
import { computed } from 'vue'
import { useData } from 'vitepress'
import ListMenu from '../../common/navigation/ListMenu.vue'

const { theme } = useData()

// config.ts の locales.xx.themeConfig.sidebar がここに入ってくる
const items = computed(() => {
  const t = theme.value as any
  return t?.sidebar ?? []
})

</script>

<template>
  <ListMenu
    :items="items"
    variant="sidebar"
    class="sidebar-menu"
  />
</template>

<style scoped>
.sidebar-menu {
  position: sticky;
  top: 8rem;
}

@media (max-width: 860px) {
  .sidebar-menu {
    position: static;
    top: auto;
  }
}
</style>