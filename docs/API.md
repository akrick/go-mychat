# MyChat API 文档

## 基础信息

### 管理后台 API
- 基础地址: `http://localhost:8081`
- 认证方式: JWT Bearer Token
- 响应格式: JSON

### 用户端 API
- 基础地址: `http://localhost:3000`
- WebSocket: `ws://localhost:8082/ws`

### 通用响应格式

**成功响应**:
```json
{
  "code": 200,
  "msg": "操作成功",
  "data": {}
}
```

**失败响应**:
```json
{
  "code": 400,
  "msg": "操作失败"
}
```

## 管理后台 API

### 认证接口

#### 登录
```http
POST /api/admin/login
Content-Type: application/json

{
  "username": "admin",
  "password": "123456"
}
```

#### 退出登录
```http
POST /api/admin/logout
Authorization: Bearer {token}
```

### 用户管理

#### 获取用户列表
```http
GET /api/admin/users?page=1&page_size=20&keyword=&status=
Authorization: Bearer {token}
```

#### 创建用户
```http
POST /api/admin/users
Authorization: Bearer {token}
Content-Type: application/json

{
  "username": "test",
  "password": "123456",
  "email": "test@example.com",
  "phone": "13800138000"
}
```

#### 更新用户
```http
PUT /api/admin/users/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "email": "new@example.com",
  "phone": "13900139000"
}
```

#### 删除用户
```http
DELETE /api/admin/users/{id}
Authorization: Bearer {token}
```

#### 重置用户密码
```http
POST /api/admin/users/{id}/password
Authorization: Bearer {token}
Content-Type: application/json

{
  "password": "newpassword123"
}
```

### 咨询师管理

#### 获取咨询师列表
```http
GET /api/admin/counselors?page=1&page_size=20&keyword=&status=
Authorization: Bearer {token}
```

#### 创建咨询师
```http
POST /api/admin/counselors
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "张三",
  "avatar": "/uploads/avatars/1.jpg",
  "specialties": "心理咨询,情绪管理",
  "price_per_minute": 5.0,
  "status": 1
}
```

#### 更新咨询师
```http
PUT /api/admin/counselors/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "李四",
  "specialties": "心理咨询,情绪管理,职业规划"
}
```

#### 删除咨询师
```http
DELETE /api/admin/counselors/{id}
Authorization: Bearer {token}
```

### 订单管理

#### 获取订单列表
```http
GET /api/admin/orders?page=1&page_size=20&status=&user_id=&counselor_id=
Authorization: Bearer {token}
```

#### 获取订单统计
```http
GET /api/admin/orders/statistics
Authorization: Bearer {token}
```

#### 更新订单状态
```http
PUT /api/admin/orders/{id}/status
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": 2
}
```

### 聊天管理

#### 获取聊天会话列表
```http
GET /api/admin/chat/sessions?page=1&page_size=20
Authorization: Bearer {token}
```

#### 获取聊天消息
```http
GET /api/admin/chat/sessions/{session_id}/messages?page=1&page_size=50
Authorization: Bearer {token}
```

#### 获取聊天统计
```http
GET /api/admin/chat/statistics
Authorization: Bearer {token}
```

#### 搜索聊天消息
```http
GET /api/admin/chat/messages/search?keyword=测试&page=1&page_size=20
Authorization: Bearer {token}
```

#### 删除聊天会话
```http
DELETE /api/admin/chat/sessions/{id}
Authorization: Bearer {token}
```

### 财务管理

#### 获取待审核提现
```http
GET /api/admin/withdraws/pending
Authorization: Bearer {token}
```

#### 审核提现
```http
POST /api/admin/withdraw/{id}/approve
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": 1,
  "remark": "审核通过"
}
```

#### 确认打款
```http
POST /api/admin/withdraw/{id}/transfer
Authorization: Bearer {token}
```

#### 获取提现列表
```http
GET /api/admin/withdraws?page=1&page_size=20&status=
Authorization: Bearer {token}
```

#### 获取财务统计
```http
GET /api/admin/finance/stats?start_date=&end_date=
Authorization: Bearer {token}
```

#### 获取营收报表
```http
GET /api/admin/finance/revenue?type=daily&start_date=&end_date=
Authorization: Bearer {token}
```

#### 获取财务报表
```http
GET /api/admin/finance/reports?type=revenue&start_date=&end_date=
Authorization: Bearer {token}
```

