<template>
  <div class="counselors-page">
    <AppHeader />
    <div class="counselors">
      <div class="page-header">
        <h2>咨询师列表</h2>
        <p class="subtitle">专业认证心理咨询师，为您提供一对一咨询服务</p>
      </div>

      <el-card class="search-card">
        <el-form :inline="true" :model="searchForm">
          <el-form-item label="关键词">
            <el-input v-model="searchForm.keyword" placeholder="搜索咨询师姓名或专长" clearable style="width: 300px" />
          </el-form-item>
          <el-form-item label="排序方式">
            <el-select v-model="searchForm.sort_by" placeholder="请选择排序" style="width: 150px" @change="handleSortChange">
              <el-option label="评分排序" value="rating" />
              <el-option label="价格排序" value="price" />
              <el-option label="服务人数" value="service_count" />
            </el-select>
          </el-form-item>
          <el-form-item v-if="searchForm.sort_by">
            <el-radio-group v-model="searchForm.sort_order" @change="handleSortChange">
              <el-radio-button value="desc">降序</el-radio-button>
              <el-radio-button value="asc">升序</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <el-row :gutter="24" v-loading="loading">
        <el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="6" v-for="counselor in counselors" :key="counselor.id">
          <el-card class="counselor-card" @click="$router.push(`/counselor/${counselor.id}`)" style="cursor: pointer">
            <!-- 推荐标签 -->
            <div v-if="counselor.is_recommended" class="recommend-badge">
              <el-icon><Star /></el-icon>
              <span>推荐咨询师</span>
            </div>

            <!-- 头像区域 -->
            <div class="counselor-avatar">
              <el-avatar :size="100" :src="counselor.avatar">
                <el-icon><User /></el-icon>
              </el-avatar>
              <div class="avatar-ring"></div>
            </div>

            <!-- 基本信息 -->
            <h3 class="counselor-name">{{ counselor.name }}</h3>
            <p class="counselor-title">{{ counselor.title }}</p>

            <!-- 擅长领域标签 -->
            <div class="specialty-tags">
              <el-tag
                v-for="(tag, index) in getSpecialtyTags(counselor.specialty)"
                :key="index"
                size="small"
                type="info"
                effect="plain"
              >
                {{ tag }}
              </el-tag>
            </div>

            <!-- 统计信息 -->
            <div class="stats-info">
              <div class="stat-item">
                <el-icon class="stat-icon"><User /></el-icon>
                <div class="stat-content">
                  <span class="stat-value">{{ counselor.service_count || 0 }}</span>
                  <span class="stat-label">服务人数</span>
                </div>
              </div>
              <div class="stat-divider"></div>
              <div class="stat-item">
                <el-icon class="stat-icon"><ChatDotRound /></el-icon>
                <div class="stat-content">
                  <span class="stat-value">{{ counselor.review_count || 0 }}</span>
                  <span class="stat-label">评价数</span>
                </div>
              </div>
              <div class="stat-divider"></div>
              <div class="stat-item">
                <el-icon class="stat-icon"><Calendar /></el-icon>
                <div class="stat-content">
                  <span class="stat-value">{{ counselor.years_exp }}年</span>
                  <span class="stat-label">从业年限</span>
                </div>
              </div>
            </div>

            <!-- 评分 -->
            <div class="rating-section">
              <el-rate v-model="counselor.rating" disabled show-score text-color="#ff9900" score-template="{value}" />
              <span v-if="counselor.is_recommended" class="recommend-text">优质推荐</span>
            </div>

            <!-- 底部信息 -->
            <div class="card-footer">
              <div class="price-info">
                <span class="price-label">咨询费</span>
                <span class="price-value">¥{{ counselor.price }}/分钟</span>
              </div>
              <el-button type="primary" size="small" @click.stop="$router.push(`/counselor/${counselor.id}`)">
                立即预约
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <div class="pagination" v-if="total > 0">
        <el-pagination
          v-model:current-page="searchForm.page"
          v-model:page-size="searchForm.page_size"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadCounselors"
          @current-change="loadCounselors"
        />
      </div>

      <el-empty v-if="!loading && counselors.length === 0" description="暂无咨询师" />
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import {
  User,
  Calendar,
  Search,
  Star,
  ChatDotRound
} from '@element-plus/icons-vue'
import { getCounselorList } from '@/api/counselor'
import { handleError } from '@/utils/errorHandler'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'

const route = useRoute()
const loading = ref(false)
const counselors = ref([])
const total = ref(0)

