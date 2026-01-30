<template>
  <div class="not-found-page">
    <AppHeader />
    <div class="not-found">
      <div class="content">
        <h1>404</h1>
        <p>抱歉，您访问的页面不存在</p>
        <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
      </div>
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
</script>

<style scoped>
.not-found-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.not-found {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.content {
  text-align: center;
}

.content h1 {
  font-size: 120px;
  margin: 0 0 20px 0;
  color: #409eff;
}

.content p {
  font-size: 18px;
  color: #666;
  margin: 0 0 30px 0;
}
</style>
