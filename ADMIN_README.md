# MyChat 管理后台使用说明

## 功能概述

本管理后台系统基于 Gin + Vue.js + ElementUI 实现，包含用户管理、管理员管理、订单管理、心理咨询师管理、聊天记录等核心模块。

## 项目结构

```
mychat/
├── admin/                      # 管理后台目录
│   ├── backend/                # 后端服务 (Go)
│   │   ├── main.go            # 主程序入口
│   │   ├── go.mod             # Go 模块文件
│   │   ├── go.sum             # 依赖锁定
│   │   ├── database/          # 数据库相关
│   │   ├── models/            # 数据模型
│   │   ├── handlers/          # 处理器
│   │   ├── middleware/        # 中间件
│   │   ├── utils/             # 工具函数
│   │   ├── websocket/         # WebSocket
│   │   └── cache/             # 缓存
│   │
│   └── frontend/              # 前端应用 (Vue 3)
│       ├── index.html         # HTML 入口
│       ├── package.json       # 依赖配置
│       ├── vite.config.js     # Vite 配置
│       └── src/
│           ├── main.js        # 主入口
│           ├── api/           # API 接口
│           ├── router/        # 路由
│           ├── stores/        # 状态管理
│           ├── views/         # 页面组件
│           ├── layout/        # 布局
│           ├── utils/         # 工具
│           └── styles/        # 样式
│
├── database/                   # 数据库连接
├── models/                     # 数据模型（主项目）
├── handlers/                   # 处理器（主项目）
└── ...
```

## 系统架构

### 后端技术栈
- **框架**: Gin Web Framework
- **ORM**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **认证**: JWT
- **WebSocket**: 实时通信

### 前端技术栈
- **框架**: Vue 3
- **UI组件**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router
- **构建工具**: Vite

## 快速开始

### 1. 后端服务

```bash
cd admin/backend

# 安装依赖
go mod tidy

# 运行服务
go run main.go
```

后端服务将启动在端口 **:8081**

### 2. 前端应用

```bash
cd admin/frontend

# 安装依赖
npm install

# 运行开发服务器
npm run dev
```

前端应用将启动在端口 **:3000**，并自动代理 API 请求到后端

## 核心功能模块

### 1. 用户管理 (`/api/admin/users`)
- **查看用户列表**: 支持分页、搜索关键词、状态筛选
- **创建用户**: 支持设置用户名、密码、邮箱、手机号、头像、状态、管理员权限
- **编辑用户**: 更新用户信息
- **删除用户**: 删除指定用户
- **重置密码**: 管理员可重置用户密码

### 2. 咨询师管理 (`/api/counselor/*`)
- **查看咨询师列表**: 支持分页、搜索
- **创建咨询师**: 设置姓名、职称、头像、个人简介、擅长领域、单价、从业年限、评分
- **更新咨询师**: 修改咨询师信息
- **删除咨询师**: 删除指定咨询师
- **查看详情**: 获取咨询师完整信息

### 3. 订单管理 (`/api/admin/orders`)
- **查看订单列表**: 支持分页、搜索订单号、状态筛选
- **查看订单统计**: 总订单数、待支付/已支付/已完成/已取消/已退款统计
- **更新订单状态**: 管理员可更新订单状态

### 4. 聊天记录管理 (`/api/admin/chat/*`)
- **查看会话列表**: 查看所有聊天会话，支持状态筛选、关键词搜索
- **查看聊天消息**: 查看指定会话的所有消息记录
- **聊天统计**: 总会话数、活跃会话数、总消息数、今日消息数、总计费金额
- **搜索消息**: 全局搜索聊天消息内容
- **删除会话**: 删除指定会话及其所有消息

### 5. 财务管理
- **提现审核** (`/api/admin/withdraws/pending`):
  - 查看待审核提现申请
  - 通过提现申请（解冻并扣除余额）
  - 拒绝提现申请（解冻并退回余额，需填写拒绝原因）

- **财务统计** (`/api/admin/statistics`):
  - 总用户数
  - 在线用户数
  - 咨询师数
  - 总订单数
  - 今日订单数
  - 总交易额
  - 今日收入
  - 活跃会话数

### 6. 系统管理 (`/api/admin/*`)
- **管理员登录**: 使用管理员账号登录系统
- **获取管理员信息**: 获取当前登录的管理员信息
- **获取权限列表**: 获取当前管理员的权限列表
- **获取统计数据**: 获取系统整体统计数据
- **会话统计**: 获取WebSocket会话统计信息
- **在线用户**: 获取当前在线的所有用户
- **广播消息**: 向所有在线用户广播系统消息

