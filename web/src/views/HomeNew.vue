<template>
  <div class="home">
    <AppHeader />

    <!-- 轮播图 -->
    <div class="carousel-section">
      <el-carousel :interval="5000" height="500px" indicator-position="outside">
        <el-carousel-item v-for="(item, index) in carouselItems" :key="index">
          <div class="carousel-item" :style="{ backgroundImage: `url(${item.image})` }">
            <div class="carousel-content">
              <h1>{{ item.title }}</h1>
              <p>{{ item.subtitle }}</p>
              <el-button type="primary" size="large" @click="handleAction(item.action)">
                {{ item.btnText }}
              </el-button>
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>
    </div>

    <!-- 服务优势 -->
    <div class="features-section">
      <div class="container">
        <div class="section-header">
          <h2>为什么选择我们</h2>
          <p>专业团队 · 隐私保护 · 灵活便捷</p>
        </div>
        <el-row :gutter="30">
          <el-col :span="6" v-for="feature in features" :key="feature.title">
            <div class="feature-card">
              <div class="feature-icon" :style="{ background: feature.bgColor }">
                <el-icon :size="40" :color="feature.iconColor">
                  <component :is="feature.icon" />
                </el-icon>
              </div>
              <h3>{{ feature.title }}</h3>
              <p>{{ feature.desc }}</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- 服务产品 -->
    <div class="services-section">
      <div class="container">
        <div class="section-header">
          <h2>服务产品</h2>
          <p>多样化咨询产品，满足不同需求</p>
        </div>
        <el-row :gutter="30">
          <el-col :span="6" v-for="service in services" :key="service.title">
            <div class="service-card">
              <div class="service-image">
                <img :src="service.image" :alt="service.title" />
              </div>
              <div class="service-content">
                <h3>{{ service.title }}</h3>
                <p class="service-desc">{{ service.desc }}</p>
                <div class="service-price">
                  <span class="price">¥{{ service.price }}</span>
                  <span class="unit">/{{ service.unit }}</span>
                </div>
                <el-button type="primary" class="btn-booking" @click="handleBooking(service)">
                  立即预约
                </el-button>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- 咨询师入驻入口 -->
    <div class="counselor-join-section">
      <div class="container">
        <el-card class="join-card">
          <el-row :gutter="40" align="middle">
            <el-col :span="12">
              <div class="join-content">
                <h2>加入我们，成为专业咨询师</h2>
                <p>如果您持有心理咨询师相关资质证书，欢迎加入我们的平台，为更多需要帮助的人提供专业的心理咨询服务</p>
                <div class="join-benefits">
                  <div class="benefit-item">
                    <el-icon><Check /></el-icon>
                    <span>专业平台，稳定客源</span>
                  </div>
                  <div class="benefit-item">
                    <el-icon><Check /></el-icon>
                    <span>灵活工作时间</span>
                  </div>
                  <div class="benefit-item">
                    <el-icon><Check /></el-icon>
                    <span>高额收入分成</span>
                  </div>
                  <div class="benefit-item">
                    <el-icon><Check /></el-icon>
                    <span>专业培训支持</span>
                  </div>
                </div>
              </div>
            </el-col>
            <el-col :span="12" class="join-action">
              <h3>立即申请入驻</h3>
              <p>填写基本信息，上传资质证书，审核通过后即可开始接单</p>
              <el-button type="primary" size="large" @click="handleJoinCounselor">
                申请入驻
              </el-button>
            </el-col>
          </el-row>
        </el-card>
      </div>
    </div>

    <!-- 咨询流程 -->
    <div class="process-section">
      <div class="container">
        <div class="section-header">
          <h2>咨询流程</h2>
          <p>简单四步，开启您的心理咨询之旅</p>
        </div>
        <div class="process-steps">
          <div class="step-item" v-for="(step, index) in steps" :key="index">
            <div class="step-number">{{ index + 1 }}</div>
            <div class="step-icon">
              <el-icon :size="48" :color="step.color">
                <component :is="step.icon" />
              </el-icon>
            </div>
            <h4>{{ step.title }}</h4>
            <p>{{ step.desc }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 热门咨询师 -->
    <div class="counselors-section">
      <div class="container">
        <div class="section-header">
          <h2>热门咨询师</h2>
          <el-button text @click="$router.push('/counselors')">查看全部 →</el-button>
        </div>
        <el-row :gutter="30" v-loading="loading">
          <el-col :span="6" v-for="counselor in counselors" :key="counselor.id">
            <div class="counselor-card" @click="handleCounselorClick(counselor.id)">
              <div class="counselor-avatar">
                <el-avatar :size="100" :src="counselor.avatar">
                  <el-icon><User /></el-icon>
                </el-avatar>
              </div>
              <h3>{{ counselor.name }}</h3>
              <p class="title">{{ counselor.title }}</p>
              <el-tag size="small" type="info">{{ counselor.specialty }}</el-tag>
              <div class="rating">
                <el-rate v-model="counselor.rating" disabled show-score text-color="#ff9900" />
              </div>
              <p class="price">¥{{ counselor.price }}/分钟</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- 数据统计 -->
    <div class="stats-section">
      <div class="container">
        <el-row :gutter="40">
          <el-col :span="6" v-for="stat in stats" :key="stat.label">
            <div class="stat-item">
              <div class="stat-number">{{ stat.value }}</div>
              <div class="stat-label">{{ stat.label }}</div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <AppFooter />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, UserFilled, Clock, Lock, Star, Search, Calendar, ChatDotRound, CircleCheck, Check } from '@element-plus/icons-vue'
