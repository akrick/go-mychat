# Web 前端项目说明

## 项目结构

```
src/
├── api/           # API 接口
├── assets/        # 静态资源
├── components/     # 公共组件
│   ├── AppHeader.vue
│   └── AppFooter.vue
├── composables/   # 组合式函数
│   ├── useLoading.js
│   └── usePagination.js
├── constants/     # 常量定义
│   └── index.js
├── layouts/       # 布局组件
│   └── MainLayout.vue
├── router/        # 路由配置
│   └── index.js
├── stores/        # 状态管理
│   └── user.js
├── utils/         # 工具函数
│   ├── request.js
│   ├── errorHandler.js
│   ├── validator.js
│   └── formatter.js
└── views/         # 页面组件
    ├── Home.vue
    ├── HomeNew.vue
    ├── Login.vue
    ├── Register.vue
    ├── FAQ.vue
    ├── PrivacyPolicy.vue
    ├── TermsOfService.vue
    ├── BookingGuide.vue
    └── ...
```

## 常量使用

### 导入常量
```javascript
import {
  ROUTE_NAMES,
  ORDER_STATUS,
  ORDER_STATUS_TEXT,
  ORDER_STATUS_COLOR,
  SESSION_DURATION,
  PAGINATION,
  REGEX,
  MESSAGE_TYPE,
  SUCCESS_MESSAGES,
  ERROR_MESSAGES,
  TIME_FORMAT
} from '@/constants'
```

### 使用示例

#### 路由名称
```javascript
router.push({ name: ROUTE_NAMES.HOME })
```

#### 订单状态
```javascript
const statusText = ORDER_STATUS_TEXT[order.status]
const statusColor = ORDER_STATUS_COLOR[order.status]
```

#### 会话时长
```javascript
<el-option
  v-for="duration in Object.values(SESSION_DURATION)"
  :key="duration.value"
  :label="duration.label"
  :value="duration.value"
/>
```

#### 正则验证
```javascript
if (!REGEX.PHONE.test(phone)) {
  ElMessage.error('手机号格式不正确')
}
```

#### 成功/错误消息
```javascript
ElMessage.success(SUCCESS_MESSAGES.LOGIN)
ElMessage.error(ERROR_MESSAGES.NETWORK_ERROR)
```

## Composables 使用

### useLoading - 加载状态管理

```javascript
import { useLoading } from '@/composables/useLoading'

const { loading, error, executeAsync } = useLoading()

async function handleSave() {
  await executeAsync(async () => {
    await saveData()
  })
}
```

### usePagination - 分页管理

```javascript
import { usePagination } from '@/composables/usePagination'

const {
  pagination,
  total,
  loading,
  hasData,
  changePage,
  changePageSize,
  resetPage
} = usePagination({
  pageSize: 20,
  loadData: loadList
})
```

## 工具函数使用

### validator.js - 验证函数

```javascript
import {
  validatePhone,
  validateEmail,
  validatePassword,
  validateImageFile,
  validateConfirmPassword
} from '@/utils/validator'

// 验证手机号
if (!validatePhone(phone)) {
  ElMessage.error('手机号格式不正确')
}

// 验证图片文件
const result = validateImageFile(file)
if (!result.valid) {
  ElMessage.error(result.message)
}
```

### formatter.js - 格式化函数

```javascript
import {
  formatDate,
  formatDateTime,
  formatCurrency,
  formatSeconds,
  formatDuration,
  formatRelativeTime,
  maskPhone,
  maskIdCard,
  copyToClipboard,
  downloadFile
} from '@/utils/formatter'

// 格式化日期
const dateStr = formatDate(date, 'YYYY-MM-DD')

// 格式化金额
const price = formatCurrency(1234.56, 2, '¥') // ¥1234.56

// 格式化秒数
const time = formatSeconds(125) // 02:05

// 隐藏手机号
const masked = maskPhone('13800138000') // 138****8000
```

### errorHandler.js - 错误处理

```javascript
import {
  handleError,
  handleNetworkError,
  handleValidationError,
  showSuccess,
  showWarning,
  showInfo
} from '@/utils/errorHandler'

// 统一错误处理
try {
  await someApi()
} catch (error) {
  handleError(error)
}

// 显示成功消息
showSuccess('操作成功')
```

## 页面列表

### 公开页面
- `/` - 首页
- `/login` - 登录
- `/register` - 注册
- `/services` - 服务介绍
- `/about` - 关于我们
- `/faq` - 常见问题
- `/privacy` - 隐私政策
- `/terms` - 服务条款
- `/guide` - 预约指南

### 需要登录的页面
- `/counselors` - 咨询师列表
- `/counselor/:id` - 咨询师详情
- `/orders` - 我的订单
- `/chat/:sessionId` - 咨询会话
- `/profile` - 个人中心

## 统一布局

所有页面统一使用 `AppHeader` 和 `AppFooter` 组件：

```vue
<template>
  <div class="page-name">
    <AppHeader />
    <div class="page-content">
      <!-- 页面内容 -->
    </div>
    <AppFooter />
  </div>
</template>

<script setup>
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
</script>

<style scoped>
.page-name {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.page-content {
  flex: 1;
  padding: 20px 0;
}
</style>
```

## 路由守卫

路由守卫已配置为：
1. 需要认证的页面会检查 token
2. 如果 token 无效，会自动重新获取用户信息
3. 如果获取失败，自动退出登录并跳转到登录页
4. 已登录用户访问登录/注册页面会跳转到首页

## 状态管理

### userStore
```javascript
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

// 获取信息
const token = userStore.token
const userInfo = userStore.userInfo

// 设置信息
userStore.setToken(token)
userStore.setUserInfo(userInfo)

// 获取用户信息
await userStore.fetchUserInfo()

// 退出登录
userStore.logout()
```

## API 调用

所有 API 统一在 `api/` 目录下管理：

```javascript
import { login, register, getUserInfo } from '@/api/auth'
import { getCounselorList, getCounselorDetail } from '@/api/counselor'
import { getOrderList, createOrder, cancelOrder } from '@/api/order'
import { startChatSession, sendMessage, getMessages } from '@/api/chat'
import { updateProfile, changePassword, uploadAvatar } from '@/api/profile'
```

## 样式规范

### 颜色
- 主色：`#409eff`
- 渐变：`linear-gradient(135deg, #667eea 0%, #764ba2 100%)`
- 成功：`#67c23a`
- 警告：`#e6a23c`
- 危险：`#f56c6c`
- 信息：`#909399`

### 间距
- 小间距：8px
- 中间距：16px
- 大间距：24px
- 超大间距：32px

### 圆角
- 小圆角：4px
- 中圆角：8px
- 大圆角：12px

## 响应式断点

- xs: < 768px
- sm: >= 768px
- md: >= 992px
- lg: >= 1200px
- xl: >= 1920px

## 注意事项

1. **所有页面都必须使用 AppHeader 和 AppFooter**
2. **使用常量替代魔法数字和硬编码字符串**
3. **使用 Composables 复用逻辑**
4. **统一使用工具函数进行格式化和验证**
5. **错误处理统一使用 errorHandler.js 中的函数**
6. **组件命名使用 PascalCase**
7. **文件命名使用 PascalCase（组件）或 camelCase（工具函数）**

## 待完善功能

- [ ] WebSocket 实时通信
- [ ] 支付功能集成
- [ ] 文件上传功能
- [ ] 评价系统
- [ ] 通知系统
- [ ] 收藏功能
- [ ] 搜索功能
