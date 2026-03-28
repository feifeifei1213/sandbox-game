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
      component: () => import('@/views/sandbox-game/player/report/PlayerReportPage.vue'),
    },
    {
      path: '/sandbox-game/admin',
      component: () => import('@/views/sandbox-game/admin/AdminLayout.vue'),
      children: [
        {
          path: '',
          redirect: '/sandbox-game/admin/summary',
        },
        {
          path: 'summary',
          name: 'sandbox-admin-summary',
          component: () => import('@/views/sandbox-game/admin/summary/AdminSummaryPage.vue'),
        },
        {
          path: 'control',
          name: 'sandbox-admin-control',
          component: () => import('@/views/sandbox-game/admin/control/AdminControlPage.vue'),
        },
        {
          path: 'baseline',
          name: 'sandbox-admin-baseline',
          component: () => import('@/views/sandbox-game/admin/baseline/AdminBaselinePage.vue'),
        },
        {
          path: 'group-data',
          name: 'sandbox-admin-group-data',
          component: () => import('@/views/sandbox-game/admin/group-data/AdminGroupDataPage.vue'),
        },
      ],
    },
  ],
})

export default router
