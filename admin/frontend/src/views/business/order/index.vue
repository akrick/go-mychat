<template>
  <div class="order-container">
    <el-card>
      <el-form :inline="true" :model="queryParams" class="demo-form-inline">
        <el-form-item label="搜索">
          <el-input v-model="queryParams.keyword" placeholder="订单号/备注" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="全部" clearable>
            <el-option label="待支付" :value="0" />
            <el-option label="已支付" :value="1" />
            <el-option label="已完成" :value="2" />
            <el-option label="已取消" :value="3" />
            <el-option label="已退款" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" border style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="order_no" label="订单号" width="160" />
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
        <el-table-column prop="duration" label="时长(分钟)" width="100" />
        <el-table-column prop="amount" label="金额(元)" width="100">
          <template #default="{ row }">
            ¥{{ row.amount.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="schedule_time" label="预约时间" width="180" />
        <el-table-column prop="pay_time" label="支付时间" width="180" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleUpdateStatus(row)">更新状态</el-button>
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

    <!-- 更新状态对话框 -->
    <el-dialog v-model="statusDialogVisible" title="更新订单状态" width="500px">
      <el-form :model="statusForm" ref="statusFormRef" label-width="100px">
        <el-form-item label="订单号">
          <span>{{ currentOrder?.order_no }}</span>
        </el-form-item>
        <el-form-item label="当前状态">
          <el-tag :type="getStatusType(currentOrder?.status)">
            {{ getStatusText(currentOrder?.status) }}
          </el-tag>
        </el-form-item>
        <el-form-item label="新状态">
          <el-select v-model="statusForm.status">
            <el-option label="待支付" :value="0" />
            <el-option label="已支付" :value="1" />
            <el-option label="已完成" :value="2" />
            <el-option label="已取消" :value="3" />
            <el-option label="已退款" :value="4" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="statusDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleStatusSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getOrderList, updateOrderStatus } from '@/api/adminOrder'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const statusDialogVisible = ref(false)
const statusFormRef = ref()
const currentOrder = ref(null)

const queryParams = reactive({
  page: 1,
  page_size: 20,
  keyword: '',
  status: ''
})

const statusForm = reactive({
  status: 1
})

const handleQuery = async () => {
  loading.value = true
  try {
    const res = await getOrderList(queryParams)
    tableData.value = res.orders
    total.value = res.total
  } catch (error) {
    ElMessage.error(error.message || '获取订单列表失败')
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

const getStatusType = (status) => {
  const types = {
    0: 'info',
    1: 'warning',
    2: 'success',
    3: 'danger',
    4: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    0: '待支付',
    1: '已支付',
    2: '已完成',
    3: '已取消',
    4: '已退款'
  }
  return texts[status] || '未知'
}

const handleUpdateStatus = (row) => {
  currentOrder.value = row
  statusForm.status = row.status
  statusDialogVisible.value = true
}

const handleStatusSubmit = async () => {
  try {
    await updateOrderStatus(currentOrder.value.id, statusForm)
    ElMessage.success('更新成功')
    statusDialogVisible.value = false
    handleQuery()
  } catch (error) {
    ElMessage.error(error.message || '更新失败')
  }
}

onMounted(() => {
  handleQuery()
})
</script>

<style scoped>
.order-container {
  padding: 20px;
}

.demo-form-inline {
  margin-bottom: 20px;
}

.el-pagination {
  margin-top: 20px;
}
</style>
