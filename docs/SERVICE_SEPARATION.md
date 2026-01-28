# MyChat 服务分离说明

> 分离日期: 2026-01-28
> 目标: 将用户端 API 和管理后台 API 完全分离

---

## 📋 服务架构概览

MyChat 项目现在包含 **3 个独立的服务**：

| 服务 | 端口 | 职责 | 路径 |
|------|------|------|------|
| **用户端 API** | 8080 | 用户、咨询师、订单、支付、聊天 | `api/` |
| **管理后台 API** | 8081 | 用户管理、角色权限、订单审核、财务管理 | `admin/backend/` |
| **WebSocket 服务** | 8082 | 实时消息推送、在线状态 | `websocket/` |

---

## 🔥 用户端 API (api/) - 端口 8080

### 服务职责
- ✅ 用户注册/登录
- ✅ 用户信息管理
- ✅ 咨询师列表和详情（只读）
- ✅ 订单创建和管理（用户端）
- ✅ 支付处理
- ✅ 聊天会话管理
- ✅ 账单和提现（用户端）
- ✅ 通知管理
- ✅ 文件上传

### 主要接口
```
# 认证
POST   /api/register
POST   /api/login
POST   /api/token/refresh

# 用户
GET    /api/user/info
PUT    /api/user/profile
POST   /api/user/password
POST   /api/upload/avatar

# 咨询师（只读）
GET    /api/counselor/list
GET    /api/counselor/:id

# 订单
POST   /api/order/create
GET    /api/order/:id
GET    /api/order/list
PUT    /api/order/:id/status
POST   /api/order/:id/cancel
GET    /api/counselor/orders
GET    /api/order/:id/validate
GET    /api/order/:id/timeline

# 支付
POST   /api/payment/create
GET    /api/payment/:id
GET    /api/payment/list
POST   /api/payment/refund
POST   /api/payment/wechat/callback
POST   /api/payment/alipay/callback

# 聊天
POST   /api/chat/session/:order_id/start
GET    /api/chat/session/:session_id/messages
POST   /api/chat/session/:session_id/end
GET    /api/chat/sessions

# 账单
GET    /api/chat/billings
GET    /api/chat/counselor/billings
GET    /api/chat/counselor/account
POST   /api/chat/counselor/withdraw
GET    /api/chat/counselor/withdraws

# 通知
GET    /api/notification/list
POST   /api/notification/:id/read
POST   /api/notification/read-all
DELETE /api/notification/:id

# 文件
POST   /api/upload
GET    /api/file/:id
DELETE /api/file/:id

# 健康检查
GET    /health

# 文档
GET    /swagger/*
```

### 目录结构
```
api/
├── main.go              # 主入口
├── handlers/            # 处理器 (9个)
│   ├── auth.go         # 认证
│   ├── user.go         # 用户
│   ├── counselor.go    # 咨询师
│   ├── order.go        # 订单
│   ├── payment.go      # 支付
│   ├── chat.go         # 聊天
│   ├── notification.go # 通知
│   ├── profile.go      # 个人资料
│   └── upload.go       # 上传
├── models/             # 模型 (6个)
├── cache/              # 缓存 (5个)
├── database/           # 数据库 (1个)
├── middleware/         # 中间件 (3个)
└── utils/              # 工具 (5个)
```

---

## 🔐 管理后台 API (admin/backend/) - 端口 8081

### 服务职责
- ✅ 管理员登录/登出
- ✅ 用户管理（CRUD、密码重置）
- ✅ 咨询师管理（CRUD、账户管理）
- ✅ 订单管理（审核、状态更新）
- ✅ 聊天记录管理
- ✅ 财务管理（提现审核、打款、统计）
- ✅ 角色权限管理（RBAC）
- ✅ 菜单管理
- ✅ 系统配置管理
- ✅ 系统日志管理
- ✅ 低代码平台

