import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/sandbox-game/player/operating',
    },
    {
      path: '/sandbox-game/player/operating',
      name: 'sandbox-player-operating',
      component: () => import('@/views/sandbox-game/player/operating/PlayerOperatingPage.vue'),
    },
    {
      path: '/sandbox-game/player/report',
      name: 'sandbox-player-report',
      component: () => import('@/views/sandbox-game/player/report/PlayerReportPlaceholder.vue'),
    },
  ],
})

export default router