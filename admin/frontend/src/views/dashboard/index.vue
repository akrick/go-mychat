<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card user-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)">
              <el-icon :size="32" color="#fff"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ formatNumber(statistics.userCount || 0) }}</div>
              <div class="stat-label">总用户数</div>
              <div class="stat-trend" :class="{ 'up': statistics.userTrend >= 0, 'down': statistics.userTrend < 0 }">
                <el-icon><TrendCharts v-if="statistics.userTrend >= 0" /><Bottom v-else /></el-icon>
                <span>{{ Math.abs(statistics.userTrend || 0) }}%</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card counselor-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%)">
              <el-icon :size="32" color="#fff"><Avatar /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ formatNumber(statistics.counselorCount || 0) }}</div>
              <div class="stat-label">咨询师数</div>
              <div class="stat-trend" :class="{ 'up': statistics.counselorTrend >= 0, 'down': statistics.counselorTrend < 0 }">
                <el-icon><TrendCharts v-if="statistics.counselorTrend >= 0" /><Bottom v-else /></el-icon>
                <span>{{ Math.abs(statistics.counselorTrend || 0) }}%</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card order-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)">
              <el-icon :size="32" color="#fff"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ formatNumber(statistics.orderCount || 0) }}</div>
              <div class="stat-label">订单总数</div>
              <div class="stat-trend" :class="{ 'up': statistics.orderTrend >= 0, 'down': statistics.orderTrend < 0 }">
                <el-icon><TrendCharts v-if="statistics.orderTrend >= 0" /><Bottom v-else /></el-icon>
                <span>{{ Math.abs(statistics.orderTrend || 0) }}%</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card revenue-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)">
              <el-icon :size="32" color="#fff"><Money /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">¥{{ formatNumber(statistics.totalRevenue || 0) }}</div>
              <div class="stat-label">总营收</div>
              <div class="stat-trend" :class="{ 'up': statistics.revenueTrend >= 0, 'down': statistics.revenueTrend < 0 }">
                <el-icon><TrendCharts v-if="statistics.revenueTrend >= 0" /><Bottom v-else /></el-icon>
                <span>{{ Math.abs(statistics.revenueTrend || 0) }}%</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="charts-row">
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>订单趋势</span>
              <el-radio-group v-model="orderChartPeriod" size="small" @change="loadOrderChart">
                <el-radio-button label="week">本周</el-radio-button>
                <el-radio-button label="month">本月</el-radio-button>
                <el-radio-button label="year">全年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="orderChartRef" style="height: 350px"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>订单状态</span>
            </div>
          </template>
          <div ref="statusChartRef" style="height: 350px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="charts-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>营收趋势</span>
              <el-radio-group v-model="revenueChartPeriod" size="small" @change="loadRevenueChart">
                <el-radio-button label="week">本周</el-radio-button>
                <el-radio-button label="month">本月</el-radio-button>
                <el-radio-button label="year">全年</el-radio-button>
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
              <span>咨询师排名</span>
              <el-button link type="primary" @click="loadCounselorRanking">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </template>
          <el-table :data="counselorRanking" stripe height="350px">
            <el-table-column type="index" label="排名" width="80" align="center">
              <template #default="{ $index }">
                <el-tag :type="$index < 3 ? 'danger' : $index < 6 ? 'warning' : 'info'" size="small">
                  {{ $index + 1 }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="name" label="咨询师" min-width="120">
              <template #default="{ row }">
                <div style="display: flex; align-items: center; gap: 10px;">
                  <el-avatar :size="32">{{ row.name?.charAt(0) }}</el-avatar>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="orderCount" label="订单数" width="100" align="center" />
            <el-table-column prop="revenue" label="营收" width="120" align="right">
              <template #default="{ row }">
                <span style="color: #67c23a; font-weight: bold;">¥{{ row.revenue }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="rating" label="评分" width="100" align="center">
              <template #default="{ row }">
                <el-rate v-model="row.rating" disabled show-score text-color="#ff9900" />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷操作 -->
    <el-row :gutter="20" class="actions-row">
      <el-col :span="6">
        <el-card class="action-card" @click="goToUserManagement">
          <div class="action-content">
            <el-icon :size="40" color="#667eea"><User /></el-icon>
            <h3>用户管理</h3>
            <p>管理平台用户</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="action-card" @click="goToCounselorManagement">
          <div class="action-content">
            <el-icon :size="40" color="#f5576c"><Avatar /></el-icon>
            <h3>咨询师管理</h3>
            <p>管理咨询师信息</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="action-card" @click="goToOrderManagement">
          <div class="action-content">
            <el-icon :size="40" color="#4facfe"><Document /></el-icon>
            <h3>订单管理</h3>
            <p>查看和处理订单</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="action-card" @click="goToChatManagement">
          <div class="action-content">
            <el-icon :size="40" color="#43e97b"><ChatDotRound /></el-icon>
            <h3>聊天记录</h3>
            <p>查看聊天消息</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import * as echarts from 'echarts'
import { useRouter } from 'vue-router'
import request from '@/utils/request'
import { User, Avatar, Document, Money, ChatDotRound, TrendCharts, Bottom, Refresh } from '@element-plus/icons-vue'

const router = useRouter()
const orderChartRef = ref(null)
const revenueChartRef = ref(null)
const statusChartRef = ref(null)
const statistics = ref({})
const counselorRanking = ref([])
const orderChartPeriod = ref('week')
const revenueChartPeriod = ref('month')

let orderChart = null
let revenueChart = null
let statusChart = null

const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toLocaleString()
}

const loadStatistics = async () => {
  try {
    const res = await request.get('/admin/statistics')
    console.log('统计数据响应:', res)
    // 兼容不同的响应格式
    const data = res.data || res || {}
    // 转换snake_case到camelCase
    statistics.value = {
      userCount: data.user_count || 0,
      userTrend: data.user_trend || 0,
      counselorCount: data.counselor_count || 0,
      counselorTrend: data.counselor_trend || 0,
      orderCount: data.order_count || 0,
      orderTrend: data.order_trend || 0,
      totalRevenue: data.total_amount || 0,
      revenueTrend: data.revenue_trend || 0,
      todayOrderCount: data.today_order_count || 0,
      todayAmount: data.today_amount || 0,
      sessionCount: data.session_count || 0,
      activeSessionCount: data.active_session_count || 0
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const loadCounselorRanking = async () => {
  try {
    const res = await request.get('/admin/stats/counselor/ranking')
    counselorRanking.value = res || []
  } catch (error) {
    console.error('获取咨询师排名失败:', error)
  }
}

const loadOrderChart = async () => {
  try {
    const res = await request.get('/admin/stats/order/trend', { period: orderChartPeriod.value })
    updateOrderChart(res)
  } catch (error) {
    console.error('获取订单趋势失败:', error)
  }
}

const loadRevenueChart = async () => {
  try {
    const res = await request.get('/admin/finance/revenue', { group_by: revenueChartPeriod.value })
    updateRevenueChart(res)
  } catch (error) {
    console.error('获取营收趋势失败:', error)
  }
}

const updateOrderChart = (data) => {
  const dates = data?.dates || ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
  const values = data?.values || [120, 200, 150, 80, 70, 110, 130]

  if (orderChart) {
    const option = {
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        borderColor: '#667eea',
        borderWidth: 1,
        textStyle: { color: '#fff' }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: dates,
        boundaryGap: false,
        axisLabel: { color: '#909399' }
      },
      yAxis: {
        type: 'value',
        axisLabel: { color: '#909399' },
        splitLine: { lineStyle: { color: '#f0f2f5' } }
      },
      series: [{
        data: values,
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        lineStyle: { color: '#667eea', width: 3 },
        itemStyle: { color: '#667eea' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(102, 126, 234, 0.3)' },
            { offset: 1, color: 'rgba(102, 126, 234, 0.05)' }
          ])
        }
      }]
    }
    orderChart.setOption(option)
  }
}

const updateRevenueChart = (data) => {
  const dates = data?.dates || ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
  const values = data?.values || [1200, 2000, 1500, 800, 700, 1100, 1300]

  if (revenueChart) {
    const option = {
      tooltip: {
        trigger: 'axis',
        formatter: '{b}: ¥{c}'
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: dates,
        axisLabel: { color: '#909399' }
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          formatter: '¥{value}',
          color: '#909399'
        },
        splitLine: { lineStyle: { color: '#f0f2f5' } }
      },
      series: [{
        data: values,
        type: 'bar',
        barWidth: '40%',
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#43e97b' },
            { offset: 1, color: '#38f9d7' }
          ]),
          borderRadius: [8, 8, 0, 0]
        }
      }]
    }
    revenueChart.setOption(option)
  }
}

