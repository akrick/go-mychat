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
          <el-form :model="paymentConfig" label-width="150px">
            <el-divider content-position="left">支付宝配置</el-divider>
            <el-form-item label="支付宝AppID">
              <el-input v-model="paymentConfig.alipay_app_id" placeholder="请输入支付宝AppID" />
            </el-form-item>
            <el-form-item label="支付宝私钥">
              <el-input
                v-model="paymentConfig.alipay_private_key"
                type="textarea"
                :rows="3"
                placeholder="请输入应用私钥"
              />
            </el-form-item>
            <el-form-item label="支付宝公钥">
              <el-input
                v-model="paymentConfig.alipay_public_key"
                type="textarea"
                :rows="3"
                placeholder="请输入支付宝公钥"
              />
            </el-form-item>
            <el-form-item label="异步通知地址">
              <el-input v-model="paymentConfig.alipay_notify_url" placeholder="http://your-domain.com/api/payment/alipay/notify" />
            </el-form-item>
            <el-form-item label="同步跳转地址">
              <el-input v-model="paymentConfig.alipay_return_url" placeholder="http://your-domain.com/payment/success" />
            </el-form-item>
            <el-form-item label="沙箱环境">
              <el-switch v-model="paymentConfig.alipay_sandbox" />
              <span class="form-tip">开启后使用沙箱环境测试</span>
            </el-form-item>

            <el-divider content-position="left">微信支付配置</el-divider>
            <el-form-item label="微信AppID">
              <el-input v-model="paymentConfig.wechat_app_id" placeholder="请输入微信AppID" />
            </el-form-item>
            <el-form-item label="微信商户号">
              <el-input v-model="paymentConfig.wechat_mch_id" placeholder="请输入微信商户号" />
            </el-form-item>
            <el-form-item label="微信API密钥">
              <el-input v-model="paymentConfig.wechat_api_key" type="password" placeholder="请输入微信API密钥" />
            </el-form-item>
            <el-form-item label="微信证书路径">
              <el-input v-model="paymentConfig.wechat_cert_path" placeholder="/path/to/cert/apiclient_cert.p12" />
            </el-form-item>
            <el-form-item label="异步通知地址">
              <el-input v-model="paymentConfig.wechat_notify_url" placeholder="http://your-domain.com/api/payment/wechat/notify" />
            </el-form-item>
            <el-form-item label="沙箱环境">
              <el-switch v-model="paymentConfig.wechat_sandbox" />
              <span class="form-tip">开启后使用沙箱环境测试</span>
            </el-form-item>

            <el-divider content-position="left">提现配置</el-divider>
            <el-form-item label="提现最小金额">
              <el-input-number v-model="paymentConfig.min_withdraw_amount" :min="1" :max="10000" :step="10" />
              <span class="form-tip">元</span>
            </el-form-item>
            <el-form-item label="提现手续费率">
              <el-input-number v-model="paymentConfig.withdraw_fee_rate" :min="0" :max="100" :step="0.1" :precision="2" />
              <span class="form-tip">%</span>
            </el-form-item>
            <el-form-item label="提现到账时间">
              <el-input v-model="paymentConfig.withdraw_arrival_time" placeholder="T+1" />
              <span class="form-tip">如：T+1、T+0</span>
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

        <el-tab-pane label="微信公众号" name="wechat_mp">
          <el-form :model="wechatMpConfig" label-width="150px">
            <el-divider content-position="left">基本配置</el-divider>
            <el-form-item label="AppID">
              <el-input v-model="wechatMpConfig.appid" placeholder="请输入公众号AppID" />
            </el-form-item>
            <el-form-item label="AppSecret">
              <el-input v-model="wechatMpConfig.appsecret" type="password" placeholder="请输入公众号AppSecret" />
            </el-form-item>
            <el-form-item label="原始ID">
              <el-input v-model="wechatMpConfig.original_id" placeholder="gh_xxxxx" />
            </el-form-item>
            <el-form-item label="公众号名称">
              <el-input v-model="wechatMpConfig.name" placeholder="请输入公众号名称" />
            </el-form-item>

            <el-divider content-position="left">服务器配置</el-divider>
            <el-form-item label="Token">
              <el-input v-model="wechatMpConfig.token" placeholder="服务器配置Token" />
            </el-form-item>
            <el-form-item label="EncodingAESKey">
              <el-input v-model="wechatMpConfig.encoding_aes_key" type="textarea" :rows="2" placeholder="消息加解密密钥" />
            </el-form-item>
            <el-form-item label="消息加解密方式">
              <el-radio-group v-model="wechatMpConfig.msg_encrypt_mode">
                <el-radio label="plain">明文模式</el-radio>
                <el-radio label="aes">兼容模式</el-radio>
                <el-radio label="safe">安全模式</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="服务器URL">
              <el-input v-model="wechatMpConfig.server_url" placeholder="http://your-domain.com/api/wechat/mp/callback" />
            </el-form-item>

            <el-divider content-position="left">业务配置</el-divider>
            <el-form-item label="关注欢迎语">
              <el-input v-model="wechatMpConfig.welcome_msg" type="textarea" :rows="3" placeholder="用户关注时的欢迎语" />
            </el-form-item>
            <el-form-item label="默认回复">
              <el-input v-model="wechatMpConfig.default_reply" type="textarea" :rows="3" placeholder="无法匹配关键词时的默认回复" />
            </el-form-item>
            <el-form-item label="启用自动回复">
              <el-switch v-model="wechatMpConfig.auto_reply_enabled" />
            </el-form-item>
            <el-form-item label="启用菜单">
              <el-switch v-model="wechatMpConfig.menu_enabled" />
            </el-form-item>

            <el-divider content-position="left">JS-SDK配置</el-divider>
            <el-form-item label="启用JS-SDK">
              <el-switch v-model="wechatMpConfig.js_sdk_enabled" />
            </el-form-item>
            <el-form-item label="JS接口安全域名">
              <el-input v-model="wechatMpConfig.js_domain" placeholder="your-domain.com" />
            </el-form-item>
            <el-form-item label="授权域名">
              <el-input v-model="wechatMpConfig.oauth_domain" placeholder="your-domain.com" />
            </el-form-item>

            <el-divider content-position="left">模板消息</el-divider>
            <el-form-item label="启用模板消息">
              <el-switch v-model="wechatMpConfig.template_enabled" />
            </el-form-item>
            <el-form-item label="订单通知模板ID">
              <el-input v-model="wechatMpConfig.template_order" placeholder="订单模板ID" />
            </el-form-item>
            <el-form-item label="支付成功模板ID">
              <el-input v-model="wechatMpConfig.template_payment" placeholder="支付成功模板ID" />
            </el-form-item>
            <el-form-item label="提现通知模板ID">
              <el-input v-model="wechatMpConfig.template_withdraw" placeholder="提现通知模板ID" />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="handleSaveConfig('wechat_mp')">保存公众号配置</el-button>
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
import { getSystemConfigs, createSystemConfig, updateSystemConfig } from '@/api/system'

