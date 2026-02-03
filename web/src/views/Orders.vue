<template>
  <div class="orders-page">
    <AppHeader />
    <div class="orders-content">
      <div class="page-header">
        <div class="header-content">
          <h1>我的订单</h1>
          <p class="header-desc">查看和管理您的咨询订单</p>
        </div>
        <div class="header-actions">
          <el-button type="primary" @click="$router.push('/counselors')">
            <el-icon><Plus /></el-icon>
            预约咨询
          </el-button>
        </div>
      </div>

      <el-card class="order-card">
        <el-tabs v-model="activeTab" @tab-change="handleTabChange" class="order-tabs">
          <el-tab-pane label="全部订单" name="">
          <template #label>
            <el-icon><List /></el-icon>
            <span>全部订单</span>
          </template>
        </el-tab-pane>
        <el-tab-pane label="待支付" name="0">
          <template #label>
            <el-icon><Wallet /></el-icon>
            <span>待支付</span>
          </template>
        </el-tab-pane>
        <el-tab-pane label="已支付" name="1">
          <template #label>
            <el-icon><Select /></el-icon>
            <span>已支付</span>
          </template>
        </el-tab-pane>
        <el-tab-pane label="已完成" name="2">
          <template #label>
            <el-icon><CircleCheck /></el-icon>
            <span>已完成</span>
          </template>
        </el-tab-pane>
        <el-tab-pane label="已取消" name="3">
          <template #label>
            <el-icon><CircleClose /></el-icon>
            <span>已取消</span>
          </template>
        </el-tab-pane>
        <el-tab-pane label="已退款" name="4">
          <template #label>
            <el-icon><RefreshLeft /></el-icon>
            <span>已退款</span>
          </template>
        </el-tab-pane>
        </el-tabs>

        <div v-loading="loading">
          <div v-if="orders.length > 0">
            <div class="order-item" v-for="order in orders" :key="order.id">
              <div class="order-header">
                <div class="order-info">
                  <span class="order-no">订单号：{{ order.order_no }}</span>
                  <span class="order-time">{{ formatDateTime(order.created_at) }}</span>
                </div>
                <el-tag :type="getStatusType(order.status)">{{ getStatusText(order.status) }}</el-tag>
              </div>
              <div class="order-body">
                <el-avatar :size="60" :src="order.counselor?.avatar">
                  <el-icon><User /></el-icon>
                </el-avatar>
                <div class="counselor-info">
                  <h3>{{ order.counselor?.name }}</h3>
                  <p class="specialty">{{ order.counselor?.specialty }}</p>
                </div>
                <div class="order-details">
                  <p>咨询时长：{{ order.duration }}分钟</p>
                  <p>预约时间：{{ formatDateTime(order.schedule_time) }}</p>
                  <p v-if="order.notes">备注：{{ order.notes }}</p>
                </div>
                <div class="order-amount">
                  <p class="amount">¥{{ order.amount }}</p>
                </div>
              </div>
              <div class="order-footer">
                <el-button
                  v-if="order.status === 1"
                  type="primary"
                  @click="handleEnterChat(order)"
                >
                  进入咨询
                </el-button>
                <el-button
                  v-if="order.status === 0 || order.status === 1"
                  @click="handleCancelOrder(order)"
                >
                  取消订单
                </el-button>
                <el-button
                  v-if="order.status === 0"
                  type="success"
                  @click="showPayDialog(order)"
                >
                  立即支付
                </el-button>
                <el-button
                  v-if="order.status === 2 && !order.has_review"
                  type="primary"
                  @click="showReviewDialog(order)"
                >
                  评价咨询师
                </el-button>
              </div>
            </div>

            <div class="pagination" v-if="total > 0">
              <el-pagination
                v-model:current-page="pagination.page"
                v-model:page-size="pagination.page_size"
                :page-sizes="PAGINATION.PAGE_SIZE_OPTIONS"
                :total="total"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="loadOrders"
                @current-change="loadOrders"
              />
            </div>
          </div>
          <el-empty v-else description="暂无订单" :image-size="200">
            <el-button type="primary" @click="$router.push('/counselors')">
              去预约咨询
            </el-button>
          </el-empty>
        </div>
      </el-card>
    </div>
    <AppFooter />

    <!-- 支付对话框 -->
    <el-dialog v-model="payDialogVisible" title="选择支付方式" width="500px" :close-on-click-modal="false">
      <div v-if="currentOrder" class="pay-dialog-content">
        <div class="order-summary">
          <h4>订单信息</h4>
          <div class="order-info-item">
            <span class="label">订单号：</span>
            <span class="value">{{ currentOrder.order_no }}</span>
          </div>
          <div class="order-info-item">
            <span class="label">咨询师：</span>
            <span class="value">{{ currentOrder.counselor?.name }}</span>
          </div>
          <div class="order-info-item">
            <span class="label">咨询时长：</span>
            <span class="value">{{ currentOrder.duration }}分钟</span>
          </div>
          <div class="order-info-item">
            <span class="label">预约时间：</span>
            <span class="value">{{ formatDateTime(currentOrder.schedule_time) }}</span>
          </div>
          <div class="order-amount-display">
            <span class="label">支付金额：</span>
            <span class="amount">¥{{ currentOrder.amount }}</span>
          </div>
        </div>

        <el-divider />

        <div class="payment-methods">
          <h4>选择支付方式</h4>
          <el-radio-group v-model="selectedPaymentMethod" class="payment-group">
            <el-radio-button value="wechat">
              <div class="payment-option">
                <div class="payment-icon wechat">
                  <svg viewBox="0 0 24 24" width="24" height="24">
                    <path fill="#07C160" d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 0 1 .598.082l1.584.926a.272.272 0 0 0 .14.045c.133 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 0 1-.023-.156.49.49 0 0 1 .201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-6.656-6.088V8.89c-.135-.01-.269-.027-.407-.03zm-2.53 3.274c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982zm4.844 0c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982z"/>
                  </svg>
                </div>
                <span class="payment-name">微信支付</span>
              </div>
            </el-radio-button>
            <el-radio-button value="alipay">
              <div class="payment-option">
                <div class="payment-icon alipay">
                  <svg viewBox="0 0 24 24" width="24" height="24">
                    <path fill="#1677FF" d="M16.35 8.98c-1.36 0-2.48-1.11-2.48-2.48 0-1.36 1.12-2.48 2.48-2.48 1.36 0 2.48 1.12 2.48 2.48 0 1.37-1.12 2.48-2.48 2.48zm-8.7 0c-1.36 0-2.48-1.11-2.48-2.48 0-1.36 1.12-2.48 2.48-2.48 1.36 0 2.48 1.12 2.48 2.48 0 1.37-1.12 2.48-2.48 2.48zM12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm0-14c-3.31 0-6 2.69-6 6s2.69 6 6 6 6-2.69 6-6-2.69-6-6-6z"/>
                  </svg>
                </div>
                <span class="payment-name">支付宝</span>
              </div>
            </el-radio-button>
          </el-radio-group>
        </div>

        <div class="pay-tips">
          <el-icon><InfoFilled /></el-icon>
          <span>支付成功后将自动进入咨询环节，请确保网络连接正常</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="payDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="paying" @click="handlePay">
          立即支付 ¥{{ currentOrder?.amount }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 评价对话框 -->
    <el-dialog v-model="reviewDialogVisible" title="评价咨询师" width="500px">
      <div v-if="currentOrder" class="review-dialog-content">
        <div class="counselor-info">
          <el-avatar :size="60" :src="currentOrder.counselor?.avatar">
            <el-icon><User /></el-icon>
          </el-avatar>
          <div class="counselor-name">{{ currentOrder.counselor?.name }}</div>
        </div>

        <el-form :model="reviewForm" label-width="80px">
          <el-form-item label="总体评价">
            <el-rate v-model="reviewForm.rating" show-text />
          </el-form-item>
          <el-form-item label="评价内容">
            <el-input
              v-model="reviewForm.comment"
              type="textarea"
              :rows="4"
              placeholder="请输入您对咨询师的评价..."
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
          <el-form-item label="匿名评价">
            <el-switch v-model="reviewForm.is_anonymous" />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="reviewDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submittingReview" @click="handleSubmitReview">
          提交评价
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  List,
  Wallet,
  Select,
  CircleCheck,
  CircleClose,
  RefreshLeft,
  User,
  InfoFilled
} from '@element-plus/icons-vue'
import { getOrderList, cancelOrder as cancelOrderAPI } from '@/api/order'
import { getOrderSessionId } from '@/api/chat'
import { createPayment } from '@/api/payment'
import { createCounselorReview } from '@/api/counselor'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import { formatDateTime } from '@/utils/formatter'
import { handleError, showSuccess } from '@/utils/errorHandler'
import { ORDER_STATUS_TEXT, ORDER_STATUS_COLOR, ORDER_STATUS, SUCCESS_MESSAGES, PAGINATION } from '@/constants'

