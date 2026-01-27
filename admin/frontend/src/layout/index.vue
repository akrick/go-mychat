<template>
  <div class="layout-container">
    <el-container>
      <el-aside :width="isCollapse ? '64px' : '200px'" class="sidebar">
        <div class="logo" :class="{ 'is-collapse': isCollapse }">
          <span v-if="!isCollapse">MyChat</span>
          <span v-else>MC</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          :collapse="isCollapse"
          :unique-opened="true"
          router
        >
          <template v-for="item in menuList" :key="item.path">
            <el-sub-menu v-if="item.children && item.children.length > 0" :index="item.path">
              <template #title>
                <el-icon><component :is="item.meta.icon" /></el-icon>
                <span>{{ item.meta.title }}</span>
              </template>
              <el-menu-item
                v-for="child in item.children"
                :key="child.path"
                :index="child.path"
              >
                <el-icon><component :is="child.meta.icon" /></el-icon>
                <span>{{ child.meta.title }}</span>
              </el-menu-item>
            </el-sub-menu>
            <el-menu-item v-else :index="item.path">
              <el-icon><component :is="item.meta.icon" /></el-icon>
              <span>{{ item.meta.title }}</span>
            </el-menu-item>
          </template>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header class="header">
          <div class="header-left">
            <el-icon class="collapse-btn" @click="toggleCollapse">
              <Fold v-if="!isCollapse" />
              <Expand v-else />
            </el-icon>
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
              <el-breadcrumb-item v-if="currentRoute.meta.title">
                {{ currentRoute.meta.title }}
              </el-breadcrumb-item>
            </el-breadcrumb>
          </div>

          <div class="header-right">
            <el-tooltip content="全屏" placement="bottom">
              <el-icon class="action-icon" @click="toggleFullScreen">
                <FullScreen />
              </el-icon>
            </el-tooltip>

            <el-tooltip content="刷新" placement="bottom">
              <el-icon class="action-icon" @click="handleRefresh">
                <Refresh />
              </el-icon>
            </el-tooltip>

            <el-dropdown @command="handleCommand">
              <div class="user-info">
                <el-avatar :size="32" :src="userStore.userAvatar">
                  {{ userStore.userName.charAt(0).toUpperCase() }}
                </el-avatar>
                <span class="username">{{ userStore.userName }}</span>
                <el-icon><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">
                    <el-icon><User /></el-icon>
                    个人中心
                  </el-dropdown-item>
                  <el-dropdown-item divided command="logout">
                    <el-icon><SwitchButton /></el-icon>
                    退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>

        <el-main class="main">
          <router-view v-slot="{ Component }">
            <transition name="fade-transform" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import {
  Fold, Expand, FullScreen, ArrowDown, User, SwitchButton, Refresh
} from '@element-plus/icons-vue'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)
const menuList = ref([])

const currentRoute = computed(() => route)
const activeMenu = computed(() => route.path)

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const toggleFullScreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}

const handleRefresh = () => {
  location.reload()
}

const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      ElMessageBox.alert('个人中心功能开发中...', '提示')
      break
    case 'logout':
      ElMessageBox.confirm('确定退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        userStore.logout()
        router.push('/login')
      })
      break
  }
}

const loadMenus = () => {
  // 从路由配置动态生成菜单
  // 这里使用硬编码菜单，也可以从后端API获取
  menuList.value = [
    {
      path: '/dashboard',
      meta: { title: '数据看板', icon: 'DataAnalysis' }
    },
    {
      path: '/system',
      meta: { title: '系统管理', icon: 'Setting' },
      children: [
        { path: '/user', meta: { title: '用户管理', icon: 'User' } },
        { path: '/counselor', meta: { title: '咨询师管理', icon: 'UserFilled' } },
        { path: '/roles', meta: { title: '角色管理', icon: 'Avatar' } },
        { path: '/permissions', meta: { title: '权限管理', icon: 'Key' } },
        { path: '/menus', meta: { title: '菜单管理', icon: 'Menu' } }
      ]
    },
    {
      path: '/business',
      meta: { title: '业务管理', icon: 'Briefcase' },
      children: [
        { path: '/order', meta: { title: '订单管理', icon: 'Document' } },
        { path: '/chat', meta: { title: '聊天记录', icon: 'ChatDotRound' } }
      ]
    },
    {
      path: '/finance',
      meta: { title: '财务管理', icon: 'Wallet' },
      children: [
        { path: '/withdraw', meta: { title: '提现审核', icon: 'Coin' } },
        { path: '/statistics', meta: { title: '财务统计', icon: 'DataLine' } },
        { path: '/reports', meta: { title: '财务报表', icon: 'TrendCharts' } }
      ]
    },
    {
      path: '/lowcode',
      meta: { title: '低代码平台', icon: 'MagicStick' },
      children: [
        { path: '/lowcode/forms', meta: { title: '表单设计', icon: 'Edit' } },
        { path: '/lowcode/pages', meta: { title: '页面设计', icon: 'Grid' } },
        { path: '/lowcode/data', meta: { title: '数据管理', icon: 'Database' } }
      ]
    }
  ]
}

