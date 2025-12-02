<script setup lang="ts">
import { computed, ref } from 'vue'
import { withBase } from 'vitepress'

type NavItem = {
  text: string
  link?: string
  items?: NavItem[]
}

const props = defineProps<{
  items: NavItem[]
  /** 再帰レベル（1 がルート） */
  level?: number
}>()

const openState = ref<Record<string, boolean>>({})

const normalizedItems = computed(() => props.items ?? [])
const level = computed(() => props.level ?? 1)

const hasItems = (item: NavItem) => Array.isArray(item.items) && item.items.length > 0

const keyFor = (index: number) => `${level.value}-${index}`

const isOpen = (index: number) => !!openState.value[keyFor(index)]

const toggle = (index: number) => {
  const key = keyFor(index)
  openState.value[key] = !openState.value[key]
}
</script>

<template>
  <ul
    class="sidebar-nav-list"
    :class="{ 'sidebar-nav-root': level === 1 }"
  >
    <li
      v-for="(item, index) in normalizedItems"
      :key="item.text ?? index"
      class="sidebar-nav-item"
    >
      <!-- 第1階層: デフォルト表示・トグルなし -->
      <template v-if="level === 1">
        <div class="sidebar-nav-section">
          <span class="sidebar-nav-section-title">
            {{ item.text }}
          </span>
        </div>
        <ListMenu
          v-if="hasItems(item)"
          :items="item.items || []"
          :level="level + 1"
        />
      </template>

      <!-- 第2階層以降: トグルで展開／折りたたみ -->
      <template v-else>
        <template v-if="hasItems(item)">
          <button
            type="button"
            class="sidebar-nav-link sidebar-nav-toggle"
            :aria-expanded="isOpen(index)"
            @click="toggle(index)"
          >
            <span class="sidebar-nav-text">
              {{ item.text }}
            </span>
            <IconLoader
              name="icon-chevron"
              :width="14"
              :height="14"
              fill="transparent"
              :rotate="isOpen(index) ? -90 : 90"
              class="sidebar-nav-toggle-icon"
              aria-hidden="true"
            />
          </button>

          <transition name="sidebar-nav-collapse">
            <div
              v-if="isOpen(index)"
              class="sidebar-nav-children"
            >
              <ListMenu
                :items="item.items || []"
                :level="level + 1"
              />
            </div>
          </transition>
        </template>
        <template v-else>
          <a
            v-if="item.link"
            class="sidebar-nav-link sidebar-nav-leaf"
            :href="withBase(item.link)"
          >
            {{ item.text }}
          </a>
          <span
            v-else
            class="sidebar-nav-link sidebar-nav-leaf sidebar-nav-leaf--static"
          >
            {{ item.text }}
          </span>
        </template>
      </template>
    </li>
  </ul>
</template>

<style scoped>
.sidebar-nav-list {
  list-style: none;
  margin: 0;
  margin-top: 1.25rem;
  margin-bottom: 1rem;
  padding-left: 0;

  display: flex;
  flex-direction: column;
  gap: .65rem;

  &.sidebar-nav-root {
    display: flex;
    flex-direction: column;
    gap: 3rem;
    margin-top: 0;
  }
}

.sidebar-nav-item {
  margin: 0;

  .sidebar-nav-section {
    .sidebar-nav-section-title {
      display: block;
      font-size: 1.25rem;
      font-weight: 700;
      color: var(--text-option);
      text-transform: uppercase;
      letter-spacing: 0.08em;
      margin-bottom: 0.5rem;
    }
  }

  .sidebar-nav-link {
    width: 100%;
    display: flex;
    align-items: center;
    text-align: left;
    gap: 0.5rem;

    font-size: 1.4rem;
    font-weight: 600;
    color: var(--text-main, var(--text-bold));
    text-decoration: none;

    padding: 0.25rem 0;
    border: none;
    background: transparent;
    cursor: pointer;

    &:hover {
      color: var(--accent, var(--text-bold));
    }

    &.sidebar-nav-leaf {
      /* padding-left: 0.75rem; */

      &.sidebar-nav-leaf--static {
        cursor: pointer;
      }
    }

    .sidebar-nav-toggle-icon {
      flex-shrink: 0;
      margin-left: 0.25rem;
    }
  }

  .sidebar-nav-children {
    margin-top: 0.35rem;
    margin-left: 2rem;
    padding-left: 1.2rem;
    border-left: 1px solid var(--stroke-light, rgba(148, 163, 184, 0.4));

    * {
      font-weight: 400;
    }
  }
}

.sidebar-nav-collapse-enter-active,
.sidebar-nav-collapse-leave-active {
  transition: all 0.12s ease-out;
}

.sidebar-nav-collapse-enter-from,
.sidebar-nav-collapse-leave-to {
  opacity: 0;
  transform: translateY(-2px);
}
</style>