#### 获取咨询师账户列表
```http
GET /api/admin/finance/accounts?page=1&page_size=20&keyword=
Authorization: Bearer {token}
```

#### 获取咨询师账户详情
```http
GET /api/admin/finance/accounts/{id}
Authorization: Bearer {token}
```

### 系统管理

#### 获取用户信息
```http
GET /api/admin/user/info
Authorization: Bearer {token}
```

#### 获取用户权限
```http
GET /api/admin/user/permissions
Authorization: Bearer {token}
```

#### 获取会话统计
```http
GET /api/admin/session/stats
Authorization: Bearer {token}
```

#### 获取在线用户
```http
GET /api/admin/online/users
Authorization: Bearer {token}
```

#### 广播系统消息
```http
POST /api/admin/broadcast
Authorization: Bearer {token}
Content-Type: application/json

{
  "message": "系统维护通知",
  "type": "system"
}
```

#### 获取系统日志
```http
GET /api/admin/logs?page=1&page_size=20&level=&module=&start_date=&end_date=
Authorization: Bearer {token}
```

#### 获取系统配置
```http
GET /api/admin/configs?category=
Authorization: Bearer {token}
```

#### 创建系统配置
```http
POST /api/admin/configs
Authorization: Bearer {token}
Content-Type: application/json

{
  "key": "custom_config",
  "value": "custom_value",
  "category": "basic",
  "label": "自定义配置",
  "type": "string"
}
```

#### 更新系统配置
```http
PUT /api/admin/configs/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "value": "new_value"
}
```

#### 删除系统配置
```http
DELETE /api/admin/configs/{id}
Authorization: Bearer {token}
```

### 权限管理 (RBAC)

#### 获取角色列表
```http
GET /api/admin/roles?page=1&page_size=20&keyword=
Authorization: Bearer {token}
```

#### 创建角色
```http
POST /api/admin/roles
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "财务经理",
  "description": "财务管理权限",
  "status": 1
}
```

#### 更新角色
```http
PUT /api/admin/roles/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "高级财务经理",
  "description": "高级财务管理权限"
}
```

#### 删除角色
```http
DELETE /api/admin/roles/{id}
Authorization: Bearer {token}
```

#### 获取角色权限
```http
GET /api/admin/roles/{id}/permissions
Authorization: Bearer {token}
```

#### 分配权限
```http
PUT /api/admin/roles/{id}/permissions
Authorization: Bearer {token}
Content-Type: application/json

{
  "permissions": [1, 2, 3, 4, 5]
}
```

#### 获取权限树
```http
GET /api/admin/permissions/tree
Authorization: Bearer {token}
```

#### 获取权限列表
```http
GET /api/admin/permissions?page=1&page_size=20&keyword=
Authorization: Bearer {token}
```

#### 创建权限
```http
POST /api/admin/permissions
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "财务报表",
  "code": "finance:reports:view",
  "type": "menu",
  "parent_id": 0,
  "sort": 1
}
```

#### 更新权限
```http
PUT /api/admin/permissions/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "高级财务报表"
}
```

#### 删除权限
```http
DELETE /api/admin/permissions/{id}
Authorization: Bearer {token}
```

### 菜单管理

#### 获取菜单树
```http
GET /api/admin/menus/tree
Authorization: Bearer {token}
```

#### 获取菜单列表
```http
GET /api/admin/menus?page=1&page_size=20&keyword=
Authorization: Bearer {token}
```

#### 创建菜单
```http
POST /api/admin/menus
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "财务管理",
  "path": "/finance",
  "icon": "Wallet",
  "parent_id": 0,
  "sort": 1,
  "status": 1
}
```

#### 更新菜单
```http
PUT /api/admin/menus/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "高级财务管理",
  "path": "/finance/advanced"
}
```

#### 删除菜单
```http
DELETE /api/admin/menus/{id}
Authorization: Bearer {token}
```

### 低代码平台

#### 获取表单列表
```http
GET /api/admin/lowcode/forms?page=1&page_size=20&keyword=
Authorization: Bearer {token}
```

#### 保存表单
```http
POST /api/admin/lowcode/forms
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "用户反馈表单",
  "description": "收集用户反馈信息",
  "schema": {
    "fields": [...]
  }
}
```

#### 更新表单
```http
PUT /api/admin/lowcode/forms/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "用户反馈表单（更新版）"
}
```

#### 删除表单
```http
DELETE /api/admin/lowcode/forms/{id}
Authorization: Bearer {token}
```

