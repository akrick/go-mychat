<template>
  <div class="main-layout">
    <AppHeader />
    <div class="main-content">
      <slot></slot>
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
}
</style>