import { getCounselorList } from '@/api/counselor'
import { useUserStore } from '@/stores/user'
import { handleError, showWarning } from '@/utils/errorHandler'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const counselors = ref([])

// 轮播图数据
const carouselItems = ref([
  {
    image: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=1920&h=500&fit=crop',
    title: '专业的心理咨询',
    subtitle: '连接您与专业心理咨询师，开启心灵成长之旅',
    btnText: '寻找咨询师',
    action: 'counselors'
  },
  {
    image: 'https://images.unsplash.com/photo-1526506118085-60ce8714f8c5?w=1920&h=500&fit=crop',
    title: '用心倾听每一刻',
    subtitle: '我们在这里，倾听您的心声，陪伴您度过难关',
    btnText: '了解更多',
    action: 'about'
  },
  {
    image: 'https://images.unsplash.com/photo-1493836512294-502baa1986e2?w=1920&h=500&fit=crop',
    title: '隐私安全可靠',
    subtitle: '严格保护您的隐私，安全可靠的心理咨询服务',
    btnText: '立即预约',
    action: 'booking'
  }
])

// 服务优势
const features = ref([
  {
    icon: UserFilled,
    title: '专业团队',
    desc: '持证心理咨询师，丰富的临床经验',
    bgColor: '#e8f4fd',
    iconColor: '#409eff'
  },
  {
    icon: Clock,
    title: '灵活预约',
    desc: '随时随地在线预约，便捷高效',
    bgColor: '#e8f8f3',
    iconColor: '#67c23a'
  },
  {
    icon: Lock,
    title: '隐私保护',
    desc: '严格保密机制，保护您的隐私',
    bgColor: '#fef9e8',
    iconColor: '#e6a23c'
  },
  {
    icon: Star,
    title: '品质保障',
    desc: '优质服务体验，满意度高',
    bgColor: '#fee8f0',
    iconColor: '#f56c6c'
  }
])

