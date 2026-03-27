<template>
  <div class="placeholder-page">
    <div class="placeholder-card">
      <h1>玩家财报页</h1>
      <p>正式前端工程已经起好，这一轮先把玩家经营页接到真实接口并跑通。</p>
      <p>接下来会在同一套工程结构里继续接财报页、税率下拉和绿色手工项。</p>
      <RouterLink class="back-link" to="/sandbox-game/player/operating">返回经营页</RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.placeholder-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
}

.placeholder-card {
  width: min(680px, 100%);
  border: 1px solid var(--line);
  border-radius: 20px;
  background: #ffffff;
  padding: 32px;
  box-shadow: var(--shadow);
}

.placeholder-card h1 {
  margin-top: 0;
}

.back-link {
  display: inline-flex;
  margin-top: 16px;
  padding: 10px 14px;
  border-radius: 12px;
  background: var(--accent);
  color: #ffffff;
}
</style>