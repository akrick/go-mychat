<template>
  <div class="statistics-container">
    <el-row :gutter="20" style="margin-bottom: 20px;">
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="总用户数" :value="statistics.user_count">
            <template #suffix>
              <el-icon><User /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="在线用户" :value="statistics.online_user_count">
            <template #suffix>
              <el-icon><Connection /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="咨询师数" :value="statistics.counselor_count">
            <template #suffix>
              <el-icon><UserFilled /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="总订单数" :value="statistics.order_count">
            <template #suffix>
              <el-icon><Document /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-bottom: 20px;">
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="今日订单" :value="statistics.today_order_count">
            <template #suffix>
              <el-icon><Tickets /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="总交易额" :value="statistics.total_amount" :precision="2" prefix="¥">
            <template #suffix>
              <el-icon><Wallet /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="今日收入" :value="statistics.today_amount" :precision="2" prefix="¥">
            <template #suffix>
              <el-icon><Money /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="活跃会话" :value="statistics.active_session_count">
            <template #suffix>
              <el-icon><ChatDotRound /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>会话统计</span>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="总会话数">{{ statistics.session_count }}</el-descriptions-item>
            <el-descriptions-item label="活跃会话">{{ statistics.active_session_count }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>在线用户</span>
          </template>
          <el-button type="primary" @click="loadOnlineUsers" style="margin-bottom: 10px;">刷新</el-button>
          <el-table :data="onlineUsers" border max-height="400">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="username" label="用户名" />
            <el-table-column prop="email" label="邮箱" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAdminStatistics, getOnlineUsers } from '@/api/statistics'

const statistics = reactive({
  user_count: 0,
  online_user_count: 0,
  counselor_count: 0,
  order_count: 0,
  today_order_count: 0,
  total_amount: 0,
  today_amount: 0,
  session_count: 0,
  active_session_count: 0
})

const onlineUsers = ref([])

const loadStatistics = async () => {
  try {
    const res = await getAdminStatistics()
    Object.assign(statistics, res.data || res)
  } catch (error) {
    ElMessage.error(error.message || '获取统计数据失败')
  }
}

const loadOnlineUsers = async () => {
  try {
    const res = await getOnlineUsers()
    onlineUsers.value = res.data?.users || res.users || []
  } catch (error) {
    ElMessage.error(error.message || '获取在线用户失败')
  }
}

onMounted(() => {
  loadStatistics()
  loadOnlineUsers()
})
</script>

<style scoped>
.statistics-container {
  padding: 20px;
}
</style>
