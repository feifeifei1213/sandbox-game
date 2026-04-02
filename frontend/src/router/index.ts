import { createRouter, createWebHistory, type RouteLocationNormalized, type RouteRecordRaw } from 'vue-router'

import { pinia } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import type { AuthRoleType } from '@/types/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/sandbox-game/login',
  },
  {
    path: '/sandbox-game/login',
    name: 'sandbox-login',
    component: () => import('@/views/sandbox-game/login/LoginPage.vue'),
    meta: {
      public: true,
    },
  },
  {
    path: '/sandbox-game/player/operating',
    name: 'sandbox-player-operating',
    component: () => import('@/views/sandbox-game/player/operating/PlayerOperatingPage.vue'),
    meta: {
      requiresAuth: true,
      roleType: 'GROUP',
    },
  },
  {
    path: '/sandbox-game/player/report',
    name: 'sandbox-player-report',
    component: () => import('@/views/sandbox-game/player/report/PlayerReportPage.vue'),
    meta: {
      requiresAuth: true,
      roleType: 'GROUP',
    },
  },
  {
    path: '/sandbox-game/admin',
    component: () => import('@/views/sandbox-game/admin/AdminLayout.vue'),
    meta: {
      requiresAuth: true,
      roleType: 'ADMIN',
    },
    children: [
      {
        path: '',
        redirect: '/sandbox-game/admin/summary',
      },
      {
        path: 'setup',
        name: 'sandbox-admin-setup',
        component: () => import('@/views/sandbox-game/admin/setup/AdminSetupPage.vue'),
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
      {
        path: 'notices',
        name: 'sandbox-admin-notices',
        component: () => import('@/views/sandbox-game/admin/notice/AdminNoticePage.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore(pinia)
  authStore.restoreSession()

  if (isPlayerPreviewRoute(to)) {
    return true
  }

  if (to.path === '/sandbox-game/login') {
    try {
      const authenticated = await authStore.ensureAuthenticated()
      if (authenticated) {
        return authStore.resolveDefaultRoute()
      }
    } catch {
      authStore.clearSession()
    }
    return true
  }

  if (!Boolean(to.meta.requiresAuth)) {
    return true
  }

  try {
    const authenticated = await authStore.ensureAuthenticated()
    if (!authenticated) {
      return { path: '/sandbox-game/login' }
    }
  } catch {
    return true
  }

  const requiredRoleType = getRequiredRoleType(to)
  if (requiredRoleType && authStore.currentUser?.roleType !== requiredRoleType) {
    return authStore.resolveDefaultRoute()
  }

  if (authStore.currentUser?.roleType === 'ADMIN') {
    const defaultRoute = authStore.resolveDefaultRoute()
    if (to.path === '/sandbox-game/admin') {
      return defaultRoute
    }
    if (defaultRoute === '/sandbox-game/admin/setup' && to.path !== defaultRoute) {
      return defaultRoute
    }
    if (defaultRoute !== '/sandbox-game/admin/setup' && to.path === '/sandbox-game/admin/setup') {
      return defaultRoute
    }
  }

  return true
})

function isPlayerPreviewRoute(to: RouteLocationNormalized) {
  if (!import.meta.env.DEV) {
    return false
  }
  if (!['/sandbox-game/player/operating', '/sandbox-game/player/report'].includes(to.path)) {
    return false
  }
  const rawPreview = Array.isArray(to.query.preview) ? to.query.preview[0] : to.query.preview
  return rawPreview === '1' || rawPreview === 'true'
}

function getRequiredRoleType(to: RouteLocationNormalized): AuthRoleType | undefined {
  return to.meta.roleType === 'ADMIN' || to.meta.roleType === 'GROUP'
    ? to.meta.roleType
    : undefined
}

export default router