const router = useRouter()
const loading = ref(false)
const orders = ref([])
const total = ref(0)
const activeTab = ref('')
const payDialogVisible = ref(false)
const reviewDialogVisible = ref(false)
const currentOrder = ref(null)
const paying = ref(false)
const submittingReview = ref(false)
const selectedPaymentMethod = ref('wechat')

const pagination = reactive({
  page: PAGINATION.DEFAULT_PAGE,
  page_size: PAGINATION.DEFAULT_PAGE_SIZE
})

onMounted(() => {
  loadOrders()
})

const loadOrders = async () => {
  loading.value = true
  try {
    const params = { ...pagination }
    if (activeTab.value !== '') {
      params.status = activeTab.value
    }
    const res = await getOrderList(params)
    orders.value = res.data.orders
    total.value = res.data.total
  } catch (error) {
    handleError(error, '加载订单失败')
  } finally {
    loading.value = false
  }
}

const handleTabChange = () => {
  pagination.page = 1
  loadOrders()
}



// 订单状态映射 - 用于兼容后端返回的数字状态码
const ORDER_STATUS_MAP = {
  0: ORDER_STATUS.PENDING,
  1: ORDER_STATUS.PAID,
  2: ORDER_STATUS.COMPLETED,
  3: ORDER_STATUS.CANCELLED,
  4: ORDER_STATUS.REFUNDED
}

