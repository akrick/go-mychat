<template>
  <div class="guide-page">
    <AppHeader />
    <div class="guide">
      <div class="page-header">
        <div class="header-content">
          <h1>预约指南</h1>
          <p class="header-desc">轻松完成咨询师预约，开始您的心理成长之旅</p>
        </div>
      </div>

      <el-timeline class="guide-timeline">
        <el-timeline-item timestamp="步骤 1" placement="top">
          <el-card>
            <template #header>
              <div class="step-header">
                <el-icon><User /></el-icon>
                <h3>注册登录</h3>
              </div>
            </template>
            <div class="step-content">
              <p>首次使用需要注册账号</p>
              <el-steps :active="1" simple class="mini-steps">
                <el-step title="点击注册" />
                <el-step title="填写信息" />
                <el-step title="完成注册" />
              </el-steps>
              <div class="step-tips">
                <el-icon><InfoFilled /></el-icon>
                <span>提示：请填写真实信息，便于后续联系和服务</span>
              </div>
            </div>
          </el-card>
        </el-timeline-item>

        <el-timeline-item timestamp="步骤 2" placement="top">
          <el-card>
            <template #header>
              <div class="step-header">
                <el-icon><Search /></el-icon>
                <h3>浏览咨询师</h3>
              </div>
            </template>
            <div class="step-content">
              <p>查看咨询师列表，选择合适的咨询师</p>
              <ul class="feature-list">
                <li><el-icon><Check /></el-icon> 查看咨询师专业背景和经验</li>
                <li><el-icon><Check /></el-icon> 了解咨询师擅长领域</li>
                <li><el-icon><Check /></el-icon> 参考其他用户的评价</li>
                <li><el-icon><Check /></el-icon> 对比咨询收费标准</li>
              </ul>
              <el-button type="primary" @click="$router.push('/counselors')">
                前往咨询师列表
              </el-button>
            </div>
          </el-card>
        </el-timeline-item>

        <el-timeline-item timestamp="步骤 3" placement="top">
          <el-card>
            <template #header>
              <div class="step-header">
                <el-icon><User /></el-icon>
                <h3>查看详情</h3>
              </div>
            </template>
            <div class="step-content">
              <p>点击咨询师进入详情页，了解更多信息</p>
              <el-descriptions :column="1" border>
                <el-descriptions-item label="专业资质">
                  国家二级心理咨询师
                </el-descriptions-item>
                <el-descriptions-item label="从业年限">
                  5年以上
                </el-descriptions-item>
                <el-descriptions-item label="咨询风格">
                  亲和、耐心、专业
                </el-descriptions-item>
                <el-descriptions-item label="收费标准">
                  ¥2-5/分钟
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </el-card>
        </el-timeline-item>

        <el-timeline-item timestamp="步骤 4" placement="top">
          <el-card>
            <template #header>
              <div class="step-header">
                <el-icon><Calendar /></el-icon>
                <h3>选择时间和时长</h3>
              </div>
            </template>
            <div class="step-content">
              <p>点击"预约咨询"按钮，选择咨询时长和预约时间</p>
              <div class="time-options">
                <h4>咨询时长选择：</h4>
                <div class="time-cards">
                  <el-card class="time-card">
                    <div class="time-value">30</div>
                    <div class="time-unit">分钟</div>
                    <div class="time-desc">适合初步沟通</div>
                  </el-card>
                  <el-card class="time-card recommended">
                    <div class="badge">推荐</div>
                    <div class="time-value">60</div>
                    <div class="time-unit">分钟</div>
                    <div class="time-desc">标准咨询时长</div>
                  </el-card>
                  <el-card class="time-card">
                    <div class="time-value">90</div>
                    <div class="time-unit">分钟</div>
                    <div class="time-desc">深度咨询</div>
                  </el-card>
                  <el-card class="time-card">
                    <div class="time-value">120</div>
                    <div class="time-unit">分钟</div>
                    <div class="time-desc">综合治疗</div>
                  </el-card>
                </div>
              </div>
              <div class="step-tips">
                <el-icon><WarningFilled /></el-icon>
                <span>提示：首次咨询建议选择60分钟，便于充分交流</span>
              </div>
            </div>
          </el-card>
        </el-timeline-item>

        <el-timeline-item timestamp="步骤 5" placement="top">
          <el-card>
            <template #header>
              <div class="step-header">
                <el-icon><Wallet /></el-icon>
                <h3>确认并支付</h3>
              </div>
            </template>
            <div class="step-content">
              <p>核对订单信息，完成支付</p>
              <div class="order-preview">
                <div class="preview-item">
                  <span>咨询师：</span>
                  <span>张三博士</span>
                </div>
                <div class="preview-item">
                  <span>咨询时长：</span>
                  <span>60分钟</span>
                </div>
                <div class="preview-item">
                  <span>预约时间：</span>
                  <span>2026-01-30 14:00</span>
                </div>
                <div class="preview-item total">
                  <span>订单金额：</span>
                  <span class="amount">¥300.00</span>
                </div>
              </div>
              <div class="payment-methods">
                <h4>支持的支付方式：</h4>
                <div class="payment-icons">
                  <el-tag type="success">微信支付</el-tag>
                  <el-tag type="primary">支付宝</el-tag>
                </div>
              </div>
            </div>
          </el-card>
        </el-timeline-item>

        <el-timeline-item timestamp="步骤 6" placement="top">
          <el-card>
            <template #header>
              <div class="step-header">
                <el-icon><ChatDotRound /></el-icon>
                <h3>开始咨询</h3>
              </div>
            </template>
            <div class="step-content">
              <p>预约时间开始后，在订单中点击"进入咨询"即可开始</p>
              <div class="step-tips">
                <el-icon><SuccessFilled /></el-icon>
                <span>咨询开始前5分钟，系统会发送提醒通知</span>
              </div>
              <el-alert
                title="咨询注意事项"
                type="info"
                :closable="false"
              >
                <template #default>
                  <ul>
                    <li>确保网络连接稳定</li>
                    <li>选择安静、私密的环境</li>
                    <li>准备好想要讨论的问题</li>
                    <li>保持开放和诚实的态度</li>
                  </ul>
                </template>
              </el-alert>
            </div>
          </el-card>
        </el-timeline-item>
      </el-timeline>

      <el-card class="faq-card">
        <template #header>
          <h3><el-icon><QuestionFilled /></el-icon> 常见问题</h3>
        </template>
        <el-collapse v-model="activeFaq">
          <el-collapse-item title="预约后可以取消吗？" name="1">
            <p>可以。在咨询开始前24小时以上可免费取消，24小时内取消将收取50%费用。</p>
          </el-collapse-item>
          <el-collapse-item title="咨询时间可以延长吗？" name="2">
            <p>可以在订单中选择更长的时长。如需临时延长，可在当前会话结束后再次预约。</p>
          </el-collapse-item>
          <el-collapse-item title="可以更换咨询师吗？" name="3">
            <p>可以。如果与当前咨询师不匹配，可以取消预约并重新选择其他咨询师。</p>
          </el-collapse-item>
          <el-collapse-item title="咨询记录会保存吗？" name="4">
            <p>是的，您的所有咨询记录都会安全保存，您可以在个人中心随时查看。</p>
          </el-collapse-item>
        </el-collapse>
      </el-card>

      <div class="action-section">
        <el-button type="primary" size="large" @click="$router.push('/counselors')">
          <el-icon><Plus /></el-icon>
          立即预约
        </el-button>
        <el-button size="large" @click="$router.push('/faq')">
          <el-icon><QuestionFilled /></el-icon>
          查看更多问题
        </el-button>
      </div>
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { User, InfoFilled, Search, Check, Calendar, WarningFilled, Wallet, Plus, QuestionFilled, SuccessFilled, ChatDotRound } from '@element-plus/icons-vue'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'

