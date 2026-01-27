<template>
  <div class="form-designer">
    <div class="toolbar">
      <el-button type="primary" @click="handleSave">
        <el-icon><Document /></el-icon>
        保存表单
      </el-button>
      <el-button @click="handlePreview">
        <el-icon><View /></el-icon>
        预览
      </el-button>
      <el-button @click="handleClear">
        <el-icon><Delete /></el-icon>
        清空
      </el-button>
      <div class="form-title">
        <el-input v-model="formConfig.title" placeholder="表单标题" />
      </div>
    </div>
    
    <div class="designer-container">
      <!-- 左侧组件库 -->
      <div class="component-panel">
        <div class="panel-title">基础组件</div>
        <div class="component-list">
          <div
            v-for="item in baseComponents"
            :key="item.type"
            class="component-item"
            draggable="true"
            @dragstart="handleDragStart(item)"
          >
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.label }}</span>
          </div>
        </div>
        
        <div class="panel-title">高级组件</div>
        <div class="component-list">
          <div
            v-for="item in advancedComponents"
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
        <div v-if="formItems.length === 0" class="empty-tips">
          从左侧拖拽组件到此处设计表单
        </div>
        <draggable
          v-else
          v-model="formItems"
          item-key="id"
          :animation="200"
          @end="handleDragEnd"
        >
          <template #item="{ element, index }">
            <div
              class="form-item-wrapper"
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
                  placeholder="字段标题"
                  size="small"
                />
              </div>
              
              <div class="item-field">
                <el-input
                  v-if="element.type === 'input'"
                  v-model="element.defaultValue"
                  placeholder="默认值"
                  disabled
                />
                <el-input
                  v-else-if="element.type === 'textarea'"
                  v-model="element.defaultValue"
                  type="textarea"
                  :rows="2"
                  placeholder="默认值"
                  disabled
                />
                <el-input-number
                  v-else-if="element.type === 'number'"
                  v-model="element.defaultValue"
                  disabled
                />
                <el-date-picker
                  v-else-if="element.type === 'date'"
                  v-model="element.defaultValue"
                  type="date"
                  placeholder="选择日期"
                  disabled
                  style="width: 100%"
                />
                <el-time-picker
                  v-else-if="element.type === 'time'"
                  v-model="element.defaultValue"
                  placeholder="选择时间"
                  disabled
                  style="width: 100%"
                />
                <el-select
                  v-else-if="element.type === 'select'"
                  v-model="element.defaultValue"
                  placeholder="请选择"
                  disabled
                  style="width: 100%"
                >
                  <el-option
                    v-for="opt in element.options"
                    :key="opt.value"
                    :label="opt.label"
                    :value="opt.value"
                  />
                </el-select>
                <el-radio-group v-else-if="element.type === 'radio'" disabled>
                  <el-radio
                    v-for="opt in element.options"
                    :key="opt.value"
                    :label="opt.value"
                  >
                    {{ opt.label }}
                  </el-radio>
                </el-radio-group>
                <el-checkbox-group v-else-if="element.type === 'checkbox'" disabled>
                  <el-checkbox
                    v-for="opt in element.options"
                    :key="opt.value"
                    :label="opt.value"
                  >
                    {{ opt.label }}
                  </el-checkbox>
                </el-checkbox-group>
                <el-switch v-else-if="element.type === 'switch'" disabled />
                <el-upload
                  v-else-if="element.type === 'upload'"
                  disabled
                  action="#"
                  list-type="picture-card"
                >
                  <el-icon><Plus /></el-icon>
                </el-upload>
              </div>
            </div>
          </template>
        </draggable>
      </div>
      
      <!-- 右侧属性面板 -->
      <div class="property-panel">
        <div class="panel-title">属性配置</div>
        <div v-if="activeIndex >= 0 && formItems[activeIndex]" class="property-form">
          <el-form label-width="80px" size="small">
            <el-form-item label="字段名">
              <el-input v-model="formItems[activeIndex].field" />
            </el-form-item>
            <el-form-item label="占位符">
              <el-input v-model="formItems[activeIndex].placeholder" />
            </el-form-item>
            <el-form-item label="默认值">
              <el-input v-model="formItems[activeIndex].defaultValue" />
            </el-form-item>
            <el-form-item label="是否必填">
              <el-switch v-model="formItems[activeIndex].required" />
            </el-form-item>
            <el-form-item label="是否禁用">
              <el-switch v-model="formItems[activeIndex].disabled" />
            </el-form-item>
            
            <!-- 选项配置 -->
            <template v-if="['select', 'radio', 'checkbox'].includes(formItems[activeIndex].type)">
              <el-form-item label="选项配置">
                <div
                  v-for="(opt, idx) in formItems[activeIndex].options"
                  :key="idx"
                  class="option-item"
                >
                  <el-input v-model="opt.label" placeholder="显示值" size="small" />
                  <el-input v-model="opt.value" placeholder="实际值" size="small" />
                  <el-button
                    type="danger"
                    :icon="Delete"
                    circle
                    size="small"
                    @click="handleDeleteOption(activeIndex, idx)"
                  />
                </div>
                <el-button
                  type="primary"
                  :icon="Plus"
                  size="small"
                  @click="handleAddOption(activeIndex)"
                >
                  添加选项
                </el-button>
              </el-form-item>
            </template>
          </el-form>
        </div>
        <div v-else class="empty-tips">
          请选择表单组件配置属性
        </div>
      </div>
    </div>
    
    <!-- 预览对话框 -->
    <el-dialog v-model="previewVisible" title="表单预览" width="600px">
      <el-form :model="previewForm" label-width="100px">
        <el-form-item
          v-for="item in formItems"
          :key="item.id"
          :label="item.label"
          :required="item.required"
        >
          <el-input
            v-if="item.type === 'input'"
            v-model="previewForm[item.field]"
            :placeholder="item.placeholder"
          />
          <el-input
            v-else-if="item.type === 'textarea'"
            v-model="previewForm[item.field]"
            type="textarea"
            :rows="3"
            :placeholder="item.placeholder"
          />
          <el-input-number
            v-else-if="item.type === 'number'"
            v-model="previewForm[item.field]"
          />
          <el-date-picker
            v-else-if="item.type === 'date'"
            v-model="previewForm[item.field]"
            type="date"
            style="width: 100%"
          />
          <el-select
            v-else-if="item.type === 'select'"
            v-model="previewForm[item.field]"
            style="width: 100%"
          >
            <el-option
              v-for="opt in item.options"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </el-select>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Document, View, Delete, Plus, Edit, Calendar, Timer, List, Rank,
  CheckBox, Switch, Upload
} from '@element-plus/icons-vue'
import draggable from 'vuedraggable'
import { saveFormDesign } from '@/api/lowcode'

