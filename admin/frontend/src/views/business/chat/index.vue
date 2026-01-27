<template>
  <div class="chat-container">
    <el-card>
      <el-form :inline="true" :model="queryParams" class="demo-form-inline">
        <el-form-item label="搜索">
          <el-input v-model="queryParams.keyword" placeholder="用户名/咨询师名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="全部" clearable>
            <el-option label="待开始" :value="0" />
            <el-option label="进行中" :value="1" />
            <el-option label="已结束" :value="2" />
            <el-option label="已超时" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" border style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户" width="120">
          <template #default="{ row }">
            {{ row.user?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="咨询师" width="120">
          <template #default="{ row }">
            {{ row.counselor?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getSessionStatusType(row.status)">
              {{ getSessionStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="时长(秒)" width="100" />
        <el-table-column prop="total_amount" label="金额(元)" width="100">
          <template #default="{ row }">
            ¥{{ row.total_amount?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="start_time" label="开始时间" width="180" />
        <el-table-column prop="end_time" label="结束时间" width="180" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="handleViewMessages(row)">查看消息</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="queryParams.page"
        v-model:page-size="queryParams.page_size"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleQuery"
        @current-change="handleQuery"
      />
    </el-card>

    <!-- 消息列表对话框 -->
    <el-dialog v-model="messagesDialogVisible" title="聊天消息" width="800px" top="5vh">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="用户">{{ currentSession?.user?.username }}</el-descriptions-item>
        <el-descriptions-item label="咨询师">{{ currentSession?.counselor?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getSessionStatusType(currentSession?.status)">
            {{ getSessionStatusText(currentSession?.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="时长">{{ currentSession?.duration }}秒</el-descriptions-item>
      </el-descriptions>

      <el-divider>消息记录</el-divider>

      <div class="messages-container" v-loading="messagesLoading">
        <div v-if="messages.length === 0" class="empty">暂无消息</div>
        <div v-for="msg in messages" :key="msg.id" class="message-item">
          <div class="message-header">
            <span class="sender">{{ msg.sender_type === 'user' ? currentSession?.user?.username : currentSession?.counselor?.name }}</span>
            <span class="time">{{ msg.created_at }}</span>
          </div>
          <div class="message-content">
            <span v-if="msg.content_type === 'text'">{{ msg.content }}</span>
            <span v-else-if="msg.content_type === 'image'">[图片]</span>
            <span v-else-if="msg.content_type === 'file'">[文件]</span>
          </div>
        </div>
      </div>

      <el-pagination
        v-model:current-page="messagesParams.page"
        v-model:page-size="messagesParams.page_size"
        :total="messagesTotal"
        :page-sizes="[20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadMessages"
        @current-change="loadMessages"
      />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getChatSessions, getChatMessages, deleteChatSession } from '@/api/adminChat'

const loading = ref(false)
const messagesLoading = ref(false)
const tableData = ref([])
const total = ref(0)
const messagesDialogVisible = ref(false)
const messages = ref([])
const messagesTotal = ref(0)
const currentSession = ref(null)

const queryParams = reactive({
  page: 1,
  page_size: 20,
  keyword: '',
  status: ''
})

const messagesParams = reactive({
  page: 1,
  page_size: 50
})

const handleQuery = async () => {
  loading.value = true
  try {
    const res = await getChatSessions(queryParams)
    tableData.value = res.sessions
    total.value = res.total
  } catch (error) {
    ElMessage.error(error.message || '获取会话列表失败')
  } finally {
    loading.value = false
  }
}

const resetQuery = () => {
  queryParams.keyword = ''
  queryParams.status = ''
  queryParams.page = 1
  handleQuery()
}

const getSessionStatusType = (status) => {
  const types = {
    0: 'info',
    1: 'success',
    2: 'primary',
    3: 'danger'
  }
  return types[status] || 'info'
}

const getSessionStatusText = (status) => {
  const texts = {
    0: '待开始',
    1: '进行中',
    2: '已结束',
    3: '已超时'
  }
  return texts[status] || '未知'
}

const handleViewMessages = async (row) => {
  currentSession.value = row
  messagesParams.page = 1
  messagesDialogVisible.value = true
  await loadMessages()
}

const loadMessages = async () => {
  messagesLoading.value = true
  try {
    const res = await getChatMessages(currentSession.value.id, messagesParams)
    messages.value = res.messages
    messagesTotal.value = res.total
  } catch (error) {
    ElMessage.error(error.message || '获取消息列表失败')
  } finally {
    messagesLoading.value = false
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该会话及其所有消息吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteChatSession(row.id)
      ElMessage.success('删除成功')
      handleQuery()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  })
}

onMounted(() => {
  handleQuery()
})
</script>

<style scoped>
.chat-container {
  padding: 20px;
}

.demo-form-inline {
  margin-bottom: 20px;
}

.el-pagination {
  margin-top: 20px;
}

.messages-container {
  max-height: 500px;
  overflow-y: auto;
  margin-bottom: 20px;
}

.message-item {
  padding: 12px;
  margin-bottom: 10px;
  background: #f5f7fa;
  border-radius: 4px;
}

.message-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 12px;
  color: #909399;
}

.sender {
  font-weight: bold;
  color: #409eff;
}

.message-content {
  color: #303133;
  white-space: pre-wrap;
  word-break: break-all;
}

.empty {
  text-align: center;
  padding: 40px;
  color: #909399;
}
</style>