### 7. RBAC权限管理 (`/api/admin/roles`, `/api/admin/permissions`)
- **角色管理**:
  - 创建/编辑/删除角色
  - 为角色分配权限
  - 查看角色权限列表

- **权限管理**:
  - 树形权限结构
  - 创建/编辑/删除权限
  - 权限类型：菜单、按钮、接口

### 8. 低代码模块 (`/api/admin/lowcode/*`)
- **表单设计**: 可视化拖拽表单设计器
- **页面设计**: 低代码页面构建
- **数据管理**: 表单数据管理

## API 接口文档

### 基础信息
- **后端地址**: `http://localhost:8081`
- **API 前缀**: `/api/admin`
- **认证方式**: JWT Bearer Token

### 用户管理相关

#### 获取用户列表
```
GET /api/admin/users
参数:
  - page: 页码（默认1）
  - page_size: 每页数量（默认20）
  - keyword: 搜索关键词
  - status: 状态（0-禁用，1-正常）
```

#### 创建用户
```
POST /api/admin/users
Body:
{
  "username": "用户名",
  "password": "密码",
  "email": "邮箱",
  "phone": "手机号",
  "avatar": "头像URL",
  "status": 1,
  "is_admin": false
}
```

#### 更新用户
```
PUT /api/admin/users/{id}
Body:
{
  "email": "邮箱",
  "phone": "手机号",
  "avatar": "头像URL",
  "status": 1,
  "is_admin": false
}
```

#### 删除用户
```
DELETE /api/admin/users/{id}
```

#### 重置密码
```
POST /api/admin/users/{id}/password
Body:
{
  "password": "新密码"
}
```

### 咨询师管理相关

#### 获取咨询师列表
```
GET /api/counselor/list
参数:
  - page: 页码
  - page_size: 每页数量
```

#### 创建咨询师
```
POST /api/counselor/create
Body:
{
  "name": "姓名",
  "title": "职称",
  "avatar": "头像URL",
  "bio": "个人简介",
  "specialty": "擅长领域",
  "price": 5.00,
  "years_exp": 10,
  "rating": 5.00
}
```

### 订单管理相关

#### 获取订单列表（管理员）
```
GET /api/admin/orders
参数:
  - page: 页码
  - page_size: 每页数量
  - status: 订单状态（0-待支付，1-已支付，2-已完成，3-已取消，4-已退款）
  - keyword: 搜索关键词（订单号/备注）
```

#### 获取订单统计
```
GET /api/admin/orders/statistics
返回:
{
  "total_orders": 总订单数,
  "pending_orders": 待支付订单数,
  "paid_orders": 已支付订单数,
  "completed_orders": 已完成订单数,
  "cancelled_orders": 已取消订单数,
  "total_amount": 总交易额,
  "today_amount": 今日收入,
  "this_month_amount": 本月收入
}
```

#### 更新订单状态
```
PUT /api/admin/orders/{id}/status
Body:
{
  "status": 1
}
```

### 聊天记录管理相关

#### 获取会话列表（管理员）
```
GET /api/admin/chat/sessions
参数:
  - page: 页码
  - page_size: 每页数量
  - status: 会话状态（0-待开始，1-进行中，2-已结束，3-已超时）
  - keyword: 搜索关键词（用户名/咨询师名）
```

#### 获取会话消息列表
```
GET /api/admin/chat/sessions/{session_id}/messages
参数:
  - page: 页码
  - page_size: 每页数量
```

#### 获取聊天统计
```
GET /api/admin/chat/statistics
返回:
{
  "total_sessions": 总会话数,
  "active_sessions": 活跃会话数,
  "total_messages": 总消息数,
  "today_messages": 今日消息数,
  "total_billing": 总计费金额,
  "recent_messages": 最近消息列表
}
```

### 系统管理相关

#### 管理员登录
```
POST /api/admin/login
Body:
{
  "username": "管理员用户名",
  "password": "密码"
}
返回:
{
  "code": 200,
  "msg": "登录成功",
  "data": {
    "token": "JWT token",
    "user": {用户信息}
  }
}
```

