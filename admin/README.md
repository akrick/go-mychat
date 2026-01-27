# MyChat 管理后台

基于 Vue 3 + Element Plus + Go + Gin 构建的现代化管理后台系统。

## 目录结构

```
admin/
├── backend/                    # 后端服务
│   ├── main.go               # 服务入口
│   ├── go.mod                # Go 模块文件
│   ├── go.sum                # 依赖锁定
│   ├── database/             # 数据库配置
│   ├── models/               # 数据模型
│   ├── handlers/             # 业务处理器
│   ├── middleware/           # 中间件
│   ├── utils/                # 工具函数
│   ├── websocket/            # WebSocket 相关
│   └── cache/                # 缓存操作
│
├── frontend/                   # 前端应用
│   ├── src/
│   │   ├── api/             # API 接口
│   │   ├── router/          # 路由配置
│   │   ├── stores/          # 状态管理
│   │   ├── views/           # 页面组件
│   │   ├── layout/          # 布局组件
│   │   ├── utils/           # 工具函数
│   │   └── styles/          # 样式文件
│   ├── package.json         # 依赖配置
│   ├── vite.config.js       # Vite 配置
│   └── index.html           # HTML 入口
│
└── README.md                # 本文档
```

## 功能模块

### 1. 数据看板
- 实时统计用户、咨询师、订单数量
- 营收趋势图表展示
- 订单状态分布
- 咨询师排名

### 2. 系统管理
#### 用户管理
- 用户列表查询（分页、搜索、筛选）
- 用户信息 CRUD
- 密码重置
- 用户状态管理

#### 角色管理
- 角色列表查询
- 角色权限配置
- 角色状态管理

#### 权限管理
- 权限列表查询
- 权限分组管理
- 权限状态管理

#### 菜单管理
- 菜单树形结构
- 菜单类型（目录/菜单/按钮）
- 菜单排序

#### 咨询师管理
- 咨询师列表查询
- 咨询师信息管理
- 咨询师状态控制

### 3. 业务管理
#### 订单管理
- 订单列表查询
- 订单状态更新
- 订单详情查看

#### 聊天记录管理
- 会话列表查询
- 聊天消息查看
- 消息搜索
- 会话统计

### 4. 财务管理
#### 提现审核
- 待审核提现列表
- 提现审核（通过/拒绝）
- 提现记录查询

#### 财务统计
- 营收统计（总营收、咨询师收益、平台佣金）
- 营收趋势图表
- 提现统计

#### 财务报表
- 营收报表（按天/月/年）
- 订单统计
- 财务数据导出

### 5. 低代码平台
#### 表单设计
- 可视化表单设计器
- 表单组件拖拽
- 表单预览

#### 页面设计
- 页面布局设计
- 组件配置
- 页面预览

#### 数据管理
- 表单数据查询
- 数据统计
- 数据导出

## 技术栈

### 后端
- **Go 1.21+** - 编程语言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **JWT** - 认证授权
- **MySQL** - 关系型数据库
- **Redis** - 缓存数据库

### 前端
- **Vue 3** - 前端框架
- **Element Plus** - UI 组件库
- **Vue Router 4** - 路由管理
- **Pinia** - 状态管理
- **Axios** - HTTP 客户端
- **ECharts** - 图表库
- **Vite** - 构建工具

## 快速开始

### 后端启动

```bash
cd admin/backend
go mod tidy
go run main.go
```

服务启动在 `http://localhost:8081`

### 前端启动

```bash
cd admin/frontend
npm install
npm run dev
```

服务启动在 `http://localhost:3000`

### 默认账号

- 用户名: `admin`
- 密码: `admin123`

## API 接口

### 认证相关
- `POST /api/admin/login` - 管理员登录
- `POST /api/admin/logout` - 管理员登出
- `GET /api/admin/user/info` - 获取管理员信息
- `GET /api/admin/user/permissions` - 获取管理员权限

### 用户管理
- `GET /api/admin/users` - 获取用户列表
- `POST /api/admin/users` - 创建用户
- `PUT /api/admin/users/:id` - 更新用户
- `DELETE /api/admin/users/:id` - 删除用户
- `POST /api/admin/users/:id/password` - 重置密码

### 角色管理
- `GET /api/admin/roles` - 获取角色列表
- `POST /api/admin/roles` - 创建角色
- `PUT /api/admin/roles/:id` - 更新角色
- `DELETE /api/admin/roles/:id` - 删除角色
- `GET /api/admin/roles/:id/permissions` - 获取角色权限
- `PUT /api/admin/roles/:id/permissions` - 更新角色权限

### 权限管理
- `GET /api/admin/permissions` - 获取权限列表
- `POST /api/admin/permissions` - 创建权限
- `PUT /api/admin/permissions/:id` - 更新权限
- `DELETE /api/admin/permissions/:id` - 删除权限

