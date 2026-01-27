<template>
  <div class="data-table">
    <!-- 工具栏 -->
    <div v-if="showToolbar" class="table-toolbar">
      <div class="toolbar-left">
        <el-checkbox
          v-if="showSelection"
          v-model="selectAll"
          :indeterminate="isIndeterminate"
          @change="handleSelectAll"
        >
          全选
        </el-checkbox>
        <slot name="toolbar-left"></slot>
      </div>
      <div class="toolbar-right">
        <slot name="toolbar-right">
          <el-input
            v-if="showSearch"
            v-model="searchKeyword"
            placeholder="搜索..."
            clearable
            :prefix-icon="Search"
            style="width: 200px"
            @input="handleSearch"
          />
        </slot>
      </div>
    </div>

    <!-- 表格 -->
    <el-table
      ref="tableRef"
      v-loading="loading"
      :data="data"
      :height="height"
      :max-height="maxHeight"
      :stripe="stripe"
      :border="border"
      :size="size"
      @selection-change="handleSelectionChange"
      @sort-change="handleSortChange"
    >
      <el-table-column
        v-if="showSelection"
        type="selection"
        width="55"
        :reserve-selection="reserveSelection"
      />
      <el-table-column
        v-if="showIndex"
        type="index"
        label="序号"
        width="80"
        :index="indexMethod"
      />
      <slot></slot>
      <!-- 操作列 -->
      <el-table-column
        v-if="showActions"
        label="操作"
        :width="actionWidth"
        :fixed="actionFixed"
      >
        <template #default="{ row, $index }">
          <slot name="actions" :row="row" :index="$index"></slot>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div v-if="showPagination" class="table-pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="pageSizes"
        :layout="paginationLayout"
        :background="true"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 空状态 -->
    <div v-if="!loading && data.length === 0" class="table-empty">
      <slot name="empty">
        <el-empty :description="emptyText"></el-empty>
      </slot>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Search } from '@element-plus/icons-vue'

const props = defineProps({
  data: {
    type: Array,
    default: () => []
  },
  loading: Boolean,
  // 表格属性
  height: [String, Number],
  maxHeight: [String, Number],
  stripe: {
    type: Boolean,
    default: true
  },
  border: {
    type: Boolean,
    default: true
  },
  size: {
    type: String,
    default: 'default', // large, default, small
    validator: (value) => ['large', 'default', 'small'].includes(value)
  },
  // 功能开关
  showToolbar: {
    type: Boolean,
    default: true
  },
  showSelection: {
    type: Boolean,
    default: false
  },
  showIndex: {
    type: Boolean,
    default: false
  },
  showSearch: {
    type: Boolean,
    default: true
  },
  showActions: {
    type: Boolean,
    default: true
  },
  showPagination: {
    type: Boolean,
    default: true
  },
  // 选择相关
  reserveSelection: Boolean,
  // 索引方法
  indexMethod: Function,
  // 操作列
  actionWidth: {
    type: String,
    default: '200'
  },
  actionFixed: {
    type: String,
    default: 'right'
  },
  // 分页
  total: {
    type: Number,
    default: 0
  },
  currentPage: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  pageSizes: {
    type: Array,
    default: () => [10, 20, 50, 100]
  },
  paginationLayout: {
    type: String,
    default: 'total, sizes, prev, pager, next, jumper'
  },
  // 空状态
  emptyText: {
    type: String,
    default: '暂无数据'
  }
})

const emit = defineEmits([
  'update:currentPage',
  'update:pageSize',
  'search',
  'selection-change',
  'sort-change'
])

const tableRef = ref(null)
const selectedRows = ref([])
const selectAll = ref(false)
const isIndeterminate = computed(() =>
  selectedRows.value.length > 0 && selectedRows.value.length < props.data.length
)
const searchKeyword = ref('')

const handleSelectAll = (val) => {
  if (val) {
    tableRef.value?.toggleAllSelection()
  } else {
    tableRef.value?.clearSelection()
  }
}

const handleSelectionChange = (selection) => {
  selectedRows.value = selection
  selectAll.value = selection.length === props.data.length
  emit('selection-change', selection)
}

const handleSortChange = (sort) => {
  emit('sort-change', sort)
}

const handleSearch = (value) => {
  emit('search', value)
}

const handleSizeChange = (size) => {
  emit('update:pageSize', size)
}

const handleCurrentChange = (page) => {
  emit('update:currentPage', page)
}

// 暴露方法
defineExpose({
  clearSelection: () => tableRef.value?.clearSelection(),
  toggleRowSelection: (row, selected) =>
    tableRef.value?.toggleRowSelection(row, selected),
  toggleAllSelection: () => tableRef.value?.toggleAllSelection(),
  getSelectionRows: () => tableRef.value?.getSelectionRows()
})
</script>

<style scoped lang="scss">
.data-table {
  .table-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    padding: 16px;
    background: #f5f7fa;
    border-radius: 4px;

    .toolbar-left {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .toolbar-right {
      display: flex;
      align-items: center;
      gap: 12px;
    }
  }

  .table-pagination {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
    padding: 16px 0;
  }

  .table-empty {
    padding: 40px 0;
  }
}
</style>
