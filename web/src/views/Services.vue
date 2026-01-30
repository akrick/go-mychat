<template>
  <div class="services-page">
    <AppHeader />

    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <h1>服务介绍</h1>
        <p>专业的心理咨询服务，满足您的多样化需求</p>
      </div>
    </div>

    <!-- 服务详情 -->
    <div class="services-detail">
      <div class="container">
        <el-row :gutter="40">
          <el-col :span="12" v-for="service in serviceList" :key="service.id">
            <div class="service-item">
              <div class="service-icon" :style="{ background: service.bgColor }">
                <el-icon :size="50" :color="service.iconColor">
                  <component :is="service.icon" />
                </el-icon>
              </div>
              <div class="service-info">
                <h3>{{ service.title }}</h3>
                <p class="service-brief">{{ service.brief }}</p>
                <div class="service-features">
                  <div v-for="(feature, index) in service.features" :key="index" class="feature-tag">
                    <el-icon><Check /></el-icon>
                    {{ feature }}
                  </div>
                </div>
                <div class="service-price">
                  <span class="price">¥{{ service.price }}</span>
                  <span class="unit">/ {{ service.unit }}</span>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- 服务流程 -->
    <div class="service-process">
      <div class="container">
        <div class="section-header">
          <h2>服务流程</h2>
          <p>简单便捷，轻松获得专业心理支持</p>
        </div>
        <div class="process-timeline">
          <el-timeline>
            <el-timeline-item
              v-for="(step, index) in processSteps"
              :key="index"
              :timestamp="step.time"
              placement="top"
              :color="step.color"
            >
              <div class="timeline-content">
                <h4>{{ step.title }}</h4>
                <p>{{ step.desc }}</p>
              </div>
            </el-timeline-item>
          </el-timeline>
        </div>
      </div>
    </div>

    <!-- 常见问题 -->
    <div class="faq-section">
      <div class="container">
        <div class="section-header">
          <h2>常见问题</h2>
          <p>您关心的问题，这里都有答案</p>
        </div>
        <el-collapse v-model="activeFaq">
          <el-collapse-item
            v-for="(faq, index) in faqs"
            :key="index"
            :title="faq.question"
            :name="index"
          >
            <div class="faq-answer">{{ faq.answer }}</div>
          </el-collapse-item>
        </el-collapse>
      </div>
    </div>

    <!-- 联系我们 -->
    <div class="contact-section">
      <div class="container">
        <div class="section-header">
          <h2>联系我们</h2>
          <p>有任何问题，随时联系我们</p>
        </div>
        <el-row :gutter="40">
          <el-col :span="8">
            <div class="contact-item">
              <div class="contact-icon">
                <el-icon :size="40" color="#409eff"><Phone /></el-icon>
              </div>
              <h4>客服热线</h4>
              <p>400-888-9999</p>
              <p class="contact-desc">工作日 9:00-21:00</p>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="contact-item">
              <div class="contact-icon">
                <el-icon :size="40" color="#67c23a"><Message /></el-icon>
              </div>
              <h4>在线咨询</h4>
              <p>service@mychat.com</p>
              <p class="contact-desc">24小时内回复</p>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="contact-item">
              <div class="contact-icon">
                <el-icon :size="40" color="#e6a23c"><Location /></el-icon>
              </div>
              <h4>公司地址</h4>
              <p>北京市朝阳区科技园区</p>
              <p class="contact-desc">心理咨询中心</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- CTA区域 -->
    <div class="cta-section">
      <div class="container">
        <div class="cta-content">
          <h2>准备好开始您的咨询之旅了吗？</h2>
          <p>选择我们，让专业的心理咨询师陪伴您成长</p>
          <el-button type="primary" size="large" @click="handleBooking">
            立即预约咨询
          </el-button>
        </div>
      </div>
    </div>

    <AppFooter />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { UserFilled, Avatar, Briefcase, House, Check, Phone, Message, Location } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'

const router = useRouter()
const userStore = useUserStore()
const activeFaq = ref([0])

// 服务列表
const serviceList = ref([
  {
    id: 1,
    icon: UserFilled,
    title: '个人心理咨询',
    brief: '一对一专业心理咨询，帮助您解决情绪困扰，提升自我认知，实现个人成长',
    price: '300',
    unit: '小时',
    bgColor: '#e8f4fd',
    iconColor: '#409eff',
    features: ['持证心理咨询师', '严格保密', '灵活预约', '视频/文字咨询']
  },
  {
    id: 2,
    icon: Avatar,
    title: '情感咨询',
    brief: '专业情感指导，帮助您改善亲密关系，解决情感困惑，建立健康的人际关系',
    price: '400',
    unit: '小时',
    bgColor: '#fee8f0',
    iconColor: '#f56c6c',
    features: ['情感关系分析', '沟通技巧指导', '心理支持', '长期陪伴']
  },
  {
    id: 3,
    icon: Briefcase,
    title: '职业规划',
    brief: '专业的职业发展方向指导，帮助您明确职业目标，规划职业发展路径',
    price: '500',
    unit: '小时',
    bgColor: '#fef9e8',
    iconColor: '#e6a23c',
    features: ['职业测评', '发展方向规划', '求职指导', '职业困惑解答']
  },
  {
    id: 4,
    icon: House,
    title: '家庭治疗',
    brief: '家庭关系改善与治疗，促进家庭和谐，解决家庭矛盾，改善家庭互动模式',
    price: '600',
    unit: '小时',
    bgColor: '#e8f8f3',
    iconColor: '#67c23a',
    features: ['家庭系统治疗', '亲子关系指导', '婚姻咨询', '家庭矛盾调解']
  }
])

