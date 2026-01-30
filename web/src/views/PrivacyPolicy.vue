<template>
  <div class="privacy-page">
    <AppHeader />
    <div class="privacy">
      <div class="page-header">
        <div class="header-content">
          <h1>隐私政策</h1>
          <p class="header-desc">保护您的个人信息安全</p>
        </div>
      </div>

      <el-card class="privacy-card">
        <div class="privacy-content">
          <div class="update-time">
            <el-icon><Clock /></el-icon>
            最后更新时间：2026年1月29日
          </div>

          <section>
            <h2>1. 引言</h2>
            <p>
              本隐私政策旨在向您说明我们如何收集、使用、存储和保护您的个人信息。
              使用本平台服务即表示您同意本隐私政策的条款。
            </p>
          </section>

          <section>
            <h2>2. 信息收集</h2>
            <h3>2.1 我们收集的信息包括：</h3>
            <ul>
              <li>账户信息：用户名、密码、手机号、邮箱等</li>
              <li>个人资料：头像、昵称等您主动提供的信息</li>
              <li>咨询记录：您与咨询师的聊天内容和咨询记录</li>
              <li>订单信息：您的预约和支付记录</li>
              <li>设备信息：访问平台的设备类型、浏览器信息等</li>
              <li>日志信息：访问时间、访问页面、操作记录等</li>
            </ul>
          </section>

          <section>
            <h2>3. 信息使用</h2>
            <p>我们使用收集的信息用于：</p>
            <ul>
              <li>提供和改进我们的服务</li>
              <li>处理您的预约和支付订单</li>
              <li>与您沟通，提供客户服务</li>
              <li>维护平台安全，防止欺诈行为</li>
              <li>遵守法律法规要求</li>
            </ul>
          </section>

          <section>
            <h2>4. 信息共享</h2>
            <p>我们不会在未经您同意的情况下出售、出租或共享您的个人信息，除非：</p>
            <ul>
              <li>获得您的明确同意</li>
              <li>法律法规要求或政府部门要求</li>
              <li>为保护我们的权利和财产</li>
              <li>与可信赖的第三方服务提供商共享（仅限于提供服务所需）</li>
            </ul>
          </section>

          <section>
            <h2>5. 信息存储</h2>
            <p>
              我们采用行业标准的安全措施来保护您的信息，包括：
            </p>
            <ul>
              <li>加密存储敏感信息</li>
              <li>限制对个人信息的访问权限</li>
              <li>定期进行安全审查和漏洞扫描</li>
              <li>安全的数据传输协议</li>
            </ul>
          </section>

          <section>
            <h2>6. Cookie 使用</h2>
            <p>
              我们使用 Cookie 来改善用户体验，分析网站流量。
              您可以通过浏览器设置禁用 Cookie，但这可能会影响某些功能的使用。
            </p>
          </section>

          <section>
            <h2>7. 您的权利</h2>
            <p>您对自己的个人信息享有以下权利：</p>
            <ul>
              <li>访问权：查看我们收集的关于您的信息</li>
              <li>更正权：要求更正不准确的信息</li>
              <li>删除权：要求删除您的个人信息</li>
              <li>撤回同意：随时撤回您对信息使用的同意</li>
              <li>数据携带：以结构化格式导出您的数据</li>
            </ul>
          </section>

          <section>
            <h2>8. 咨询保密</h2>
            <p>
              您的咨询内容将严格保密，仅咨询师和授权人员可以访问。
              我们承诺不会将您的咨询内容用于商业目的或向第三方披露。
            </p>
          </section>

          <section>
            <h2>9. 儿童隐私</h2>
            <p>
              本平台不面向16岁以下的儿童。如果我们发现收集了未成年人的信息，
              将会立即删除。
            </p>
          </section>

          <section>
            <h2>10. 政策变更</h2>
            <p>
              我们可能会不时更新本隐私政策。
              重大变更时，我们会通过网站公告或邮件通知您。
              建议您定期查看本页面以了解最新信息。
            </p>
          </section>

          <section>
            <h2>11. 联系我们</h2>
            <p>如果您对本隐私政策有任何疑问或意见，请通过以下方式联系我们：</p>
            <div class="contact-info">
              <p><el-icon><Message /></el-icon> 邮箱：support@mychat.com</p>
              <p><el-icon><Phone /></el-icon> 电话：400-888-9999</p>
              <p><el-icon><Location /></el-icon> 地址：中国 北京市 朝阳区</p>
            </div>
          </section>
        </div>
      </el-card>
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import { Clock, Message, Phone, Location } from '@element-plus/icons-vue'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
</script>

<style scoped>
.privacy-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.privacy {
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

.privacy-card {
  max-width: 900px;
  margin: 0 auto;
}

.update-time {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 32px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  color: #666;
}

.privacy-content section {
  margin-bottom: 32px;
}

.privacy-content h2 {
  color: #333;
  font-size: 24px;
  margin-bottom: 16px;
  border-bottom: 2px solid #667eea;
  padding-bottom: 8px;
}

.privacy-content h3 {
  color: #555;
  font-size: 18px;
  margin: 20px 0 12px 0;
}

.privacy-content p {
  color: #666;
  line-height: 1.8;
  margin-bottom: 16px;
}

.privacy-content ul {
  list-style: none;
  padding: 0;
}

.privacy-content ul li {
  color: #666;
  line-height: 1.8;
  padding-left: 24px;
  position: relative;
  margin-bottom: 8px;
}

.privacy-content ul li::before {
  content: '•';
  color: #667eea;
  font-weight: bold;
  position: absolute;
  left: 0;
}

.contact-info {
  background: #f5f7fa;
  padding: 24px;
  border-radius: 8px;
  margin-top: 16px;
}

.contact-info p {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.contact-info p:last-child {
  margin-bottom: 0;
}

@media (max-width: 768px) {
  .page-header {
    padding: 24px;
  }

  .page-header h1 {
    font-size: 24px;
  }

  .privacy-content h2 {
    font-size: 20px;
  }
}
</style>
