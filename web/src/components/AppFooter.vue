<template>
  <div class="app-footer">
    <div class="footer-top">
      <div class="container">
        <el-row :gutter="40">
          <el-col :xs="24" :sm="12" :md="6">
            <div class="footer-section">
              <div class="logo">
                <el-icon size="40" color="#409eff"><ChatDotRound /></el-icon>
                <div class="logo-text">
                  <span class="logo-name">心理咨询</span>
                </div>
              </div>
              <p class="company-desc">
                专业的心理咨询服务平台，致力于为用户提供高品质的心理咨询服务，用心守护每一位用户的心理健康。
              </p>
              <div class="social-links">
                <el-icon class="social-icon"><ChatLineSquare /></el-icon>
                <el-icon class="social-icon"><Promotion /></el-icon>
                <el-icon class="social-icon"><Position /></el-icon>
              </div>
            </div>
          </el-col>
          <el-col :xs="12" :sm="6" :md="5">
            <div class="footer-section">
              <h3>快速链接</h3>
              <ul class="footer-links">
                <li><router-link to="/">首页</router-link></li>
                <li><router-link to="/counselors">咨询师</router-link></li>
                <li><router-link to="/orders">我的订单</router-link></li>
                <li><router-link to="/profile">个人中心</router-link></li>
              </ul>
            </div>
          </el-col>
          <el-col :xs="12" :sm="6" :md="5">
            <div class="footer-section">
              <h3>服务项目</h3>
              <ul class="footer-links">
                <li><router-link to="/services">个人心理咨询</router-link></li>
                <li><router-link to="/services">情感咨询</router-link></li>
                <li><router-link to="/services">职业规划</router-link></li>
                <li><router-link to="/services">家庭治疗</router-link></li>
                <li><router-link to="/services">情绪管理</router-link></li>
              </ul>
            </div>
          </el-col>
          <el-col :xs="24" :sm="12" :md="5">
            <div class="footer-section">
              <h3>帮助中心</h3>
              <ul class="footer-links">
                <li><a href="#" @click.prevent="handleShowInfo('faq')">常见问题</a></li>
                <li><a href="#" @click.prevent="handleShowInfo('guide')">预约指南</a></li>
                <li><a href="#" @click.prevent="handleShowInfo('privacy')">隐私政策</a></li>
                <li><a href="#" @click.prevent="handleShowInfo('terms')">服务条款</a></li>
                <li><a href="#" @click.prevent="handleShowInfo('contact')">联系客服</a></li>
              </ul>
            </div>
          </el-col>
          <el-col :xs="24" :sm="12" :md="3">
            <div class="footer-section">
              <h3>联系我们</h3>
              <div class="contact-info">
                <p class="contact-item">
                  <el-icon><Phone /></el-icon>
                  <span>400-888-9999</span>
                </p>
                <p class="contact-item">
                  <el-icon><Message /></el-icon>
                  <span>service@mychat.com</span>
                </p>
                <p class="contact-item">
                  <el-icon><Location /></el-icon>
                  <span>北京市朝阳区</span>
                </p>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
    <div class="footer-bottom">
      <div class="container">
        <p>&copy; 2024 心理咨询平台 版权所有</p>
        <div class="footer-links-bottom">
          <a href="#" @click.prevent="handleShowInfo('privacy')">隐私政策</a>
          <a href="#" @click.prevent="handleShowInfo('terms')">服务条款</a>
          <a href="#" @click.prevent="handleShowInfo('filing')">备案信息</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { showInfo } from '@/utils/errorHandler'

const router = useRouter()

const handleShowInfo = (type) => {
  if (type === 'privacy') {
    router.push('/privacy')
  } else if (type === 'terms') {
    router.push('/terms')
  } else if (type === 'filing') {
    showInfo('备案信息：京ICP备XXXXXXXX号-X')
  } else if (type === 'faq') {
    router.push('/faq')
  } else if (type === 'guide') {
    router.push('/guide')
  } else if (type === 'contact') {
    showInfo('客服邮箱：service@mychat.com')
  }
}
</script>

<style scoped>
.app-footer {
  background: #2c3e50;
  color: white;
  margin-top: auto;
}

.footer-top {
  padding: 60px 0 40px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.footer-section h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 24px;
  color: white;
}

.footer-section .logo {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.footer-section .logo-name {
  font-size: 20px;
  font-weight: bold;
  color: white;
}

.company-desc {
  color: #95a5a6;
  line-height: 1.8;
  margin-bottom: 20px;
  font-size: 14px;
}

.social-links {
  display: flex;
  gap: 12px;
}

.social-icon {
  font-size: 24px;
  color: #95a5a6;
  cursor: pointer;
  transition: color 0.3s;
}

.social-icon:hover {
  color: #409eff;
}

.footer-links {
  list-style: none;
  padding: 0;
  margin: 0;
}

.footer-links li {
  margin-bottom: 12px;
}

.footer-links a {
  color: #95a5a6;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.3s;
}

.footer-links a:hover {
  color: #409eff;
}

.contact-info {
  color: #95a5a6;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-size: 14px;
}

.contact-item .el-icon {
  color: #409eff;
}

.footer-bottom {
  background: #1a252f;
  padding: 20px 0;
}

.footer-bottom .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.footer-bottom p {
  margin: 0;
  color: #7f8c8d;
  font-size: 14px;
}

.footer-links-bottom {
  display: flex;
  gap: 24px;
}

.footer-links-bottom a {
  color: #7f8c8d;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.3s;
}

.footer-links-bottom a:hover {
  color: #409eff;
}

@media (max-width: 768px) {
  .footer-top {
    padding: 40px 0 20px;
  }

  .footer-bottom .container {
    flex-direction: column;
    text-align: center;
  }

  .footer-links-bottom {
    flex-wrap: wrap;
    justify-content: center;
    gap: 16px;
  }
}
</style>
