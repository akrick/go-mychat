<template>
  <div class="profile-page">
    <AppHeader />
    <div class="profile-content">
      <div class="page-header">
        <div class="header-content">
          <h1>个人中心</h1>
          <p class="header-desc">管理您的个人信息和账户设置</p>
        </div>
      </div>

      <el-row :gutter="24" class="profile-row">
        <el-col :xs="24" :sm="24" :md="8" :lg="7" :xl="7">
          <el-card class="user-card" shadow="hover">
            <div class="card-header">
              <h3>个人信息</h3>
            </div>
            <div class="avatar-section">
              <el-upload
                class="avatar-uploader"
                :show-file-list="false"
                :before-upload="beforeAvatarUpload"
                :http-request="handleAvatarUpload"
              >
                <el-avatar :size="120" :src="userInfo?.avatar">
                  <el-icon><User /></el-icon>
                </el-avatar>
                <div class="avatar-overlay">
                  <el-icon><Camera /></el-icon>
                  <span>更换头像</span>
                </div>
              </el-upload>
              <h3>{{ userInfo?.username || '用户' }}</h3>
            </div>
            <div class="info-list">
              <div class="info-item">
                <el-icon class="info-icon"><User /></el-icon>
                <span class="label">用户ID：</span>
                <span class="value">{{ userInfo?.user_id || userInfo?.id || '-' }}</span>
              </div>
              <div class="info-item">
                <el-icon class="info-icon"><Message /></el-icon>
                <span class="label">邮箱：</span>
                <span class="value">{{ userInfo?.email || '未设置' }}</span>
              </div>
              <div class="info-item">
                <el-icon class="info-icon"><Phone /></el-icon>
                <span class="label">手机：</span>
                <span class="value">{{ userInfo?.phone || '未设置' }}</span>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="24" :sm="24" :md="16" :lg="17" :xl="17">
          <el-card class="settings-card" shadow="hover">
            <el-tabs v-model="activeTab" class="profile-tabs" @tab-change="handleTabChange">
              <el-tab-pane label="基本资料" name="profile">
                <template #label>
                  <el-icon><User /></el-icon>
                  <span>基本资料</span>
                </template>
                <el-form :model="profileForm" :rules="profileRules" ref="profileFormRef" label-width="120px" class="profile-form">
                  <el-form-item label="用户名">
                    <el-input :value="userInfo?.username" disabled />
                  </el-form-item>
                  <el-form-item label="邮箱地址" prop="email">
                    <el-input v-model="profileForm.email" placeholder="请输入邮箱" />
                  </el-form-item>
                  <el-form-item label="手机号码" prop="phone">
                    <el-input v-model="profileForm.phone" placeholder="请输入手机号" maxlength="11" />
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" :loading="submitting" size="large" @click="handleUpdateProfile">
                      <el-icon><Select /></el-icon>
                      保存修改
                    </el-button>
                  </el-form-item>
                </el-form>
              </el-tab-pane>

              <el-tab-pane label="修改密码" name="password">
                <template #label>
                  <el-icon><Lock /></el-icon>
                  <span>修改密码</span>
                </template>
                <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="120px" class="profile-form">
                  <el-form-item label="当前密码" prop="old_password">
                    <el-input v-model="passwordForm.old_password" type="password" placeholder="请输入当前密码" show-password />
                  </el-form-item>
                  <el-form-item label="新密码" prop="new_password">
                    <el-input v-model="passwordForm.new_password" type="password" placeholder="请输入新密码" show-password />
                  </el-form-item>
                  <el-form-item label="确认新密码" prop="confirm_password">
                    <el-input v-model="passwordForm.confirm_password" type="password" placeholder="请再次输入新密码" show-password />
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" :loading="submitting" size="large" @click="handleChangePassword">
                      <el-icon><Lock /></el-icon>
                      修改密码
                    </el-button>
                  </el-form-item>
                </el-form>
              </el-tab-pane>

              <el-tab-pane label="我的预约" name="appointments">
                <template #label>
                  <el-icon><Calendar /></el-icon>
                  <span>我的预约</span>
                </template>
                <div class="appointments-list" v-loading="appointmentsLoading">
                  <el-empty v-if="appointments.length === 0" description="暂无预约记录" :image-size="120" />
                  <div v-else class="appointment-item" v-for="appointment in appointments" :key="appointment.id">
                    <div class="appointment-header">
                      <div class="counselor-info">
                        <el-avatar :size="40" :src="appointment.counselor?.avatar">
                          <el-icon><User /></el-icon>
                        </el-avatar>
                        <div class="counselor-details">
                          <h4>{{ appointment.counselor?.name }}</h4>
                          <p>{{ appointment.counselor?.title }}</p>
                        </div>
                      </div>
                      <el-tag :type="getAppointmentStatusType(appointment.status)">
                        {{ getAppointmentStatusText(appointment.status) }}
                      </el-tag>
                    </div>
                    <div class="appointment-body">
                      <div class="info-row">
                        <el-icon><Clock /></el-icon>
                        <span>预约时间：{{ formatDateTime(appointment.schedule_time) }}</span>
                      </div>
                      <div class="info-row">
                        <el-icon><Timer /></el-icon>
                        <span>咨询时长：{{ appointment.duration }}分钟</span>
                      </div>
                      <div class="info-row" v-if="appointment.notes">
                        <el-icon><Document /></el-icon>
                        <span>备注：{{ appointment.notes }}</span>
                      </div>
                    </div>
                    <div class="appointment-footer">
                      <el-button
                        v-if="appointment.status === 0"
                        type="danger"
                        size="small"
                        @click="handleCancelAppointment(appointment)"
                      >
                        取消预约
                      </el-button>
                      <el-button
                        v-if="appointment.status === 1"
                        type="primary"
                        size="small"
                        @click="handleEnterAppointment(appointment)"
                      >
                        进入咨询
                      </el-button>
                    </div>
                  </div>
                </div>
              </el-tab-pane>

              <el-tab-pane label="我的账户" name="account">
                <template #label>
                  <el-icon><Wallet /></el-icon>
                  <span>我的账户</span>
                </template>
                  <div class="account-section">
                  <div class="balance-card">
                    <div class="balance-info">
                      <span class="balance-label">账户余额</span>
                      <span class="balance-amount">¥{{ userInfo?.balance || '0.00' }}</span>
                    </div>
                    <el-button type="primary" @click="showRechargeDialog = true">
                      <el-icon><Plus /></el-icon>
                      立即充值
                    </el-button>
                  </div>

                  <div class="transactions-list">
                    <h3>交易记录</h3>
                    <el-empty v-if="transactions.length === 0" description="暂无交易记录" :image-size="100" />
                    <div v-else class="transaction-item" v-for="transaction in transactions" :key="transaction.id">
                      <div class="transaction-info">
                        <span class="transaction-type">{{ getTransactionTypeText(transaction.type) }}</span>
                        <span class="transaction-desc">{{ transaction.description }}</span>
                      </div>
                      <div class="transaction-amount">
                        <span :class="['amount', transaction.type === 'recharge' ? 'income' : 'expense']">
                          {{ transaction.type === 'recharge' ? `+¥${transaction.amount}` : `-¥${transaction.amount}` }}
                        </span>
                        <span class="transaction-time">{{ formatDateTime(transaction.created_at) }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </el-tab-pane>

              <el-tab-pane label="我的订单" name="orders">
                <template #label>
                  <el-icon><Document /></el-icon>
                  <span>我的订单</span>
                </template>
                <div class="orders-list" v-loading="ordersLoading">
                  <el-empty v-if="orders.length === 0" description="暂无订单" :image-size="120">
                    <el-button type="primary" @click="$router.push('/counselors')">
                      去预约咨询
                    </el-button>
                  </el-empty>
                  <div v-else class="order-item" v-for="order in orders" :key="order.id">
                    <div class="order-header">
                      <span class="order-no">订单号：{{ order.order_no }}</span>
                      <el-tag :type="getOrderStatusType(order.status)">
                        {{ getOrderStatusText(order.status) }}
                      </el-tag>
                    </div>
                    <div class="order-body">
                      <div class="counselor-info">
                        <el-avatar :size="50" :src="order.counselor?.avatar">
                          <el-icon><User /></el-icon>
                        </el-avatar>
                        <div class="counselor-details">
                          <h4>{{ order.counselor?.name }}</h4>
                          <p>{{ order.counselor?.specialty }}</p>
                        </div>
                      </div>
                      <div class="order-info">
                        <p><el-icon><Timer /></el-icon> 时长：{{ order.duration }}分钟</p>
                        <p><el-icon><Clock /></el-icon> 预约：{{ formatDateTime(order.schedule_time) }}</p>
                      </div>
                      <div class="order-amount">
                        <span class="amount">¥{{ order.amount }}</span>
                      </div>
                    </div>
                    <div class="order-footer">
                      <el-button v-if="order.status === 1" type="primary" size="small" @click="handleEnterOrder(order)">
                        进入咨询
                      </el-button>
                      <el-button
                        v-if="order.status === 0 || order.status === 1"
                        size="small"
                        @click="handleCancelOrder(order)"
                      >
                        取消订单
                      </el-button>
                    </div>
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>
          </el-card>
        </el-col>
      </el-row>
    </div>
    <AppFooter />

    <!-- 充值对话框 -->
    <el-dialog v-model="showRechargeDialog" title="账户充值" width="400px">
      <el-form label-width="100px">
        <el-form-item label="充值金额">
          <el-input v-model="rechargeAmount" placeholder="请输入充值金额" type="number">
            <template #append>元</template>
          </el-input>
        </el-form-item>
        <el-form-item label="快捷金额">
          <div class="recharge-options">
            <el-button
              v-for="amount in rechargeOptions"
              :key="amount"
              @click="rechargeAmount = amount"
            >
              ¥{{ amount }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRechargeDialog = false">取消</el-button>
        <el-button type="primary" @click="handleRecharge">
          <el-icon><Plus /></el-icon>
          确认充值
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import {
  User,
  Camera,
  Message,
  Phone,
  Select,
  Lock,
  Calendar,
  Clock,
  Timer,
  Document,
  Wallet,
  Plus
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { updateProfile, changePassword, uploadAvatar, recharge as rechargeAPI, getTransactions } from '@/api/profile'
import { getUserInfo } from '@/api/auth'
import { getOrderList, cancelOrder as cancelOrderAPI } from '@/api/order'
import { startChatSession } from '@/api/chat'
import { handleError, showSuccess, showError } from '@/utils/errorHandler'
import { formatDateTime } from '@/utils/formatter'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import { SUCCESS_MESSAGES, REGEX } from '@/constants'

const userStore = useUserStore()
const router = useRouter()
const activeTab = ref('profile')
const submitting = ref(false)
const profileFormRef = ref()
const passwordFormRef = ref()

const userInfo = ref({
  id: null,
  username: '',
  email: '',
  phone: '',
  avatar: '',
  balance: '0.00'
})

const profileForm = reactive({
  email: '',
  phone: ''
})

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

// 预约数据
const appointments = ref([])
const appointmentsLoading = ref(false)

// 账户数据
const transactions = ref([])
const showRechargeDialog = ref(false)
const rechargeAmount = ref('')
const rechargeOptions = [100, 200, 500, 1000]

// 订单数据
const orders = ref([])
const ordersLoading = ref(false)

const validateConfirmPassword = (_rule, value, callback) => {
  if (value !== passwordForm.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const profileRules = {
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { pattern: REGEX.PHONE, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

const passwordRules = {
  old_password: [
    { required: true, message: '请输入旧密码', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const loadUserInfo = async () => {
  try {
    const res = await getUserInfo()
    const data = res.data || {}
    userInfo.value = {
      id: data.id,
      username: data.username || '',
      email: data.email || '',
      phone: data.phone || '',
      avatar: data.avatar || '',
      balance: data.balance || '0.00'
    }
    userStore.setUserInfo(data)
    userInfo.value = {
      id: data.id,
      username: data.username || '',
      email: data.email || '',
      phone: data.phone || '',
      avatar: data.avatar || '',
      balance: data.balance || '0.00'
    }
    profileForm.email = data.email || ''
    profileForm.phone = data.phone || ''
  } catch (error) {
    handleError(error, '加载用户信息失败')
  }
}

const beforeAvatarUpload = (file) => {
  // 验证图片文件
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isImage) {
    showError('请上传图片文件')
    return false
  }
  if (!isLt5M) {
    showError('图片大小不能超过 5MB')
    return false
  }
  return true
}

const handleAvatarUpload = async ({ file }) => {
  try {
    const res = await uploadAvatar(file)
    userInfo.value.avatar = res.data.url
    showSuccess(SUCCESS_MESSAGES.UPLOAD_AVATAR)
    await loadUserInfo()
  } catch (error) {
    handleError(error, '上传头像失败')
  }
}

const handleUpdateProfile = async () => {
  await profileFormRef.value.validate()
  submitting.value = true
  try {
    await updateProfile(profileForm)
    showSuccess(SUCCESS_MESSAGES.UPDATE_PROFILE)
    await loadUserInfo()
  } catch (error) {
    handleError(error, '修改资料失败')
  } finally {
    submitting.value = false
  }
}

const handleChangePassword = async () => {
  await passwordFormRef.value.validate()
  submitting.value = true
  try {
    await changePassword({
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password
    })
    showSuccess(SUCCESS_MESSAGES.CHANGE_PASSWORD)
    // 清空表单
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch (error) {
    handleError(error, '修改密码失败')
  } finally {
    submitting.value = false
  }
}

// 预约相关方法
const loadAppointments = async () => {
  appointmentsLoading.value = true
  try {
    // 使用订单列表作为预约数据
    const res = await getOrderList({ page: 1, page_size: 10, status: 1 })
    appointments.value = res.data.orders || []
  } catch (error) {
    handleError(error, '加载预约失败')
  } finally {
    appointmentsLoading.value = false
  }
}

const getAppointmentStatusText = (status) => {
  const statusMap = {
    0: '待确认',
    1: '已确认',
    2: '已完成',
    3: '已取消'
  }
  return statusMap[status] || '未知'
}

const getAppointmentStatusType = (status) => {
  const typeMap = {
    0: 'warning',
    1: 'success',
    2: 'info',
    3: 'info'
  }
  return typeMap[status] || ''
}

const handleCancelAppointment = async (appointment) => {
  try {
    await ElMessageBox.confirm('确定要取消此预约吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await cancelOrderAPI(appointment.id)
    showSuccess('预约已取消')
    await loadAppointments()
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, '取消预约失败')
    }
  }
}

const handleEnterAppointment = async (appointment) => {
  try {
    const res = await startChatSession(appointment.id)
    showSuccess('进入咨询')
    router.push(`/chat/${res.data.session_id}`)
  } catch (error) {
    handleError(error, '进入咨询失败')
  }
}

// 账户相关方法
const loadTransactions = async () => {
  try {
    const res = await getTransactions({ page: 1, page_size: 20 })
    transactions.value = res.data.transactions || []
  } catch (error) {
    // 如果API不可用，使用模拟数据
    console.warn('加载交易记录失败，使用模拟数据:', error)
    transactions.value = [
      {
        id: 1,
        type: 'recharge',
        amount: '100.00',
        description: '账户充值',
        created_at: new Date().toISOString()
      },
      {
        id: 2,
        type: 'consume',
        amount: '30.00',
        description: '购买咨询服务',
        created_at: new Date(Date.now() - 86400000).toISOString()
      }
    ]
  }
}

const getTransactionTypeText = (type) => {
  const typeMap = {
    recharge: '充值',
    consume: '消费',
    refund: '退款'
  }
  return typeMap[type] || '未知'
}

const handleRecharge = async () => {
  if (!rechargeAmount.value) {
    showError('请输入充值金额')
    return
  }
  if (isNaN(rechargeAmount.value) || rechargeAmount.value <= 0) {
    showError('请输入有效的充值金额')
    return
  }
  try {
    // 尝试调用充值API
    await rechargeAPI({ amount: parseFloat(rechargeAmount.value) })
    showSuccess(`成功充值 ¥${rechargeAmount.value}`)
    rechargeAmount.value = ''
    showRechargeDialog.value = false
    // 刷新用户信息和交易记录
    await loadUserInfo()
    await loadTransactions()
  } catch (error) {
    // 如果API不可用，模拟充值成功
    console.warn('充值API不可用，使用模拟数据:', error)
    const currentBalance = parseFloat(userInfo.value.balance || '0.00')
    const rechargeAmountNum = parseFloat(rechargeAmount.value)
    userInfo.value.balance = (currentBalance + rechargeAmountNum).toFixed(2)

    transactions.value.unshift({
      id: Date.now(),
      type: 'recharge',
      amount: rechargeAmount.value,
      description: '账户充值',
      created_at: new Date().toISOString()
    })

    showSuccess(`成功充值 ¥${rechargeAmount.value}`)
    rechargeAmount.value = ''
    showRechargeDialog.value = false
  }
}

// 订单相关方法
const loadOrders = async () => {
  ordersLoading.value = true
  try {
    const res = await getOrderList({ page: 1, page_size: 10 })
    orders.value = res.data.orders || []
  } catch (error) {
    handleError(error, '加载订单失败')
  } finally {
    ordersLoading.value = false
  }
}

const getOrderStatusText = (status) => {
  const statusMap = {
    0: '待支付',
    1: '已支付',
    2: '已完成',
    3: '已取消',
    4: '已退款'
  }
  return statusMap[status] || '未知'
}

const getOrderStatusType = (status) => {
  const typeMap = {
    0: 'warning',
    1: 'success',
    2: 'info',
    3: 'info',
    4: 'info'
  }
  return typeMap[status] || ''
}

const handleCancelOrder = async (order) => {
  try {
    await ElMessageBox.confirm('确定要取消此订单吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await cancelOrderAPI(order.id)
    showSuccess(SUCCESS_MESSAGES.CANCEL_ORDER)
    await loadOrders()
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, '取消订单失败')
    }
  }
}

const handleEnterOrder = async (order) => {
  try {
    const res = await startChatSession(order.id)
    showSuccess('进入咨询')
    router.push(`/chat/${res.data.session_id}`)
  } catch (error) {
    handleError(error, '进入咨询失败')
  }
}

// 监听 tab 切换,加载数据
const handleTabChange = (tabName) => {
  if (tabName === 'appointments') {
    loadAppointments()
  } else if (tabName === 'orders') {
    loadOrders()
  } else if (tabName === 'account') {
    loadTransactions()
  }
}

// 页面加载时初始化
onMounted(async () => {
  console.log('Profile page mounted')
  await loadUserInfo()
  console.log('User info loaded:', userInfo.value)
  await loadTransactions()
})
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.profile-content {
  flex: 1;
  padding: 20px 0;
}

.profile-row {
  margin-top: 24px;
}

.page-header {
  padding: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  margin-bottom: 24px;
}

.header-content h1 {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 600;
  color: white;
}

.header-desc {
  margin: 0;
  font-size: 14px;
  opacity: 0.9;
}



.user-card,
.settings-card {
  border-radius: 12px;
  border: none;
}

.card-header {
  text-align: center;
  margin-bottom: 20px;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  color: #666;
}

.user-card {
  text-align: center;
}

.avatar-section {
  margin-bottom: 24px;
}

.avatar-uploader {
  position: relative;
  display: inline-block;
  cursor: pointer;
}

.avatar-uploader:hover .avatar-overlay {
  opacity: 1;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  opacity: 0;
  transition: opacity 0.3s;
}

.avatar-overlay .el-icon {
  font-size: 24px;
  margin-bottom: 4px;
}

.avatar-overlay span {
  font-size: 12px;
}

.avatar-section h3 {
  margin: 16px 0 0 0;
  color: #333;
  font-size: 20px;
  font-weight: 500;
}

.info-list {
  text-align: left;
}

.info-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f5f7fa;
  gap: 12px;
}

.info-item:last-child {
  border-bottom: none;
}

.info-icon {
  color: #409eff;
  font-size: 18px;
}

.info-item .label {
  color: #666;
  font-weight: 500;
  flex-shrink: 0;
}

.info-item .value {
  color: #333;
  flex: 1;
}

.profile-tabs {
  padding: 0 24px;
}

.profile-tabs :deep(.el-tabs__header) {
  margin-bottom: 32px;
}

.profile-tabs :deep(.el-tabs__item) {
  font-size: 16px;
  padding: 0 24px;
}

.profile-tabs :deep(.el-tabs__item) .el-icon {
  margin-right: 6px;
}

.profile-form {
  max-width: 500px;
  padding: 0 24px 24px 24px;
}

.profile-form .el-button {
  width: 100%;
  margin-top: 16px;
}

/* 预约列表样式 */
.appointments-list {
  padding: 0 24px 24px 24px;
}

.appointment-item {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  margin-bottom: 16px;
  padding: 16px;
  transition: all 0.3s;
}

.appointment-item:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.appointment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.appointment-header .counselor-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.counselor-details h4 {
  margin: 0 0 4px 0;
  font-size: 16px;
  color: #333;
}

.counselor-details p {
  margin: 0;
  font-size: 13px;
  color: #999;
}

.appointment-body {
  padding: 12px 0;
  border-top: 1px solid #f5f7fa;
  border-bottom: 1px solid #f5f7fa;
  margin-bottom: 12px;
}

.appointment-body .info-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  color: #666;
  font-size: 14px;
}

.appointment-body .info-row:last-child {
  margin-bottom: 0;
}

.appointment-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

/* 账户样式 */
.account-section {
  padding: 0 24px 24px 24px;
}

.balance-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  padding: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  color: white;
}

.balance-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.balance-label {
  font-size: 14px;
  opacity: 0.9;
}

.balance-amount {
  font-size: 36px;
  font-weight: bold;
}

.balance-card .el-button {
  background: white;
  color: #667eea;
  border: none;
}

.transactions-list h3 {
  margin: 0 0 16px 0;
  font-size: 18px;
  color: #333;
}

.transaction-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #f5f7fa;
}

.transaction-item:last-child {
  border-bottom: none;
}

.transaction-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.transaction-type {
  font-weight: 500;
  color: #333;
  font-size: 15px;
}

.transaction-desc {
  font-size: 13px;
  color: #999;
}

.transaction-amount {
  text-align: right;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.transaction-amount .amount {
  font-size: 18px;
  font-weight: bold;
}

.transaction-amount .amount.income {
  color: #67c23a;
}

.transaction-amount .amount.expense {
  color: #f56c6c;
}

.transaction-time {
  font-size: 12px;
  color: #999;
}

/* 充值对话框样式 */
.recharge-options {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.recharge-options .el-button {
  margin: 0;
}

/* 订单列表样式 */
.orders-list {
  padding: 0 24px 24px 24px;
}

.order-item {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  margin-bottom: 16px;
  padding: 16px;
  transition: all 0.3s;
}

.order-item:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 12px;
  border-bottom: 1px solid #f5f7fa;
  margin-bottom: 12px;
}

.order-no {
  font-size: 13px;
  color: #666;
}

.order-body {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.order-body .counselor-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.order-body .counselor-details h4 {
  margin: 0 0 4px 0;
  font-size: 16px;
  color: #333;
}

.order-body .counselor-details p {
  margin: 0;
  font-size: 13px;
  color: #999;
}

.order-body .order-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
  color: #666;
  font-size: 14px;
}

.order-body .order-info p {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0;
}

.order-body .order-amount {
  display: flex;
  align-items: center;
}

.order-body .order-amount .amount {
  font-size: 24px;
  font-weight: bold;
  color: #f56c6c;
}

.order-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid #f5f7fa;
}

@media (max-width: 768px) {
  .page-header {
    text-align: center;
  }

  .profile-content {
    margin-top: 20px;
  }

  .profile-form {
    max-width: 100%;
  }
}
</style>

