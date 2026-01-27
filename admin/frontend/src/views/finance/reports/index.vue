<template>
  <div class="finance-reports-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>财务统计</span>
          <el-button type="primary" @click="handleQuery">刷新数据</el-button>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #409eff">
                <el-icon :size="32" color="#fff"><Money /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">¥{{ formatNumber(stats.total_revenue || 0) }}</div>
                <div class="stat-label">总营收</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #67c23a">
                <el-icon :size="32" color="#fff"><Wallet /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">¥{{ formatNumber(stats.total_counselor_earning || 0) }}</div>
                <div class="stat-label">咨询师收益</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #e6a23c">
                <el-icon :size="32" color="#fff"><Stamp /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">¥{{ formatNumber(stats.total_commission || 0) }}</div>
                <div class="stat-label">平台佣金</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #f56c6c">
                <el-icon :size="32" color="#fff"><WalletFilled /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">¥{{ formatNumber(stats.total_withdrawn || 0) }}</div>
                <div class="stat-label">已提现</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span="12">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>营收趋势</span>
                <el-radio-group v-model="groupBy" @change="handleQuery">
                  <el-radio-button label="day">按天</el-radio-button>
                  <el-radio-button label="month">按月</el-radio-button>
                  <el-radio-button label="year">按年</el-radio-button>
                </el-radio-group>
              </div>
            </template>
            <div ref="revenueChartRef" style="height: 350px"></div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>今日数据</span>
              </div>
            </template>
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="today-stat">
                  <div class="today-label">今日营收</div>
                  <div class="today-value">¥{{ formatNumber(stats.today_revenue || 0) }}</div>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="today-stat">
                  <div class="today-label">今日订单</div>
                  <div class="today-value">{{ stats.today_orders || 0 }}</div>
                </div>
              </el-col>
            </el-row>
            <el-divider />
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="today-stat">
                  <div class="today-label">待审核提现</div>
                  <div class="today-value">{{ stats.pending_withdraws || 0 }}</div>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="today-stat">
                  <div class="today-label">已审核提现</div>
                  <div class="today-value">{{ stats.approved_withdraws || 0 }}</div>
                </div>
              </el-col>
            </el-row>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import * as echarts from 'echarts'
import { Money, Wallet, Stamp, WalletFilled } from '@element-plus/icons-vue'
import { getFinanceStats, getRevenueReport } from '@/api/finance'

const revenueChartRef = ref(null)
const groupBy = ref('day')
const stats = ref({})
let revenueChart = null

const formatNumber = (num) => {
  return parseFloat(num).toFixed(2)
}

const loadStats = async () => {
  try {
    const res = await getFinanceStats()
    stats.value = res || {}
  } catch (error) {
    console.error(error)
  }
}

const loadRevenueReport = async () => {
  try {
    const res = await getRevenueReport({ group_by: groupBy.value })
    const revenueData = res?.revenue_data || []
    
    if (revenueChart) {
      const option = {
        tooltip: { trigger: 'axis' },
        xAxis: {
          type: 'category',
          data: revenueData.map(item => item.date)
        },
        yAxis: { type: 'value' },
        series: [{
          name: '营收',
          data: revenueData.map(item => item.amount),
          type: 'line',
          smooth: true,
          areaStyle: { opacity: 0.3 },
          itemStyle: { color: '#409eff' }
        }]
      }
      revenueChart.setOption(option)
    }
  } catch (error) {
    console.error(error)
  }
}

const handleQuery = async () => {
  await loadStats()
  await loadRevenueReport()
}

const handleResize = () => {
  revenueChart?.resize()
}

onMounted(async () => {
  await handleQuery()
  revenueChart = echarts.init(revenueChartRef.value)
  await loadRevenueReport()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  revenueChart?.dispose()
  window.removeEventListener('resize', handleResize)
})
</script>

<style lang="scss" scoped>
.finance-reports-container {
  padding: 20px;
}

.stat-card {
  margin-bottom: 20px;
  
  :deep(.el-card__body) {
    padding: 20px;
  }
  
  .stat-content {
    display: flex;
    align-items: center;
    
    .stat-icon {
      width: 64px;
      height: 64px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-right: 16px;
    }
    
    .stat-info {
      flex: 1;
      
      .stat-value {
        font-size: 24px;
        font-weight: bold;
        color: #333;
        margin-bottom: 8px;
      }
      
      .stat-label {
        font-size: 14px;
        color: #909399;
      }
    }
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
}

.today-stat {
  text-align: center;
  padding: 20px 0;
  
  .today-label {
    font-size: 14px;
    color: #909399;
    margin-bottom: 8px;
  }
  
  .today-value {
    font-size: 28px;
    font-weight: bold;
    color: #409eff;
  }
}
</style>
