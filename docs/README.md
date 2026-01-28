# MyChat 项目文档

> 一个完整的在线咨询聊天系统，包含用户端、管理后台和实时通信功能

## 📁 项目结构

```
mychat/
├── api/                  # 用户端 API 服务 (端口 8080)
├── admin/               # 管理后台
│   ├── backend/         # 管理后台 API 服务 (端口 8081)
│   └── frontend/       # 管理后台前端 (Vue + Element Plus)
├── websocket/           # WebSocket 服务 (端口 8082)
├── docs/                # 项目文档
├── uploads/             # 文件上传目录
└── .gitignore          # Git 忽略配置
```

## 🏗️ 服务架构

| 服务 | 端口 | 技术栈 | 职责 |
|------|------|--------|------|
| **用户端 API** | 8080 | Go + Gin + GORM | 用户、咨询师、订单、支付、聊天 |
| **管理后台 API** | 8081 | Go + Gin + GORM | 用户管理、订单审核、财务管理、RBAC |
| **管理后台前端** | 3000 | Vue + Element Plus | 后台管理系统界面 |
| **WebSocket 服务** | 8082 | Go + Gorilla WebSocket | 实时消息推送、在线状态 |

## 📚 文档索引

### 核心文档
- **[服务分离说明](./SERVICE_SEPARATION.md)** - 三大服务的职责和接口说明
- **[API 文档](./API.md)** - 完整的 API 接口文档
- **[数据库文档](./DATABASE.md)** - 数据库表结构和设计说明

### 清理文档
- **[API 目录清理](./API_CLEANUP.md)** - api 目录清理记录
- **[项目清理总结](../CLEANUP_SUMMARY.md)** - 整体项目整理记录

## 🚀 快速开始

### 环境要求
- Go 1.20+
- Node.js 16+
- MySQL 5.7+
- Redis 6.0+ (可选)

### 1. 初始化数据库
```bash
mysql -u root -p mychat < api/init_data.sql
```

### 2. 启动服务

#### 用户端 API (8080)
```bash
cd api
go run main.go
```
- Swagger 文档: http://localhost:8080/swagger/index.html
- 健康检查: http://localhost:8080/health

#### 管理后台 API (8081)
```bash
cd admin/backend
go run main.go
```
- 测试接口: http://localhost:8081/api/test

#### WebSocket 服务 (8082)
```bash
cd websocket
go run main.go
```

#### 管理后台前端 (3000)
```bash
cd admin/frontend
npm install
npm run dev
```
- 访问地址: http://localhost:3000
- 默认账号: admin / 123456

## 🔐 默认账户

### 管理员
- 用户名: `admin`
- 密码: `123456`

### 测试用户
- 用户名: `testuser`
- 密码: `123456`

## 📋 功能模块

### 用户端功能
- ✅ 用户注册/登录
- ✅ 用户信息管理
- ✅ 咨询师列表和详情
- ✅ 订单创建和管理
- ✅ 在线支付（微信、支付宝）
- ✅ 实时聊天
- ✅ 账单和提现
- ✅ 通知管理
- ✅ 文件上传

### 管理后台功能
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

## 🛠️ 技术栈

### 后端
- **框架**: Gin Web Framework
- **ORM**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **认证**: JWT
- **实时通信**: WebSocket (Gorilla WebSocket)
- **支付**: 微信支付、支付宝

### 前端
- **框架**: Vue 3
- **UI 组件**: Element Plus
- **HTTP 请求**: Axios
- **状态管理**: Pinia
- **路由**: Vue Router

## 📊 数据库

### 主要表结构
- users - 用户表
- counselors - 咨询师表
- orders - 订单表
- chat_sessions - 聊天会话表
- chat_messages - 聊天消息表
- payments - 支付表
- withdrawals - 提现表
- notifications - 通知表
- roles - 角色表
- permissions - 权限表
- menus - 菜单表
- system_configs - 系统配置表

详细表结构请参考 [数据库文档](./DATABASE.md)

## 🔗 API 文档

### 用户端 API (8080)
- [Swagger 文档](http://localhost:8080/swagger/index.html)
- 详细接口说明参考 [服务分离说明](./SERVICE_SEPARATION.md)

### 管理后台 API (8081)
- 详细接口说明参考 [服务分离说明](./SERVICE_SEPARATION.md)

### WebSocket 消息格式
- 连接: `ws://localhost:8082/ws`
- 消息格式: JSON
- 详见 [API 文档](./API.md)

## 🔧 开发指南

### 添加新的 API 接口

#### 用户端 API (api/)
1. 在 `api/handlers/` 创建处理函数
2. 在 `api/main.go` 注册路由
3. 更新 Swagger 注释
4. 重新生成文档: `swag init`

#### 管理后台 API (admin/backend/)
1. 在 `admin/backend/handlers/` 创建处理函数
2. 在 `admin/backend/main.go` 注册路由

### 添加新的数据表
1. 在 MySQL 中创建表
2. 在对应 `models/` 目录创建模型文件
3. 在 `database/db.go` 中添加 AutoMigrate

## 📝 项目规范

### 代码规范
- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 函数命名使用驼峰命名法
- 接口注释遵循 Swagger 格式

### Git 提交规范
```
feat: 新功能
fix: 修复 bug
docs: 文档更新
style: 代码格式调整
refactor: 重构
test: 测试相关
chore: 构建/工具相关
```

### 分支策略
- `master` - 主分支，生产环境
- `develop` - 开发分支
- `feature/*` - 功能分支
- `bugfix/*` - 修复分支

## 🐛 常见问题

### 1. 服务启动失败
- 检查端口是否被占用
- 确认数据库连接配置正确
- 确认 Redis 是否启动（可选）

### 2. 数据库连接失败
- 检查 MySQL 是否启动
- 确认数据库配置在 `database/db.go`
- 确认用户名密码正确

### 3. 前端无法连接 API
- 检查跨域配置 (CORS)
- 确认 API 服务已启动
- 检查防火墙设置

### 4. WebSocket 连接失败
- 确认 WebSocket 服务已启动 (8082)
- 检查 WebSocket URL 是否正确
- 确认网络连接正常

## 📄 License

MIT License

## 👥 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 联系方式

如有问题，请通过以下方式联系：
- 提交 Issue
- 发送邮件

---

**最后更新**: 2026-01-28
