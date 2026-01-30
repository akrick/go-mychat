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
                  @click="handlePay(order)"
                >
                  立即支付
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
  User
} from '@element-plus/icons-vue'
import { getOrderList, cancelOrder as cancelOrderAPI } from '@/api/order'
import { startChatSession } from '@/api/chat'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import { formatDateTime } from '@/utils/formatter'
import { handleError } from '@/utils/errorHandler'
import { ORDER_STATUS_TEXT, ORDER_STATUS_COLOR, ORDER_STATUS, SUCCESS_MESSAGES, PAGINATION } from '@/constants'

const router = useRouter()
const loading = ref(false)
const orders = ref([])
const total = ref(0)
const activeTab = ref('')

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

const handlePay = (_order) => {
  ElMessage.info('支付功能开发中，请联系管理员')
}

const handleEnterChat = async (order) => {
  try {
    const res = await startChatSession(order.id)
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
</style>
