<template>
  <div class="counselor-detail-page">
    <AppHeader />
    <div class="counselor-detail" v-loading="loading">
      <el-page-header @back="$router.back()" content="咨询师详情" class="page-header" />

      <el-row :gutter="30" v-if="counselor">
        <el-col :span="8">
          <el-card class="profile-card">
            <!-- 推荐标签 -->
            <div v-if="counselor.is_recommended" class="recommend-badge">
              <el-icon><Star /></el-icon>
              <span>优质推荐</span>
            </div>

            <div class="avatar">
              <el-avatar :size="120" :src="counselor.avatar">
                <el-icon><User /></el-icon>
              </el-avatar>
            </div>
            <h2>{{ counselor.name }}</h2>
            <p class="title">{{ counselor.title }}</p>
            <div class="rating">
              <el-rate v-model="counselor.rating" disabled show-score text-color="#ff9900" />
            </div>

            <!-- 统计信息 -->
            <div class="stats-section">
              <div class="stat-item">
                <div class="stat-value">{{ counselor.service_count || 0 }}</div>
                <div class="stat-label">服务人数</div>
              </div>
              <div class="stat-divider"></div>
              <div class="stat-item">
                <div class="stat-value">{{ counselor.review_count || 0 }}</div>
                <div class="stat-label">评价数</div>
              </div>
            </div>

            <div class="info">
              <div class="info-item">
                <el-icon><Calendar /></el-icon>
                <span>{{ counselor.years_exp }}年经验</span>
              </div>
              <div class="info-item">
                <el-icon><Trophy /></el-icon>
                <span>擅长：{{ counselor.specialty }}</span>
              </div>
            </div>
            <div class="price-section">
              <p class="price">¥{{ counselor.price }}</p>
              <p class="unit">/分钟</p>
            </div>
            <el-button type="primary" size="large" @click="showOrderDialog = true" style="width: 100%">
              预约咨询
            </el-button>
          </el-card>
        </el-col>

        <el-col :span="16">
          <el-card class="intro-card">
            <template #header>
              <h3>个人简介</h3>
            </template>
            <p class="bio">{{ counselor.bio || '暂无简介' }}</p>
          </el-card>

          <!-- 用户评价 -->
          <el-card class="reviews-card" style="margin-top: 20px">
            <template #header>
              <div class="card-header">
                <h3>用户评价</h3>
                <span class="review-count">({{ reviews.length }}条)</span>
              </div>
            </template>
            <div v-if="reviews.length > 0" class="reviews-list">
              <div v-for="review in reviews" :key="review.id" class="review-item">
                <div class="review-header">
                  <el-avatar :size="32" :src="review.user?.avatar">
                    <el-icon><User /></el-icon>
                  </el-avatar>
                  <div class="review-user">
                    <span class="username">{{ review.is_anonymous ? '匿名用户' : review.user?.username || '用户' }}</span>
                    <el-rate v-model="review.rating" disabled size="small" />
                  </div>
                  <span class="review-time">{{ formatTime(review.created_at) }}</span>
                </div>
                <p class="review-content">{{ review.comment }}</p>
              </div>
            </div>
            <el-empty v-else description="暂无评价" />
          </el-card>
        </el-col>
      </el-row>

    <!-- 预约对话框 -->
    <el-dialog v-model="showOrderDialog" title="预约咨询" width="500px">
      <el-form :model="orderForm" :rules="orderRules" ref="orderFormRef" label-width="100px">
        <el-form-item label="咨询师">
          <el-input v-model="counselor?.name" disabled />
        </el-form-item>
        <el-form-item label="咨询时长" prop="duration">
          <el-select v-model="orderForm.duration" placeholder="请选择时长" style="width: 100%">
            <el-option
              v-for="duration in Object.values(SESSION_DURATION)"
              :key="duration.value"
              :label="duration.label"
              :value="duration.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="预约时间" prop="schedule_time">
          <el-date-picker
            v-model="orderForm.schedule_time"
            type="datetime"
            placeholder="选择日期时间"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DDTHH:mm:ss"
            style="width: 100%"
            :disabled-date="disabledDate"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="orderForm.notes"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息（可选）"
          />
        </el-form-item>
        <el-form-item label="订单金额">
          <span class="total-price">¥{{ totalAmount }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showOrderDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmitOrder">提交订单</el-button>
      </template>
    </el-dialog>
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { User, Calendar, Trophy, Star } from '@element-plus/icons-vue'
import { getCounselorDetail, getCounselorReviews } from '@/api/counselor'
import { createOrder } from '@/api/order'
import { handleError, showSuccess } from '@/utils/errorHandler'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import { SESSION_DURATION, SUCCESS_MESSAGES } from '@/constants'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const counselor = ref(null)
const reviews = ref([])
const showOrderDialog = ref(false)
const submitting = ref(false)
const orderFormRef = ref()

const orderForm = ref({
  counselor_id: null,
  duration: SESSION_DURATION.MINUTES_60.value,
  schedule_time: '',
  notes: ''
})

const totalAmount = computed(() => {
  if (counselor.value && orderForm.value.duration) {
    return (counselor.value.price * orderForm.value.duration).toFixed(2)
  }
  return '0.00'
})

const orderRules = {
  duration: [{ required: true, message: '请选择咨询时长', trigger: 'change' }],
  schedule_time: [{ required: true, message: '请选择预约时间', trigger: 'change' }]
}

const disabledDate = (time) => {
  return time.getTime() < Date.now() - 24 * 60 * 60 * 1000
}

const formatTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) {
    const hours = Math.floor(diff / (1000 * 60 * 60))
    if (hours === 0) {
      const minutes = Math.floor(diff / (1000 * 60))
      if (minutes === 0) return '刚刚'
      return `${minutes}分钟前`
    }
    return `${hours}小时前`
  } else if (days < 7) {
    return `${days}天前`
  } else {
    return date.toLocaleDateString('zh-CN')
  }
}

