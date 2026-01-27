<template>
  <div class="page-designer">
    <div class="toolbar">
      <el-button type="primary" @click="handleSave">
        <el-icon><Document /></el-icon>
        保存页面
      </el-button>
      <el-button @click="handlePreview">
        <el-icon><View /></el-icon>
        预览
      </el-button>
      <el-button @click="handleClear">
        <el-icon><Delete /></el-icon>
        清空
      </el-button>
      <div class="page-title">
        <el-input v-model="pageConfig.title" placeholder="页面标题" />
      </div>
    </div>

    <div class="designer-container">
      <!-- 左侧组件库 -->
      <div class="component-panel">
        <div class="panel-title">布局组件</div>
        <div class="component-list">
          <div
            v-for="item in layoutComponents"
            :key="item.type"
            class="component-item"
            draggable="true"
            @dragstart="handleDragStart(item)"
          >
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.label }}</span>
          </div>
        </div>

        <div class="panel-title">展示组件</div>
        <div class="component-list">
          <div
            v-for="item in displayComponents"
            :key="item.type"
            class="component-item"
            draggable="true"
            @dragstart="handleDragStart(item)"
          >
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.label }}</span>
          </div>
        </div>

        <div class="panel-title">表单组件</div>
        <div class="component-list">
          <div
            v-for="item in formComponents"
            :key="item.type"
            class="component-item"
            draggable="true"
            @dragstart="handleDragStart(item)"
          >
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.label }}</span>
          </div>
        </div>
      </div>

      <!-- 中间画布 -->
      <div
        class="canvas-panel"
        @drop="handleDrop"
        @dragover.prevent
      >
        <div v-if="pageItems.length === 0" class="empty-tips">
          从左侧拖拽组件到此处设计页面
        </div>
        <draggable
          v-else
          v-model="pageItems"
          item-key="id"
          :animation="200"
          @end="handleDragEnd"
        >
          <template #item="{ element, index }">
            <div
              class="page-item-wrapper"
              :class="{ 'is-active': activeIndex === index }"
              @click="activeIndex = index"
            >
              <div class="item-actions">
                <el-button
                  type="danger"
                  :icon="Delete"
                  circle
                  size="small"
                  @click.stop="handleDeleteItem(index)"
                />
              </div>

              <div class="item-label">
                <el-input
                  v-model="element.label"
                  placeholder="组件标题"
                  size="small"
                />
              </div>

              <div class="item-preview">
                <!-- 卡片组件 -->
                <el-card v-if="element.type === 'card'" shadow="hover">
                  <template #header>
                    <span>{{ element.label }}</span>
                  </template>
                  <div class="card-content">
                    卡片内容区域
                  </div>
                </el-card>

                <!-- 统计卡片 -->
                <el-card v-else-if="element.type === 'statcard'" class="stat-card" shadow="hover">
                  <div class="stat-value">{{ element.defaultValue || '0' }}</div>
                  <div class="stat-label">{{ element.label }}</div>
                  <div class="stat-trend" :class="element.trend || 'up'">
                    <el-icon><Top /></el-icon>
                    <span>{{ element.trendValue || '12%' }}</span>
                  </div>
                </el-card>

                <!-- 表格组件 -->
                <el-card v-else-if="element.type === 'table'" shadow="hover">
                  <el-table :data="[]">
                    <el-table-column
                      v-for="col in element.columns || [{ prop: 'name', label: '名称' }]"
                      :key="col.prop"
                      :prop="col.prop"
                      :label="col.label"
                    />
                  </el-table>
                </el-card>

                <!-- 图表组件 -->
                <el-card v-else-if="element.type === 'chart'" shadow="hover">
                  <div ref="chartRef" style="height: 300px;"></div>
                </el-card>

                <!-- 表单组件 -->
                <el-card v-else-if="element.type === 'form'" shadow="hover">
                  <el-form :model="{}" label-width="100px">
                    <el-form-item label="表单项">
                      <el-input placeholder="请输入" disabled />
                    </el-form-item>
                  </el-form>
                </el-card>

                <!-- 按钮组件 -->
                <el-button v-else-if="element.type === 'button'" type="primary">
                  {{ element.label }}
                </el-button>

                <!-- 标签页 -->
                <el-tabs v-else-if="element.type === 'tabs'" type="card">
                  <el-tab-pane label="标签页1" />
                  <el-tab-pane label="标签页2" />
                </el-tabs>

                <!-- 时间线 -->
                <el-timeline v-else-if="element.type === 'timeline'">
                  <el-timeline-item timestamp="2024-01-01">事件1</el-timeline-item>
                  <el-timeline-item timestamp="2024-01-02">事件2</el-timeline-item>
                </el-timeline>

                <!-- 分割线 -->
                <el-divider v-else-if="element.type === 'divider'">
                  {{ element.label }}
                </el-divider>

                <!-- 折叠面板 -->
                <el-collapse v-else-if="element.type === 'collapse'">
                  <el-collapse-item title="面板1" name="1">
                    内容1
                  </el-collapse-item>
                  <el-collapse-item title="面板2" name="2">
                    内容2
                  </el-collapse-item>
                </el-collapse>

                <!-- 警告提示 -->
                <el-alert
                  v-else-if="element.type === 'alert'"
                  :title="element.label"
                  type="info"
                  show-icon
                />

                <!-- 默认显示 -->
                <div v-else class="default-component">
                  {{ element.label }}
                </div>
              </div>
            </div>
          </template>
        </draggable>
      </div>

      <!-- 右侧属性面板 -->
      <div class="property-panel">
        <div class="panel-title">属性配置</div>
        <div v-if="activeIndex >= 0 && pageItems[activeIndex]" class="property-form">
          <el-form label-width="80px" size="small">
            <el-form-item label="组件标题">
              <el-input v-model="pageItems[activeIndex].label" />
            </el-form-item>
            <el-form-item label="组件类型">
              <el-input :value="pageItems[activeIndex].type" disabled />
            </el-form-item>

            <!-- 统计卡片特殊属性 -->
            <template v-if="pageItems[activeIndex].type === 'statcard'">
              <el-form-item label="数值">
                <el-input v-model="pageItems[activeIndex].defaultValue" />
              </el-form-item>
              <el-form-item label="趋势">
                <el-select v-model="pageItems[activeIndex].trend">
                  <el-option label="上升" value="up" />
                  <el-option label="下降" value="down" />
                </el-select>
              </el-form-item>
              <el-form-item label="趋势值">
                <el-input v-model="pageItems[activeIndex].trendValue" />
              </el-form-item>
            </template>

            <!-- 表格列配置 -->
            <template v-if="pageItems[activeIndex].type === 'table'">
              <el-form-item label="表格列">
                <div
                  v-for="(col, idx) in pageItems[activeIndex].columns"
                  :key="idx"
                  class="column-item"
                >
                  <el-input v-model="col.label" placeholder="列标题" size="small" />
                  <el-input v-model="col.prop" placeholder="字段名" size="small" />
                  <el-button
                    type="danger"
                    :icon="Delete"
                    circle
                    size="small"
                    @click="handleDeleteColumn(activeIndex, idx)"
                  />
                </div>
                <el-button
                  type="primary"
                  :icon="Plus"
                  size="small"
                  @click="handleAddColumn(activeIndex)"
                >
                  添加列
                </el-button>
              </el-form-item>
            </template>

            <el-form-item label="背景色">
              <el-input v-model="pageItems[activeIndex].backgroundColor" placeholder="#ffffff" />
            </el-form-item>
            <el-form-item label="圆角">
              <el-input-number v-model="pageItems[activeIndex].borderRadius" :min="0" :max="20" />
            </el-form-item>
          </el-form>
        </div>
        <div v-else class="empty-tips">
          请选择组件配置属性
        </div>
      </div>
    </div>

    <!-- 预览对话框 -->
    <el-dialog v-model="previewVisible" title="页面预览" width="80%" fullscreen>
      <div class="preview-container">
        <el-row :gutter="20">
          <el-col
            v-for="item in pageItems"
            :key="item.id"
            :span="item.span || 12"
          >
            <el-card shadow="hover">
              <template #header>
                <span>{{ item.label }}</span>
              </template>
              <div>{{ item.type }}</div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Document, View, Delete, Plus, Grid, DataAnalysis,
  Picture, Calendar, Tickets, Grid as GridIcon,
  Top, ChatDotRound, Warning
} from '@element-plus/icons-vue'
import draggable from 'vuedraggable'
import { savePageDesign } from '@/api/lowcode'

