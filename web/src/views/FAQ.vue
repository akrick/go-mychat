<template>
  <div class="faq-page">
    <AppHeader />
    <div class="faq">
      <div class="page-header">
        <div class="header-content">
          <h1>常见问题</h1>
          <p class="header-desc">为您解答使用过程中的常见疑问</p>
        </div>
      </div>

      <el-row :gutter="20" class="faq-container">
        <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
          <el-card class="category-card">
            <template #header>
              <h3>问题分类</h3>
            </template>
            <el-menu
              :default-active="activeCategory"
              @select="handleCategoryChange"
              class="category-menu"
            >
              <el-menu-item index="all">
                <el-icon><List /></el-icon>
                全部问题
              </el-menu-item>
              <el-menu-item index="account">
                <el-icon><User /></el-icon>
                账户相关
              </el-menu-item>
              <el-menu-item index="booking">
                <el-icon><Calendar /></el-icon>
                预约咨询
              </el-menu-item>
              <el-menu-item index="payment">
                <el-icon><Wallet /></el-icon>
                支付问题
              </el-menu-item>
              <el-menu-item index="chat">
                <el-icon><ChatDotRound /></el-icon>
                咨询会话
              </el-menu-item>
              <el-menu-item index="other">
                <el-icon><QuestionFilled /></el-icon>
                其他问题
              </el-menu-item>
            </el-menu>
          </el-card>
        </el-col>

        <el-col :xs="24" :sm="24" :md="18" :lg="18" :xl="18">
          <div class="faq-list">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索问题关键词..."
              prefix-icon="Search"
              class="search-input"
              @input="handleSearch"
            />

            <el-collapse v-model="activeItems" class="faq-collapse">
              <el-collapse-item
                v-for="faq in filteredFaqs"
                :key="faq.id"
                :name="faq.id"
              >
                <template #title>
                  <div class="faq-title">
                    <el-icon><QuestionFilled /></el-icon>
                    {{ faq.question }}
                  </div>
                </template>
                <div class="faq-answer">{{ faq.answer }}</div>
              </el-collapse-item>
            </el-collapse>

            <el-empty
              v-if="filteredFaqs.length === 0"
              description="没有找到相关问题"
            />
          </div>
        </el-col>
      </el-row>

      <el-card class="contact-card">
        <div class="contact-content">
          <h3>没有找到答案？</h3>
          <p>如果您的问题未在上述内容中得到解答，欢迎联系我们</p>
          <div class="contact-buttons">
            <el-button type="primary" @click="$router.push('/about')">
              <el-icon><Phone /></el-icon>
              联系客服
            </el-button>
            <el-button @click="$router.push('/register')">
              <el-icon><User /></el-icon>
              注册账号
            </el-button>
          </div>
        </div>
      </el-card>
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { List, User, Calendar, Wallet, ChatDotRound, QuestionFilled, Phone } from '@element-plus/icons-vue'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'

const activeCategory = ref('all')
const searchKeyword = ref('')
const activeItems = ref([])

const faqData = [
  {
    id: 1,
    category: 'account',
    question: '如何注册账号？',
    answer: '点击页面右上角的"注册"按钮，填写用户名、密码等信息，即可完成注册。注册后请及时完善个人信息。'
  },
  {
    id: 2,
    category: 'account',
    question: '忘记密码怎么办？',
    answer: '在登录页面点击"忘记密码"，通过注册手机号或邮箱找回密码。如果仍有问题，请联系客服。'
  },
  {
    id: 3,
    category: 'booking',
    question: '如何预约咨询师？',
    answer: '在"咨询师"页面浏览咨询师列表，选择合适的咨询师进入详情页，点击"预约咨询"按钮，选择咨询时长和预约时间即可。'
  },
  {
    id: 4,
    category: 'booking',
    question: '预约后可以取消吗？',
    answer: '可以在"我的订单"中查看订单状态，对于未支付的订单可以直接取消。已支付的订单需要联系客服处理。'
  },
  {
    id: 5,
    category: 'payment',
    question: '支持哪些支付方式？',
    answer: '目前支持微信支付和支付宝支付。我们正在陆续接入更多支付方式，敬请期待。'
  },
  {
    id: 6,
    category: 'payment',
    question: '支付失败怎么办？',
    answer: '请检查账户余额是否充足，或更换支付方式重试。如持续失败，请联系客服并提供订单号。'
  },
  {
    id: 7,
    category: 'chat',
    question: '如何开始咨询会话？',
    answer: '订单支付成功后，在订单列表中点击"进入咨询"即可开始与咨询师的实时聊天。'
  },
  {
    id: 8,
    category: 'chat',
    question: '咨询时长有限制吗？',
    answer: '您可以在预约时选择30分钟、60分钟、90分钟或120分钟的咨询时长。咨询结束后系统会自动关闭会话。'
  },
  {
    id: 9,
    category: 'chat',
    question: '咨询记录会保存吗？',
    answer: '是的，您的所有咨询记录都会安全保存在系统中，您可以随时在个人中心查看历史咨询记录。'
  },
  {
    id: 10,
    category: 'other',
    question: '个人信息会保密吗？',
    answer: '我们严格遵守隐私保护法规，您的所有个人信息和咨询内容都会严格保密，不会向第三方透露。'
  }
]

const filteredFaqs = computed(() => {
  let result = faqData

  if (activeCategory.value !== 'all') {
    result = result.filter(faq => faq.category === activeCategory.value)
  }

  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(faq =>
      faq.question.toLowerCase().includes(keyword) ||
      faq.answer.toLowerCase().includes(keyword)
    )
  }

  return result
})

const handleCategoryChange = (category) => {
  activeCategory.value = category
  activeItems.value = []
}

const handleSearch = () => {
  activeItems.value = []
}
</script>

<style scoped>
.faq-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.faq {
  flex: 1;
  padding: 20px 0;
}

.page-header {
  padding: 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  margin-bottom: 24px;
  text-align: center;
}

.page-header h1 {
  margin: 0 0 12px 0;
  font-size: 32px;
  font-weight: 600;
}

.header-desc {
  margin: 0;
  font-size: 16px;
  opacity: 0.9;
}

.faq-container {
  margin-bottom: 24px;
}

.category-card h3 {
  margin: 0;
  color: #333;
}

.category-menu {
  border: none;
}

.faq-list {
  background: white;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.search-input {
  margin-bottom: 24px;
}

.faq-collapse {
  border: none;
}

.faq-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.faq-answer {
  padding: 16px;
  color: #666;
  line-height: 1.8;
  background: #f5f7fa;
  border-radius: 8px;
  margin-top: 12px;
}

.contact-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.contact-content {
  text-align: center;
}

.contact-content h3 {
  margin: 0 0 12px 0;
  font-size: 24px;
}

.contact-content p {
  margin: 0 0 24px 0;
  opacity: 0.9;
}

.contact-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.contact-buttons .el-button {
  color: white;
  border-color: rgba(255, 255, 255, 0.5);
}

.contact-buttons .el-button--primary {
  background: white;
  color: #667eea;
  border-color: white;
}

.contact-buttons .el-button:hover {
  opacity: 0.9;
}

@media (max-width: 768px) {
  .page-header {
    padding: 24px;
  }

  .page-header h1 {
    font-size: 24px;
  }

  .faq-container > .el-col {
    margin-bottom: 20px;
  }

  .contact-buttons {
    flex-direction: column;
  }
}
</style>