onMounted(async () => {
  await loadCounselorDetail()
  await loadCounselorReviews()
})

const loadCounselorDetail = async () => {
  loading.value = true
  try {
    const res = await getCounselorDetail(route.params.id)
    counselor.value = res.data
    orderForm.value.counselor_id = res.data.id
  } catch (error) {
    handleError(error, '加载咨询师详情失败')
  } finally {
    loading.value = false
  }
}

const loadCounselorReviews = async () => {
  try {
    const res = await getCounselorReviews(route.params.id)
    reviews.value = res.data.reviews || []
  } catch (error) {
    // 评价加载失败不影响主流程
    console.error('加载评价失败:', error)
  }
}

const handleSubmitOrder = async () => {
  await orderFormRef.value.validate()
  submitting.value = true
  try {
    await createOrder(orderForm.value)
    showSuccess(SUCCESS_MESSAGES.CREATE_ORDER)
    showOrderDialog.value = false
    router.push('/orders')
  } catch (error) {
    handleError(error, '创建订单失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.counselor-detail-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.counselor-detail {
  flex: 1;
  padding: 20px 0;
}

.page-header {
  margin-bottom: 20px;
}

.profile-card {
  text-align: center;
}

.avatar {
  margin-bottom: 20px;
}

.profile-card h2 {
  margin: 0 0 8px 0;
  color: #333;
}

.profile-card .title {
  margin: 0 0 16px 0;
  color: #666;
  font-size: 14px;
}

.profile-card .rating {
  margin-bottom: 20px;
}

.profile-card .info {
  margin: 20px 0;
  text-align: left;
}

.profile-card .info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  color: #666;
}

.price-section {
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 4px;
  margin: 20px 0;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 4px;
}

.price-section .price {
  font-size: 32px;
  font-weight: bold;
  color: #f56c6c;
  margin: 0;
}

.price-section .unit {
  color: #999;
  font-size: 14px;
}

.intro-card h3 {
  margin: 0;
  color: #333;
}

.intro-card .bio {
  line-height: 1.8;
  color: #666;
  white-space: pre-wrap;
}

.total-price {
  font-size: 24px;
  font-weight: bold;
  color: #f56c6c;
}

.recommend-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(240, 147, 251, 0.4);
}

.recommend-badge .el-icon {
  font-size: 14px;
}

.stats-section {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin: 20px 0;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 12px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #409eff;
  line-height: 1.2;
}

.stat-label {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.stat-divider {
  width: 1px;
  height: 32px;
  background: #dcdfe6;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-header h3 {
  margin: 0;
  color: #333;
}

.review-count {
  color: #909399;
  font-size: 14px;
}

.reviews-list {
  max-height: 500px;
  overflow-y: auto;
}

.review-item {
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.review-item:last-child {
  border-bottom: none;
}

.review-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.review-user {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.review-user .username {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.review-time {
  font-size: 12px;
  color: #909399;
}

.review-content {
  margin: 0;
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}
</style>