const layoutComponents = [
  { type: 'card', label: '卡片', icon: Picture },
  { type: 'statcard', label: '统计卡片', icon: DataAnalysis },
  { type: 'divider', label: '分割线', icon: GridIcon },
  { type: 'tabs', label: '标签页', icon: Tickets },
  { type: 'collapse', label: '折叠面板', icon: GridIcon }
]

const displayComponents = [
  { type: 'table', label: '表格', icon: Grid },
  { type: 'chart', label: '图表', icon: DataAnalysis },
  { type: 'timeline', label: '时间线', icon: Calendar },
  { type: 'alert', label: '警告提示', icon: Warning }
]

const formComponents = [
  { type: 'form', label: '表单', icon: ChatDotRound },
  { type: 'button', label: '按钮', icon: Plus }
]

const pageItems = ref([])
const activeIndex = ref(-1)
const previewVisible = ref(false)


const pageConfig = reactive({
  title: '',
  description: ''
})

let draggedComponent = null

const handleDragStart = (item) => {
  draggedComponent = item
}

const handleDrop = (e) => {
  e.preventDefault()
  if (!draggedComponent) return

  const newItem = {
    id: Date.now(),
    type: draggedComponent.type,
    label: draggedComponent.label,
    span: 12,
    backgroundColor: '',
    borderRadius: 4,
    defaultValue: '',
    columns: draggedComponent.type === 'table' ? [{ prop: 'name', label: '名称' }] : [],
    trend: 'up',
    trendValue: '12%'
  }

  pageItems.value.push(newItem)
  activeIndex.value = pageItems.value.length - 1
  draggedComponent = null
}