const activeTab = ref('basic')
const loading = ref(false)

// 配置标签映射
const configLabels = {
  system_name: '系统名称',
  system_description: '系统简介',
  site_domain: '网站域名',
  contact_phone: '联系电话',
  contact_email: '联系邮箱',
  icp_code: 'ICP备案号',
  allow_register: '允许注册',
  default_avatar: '默认头像',
  password_min_length: '密码最小长度',
  login_max_attempts: '账户锁定阈值',
  login_lock_time: '锁定时长(分钟)',
  free_chat_duration: '免费聊天时长',
  billing_cycle: '计费周期(秒)',
  session_timeout: '会话超时(分钟)',
  message_retention_days: '消息保留天数',
  allow_image_message: '允许发图',
  allow_voice_message: '允许发语音',
  alipay_app_id: '支付宝AppID',
  alipay_private_key: '支付宝私钥',
  alipay_public_key: '支付宝公钥',
  alipay_notify_url: '支付宝异步通知地址',
  alipay_return_url: '支付宝同步跳转地址',
  alipay_sandbox: '支付宝沙箱环境',
  wechat_app_id: '微信AppID',
  wechat_mch_id: '微信商户号',
  wechat_api_key: '微信API密钥',
  wechat_cert_path: '微信证书路径',
  wechat_notify_url: '微信异步通知地址',
  wechat_sandbox: '微信沙箱环境',
  min_withdraw_amount: '提现最小金额',
  withdraw_fee_rate: '提现手续费率',
  withdraw_arrival_time: '提现到账时间',
  email_enabled: '启用邮件通知',
  smtp_host: 'SMTP服务器',
  smtp_port: 'SMTP端口',
  smtp_username: 'SMTP用户名',
  smtp_password: 'SMTP密码',
  sms_enabled: '启用短信通知',
  sms_provider: '短信服务商',
  storage_type: '存储方式',
  upload_path: '上传路径',
  allowed_file_types: '允许文件类型',
  max_file_size: '最大文件大小(MB)',
  image_compress: '图片压缩',
  compress_quality: '压缩质量',
  // 微信公众号配置
  appid: '公众号AppID',
  appsecret: '公众号AppSecret',
  original_id: '原始ID',
  name: '公众号名称',
  token: 'Token',
  encoding_aes_key: '消息加解密密钥',
  msg_encrypt_mode: '消息加解密方式',
  server_url: '服务器URL',
  welcome_msg: '关注欢迎语',
  default_reply: '默认回复',
  auto_reply_enabled: '启用自动回复',
  menu_enabled: '启用菜单',
  js_sdk_enabled: '启用JS-SDK',
  js_domain: 'JS接口安全域名',
  oauth_domain: '授权域名',
  template_enabled: '启用模板消息',
  template_order: '订单通知模板ID',
  template_payment: '支付成功模板ID',
  template_withdraw: '提现通知模板ID'
}

const basicConfig = reactive({
  system_name: '我的聊天系统',
  system_description: '专业的在线心理咨询平台',
  site_domain: 'http://localhost:5173',
  contact_phone: '400-888-8888',
  contact_email: 'support@example.com',
  icp_code: '京ICP备12345678号'
})