const activeFaq = ref([])
</script>

<style scoped>
.guide-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.guide {
  flex: 1;
  padding: 20px 0;
}

.page-header {
  padding: 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  margin-bottom: 24px;
  text-align: center;
}

.page-header h1 {
  margin: 0 0 12px 0;
  font-size: 32px;
  font-weight: 600;
}

.header-desc {
  margin: 0;
  font-size: 16px;
  opacity: 0.9;
}

.guide-timeline {
  max-width: 900px;
  margin: 0 auto 32px;
}

.step-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.step-header h3 {
  margin: 0;
  color: #333;
}

.step-content p {
  color: #666;
  margin-bottom: 16px;
}

.mini-steps {
  margin: 16px 0;
}

.step-tips {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: #ecf5ff;
  border-radius: 4px;
  color: #409eff;
  font-size: 14px;
}

.feature-list {
  list-style: none;
  padding: 0;
  margin-bottom: 16px;
}

.feature-list li {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  color: #666;
}

.feature-list li .el-icon {
  color: #67c23a;
}

.time-options {
  margin: 20px 0;
}

.time-options h4 {
  margin-bottom: 16px;
  color: #333;
}

.time-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.time-card {
  text-align: center;
  padding: 16px;
  position: relative;
  cursor: pointer;
  transition: all 0.3s;
}

.time-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

.time-card.recommended {
  border: 2px solid #409eff;
}

.badge {
  position: absolute;
  top: -10px;
  right: -10px;
  background: #409eff;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.time-value {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 4px;
}

.time-unit {
  color: #999;
  font-size: 14px;
  margin-bottom: 8px;
}

.time-desc {
  color: #666;
  font-size: 13px;
}

.order-preview {
  background: #f5f7fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.preview-item {
  display: flex;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #e4e7ed;
}

.preview-item:last-child {
  border-bottom: none;
}

.preview-item.total {
  font-weight: bold;
  font-size: 18px;
}

.amount {
  color: #f56c6c;
  font-size: 24px;
}

.payment-methods {
  margin-top: 20px;
}

.payment-methods h4 {
  margin-bottom: 12px;
  color: #333;
}

.payment-icons {
  display: flex;
  gap: 12px;
}

.faq-card {
  max-width: 900px;
  margin: 0 auto 32px;
}

.faq-card h3 {
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #333;
}

.faq-card p {
  color: #666;
  line-height: 1.8;
}

.faq-card ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.faq-card ul li {
  padding: 4px 0;
  color: #666;
}

.action-section {
  text-align: center;
  padding: 32px 0;
  max-width: 900px;
  margin: 0 auto;
}

.action-section .el-button {
  margin: 0 8px;
}

@media (max-width: 768px) {
  .page-header {
    padding: 24px;
  }

  .page-header h1 {
    font-size: 24px;
  }

  .time-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .action-section {
    padding: 24px 0;
  }

  .action-section .el-button {
    margin: 8px 4px;
    width: calc(50% - 12px);
  }
}
</style>
