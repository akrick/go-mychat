<template>
  <div class="about-page">
    <AppHeader />

    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <h1>关于我们</h1>
        <p>专业、可靠、贴心的心理咨询服务平台</p>
      </div>
    </div>

    <!-- 公司简介 -->
    <div class="about-intro">
      <div class="container">
        <el-row :gutter="60">
          <el-col :span="12">
            <div class="intro-text">
              <h2>关于我们</h2>
              <p>MyChat 心理咨询平台成立于2020年，是国内领先的专业心理健康服务提供商。我们致力于通过互联网技术，让更多人能够便捷地获得高质量的心理咨询服务。</p>
              <p>我们的团队由心理学博士、资深心理咨询师、互联网技术专家组成，拥有丰富的临床经验和技术实力。我们相信，每个人都需要专业的心理支持，我们希望成为您心灵成长的伙伴。</p>
              <div class="intro-stats">
                <div class="stat-item">
                  <div class="stat-number">500+</div>
                  <div class="stat-label">专业咨询师</div>
                </div>
                <div class="stat-item">
                  <div class="stat-number">10万+</div>
                  <div class="stat-label">服务用户</div>
                </div>
                <div class="stat-item">
                  <div class="stat-number">98%</div>
                  <div class="stat-label">满意度</div>
                </div>
              </div>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="intro-image">
              <img src="https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=600&h=500&fit=crop" alt="关于我们" />
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- 我们的使命 -->
    <div class="mission-section">
      <div class="container">
        <div class="section-header">
          <h2>我们的使命与愿景</h2>
          <p>致力于让每个人都能获得专业的心理支持</p>
        </div>
        <el-row :gutter="40">
          <el-col :span="8">
            <div class="mission-card">
              <div class="mission-icon">
                <el-icon :size="50" color="#409eff"><Star /></el-icon>
              </div>
              <h3>我们的使命</h3>
              <p>通过提供专业、便捷、平价的心理咨询服务，帮助更多人解决心理困扰，提升心理健康水平，实现自我成长。</p>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="mission-card">
              <div class="mission-icon">
                <el-icon :size="50" color="#67c23a"><Trophy /></el-icon>
              </div>
              <h3>我们的愿景</h3>
              <p>成为中国最值得信赖的心理咨询服务平台，让心理健康服务触手可及，推动社会心理健康水平不断提升。</p>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="mission-card">
              <div class="mission-icon">
                <el-icon :size="50" color="#e6a23c"><Medal /></el-icon>
              </div>
              <h3>我们的价值观</h3>
              <p>专业、真诚、尊重、保密。我们以专业的态度对待每一位用户，真诚倾听，尊重差异，严格保密。</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- 团队介绍 -->
    <div class="team-section">
      <div class="container">
        <div class="section-header">
          <h2>核心团队</h2>
          <p>由心理学专家和互联网精英组成的专业团队</p>
        </div>
        <el-row :gutter="30">
          <el-col :span="6" v-for="member in teamMembers" :key="member.id">
            <div class="team-card">
              <div class="team-avatar">
                <img :src="member.avatar" :alt="member.name" />
              </div>
              <h3>{{ member.name }}</h3>
              <p class="role">{{ member.role }}</p>
              <p class="intro">{{ member.intro }}</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- 发展历程 -->
    <div class="timeline-section">
      <div class="container">
        <div class="section-header">
          <h2>发展历程</h2>
          <p>不断进步，持续创新</p>
        </div>
        <div class="timeline-container">
          <el-timeline>
            <el-timeline-item
              v-for="(event, index) in timeline"
              :key="index"
              :timestamp="event.year"
              placement="top"
              :color="event.color"
              size="large"
            >
              <div class="timeline-content">
                <h4>{{ event.title }}</h4>
                <p>{{ event.desc }}</p>
              </div>
            </el-timeline-item>
          </el-timeline>
        </div>
      </div>
    </div>

    <!-- 合作伙伴 -->
    <div class="partners-section">
      <div class="container">
        <div class="section-header">
          <h2>合作伙伴</h2>
          <p>与行业领先机构携手共创</p>
        </div>
        <el-row :gutter="40">
          <el-col :span="4" v-for="partner in partners" :key="partner.name">
            <div class="partner-item">
              <div class="partner-logo">
                <el-icon :size="40" color="#666"><OfficeBuilding /></el-icon>
              </div>
              <p>{{ partner.name }}</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- 联系我们 -->
    <div class="contact-section">
      <div class="container">
        <div class="section-header">
          <h2>联系我们</h2>
          <p>期待与您的沟通与合作</p>
        </div>
        <el-row :gutter="40">
          <el-col :span="8">
            <div class="contact-card">
              <div class="contact-icon">
                <el-icon :size="36" color="#409eff"><Phone /></el-icon>
              </div>
              <h4>客服热线</h4>
              <p class="contact-value">400-888-9999</p>
              <p class="contact-time">工作日 9:00-21:00</p>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="contact-card">
              <div class="contact-icon">
                <el-icon :size="36" color="#67c23a"><Message /></el-icon>
              </div>
              <h4>邮箱地址</h4>
              <p class="contact-value">contact@mychat.com</p>
              <p class="contact-time">24小时内回复</p>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="contact-card">
              <div class="contact-icon">
                <el-icon :size="36" color="#e6a23c"><Location /></el-icon>
              </div>
              <h4>公司地址</h4>
              <p class="contact-value">北京市朝阳区科技园区</p>
              <p class="contact-time">心理咨询中心大厦8层</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <AppFooter />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Star, Trophy, Medal, OfficeBuilding, Phone, Location, Message } from '@element-plus/icons-vue'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'