// 服务产品
const services = ref([
  {
    title: '个人心理咨询',
    desc: '一对一专业心理咨询，解决个人情绪困扰',
    price: '300',
    unit: '小时',
    image: 'https://images.unsplash.com/photo-1544717305-2782549b5136?w=400&h=300&fit=crop'
  },
  {
    title: '情感咨询',
    desc: '专业情感指导，改善亲密关系',
    price: '400',
    unit: '小时',
    image: 'https://images.unsplash.com/photo-1516589178581-6cd7833ae3b2?w=400&h=300&fit=crop'
  },
  {
    title: '职业规划',
    desc: '职业发展方向指导，实现职业目标',
    price: '500',
    unit: '小时',
    image: 'https://images.unsplash.com/photo-1486312338219-ce68d2c6f44d?w=400&h=300&fit=crop'
  },
  {
    title: '家庭治疗',
    desc: '家庭关系改善，促进家庭和谐',
    price: '600',
    unit: '小时',
    image: 'https://images.unsplash.com/photo-1529156069898-49953e39b3ac?w=400&h=300&fit=crop'
  }
])

// 咨询流程
const steps = ref([
  {
    icon: Search,
    title: '选择咨询师',
    desc: '根据需求选择合适的咨询师',
    color: '#409eff'
  },
  {
    icon: Calendar,
    title: '在线预约',
    desc: '选择时间，提交预约申请',
    color: '#67c23a'
  },
  {
    icon: ChatDotRound,
    title: '在线咨询',
    desc: '通过文字或语音进行咨询',
    color: '#e6a23c'
  },
  {
    icon: CircleCheck,
    title: '完成咨询',
    desc: '获得专业建议和心理支持',
    color: '#f56c6c'
  }
])

// 数据统计
const stats = ref([
  { value: '500+', label: '专业咨询师' },
  { value: '10万+', label: '服务用户' },
  { value: '98%', label: '满意度' },
  { value: '24/7', label: '在线服务' }
])

onMounted(async () => {
  await loadCounselors()
})

const loadCounselors = async () => {
  loading.value = true
  try {
    const res = await getCounselorList({ page: 1, page_size: 4 })
    counselors.value = res.data.counselors || []
  } catch (error) {
    handleError(error, '加载咨询师失败')
    counselors.value = []
  } finally {
    loading.value = false
  }
}

const handleAction = (action) => {
  if (action === 'counselors') {
    router.push('/counselors')
  } else if (action === 'about') {
    router.push('/about')
  } else if (action === 'booking') {
    if (!userStore.token) {
      showWarning('请先登录')
      router.push('/login')
    } else {
      router.push('/counselors')
    }
  }
}

const handleBooking = (service) => {
  if (!userStore.token) {
    showWarning('请先登录')
    router.push('/login')
  } else {
    router.push({ path: '/counselors', query: { serviceType: service.title } })
  }
}

const handleCounselorClick = (counselorId) => {
  if (!userStore.token) {
    showWarning('请先登录')
    router.push('/login')
  } else {
    router.push(`/counselor/${counselorId}`)
  }
}

const handleJoinCounselor = () => {
  if (!userStore.token) {
    showWarning('请先登录')
    router.push('/login')
  } else {
    router.push('/counselor-application')
  }
}
</script>

<style scoped>
.home {
  min-height: 100vh;
  background: #f5f7fa;
}

.carousel-section {
  margin-bottom: 60px;
}

.carousel-item {
  width: 100%;
  height: 500px;
  background-size: cover;
  background-position: center;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.carousel-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
}

.carousel-content {
  position: relative;
  z-index: 1;
  text-align: center;
  color: white;
  padding: 20px;
}

.carousel-content h1 {
  font-size: 48px;
  margin-bottom: 20px;
  font-weight: bold;
}

.carousel-content p {
  font-size: 20px;
  margin-bottom: 30px;
  opacity: 0.95;
}

.features-section {
  padding: 60px 0;
  background: white;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.section-header {
  text-align: center;
  margin-bottom: 50px;
}

.section-header h2 {
  font-size: 36px;
  color: #333;
  margin-bottom: 12px;
}

.section-header p {
  color: #999;
  font-size: 16px;
}

.feature-card {
  text-align: center;
  padding: 30px;
  transition: all 0.3s;
}

.feature-card:hover {
  transform: translateY(-8px);
}

.feature-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
}

