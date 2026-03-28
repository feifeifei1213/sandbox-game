<template>
  <div class="admin-page">
    <div class="shell">
      <header class="page-header">
        <div>
          <p class="eyebrow">Sandbox Game / Admin</p>
          <h1>管理员端</h1>
          <p class="subtext">首版聚焦汇总、年度控制、初始基线与组数据四个入口，优先服务现场主持人操作。</p>
        </div>
        <div class="header-pills">
          <span class="pill">最终年份：{{ config?.finalYear ?? '--' }}</span>
          <span class="pill">当前开放：{{ config?.currentOpenYear ?? '--' }}</span>
          <span class="pill">共享基线：{{ config?.initialBaselineSubmitted ? '已提交' : '未提交' }}</span>
          <span class="pill">下一年：{{ config?.nextOpenableYear ?? '--' }}</span>
        </div>
      </header>

      <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
        {{ pageMessage.text }}
      </section>

      <div class="workspace">
        <aside class="nav-panel">
          <AdminNav />
        </aside>
        <main class="main-panel">
          <RouterView />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { RouterView } from 'vue-router'

import AdminNav from '@/components/sandbox-game/admin/AdminNav.vue'
import { useAdminShellStore } from '@/stores/admin-shell'

const shellStore = useAdminShellStore()
const { config, pageMessage } = storeToRefs(shellStore)

onMounted(async () => {
  if (config.value) {
    return
  }
  try {
    await shellStore.bootstrap()
  } catch {
    // 错误消息由 store 统一展示。
  }
})
</script>

<style scoped>
.admin-page {
  min-height: 100vh;
  padding: 20px;
}

.shell {
  width: min(1760px, calc(100vw - 24px));
  margin: 0 auto;
  background: var(--shell-bg);
  border: 1px solid #dbe2ea;
  border-radius: 22px;
  box-shadow: var(--shadow);
  overflow: hidden;
}

.page-header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 22px 24px 18px;
  background: linear-gradient(180deg, #fbfdff 0%, #eef4fb 100%);
  border-bottom: 1px solid #dde5ef;
}

.page-header h1 {
  margin: 6px 0 8px;
  font-size: 28px;
}

.eyebrow {
  margin: 0;
  color: var(--accent);
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.subtext {
  margin: 0;
  color: var(--muted);
}

.header-pills {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  align-content: flex-start;
  gap: 10px;
}

.pill {
  display: inline-flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: #ffffff;
  font-size: 13px;
}

.message-bar {
  margin: 16px 18px 0;
  padding: 12px 14px;
  border-radius: 14px;
  font-size: 14px;
}

.message-bar.success {
  background: #edfdf3;
  color: var(--success);
  border: 1px solid #b7e2c5;
}

.message-bar.error {
  background: #fff5f5;
  color: var(--danger);
  border: 1px solid #efc4c4;
}

.message-bar.info {
  background: var(--accent-soft);
  color: var(--accent);
  border: 1px solid #cbdcff;
}

.workspace {
  display: grid;
  grid-template-columns: 250px minmax(0, 1fr);
  gap: 16px;
  padding: 16px;
  align-items: start;
}

.nav-panel {
  position: sticky;
  top: 18px;
  border: 1px solid var(--line);
  border-radius: 18px;
  background: linear-gradient(180deg, #fbfcfe 0%, #f4f7fb 100%);
  padding: 14px;
}

.main-panel {
  min-width: 0;
}

@media (max-width: 1180px) {
  .workspace {
    grid-template-columns: 1fr;
  }

  .nav-panel {
    position: static;
  }
}

@media (max-width: 1024px) {
  .admin-page {
    padding: 12px;
  }

  .shell {
    width: 100%;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-pills {
    justify-content: flex-start;
  }
}
</style>
