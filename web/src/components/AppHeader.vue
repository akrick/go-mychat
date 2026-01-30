<template>
  <div class="app-header">
    <div class="header-top">
      <div class="container">
        <div class="header-left">
          <span class="slogan">专业心理咨询平台，用心守护您的心理健康</span>
        </div>
        <div class="header-right">
          <span class="phone">
            <el-icon><Phone /></el-icon>
            客服热线：400-888-9999
          </span>
        </div>
      </div>
    </div>
    <div class="header-main">
      <div class="container">
        <div class="logo" @click="$router.push('/')">
          <el-icon size="40" color="#409eff"><ChatDotRound /></el-icon>
          <div class="logo-text">
            <span class="logo-name">心理咨询</span>
            <span class="logo-en">Psychological Counseling</span>
          </div>
        </div>
        <nav class="nav-menu">
          <router-link to="/" class="nav-item" active-class="active">首页</router-link>
          <router-link to="/counselors" class="nav-item">咨询师</router-link>
          <router-link to="/services" class="nav-item">服务介绍</router-link>
          <router-link to="/about" class="nav-item">关于我们</router-link>
        </nav>
        <div class="header-actions">
          <div v-if="userStore.token" class="user-login">
            <el-dropdown @command="handleCommand">
              <div class="user-info">
                <el-avatar :size="32" :src="userStore.userInfo?.avatar">
                  <el-icon><User /></el-icon>
                </el-avatar>
                <span class="username">{{ userStore.userInfo?.username || '用户' }}</span>
                <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="orders">
                    <el-icon><Document /></el-icon>
                    我的订单
                  </el-dropdown-item>
                  <el-dropdown-item command="profile">
                    <el-icon><UserFilled /></el-icon>
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
          <div v-else class="login-actions">
            <router-link to="/login" class="btn-login">登录</router-link>
            <router-link to="/register" class="btn-register">注册</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import {
  Phone,
  ChatDotRound,
  User,
  ArrowDown,
  Document,
  UserFilled,
  SwitchButton
} from '@element-plus/icons-vue'
import { showSuccess } from '@/utils/errorHandler'

const userStore = useUserStore()
const router = useRouter()

const handleCommand = (command) => {
  console.log('Dropdown command:', command)
  if (command === 'orders') {
    router.push('/orders')
  } else if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      userStore.logout()
      showSuccess('已退出登录')
      router.push('/')
    }).catch(() => {
      // 用户取消
    })
  }
}
</script>

<style scoped>
.app-header {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.header-top {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 8px 0;
  font-size: 14px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.header-top .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-main {
  padding: 16px 0;
}

.header-main .container {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  flex-shrink: 0;
}

.logo-text {
  display: flex;
  flex-direction: column;
}

.logo-name {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  line-height: 1.2;
}

.logo-en {
  font-size: 10px;
  color: #999;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.nav-menu {
  display: flex;
  gap: 32px;
  flex: 1;
  justify-content: center;
}

.nav-item {
  font-size: 15px;
  color: #666;
  text-decoration: none;
  position: relative;
  padding: 8px 0;
  transition: color 0.3s;
  white-space: nowrap;
}

.nav-item:hover {
  color: #409eff;
}

.nav-item.active {
  color: #409eff;
  font-weight: 500;
}

.nav-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: #409eff;
}

.header-actions {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  cursor: pointer;
  border-radius: 20px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f7fa;
}

.username {
  font-size: 14px;
  color: #333;
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dropdown-icon {
  font-size: 12px;
  color: #999;
}

.login-actions {
  display: flex;
  gap: 12px;
}

.btn-login {
  padding: 8px 20px;
  color: #409eff;
  border: 1px solid #409eff;
  border-radius: 20px;
  text-decoration: none;
  transition: all 0.3s;
  font-size: 14px;
}

.btn-login:hover {
  background: #409eff;
  color: white;
}

.btn-register {
  padding: 8px 20px;
  color: white;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  text-decoration: none;
  transition: all 0.3s;
  font-size: 14px;
}

.btn-register:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.phone {
  display: flex;
  align-items: center;
  gap: 6px;
}

@media (max-width: 992px) {
  .nav-menu {
    gap: 20px;
  }

  .nav-item {
    font-size: 14px;
  }
}

@media (max-width: 768px) {
  .header-top {
    display: none;
  }

  .nav-menu {
    gap: 12px;
  }

  .nav-item {
    font-size: 13px;
  }

  .logo-name {
    font-size: 20px;
  }

  .username {
    display: none;
  }
}
</style>