### 主要接口
```
# 认证
POST   /api/admin/login
POST   /api/admin/logout
GET    /api/admin/user/info
GET    /api/admin/user/permissions

# 用户管理
GET    /api/admin/users
POST   /api/admin/users
PUT    /api/admin/users/:id
DELETE /api/admin/users/:id
POST   /api/admin/users/:id/password

# 咨询师管理
GET    /api/admin/counselors
POST   /api/admin/counselors
PUT    /api/admin/counselors/:id
DELETE /api/admin/counselors/:id

# 订单管理
GET    /api/admin/orders
GET    /api/admin/orders/statistics
PUT    /api/admin/orders/:id/status

# 聊天管理
GET    /api/admin/chat/sessions
GET    /api/admin/chat/sessions/:session_id/messages
GET    /api/admin/chat/statistics
GET    /api/admin/chat/messages/search
DELETE /api/admin/chat/sessions/:id

# 财务管理
GET    /api/admin/withdraws/pending
POST   /api/admin/withdraw/:id/approve
POST   /api/admin/withdraw/:id/transfer
GET    /api/admin/withdraws
GET    /api/admin/finance/stats
GET    /api/admin/finance/revenue
GET    /api/admin/finance/reports
GET    /api/admin/finance/accounts
GET    /api/admin/finance/accounts/:id

# 系统管理
GET    /api/admin/statistics
GET    /api/admin/session/stats
GET    /api/admin/online/users
POST   /api/admin/broadcast
GET    /api/admin/logs
GET    /api/admin/configs
POST   /api/admin/configs
PUT    /api/admin/configs/:id
DELETE /api/admin/configs/:id

# RBAC
GET    /api/admin/roles
POST   /api/admin/roles
PUT    /api/admin/roles/:id
DELETE /api/admin/roles/:id
GET    /api/admin/roles/:id/permissions
PUT    /api/admin/roles/:id/permissions

GET    /api/admin/permissions/tree
GET    /api/admin/permissions
POST   /api/admin/permissions
PUT    /api/admin/permissions/:id
DELETE /api/admin/permissions/:id

# 菜单管理
GET    /api/admin/menus/tree
GET    /api/admin/menus
POST   /api/admin/menus
PUT    /api/admin/menus/:id
DELETE /api/admin/menus/:id

# 低代码平台
GET    /api/admin/lowcode/forms
POST   /api/admin/lowcode/forms
GET    /api/admin/lowcode/forms/:id
DELETE /api/admin/lowcode/forms/:id
GET    /api/admin/lowcode/forms/:id/data
POST   /api/admin/lowcode/forms/:id/submit

GET    /api/admin/lowcode/pages
POST   /api/admin/lowcode/pages
GET    /api/admin/lowcode/pages/:id
DELETE /api/admin/lowcode/pages/:id
GET    /api/admin/lowcode/pages/:id/preview

# 测试
GET    /api/test
```

### 目录结构
```
admin/backend/
├── main.go              # 主入口
├── init_system.go       # 系统初始化
├── handlers/            # 处理器 (23个)
│   ├── admin_chat.go   # 管理聊天
│   ├── admin_order.go  # 管理订单
│   ├── admin.go        # 管理员
│   ├── auth.go         # 认证
│   ├── chat.go         # 聊天
│   ├── config.go       # 配置
│   ├── counselor.go    # 咨询师
│   ├── finance.go      # 财务
│   ├── lowcode.go      # 低代码
│   ├── menu.go         # 菜单
│   ├── notification.go # 通知
│   ├── order.go        # 订单
│   ├── payment.go      # 支付
│   ├── rbac.go         # 权限
│   ├── stats.go        # 统计
│   ├── system.go       # 系统
│   ├── upload.go       # 上传
│   └── user.go         # 用户
├── models/             # 模型 (11个)
├── cache/              # 缓存 (5个)
├── database/           # 数据库 (1个)
├── middleware/         # 中间件 (1个)
├── websocket/          # WebSocket (4个)
│   ├── hub.go
│   ├── manager.go
│   ├── message.go
│   └── stats.go
└── utils/              # 工具
```

---

## 🔄 WebSocket 服务 (websocket/) - 端口 8082

### 服务职责
- ✅ WebSocket 连接管理
- ✅ 实时消息推送
- ✅ 在线用户状态
- ✅ 会话管理
- ✅ 消息广播

### 功能
- 聊天消息实时推送
- 在线状态同步
- 系统消息广播
- 会话统计

---

