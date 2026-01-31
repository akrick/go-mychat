<template>
  <div class="image-upload-container">
    <el-upload
      class="avatar-uploader"
      :action="uploadUrl"
      :headers="uploadHeaders"
      :show-file-list="false"
      :on-success="handleSuccess"
      :on-error="handleError"
      :on-progress="handleProgress"
      :before-upload="beforeUpload"
      :accept="accept"
      name="file"
    >
      <div v-if="imageUrl" class="avatar-wrapper">
        <img :src="imageUrl" class="avatar" @error="handleImageError" @click="handlePreview" />
        <div v-if="!imageLoadSuccess" class="avatar-error">
          <el-icon><Picture /></el-icon>
        </div>
      </div>
      <div v-else class="avatar-uploader-icon">
        <el-icon v-if="uploading"><Loading /></el-icon>
        <el-icon v-else><Plus /></el-icon>
      </div>
    </el-upload>

    <div class="upload-actions" v-if="imageUrl && imageLoadSuccess">
      <el-button link type="primary" size="small" @click="handlePreview">预览</el-button>
      <el-button link type="danger" size="small" @click="handleRemove">删除</el-button>
    </div>
    <div class="upload-tips" v-if="tips">{{ tips }}</div>

    <!-- 图片预览对话框 -->
    <el-dialog v-model="previewVisible" title="图片预览" width="60%">
      <img :src="previewUrl" style="width: 100%; max-height: 70vh; object-fit: contain;" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Loading, Picture } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  tips: {
    type: String,
    default: ''
  },
  accept: {
    type: String,
    default: 'image/jpeg,image/jpg,image/png,image/webp'
  },
  maxSize: {
    type: Number,
    default: 2 // MB
  }
})

const emit = defineEmits(['update:modelValue'])

const userStore = useUserStore()

const uploadUrl = computed(() => {
  // 使用相对路径，由 Vite 代理处理转发
  return '/api/admin/upload/image'
})
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${userStore.token}`
}))

const imageUrl = ref(props.modelValue)
const uploading = ref(false)
const imageLoadSuccess = ref(true)
const previewVisible = ref(false)
const previewUrl = ref('')

watch(() => props.modelValue, (newVal) => {
  imageUrl.value = newVal
  imageLoadSuccess.value = true
})

const handleSuccess = (response, file) => {
  uploading.value = false
  console.log('Upload success:', response)

  if (response && response.code === 200) {
    // 处理返回的文件URL
    const fileUrl = response.data?.file_url || response.file_url
    if (fileUrl) {
      imageUrl.value = fileUrl
      imageLoadSuccess.value = true
      emit('update:modelValue', fileUrl)
      ElMessage.success('上传成功')
    } else {
      ElMessage.error('上传失败：未返回文件URL')
    }
  } else {
    ElMessage.error(response?.msg || response?.message || '上传失败')
  }
}

const handleError = (error) => {
  uploading.value = false
  console.error('Upload error:', error)
  console.error('Error details:', {
    status: error.status,
    message: error.message,
    response: error.response
  })

  if (error.status === 401) {
    ElMessage.error('未授权，请先登录')
  } else if (error.status === 403) {
    ElMessage.error('无权限上传文件')
  } else {
    ElMessage.error('上传失败：' + (error.message || '请重试'))
  }
}

const handleProgress = (event) => {
  uploading.value = true
  console.log('Upload progress:', event.percent)
}

const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLtMaxSize = file.size / 1024 / 1024 < props.maxSize

  if (!isImage) {
    ElMessage.error('只能上传图片文件！')
    return false
  }
  if (!isLtMaxSize) {
    ElMessage.error(`图片大小不能超过 ${props.maxSize}MB！`)
    return false
  }
  uploading.value = true
  return true
}

const handleRemove = () => {
  imageUrl.value = ''
  imageLoadSuccess.value = false
  emit('update:modelValue', '')
  ElMessage.success('已删除图片')
}

const handlePreview = () => {
  if (!imageUrl.value) return
  previewUrl.value = imageUrl.value
  previewVisible.value = true
}

const handleImageError = () => {
  console.error('Image load failed:', imageUrl.value)
  imageLoadSuccess.value = false
  ElMessage.error('图片加载失败')
}
</script>

<style scoped>
.image-upload-container {
  display: inline-block;
}

.avatar-uploader {
  display: inline-block;
}

.avatar-uploader :deep(.el-upload) {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.3s;
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-uploader :deep(.el-upload:hover) {
  border-color: #409eff;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-wrapper {
  position: relative;
  width: 100px;
  height: 100px;
  overflow: hidden;
}

.avatar {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
  cursor: pointer;
  transition: all 0.3s;
}

.avatar:hover {
  transform: scale(1.05);
}

.avatar-error {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
}

.avatar-error .el-icon {
  font-size: 32px;
}

.upload-actions {
  margin-top: 8px;
  display: flex;
  gap: 8px;
  justify-content: center;
}

.upload-tips {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
  line-height: 1.5;
  text-align: center;
}
</style>
