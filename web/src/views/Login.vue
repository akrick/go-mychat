<template>
  <div class="login-page">
    <AppHeader />
    <div class="login-container">
      <el-card class="login-card">
        <template #header>
          <div class="card-header">
            <h2>心理咨询平台</h2>
            <p>用户登录</p>
          </div>
        </template>
        <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-width="0">
          <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名"
            size="large"
            prefix-icon="User"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            size="large"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleLogin"
            style="width: 100%"
          >
            登录
          </el-button>
        </el-form-item>
        <div class="login-footer">
          <span>还没有账号？</span>
          <router-link to="/register">立即注册</router-link>
        </div>
      </el-form>
      </el-card>
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '@/api/auth'
import { useUserStore } from '@/stores/user'
import { handleError, showSuccess } from '@/utils/errorHandler'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'

const router = useRouter()
const userStore = useUserStore()
const loginFormRef = ref()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  await loginFormRef.value.validate()
  loading.value = true
  try {
    const res = await login(loginForm)
    userStore.setToken(res.data.token)
    userStore.setUserInfo(res.data.user)
    showSuccess('登录成功')
    router.push('/')
  } catch (error) {
    handleError(error, '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.login-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 400px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.card-header {
  text-align: center;
}

.card-header h2 {
  margin: 0 0 8px 0;
  color: #333;
}

.card-header p {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.login-footer {
  text-align: center;
  margin-top: 16px;
  font-size: 14px;
}

.login-footer a {
  color: #409eff;
  text-decoration: none;
  margin-left: 8px;
}

.login-footer a:hover {
  text-decoration: underline;
}
</style>
