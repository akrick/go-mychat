<template>
  <div class="profile-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>个人中心</span>
        </div>
      </template>

      <div class="profile-content">
        <div class="avatar-section">
          <el-upload
            class="avatar-uploader"
            :action="uploadUrl"
            :show-file-list="false"
            :on-success="handleAvatarSuccess"
            :before-upload="beforeAvatarUpload"
          >
            <el-avatar :size="100" :src="userStore.userAvatar || defaultAvatar">
              <el-icon :size="50" color="#fff"><Camera /></el-icon>
            </el-avatar>
            <template #tip>
              <div class="el-upload__tip">点击更换头像</div>
            </template>
          </el-upload>
        </div>

        <div class="info-section">
          <el-descriptions title="基本信息" :column="2" border>
            <el-descriptions-item label="用户名">
              {{ userStore.userName }}
            </el-descriptions-item>
            <el-descriptions-item label="用户ID">
              {{ userStore.userInfo?.id }}
            </el-descriptions-item>
            <el-descriptions-item label="邮箱">
              {{ userStore.userInfo?.email || '未设置' }}
            </el-descriptions-item>
            <el-descriptions-item label="手机号">
              {{ userStore.userInfo?.phone || '未设置' }}
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag type="success">正常</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="注册时间">
              {{ formatDate(userStore.userInfo?.created_at) }}
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="form-section">
          <el-divider>修改密码</el-divider>
          <el-form
            ref="passwordFormRef"
            :model="passwordForm"
            :rules="passwordRules"
            label-width="100px"
            style="max-width: 500px"
          >
            <el-form-item label="旧密码" prop="oldPassword">
              <el-input
                v-model="passwordForm.oldPassword"
                type="password"
                placeholder="请输入旧密码"
                show-password
              />
            </el-form-item>
            <el-form-item label="新密码" prop="newPassword">
              <el-input
                v-model="passwordForm.newPassword"
                type="password"
                placeholder="请输入新密码(6-20位)"
                show-password
              />
            </el-form-item>
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input
                v-model="passwordForm.confirmPassword"
                type="password"
                placeholder="请再次输入新密码"
                show-password
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleChangePassword" :loading="passwordLoading">
                修改密码
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useUserStore } from '@/stores/user'
import { Camera } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import dayjs from 'dayjs'

const userStore = useUserStore()

const passwordFormRef = ref(null)
const passwordLoading = ref(false)
const uploadUrl = import.meta.env.VITE_API_BASE_URL + '/admin/upload'
const defaultAvatar = 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxMDAiIGhlaWdodD0iMTAwIj48cmVjdCB3aWR0aD0iNjUiIGhlaWdodD0iNjUiIHJ4PSI1MDAiIHJ5PSI3MCIgZmlsbD0iI2QxMjZmZiIi8+PC9zdmc+'

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (_rule, value, callback) => {
  if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入旧密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度6-20位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const handleAvatarSuccess = (response) => {
  if (response.code === 200) {
    userStore.userInfo.avatar = response.data.url
    ElMessage.success('头像更新成功')
  }
}

const beforeAvatarUpload = (file) => {
  const isJPG = file.type === 'image/jpeg'
  const isPNG = file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG && !isPNG) {
    ElMessage.error('上传头像图片只能是JPG/PNG格式!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('上传头像图片大小不能超过2MB!')
    return false
  }
  return true
}

const handleChangePassword = async () => {
  await passwordFormRef.value.validate()

  try {
    passwordLoading.value = true
    await request.post('/admin/user/password', {
      old_password: passwordForm.oldPassword,
      new_password: passwordForm.newPassword
    })
    ElMessage.success('密码修改成功，请重新登录')
    passwordFormRef.value.resetFields()
  } catch (error) {
    console.error(error)
  } finally {
    passwordLoading.value = false
  }
}
</script>

<style lang="scss" scoped>
.profile-container {
  padding: 20px;

  .card-header {
    font-weight: bold;
    font-size: 16px;
  }

  .profile-content {
    .avatar-section {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 40px 0;
      border-bottom: 1px solid #ebeef5;
      margin-bottom: 30px;

      .avatar-uploader {
        :deep(.el-upload) {
          border: none;
          background: transparent;
        }

        :deep(.el-avatar) {
          border: 4px solid #e5e7eb;
          cursor: pointer;
          transition: all 0.3s;

          &:hover {
            border-color: #667eea;
            transform: scale(1.05);
          }
        }

        .el-upload__tip {
          margin-top: 10px;
          color: #909399;
          font-size: 13px;
        }
      }
    }
  }

  .info-section {
    margin-bottom: 30px;

    :deep(.el-descriptions) {
      .el-descriptions__label {
        font-weight: 500;
        color: #606266;
      }
    }
  }

  .form-section {
    max-width: 600px;

    .el-divider {
      margin: 30px 0 20px 0;
      font-size: 14px;
      font-weight: bold;
      color: #303133;
    }
  }
}
</style>
