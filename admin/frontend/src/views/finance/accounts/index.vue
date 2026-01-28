<template>
  <div class="accounts-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>咨询师账户管理</span>
          <el-button type="primary" @click="handleQuery">刷新</el-button>
        </div>
      </template>

      <el-table :data="tableData" border style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="咨询师" width="150">
          <template #default="{ row }">
            {{ row.counselor?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="total_income" label="总收入(元)" width="120">
          <template #default="{ row }">
            <span style="color: #67c23a; font-weight: bold;">¥{{ row.total_income.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="可用余额(元)" width="120">
          <template #default="{ row }">
            <span style="color: #409eff; font-weight: bold;">¥{{ row.balance.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="withdrawn" label="已提现(元)" width="120">
          <template #default="{ row }">
            <span>¥{{ row.withdrawn.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="frozen_amount" label="冻结金额(元)" width="120">
          <template #default="{ row }">
            <span style="color: #e6a23c;">¥{{ row.frozen_amount.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="handleViewDetail(row)">详情</el-button>
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="账户详情" width="800px">
      <el-descriptions v-if="currentAccount" :column="2" border>
        <el-descriptions-item label="咨询师">{{ currentAccount.counselor?.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="总收入">¥{{ currentAccount.total_income?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="可用余额">¥{{ currentAccount.balance?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="已提现">¥{{ currentAccount.withdrawn?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="冻结金额">¥{{ currentAccount.frozen_amount?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="提现次数">{{ accountDetail?.withdraw_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="提现总额">¥{{ accountDetail?.withdraw_total?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentAccount.created_at }}</el-descriptions-item>
      </el-descriptions>

      <el-divider>待审核提现</el-divider>

      <el-table v-if="accountDetail?.pending_withdraws?.length" :data="accountDetail.pending_withdraws" border max-height="300">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            ¥{{ row.amount.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="bank_name" label="开户行" width="180" />
        <el-table-column prop="bank_account" label="银行账号" width="180" />
        <el-table-column prop="created_at" label="申请时间" width="180" />
      </el-table>
      <el-empty v-else description="暂无待审核提现" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getCounselorAccountList, getCounselorAccountDetail } from '@/api/finance'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const detailDialogVisible = ref(false)
const currentAccount = ref(null)
const accountDetail = ref(null)

const queryParams = reactive({
  page: 1,
  page_size: 20
})

const handleQuery = async () => {
  loading.value = true
  try {
    const res = await getCounselorAccountList(queryParams)
    tableData.value = res.data?.accounts || res.accounts || []
    total.value = res.data?.total || res.total || 0
  } catch (error) {
    ElMessage.error(error.message || '获取账户列表失败')
  } finally {
    loading.value = false
  }
}

const handleViewDetail = async (row) => {
  try {
    const res = await getCounselorAccountDetail(row.counselor_id)
    currentAccount.value = row
    accountDetail.value = res.data || res
    detailDialogVisible.value = true
  } catch (error) {
    ElMessage.error(error.message || '获取账户详情失败')
  }
}

onMounted(() => {
  handleQuery()
})
</script>

<style scoped>
.accounts-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
}

.el-pagination {
  margin-top: 20px;
}
</style>