### 菜单管理
- `GET /api/admin/menus` - 获取菜单列表
- `POST /api/admin/menus` - 创建菜单
- `PUT /api/admin/menus/:id` - 更新菜单
- `DELETE /api/admin/menus/:id` - 删除菜单

### 咨询师管理
- `GET /api/admin/counselors` - 获取咨询师列表
- `POST /api/admin/counselors` - 创建咨询师
- `PUT /api/admin/counselors/:id` - 更新咨询师
- `DELETE /api/admin/counselors/:id` - 删除咨询师

### 订单管理
- `GET /api/admin/orders` - 获取订单列表
- `GET /api/admin/orders/statistics` - 获取订单统计
- `PUT /api/admin/orders/:id/status` - 更新订单状态

### 聊天管理
- `GET /api/admin/chat/sessions` - 获取聊天会话列表
- `GET /api/admin/chat/sessions/:id/messages` - 获取会话消息
- `GET /api/admin/chat/statistics` - 获取聊天统计
- `GET /api/admin/chat/messages/search` - 搜索聊天消息
- `DELETE /api/admin/chat/sessions/:id` - 删除聊天会话

### 财务管理
- `GET /api/admin/withdraws/pending` - 获取待审核提现
- `GET /api/admin/withdraws` - 获取提现记录
- `POST /api/admin/withdraw/:id/approve` - 审核提现
- `GET /api/admin/finance/stats` - 获取财务统计
- `GET /api/admin/finance/revenue` - 获取营收报表

### 统计数据
- `GET /api/admin/statistics` - 获取管理员统计数据
- `GET /api/admin/session/stats` - 获取会话统计
- `GET /api/admin/online/users` - 获取在线用户
- `POST /api/admin/broadcast` - 广播系统消息

### 低代码平台
- `GET /api/admin/lowcode/forms` - 获取表单列表
- `POST /api/admin/lowcode/forms` - 创建表单
- `PUT /api/admin/lowcode/forms/:id` - 更新表单
- `DELETE /api/admin/lowcode/forms/:id` - 删除表单
- `GET /api/admin/lowcode/forms/:id` - 获取表单详情
- `GET /api/admin/lowcode/forms/:id/data` - 获取表单数据
- `POST /api/admin/lowcode/forms/:id/submit` - 提交表单数据
- `DELETE /api/admin/lowcode/data/:id` - 删除表单数据
- `GET /api/admin/lowcode/pages` - 获取页面列表
- `POST /api/admin/lowcode/pages` - 创建页面
- `PUT /api/admin/lowcode/pages/:id` - 更新页面
- `DELETE /api/admin/lowcode/pages/:id` - 删除页面
- `GET /api/admin/lowcode/pages/:id` - 获取页面详情

## 环境变量

创建 `admin/backend/.env` 文件：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=mychat

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT 配置
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRES=24h

# 服务配置
SERVER_PORT=8081
SERVER_MODE=release
```

## 构建部署

### 后端构建

```bash
cd admin/backend
go build -o admin-backend .
```

### 前端构建

```bash
cd admin/frontend
npm run build
```

构建产物在 `admin/frontend/dist` 目录

### Docker 部署

```bash
cd admin
docker-compose up -d
```

## 开发指南

### 添加新的页面

1. 在 `frontend/src/views/` 创建页面组件
2. 在 `frontend/src/api/` 创建对应的 API 接口文件
3. 在 `frontend/src/router/index.js` 添加路由配置

### 添加新的 API

1. 在 `backend/handlers/` 创建处理器函数
2. 在 `backend/main.go` 注册路由
3. 添加 Swagger 注释

### 添加新的数据模型

1. 在 `backend/models/` 创建模型文件
2. 定义 GORM 结构体
3. 运行数据库迁移

## 权限系统

### 角色权限

系统采用 RBAC (Role-Based Access Control) 权限模型：

- **用户** - 可以拥有多个角色
- **角色** - 可以拥有多个权限
- **权限** - 定义具体的操作权限

### 权限标识格式

```
模块:操作

例如:
- system:user:view    - 查看用户
- system:user:create  - 创建用户
- system:user:update  - 更新用户
- system:user:delete  - 删除用户
```

## 常见问题

### 1. 登录失败

检查：
- 用户名和密码是否正确
- 账号是否被禁用
- 是否为管理员账号

### 2. API 请求失败

检查：
- 后端服务是否正常启动
- API 地址是否正确
- JWT Token 是否有效

### 3. 页面显示异常

检查：
- 浏览器控制台是否有错误
- API 响应数据格式是否正确
- Element Plus 组件是否正确引用

## 更新日志

### v1.0.0 (2024-01-26)
- ✅ 完成基础架构搭建
- ✅ 实现用户管理功能
- ✅ 实现角色权限管理
- ✅ 实现菜单管理
- ✅ 实现订单管理
- ✅ 实现聊天记录管理
- ✅ 实现财务管理
- ✅ 实现低代码平台基础功能

## 技术支持

如有问题，请提交 Issue 或联系开发团队。

## 许可证

MIT License
