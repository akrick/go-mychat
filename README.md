# MyChat 在线咨询系统

> 一个完整的在线咨询聊天系统，包含用户端、管理后台和实时通信功能

## 🌟 项目简介

MyChat 是一个功能完整的在线咨询平台，支持用户与咨询师进行实时聊天、订单管理、在线支付等功能。

### 核心功能
- 👤 用户注册/登录、个人资料管理
- 👨‍⚕️ 咨询师管理、在线状态
- 💬 实时聊天、消息推送
- 📦 订单管理、状态跟踪
- 💳 在线支付（微信、支付宝）
- 💰 账单管理、提现申请
- 🔐 角色权限管理 (RBAC)
- 📊 统计报表、数据分析
- 🎨 低代码平台

## 🏗️ 技术架构

### 服务架构

```
┌─────────────────────────────────────────────────────┐
│                   用户端应用                          │
│  (Vue/React/小程序/H5)                                │
└────────────────────┬────────────────────────────────┘
                     │
        ┌────────────┴────────────┐
        │                         │
        ↓ HTTP REST API           ↓ WebSocket
┌──────────────────┐      ┌──────────────────┐
│ 用户端 API (8080) │      │WebSocket (8082)  │
│ • 用户/咨询师      │      │ • 实时消息       │
│ • 订单/支付        │      │ • 在线状态       │
└────────┬─────────┘      └──────────────────┘
         │
         └──────────┬─────────────────┘
                    │
    ┌───────────────┴────────────────┐
    │   MySQL + Redis               │
    └───────────────┬────────────────┘
                    │
                    ↓
┌──────────────────────────────────┐
│   管理后台 (admin/backend)  │
│   (端口 8081)               │
│  • 用户管理                 │
│  • 订单审核                 │
│  • 财务管理                 │
│  • 系统配置                 │
│  • RBAC 权限               │
└──────────────────────────────────┘
```

### 技术栈

#### 后端
- **语言**: Go 1.20+
- **框架**: Gin Web Framework
- **ORM**: GORM
- **数据库**: MySQL 5.7+
- **缓存**: Redis 6.0+ (可选)
- **认证**: JWT
- **实时通信**: Gorilla WebSocket
- **支付**: 微信支付、支付宝

#### 前端
- **框架**: Vue 3
- **UI**: Element Plus
- **HTTP**: Axios
- **状态**: Pinia
- **路由**: Vue Router

## 📁 项目结构

```
mychat/
├── api/                      # 用户端 API (端口 8080)
│   ├── handlers/              # 请求处理
│   ├── models/               # 数据模型
│   ├── cache/                # 缓存层
│   ├── database/             # 数据库连接
│   ├── middleware/           # 中间件
│   ├── utils/                # 工具类
│   └── main.go               # 主入口
├── admin/                    # 管理后台
│   ├── backend/              # 管理后台 API (端口 8081)
│   │   ├── handlers/        # 请求处理
│   │   ├── models/          # 数据模型
│   │   ├── cache/           # 缓存层
│   │   ├── database/        # 数据库连接
│   │   ├── middleware/      # 中间件
│   │   ├── utils/           # 工具类
│   │   ├── websocket/       # WebSocket Hub
│   │   └── main.go         # 主入口
│   └── frontend/            # 管理后台前端 (端口 3000)
│       └── src/             # Vue 源码
├── websocket/                # WebSocket 服务 (端口 8082)
│   ├── hub.go              # WebSocket Hub
│   ├── manager.go          # 连接管理
│   ├── message.go          # 消息处理
│   └── stats.go            # 统计信息
├── docs/                    # 项目文档
│   ├── README.md           # 文档索引
│   ├── SERVICE_SEPARATION.md  # 服务分离说明
│   ├── API.md              # API 文档
│   ├── DATABASE.md          # 数据库文档
│   ├── PROJECT_STATUS.md    # 项目状态
│   └── IMPROVEMENT_PLAN.md # 完善计划
├── uploads/                 # 文件上传目录
├── start-all.bat            # 一键启动所有服务
├── stop-all.bat             # 一键停止所有服务
└── .gitignore              # Git 忽略配置
```

