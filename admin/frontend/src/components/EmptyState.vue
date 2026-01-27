<template>
  <div class="empty-state">
    <el-empty
      :description="description"
      :image-size="imageSize"
    >
      <template v-if="$slots.default" #default>
        <slot></slot>
      </template>
      <template v-if="$slots.extra" #extra>
        <slot name="extra"></slot>
      </template>
    </el-empty>
  </div>
</template>

<script setup>
defineProps({
  description: {
    type: String,
    default: '暂无数据'
  },
  imageSize: {
    type: Number,
    default: 100
  }
})
</script>

<style scoped lang="scss">
.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  background: #fff;
  border-radius: 4px;
}
</style>
