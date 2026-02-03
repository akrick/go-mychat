<template>
  <div class="application-page">
    <AppHeader />

    <div class="page-header">
      <div class="container">
        <h1>咨询师入驻申请</h1>
        <p>加入我们，开启专业心理咨询之旅</p>
      </div>
    </div>

    <div class="container">
      <div class="application-wrapper">
        <!-- 步骤指示器 -->
        <el-steps :active="activeStep" finish-status="success" align-center>
          <el-step title="填写基本信息" />
          <el-step title="上传资质证书" />
          <el-step title="提交申请" />
        </el-steps>

        <div class="step-content">
          <!-- 步骤1: 填写基本信息 -->
          <el-card v-if="activeStep === 0" class="step-card">
            <h3>基本信息</h3>
            <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="真实姓名" prop="name">
                    <el-input v-model="form.name" placeholder="请输入真实姓名" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="性别" prop="gender">
                    <el-radio-group v-model="form.gender">
                      <el-radio label="男">男</el-radio>
                      <el-radio label="女">女</el-radio>
                    </el-radio-group>
                  </el-form-item>
                </el-col>
              </el-row>

              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="联系电话" prop="phone">
                    <el-input v-model="form.phone" placeholder="请输入联系电话" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="邮箱地址" prop="email">
                    <el-input v-model="form.email" placeholder="请输入邮箱地址" />
                  </el-form-item>
                </el-col>
              </el-row>

              <el-form-item label="职称" prop="title">
                <el-select v-model="form.title" placeholder="请选择职称" style="width: 100%">
                  <el-option label="国家二级心理咨询师" value="国家二级心理咨询师" />
                  <el-option label="国家三级心理咨询师" value="国家三级心理咨询师" />
                  <el-option label="心理治疗师" value="心理治疗师" />
                  <el-option label="精神科医师" value="精神科医师" />
                  <el-option label="临床心理学家" value="临床心理学家" />
                  <el-option label="其他" value="其他" />
                </el-select>
              </el-form-item>

              <el-form-item label="从业年限" prop="years_exp">
                <el-input-number
                  v-model="form.years_exp"
                  :min="1"
                  :max="50"
                  controls-position="right"
                  style="width: 200px"
                />
                <span style="margin-left: 10px">年</span>
              </el-form-item>

              <el-form-item label="擅长领域" prop="specialty">
                <el-select
                  v-model="specialtyTags"
                  multiple
                  placeholder="请选择擅长领域（可多选）"
                  style="width: 100%"
                >
                  <el-option label="情感咨询" value="情感咨询" />
                  <el-option label="婚姻家庭" value="婚姻家庭" />
                  <el-option label="亲子教育" value="亲子教育" />
                  <el-option label="职场发展" value="职场发展" />
                  <el-option label="情绪管理" value="情绪管理" />
                  <el-option label="人际关系" value="人际关系" />
                  <el-option label="青少年心理" value="青少年心理" />
                  <el-option label="焦虑抑郁" value="焦虑抑郁" />
                  <el-option label="睡眠障碍" value="睡眠障碍" />
                  <el-option label="性心理" value="性心理" />
                  <el-option label="创伤疗愈" value="创伤疗愈" />
                  <el-option label="其他" value="其他" />
                </el-select>
              </el-form-item>

              <el-form-item label="个人简介" prop="bio">
                <el-input
                  v-model="form.bio"
                  type="textarea"
                  :rows="6"
                  placeholder="请详细介绍您的专业背景、咨询风格、擅长问题等，至少50字"
                  maxlength="500"
                  show-word-limit
                />
              </el-form-item>
            </el-form>

            <div class="step-actions">
              <el-button type="primary" @click="nextStep">下一步</el-button>
            </div>
          </el-card>

          <!-- 步骤2: 上传资质证书 -->
          <el-card v-if="activeStep === 1" class="step-card">
            <h3>资质证书</h3>
            <p class="upload-tip">请上传您的心理咨询师资质证书，至少上传一张，最多三张</p>

            <el-upload
              v-model:file-list="certificateList"
              class="certificate-upload"
              :action="uploadUrl"
              :headers="uploadHeaders"
              list-type="picture-card"
              :on-success="handleUploadSuccess"
              :before-upload="beforeUpload"
              :on-remove="handleRemove"
              :limit="3"
              accept="image/jpeg,image/jpg,image/png"
            >
              <el-icon><Plus /></el-icon>
              <template #tip>
                <div class="el-upload__tip">
                  支持 JPG/JPEG/PNG 格式，单张图片不超过 5MB
                </div>
              </template>
            </el-upload>

            <div class="step-actions">
              <el-button @click="prevStep">上一步</el-button>
              <el-button type="primary" @click="nextStep">下一步</el-button>
            </div>
          </el-card>

          <!-- 步骤3: 确认提交 -->
          <el-card v-if="activeStep === 2" class="step-card">
            <h3>确认申请信息</h3>

            <div class="confirm-info">
              <el-descriptions :column="2" border>
                <el-descriptions-item label="真实姓名">{{ form.name }}</el-descriptions-item>
                <el-descriptions-item label="性别">{{ form.gender }}</el-descriptions-item>
                <el-descriptions-item label="联系电话">{{ form.phone }}</el-descriptions-item>
                <el-descriptions-item label="邮箱地址">{{ form.email }}</el-descriptions-item>
                <el-descriptions-item label="职称">{{ form.title }}</el-descriptions-item>
                <el-descriptions-item label="从业年限">{{ form.years_exp }}年</el-descriptions-item>
                <el-descriptions-item label="擅长领域" :span="2">
                  <el-tag
                    v-for="(tag, index) in specialtyTags"
                    :key="index"
                    style="margin-right: 5px; margin-bottom: 5px"
                  >
                    {{ tag }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="个人简介" :span="2">{{ form.bio }}</el-descriptions-item>
                <el-descriptions-item label="资质证书" :span="2">
                  <div class="certificate-preview">
                    <el-image
                      v-for="(cert, index) in certificateList"
                      :key="index"
                      :src="cert.url"
                      :preview-src-list="certificateList.map(c => c.url)"
                      style="width: 100px; height: 100px; margin-right: 10px"
                      fit="cover"
                    />
                  </div>
                </el-descriptions-item>
              </el-descriptions>
            </div>

            <el-alert
              type="warning"
              :closable="false"
              style="margin: 20px 0"
            >
              <template #title>
                <div>提交后，工作人员将在 1-3 个工作日内完成审核，请耐心等待。审核结果将通过系统通知。</div>
              </template>
            </el-alert>

            <div class="step-actions">
              <el-button @click="prevStep">上一步</el-button>
              <el-button type="primary" @click="handleSubmit" :loading="submitting">提交申请</el-button>
            </div>
          </el-card>

          <!-- 提交成功 -->
          <el-card v-if="applicationSuccess" class="step-card success-card">
            <el-result icon="success" title="申请提交成功" sub-title="我们将尽快审核您的申请，请耐心等待">
              <template #extra>
                <el-button type="primary" @click="goToProfile">查看我的申请</el-button>
                <el-button @click="$router.push('/')">返回首页</el-button>
              </template>
            </el-result>
          </el-card>
        </div>
      </div>
    </div>

    <AppFooter />
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { createApplication } from '@/api/application'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import { showSuccess, showError } from '@/utils/errorHandler'

const router = useRouter()
const userStore = useUserStore()

// 检查登录
if (!userStore.token) {
  ElMessage.warning('请先登录')
  router.push('/login')
}

// 步骤控制
const activeStep = ref(0)
const applicationSuccess = ref(false)

// 表单数据
const form = reactive({
  name: '',
  gender: '男',
  phone: '',
  email: userStore.userInfo?.email || '',
  title: '',
  years_exp: 1,
  bio: '',
  certificate_img1: '',
  certificate_img2: '',
  certificate_img3: ''
})

const specialtyTags = ref([])

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '姓名长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  title: [
    { required: true, message: '请选择职称', trigger: 'change' }
  ],
  years_exp: [
    { type: 'number', min: 1, max: 50, message: '从业年限在 1-50 年', trigger: 'blur' }
  ],
  specialty: [
    { required: true, message: '请选择擅长领域', trigger: 'change' }
  ],
  bio: [
    { required: true, message: '请输入个人简介', trigger: 'blur' },
    { min: 50, message: '个人简介至少 50 个字', trigger: 'blur' }
  ]
}

