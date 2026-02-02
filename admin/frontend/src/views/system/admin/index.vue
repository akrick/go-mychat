<template>
  <div class="admin-container">
    <el-card>
      <el-form :inline="true" :model="queryParams" class="demo-form-inline">
        <el-form-item label="搜索">
          <el-input v-model="queryParams.keyword" placeholder="用户名/邮箱/手机" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="全部" clearable>
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
          <el-button type="success" @click="handleAdd">新增管理员</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" border style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="real_name" label="真实姓名" width="120" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="phone" label="手机号" width="120" />
        <el-table-column prop="role" label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="row.role === 'super_admin' ? 'danger' : 'primary'">
              {{ row.role === 'super_admin' ? '超级管理员' : '管理员' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="头像" width="80">
          <template #default="{ row }">
            <el-avatar v-if="row.avatar" :src="getAvatarUrl(row.avatar)" />
            <el-avatar v-else>A</el-avatar>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="primary" @click="handleResetPassword(row)">重置密码</el-button>
            <el-button
              size="small"
              :type="row.status === 1 ? 'warning' : 'success'"
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
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

    <!-- 新增/编辑管理员对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" :disabled="!!form.id" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!form.id">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="真实姓名" prop="real_name">
          <el-input v-model="form.real_name" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色">
            <el-option label="管理员" value="admin" />
            <el-option label="超级管理员" value="super_admin" />
          </el-select>
        </el-form-item>
        <el-form-item label="头像" prop="avatar">
          <ImageUpload v-model="form.avatar" tips="支持 JPG、PNG 格式，大小不超过 2MB" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">正常</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 重置密码对话框 -->
    <el-dialog v-model="passwordDialogVisible" title="重置密码" width="500px">
      <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
        <el-form-item label="新密码" prop="password">
          <el-input v-model="passwordForm.password" type="password" placeholder="请输入新密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handlePasswordSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAdministratorList,
  createAdministrator,
  updateAdministrator,
  deleteAdministrator,
  resetAdministratorPassword,
  toggleAdministratorStatus
} from '@/api/administrator'
import ImageUpload from '@/components/ImageUpload.vue'

// 获取完整的头像URL
const getAvatarUrl = (avatar) => {
  if (!avatar) return ''
  if (avatar.startsWith('http://') || avatar.startsWith('https://')) {
    return avatar
  }
  // 如果是相对路径，添加后端服务器地址
  return import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081' + avatar
}

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const passwordDialogVisible = ref(false)
const formRef = ref()
const passwordFormRef = ref()

const queryParams = reactive({
  page: 1,
  page_size: 20,
  keyword: '',
  status: ''
})

const form = reactive({
  id: null,
  username: '',
  password: '',
  real_name: '',
  email: '',
  phone: '',
  avatar: '',
  role: 'admin',
  status: 1
})

const passwordForm = reactive({
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  email: [{ type: 'email', message: '请输入正确的邮箱', trigger: 'blur' }],
  phone: [{ pattern: /^1\d{10}$/, message: '请输入正确的手机号', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }]
}

const passwordRules = {
  password: [{ required: true, message: '请输入新密码', trigger: 'blur' }]
}

const handleQuery = async () => {
  loading.value = true
  try {
    const res = await getAdministratorList(queryParams)
    tableData.value = res.data?.administrators || []
    total.value = res.data?.total || 0
  } catch (error) {
    ElMessage.error(error.message || '获取管理员列表失败')
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

const handleAdd = () => {
  dialogTitle.value = '新增管理员'
  Object.assign(form, {
    id: null,
    username: '',
    password: '',
    real_name: '',
    email: '',
    phone: '',
    avatar: '',
    role: 'admin',
    status: 1
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑管理员'
  Object.assign(form, {
    id: row.id,
    username: row.username,
    password: '',
    real_name: row.real_name,
    email: row.email,
    phone: row.phone,
    avatar: row.avatar,
    role: row.role,
    status: row.status
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    if (form.id) {
      await updateAdministrator(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await createAdministrator(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    handleQuery()
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该管理员吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteAdministrator(row.id)
      ElMessage.success('删除成功')
      handleQuery()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  })
}

const handleResetPassword = (row) => {
  Object.assign(passwordForm, { password: '' })
  passwordForm.currentAdminId = row.id
  passwordDialogVisible.value = true
}

const handlePasswordSubmit = async () => {
  try {
    await passwordFormRef.value.validate()
    await resetAdministratorPassword(passwordForm.currentAdminId, passwordForm)
    ElMessage.success('密码重置成功')
    passwordDialogVisible.value = false
  } catch (error) {
    ElMessage.error(error.message || '重置失败')
  }
}

const handleToggleStatus = (row) => {
  const newStatus = row.status === 1 ? 0 : 1
  const action = newStatus === 1 ? '启用' : '禁用'
  ElMessageBox.confirm(`确定要${action}该管理员吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await toggleAdministratorStatus(row.id, { status: newStatus })
      ElMessage.success(`${action}成功`)
      handleQuery()
    } catch (error) {
      ElMessage.error(error.message || `${action}失败`)
    }
  })
}

onMounted(() => {
  handleQuery()
})
</script>

<style scoped>
.admin-container {
  padding: 20px;
}

.demo-form-inline {
  margin-bottom: 20px;
}

.el-pagination {
  margin-top: 20px;
}
</style>