const getStatusText = (status) => {
  const statusKey = ORDER_STATUS_MAP[status]
  return ORDER_STATUS_TEXT[statusKey] || '未知'
}

const getStatusType = (status) => {
  const statusKey = ORDER_STATUS_MAP[status]
  return ORDER_STATUS_COLOR[statusKey] || ''
}

const handleCancelOrder = async (order) => {
  try {
    await ElMessageBox.confirm('确定要取消此订单吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await cancelOrderAPI(order.id)
    ElMessage.success(SUCCESS_MESSAGES.CANCEL_ORDER)
    loadOrders()
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, '取消订单失败')
    }
  }
}

const handlePay = async () => {
  if (!currentOrder.value) return

  paying.value = true
  try {
    const res = await createPayment({
      order_id: currentOrder.value.id,
      payment_method: selectedPaymentMethod.value,
      trade_type: 'NATIVE',
      client_ip: '127.0.0.1',
      return_url: window.location.origin + '/orders'
    })

    if (res.data.pay_url) {
      // 如果返回支付URL，在新窗口打开
      window.open(res.data.pay_url, '_blank')
    }

    ElMessage.success('支付订单已创建')
    payDialogVisible.value = false
    loadOrders()
  } catch (error) {
    handleError(error, '创建支付订单失败')
  } finally {
    paying.value = false
  }
}

const showPayDialog = (order) => {
  currentOrder.value = order
  payDialogVisible.value = true
}

const showReviewDialog = (order) => {
  currentOrder.value = order
  reviewForm.value = {
    counselor_id: order.counselor_id,
    order_id: order.id,
    rating: 5,
    comment: '',
    is_anonymous: false
  }
  reviewDialogVisible.value = true
}

