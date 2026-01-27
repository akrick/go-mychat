<template>
  <div class="roles-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增角色
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :model="queryForm" inline class="search-form">
        <el-form-item label="角色名称">
          <el-input v-model="queryForm.name" placeholder="请输入角色名称" clearable @clear="handleQuery" @keyup.enter="handleQuery" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.status" placeholder="请选择状态" clearable @clear="handleQuery" @change="handleQuery">
            <el-option label="全部" :value="-1" />
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" stripe v-loading="loading" border>
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="code" label="角色代码" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="280" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleViewUsers(row)">
              查看用户
            </el-button>
            <el-button type="primary" link size="small" @click="handleAssignPermissions(row)">
              分配权限
            </el-button>
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              删除
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
        @current-change="handleQuery"
        @size-change="handleQuery"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色代码" prop="code">
          <el-input v-model="formData.code" placeholder="请输入角色代码，如：admin、user" :disabled="isEdit">
            <template #append>
              <el-button v-if="!isEdit" @click="generateCode">自动生成</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="角色描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入角色描述" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="999" controls-position="right" />
          <span class="form-tip">数值越小排序越靠前</span>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 分配权限对话框 -->
    <el-dialog
      v-model="permissionDialogVisible"
      title="分配权限"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-alert
        title="权限选择"
        type="info"
        :closable="false"
        style="margin-bottom: 16px"
      >
        勾选权限后，拥有该角色的用户将获得对应的操作权限
      </el-alert>

      <el-form label-width="80px">
        <el-form-item label="权限类型">
          <el-radio-group v-model="permissionFilter" @change="handlePermissionFilterChange">
            <el-radio-button label="">全部</el-radio-button>
            <el-radio-button label="menu">菜单</el-radio-button>
            <el-radio-button label="button">按钮</el-radio-button>
            <el-radio-button label="api">接口</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <el-scrollbar max-height="400px">
        <el-tree
          ref="permissionTreeRef"
          :data="filteredPermissionTree"
          :props="treeProps"
          node-key="id"
          show-checkbox
          default-expand-all
          :check-strictly="false"
        >
          <template #default="{ node, data }">
            <span class="custom-tree-node">
              <el-tag v-if="data.type === 'menu'" type="success" size="small" style="margin-right: 8px">菜单</el-tag>
              <el-tag v-else-if="data.type === 'button'" type="warning" size="small" style="margin-right: 8px">按钮</el-tag>
              <el-tag v-else-if="data.type === 'api'" type="info" size="small" style="margin-right: 8px">接口</el-tag>
              <span>{{ node.label }}</span>
            </span>
          </template>
        </el-tree>
      </el-scrollbar>

      <template #footer>
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAssignSubmit" :loading="assignLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 查看用户对话框 -->
    <el-dialog
      v-model="userDialogVisible"
      title="角色用户列表"
      width="800px"
    >
      <el-table :data="roleUsers" v-loading="userLoading" border>
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="id" label="用户ID" width="80" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="nickname" label="昵称" width="150" />
        <el-table-column prop="email" label="邮箱" min-width="200" />
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="userPage"
        :page-size="10"
        :total="userTotal"
        layout="total, prev, pager, next"
        @current-change="loadRoleUsers"
        style="margin-top: 16px; justify-content: flex-end"
      />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh } from '@element-plus/icons-vue'
import { getRoleList, createRole, updateRole, deleteRole, getRolePermissions, assignPermissions } from '@/api/role'
import { getPermissionTree } from '@/api/permission'

const loading = ref(false)
const submitLoading = ref(false)
const assignLoading = ref(false)
const dialogVisible = ref(false)
const permissionDialogVisible = ref(false)
const userDialogVisible = ref(false)
const userLoading = ref(false)
const isEdit = ref(false)
const dialogTitle = ref('')
const formRef = ref(null)
const permissionTreeRef = ref(null)
const tableData = ref([])
const total = ref(0)
const permissionTree = ref([])
const roleUsers = ref([])
const userTotal = ref(0)
const userPage = ref(1)
const permissionFilter = ref('')

const treeProps = {
  children: 'children',
  label: 'name'
}

const queryForm = reactive({
  page: 1,
  pageSize: 20,
  name: '',
  status: -1
})

const formData = reactive({
  id: null,
  name: '',
  code: '',
  description: '',
  sort: 0,
  status: 1
})

const formRules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 50, message: '角色名称长度为2-50个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入角色代码', trigger: 'blur' },
    { pattern: /^[a-z][a-z0-9_]*$/, message: '角色代码只能包含小写字母、数字和下划线，且以字母开头', trigger: 'blur' }
  ]
}

const filteredPermissionTree = computed(() => {
  if (!permissionFilter.value) {
    return permissionTree.value
  }
  return filterPermissionByType(permissionTree.value, permissionFilter.value)
})

const filterPermissionByType = (tree, type) => {
  return tree.map(node => {
    const newNode = { ...node }
    if (node.children && node.children.length > 0) {
      newNode.children = filterPermissionByType(node.children, type)
    }
    // 保留匹配类型或其子节点有匹配类型的节点
    if (node.type === type || (newNode.children && newNode.children.some(child => child.type === type))) {
      return newNode
    }
    // 如果不匹配且有子节点，可能需要保留作为父节点
    if (newNode.children && newNode.children.length > 0) {
      return newNode
    }
    return null
  }).filter(node => node !== null)
}