## 🏗️ 服务依赖关系

```
┌─────────────────────────────────────────────────────────┐
│                   用户端应用                            │
│  (Vue/React/小程序/H5)                                │
└────────────────────┬────────────────────────────────────┘
                     │
    ┌────────────────┴────────────────┐
    │                                 │
    ↓ HTTP REST API                  ↓ WebSocket
┌──────────────────┐         ┌──────────────────┐
│ 用户端 API (8080)│         │WebSocket (8082)  │
│ • 用户/咨询师     │         │ • 实时消息       │
│ • 订单/支付       │         │ • 在线状态       │
│ • 聊天会话        │         └────────┬─────────┘
└────────┬─────────┘                  │
         │                            │
         └──────────┬─────────────────┘
                    │
                    ↓ 共享
    ┌──────────────────────────────┐
    │        共享资源             │
    │  • MySQL 数据库            │
    │  • Redis 缓存              │
    └──────────────────────────────┘
                    │
                    ↓
┌──────────────────────────────┐
│   管理后台 (admin/backend)  │
│   (端口 8081)               │
│  • 用户管理                 │
│  • 订单审核                 │
│  • 财务管理                 │
│  • 系统配置                 │
│  • RBAC 权限               │
└──────────────────────────────┘
```

---

## 🚀 启动服务

### 1. 用户端 API (8080)
```bash
cd api
go run main.go
```

### 2. 管理后台 API (8081)
```bash
cd admin/backend
go run main.go
```

### 3. WebSocket 服务 (8082)
```bash
cd websocket
go run main.go
```

---

## 📊 服务分离对比

### 分离前
```
api/
├── 19 个 handlers (包含用户端 + 管理后台)
├── 10 个 models
├── 70+ 个接口 (用户端 + 管理后台)
└── 单一服务 (8080)
```

### 分离后
```
api/ (用户端)
├── 9 个 handlers
├── 6 个 models
├── 30+ 个接口 (纯用户端)
└── 独立服务 (8080)

admin/backend/ (管理后台)
├── 23 个 handlers
├── 11 个 models
├── 70+ 个接口 (纯管理后台)
└── 独立服务 (8081)

websocket/ (实时通信)
├── 4 个文件
└── 独立服务 (8082)
```

---

## ✅ 分离优势

### 1. **职责清晰**
- 用户端 API 专注于用户和咨询师业务
- 管理后台专注于管理和运营
- WebSocket 专注于实时通信

### 2. **独立部署**
- 三个服务可以独立部署和扩展
- 故障隔离，互不影响
- 灵活分配资源

### 3. **性能优化**
- 管理后台和用户端分离，避免相互影响
- WebSocket 独立服务，专门处理实时通信

### 4. **维护便利**
- 代码模块化，便于团队协作
- 降低单个服务的复杂度
- 便于测试和调试

### 5. **安全增强**
- 管理后台可以独立配置安全策略
- 用户端和管理后台隔离

---

## 📝 注意事项

1. **数据库共享**: 三个服务共享同一个 MySQL 数据库
2. **缓存共享**: 三个服务可以共享同一个 Redis 实例
3. **跨域配置**: 每个服务都需要配置 CORS
4. **JWT 验证**: 用户端和管理后台使用不同的 JWT 验证逻辑
5. **前端配置**: 前端需要分别配置三个服务的地址

---

## 🔄 迁移指南

### 如果之前使用旧的 api 服务

1. **前端配置更新**
```javascript
// 旧配置
const API_BASE = 'http://localhost:8080'

// 新配置
const USER_API_BASE = 'http://localhost:8080'      // 用户端
const ADMIN_API_BASE = 'http://localhost:8081'     // 管理后台
const WS_BASE = 'ws://localhost:8082'             // WebSocket
```

2. **接口路由更新**
```javascript
// 旧路由
const url = '/api/admin/users'

// 新路由
const url = 'http://localhost:8081/api/admin/users'
```

3. **WebSocket 连接更新**
```javascript
// 旧连接
const ws = new WebSocket('ws://localhost:8080/ws')

// 新连接
const ws = new WebSocket('ws://localhost:8082/ws')
```

---

## 📄 License

MIT License
