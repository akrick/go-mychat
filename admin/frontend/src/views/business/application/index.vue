<template>
  <div class="application-container">
    <el-card>
      <el-form :inline="true" :model="queryParams" class="demo-form-inline">
        <el-form-item label="搜索">
          <el-input v-model="queryParams.keyword" placeholder="姓名/手机号/邮箱" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="全部" clearable>
            <el-option label="待审核" :value="0" />
            <el-option label="审核通过" :value="1" />
            <el-option label="审核拒绝" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" border style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="申请人" width="120">
          <template #default="{ row }">
            {{ row.name }}
          </template>
        </el-table-column>
        <el-table-column prop="gender" label="性别" width="80" />
        <el-table-column prop="phone" label="联系电话" width="130" />
        <el-table-column prop="email" label="邮箱" width="180" />
        <el-table-column prop="title" label="职称" width="140" />
        <el-table-column prop="years_exp" label="从业年限" width="100">
          <template #default="{ row }">
            {{ row.years_exp }}年
          </template>
        </el-table-column>
        <el-table-column label="擅长领域" width="200">
          <template #default="{ row }">
            <el-tag
              v-for="(tag, index) in getSpecialtyTags(row.specialty)"
              :key="index"
              size="small"
              style="margin-right: 5px; margin-bottom: 5px"
            >
              {{ tag }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="180" />
        <el-table-column prop="reviewed_at" label="审核时间" width="180">
          <template #default="{ row }">
            {{ row.reviewed_at || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">查看</el-button>
            <el-button
              v-if="row.status === 0"
              size="small"
              type="primary"
              @click="handleReview(row)"
            >
              审核
            </el-button>
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

    <!-- 查看详情对话框 -->
    <el-dialog v-model="viewDialogVisible" title="申请详情" width="800px" class="detail-dialog">
      <el-descriptions v-if="currentApplication" :column="2" border>
        <el-descriptions-item label="申请人">{{ currentApplication.name }}</el-descriptions-item>
        <el-descriptions-item label="性别">{{ currentApplication.gender }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ currentApplication.phone }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ currentApplication.email }}</el-descriptions-item>
        <el-descriptions-item label="职称">{{ currentApplication.title }}</el-descriptions-item>
        <el-descriptions-item label="从业年限">{{ currentApplication.years_exp }}年</el-descriptions-item>
        <el-descriptions-item label="擅长领域" :span="2">
          <el-tag
            v-for="(tag, index) in getSpecialtyTags(currentApplication.specialty)"
            :key="index"
            style="margin-right: 5px; margin-bottom: 5px"
          >
            {{ tag }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="个人简介" :span="2">
          {{ currentApplication.bio }}
        </el-descriptions-item>
        <el-descriptions-item label="状态" :span="2">
          <el-tag :type="getStatusType(currentApplication.status)">
            {{ getStatusText(currentApplication.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="申请时间" :span="2">
          {{ currentApplication.created_at }}
        </el-descriptions-item>
        <el-descriptions-item v-if="currentApplication.reviewed_at" label="审核时间" :span="2">
          {{ currentApplication.reviewed_at }}
        </el-descriptions-item>
        <el-descriptions-item v-if="currentApplication.reviewer" label="审核人" :span="2">
          {{ currentApplication.reviewer.username }}
        </el-descriptions-item>
        <el-descriptions-item v-if="currentApplication.reject_reason" label="拒绝原因" :span="2">
          {{ currentApplication.reject_reason }}
        </el-descriptions-item>
      </el-descriptions>

      <el-divider>资质证书</el-divider>
      <div class="certificate-list">
        <el-image
          v-for="(cert, index) in getCertificates(currentApplication)"
          :key="index"
          :src="cert"
          :preview-src-list="getCertificates(currentApplication)"
          style="width: 200px; height: 280px; margin-right: 15px; margin-bottom: 15px"
          fit="cover"
        >
          <template #error>
            <div class="image-error">
              <el-icon><Picture /></el-icon>
              <span>暂无图片</span>
            </div>
          </template>
        </el-image>
        <div v-if="!hasCertificates(currentApplication)" class="no-certificates">
          暂无资质证书
        </div>
      </div>

      <template #footer>
        <el-button @click="viewDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 审核对话框 -->
    <el-dialog v-model="reviewDialogVisible" title="审核入驻申请" width="600px">
      <el-descriptions v-if="currentApplication" :column="2" border>
        <el-descriptions-item label="申请人">{{ currentApplication.name }}</el-descriptions-item>
        <el-descriptions-item label="职称">{{ currentApplication.title }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ currentApplication.phone }}</el-descriptions-item>
        <el-descriptions-item label="从业年限">{{ currentApplication.years_exp }}年</el-descriptions-item>
        <el-descriptions-item label="擅长领域" :span="2">
          {{ currentApplication.specialty }}
        </el-descriptions-item>
      </el-descriptions>

      <el-form :model="reviewForm" ref="reviewFormRef" label-width="100px" style="margin-top: 20px">
        <el-form-item label="审核结果" required>
          <el-radio-group v-model="reviewForm.status">
            <el-radio :value="1">
              <el-text type="success">通过</el-text>
            </el-radio>
            <el-radio :value="2">
              <el-text type="danger">拒绝</el-text>
            </el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item
          v-if="reviewForm.status === 2"
          label="拒绝原因"
          required
        >
          <el-input
            v-model="reviewForm.reject_reason"
            type="textarea"
            :rows="4"
            placeholder="请填写拒绝原因"
          />
        </el-form-item>
      </el-form>

      <el-alert
        v-if="reviewForm.status === 1"
        type="success"
        :closable="false"
        style="margin-bottom: 15px"
      >
        审核通过后，系统将自动创建咨询师账户，用户即可开始接单。
      </el-alert>

      <template #footer>
        <el-button @click="reviewDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleReviewSubmit" :loading="reviewing">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Picture } from '@element-plus/icons-vue'
import { getApplicationList, reviewApplication } from '@/api/adminApplication'

const loading = ref(false)
const reviewing = ref(false)
const tableData = ref([])
const total = ref(0)

const queryParams = reactive({
  page: 1,
  page_size: 20,
  status: '',
  keyword: ''
})

const viewDialogVisible = ref(false)
const reviewDialogVisible = ref(false)
const currentApplication = ref(null)
const reviewFormRef = ref(null)

const reviewForm = reactive({
  status: 1,
  reject_reason: ''
})

onMounted(() => {
  handleQuery()
})

const handleQuery = () => {
  loading.value = true
  const params = { ...queryParams }
  if (!params.status) {
    delete params.status
  }
  if (!params.keyword) {
    delete params.keyword
  }
  getApplicationList(params)
    .then(res => {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    })
    .catch(err => {
      ElMessage.error(err.msg || '查询失败')
      tableData.value = []
      total.value = 0
    })
    .finally(() => {
      loading.value = false
    })
}

const resetQuery = () => {
  queryParams.keyword = ''
  queryParams.status = ''
  queryParams.page = 1
  handleQuery()
}

const getSpecialtyTags = (specialty) => {
  if (!specialty) return []
  return specialty.split(/[,，]/).map(t => t.trim()).filter(t => t)
}

const getStatusText = (status) => {
  const map = {
    0: '待审核',
    1: '审核通过',
    2: '审核拒绝'
  }
  return map[status] || '未知'
}

const getStatusType = (status) => {
  const map = {
    0: 'warning',
    1: 'success',
    2: 'danger'
  }
  return map[status] || 'info'
}

const getCertificates = (application) => {
  const certs = []
  if (application.certificate_img1) certs.push(application.certificate_img1)
  if (application.certificate_img2) certs.push(application.certificate_img2)
  if (application.certificate_img3) certs.push(application.certificate_img3)
  return certs
}

const hasCertificates = (application) => {
  return application.certificate_img1 || application.certificate_img2 || application.certificate_img3
}

const handleView = (row) => {
  currentApplication.value = row
  viewDialogVisible.value = true
}

const handleReview = (row) => {
  currentApplication.value = row
  reviewForm.status = 1
  reviewForm.reject_reason = ''
  reviewDialogVisible.value = true
}

const handleReviewSubmit = () => {
  if (reviewForm.status === 2 && !reviewForm.reject_reason.trim()) {
    ElMessage.warning('请填写拒绝原因')
    return
  }

  ElMessageBox.confirm(
    reviewForm.status === 1
      ? '确认通过该入驻申请？审核通过后将自动创建咨询师账户。'
      : '确认拒绝该入驻申请？',
    '确认操作',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    reviewing.value = true
    reviewApplication(currentApplication.value.id, reviewForm)
      .then(() => {
        ElMessage.success('审核成功')
        reviewDialogVisible.value = false
        handleQuery()
      })
      .catch(err => {
        ElMessage.error(err.msg || '审核失败')
      })
      .finally(() => {
        reviewing.value = false
      })
  }).catch(() => {})
}
</script>

<style scoped>
.application-container {
  padding: 20px;
}

.demo-form-inline {
  margin-bottom: 20px;
}

.certificate-list {
  display: flex;
  flex-wrap: wrap;
  padding: 10px 0;
}

.no-certificates {
  width: 100%;
  padding: 40px;
  text-align: center;
  color: #999;
  background: #f5f7fa;
  border-radius: 4px;
}

.image-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
}

.image-error .el-icon {
  font-size: 48px;
  margin-bottom: 8px;
}

.detail-dialog :deep(.el-descriptions__label) {
  font-weight: 600;
}
</style>
