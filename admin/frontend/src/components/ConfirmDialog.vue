<template>
  <el-dialog
    v-model="visible"
    :title="title"
    :width="width"
    :before-close="handleClose"
  >
    <div class="confirm-content">
      <el-icon class="confirm-icon" :class="iconType">
        <component :is="iconComponent" />
      </el-icon>
      <p class="confirm-message">{{ message }}</p>
      <p v-if="detail" class="confirm-detail">{{ detail }}</p>
    </div>
    <template #footer>
      <el-button @click="handleClose">{{ cancelText }}</el-button>
      <el-button :type="confirmType" :loading="loading" @click="handleConfirm">
        {{ confirmText }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed } from 'vue'
import { Warning, InfoFilled, CircleCloseFilled, CircleCheckFilled } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: Boolean,
  title: {
    type: String,
    default: '操作确认'
  },
  message: {
    type: String,
    default: '确定要进行此操作吗？'
  },
  detail: String,
  type: {
    type: String,
    default: 'warning', // warning, info, error, success
    validator: (value) => ['warning', 'info', 'error', 'success'].includes(value)
  },
  confirmText: {
    type: String,
    default: '确定'
  },
  cancelText: {
    type: String,
    default: '取消'
  },
  loading: Boolean,
  width: {
    type: String,
    default: '420px'
  }
})

const emit = defineEmits(['update:modelValue', 'confirm'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const iconType = computed(() => `icon-${props.type}`)

const iconComponent = computed(() => {
  const icons = {
    warning: Warning,
    info: InfoFilled,
    error: CircleCloseFilled,
    success: CircleCheckFilled
  }
  return icons[props.type] || Warning
})

const confirmType = computed(() => {
  const types = {
    warning: 'primary',
    info: 'primary',
    error: 'danger',
    success: 'success'
  }
  return types[props.type] || 'primary'
})

const handleClose = () => {
  if (!props.loading) {
    visible.value = false
  }
}

const handleConfirm = () => {
  emit('confirm')
}
</script>

<style scoped lang="scss">
.confirm-content {
  text-align: center;
  padding: 20px 0;

  .confirm-icon {
    font-size: 48px;
    margin-bottom: 16px;

    &.icon-warning {
      color: #e6a23c;
    }

    &.icon-info {
      color: #909399;
    }

    &.icon-error {
      color: #f56c6c;
    }

    &.icon-success {
      color: #67c23a;
    }
  }

  .confirm-message {
    font-size: 16px;
    color: #303133;
    margin: 0;
  }

  .confirm-detail {
    font-size: 14px;
    color: #909399;
    margin: 12px 0 0 0;
  }
}
</style>