const baseComponents = [
  { type: 'input', label: '单行文本', icon: Edit },
  { type: 'textarea', label: '多行文本', icon: Edit },
  { type: 'number', label: '数字输入', icon: Edit },
  { type: 'date', label: '日期选择', icon: Calendar },
  { type: 'time', label: '时间选择', icon: Timer },
  { type: 'select', label: '下拉选择', icon: List },
  { type: 'radio', label: '单选框', icon: Rank },
  { type: 'checkbox', label: '多选框', icon: CheckBox },
  { type: 'switch', label: '开关', icon: Switch },
  { type: 'upload', label: '文件上传', icon: Upload }
]

const advancedComponents = [
  { type: 'richtext', label: '富文本', icon: Edit },
  { type: 'cascader', label: '级联选择', icon: List },
  { type: 'transfer', label: '穿梭框', icon: List }
]

const formItems = ref([])
const activeIndex = ref(-1)
const previewVisible = ref(false)
const previewForm = reactive({})

const formConfig = reactive({
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
    field: `field_${Date.now()}`,
    placeholder: '请输入',
    defaultValue: '',
    required: false,
    disabled: false,
    options: draggedComponent.type === 'select' || draggedComponent.type === 'radio' || draggedComponent.type === 'checkbox'
      ? [{ label: '选项1', value: '1' }, { label: '选项2', value: '2' }]
      : []
  }
  
  formItems.value.push(newItem)
  activeIndex.value = formItems.value.length - 1
  draggedComponent = null
}

const handleDragEnd = () => {
  activeIndex.value = -1
}

const handleDeleteItem = (index) => {
  formItems.value.splice(index, 1)
  activeIndex.value = -1
}

const handleAddOption = (itemIndex) => {
  formItems.value[itemIndex].options.push({
    label: `选项${formItems.value[itemIndex].options.length + 1}`,
    value: `${formItems.value[itemIndex].options.length + 1}`
  })
}

const handleDeleteOption = (itemIndex, optIndex) => {
  formItems.value[itemIndex].options.splice(optIndex, 1)
}

const handleSave = async () => {
  if (!formConfig.title) {
    ElMessage.warning('请输入表单标题')
    return
  }
  
  try {
    await saveFormDesign({
      title: formConfig.title,
      description: formConfig.description,
      items: formItems.value
    })
    ElMessage.success('保存成功')
  } catch (error) {
    console.error(error)
  }
}

const handlePreview = () => {
  Object.keys(previewForm).forEach(key => delete previewForm[key])
  formItems.value.forEach(item => {
    if (item.field) {
      previewForm[item.field] = item.defaultValue || ''
    }
  })
  previewVisible.value = true
}

const handleClear = () => {
  formItems.value = []
  activeIndex.value = -1
}
</script>

<style lang="scss" scoped>
.form-designer {
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
    
    .form-title {
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
    
    .form-item-wrapper {
      position: relative;
      padding: 16px;
      margin-bottom: 16px;
      border: 2px dashed #dcdfe6;
      border-radius: 4px;
      background: #fafafa;
      cursor: pointer;
      
      &:hover,
      &.is-active {
        border-color: #409eff;
        background: #ecf5ff;
      }
      
      .item-actions {
        position: absolute;
        right: 8px;
        top: 8px;
        display: none;
      }
      
      &:hover .item-actions,
      &.is-active .item-actions {
        display: block;
      }
      
      .item-label {
        margin-bottom: 8px;
      }
    }
  }
  
  .property-form {
    padding: 16px;
    
    .option-item {
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
}
</style>
