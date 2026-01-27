<template>
  <div class="login-container">
    <div class="login-bg"></div>
    <div class="login-box">
      <div class="login-header">
        <div class="logo">
          <el-icon :size="40" color="#667eea"><ChatDotRound /></el-icon>
        </div>
        <h2>MyChat管理后台</h2>
        <p>专业的在线心理咨询管理平台</p>
      </div>

      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            size="large"
            clearable
            :prefix-icon="User"
          >
          </el-input>
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            show-password
            :prefix-icon="Lock"
          >
          </el-input>
        </el-form-item>

        <el-form-item>
          <div class="form-options">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleLogin"
            class="login-btn"
          >
            <template v-if="!loading">
              <span>登 录</span>
            </template>
            <template v-else>
              <span>登录中...</span>
            </template>
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <div class="tips">
          <el-icon><InfoFilled /></el-icon>
          <span>默认账号: admin / admin123</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { User, Lock, ChatDotRound, InfoFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loginFormRef = ref(null)
const loading = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度3-20字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度6-20字符', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  try {
    await loginFormRef.value.validate()

    loading.value = true
    
    // 1. 登录
    console.log('开始登录...')
    await userStore.login(loginForm)
    console.log('登录成功')
    
    // 2. 获取用户信息（失败不影响跳转）
    try {
      await userStore.getUserInfo()
    } catch (e) {
      console.warn('获取用户信息失败，但不影响登录:', e)
    }
    
    // 3. 获取权限（失败不影响跳转，已有默认值）
    try {
      await userStore.getPermissions()
    } catch (e) {
      console.warn('获取权限失败，使用默认权限:', e)
    }

    if (rememberMe.value) {
      localStorage.setItem('remember_username', loginForm.username)
    }

    const redirect = route.query.redirect || '/'
    console.log('准备跳转到:', redirect)
    
    // 4. 跳转到目标页面
    await router.push(redirect)
    console.log('跳转完成')

    ElMessage.success('登录成功！欢迎回来')
  } catch (error) {
    loading.value = false
    console.error('登录错误详情:', error)
    // 错误已经在 request.js 拦截器中处理，这里不需要重复处理
  }
}

onMounted(() => {
  // 自动填充记住的用户名
  const rememberedUsername = localStorage.getItem('remember_username')
  if (rememberedUsername) {
    loginForm.username = rememberedUsername
    rememberMe.value = true
  }

  // 添加页面加载动画
  const loginBox = document.querySelector('.login-box')
  loginBox?.classList.add('animate-in')
})
</script>

<style lang="scss" scoped>
.login-container {
  position: relative;
  width: 100%;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;

  .login-bg {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    opacity: 0.1;
    background-image:
      radial-gradient(circle at 20% 50%, rgba(255, 255, 255, 0.1) 0%, transparent 50%),
      radial-gradient(circle at 80% 20%, rgba(255, 255, 255, 0.1) 0%, transparent 50%),
      radial-gradient(circle at 40% 80%, rgba(255, 255, 255, 0.1) 0%, transparent 50%);
  }

  .login-box {
    position: relative;
    width: 420px;
    padding: 50px 45px;
    background: rgba(255, 255, 255, 0.95);
    border-radius: 16px;
    box-shadow:
      0 20px 60px rgba(0, 0, 0, 0.3),
      0 10px 20px rgba(0, 0, 0, 0.2);
    backdrop-filter: blur(10px);
    animation: fadeInUp 0.6s ease-out;

    &.animate-in {
      animation: fadeInUp 0.6s ease-out;
    }

    .login-header {
      text-align: center;
      margin-bottom: 35px;

      .logo {
        margin-bottom: 20px;
        display: flex;
        justify-content: center;
        align-items: center;
      }

      h2 {
        font-size: 28px;
        font-weight: 600;
        color: #333;
        margin-bottom: 12px;
        letter-spacing: 0.5px;
      }

      p {
        color: #666;
        font-size: 14px;
        margin: 0;
      }
    }

    .login-form {
      :deep(.el-input__wrapper) {
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        border-radius: 8px;

        &:hover {
          box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
        }
      }

      .form-options {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;

        :deep(.el-checkbox__label) {
          color: #666;
        }
      }

      .login-btn {
        width: 100%;
        height: 46px;
        font-size: 16px;
        font-weight: 500;
        border-radius: 8px;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        border: none;
        transition: all 0.3s;

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
        }

        &:active {
          transform: translateY(0);
        }
      }
    }

    .login-footer {
      margin-top: 25px;
      padding-top: 20px;
      border-top: 1px solid #eee;

      .tips {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        color: #667eea;
        font-size: 13px;
        background: rgba(102, 126, 234, 0.1);
        padding: 10px 16px;
        border-radius: 6px;
      }
    }
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 响应式设计
@media (max-width: 768px) {
  .login-container {
    padding: 20px;

    .login-box {
      width: 100%;
      padding: 35px 25px;
    }
  }
}
</style>