const userConfig = reactive({
  allow_register: true,
  default_avatar: '/uploads/default-avatar.png',
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
  alipay_private_key: '',
  alipay_public_key: '',
  alipay_notify_url: '',
  alipay_return_url: '',
  alipay_sandbox: false,
  wechat_app_id: '',
  wechat_mch_id: '',
  wechat_api_key: '',
  wechat_cert_path: '',
  wechat_notify_url: '',
  wechat_sandbox: false,
  min_withdraw_amount: 100,
  withdraw_fee_rate: 1,
  withdraw_arrival_time: 'T+1'
})

const notificationConfig = reactive({
  email_enabled: false,
  smtp_host: 'smtp.example.com',
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

const wechatMpConfig = reactive({
  appid: '',
  appsecret: '',
  original_id: '',
  name: '',
  token: '',
  encoding_aes_key: '',
  msg_encrypt_mode: 'plain',
  server_url: '',
  welcome_msg: '欢迎关注我们的公众号！',
  default_reply: '抱歉，我没有理解您的消息，请输入关键词或联系客服。',
  auto_reply_enabled: true,
  menu_enabled: true,
  js_sdk_enabled: false,
  js_domain: '',
  oauth_domain: '',
  template_enabled: false,
  template_order: '',
  template_payment: '',
  template_withdraw: ''
})

// 存储已有的配置项ID
const configIds = reactive({})

const loadConfigs = async () => {
  loading.value = true
  try {
    const res = await getSystemConfigs({ page: 1, pageSize: 100 })
    const configs = res.list || []

    // 加载配置到对应的表单，并记录配置ID
    configs.forEach(config => {
      const key = config.key
      configIds[key] = config.id

      try {
        const value = JSON.parse(config.value || '{}')

        switch (config.category) {
          case 'basic':
            if (value !== null && typeof value === 'object') {
              basicConfig[key] = value
            } else {
              basicConfig[key] = value
            }
            break
          case 'user':
            if (value !== null && typeof value === 'object') {
              userConfig[key] = value
            } else {
              userConfig[key] = value
            }
            break
          case 'chat':
            if (value !== null && typeof value === 'object') {
              chatConfig[key] = value
            } else {
              chatConfig[key] = value
            }
            break
          case 'payment':
            if (value !== null && typeof value === 'object') {
              paymentConfig[key] = value
            } else {
              paymentConfig[key] = value
            }
            break
          case 'notification':
            if (value !== null && typeof value === 'object') {
              notificationConfig[key] = value
            } else {
              notificationConfig[key] = value
            }
            break
          case 'storage':
            if (value !== null && typeof value === 'object') {
              storageConfig[key] = value
            } else {
              storageConfig[key] = value
            }
            break
          case 'wechat_mp':
            if (value !== null && typeof value === 'object') {
              wechatMpConfig[key] = value
            } else {
              wechatMpConfig[key] = value
            }
            break
        }
      } catch (e) {
        // 如果解析失败，直接使用字符串值
        switch (config.category) {
          case 'basic':
            basicConfig[key] = config.value
            break
          case 'user':
            userConfig[key] = config.value
            break
          case 'chat':
            chatConfig[key] = config.value
            break
          case 'payment':
            paymentConfig[key] = config.value
            break
          case 'notification':
            notificationConfig[key] = config.value
            break
          case 'storage':
            storageConfig[key] = config.value
            break
          case 'wechat_mp':
            wechatMpConfig[key] = config.value
            break
        }
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
    let configCategory = ''

    switch (category) {
      case 'basic':
        configData = basicConfig
        configCategory = 'basic'
        break
      case 'user':
        configData = userConfig
        configCategory = 'user'
        break
      case 'chat':
        configData = chatConfig
        configCategory = 'chat'
        break
      case 'payment':
        configData = paymentConfig
        configCategory = 'payment'
        break
      case 'notification':
        configData = notificationConfig
        configCategory = 'notification'
        break
      case 'storage':
        configData = storageConfig
        configCategory = 'storage'
        break
      case 'wechat_mp':
        configData = wechatMpConfig
        configCategory = 'wechat_mp'
        break
    }

    // 保存每个配置项
    for (const [key, value] of Object.entries(configData)) {
      const configKey = `${category}.${key}`
      const label = configLabels[key] || key

      if (configIds[configKey]) {
        // 更新现有配置
        await updateSystemConfig(configIds[configKey], {
          key: configKey,
          value: JSON.stringify(value),
          category: configCategory
        })
      } else {
        // 创建新配置
        const res = await createSystemConfig({
          key: configKey,
          value: JSON.stringify(value),
          category: configCategory,
          label: label,
          type: typeof value === 'boolean' ? 'switch' : typeof value === 'number' ? 'number' : 'string',
          is_system: true,
          sort: 0,
          remark: '系统配置'
        })
        configIds[configKey] = res.data?.id || res.id
      }
    }

    ElMessage.success('配置保存成功')
  } catch (error) {
    console.error('保存配置错误:', error)
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
    await handleSaveConfig('wechat_mp')
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
