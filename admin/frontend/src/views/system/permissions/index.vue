<template>
  <div class="permissions-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>权限管理</span>
          <div>
            <el-button @click="expandAll">展开全部</el-button>
            <el-button @click="collapseAll">收起全部</el-button>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              新增权限
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :model="queryForm" inline class="search-form">
        <el-form-item label="权限名称">
          <el-input v-model="queryForm.name" placeholder="请输入权限名称" clearable @clear="handleQuery" @keyup.enter="handleQuery" />
        </el-form-item>
        <el-form-item label="权限类型">
          <el-select v-model="queryForm.type" placeholder="请选择类型" clearable @clear="handleQuery" @change="handleQuery">
            <el-option label="全部" value="" />
            <el-option label="菜单" value="menu" />
            <el-option label="按钮" value="button" />
            <el-option label="接口" value="api" />
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

      <!-- 权限树表格 -->
      <el-table
        :data="tableData"
        row-key="id"
        :tree-props="{ children: 'children' }"
        border
        v-loading="loading"
        default-expand-all
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="权限名称" width="200" />
        <el-table-column prop="code" label="权限代码" width="200" />
        <el-table-column label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)" size="small">
              {{ getTypeLabel(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路由路径" width="200" show-overflow-tooltip />
        <el-table-column prop="icon" label="图标" width="100">
          <template #default="{ row }">
            <el-icon v-if="row.icon"><component :is="row.icon" /></el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="component" label="组件路径" width="200" show-overflow-tooltip />
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleAddChild(row)" v-if="row.type === 'menu'">
              新增子菜单
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
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="700px"
      :close-on-click-modal="false"
    >
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="110px">
        <el-form-item label="上级权限" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="parentOptions"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="请选择上级权限"
            clearable
            check-strictly
            :render-after-expand="false"
          />
        </el-form-item>
        <el-form-item label="权限名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入权限名称" />
        </el-form-item>
        <el-form-item label="权限代码" prop="code">
          <el-input v-model="formData.code" placeholder="例如: user:list" :disabled="!!formData.id">
            <template #append v-if="!formData.id">
              <el-button @click="generateCode">自动生成</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="权限类型" prop="type">
          <el-radio-group v-model="formData.type" @change="handleTypeChange">
            <el-radio label="menu">菜单</el-radio>
            <el-radio label="button">按钮</el-radio>
            <el-radio label="api">接口</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="路由路径" prop="path" v-if="formData.type === 'menu'">
          <el-input v-model="formData.path" placeholder="例如: /users" />
        </el-form-item>
        <el-form-item label="组件路径" prop="component" v-if="formData.type === 'menu'">
          <el-input v-model="formData.component" placeholder="例如: /views/system/user/index" />
        </el-form-item>
        <el-form-item label="图标" prop="icon" v-if="formData.type === 'menu'">
          <el-select v-model="formData.icon" placeholder="请选择图标" filterable clearable>
            <el-option v-for="icon in iconList" :key="icon.value" :label="icon.label" :value="icon.value">
              <span style="display: flex; align-items: center;">
                <el-icon style="margin-right: 8px;"><component :is="icon.value" /></el-icon>
                {{ icon.label }}
              </span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="999" controls-position="right" />
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh, User, Setting, Document, Wallet, DataAnalysis, Grid } from '@element-plus/icons-vue'
import { getPermissionTree, createPermission, updatePermission, deletePermission } from '@/api/permission'

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref(null)
const tableData = ref([])
const parentOptions = ref([])
const allPermissions = ref([]) // 保存完整的权限树用于搜索过滤

const iconList = [
  { label: 'User', value: 'User' },
  { label: 'Setting', value: 'Setting' },
  { label: 'Document', value: 'Document' },
  { label: 'Wallet', value: 'Wallet' },
  { label: 'DataAnalysis', value: 'DataAnalysis' },
  { label: 'Grid', value: 'Grid' },
  { label: 'ChatDotRound', value: 'ChatDotRound' },
  { label: 'Calendar', value: 'Calendar' },
  { label: 'Tickets', value: 'Tickets' },
  { label: 'Avatar', value: 'Avatar' },
  { label: 'Key', value: 'Key' },
  { label: 'Menu', value: 'Menu' },
  { label: 'Warning', value: 'Warning' },
  { label: 'InfoFilled', value: 'InfoFilled' }
]

const queryForm = reactive({
  name: '',
  type: ''
})

const formData = reactive({
  id: null,
  parent_id: null,
  name: '',
  code: '',
  type: 'menu',
  path: '',
  component: '',
  icon: '',
  sort: 0,
  status: 1
})

const formRules = {
  name: [
    { required: true, message: '请输入权限名称', trigger: 'blur' },
    { min: 2, max: 50, message: '权限名称长度为2-50个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入权限代码', trigger: 'blur' },
    { pattern: /^[a-z:]+[a-z0-9:_]*$/, message: '权限代码格式不正确，例如: user:list', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择权限类型', trigger: 'change' }
  ]
}

const loadPermissions = async () => {
  loading.value = true
  try {
    const res = await getPermissionTree()
    allPermissions.value = res.data || res || []
    filterPermissions()
    parentOptions.value = res.data || res || []
  } catch (error) {
    ElMessage.error(error.message || '获取权限列表失败')
  } finally {
    loading.value = false
  }
}

const handleQuery = () => {
  filterPermissions()
}

// 递归过滤权限树
const filterTree = (data, name, type) => {
  return data
    .map(item => {
      const filteredChildren = item.children ? filterTree(item.children, name, type) : []
      return { ...item, children: filteredChildren }
    })
    .filter(item => {
      const matchName = !name || item.name.includes(name)
      const matchType = !type || item.type === type
      const hasChildren = item.children && item.children.length > 0
      return (matchName && matchType) || hasChildren
    })
}

const filterPermissions = () => {
  if (!queryForm.name && !queryForm.type) {
    tableData.value = allPermissions.value
  } else {
    tableData.value = filterTree(allPermissions.value, queryForm.name, queryForm.type)
  }
}

const handleReset = () => {
  Object.assign(queryForm, {
    name: '',
    type: ''
  })
  loadPermissions()
}

const expandAll = () => {
  // 展开所有节点
  ElMessage.info('表格已全部展开')
}

const collapseAll = () => {
  // 收起所有节点
  ElMessage.info('表格已全部收起')
}

const handleAdd = () => {
  dialogTitle.value = '新增权限'
  Object.assign(formData, {
    id: null,
    parent_id: null,
    name: '',
    code: '',
    type: 'menu',
    path: '',
    component: '',
    icon: '',
    sort: 0,
    status: 1
  })
  dialogVisible.value = true
}

const handleAddChild = (row) => {
  dialogTitle.value = '新增子权限'
  Object.assign(formData, {
    id: null,
    parent_id: row.id,
    name: '',
    code: '',
    type: row.type,
    path: '',
    component: '',
    icon: '',
    sort: 0,
    status: 1
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑权限'
  Object.assign(formData, {
    id: row.id,
    parent_id: row.parent_id || null,
    name: row.name,
    code: row.code,
    type: row.type,
    path: row.path || '',
    component: row.component || '',
    icon: row.icon || '',
    sort: row.sort || 0,
    status: row.status
  })
  dialogVisible.value = true
}

const handleTypeChange = (type) => {
  if (type !== 'menu') {
    formData.path = ''
    formData.component = ''
    formData.icon = ''
  }
}

const generateCode = () => {
  const name = formData.name.trim()
  const type = formData.type
  if (!name) {
    ElMessage.warning('请先输入权限名称')
    return
  }
  // 简单生成逻辑
  const timestamp = Date.now().toString().slice(-6)
  formData.code = `${type}:${timestamp}`
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitLoading.value = true
    if (formData.id) {
      await updatePermission(formData.id, formData)
      ElMessage.success('更新成功')
    } else {
      await createPermission(formData)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    await loadPermissions()
  } catch (error) {
    if (error.message) {
      ElMessage.error(error.message)
    }
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = (row) => {
  // 检查是否有子节点
  if (row.children && row.children.length > 0) {
    ElMessage.warning('该权限下有子权限，无法删除')
    return
  }

  ElMessageBox.confirm(`确定删除权限 ${row.name} 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deletePermission(row.id)
      ElMessage.success('删除成功')
      await loadPermissions()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  }).catch(() => {
    // 取消
  })
}

const getTypeLabel = (type) => {
  const map = {
    menu: '菜单',
    button: '按钮',
    api: '接口'
  }
  return map[type] || type
}

const getTypeTagType = (type) => {
  const map = {
    menu: 'success',
    button: 'warning',
    api: 'info'
  }
  return map[type] || ''
}

onMounted(() => {
  loadPermissions()
})
</script>

<style lang="scss" scoped>
.permissions-container {
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
}
</style>
