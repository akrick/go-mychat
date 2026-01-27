import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', noAuth: true }
  },
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/dashboard',
    meta: { title: '首页' },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '数据看板', icon: 'DataAnalysis' }
      },
      {
        path: 'user',
        name: 'User',
        component: () => import('@/views/system/user/index.vue'),
        meta: { title: '用户管理', icon: 'User' }
      },
      {
        path: 'counselor',
        name: 'Counselor',
        component: () => import('@/views/system/counselor/index.vue'),
        meta: { title: '咨询师管理', icon: 'UserFilled' }
      },
      {
        path: 'order',
        name: 'Order',
        component: () => import('@/views/business/order/index.vue'),
        meta: { title: '订单管理', icon: 'Document' }
      },
      {
        path: 'chat',
        name: 'Chat',
        component: () => import('@/views/business/chat/index.vue'),
        meta: { title: '聊天记录', icon: 'ChatDotRound' }
      },
      {
        path: 'roles',
        name: 'Roles',
        component: () => import('@/views/system/roles/index.vue'),
        meta: { title: '角色管理', icon: 'Avatar' }
      },
      {
        path: 'permissions',
        name: 'Permissions',
        component: () => import('@/views/system/permissions/index.vue'),
        meta: { title: '权限管理', icon: 'Key' }
      },
      {
        path: 'menus',
        name: 'Menus',
        component: () => import('@/views/system/menus/index.vue'),
        meta: { title: '菜单管理', icon: 'Menu' }
      },
      {
        path: 'lowcode/forms',
        name: 'LowcodeForms',
        component: () => import('@/views/lowcode/forms/index.vue'),
        meta: { title: '表单设计', icon: 'Edit', permission: 'lowcode:form:design' }
      },
      {
        path: 'lowcode/pages',
        name: 'LowcodePages',
        component: () => import('@/views/lowcode/pages/index.vue'),
        meta: { title: '页面设计', icon: 'Grid', permission: 'lowcode:page:design' }
      },
      {
        path: 'lowcode/data',
        name: 'LowcodeData',
        component: () => import('@/views/lowcode/data/index.vue'),
        meta: { title: '数据管理', icon: 'Database', permission: 'lowcode:data:manage' }
      },
      {
        path: 'withdraw',
        name: 'Withdraw',
        component: () => import('@/views/finance/withdraw/index.vue'),
        meta: { title: '提现审核', icon: 'Wallet' }
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: () => import('@/views/finance/statistics/index.vue'),
        meta: { title: '财务统计', icon: 'DataLine' }
      },
      {
        path: 'reports',
        name: 'Reports',
        component: () => import('@/views/finance/reports/index.vue'),
        meta: { title: '财务报表', icon: 'Document' }
      },
      {
        path: 'logs',
        name: 'Logs',
        component: () => import('@/views/system/logs/index.vue'),
        meta: { title: '系统日志', icon: 'Document', permission: 'system:log:list' }
      },
      {
        path: 'online',
        name: 'Online',
        component: () => import('@/views/system/online/index.vue'),
        meta: { title: '在线用户', icon: 'User', permission: 'system:online:list' }
      },
      {
        path: 'config',
        name: 'Config',
        component: () => import('@/views/system/config/index.vue'),
        meta: { title: '系统配置', icon: 'Setting', permission: 'system:config:manage' }
      }
    ]
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/profile/index.vue'),
    meta: { title: '个人中心', icon: 'User' }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '404', noAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()
  const token = userStore.token

  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - MyChat管理后台` : 'MyChat管理后台'

  if (to.meta.noAuth) {
    next()
  } else {
    if (token) {
      if (to.path === '/login') {
        next({ path: '/' })
      } else {
        // 检查权限
        if (to.meta.permission && !userStore.hasPermission(to.meta.permission)) {
          next({ path: '/' })
        } else {
          next()
        }
      }
    } else {
      next({ path: '/login', query: { redirect: to.fullPath } })
    }
  }
})

export default router