const updateStatusChart = () => {
  if (statusChart) {
    const option = {
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        left: 'left',
        textStyle: { color: '#606266' }
      },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['60%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: true,
          formatter: '{b}'
        },
        labelLine: { show: false },
        data: [
          { value: statistics.value.completedOrders || 1048, name: '已完成', itemStyle: { color: '#67c23a' } },
          { value: statistics.value.pendingOrders || 735, name: '进行中', itemStyle: { color: '#e6a23c' } },
          { value: statistics.value.unpaidOrders || 580, name: '待支付', itemStyle: { color: '#f56c6c' } },
          { value: statistics.value.cancelledOrders || 300, name: '已取消', itemStyle: { color: '#909399' } }
        ],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          },
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold'
          }
        }
      }]
    }
    statusChart.setOption(option)
  }
}

const goToUserManagement = () => {
  router.push('/user')
}

const goToCounselorManagement = () => {
  router.push('/counselor')
}

const goToOrderManagement = () => {
  router.push('/order')
}

const goToChatManagement = () => {
  router.push('/chat')
}

const handleResize = () => {
  orderChart?.resize()
  revenueChart?.resize()
  statusChart?.resize()
}

onMounted(async () => {
  await loadStatistics()
  await loadCounselorRanking()
  await loadOrderChart()
  await loadRevenueChart()

  orderChart = echarts.init(orderChartRef.value)
  await loadOrderChart()

  revenueChart = echarts.init(revenueChartRef.value)
  await loadRevenueChart()

  statusChart = echarts.init(statusChartRef.value)
  updateStatusChart()

  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  orderChart?.dispose()
  revenueChart?.dispose()
  statusChart?.dispose()
  window.removeEventListener('resize', handleResize)
})
</script>