// 证书上传
const certificateList = ref([])
const uploadUrl = '/api/counselor/upload-certificate'
const uploadHeaders = computed(() => ({
  Authorization: 'Bearer ' + userStore.token
}))

const beforeUpload = (file) => {
  const isImage = file.type === 'image/jpeg' || file.type === 'image/jpg' || file.type === 'image/png'
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isImage) {
    ElMessage.error('只能上传 JPG/JPEG/PNG 格式的图片')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB')
    return false
  }
  return true
}

const handleUploadSuccess = (response, file, fileList) => {
  if (response.code === 200) {
    // 将上传成功的文件URL保存到表单
    const index = fileList.findIndex(f => f.uid === file.uid)
    if (index >= 0) {
      fileList[index].url = response.data.url
    }
    showSuccess('上传成功')
  } else {
    showError(response.msg || '上传失败')
    // 从列表中移除
    const index = fileList.findIndex(f => f.uid === file.uid)
    if (index >= 0) {
      fileList.splice(index, 1)
    }
  }
}

const handleRemove = (file, fileList) => {
  // 更新证书列表
  certificateList.value = fileList
}

// 步骤控制
const formRef = ref(null)

const nextStep = () => {
  if (activeStep.value === 0) {
    formRef.value.validate((valid) => {
      if (valid) {
        if (specialtyTags.value.length === 0) {
          showError('请选择擅长领域')
          return
        }
        form.specialty = specialtyTags.value.join(',')
        activeStep.value++
      }
    })
  } else if (activeStep.value === 1) {
    if (certificateList.value.length === 0) {
      showError('至少需要上传一张资质证书')
      return
    }
    // 保存证书URL
    form.certificate_img1 = certificateList.value[0]?.url || ''
    form.certificate_img2 = certificateList.value[1]?.url || ''
    form.certificate_img3 = certificateList.value[2]?.url || ''
    activeStep.value++
  }
}

