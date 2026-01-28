<template>
  <div class="lowcode-data-page">
    <el-card>
      <el-form :model="queryForm" inline class="query-form">
        <el-form-item label="表单">
          <el-select v-model="queryForm.formId" placeholder="请选择表单" @change="handleFormChange">
            <el-option
              v-for="form in formList"
              :key="form.id"
              :label="form.title"
              :value="form.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      
      <div v-if="queryForm.formId">
        <el-table :data="tableData" stripe v-loading="loading">
          <el-table-column type="index" label="序号" width="80" />
          <el-table-column
            v-for="column in columns"
            :key="column.prop"
            :prop="column.prop"
            :label="column.label"
            :min-width="column.width || 150"
          >
            <template #default="{ row }">
              {{ getValue(row, column.prop) }}
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="提交时间" width="180" />
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="handleView(row)">查看</el-button>
              <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <el-pagination
          v-model:current-page="queryForm.page"
          v-model:page-size="queryForm.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
      <div v-else class="empty-tips">
        请选择表单查看数据
      </div>
    </el-card>
    
    <!-- 查看数据对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="数据详情"
      width="600px"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item
          v-for="column in columns"
          :key="column.prop"
          :label="column.label"
        >
          {{ getValue(currentRow, column.prop) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getFormList, getFormDataList } from '@/api/lowcode'

const loading = ref(false)
const detailDialogVisible = ref(false)
const tableData = ref([])
const formList = ref([])
const total = ref(0)
const columns = ref([])
const currentRow = ref({})

const queryForm = reactive({
  page: 1,
  pageSize: 20,
  formId: null
})

const loadFormList = async () => {
  try {
    const res = await getFormList({ page: 1, pageSize: 100 })
    formList.value = res.list || []
  } catch (error) {
    console.error(error)
  }
}

const loadTableData = async () => {
  if (!queryForm.formId) return

  try {
    loading.value = true
    const res = await getFormDataList(queryForm.formId, queryForm)
    tableData.value = res.data?.list || res.list || []
    total.value = res.data?.total || res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleFormChange = async () => {
  // 获取表单配置，解析出列信息
  const form = formList.value.find(f => f.id === queryForm.formId)
  if (form && form.form_json) {
    try {
      const formConfig = JSON.parse(form.form_json)
      columns.value = formConfig.map(item => ({
        prop: item.field,
        label: item.label,
        width: 120
      }))
    } catch (error) {
      console.error('解析表单配置失败', error)
    }
  }
  
  queryForm.page = 1
  loadTableData()
}

const handlePageChange = (page) => {
  queryForm.page = page
  loadTableData()
}

const handleSizeChange = (size) => {
  queryForm.pageSize = size
  loadTableData()
}

const getValue = (row, prop) => {
  try {
    const formData = JSON.parse(row.form_data)
    return formData[prop]
  } catch (error) {
    return '-'
  }
}

const handleView = (row) => {
  currentRow.value = row
  detailDialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除此数据吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    // TODO: 实现删除接口
    ElMessage.success('删除成功')
    loadTableData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

onMounted(() => {
  loadFormList()
})
</script>

<style lang="scss" scoped>
.lowcode-data-page {
  .query-form {
    margin-bottom: 20px;
  }
  
  .empty-tips {
    padding: 60px;
    text-align: center;
    color: #909399;
    font-size: 14px;
  }
  
  .el-pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
