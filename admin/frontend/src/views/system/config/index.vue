<template>
  <div class="config-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>系统配置</span>
          <el-button type="primary" size="small" @click="handleSaveAll">
            <el-icon><Check /></el-icon>
            保存全部
          </el-button>
        </div>
      </template>

      <!-- 配置分类Tab -->
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本配置" name="basic">
          <el-form :model="basicConfig" label-width="120px">
            <el-form-item label="系统名称">
              <el-input v-model="basicConfig.system_name" placeholder="请输入系统名称" />
            </el-form-item>
            <el-form-item label="系统简介">
              <el-input
                v-model="basicConfig.system_description"
                type="textarea"
                :rows="3"
                placeholder="请输入系统简介"
              />
            </el-form-item>
            <el-form-item label="网站域名">
              <el-input v-model="basicConfig.site_domain" placeholder="请输入网站域名" />
            </el-form-item>
            <el-form-item label="联系电话">
              <el-input v-model="basicConfig.contact_phone" placeholder="请输入联系电话" />
            </el-form-item>
            <el-form-item label="联系邮箱">
              <el-input v-model="basicConfig.contact_email" placeholder="请输入联系邮箱" />
            </el-form-item>
            <el-form-item label="ICP备案号">
              <el-input v-model="basicConfig.icp_code" placeholder="请输入ICP备案号" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveConfig('basic')">保存基本配置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="用户配置" name="user">
          <el-form :model="userConfig" label-width="120px">
            <el-form-item label="允许注册">
              <el-switch v-model="userConfig.allow_register" />
              <span class="form-tip">关闭后新用户无法注册</span>
            </el-form-item>
            <el-form-item label="默认头像">
              <el-input v-model="userConfig.default_avatar" placeholder="请输入默认头像URL" />
            </el-form-item>
            <el-form-item label="密码最小长度">
              <el-input-number v-model="userConfig.password_min_length" :min="6" :max="20" />
            </el-form-item>
            <el-form-item label="账户锁定阈值">
              <el-input-number v-model="userConfig.login_max_attempts" :min="3" :max="10" />
              <span class="form-tip">连续登录失败次数</span>
            </el-form-item>
            <el-form-item label="锁定时长(分钟)">
              <el-input-number v-model="userConfig.login_lock_time" :min="5" :max="60" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveConfig('user')">保存用户配置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="聊天配置" name="chat">
          <el-form :model="chatConfig" label-width="120px">
            <el-form-item label="免费聊天时长">
              <el-input-number v-model="chatConfig.free_chat_duration" :min="0" :max="60" />
              <span class="form-tip">分钟</span>
            </el-form-item>
            <el-form-item label="计费周期(秒)">
              <el-input-number v-model="chatConfig.billing_cycle" :min="10" :max="300" />
            </el-form-item>
            <el-form-item label="会话超时(分钟)">
              <el-input-number v-model="chatConfig.session_timeout" :min="10" :max="1440" />
            </el-form-item>
            <el-form-item label="消息保留天数">
              <el-input-number v-model="chatConfig.message_retention_days" :min="7" :max="365" />
            </el-form-item>
            <el-form-item label="允许发图">
              <el-switch v-model="chatConfig.allow_image_message" />
            </el-form-item>
            <el-form-item label="允许发语音">
              <el-switch v-model="chatConfig.allow_voice_message" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveConfig('chat')">保存聊天配置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="支付配置" name="payment">
          <el-form :model="paymentConfig" label-width="120px">
            <el-form-item label="支付宝AppID">
              <el-input v-model="paymentConfig.alipay_app_id" placeholder="请输入支付宝AppID" />
            </el-form-item>
            <el-form-item label="支付宝公钥">
              <el-input
                v-model="paymentConfig.alipay_public_key"
                type="textarea"
                :rows="3"
                placeholder="请输入支付宝公钥"
              />
            </el-form-item>
            <el-form-item label="微信AppID">
              <el-input v-model="paymentConfig.wechat_app_id" placeholder="请输入微信AppID" />
            </el-form-item>
            <el-form-item label="微信商户号">
              <el-input v-model="paymentConfig.wechat_mch_id" placeholder="请输入微信商户号" />
            </el-form-item>
            <el-form-item label="微信API密钥">
              <el-input v-model="paymentConfig.wechat_api_key" type="password" placeholder="请输入微信API密钥" />
            </el-form-item>
            <el-form-item label="提现最小金额">
              <el-input-number v-model="paymentConfig.min_withdraw_amount" :min="1" :max="10000" />
              <span class="form-tip">元</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveConfig('payment')">保存支付配置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="通知配置" name="notification">
          <el-form :model="notificationConfig" label-width="120px">
            <el-form-item label="启用邮件通知">
              <el-switch v-model="notificationConfig.email_enabled" />
            </el-form-item>
            <el-form-item label="SMTP服务器">
              <el-input v-model="notificationConfig.smtp_host" placeholder="请输入SMTP服务器地址" />
            </el-form-item>
            <el-form-item label="SMTP端口">
              <el-input-number v-model="notificationConfig.smtp_port" :min="1" :max="65535" />
            </el-form-item>
            <el-form-item label="SMTP用户名">
              <el-input v-model="notificationConfig.smtp_username" placeholder="请输入SMTP用户名" />
            </el-form-item>
            <el-form-item label="SMTP密码">
              <el-input v-model="notificationConfig.smtp_password" type="password" placeholder="请输入SMTP密码" />
            </el-form-item>
            <el-form-item label="启用短信通知">
              <el-switch v-model="notificationConfig.sms_enabled" />
            </el-form-item>
            <el-form-item label="短信服务商">
              <el-select v-model="notificationConfig.sms_provider" placeholder="请选择短信服务商">
                <el-option label="阿里云" value="aliyun" />
                <el-option label="腾讯云" value="tencent" />
                <el-option label="华为云" value="huawei" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveConfig('notification')">保存通知配置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="存储配置" name="storage">
          <el-form :model="storageConfig" label-width="120px">
            <el-form-item label="存储方式">
              <el-radio-group v-model="storageConfig.storage_type">
                <el-radio label="local">本地存储</el-radio>
                <el-radio label="oss">阿里云OSS</el-radio>
                <el-radio label="cos">腾讯云COS</el-radio>
                <el-radio label="s3">AWS S3</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="上传路径">
              <el-input v-model="storageConfig.upload_path" placeholder="请输入文件上传路径" />
            </el-form-item>
            <el-form-item label="允许文件类型">
              <el-input v-model="storageConfig.allowed_file_types" placeholder="jpg,png,gif,pdf,doc" />
            </el-form-item>
            <el-form-item label="最大文件大小(MB)">
              <el-input-number v-model="storageConfig.max_file_size" :min="1" :max="100" />
            </el-form-item>
            <el-form-item label="图片压缩">
              <el-switch v-model="storageConfig.image_compress" />
            </el-form-item>
            <el-form-item label="压缩质量">
              <el-slider v-model="storageConfig.compress_quality" :min="10" :max="100" :disabled="!storageConfig.image_compress" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveConfig('storage')">保存存储配置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Check } from '@element-plus/icons-vue'