.feature-card h3 {
  font-size: 20px;
  color: #333;
  margin-bottom: 12px;
}

.feature-card p {
  color: #666;
  line-height: 1.6;
}

.services-section {
  padding: 60px 0;
}

.service-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
}

.service-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.service-image {
  height: 200px;
  overflow: hidden;
}

.service-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.service-card:hover .service-image img {
  transform: scale(1.1);
}

.service-content {
  padding: 24px;
  text-align: center;
}

.service-content h3 {
  font-size: 18px;
  color: #333;
  margin-bottom: 12px;
}

.service-desc {
  color: #666;
  line-height: 1.6;
  margin-bottom: 16px;
  font-size: 14px;
  min-height: 48px;
}

.service-price {
  margin-bottom: 16px;
}

.service-price .price {
  font-size: 28px;
  font-weight: bold;
  color: #f56c6c;
}

.service-price .unit {
  color: #999;
  font-size: 14px;
}

.btn-booking {
  width: 100%;
}

.process-section {
  padding: 60px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.process-section .section-header h2,
.process-section .section-header p {
  color: white;
}

.process-steps {
  display: flex;
  justify-content: space-between;
  gap: 20px;
}

.step-item {
  flex: 1;
  text-align: center;
  position: relative;
}

.step-number {
  font-size: 80px;
  font-weight: bold;
  color: rgba(255, 255, 255, 0.15);
  position: absolute;
  top: -30px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 0;
}

.step-icon {
  position: relative;
  z-index: 1;
  margin-bottom: 16px;
}

.step-item h4 {
  font-size: 18px;
  margin-bottom: 12px;
  position: relative;
  z-index: 1;
}

.step-item p {
  font-size: 14px;
  opacity: 0.9;
  position: relative;
  z-index: 1;
}

.counselor-join-section {
  padding: 60px 0;
}

.join-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  overflow: hidden;
}

.join-card :deep(.el-card__body) {
  padding: 40px;
}

.join-content h2 {
  font-size: 32px;
  color: white;
  margin-bottom: 16px;
}

.join-content > p {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 30px;
  line-height: 1.6;
}

.join-benefits {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.benefit-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: white;
  font-size: 15px;
  background: rgba(255, 255, 255, 0.2);
  padding: 10px 20px;
  border-radius: 20px;
}

.benefit-item .el-icon {
  font-size: 18px;
}

.join-action {
  text-align: center;
  padding-left: 40px;
}

.join-action h3 {
  font-size: 28px;
  color: white;
  margin-bottom: 12px;
}

.join-action > p {
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 24px;
}

.join-action .el-button {
  min-width: 160px;
  font-size: 16px;
  padding: 14px 40px;
}

.counselors-section {
  padding: 60px 0;
  background: white;
}

.counselors-section .section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.counselor-card {
  text-align: center;
  padding: 30px;
  border-radius: 12px;
  background: #f8f9fa;
  transition: all 0.3s;
  cursor: pointer;
}

.counselor-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.counselor-avatar {
  margin-bottom: 20px;
}

.counselor-card h3 {
  font-size: 20px;
  color: #333;
  margin-bottom: 8px;
}

.counselor-card .title {
  color: #666;
  font-size: 14px;
  margin-bottom: 12px;
}

.counselor-card .rating {
  margin: 16px 0;
}

.counselor-card .price {
  font-size: 24px;
  font-weight: bold;
  color: #f56c6c;
  margin: 16px 0 0;
}

.stats-section {
  padding: 80px 0;
  background: linear-gradient(135deg, #409eff 0%, #36d1dc 100%);
  color: white;
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 48px;
  font-weight: bold;
  margin-bottom: 12px;
}

.stat-label {
  font-size: 16px;
  opacity: 0.9;
}

:deep(.el-carousel__indicator--horizontal .el-carousel__button) {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

:deep(.el-carousel__indicator--horizontal.is-active .el-carousel__button) {
  background: #409eff;
}
</style>