#### 获取表单详情
```http
GET /api/admin/lowcode/forms/{id}
Authorization: Bearer {token}
```

#### 获取表单数据
```http
GET /api/admin/lowcode/forms/{id}/data?page=1&page_size=20
Authorization: Bearer {token}
```

#### 提交表单数据
```http
POST /api/admin/lowcode/forms/{id}/submit
Authorization: Bearer {token}
Content-Type: application/json

{
  "data": {
    "field1": "value1",
    "field2": "value2"
  }
}
```

#### 获取页面列表
```http
GET /api/admin/lowcode/pages?page=1&page_size=20&keyword=
Authorization: Bearer {token}
```

#### 保存页面
```http
POST /api/admin/lowcode/pages
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "用户中心页面",
  "description": "用户个人中心",
  "schema": {
    "components": [...]
  }
}
```

#### 更新页面
```http
PUT /api/admin/lowcode/pages/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "高级用户中心"
}
```

#### 删除页面
```http
DELETE /api/admin/lowcode/pages/{id}
Authorization: Bearer {token}
```

#### 获取页面详情
```http
GET /api/admin/lowcode/pages/{id}
Authorization: Bearer {token}
```

#### 页面预览
```http
GET /api/admin/lowcode/pages/{id}/preview
Authorization: Bearer {token}
```

### 统计数据

#### 获取统计数据
```http
GET /api/admin/statistics
Authorization: Bearer {token}
```

## 用户端 API

### 公开接口

#### 获取咨询师列表
```http
GET /api/counselors?page=1&page_size=20&keyword=&specialty=
```

#### 获取咨询师详情
```http
GET /api/counselors/{id}
```

### 需要认证的接口

#### 用户注册
```http
POST /api/user/register
Content-Type: application/json

{
  "username": "test",
  "password": "123456",
  "email": "test@example.com",
  "phone": "13800138000"
}
```

#### 用户登录
```http
POST /api/user/login
Content-Type: application/json

{
  "username": "test",
  "password": "123456"
}
```

#### 创建订单
```http
POST /api/orders
Authorization: Bearer {token}
Content-Type: application/json

{
  "counselor_id": 1,
  "type": "text",
  "amount": 100
}
```

#### 支付订单
```http
POST /api/payment/create
Authorization: Bearer {token}
Content-Type: application/json

{
  "order_id": 1,
  "payment_method": "wechat"
}
```

#### WebSocket 连接
```
ws://localhost:8082/ws?token={token}&user_id={user_id}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 操作成功 |
| 400 | 参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |

## 状态码说明

### 订单状态
| 状态值 | 说明 |
|--------|------|
| 0 | 待支付 |
| 1 | 已支付 |
| 2 | 进行中 |
| 3 | 已完成 |
| 4 | 已取消 |

### 提现状态
| 状态值 | 说明 |
|--------|------|
| 0 | 待审核 |
| 1 | 已通过 |
| 2 | 已拒绝 |
| 3 | 已打款 |

### 聊天会话状态
| 状态值 | 说明 |
|--------|------|
| 0 | 待开始 |
| 1 | 进行中 |
| 2 | 已结束 |

## 数据分页

所有列表接口支持分页参数：

- `page`: 页码，从 1 开始，默认 1
- `page_size`: 每页数量，默认 20

响应格式：
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

## WebSocket 消息格式

### 客户端发送消息
```json
{
  "type": "chat",
  "data": {
    "session_id": 123,
    "content": "你好"
  }
}
```

### 服务端推送消息
```json
{
  "type": "message",
  "data": {
    "id": 456,
    "session_id": 123,
    "sender_id": 1,
    "content": "你好",
    "timestamp": 1640995200000
  }
}
```

### 消息类型
| 类型 | 说明 |
|------|------|
| chat | 聊天消息 |
| system | 系统消息 |
| notification | 通知消息 |
| typing | 正在输入 |
| online_status | 在线状态 |

## 注意事项

1. 所有需要认证的接口都必须在请求头中携带 JWT Token
2. Token 格式: `Authorization: Bearer {token}`
3. Token 有效期: 24 小时
4. 文件上传大小限制: 10MB
5. 密码最少 6 位字符
6. 用户名唯一性检查
7. 所有时间格式使用 ISO 8601 标准

## 更新日志

### v1.0.0 (2026-01-28)
- 初始版本发布
- 完成所有核心功能模块