import { getSystemConfigs, updateSystemConfig } from '@/api/system'

const activeTab = ref('basic')
const loading = ref(false)

const basicConfig = reactive({
  system_name: '',
  system_description: '',
  site_domain: '',
  contact_phone: '',
  contact_email: '',
  icp_code: ''
})

const userConfig = reactive({
  allow_register: true,
  default_avatar: '',
  password_min_length: 6,
  login_max_attempts: 5,
  login_lock_time: 30
})

const chatConfig = reactive({
  free_chat_duration: 5,
  billing_cycle: 60,
  session_timeout: 30,
  message_retention_days: 90,
  allow_image_message: true,
  allow_voice_message: true
})

const paymentConfig = reactive({
  alipay_app_id: '',
  alipay_public_key: '',
  wechat_app_id: '',
  wechat_mch_id: '',
  wechat_api_key: '',
  min_withdraw_amount: 100
})

const notificationConfig = reactive({
  email_enabled: false,
  smtp_host: '',
  smtp_port: 465,
  smtp_username: '',
  smtp_password: '',
  sms_enabled: false,
  sms_provider: 'aliyun'
})

const storageConfig = reactive({
  storage_type: 'local',
  upload_path: '/uploads',
  allowed_file_types: 'jpg,jpeg,png,gif,pdf,doc,docx,xls,xlsx',
  max_file_size: 10,
  image_compress: true,
  compress_quality: 80
})

const loadConfigs = async () => {
  loading.value = true
  try {
    const res = await getSystemConfigs({ page: 1, pageSize: 100 })
    const configs = res.list || []

    // 加载配置到对应的表单
    configs.forEach(config => {
      const value = JSON.parse(config.value || '{}')

      switch (config.category) {
        case 'basic':
          Object.assign(basicConfig, value)
          break
        case 'user':
          Object.assign(userConfig, value)
          break
        case 'chat':
          Object.assign(chatConfig, value)
          break
        case 'payment':
          Object.assign(paymentConfig, value)
          break
        case 'notification':
          Object.assign(notificationConfig, value)
          break
        case 'storage':
          Object.assign(storageConfig, value)
          break
      }
    })
  } catch (error) {
    ElMessage.error(error.message || '加载配置失败')
  } finally {
    loading.value = false
  }
}

const handleSaveConfig = async (category) => {
  loading.value = true
  try {
    let configData = {}

    switch (category) {
      case 'basic':
        configData = basicConfig
        break
      case 'user':
        configData = userConfig
        break
      case 'chat':
        configData = chatConfig
        break
      case 'payment':
        configData = paymentConfig
        break
      case 'notification':
        configData = notificationConfig
        break
      case 'storage':
        configData = storageConfig
        break
    }

    // 转换为数组调用更新接口
    for (const [key, value] of Object.entries(configData)) {
      // TODO: 调用更新接口
      await updateSystemConfig(0, { key, value: JSON.stringify(value), category })
    }

    ElMessage.success('配置保存成功')
  } catch (error) {
    ElMessage.error(error.message || '保存配置失败')
  } finally {
    loading.value = false
  }
}

const handleSaveAll = async () => {
  loading.value = true
  try {
    await handleSaveConfig('basic')
    await handleSaveConfig('user')
    await handleSaveConfig('chat')
    await handleSaveConfig('payment')
    await handleSaveConfig('notification')
    await handleSaveConfig('storage')
    ElMessage.success('所有配置保存成功')
  } catch (error) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadConfigs()
})
</script>

<style lang="scss" scoped>
.config-container {
  padding: 20px;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: bold;
    font-size: 16px;
  }

  .form-tip {
    margin-left: 10px;
    font-size: 13px;
    color: #909399;
  }

  :deep(.el-form-item) {
    margin-bottom: 22px;
  }
}
</style>