<style lang="scss" scoped>
.dashboard {
  .stat-card {
    margin-bottom: 20px;
    border: none;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    transition: all 0.3s;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    }

    :deep(.el-card__body) {
      padding: 20px;
    }

    .stat-content {
      display: flex;
      align-items: flex-start;

      .stat-icon {
        width: 64px;
        height: 64px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 16px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      }

      .stat-info {
        flex: 1;

        .stat-value {
          font-size: 28px;
          font-weight: bold;
          color: #303133;
          margin-bottom: 8px;
          line-height: 1;
        }

        .stat-label {
          font-size: 14px;
          color: #909399;
          margin-bottom: 8px;
        }

        .stat-trend {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 12px;
          padding: 4px 8px;
          border-radius: 4px;

          &.up {
            background: #f0f9ff;
            color: #67c23a;
          }

          &.down {
            background: #fef0f0;
            color: #f56c6c;
          }
        }
      }
    }
  }

  .charts-row {
    margin-top: 20px;
  }

  .actions-row {
    margin-top: 20px;

    .action-card {
      cursor: pointer;
      transition: all 0.3s;
      border: none;

      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
      }

      .action-content {
        text-align: center;
        padding: 20px;

        .el-icon {
          margin-bottom: 12px;
        }

        h3 {
          margin: 12px 0 8px 0;
          font-size: 16px;
          color: #303133;
        }

        p {
          margin: 0;
          font-size: 13px;
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
}
</style>
