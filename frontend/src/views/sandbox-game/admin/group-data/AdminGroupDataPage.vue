<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>组数据</h2>
        <p>这个入口已经进入正式前端工程，但真实查询与异常解锁动作放在下一步 `M2-12` 接入。</p>
      </div>
    </header>

    <div class="placeholder-grid">
      <section class="panel-card">
        <div class="panel-head">
          <strong>查询条件区</strong>
          <span>结构已经固定，后续直接接管理员组数据接口。</span>
        </div>
        <div class="form-grid">
          <label class="field">
            <span>小组</span>
            <select disabled>
              <option>第1组</option>
            </select>
          </label>
          <label class="field">
            <span>年份</span>
            <select disabled>
              <option>0年</option>
            </select>
          </label>
          <label class="field">
            <span>页面类型</span>
            <select disabled>
              <option>经营页</option>
            </select>
          </label>
        </div>
      </section>

      <section class="panel-card">
        <div class="panel-head">
          <strong>当前计划</strong>
          <span>下一步会把真实只读预览与异常解锁弹窗接进来。</span>
        </div>
        <ul class="note-list">
          <li>复用玩家经营页 / 财报页的只读视图结构。</li>
          <li>保留“经营 / 财报”切换入口。</li>
          <li>异常解锁按钮与原因弹窗在 `M2-12` 一并接入。</li>
        </ul>
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

import { useAdminShellStore } from '@/stores/admin-shell'

const shellStore = useAdminShellStore()

onMounted(async () => {
  if (shellStore.config) {
    return
  }
  try {
    await shellStore.bootstrap()
  } catch {
    // 错误消息由 shell store 统一处理。
  }
})
</script>

<style scoped>
.page-content {
  display: grid;
  gap: 16px;
}

.hero h2 {
  margin: 0 0 6px;
  font-size: 24px;
}

.hero p {
  margin: 0;
  color: var(--muted);
  line-height: 1.6;
}

.placeholder-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.panel-card {
  border: 1px solid var(--line);
  border-radius: 18px;
  background: #ffffff;
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.05);
}

.panel-head {
  padding: 14px 16px;
  background: #f7f9fc;
  border-bottom: 1px solid var(--line);
}

.panel-head strong {
  display: block;
  margin-bottom: 4px;
  font-size: 16px;
}

.panel-head span {
  color: var(--muted);
  font-size: 13px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  padding: 16px;
}

.field {
  display: grid;
  gap: 6px;
}

.field span {
  color: var(--muted);
  font-size: 13px;
  font-weight: 700;
}

.field select {
  border: 1px solid var(--line);
  border-radius: 12px;
  background: var(--readonly-bg);
  padding: 10px 12px;
}

.note-list {
  margin: 0;
  padding: 16px 16px 16px 32px;
  color: var(--muted);
  line-height: 1.8;
}

@media (max-width: 1080px) {
  .placeholder-grid,
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
