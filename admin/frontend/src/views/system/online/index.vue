<template>
  <div class="online-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>在线用户</span>
          <el-button type="primary" size="small" @click="handleBroadcast">
            <el-icon><Promotion /></el-icon>
            广播消息
          </el-button>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="20" class="stats-row">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.total }}</div>
            <div class="stat-label">总在线人数</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.users }}</div>
            <div class="stat-label">普通用户</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.counselors }}</div>
            <div class="stat-label">咨询师</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.admins }}</div>
            <div class="stat-label">管理员</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索栏 -->
      <el-form :model="queryForm" inline class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="queryForm.keyword" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="用户类型">
          <el-select v-model="queryForm.userType" placeholder="请选择用户类型" clearable>
            <el-option label="普通用户" value="user" />
            <el-option label="咨询师" value="counselor" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="handleRefresh">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="id" label="用户ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="nickname" label="昵称" width="120" />
        <el-table-column prop="user_type" label="用户类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getUserTypeTag(row.user_type)">
              {{ getUserTypeLabel(row.user_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip_address" label="IP地址" width="140" />
        <el-table-column prop="login_time" label="登录时间" width="180" />
        <el-table-column prop="last_active_time" label="最后活跃" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleSendTo(row)">
              发送消息
            </el-button>
            <el-button type="warning" link size="small" @click="handleKickOut(row)">
              强制下线
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="queryForm.page"
        v-model:page-size="queryForm.pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleQuery"
        @current-change="handleQuery"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <!-- 发送消息对话框 -->
    <el-dialog v-model="messageVisible" title="发送消息" width="500px">
      <el-form :model="messageForm" label-width="80px">
        <el-form-item label="接收人">
          <el-input v-model="messageForm.receiver" disabled />
        </el-form-item>
        <el-form-item label="消息内容" required>
          <el-input
            v-model="messageForm.content"
            type="textarea"
            :rows="4"
            placeholder="请输入消息内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="messageVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSendMessage" :loading="sending">发送</el-button>
      </template>
    </el-dialog>

    <!-- 广播消息对话框 -->
    <el-dialog v-model="broadcastVisible" title="广播消息" width="500px">
      <el-form :model="broadcastForm" label-width="80px">
        <el-form-item label="消息内容" required>
          <el-input
            v-model="broadcastForm.content"
            type="textarea"
            :rows="4"
            placeholder="请输入广播消息内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="broadcastVisible = false">取消</el-button>
        <el-button type="primary" @click="handleBroadcastSend" :loading="broadcasting">广播</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Promotion } from '@element-plus/icons-vue'
import { getOnlineUsers, broadcastMessage } from '@/api/system'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const messageVisible = ref(false)
const broadcastVisible = ref(false)
const sending = ref(false)
const broadcasting = ref(false)

const stats = reactive({
  total: 0,
  users: 0,
  counselors: 0,
  admins: 0
})

const queryForm = reactive({
  page: 1,
  pageSize: 20,
  keyword: '',
  userType: ''
})

const messageForm = reactive({
  userId: '',
  receiver: '',
  content: ''
})

const broadcastForm = reactive({
  content: ''
})

const handleQuery = async () => {
  loading.value = true
  try {
    const res = await getOnlineUsers(queryForm)
    tableData.value = res.data?.users || res.users || []
    total.value = res.data?.total || res.total || 0

    // 统计各类型用户
    stats.total = tableData.value.length
    stats.users = tableData.value.filter(u => !u.is_admin).length
    stats.counselors = tableData.value.filter(u => u.is_counselor).length
    stats.admins = tableData.value.filter(u => u.is_admin).length
  } catch (error) {
    ElMessage.error(error.message || '获取在线用户失败')
  } finally {
    loading.value = false
  }
}

const handleRefresh = () => {
  handleQuery()
}

const handleSendTo = (row) => {
  messageForm.userId = row.id
  messageForm.receiver = row.username
  messageForm.content = ''
  messageVisible.value = true
}

const handleSendMessage = async () => {
  if (!messageForm.content) {
    ElMessage.warning('请输入消息内容')
    return
  }

  sending.value = true
  try {
    // TODO: 调用发送消息接口
    ElMessage.success('消息发送成功')
    messageVisible.value = false
  } catch (error) {
    ElMessage.error(error.message || '发送失败')
  } finally {
    sending.value = false
  }
}

const handleKickOut = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要强制下线用户 ${row.username} 吗？`, '提示', {
      type: 'warning'
    })

    // TODO: 调用强制下线接口
    ElMessage.success('用户已下线')
    handleQuery()
  } catch {
    // 用户取消
  }
}

const handleBroadcast = () => {
  broadcastForm.content = ''
  broadcastVisible.value = true
}

const handleBroadcastSend = async () => {
  if (!broadcastForm.content) {
    ElMessage.warning('请输入广播消息内容')
    return
  }

  broadcasting.value = true
  try {
    await broadcastMessage({ content: broadcastForm.content })
    ElMessage.success('广播消息已发送')
    broadcastVisible.value = false
  } catch (error) {
    ElMessage.error(error.message || '广播失败')
  } finally {
    broadcasting.value = false
  }
}

const getUserTypeLabel = (type) => {
  if (type === 'admin') return '管理员'
  if (type === 'counselor') return '咨询师'
  return '普通用户'
}

const getUserTypeTag = (type) => {
  if (type === 'admin') return 'danger'
  if (type === 'counselor') return 'warning'
  return 'primary'
}

onMounted(() => {
  handleQuery()
  // 每30秒自动刷新
  setInterval(handleQuery, 30000)
})
</script>

<style lang="scss" scoped>
.online-container {
  padding: 20px;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: bold;
    font-size: 16px;
  }

  .stats-row {
    margin-bottom: 20px;

    .stat-card {
      text-align: center;

      .stat-value {
        font-size: 36px;
        font-weight: bold;
        color: #409eff;
        margin-bottom: 8px;
      }

      .stat-label {
        font-size: 14px;
        color: #606266;
      }
    }
  }

  .search-form {
    margin-bottom: 20px;
  }
}
</style>