#### 获取统计数据
```
GET /api/admin/statistics
返回:
{
  "user_count": 总用户数,
  "online_user_count": 在线用户数,
  "counselor_count": 咨询师数,
  "order_count": 总订单数,
  "today_order_count": 今日订单数,
  "total_amount": 总交易额,
  "today_amount": 今日收入,
  "session_count": 总会话数,
  "active_session_count": 活跃会话数
}
```

## 前端页面结构

### 系统管理
- `/user` - 用户管理页面
- `/counselor` - 咨询师管理页面
- `/roles` - 角色管理页面
- `/permissions` - 权限管理页面
- `/menus` - 菜单管理页面

### 业务管理
- `/order` - 订单管理页面
- `/chat` - 聊天记录页面

### 财务管理
- `/withdraw` - 提现审核页面
- `/statistics` - 财务统计页面

### 低代码平台
- `/lowcode/forms` - 表单设计页面
- `/lowcode/pages` - 页面设计页面
- `/lowcode/data` - 数据管理页面

## 部署说明

### 后端部署

```bash
cd admin/backend

# 编译
go build -o admin-backend .

# 运行
./admin-backend
```

### 前端部署

```bash
cd admin/frontend

# 安装依赖
npm install

# 生产环境构建
npm run build

# 构建产物在 dist 目录
```

### Nginx 配置示例

```nginx
# 后端 API
server {
    listen 8081;
    server_name api.admin.example.com;

    location / {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# 前端页面
server {
    listen 80;
    server_name admin.example.com;

    location / {
        root /path/to/admin/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 开发说明

### 后端开发

后端代码位于 `admin/backend/` 目录，使用以下技术：

- Gin Web Framework
- GORM ORM
- JWT 认证
- WebSocket 实时通信

添加新的管理模块：

1. 在 `handlers/` 目录下创建新的 handler 文件
2. 定义相关函数
3. 在 `main.go` 中添加路由

### 前端开发

前端代码位于 `admin/frontend/` 目录，使用以下技术：

- Vue 3 Composition API
- Element Plus UI 组件库
- Vue Router
- Pinia 状态管理
- Vite 构建工具

添加新的管理页面：

1. 在 `admin/frontend/src/views/` 下创建页面组件
2. 在 `admin/frontend/src/api/` 下创建 API 接口文件
3. 在 `admin/frontend/src/router/index.js` 中添加路由配置

### RBAC 权限控制

1. 在路由配置中添加 `permission` 元数据
2. 在 handler 中检查用户权限
3. 前端根据权限显示/隐藏菜单和按钮

## 注意事项

1. 所有 API 接口都需要 JWT token 认证（除登录接口外）
2. 管理员接口需要用户具有管理员权限 (`is_admin = true`)
3. 敏感操作（如删除、重置密码）需要二次确认
4. 涉及金额的操作需谨慎处理
5. 聊天消息删除会同时删除会话下的所有消息
6. 提现审核会自动更新咨询师账户余额
7. 后端默认端口为 8081，前端默认端口为 3000
8. 前端开发环境通过 Vite 代理访问后端 API

## 常见问题

### Q: 如何创建初始管理员账号？
A: 需要直接在数据库中插入用户，设置 `is_admin = true`，或者通过用户管理页面创建时勾选"管理员"选项。

### Q: 权限不生效怎么办？
A: 检查：
1. 用户角色是否正确分配
2. 角色是否包含所需权限
3. 权限代码是否与代码中一致
4. Token 是否有效

### Q: 前端页面空白？
A: 检查：
1. 控制台是否有错误
2. API 接口是否正常返回
3. Token 是否有效
4. 网络请求是否成功
5. 后端服务是否启动在 8081 端口

### Q: 如何切换环境？
A: 修改 `frontend/vite.config.js` 中的代理目标地址和 `frontend/src/utils/request.js` 中的 baseURL。

## 更新日志

### v2.0.0 (2026-01-26)
- 重构项目结构，将管理后台代码独立到 admin 目录
- admin/backend: 后端 Go 服务，独立模块管理
- admin/frontend: 前端 Vue 3 应用
- 更新所有导入路径和依赖配置
- 后端端口改为 8081，前端端口保持 3000
- 添加详细的 README 文档

### v1.0.0 (2024-01-26)
- 完成用户管理模块
- 完成咨询师管理模块
- 完成订单管理模块
- 完成聊天记录管理模块
- 完成财务管理模块
- 完成RBAC权限管理模块
- 完成低代码平台基础模块
