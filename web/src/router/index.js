import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { guest: true }
  },
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/HomeNew.vue')
  },
  {
    path: '/services',
    name: 'Services',
    component: () => import('@/views/Services.vue')
  },
  {
    path: '/about',
    name: 'About',
    component: () => import('@/views/About.vue')
  },
  {
    path: '/counselors',
    name: 'Counselors',
    component: () => import('@/views/Counselors.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/counselor/:id',
    name: 'CounselorDetail',
    component: () => import('@/views/CounselorDetail.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/orders',
    name: 'Orders',
    component: () => import('@/views/Orders.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/chat/:sessionId',
    name: 'Chat',
    component: () => import('@/views/Chat.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/faq',
    name: 'FAQ',
    component: () => import('@/views/FAQ.vue')
  },
  {
    path: '/privacy',
    name: 'Privacy',
    component: () => import('@/views/PrivacyPolicy.vue')
  },
  {
    path: '/terms',
    name: 'Terms',
    component: () => import('@/views/TermsOfService.vue')
  },
  {
    path: '/guide',
    name: 'Guide',
    component: () => import('@/views/BookingGuide.vue')
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore()
  const token = userStore.token

  if (to.meta.requiresAuth && !token) {
    // 需要认证但没有token，跳转到登录页
    next('/login')
  } else if (to.meta.requiresAuth && token) {
    // 需要认证且有token，检查用户信息
    if (!userStore.userInfo) {
      try {
        await userStore.fetchUserInfo()
        next()
      } catch (error) {
        console.error('获取用户信息失败:', error)
        userStore.logout()
        next('/login')
      }
    } else {
      next()
    }
  } else if (to.meta.guest && token) {
    // 访客页面但有token，跳转到首页
    next('/')
  } else {
    next()
  }
})

// 捕获路由错误
router.onError((error) => {
  console.error('路由导航错误:', error)
})

export default router