## 🚀 快速开始

### 环境要求
- Go 1.20+
- Node.js 16+
- MySQL 5.7+
- Redis 6.0+ (可选)

### 1. 克隆项目
```bash
git clone <repository-url>
cd mychat
```

### 2. 初始化数据库
```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE mychat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 导入初始数据
mysql -u root -p mychat < api/init_data.sql
```

### 3. 启动服务

#### 方式一：使用启动脚本（推荐）
```bash
# Windows
start-all.bat

# Linux/Mac
bash start-all.sh
```

#### 方式二：单独启动
```bash
# 用户端 API (8080)
cd api
go mod tidy
go run main.go

# 管理后台 API (8081)
cd admin/backend
go mod tidy
go run main.go

# WebSocket 服务 (8082)
cd websocket
go mod tidy
go run main.go

# 管理后台前端 (3000)
cd admin/frontend
npm install
npm run dev
```

### 4. 访问系统

#### 用户端 API
- API 地址: http://localhost:8080
- Swagger 文档: http://localhost:8080/swagger/index.html
- 健康检查: http://localhost:8080/health

#### 管理后台
- 后端 API: http://localhost:8081
- 前端界面: http://localhost:3000
- 默认账号: admin / 123456

#### WebSocket
- 连接地址: ws://localhost:8082/ws

## 🔐 默认账户

### 管理员
- 用户名: `admin`
- 密码: `123456`

### 测试用户
- 用户名: `testuser`
- 密码: `123456`

## 📚 文档

详细文档请查看 [docs/](./docs/) 目录：

- **[服务分离说明](./docs/SERVICE_SEPARATION.md)** - 详细的接口列表和职责说明
- **[API 文档](./docs/API.md)** - 完整的 API 接口文档
- **[数据库文档](./docs/DATABASE.md)** - 数据库表结构说明
- **[项目状态](./docs/PROJECT_STATUS.md)** - 当前项目状态和问题清单
- **[完善计划](./docs/IMPROVEMENT_PLAN.md)** - 后续完善计划

## 🛠️ 开发指南

### 添加新的 API 接口

#### 用户端 API (api/)
1. 在 `api/handlers/` 创建处理函数
2. 在 `api/main.go` 注册路由
3. 更新 Swagger 注释
4. 运行 `swag init` 生成文档

#### 管理后台 API (admin/backend/)
1. 在 `admin/backend/handlers/` 创建处理函数
2. 在 `admin/backend/main.go` 注册路由

### 添加新的数据表
1. 在 MySQL 中创建表
2. 在对应 `models/` 目录创建模型文件
3. 在 `database/db.go` 中添加 AutoMigrate

## 🐛 常见问题

### 1. 服务启动失败
- 检查端口是否被占用
- 确认数据库连接配置正确
- 确认 Redis 是否启动（可选）
- 运行 `go mod tidy` 安装依赖

### 2. 数据库连接失败
- 检查 MySQL 是否启动
- 确认数据库配置在 `database/db.go`
- 确认用户名密码正确
- 确认数据库已创建

### 3. 前端无法连接 API
- 检查跨域配置 (CORS)
- 确认 API 服务已启动
- 检查防火墙设置

### 4. WebSocket 连接失败
- 确认 WebSocket 服务已启动 (8082)
- 检查 WebSocket URL 是否正确
- 确认网络连接正常

## 📊 项目进度

| 模块 | 完成度 | 状态 |
|------|--------|------|
| 管理后台 API | 95% | ✅ 运行正常 |
| 管理后台前端 | 90% | ✅ 可用 |
| 用户端 API | 70% | ⚠️ 需修复 |
| WebSocket 服务 | 60% | ⚠️ 需修复 |
| 数据库设计 | 95% | ✅ 完成 |
| 文档 | 90% | ✅ 完善 |

**整体完成度**: ~82%

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

### 提交规范
```
feat: 新功能
fix: 修复 bug
docs: 文档更新
style: 代码格式调整
refactor: 重构
test: 测试相关
chore: 构建/工具相关
```

## 📄 License

MIT License

---

**项目创建**: 2026-01-28
**最后更新**: 2026-01-28
**版本**: v1.0.0
