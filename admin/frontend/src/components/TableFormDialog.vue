<template>
  <el-dialog
    v-model="visible"
    :title="title"
    :width="width"
    :before-close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="100px"
    >
      <el-row :gutter="20">
        <el-col
          v-for="field in fields"
          :key="field.prop"
          :span="field.span || 12"
        >
          <el-form-item :label="field.label" :prop="field.prop">
            <!-- 输入框 -->
            <el-input
              v-if="field.type === 'input'"
              v-model="formData[field.prop]"
              :placeholder="field.placeholder"
              :disabled="field.disabled"
              clearable
            />

            <!-- 文本域 -->
            <el-input
              v-else-if="field.type === 'textarea'"
              v-model="formData[field.prop]"
              type="textarea"
              :placeholder="field.placeholder"
              :disabled="field.disabled"
              :rows="field.rows || 3"
            />

            <!-- 数字输入框 -->
            <el-input-number
              v-else-if="field.type === 'number'"
              v-model="formData[field.prop]"
              :min="field.min"
              :max="field.max"
              :disabled="field.disabled"
              :step="field.step || 1"
              style="width: 100%"
            />

            <!-- 选择器 -->
            <el-select
              v-else-if="field.type === 'select'"
              v-model="formData[field.prop]"
              :placeholder="field.placeholder"
              :disabled="field.disabled"
              clearable
              style="width: 100%"
            >
              <el-option
                v-for="opt in field.options"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              />
            </el-select>

            <!-- 日期选择器 -->
            <el-date-picker
              v-else-if="field.type === 'date'"
              v-model="formData[field.prop]"
              type="date"
              :placeholder="field.placeholder"
              :disabled="field.disabled"
              style="width: 100%"
            />

            <!-- 日期时间选择器 -->
            <el-date-picker
              v-else-if="field.type === 'datetime'"
              v-model="formData[field.prop]"
              type="datetime"
              :placeholder="field.placeholder"
              :disabled="field.disabled"
              style="width: 100%"
            />

            <!-- 开关 -->
            <el-switch
              v-else-if="field.type === 'switch'"
              v-model="formData[field.prop]"
              :disabled="field.disabled"
            />

            <!-- 单选框组 -->
            <el-radio-group
              v-else-if="field.type === 'radio'"
              v-model="formData[field.prop]"
              :disabled="field.disabled"
            >
              <el-radio
                v-for="opt in field.options"
                :key="opt.value"
                :label="opt.value"
              >
                {{ opt.label }}
              </el-radio>
            </el-radio-group>

            <!-- 复选框组 -->
            <el-checkbox-group
              v-else-if="field.type === 'checkbox'"
              v-model="formData[field.prop]"
              :disabled="field.disabled"
            >
              <el-checkbox
                v-for="opt in field.options"
                :key="opt.value"
                :label="opt.value"
              >
                {{ opt.label }}
              </el-checkbox>
            </el-checkbox-group>

            <!-- 自定义插槽 -->
            <slot v-else :name="field.prop" :field="field"></slot>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        确定
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'

const props = defineProps({
  modelValue: Boolean,
  title: {
    type: String,
    default: '编辑'
  },
  fields: {
    type: Array,
    required: true,
    default: () => []
  },
  rules: {
    type: Object,
    default: () => ({})
  },
  loading: Boolean,
  width: {
    type: String,
    default: '800px'
  }
})

const emit = defineEmits(['update:modelValue', 'submit'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const formRef = ref(null)
const formData = reactive({})

// 初始化表单数据
const initFormData = (data = {}) => {
  props.fields.forEach(field => {
    if (data[field.prop] !== undefined) {
      formData[field.prop] = data[field.prop]
    } else if (field.defaultValue !== undefined) {
      formData[field.prop] = field.defaultValue
    } else {
      formData[field.prop] = ''
    }
  })
}

// 监听对话框打开
watch(() => props.modelValue, (val) => {
  if (val) {
    initFormData()
    formRef.value?.clearValidate()
  }
})

const handleClose = () => {
  if (!props.loading) {
    visible.value = false
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    emit('submit', { ...formData })
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

// 暴露方法
defineExpose({
  initFormData,
  resetFields: () => formRef.value?.resetFields(),
  validate: () => formRef.value?.validate()
})
</script>

<style scoped lang="scss">
:deep(.el-form-item__content) {
  .el-input,
  .el-select,
  .el-date-picker {
    width: 100%;
  }
}
</style>
