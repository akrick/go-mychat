<template>
  <div class="chat-page">
    <AppHeader />
    <div class="chat">
      <el-page-header @back="$router.back()" content="咨询会话" class="page-header" />

      <el-row :gutter="20" class="chat-container">
        <el-col :span="18">
          <el-card class="message-card">
            <template #header>
              <div class="chat-header">
                <span>咨询会话</span>
                <el-tag type="success">进行中</el-tag>
              </div>
            </template>

            <div class="message-list" ref="messageListRef">
              <div v-if="messages.length === 0" class="empty-message">
                暂无消息，开始咨询吧
              </div>
              <div
                v-for="message in messages"
                :key="message.id"
                :class="['message-item', message.sender_type === 'user' ? 'my-message' : 'other-message']"
              >
                <el-avatar :size="40" :src="message.sender?.avatar">
                  <el-icon><User /></el-icon>
                </el-avatar>
                <div class="message-content">
                  <p class="sender-name">{{ message.sender?.username }}</p>
                  <div class="message-bubble">{{ message.content }}</div>
                </div>
              </div>
            </div>

            <div class="message-input">
              <el-input
                v-model="inputMessage"
                type="textarea"
                :rows="3"
                placeholder="输入消息..."
                @keydown.enter.prevent="handleSendMessage"
              />
              <div class="input-actions">
                <span class="remaining-time" v-if="remainingTime > 0">
                  <el-icon><Timer /></el-icon>
                  剩余时间：{{ formatSeconds(remainingTime) }}
                </span>
                <el-button type="primary" :loading="sending" @click="handleSendMessage">
                  发送
                </el-button>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="6">
          <el-card class="info-card">
            <template #header>
              <h3>会话信息</h3>
            </template>
            <div class="info-content">
              <div class="info-item">
                <span class="label">咨询师：</span>
                <span class="value">{{ sessionInfo.counselor_name }}</span>
              </div>
              <div class="info-item">
                <span class="label">咨询时长：</span>
                <span class="value">{{ sessionInfo.duration }}分钟</span>
              </div>
              <div class="countdown-section">
                <h4>倒计时</h4>
                <div class="countdown" :class="{ warning: remainingTime < 300, danger: remainingTime < 60 }">
                  {{ formatSeconds(remainingTime) }}
                </div>
                <p class="countdown-tip">剩余时间</p>
              </div>
            </div>
            <el-button type="danger" @click="handleEndSession" style="width: 100%">
              结束咨询
            </el-button>
          </el-card>
        </el-col>
      </el-row>
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { User, Timer } from '@element-plus/icons-vue'
import { getMessages, sendMessage, endChatSession } from '@/api/chat'
import { formatSeconds } from '@/utils/formatter'
import { handleError, showWarning, showSuccess } from '@/utils/errorHandler'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'

const route = useRoute()
const router = useRouter()

const sessionId = ref(route.params.sessionId)
const messages = ref([])
const loading = ref(false)
const sending = ref(false)
const inputMessage = ref('')
const messageListRef = ref(null)
const countdownTimer = ref(null)
const remainingTime = ref(0)

const sessionInfo = reactive({
  counselor_name: '',
  duration: 60,
  start_time: Date.now()
})

onMounted(async () => {
  await loadMessages()
  startCountdown()
  scrollBottom()
})

onUnmounted(() => {
  if (countdownTimer.value) {
    clearInterval(countdownTimer.value)
  }
})

const loadMessages = async () => {
  loading.value = true
  try {
    const res = await getMessages(sessionId.value, { page: 1, page_size: 100 })
    messages.value = res.data.messages
    // 从消息中获取会话信息
    if (messages.value.length > 0 && messages.value[0].session) {
      sessionInfo.duration = messages.value[0].session.duration || 60
      sessionInfo.counselor_name = messages.value[0].session.counselor?.name || '咨询师'
      sessionInfo.start_time = messages.value[0].session.start_time || Date.now()
    }
    calculateRemainingTime()
  } catch (error) {
    handleError(error, '加载消息失败')
  } finally {
    loading.value = false
  }
}

const calculateRemainingTime = () => {
  const duration = sessionInfo.duration * 60
  const elapsed = (Date.now() - new Date(sessionInfo.start_time).getTime()) / 1000
  remainingTime.value = Math.max(0, duration - elapsed)
}

const startCountdown = () => {
  countdownTimer.value = setInterval(() => {
    remainingTime.value -= 1
    if (remainingTime.value <= 0) {
      clearInterval(countdownTimer.value)
      showWarning('咨询时间已结束')
      handleEndSession()
    }
  }, 1000)
}



const handleSendMessage = async () => {
  if (!inputMessage.value.trim()) {
    showWarning('请输入消息内容')
    return
  }

  sending.value = true
  try {
    await sendMessage(sessionId.value, {
      content: inputMessage.value,
      sender_type: 'user',
      content_type: 'text'
    })
    inputMessage.value = ''
    await loadMessages()
    scrollBottom()
  } catch (error) {
    handleError(error, '发送消息失败')
  } finally {
    sending.value = false
  }
}

const handleEndSession = async () => {
  try {
    await ElMessageBox.confirm('确定要结束咨询会话吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await endChatSession(sessionId.value)
    showSuccess('咨询已结束')
    router.push('/orders')
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, '结束会话失败')
    }
  }
}

const scrollBottom = () => {
  nextTick(() => {
    if (messageListRef.value) {
      messageListRef.value.scrollTop = messageListRef.value.scrollHeight
    }
  })
}
</script>

<style scoped>
.chat-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.chat {
  flex: 1;
  padding: 20px 0;
}

.page-header {
  margin-bottom: 20px;
}

.chat-container {
  margin-top: 20px;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.message-card {
  height: calc(100vh - 200px);
  display: flex;
  flex-direction: column;
}

.message-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: #f5f7fa;
}

.empty-message {
  text-align: center;
  color: #999;
  padding: 40px 0;
}

.message-item {
  display: flex;
  margin-bottom: 20px;
}

.my-message {
  flex-direction: row-reverse;
}

.my-message .message-content {
  text-align: right;
}

.my-message .message-bubble {
  background: #409eff;
  color: white;
}

.other-message .message-bubble {
  background: white;
}

.message-content {
  margin: 0 12px;
  max-width: 60%;
}

.sender-name {
  margin: 0 0 6px 0;
  font-size: 12px;
  color: #999;
}

.message-bubble {
  padding: 12px 16px;
  border-radius: 8px;
  line-height: 1.6;
  word-break: break-all;
}

.message-input {
  padding: 16px;
  border-top: 1px solid #ebeef5;
  background: white;
}

.input-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
}

.remaining-time {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #666;
  font-size: 14px;
}

.info-card h3 {
  margin: 0;
  color: #333;
}

.info-content {
  margin-bottom: 20px;
}

.info-item {
  margin-bottom: 16px;
  color: #666;
}

.info-item .label {
  display: inline-block;
  width: 80px;
  color: #999;
}

.info-item .value {
  color: #333;
  font-weight: 500;
}

.countdown-section {
  text-align: center;
  padding: 20px 0;
  border-top: 1px solid #ebeef5;
  margin-top: 20px;
}

.countdown-section h4 {
  margin: 0 0 16px 0;
  color: #666;
  font-size: 14px;
}

.countdown {
  font-size: 48px;
  font-weight: bold;
  color: #67c23a;
  margin-bottom: 8px;
}

.countdown.warning {
  color: #e6a23c;
}

.countdown.danger {
  color: #f56c6c;
}

.countdown-tip {
  margin: 0;
  color: #999;
  font-size: 13px;
}
</style>