const searchForm = reactive({
  page: 1,
  page_size: 10,
  keyword: '',
  sort_by: 'rating',
  sort_order: 'desc'
})

onMounted(() => {
  // 检查是否有服务类型参数
  if (route.query.serviceType) {
    searchForm.keyword = route.query.serviceType
  }
  loadCounselors()
})

const loadCounselors = async () => {
  loading.value = true
  try {
    const res = await getCounselorList(searchForm)
    counselors.value = res.data.counselors || []
    total.value = res.data.total || 0
  } catch (error) {
    handleError(error, '加载咨询师失败')
  } finally {
    loading.value = false
  }
}

const getSpecialtyTags = (specialty) => {
  if (!specialty) return []
  return specialty.split(/[,，]/).map(t => t.trim()).filter(t => t)
}

const handleSearch = () => {
  searchForm.page = 1
  loadCounselors()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.page = 1
  searchForm.sort_by = 'rating'
  searchForm.sort_order = 'desc'
  loadCounselors()
}

const handleSortChange = () => {
  loadCounselors()
}
</script>

<style scoped>
.counselors-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eef5 100%);
}

.counselors {
  flex: 1;
  padding: 30px 20px;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}

.page-header {
  text-align: center;
  margin-bottom: 30px;
}

.page-header h2 {
  margin: 0 0 10px 0;
  font-size: 32px;
  color: #2c3e50;
  font-weight: 600;
}

.page-header .subtitle {
  margin: 0;
  font-size: 16px;
  color: #7f8c8d;
}

.search-card {
  margin-bottom: 30px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.counselor-card {
  text-align: center;
  margin-bottom: 24px;
  border-radius: 16px;
  border: none;
  position: relative;
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  background: white;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.counselor-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 120px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  z-index: 0;
}

.counselor-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.15);
}

/* 推荐标签 */
.recommend-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  z-index: 10;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(240, 147, 251, 0.4);
}

.recommend-badge .el-icon {
  font-size: 14px;
}

/* 头像区域 */
.counselor-avatar {
  position: relative;
  z-index: 1;
  margin-bottom: 16px;
  padding-top: 30px;
}

.avatar-ring {
  position: absolute;
  top: 30px;
  left: 50%;
  transform: translateX(-50%);
  width: 100px;
  height: 100px;
  border-radius: 50%;
  border: 3px solid white;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.counselor-name {
  position: relative;
  z-index: 1;
  margin: 0 0 4px 0;
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
}

.counselor-title {
  position: relative;
  z-index: 1;
  margin: 0 0 16px 0;
  font-size: 14px;
  color: #7f8c8d;
}

/* 擅长领域标签 */
.specialty-tags {
  position: relative;
  z-index: 1;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: center;
  margin-bottom: 20px;
  padding: 0 10px;
}

/* 统计信息 */
.stats-info {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: 16px 10px;
  background: #f8f9fa;
  border-radius: 12px;
  margin-bottom: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-icon {
  font-size: 24px;
  color: #667eea;
}

.stat-content {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  line-height: 1.2;
}

.stat-label {
  font-size: 11px;
  color: #95a5a6;
  margin-top: 2px;
}

.stat-divider {
  width: 1px;
  height: 32px;
  background: #e0e0e0;
}

/* 评分区域 */
.rating-section {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 16px;
}

.rating-section .el-rate {
  flex-shrink: 0;
}

.recommend-text {
  padding: 4px 10px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

/* 底部信息 */
.card-footer {
  position: relative;
  z-index: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0 0 0;
  border-top: 1px solid #e8e8e8;
}

.price-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.price-label {
  font-size: 11px;
  color: #95a5a6;
  margin-bottom: 2px;
}

.price-value {
  font-size: 20px;
  font-weight: 700;
  color: #e74c3c;
}

.card-footer .el-button {
  font-weight: 500;
  padding: 10px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

.card-footer .el-button:hover {
  transform: scale(1.05);
}

.pagination {
  margin-top: 30px;
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .counselors {
    padding: 20px 15px;
  }

  .page-header h2 {
    font-size: 24px;
  }

  .page-header .subtitle {
    font-size: 14px;
  }

  .search-card {
    margin-bottom: 20px;
  }

  .counselor-card {
    margin-bottom: 16px;
  }

  .stats-info {
    padding: 12px 8px;
  }

  .stat-value {
    font-size: 16px;
  }

  .price-value {
    font-size: 18px;
  }
}
</style>
