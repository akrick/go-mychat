# Web 前端功能完善说明

## 已修复的问题

### 1. 变量名拼写错误
- `CounselorDetail.vue`: `submitting` 变量名修正
- `Profile.vue`: `submitting` 变量名修正

### 2. 个人中心功能完善
- 新增 `api/profile.js` 文件，包含：
  - `getProfile()` - 获取个人信息
  - `updateProfile()` - 更新个人资料
  - `changePassword()` - 修改密码
  - `uploadAvatar()` - 上传头像

- `Profile.vue` 功能增强：
  - 头像上传功能（支持图片预览和拖拽上传）
  - 个人信息修改（邮箱、手机号）
  - 密码修改（旧密码验证、新密码确认）
  - 头像悬停显示更换按钮

### 3. Chat 页面优化
- 消息加载错误处理优化
- 发送消息添加空内容验证
- 倒计时计算优化
- 会话信息获取改进

### 4. 首页登录状态检查
- 添加 `userStore` 引入
- 咨询师卡片点击前检查登录状态
- 预约按钮点击前检查登录状态
- 退出登录后跳转到首页而非登录页

### 5. 错误处理工具
- 新增 `utils/errorHandler.js` 统一错误处理：
  - `handleError()` - 通用错误处理
  - `handleValidationError()` - 表单验证错误
  - `handleNetworkError()` - 网络错误
  - `showSuccess()` - 成功提示
  - `showWarning()` - 警告提示
  - `showInfo()` - 信息提示

### 6. API 接口对接状态

#### 认证相关 (`/api/auth.js`)
✅ `login(data)` - 用户登录
✅ `register(data)` - 用户注册
✅ `getUserInfo()` - 获取用户信息

#### 咨询师相关 (`/api/counselor.js`)
✅ `getCounselorList(params)` - 获取咨询师列表
✅ `getCounselorDetail(id)` - 获取咨询师详情

#### 订单相关 (`/api/order.js`)
✅ `createOrder(data)` - 创建订单
✅ `getOrderList(params)` - 获取订单列表
✅ `getOrderDetail(id)` - 获取订单详情
✅ `cancelOrder(id)` - 取消订单
✅ `updateOrderStatus(id, status)` - 更新订单状态

#### 聊天相关 (`/api/chat.js`)
✅ `startChatSession(orderId)` - 开始聊天会话
✅ `sendMessage(sessionId, data)` - 发送消息
✅ `getMessages(sessionId, params)` - 获取消息列表
✅ `endChatSession(sessionId)` - 结束会话
✅ `getChatSessions(params)` - 获取会话列表

#### 个人中心 (`/api/profile.js`) - 新增
✅ `getProfile()` - 获取个人信息
✅ `updateProfile(data)` - 更新个人资料
✅ `changePassword(data)` - 修改密码
✅ `uploadAvatar(file)` - 上传头像

## API 接口说明

### 请求格式
所有接口返回统一格式：
```json
{
  "code": 200,
  "msg": "成功",
  "data": {}
}
```

### 错误码
- `200` - 成功
- `400` - 参数错误
- `401` - 未授权（token 无效或过期）
- `403` - 无权操作
- `404` - 资源不存在
- `500` - 服务器错误

### Token 认证
- 登录成功后返回 token
- 所有需要认证的接口需在请求头携带：
  ```
  Authorization: Bearer {token}
  ```

## 功能使用流程

### 1. 用户注册/登录
1. 访问 `/register` 注册账号
2. 访问 `/login` 登录获取 token
3. token 存储在 localStorage

### 2. 浏览咨询师
1. 首页展示热门咨询师
2. 点击"查看全部"进入咨询师列表
3. 支持关键词搜索
4. 分页加载

### 3. 预约咨询
1. 选择咨询师进入详情页
2. 点击"预约咨询"
3. 选择咨询时长
4. 选择预约时间
5. 提交订单

### 4. 管理订单
1. 订单列表展示所有订单
2. 支持状态筛选
3. 待支付订单可支付
4. 已支付订单可进入咨询

### 5. 在线咨询
1. 从订单列表点击"进入咨询"
2. 开始聊天会话
3. 发送/接收消息
4. 实时倒计时显示剩余时间
5. 时间到自动结束会话

### 6. 个人中心
1. 查看个人信息
2. 上传/更换头像
3. 修改邮箱和手机号
4. 修改登录密码

## 注意事项

1. **登录状态检查**
   - 首页访问咨询师、预约等需要先登录
   - 未登录时自动跳转到登录页

2. **错误处理**
   - 使用统一的错误处理函数
   - 401 错误自动退出登录
   - 网络错误友好提示

3. **数据获取**
   - 所有 API 调用都添加 try-catch
   - 加载状态显示 loading
   - 空数据处理

4. **表单验证**
   - 使用 Element Plus 表单验证
   - 提交前必填项验证
   - 密码二次确认

## 待完善功能

1. 支付功能（显示"开发中"提示）
2. 关于我们页面
3. 服务介绍页面
4. 消息实时推送（WebSocket）
5. 文件上传功能
6. 评价功能
7. 收藏功能