// 团队成员
const teamMembers = ref([
  {
    id: 1,
    name: '张明',
    role: '创始人 & CEO',
    intro: '心理学博士，15年临床经验',
    avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=300&h=300&fit=crop&crop=face'
  },
  {
    id: 2,
    name: '李华',
    role: '首席咨询师',
    intro: '国家二级心理咨询师',
    avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=300&h=300&fit=crop&crop=face'
  },
  {
    id: 3,
    name: '王芳',
    role: '技术总监',
    intro: '资深互联网技术专家',
    avatar: 'https://images.unsplash.com/photo-1580489944761-15a19d654956?w=300&h=300&fit=crop&crop=face'
  },
  {
    id: 4,
    name: '赵强',
    role: '运营总监',
    intro: '10年互联网运营经验',
    avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=300&h=300&fit=crop&crop=face'
  }
])

// 发展历程
const timeline = ref([
  {
    year: '2020年',
    title: '平台成立',
    desc: 'MyChat心理咨询服务平台正式成立，开启在线心理咨询服务新模式',
    color: '#409eff'
  },
  {
    year: '2021年',
    title: '用户突破',
    desc: '累计服务用户突破1万人，平台获得行业认可',
    color: '#67c23a'
  },
  {
    year: '2022年',
    title: '服务升级',
    desc: '推出语音咨询功能，服务产品线扩展至四大类',
    color: '#e6a23c'
  },
  {
    year: '2023年',
    title: '规模扩张',
    desc: '咨询师团队扩充至500人，服务用户突破5万',
    color: '#f56c6c'
  },
  {
    year: '2024年',
    title: '品牌升级',
    desc: '完成品牌升级，服务用户突破10万，满意度达到98%',
    color: '#909399'
  }
])

// 合作伙伴
const partners = ref([
  { name: '北京大学心理学院' },
  { name: '清华大学心理学系' },
  { name: '中科院心理所' },
  { name: '中国心理学会' },
  { name: '腾讯云' },
  { name: '阿里云' }
])
</script>

<style scoped>
.about-page {
  min-height: 100vh;
  background: #f5f7fa;
}

.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 80px 0 60px;
  text-align: center;
}

.page-header h1 {
  font-size: 48px;
  margin-bottom: 16px;
  font-weight: bold;
}

.page-header p {
  font-size: 18px;
  opacity: 0.9;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.about-intro {
  padding: 80px 0;
}

.intro-text h2 {
  font-size: 36px;
  color: #333;
  margin-bottom: 24px;
}

.intro-text p {
  color: #666;
  line-height: 1.8;
  margin-bottom: 16px;
  font-size: 15px;
}

.intro-stats {
  display: flex;
  gap: 40px;
  margin-top: 32px;
  padding-top: 32px;
  border-top: 1px solid #eee;
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 8px;
}

.stat-label {
  color: #666;
  font-size: 14px;
}

.intro-image img {
  width: 100%;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.mission-section {
  padding: 80px 0;
  background: white;
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

.mission-card {
  text-align: center;
  padding: 40px 24px;
  background: #f8f9fa;
  border-radius: 12px;
  transition: all 0.3s;
}

.mission-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.mission-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.mission-card h3 {
  font-size: 20px;
  color: #333;
  margin-bottom: 16px;
}

.mission-card p {
  color: #666;
  line-height: 1.6;
  font-size: 14px;
}

.team-section {
  padding: 80px 0;
}

.team-card {
  text-align: center;
  padding: 32px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
}

.team-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.team-avatar img {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
  margin-bottom: 16px;
  border: 3px solid #409eff;
}

.team-card h3 {
  font-size: 18px;
  color: #333;
  margin-bottom: 8px;
}

.team-card .role {
  color: #409eff;
  font-size: 14px;
  margin-bottom: 12px;
  font-weight: 500;
}

.team-card .intro {
  color: #666;
  font-size: 13px;
  line-height: 1.6;
}

.timeline-section {
  padding: 80px 0;
  background: white;
}

.timeline-container {
  max-width: 800px;
  margin: 0 auto;
}

.timeline-content h4 {
  font-size: 20px;
  color: #333;
  margin-bottom: 8px;
}

.timeline-content p {
  color: #666;
  line-height: 1.6;
  font-size: 15px;
}

:deep(.el-timeline-item__timestamp) {
  font-size: 18px;
  font-weight: bold;
  color: #333;
}

.partners-section {
  padding: 80px 0;
}

.partner-item {
  text-align: center;
  padding: 24px;
  background: white;
  border-radius: 8px;
  transition: all 0.3s;
}

.partner-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.partner-logo {
  margin-bottom: 12px;
}

.partner-item p {
  color: #666;
  font-size: 14px;
}

.contact-section {
  padding: 80px 0;
  background: white;
}

.contact-card {
  text-align: center;
  padding: 40px 24px;
  background: #f8f9fa;
  border-radius: 12px;
  transition: all 0.3s;
}

.contact-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.contact-card .contact-icon {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.contact-card h4 {
  font-size: 18px;
  color: #333;
  margin-bottom: 12px;
}

.contact-value {
  font-size: 20px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 8px;
}

.contact-time {
  color: #999;
  font-size: 13px;
}
</style>
