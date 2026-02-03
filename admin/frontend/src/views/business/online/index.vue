<template>
  <div class="online-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h2>在线用户管理</h2>
        <p>实时监控和管理平台在线用户</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleSendMessage">
          <el-icon><ChatDotRound /></el-icon>
          发送广播
        </el-button>
        <el-button @click="loadUsers" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon total">
              <el-icon :size="32"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.total_online || 0 }}</div>
              <div class="stat-label">在线总人数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon counselor">
              <el-icon :size="32"><Avatar /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.counselor_online || 0 }}</div>
              <div class="stat-label">在线咨询师</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon regular">
              <el-icon :size="32"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.user_online || 0 }}</div>
              <div class="stat-label">在线用户</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon muted">
              <el-icon :size="32"><MuteNotification /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.muted_counselor || 0 }}</div>
              <div class="stat-label">被禁言咨询师</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 在线用户列表 -->
    <el-card class="user-list-card">
      <template #header>
        <div class="card-header">
          <span>在线用户列表</span>
          <el-tag type="info">{{ users.length }}人在线</el-tag>
        </div>
      </template>

      <el-table :data="users" v-loading="loading" stripe>
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="username" label="用户名" min-width="120">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="32">
                <el-icon><User /></el-icon>
              </el-avatar>
              <span class="username">{{ row.username }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="role" label="角色" width="120" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.role === 'counselor'" type="success">咨询师</el-tag>
            <el-tag v-else type="primary">用户</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="online_time" label="在线时长" width="150" />
        <el-table-column prop="session_id" label="当前会话" width="120" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.session_id" type="warning" size="small">咨询中</el-tag>
            <span v-else class="empty-text">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="is_muted" label="状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_muted" type="danger">已禁言</el-tag>
            <el-tag v-else type="success">正常</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button
                type="danger"
                size="small"
                @click="handleKickUser(row)"
              >
                <el-icon><SwitchButton /></el-icon>
                踢下线
              </el-button>
              <el-button
                v-if="row.role === 'counselor'"
                :type="row.is_muted ? 'success' : 'warning'"
                size="small"
                @click="handleMuteUser(row)"
              >
                <el-icon>
                  <MuteNotification v-if="!row.is_muted" />
                  <Select v-else />
                </el-icon>
                {{ row.is_muted ? '解禁' : '禁言' }}
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && users.length === 0" description="暂无在线用户" />
    </el-card>

    <!-- 广播消息对话框 -->
    <el-dialog v-model="messageDialogVisible" title="发送系统广播消息" width="600px">
      <el-form :model="messageForm" label-width="80px">
        <el-form-item label="广播范围">
          <el-tag type="info">所有在线用户</el-tag>
        </el-form-item>
        <el-form-item label="消息内容" prop="content">
          <el-input
            v-model="messageForm.content"
            type="textarea"
            :rows="6"
            maxlength="500"
            show-word-limit
            placeholder="请输入要发送的系统消息..."
          />
        </el-form-item>
        <el-alert
          title="提示"
          type="info"
          :closable="false"
          show-icon
        >
          此消息将广播给所有在线用户，包括普通用户和咨询师
        </el-alert>
      </el-form>
      <template #footer>
        <el-button @click="messageDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSendMessageSubmit" :loading="sending">
          发送广播
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Refresh,
  User,
  Avatar,
  MuteNotification,
  SwitchButton,
  Select,
  ChatDotRound
} from '@element-plus/icons-vue'
import { getOnlineUsers, kickUser, muteUser, broadcastMessage, getOnlineStatistics } from '@/api/online'
import { handleError, showSuccess } from '@/utils/errorHandler'

const loading = ref(false)
const users = ref([])
const currentUser = ref(null)
const messageDialogVisible = ref(false)
const sending = ref(false)

const statistics = reactive({
  total_online: 0,
  counselor_online: 0,
  user_online: 0,
  muted_counselor: 0
})

const messageForm = reactive({
  content: ''
})

let refreshTimer = null

onMounted(async () => {
  console.log('在线用户页面已加载')
  await loadUsers()
  await loadStatistics()
  startAutoRefresh()
})

onUnmounted(() => {
  console.log('在线用户页面已卸载')
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})

const startAutoRefresh = () => {
  // 清除之前的定时器（如果有）
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }

  // 设置新的定时器
  refreshTimer = setInterval(() => {
    console.log('自动刷新在线用户数据')
    loadUsers()
    loadStatistics()
  }, 30000) // 30秒自动刷新
}

const loadUsers = async () => {
  loading.value = true
  try {
    const res = await getOnlineUsers()
    console.log('在线用户数据:', res)
    if (res && res.code === 200 && res.data) {
      users.value = res.data.users || []
    } else {
      ElMessage.warning(res?.msg || '获取在线用户数据异常')
      users.value = []
    }
  } catch (error) {
    console.error('加载在线用户失败:', error)
    handleError(error, '加载在线用户失败')
    users.value = []
  } finally {
    loading.value = false
  }
}

const loadStatistics = async () => {
  try {
    const res = await getOnlineStatistics()
    console.log('统计数据:', res)
    if (res && res.code === 200 && res.data) {
      Object.assign(statistics, res.data)
    } else {
      console.warn('获取统计数据异常:', res)
    }
  } catch (error) {
    console.error('加载统计信息失败:', error)
  }
}

const handleKickUser = async (user) => {
  try {
    await ElMessageBox.confirm(
      `确定要踢用户 "${user.username}" 下线吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await kickUser(user.user_id)
    showSuccess('踢下线成功')
    await loadUsers()
    await loadStatistics()
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, '踢下线失败')
    }
  }
}

const handleMuteUser = async (user) => {
  const action = user.is_muted ? '解禁' : '禁言'
  try {
    await ElMessageBox.confirm(
      `确定要${action}咨询师 "${user.username}" 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await muteUser({
      user_id: user.user_id,
      is_muted: !user.is_muted
    })

    showSuccess(`${action}成功`)
    await loadUsers()
    await loadStatistics()
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, `${action}失败`)
    }
  }
}

const handleSendMessage = () => {
  messageForm.content = ''
  messageDialogVisible.value = true
}

const handleSendMessageSubmit = async () => {
  if (!messageForm.content.trim()) {
    ElMessage.warning('请输入消息内容')
    return
  }

  sending.value = true
  try {
    await broadcastMessage({
      content: messageForm.content
    })
    showSuccess('系统广播消息已发送')
    messageDialogVisible.value = false
  } catch (error) {
    console.error('发送广播消息失败:', error)
    handleError(error, '发送广播消息失败')
  } finally {
    sending.value = false
  }
}
</script>

<style scoped>
.online-container {
  padding: 20px;
  background: #f5f7fa;
  min-height: 100vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-content h2 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 20px;
}

.header-content p {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-icon.counselor {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.stat-icon.regular {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.stat-icon.muted {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #333;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.user-list-card {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username {
  font-weight: 500;
  color: #333;
}

.empty-text {
  color: #c0c4cc;
  font-size: 14px;
}

:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.el-table__header) {
  background: #f8f9fa;
}

:deep(.el-table__row:hover) {
  background: #f0f7ff;
}
</style>