// 服务流程
const processSteps = ref([
  {
    title: '注册登录',
    desc: '完成账号注册，填写基本信息，成为平台会员',
    time: '第一步',
    color: '#409eff'
  },
  {
    title: '选择咨询师',
    desc: '浏览咨询师列表，根据专业领域、评分等选择合适的咨询师',
    time: '第二步',
    color: '#67c23a'
  },
  {
    title: '提交预约',
    desc: '选择服务类型、咨询时长和时间段，提交预约申请',
    time: '第三步',
    color: '#e6a23c'
  },
  {
    title: '在线支付',
    desc: '确认订单信息，完成在线支付，预约成功',
    time: '第四步',
    color: '#f56c6c'
  },
  {
    title: '开始咨询',
    desc: '在预约时间通过文字或语音与咨询师进行咨询',
    time: '第五步',
    color: '#909399'
  },
  {
    title: '反馈评价',
    desc: '咨询结束后，对咨询师的服务进行评价和反馈',
    time: '第六步',
    color: '#909399'
  }
])

// 常见问题
const faqs = ref([
  {
    question: '心理咨询需要多长时间？',
    answer: '心理咨询的时长因人而异。一般来说，初次咨询为1小时，后续咨询根据个人情况确定频率和次数。我们的服务以小时为单位计费，您可以根据需求选择。'
  },
  {
    question: '咨询内容会保密吗？',
    answer: '绝对保密。我们的咨询师都签署了严格的保密协议，您的所有咨询内容都受到法律保护，不会向任何第三方泄露。仅在涉及您或他人生命安全等极端情况下，才会在最小范围内进行必要的信息披露。'
  },
  {
    question: '如何选择合适的咨询师？',
    answer: '您可以根据咨询师的专业领域、受训背景、从业年限、用户评价等信息进行选择。如果您不确定，可以先选择一位进行初次咨询，体验后再决定是否继续。我们也会根据您的需求为您推荐合适的咨询师。'
  },
  {
    question: '如果对咨询效果不满意怎么办？',
    answer: '如果对咨询效果不满意，您可以随时更换咨询师。我们提供满意度评价系统，您的反馈将帮助我们不断改进服务质量。如有特殊情况，也可联系客服申请退款。'
  },
  {
    question: '咨询方式有哪些？',
    answer: '我们目前提供文字聊天和语音通话两种咨询方式。您可以在预约时选择适合自己的方式。文字咨询方便记录和回顾，语音咨询更接近面对面交流，效果更直接。'
  },
  {
    question: '心理咨询能解决所有问题吗？',
    answer: '心理咨询可以帮助您理解和处理很多心理困扰，但不是万能的。对于严重的心理疾病，如重度抑郁、焦虑症等，我们建议您先就医，配合药物治疗。心理咨询可以作为重要的辅助治疗手段。'
  }
])

const handleBooking = () => {
  if (!userStore.token) {
    ElMessage.warning('请先登录')
    router.push('/login')
  } else {
    router.push('/counselors')
  }
}
</script>

<style scoped>
.services-page {
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

.services-detail {
  padding: 80px 0;
}

.service-item {
  display: flex;
  gap: 24px;
  background: white;
  padding: 32px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  margin-bottom: 24px;
  transition: all 0.3s;
}

.service-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.service-icon {
  flex-shrink: 0;
  width: 100px;
  height: 100px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.service-info {
  flex: 1;
}

.service-info h3 {
  font-size: 24px;
  color: #333;
  margin-bottom: 12px;
}

.service-brief {
  color: #666;
  line-height: 1.6;
  margin-bottom: 16px;
  font-size: 15px;
}

.service-features {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 20px;
}

.feature-tag {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #f5f7fa;
  border-radius: 16px;
  font-size: 13px;
  color: #666;
}

.service-price {
  padding-top: 16px;
  border-top: 1px solid #eee;
}

.service-price .price {
  font-size: 32px;
  font-weight: bold;
  color: #f56c6c;
}

.service-price .unit {
  color: #999;
  font-size: 14px;
}

.service-process {
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

.process-timeline {
  max-width: 800px;
  margin: 0 auto;
}

.timeline-content h4 {
  font-size: 18px;
  color: #333;
  margin-bottom: 8px;
}

.timeline-content p {
  color: #666;
  line-height: 1.6;
}

.faq-section {
  padding: 80px 0;
}

.faq-answer {
  color: #666;
  line-height: 1.8;
  font-size: 15px;
  padding: 12px 0;
}

:deep(.el-collapse-item__header) {
  font-size: 16px;
  color: #333;
  padding: 16px 20px;
  background: white;
}

:deep(.el-collapse-item__wrap) {
  background: white;
}

.contact-section {
  padding: 80px 0;
  background: white;
}

.contact-item {
  text-align: center;
  padding: 32px;
  border-radius: 12px;
  background: #f8f9fa;
  transition: all 0.3s;
}

.contact-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.contact-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.contact-item h4 {
  font-size: 18px;
  color: #333;
  margin-bottom: 8px;
}

.contact-item p {
  color: #666;
  margin-bottom: 4px;
  font-size: 15px;
}

.contact-desc {
  font-size: 13px;
  color: #999;
  margin-top: 8px;
}

.cta-section {
  padding: 100px 0;
  background: linear-gradient(135deg, #409eff 0%, #36d1dc 100%);
  color: white;
  text-align: center;
}

.cta-content h2 {
  font-size: 36px;
  margin-bottom: 16px;
}

.cta-content > p {
  font-size: 18px;
  margin-bottom: 32px;
  opacity: 0.95;
}

.cta-section .el-button {
  padding: 16px 48px;
  font-size: 18px;
  border-radius: 30px;
}
</style>