const handleSubmitReview = async () => {
  if (!reviewForm.value.comment.trim()) {
    ElMessage.warning('请输入评价内容')
    return
  }

  submittingReview.value = true
  try {
    await createCounselorReview(reviewForm.value)
    showSuccess(SUCCESS_MESSAGES.SUBMIT_REVIEW)
    reviewDialogVisible.value = false
    loadOrders()
  } catch (error) {
    handleError(error, '提交评价失败')
  } finally {
    submittingReview.value = false
  }
}

const handleEnterChat = async (order) => {
  try {
    const res = await getOrderSessionId(order.id)
    ElMessage.success('进入咨询')
    router.push(`/chat/${res.data.session_id}`)
  } catch (error) {
    handleError(error, '进入咨询失败')
  }
}
</script>

<style scoped>
.orders-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.orders-content {
  flex: 1;
  padding: 20px 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
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

.header-actions .el-button {
  background: white;
  color: #667eea;
  border: none;
}

.header-actions .el-button:hover {
  background: #f5f7fa;
}

.order-card {
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.order-tabs {
  border-bottom: 1px solid #ebeef5;
}

.order-tabs :deep(.el-tabs__header) {
  margin: 0 0 24px 0;
}

.order-tabs :deep(.el-tabs__item) {
  font-size: 15px;
}

.order-tabs :deep(.el-tabs__item) .el-icon {
  margin-right: 4px;
}

.order-item {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  margin-bottom: 16px;
  transition: all 0.3s;
}

.order-item:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
}

.order-info {
  display: flex;
  gap: 20px;
  color: #666;
  font-size: 14px;
}

.order-no {
  font-weight: 500;
}

.order-body {
  display: flex;
  align-items: center;
  padding: 20px 16px;
  gap: 20px;
}

.counselor-info h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #333;
}

.counselor-info .specialty {
  margin: 0;
  color: #999;
  font-size: 13px;
}

.order-details {
  flex: 1;
  color: #666;
  font-size: 14px;
}

.order-details p {
  margin: 4px 0;
}

.order-amount {
  text-align: right;
  min-width: 100px;
}

.order-amount .amount {
  font-size: 24px;
  font-weight: bold;
  color: #f56c6c;
  margin: 0;
}

.order-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 12px 16px;
  border-top: 1px solid #ebeef5;
  background: #fafafa;
}

.pagination {
  margin-top: 32px;
  padding: 20px 0;
  border-top: 1px solid #ebeef5;
  text-align: center;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }

  .header-content h1 {
    font-size: 24px;
  }

  .order-body {
    flex-direction: column;
    align-items: flex-start;
  }

  .order-amount {
    text-align: left;
  }

  .order-footer {
    flex-wrap: wrap;
  }
}

/* 支付对话框样式 */
.pay-dialog-content {
  padding: 10px 0;
}

.order-summary h4 {
  margin: 0 0 16px 0;
  font-size: 16px;
  color: #333;
}

.order-info-item {
  display: flex;
  margin-bottom: 12px;
  font-size: 14px;
}

.order-info-item .label {
  width: 80px;
  color: #666;
}

.order-info-item .value {
  flex: 1;
  color: #333;
  font-weight: 500;
}

.order-amount-display {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 16px;
  padding: 12px 16px;
  background: #f5f7fa;
  border-radius: 4px;
}

.order-amount-display .label {
  font-size: 14px;
  color: #666;
}

.order-amount-display .amount {
  font-size: 24px;
  font-weight: bold;
  color: #f56c6c;
}

.payment-methods h4 {
  margin: 20px 0 16px 0;
  font-size: 16px;
  color: #333;
}

.payment-group {
  width: 100%;
  display: flex;
  gap: 16px;
}

.payment-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
}

.payment-icon {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.payment-icon.wechat {
  color: #07C160;
}

.payment-icon.alipay {
  color: #1677FF;
}

.payment-name {
  font-size: 14px;
  font-weight: 500;
}

.pay-tips {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 20px;
  padding: 12px 16px;
  background: #f0f9ff;
  border-left: 3px solid #409eff;
  border-radius: 4px;
  font-size: 13px;
  color: #606266;
}

/* 评价对话框样式 */
.review-dialog-content {
  padding: 10px 0;
}

.counselor-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 24px;
}

.counselor-info .counselor-name {
  margin-top: 12px;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}
</style>