const handleDragEnd = () => {
  activeIndex.value = -1
}

const handleDeleteItem = (index) => {
  pageItems.value.splice(index, 1)
  activeIndex.value = -1
}

const handleAddColumn = (itemIndex) => {
  pageItems.value[itemIndex].columns.push({
    label: '新列',
    prop: `field_${Date.now()}`
  })
}

const handleDeleteColumn = (itemIndex, colIndex) => {
  pageItems.value[itemIndex].columns.splice(colIndex, 1)
}

const handleSave = async () => {
  if (!pageConfig.title) {
    ElMessage.warning('请输入页面标题')
    return
  }

  try {
    await savePageDesign({
      title: pageConfig.title,
      description: pageConfig.description,
      items: pageItems.value
    })
    ElMessage.success('保存成功')
  } catch (error) {
    console.error(error)
  }
}

const handlePreview = () => {
  previewVisible.value = true
}

const handleClear = () => {
  pageItems.value = []
  activeIndex.value = -1
}
</script>

<style lang="scss" scoped>
.page-designer {
  height: 100%;
  display: flex;
  flex-direction: column;

  .toolbar {
    padding: 16px;
    background: #fff;
    border-bottom: 1px solid #e4e7ed;
    display: flex;
    align-items: center;
    gap: 12px;

    .page-title {
      flex: 1;
      max-width: 300px;
      margin-left: 20px;
    }
  }

  .designer-container {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .component-panel,
  .property-panel {
    width: 280px;
    background: #f5f7fa;
    border-right: 1px solid #e4e7ed;
    overflow-y: auto;

    .panel-title {
      padding: 12px 16px;
      font-weight: bold;
      color: #606266;
      border-bottom: 1px solid #e4e7ed;
    }
  }

  .property-panel {
    border-right: none;
    border-left: 1px solid #e4e7ed;
  }

  .component-list {
    padding: 8px;

    .component-item {
      display: flex;
      align-items: center;
      padding: 10px;
      margin-bottom: 8px;
      background: #fff;
      border: 1px solid #dcdfe6;
      border-radius: 4px;
      cursor: move;
      transition: all 0.3s;

      &:hover {
        border-color: #409eff;
        box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
      }

      .el-icon {
        margin-right: 8px;
        color: #409eff;
      }
    }
  }

  .canvas-panel {
    flex: 1;
    padding: 20px;
    background: #fff;
    overflow-y: auto;

    .empty-tips {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
      color: #909399;
      font-size: 14px;
    }

    .page-item-wrapper {
      position: relative;
      margin-bottom: 16px;
      border: 2px dashed #dcdfe6;
      border-radius: 8px;
      background: #fafafa;
      cursor: pointer;
      overflow: hidden;

      &:hover,
      &.is-active {
        border-color: #409eff;
        background: #ecf5ff;
      }

      .item-actions {
        position: absolute;
        right: 8px;
        top: 8px;
        z-index: 10;
        display: none;
      }

      &:hover .item-actions,
      &.is-active .item-actions {
        display: block;
      }

      .item-label {
        padding: 12px 16px 8px;
      }

      .item-preview {
        padding: 0 16px 16px;

        .stat-card {
          text-align: center;
          padding: 20px;

          .stat-value {
            font-size: 36px;
            font-weight: bold;
            color: #303133;
            margin-bottom: 8px;
          }

          .stat-label {
            font-size: 14px;
            color: #606266;
            margin-bottom: 12px;
          }

          .stat-trend {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 4px;
            font-size: 14px;

            &.up {
              color: #67c23a;
            }

            &.down {
              color: #f56c6c;
            }
          }
        }

        .default-component {
          padding: 20px;
          text-align: center;
          color: #909399;
        }
      }
    }
  }

  .property-form {
    padding: 16px;

    .column-item {
      display: flex;
      gap: 8px;
      margin-bottom: 8px;
    }
  }

  .empty-tips {
    padding: 40px;
    text-align: center;
    color: #909399;
  }

  .preview-container {
    padding: 20px;
    background: #f5f7fa;
    min-height: 400px;
  }
}
</style>