const loadTableData = async () => {
  try {
    loading.value = true
    const params = { ...queryForm }
    if (params.status === -1) {
      delete params.status
    }
    const res = await getRoleList(params)
    tableData.value = res.data?.list || res.list || []
    total.value = res.data?.total || res.total || 0
  } catch (error) {
    console.error(error)
    ElMessage.error(error.message || '获取角色列表失败')
  } finally {
    loading.value = false
  }
}

const loadPermissionTree = async () => {
  try {
    const res = await getPermissionTree()
    permissionTree.value = res.data || res || []
  } catch (error) {
    console.error(error)
    ElMessage.error(error.message || '获取权限树失败')
  }
}

const handleQuery = () => {
  queryForm.page = 1
  loadTableData()
}

const handleReset = () => {
  Object.assign(queryForm, {
    page: 1,
    pageSize: 20,
    name: '',
    status: -1
  })
  loadTableData()
}

const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = '新增角色'
  Object.assign(formData, {
    id: null,
    name: '',
    code: '',
    description: '',
    sort: 0,
    status: 1
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑角色'
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    code: row.code,
    description: row.description,
    sort: row.sort,
    status: row.status
  })
  dialogVisible.value = true
}

const generateCode = () => {
  // 根据角色名称生成拼音或英文代码
  const name = formData.name.trim()
  if (!name) {
    ElMessage.warning('请先输入角色名称')
    return
  }
  // 简单生成：取首字母小写 + 时间戳后6位
  const timestamp = Date.now().toString().slice(-6)
  formData.code = `role_${timestamp}`
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitLoading.value = true
    if (isEdit.value) {
      await updateRole(formData.id, formData)
      ElMessage.success('更新成功')
    } else {
      await createRole(formData)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadTableData()
  } catch (error) {
    if (error.message) {
      ElMessage.error(error.message)
    }
  } finally {
    submitLoading.value = false
  }
}

const handleAssignPermissions = async (row) => {
  currentRoleId.value = row.id
  permissionDialogVisible.value = true
  permissionFilter.value = ''

  try {
    const res = await getRolePermissions(row.id)
    const permissions = res.data || res || []
    const permissionIds = permissions.map(p => p.id)
    // 使用nextTick确保树已渲染
    await nextTick()
    if (permissionTreeRef.value) {
      permissionTreeRef.value.setCheckedKeys(permissionIds, false)
    }
  } catch (error) {
    console.error(error)
    ElMessage.error(error.message || '获取角色权限失败')
  }
}

const currentRoleId = ref(null)

const handlePermissionFilterChange = () => {
  // 权限筛选变化时，保持已勾选状态
  const checkedKeys = permissionTreeRef.value?.getCheckedKeys() || []
}

const handleAssignSubmit = async () => {
  try {
    assignLoading.value = true
    // 获取所有选中的权限ID（包括半选的父节点）
    const checkedKeys = permissionTreeRef.value.getCheckedKeys()
    const halfCheckedKeys = permissionTreeRef.value.getHalfCheckedKeys()
    const allPermissionIds = [...checkedKeys, ...halfCheckedKeys]

    await assignPermissions(currentRoleId.value, {
      permission_ids: allPermissionIds
    })
    ElMessage.success('分配成功')
    permissionDialogVisible.value = false
  } catch (error) {
    ElMessage.error(error.message || '分配权限失败')
  } finally {
    assignLoading.value = false
  }
}

const handleViewUsers = async (row) => {
  userDialogVisible.value = true
  userPage.value = 1
  currentRoleId.value = row.id
  await loadRoleUsers()
}

const loadRoleUsers = async () => {
  try {
    userLoading.value = true
    // TODO: 调用获取角色用户列表接口
    // const res = await getRoleUsers(currentRoleId.value, userPage.value, 10)
    // roleUsers.value = res.list || []
    // userTotal.value = res.total || 0
    roleUsers.value = []
    userTotal.value = 0
  } catch (error) {
    ElMessage.error(error.message || '获取角色用户失败')
  } finally {
    userLoading.value = false
  }
}

const handleDelete = async (row) => {
  // 检查是否有用户使用该角色
  if (row.user_count > 0) {
    ElMessage.warning(`该角色下有 ${row.user_count} 个用户，无法删除`)
    return
  }

  try {
    await ElMessageBox.confirm(`确定删除角色 ${row.name} 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteRole(row.id)
    ElMessage.success('删除成功')
    loadTableData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error(error.message || '删除失败')
    }
  }
}

onMounted(() => {
  loadTableData()
  loadPermissionTree()
})
</script>

<style lang="scss" scoped>
.roles-page {
  padding: 20px;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: bold;
    font-size: 16px;
  }

  .search-form {
    margin-bottom: 20px;
  }

  .form-tip {
    margin-left: 10px;
    font-size: 13px;
    color: #909399;
  }

  .custom-tree-node {
    display: flex;
    align-items: center;
    flex: 1;
  }
}
</style>