onMounted(() => {
  loadMenus()
})
</script>

<style lang="scss" scoped>
.layout-container {
  width: 100%;
  height: 100vh;

  .sidebar {
    background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
    transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 2px 0 12px rgba(0, 0, 0, 0.1);

    .logo {
      height: 64px;
      line-height: 64px;
      text-align: center;
      font-size: 22px;
      font-weight: 700;
      background: rgba(102, 126, 234, 0.1);
      color: #fff;
      overflow: hidden;
      letter-spacing: 1px;
      transition: all 0.3s;

      &:hover {
        background: rgba(102, 126, 234, 0.2);
      }

      &.is-collapse {
        font-size: 16px;
      }
    }

    :deep(.el-menu) {
      background: transparent;
      border: none;
      padding: 10px 0;

      .el-menu-item,
      .el-sub-menu__title {
        color: #a8b2d1;
        border-radius: 8px;
        margin: 4px 12px;
        height: 48px;
        line-height: 48px;
        transition: all 0.3s;

        &:hover {
          background: rgba(102, 126, 234, 0.2);
          color: #fff;
          transform: translateX(4px);
        }
      }

      .el-menu-item.is-active {
        background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
        color: #fff;
        box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
      }

      .el-sub-menu {
        .el-menu-item {
          padding-left: 45px;
        }
      }
    }
  }

  .header {
    background: #fff;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 24px;
    height: 64px;
    position: relative;
    z-index: 10;

    &::before {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      height: 1px;
      background: linear-gradient(90deg, transparent 0%, #e5e7eb 50%, transparent 100%);
    }

    .header-left {
      display: flex;
      align-items: center;
      flex: 1;

      .collapse-btn {
        font-size: 20px;
        cursor: pointer;
        margin-right: 20px;
        color: #606266;
        transition: all 0.3s;
        padding: 8px;
        border-radius: 8px;

        &:hover {
          background: #f5f7fa;
          color: #409eff;
        }
      }
    }

    .header-right {
      display: flex;
      align-items: center;
      gap: 16px;

      .action-icon {
        font-size: 18px;
        cursor: pointer;
        color: #909399;
        padding: 8px;
        border-radius: 8px;
        transition: all 0.3s;

        &:hover {
          background: #f5f7fa;
          color: #409eff;
          transform: scale(1.1);
        }
      }

      .user-info {
        display: flex;
        align-items: center;
        cursor: pointer;
        padding: 4px 12px;
        border-radius: 20px;
        transition: all 0.3s;
        gap: 8px;

        &:hover {
          background: #f5f7fa;
        }

        .username {
          margin: 0;
          color: #303133;
          font-weight: 500;
          font-size: 14px;
        }
      }
    }
  }

  .main {
    background: #f5f7fa;
    padding: 24px;
    overflow-y: auto;
    min-height: calc(100vh - 64px);
  }
}

.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

// 滚动条样式
:deep(.el-aside) {
  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 3px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }
}

:deep(.el-main) {
  &::-webkit-scrollbar {
    width: 8px;
  }

  &::-webkit-scrollbar-thumb {
    background: #dcdfe6;
    border-radius: 4px;

    &:hover {
      background: #c0c4cc;
    }
  }

  &::-webkit-scrollbar-track {
    background: #f5f7fa;
  }
}
</style>
