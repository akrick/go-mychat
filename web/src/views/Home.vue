<template>
  <div class="home-page">
    <AppHeader />
    <div class="home">
      <el-row :gutter="20" class="hero-section">
        <el-col :span="12">
          <div class="hero-content">
            <h1>专业的心理咨询平台</h1>
            <p>连接您与专业心理咨询师，开启心灵成长之旅</p>
            <div class="hero-actions">
              <el-button type="primary" size="large" @click="$router.push('/counselors')">
                寻找咨询师
              </el-button>
              <el-button size="large" @click="$router.push('/orders')">
                我的订单
              </el-button>
            </div>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="hero-image">
            <el-icon size="300"><ChatDotRound /></el-icon>
          </div>
        </el-col>
      </el-row>

      <div class="features-section">
        <h2>平台特色</h2>
        <el-row :gutter="20">
          <el-col :span="8" v-for="feature in features" :key="feature.title">
            <el-card class="feature-card">
              <div class="feature-icon">
                <el-icon :size="48" :color="feature.color">
                  <component :is="feature.icon" />
                </el-icon>
              </div>
              <h3>{{ feature.title }}</h3>
              <p>{{ feature.desc }}</p>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <div class="counselors-section">
        <div class="section-header">
          <h2>热门咨询师</h2>
          <el-button text @click="$router.push('/counselors')">查看全部 →</el-button>
        </div>
        <el-row :gutter="20" v-loading="loading">
          <el-col :span="6" v-for="counselor in counselors" :key="counselor.id">
            <el-card class="counselor-card" @click="$router.push(`/counselor/${counselor.id}`)" style="cursor: pointer">
              <div class="counselor-avatar">
                <el-avatar :size="80" :src="counselor.avatar">
                  <el-icon><User /></el-icon>
                </el-avatar>
              </div>
              <h3>{{ counselor.name }}</h3>
              <p class="title">{{ counselor.title }}</p>
              <div class="rating">
                <el-rate v-model="counselor.rating" disabled show-score text-color="#ff9900" />
              </div>
              <p class="price">¥{{ counselor.price }}/分钟</p>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ChatDotRound, User, Clock, Lock } from '@element-plus/icons-vue'
import { getCounselorList } from '@/api/counselor'
import { handleError } from '@/utils/errorHandler'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'

const loading = ref(false)
const counselors = ref([])

const features = [
  {
    icon: User,
    title: '专业团队',
    desc: '持证心理咨询师，经验丰富',
    color: '#409eff'
  },
  {
    icon: Clock,
    title: '灵活预约',
    desc: '随时随地，轻松预约',
    color: '#67c23a'
  },
  {
    icon: Lock,
    title: '隐私保护',
    desc: '严格保密，安全可靠',
    color: '#e6a23c'
  }
]

onMounted(async () => {
  await loadCounselors()
})

const loadCounselors = async () => {
  loading.value = true
  try {
    const res = await getCounselorList({ page: 1, page_size: 4 })
    counselors.value = res.data.counselors
  } catch (error) {
    handleError(error, '加载咨询师失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.home {
  flex: 1;
  padding: 20px 0;
}

.hero-section {
  margin-bottom: 60px;
  padding: 40px 0;
}

.hero-content {
  padding: 40px;
}

.hero-content h1 {
  font-size: 48px;
  margin-bottom: 16px;
  color: #333;
}

.hero-content p {
  font-size: 20px;
  color: #666;
  margin-bottom: 32px;
}

.hero-actions {
  display: flex;
  gap: 16px;
}

.hero-image {
  display: flex;
  justify-content: center;
  align-items: center;
  color: #409eff;
}

.features-section {
  margin-bottom: 60px;
}

.features-section h2 {
  text-align: center;
  margin-bottom: 40px;
  font-size: 32px;
  color: #333;
}

.feature-card {
  text-align: center;
  padding: 20px;
  height: 100%;
}

.feature-icon {
  margin-bottom: 16px;
}

.feature-card h3 {
  margin-bottom: 12px;
  color: #333;
}

.feature-card p {
  color: #666;
  line-height: 1.6;
}

.counselors-section h2 {
  margin-bottom: 32px;
  font-size: 32px;
  color: #333;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.counselor-card {
  text-align: center;
  margin-bottom: 20px;
  transition: transform 0.3s, box-shadow 0.3s;
}

.counselor-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
}

.counselor-avatar {
  margin-bottom: 16px;
}

.counselor-card h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  color: #333;
}

.counselor-card .title {
  margin: 0 0 12px 0;
  color: #666;
  font-size: 14px;
}

.counselor-card .rating {
  margin-bottom: 12px;
}

.counselor-card .price {
  margin: 0;
  font-size: 20px;
  font-weight: bold;
  color: #f56c6c;
}
</style>
