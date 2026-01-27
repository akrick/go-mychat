<template>
  <div class="withdraw-container">
    <el-card>
      <el-form :inline="true" class="demo-form-inline">
        <el-form-item>
          <el-button type="primary" @click="handleQuery">刷新</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" border style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="咨询师" width="150">
          <template #default="{ row }">
            {{ row.counselor?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="提现金额(元)" width="120">
          <template #default="{ row }">
            <span style="color: #f56c6c; font-weight: bold;">¥{{ row.amount.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="bank_name" label="开户行" width="180" />
        <el-table-column prop="bank_account" label="银行账号" width="180" />
        <el-table-column prop="account_name" label="账户名" width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="success" @click="handleApprove(row)">通过</el-button>
            <el-button size="small" type="danger" @click="handleReject(row)">拒绝</el-button>
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

    <!-- 拒绝对话框 -->
    <el-dialog v-model="rejectDialogVisible" title="拒绝提现申请" width="500px">
      <el-form :model="rejectForm" ref="rejectFormRef" label-width="100px">
        <el-form-item label="提现金额">
          <span>¥{{ currentWithdraw?.amount?.toFixed(2) }}</span>
        </el-form-item>
        <el-form-item label="拒绝原因" prop="rejected_reason">
          <el-input v-model="rejectForm.rejected_reason" type="textarea" :rows="4" placeholder="请输入拒绝原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="handleRejectSubmit">确认拒绝</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getPendingWithdraws, approveWithdraw } from '@/api/finance'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const rejectDialogVisible = ref(false)
const rejectFormRef = ref()
const currentWithdraw = ref(null)

const queryParams = reactive({
  page: 1,
  page_size: 20
})

const rejectForm = reactive({
  approved: false,
  rejected_reason: ''
})

const handleQuery = async () => {
  loading.value = true
  try {
    const res = await getPendingWithdraws(queryParams)
    tableData.value = res.withdraws
    total.value = res.total
  } catch (error) {
    ElMessage.error(error.message || '获取提现列表失败')
  } finally {
    loading.value = false
  }
}

const getStatusType = (status) => {
  const types = {
    0: 'warning',
    1: 'success',
    2: 'danger',
    3: 'success'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    0: '待审核',
    1: '已通过',
    2: '已拒绝',
    3: '已打款'
  }
  return texts[status] || '未知'
}

const handleApprove = (row) => {
  ElMessageBox.confirm(`确认通过该提现申请?金额: ¥${row.amount.toFixed(2)}`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await approveWithdraw(row.id, { approved: true })
      ElMessage.success('审核通过')
      handleQuery()
    } catch (error) {
      ElMessage.error(error.message || '审核失败')
    }
  })
}

const handleReject = (row) => {
  currentWithdraw.value = row
  rejectForm.rejected_reason = ''
  rejectDialogVisible.value = true
}

const handleRejectSubmit = async () => {
  if (!rejectForm.rejected_reason) {
    ElMessage.warning('请输入拒绝原因')
    return
  }

  try {
    await approveWithdraw(currentWithdraw.value.id, rejectForm)
    ElMessage.success('已拒绝该申请')
    rejectDialogVisible.value = false
    handleQuery()
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  }
}

onMounted(() => {
  handleQuery()
})
</script>

<style scoped>
.withdraw-container {
  padding: 20px;
}

.demo-form-inline {
  margin-bottom: 20px;
}

.el-pagination {
  margin-top: 20px;
}
</style>
