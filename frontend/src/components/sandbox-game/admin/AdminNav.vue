<template>
  <nav class="admin-nav">
    <RouterLink
      v-for="item in items"
      :key="item.to"
      :to="item.to"
      class="nav-item"
      :class="{ active: route.path === item.to }"
    >
      <strong>{{ item.label }}</strong>
      <span>{{ item.description }}</span>
    </RouterLink>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { RouterLink, useRoute } from 'vue-router'

import { useAdminShellStore } from '@/stores/admin-shell'

const route = useRoute()
const shellStore = useAdminShellStore()
const { setupStatus } = storeToRefs(shellStore)

const items = computed(() => {
  if (!setupStatus.value?.initialized) {
    return [
      { to: '/sandbox-game/admin/setup', label: '赛前配置', description: '先配置本场比赛小组数量，再初始化比赛环境' },
    ]
  }
  return [
    { to: '/sandbox-game/admin/summary', label: '汇总', description: '年度汇总区 + 最终排名区' },
    { to: '/sandbox-game/admin/control', label: '年度控制', description: '最终年份、开放下一年、阻断摘要' },
    { to: '/sandbox-game/admin/baseline', label: '初始基线', description: '共享模板录入与提交锁定' },
    { to: '/sandbox-game/admin/group-data', label: '组数据', description: '组与年份查看入口，支持异常解锁' },
    { to: '/sandbox-game/admin/notices', label: '通知与奖惩', description: '发送普通通知并按组下发奖励与罚款' },
    { to: '/sandbox-game/admin/orders', label: '订单管理', description: '配置市场开启、标段数量与释放顺序' },
  ]
})
</script>

<style scoped>
.admin-nav {
  display: grid;
  gap: 10px;
}

.nav-item {
  display: grid;
  gap: 4px;
  padding: 12px 14px;
  border-radius: 14px;
  border: 1px solid var(--line);
  background: #ffffff;
}

.nav-item strong {
  font-size: 15px;
}

.nav-item span {
  color: var(--muted);
  font-size: 12px;
  line-height: 1.5;
}

.nav-item.active {
  border-color: #b8cbf5;
  background: #eef4ff;
  box-shadow: inset 0 0 0 1px rgba(31, 95, 211, 0.06);
}
</style>