const prevStep = () => {
  if (activeStep.value > 0) {
    activeStep.value--
  }
}

// 提交申请
const submitting = ref(false)

const handleSubmit = () => {
  ElMessageBox.confirm('确认提交入驻申请吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    submitting.value = true
    createApplication(form)
      .then(() => {
        showSuccess('申请提交成功，请等待审核')
        applicationSuccess.value = true
        submitting.value = false
      })
      .catch(error => {
        showError(error.msg || '提交失败')
        submitting.value = false
      })
  }).catch(() => {})
}

const goToProfile = () => {
  router.push('/profile')
}
</script>

<style scoped>
.application-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 60px 0;
  text-align: center;
}

.page-header h1 {
  font-size: 36px;
  margin-bottom: 10px;
}

.page-header p {
  font-size: 18px;
  opacity: 0.9;
}

.container {
  max-width: 1200px;
  margin: -30px auto 0;
  padding: 0 20px;
  position: relative;
  z-index: 1;
}

.application-wrapper {
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  padding: 40px;
  margin-bottom: 40px;
}

.step-content {
  margin-top: 40px;
}

.step-card {
  border: none;
  box-shadow: none;
}

.step-card h3 {
  font-size: 20px;
  margin-bottom: 20px;
  color: #333;
}

.upload-tip {
  color: #666;
  margin-bottom: 20px;
}

.certificate-upload {
  margin-bottom: 30px;
}

.step-actions {
  text-align: center;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.step-actions .el-button {
  margin: 0 10px;
  min-width: 120px;
}

.confirm-info {
  margin-bottom: 20px;
}

.certificate-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.success-card {
  text-align: center;
  padding: 40px;
}

@media (max-width: 768px) {
  .application-wrapper {
    padding: 20px;
  }

  .page-header h1 {
    font-size: 28px;
  }

  .step-actions .el-button {
    width: 100%;
    margin: 5px 0;
  }
}
</style>
