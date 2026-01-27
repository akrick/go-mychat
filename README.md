# MyChat - 心理咨询聊天平台

基于 SOA 架构的心理咨询聊天系统，采用 Go + Vue.js 技术栈，提供用户与咨询师实时聊天、订单管理、支付结算等功能。

## 🏗️ 架构概览

本项目采用 **SOA（面向服务架构）**，将系统拆分为独立的服务单元：

```
┌─────────────────────────────────────────────────────┐
│                    客户端                         │
│  ┌─────────┐  ┌─────────┐  ┌─────────────┐    │
│  │ 用户端  │  │ 移动端  │  │  管理后台   │    │
│  └────┬────┘  └────┬────┘  └──────┬──────┘    │
└───────┼────────────┼─────────────┼────────────┘
        │            │             │
        ↓            ↓             ↓
┌─────────────────────────────────────────────────────┐
│                   网关层 (Nginx)                    │
└─────────────────────────────────────────────────────┘
        │            │             │
        ↓            ↓             ↓
┌──────────┐  ┌──────────┐  ┌─────────────┐
│API 服务  │  │WebSocket │  │管理后台 API │
│ :8080    │  │ :8082    │  │ :8081       │
└─────┬────┘  └─────┬────┘  └──────┬──────┘
      │             │              │
      └─────────────┼──────────────┘
                    ↓
          ┌────────────────┐
          │  MySQL + Redis│
          └────────────────┘
```

## 📁 项目结构

```
mychat/
├── api/                          # API 服务 (:8080)
│   ├── main.go                   # 服务入口
│   ├── handlers/                 # 业务处理器
│   ├── models/                   # 数据模型
│   ├── middleware/               # 中间件
│   ├── utils/                    # 工具函数
│   ├── cache/                    # 缓存操作
│   ├── database/                 # 数据库配置
│   └── README.md                # 服务文档
│
├── websocket/                    # WebSocket 服务 (:8082)
│   ├── main.go                   # 服务入口
│   ├── hub.go                    # 连接管理
│   ├── manager.go                # 会话管理
│   ├── message.go                # 消息处理
│   ├── models/                   # 数据模型
│   ├── cache/                    # 缓存操作
│   └── README.md                # 服务文档
│
├── admin/                        # 管理后台
│   ├── backend/                  # 管理后台后端 (:8081)
│   │   ├── main.go             # 服务入口
│   │   ├── handlers/           # 管理后台处理器
│   │   ├── models/             # 数据模型
│   │   ├── middleware/         # 中间件
│   │   ├── utils/              # 工具函数
│   │   └── README.md          # 服务文档
│   │
│   ├── frontend/                 # 管理后台前端 (:3000)
│   │   ├── src/                # 源代码
│   │   ├── package.json        # 依赖配置
│   │   ├── vite.config.js      # 构建配置
│   │   └── README.md          # 前端文档
│   │
│   ├── STRUCTURE.md             # 管理后台结构说明
│   └── README.md               # 管理后台说明
│
├── cert/                        # 证书目录
├── uploads/                     # 上传文件目录
├── README.md                    # 项目说明 (本文件)
├── SOA_ARCHITECTURE.md          # SOA 架构详细说明
├── PROJECT_STRUCTURE.md         # 项目结构说明
├── PROJECT_VERIFICATION.md      # 项目结构验证
└── ADMIN_README.md             # 管理后台使用说明
```

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Node.js 16+
- MySQL 5.7+
- Redis 6.0+

### 1. 启动 API 服务

```bash
cd api
go mod tidy
go run main.go
```

服务启动在 http://localhost:8080

**主要功能**:
- 用户认证与授权
- 订单管理
- 支付系统（微信、支付宝）
- 评价系统
- 通知系统
- 统计数据

### 2. 启动 WebSocket 服务

```bash
cd websocket
go mod tidy
go run main.go
```

服务启动在 http://localhost:8082

**主要功能**:
- 实时聊天
- 消息推送
- 会话管理
- 在线状态管理

### 3. 启动管理后台后端

```bash
cd admin/backend
go mod tidy
go run main.go
```

服务启动在 http://localhost:8081

**主要功能**:
- 用户管理
- 咨询师管理
- 订单管理
- 聊天记录管理
- 财务管理
- 权限管理（RBAC）
- 低代码平台

### 4. 启动管理后台前端

```bash
cd admin/frontend
npm install
npm run dev
```

服务启动在 http://localhost:3000

## 📚 访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| API 文档 | http://localhost:8080/swagger/index.html | Swagger API 文档 |
| API 健康检查 | http://localhost:8080/health | 服务健康状态 |
| WebSocket 服务 | ws://localhost:8082/ws | WebSocket 连接 |
| 管理后台 | http://localhost:3000 | 管理后台界面 |

## 🎯 核心功能

### 用户端
- ✅ 用户注册/登录
- ✅ 浏览咨询师列表
- ✅ 下单购买咨询服务
- ✅ 在线聊天
- ✅ 订单管理
- ✅ 支付结算
- ✅ 评价咨询师

### 咨询师端
- ✅ 接受订单
- ✅ 在线咨询服务
- ✅ 查看收入
- ✅ 提现申请
- ✅ 查看评价

### 管理后台
- ✅ 用户管理（CRUD、密码重置）
- ✅ 咨询师管理（CRUD、状态控制）
- ✅ 订单管理（列表、统计、状态更新）
- ✅ 聊天记录管理（查询、搜索、统计）
- ✅ 财务管理（提现审核、统计报表）
- ✅ 权限管理（RBAC、角色、权限）
- ✅ 低代码平台（表单设计、页面设计）

## 🔧 技术栈

### 后端

| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.21+ | 编程语言 |
| Gin | v1.9.1 | Web 框架 |
| GORM | v1.25.5 | ORM 框架 |
| MySQL | 5.7+ | 关系型数据库 |
| Redis | 6.0+ | 缓存数据库 |
| JWT | v5.2.0 | 认证授权 |
| WebSocket | - | 实时通信 |

### 前端

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.x | 前端框架 |
| Element Plus | 2.x | UI 组件库 |
| Vue Router | 4.x | 路由管理 |
| Pinia | 2.x | 状态管理 |
| Vite | 5.x | 构建工具 |
| Axios | 1.x | HTTP 客户端 |

## 📖 文档

| 文档 | 说明 |
|------|------|
| [SOA 架构说明](SOA_ARCHITECTURE.md) | 详细的 SOA 架构设计文档 |
| [项目结构说明](PROJECT_STRUCTURE.md) | 项目目录结构详解 |
| [项目结构验证](PROJECT_VERIFICATION.md) | SOA 架构验证报告 |
| [管理后台使用说明](ADMIN_README.md) | 管理后台功能详解 |
| [API 服务文档](api/README.md) | API 服务详细说明 |
| [WebSocket 服务文档](websocket/README.md) | WebSocket 服务详细说明 |
| [管理后台后端文档](admin/backend/README.md) | 管理后台 API 说明 |
| [管理后台前端文档](admin/frontend/README.md) | 管理后台前端说明 |
| [管理后台结构说明](admin/STRUCTURE.md) | 管理后台目录结构 |

## 🐳 Docker 部署

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## 🔐 环境变量

创建 `.env` 文件配置以下变量：

```env
# 数据库
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=mychat

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your_jwt_secret
JWT_EXPIRES=24h

# 微信支付
WECHAT_APP_ID=your_app_id
WECHAT_MCH_ID=your_mch_id
WECHAT_API_KEY=your_api_key
WECHAT_NOTIFY_URL=http://your-domain.com/api/payment/wechat/callback

# 支付宝
ALIPAY_APP_ID=your_app_id
ALIPAY_PRIVATE_KEY=your_private_key
ALIPAY_PUBLIC_KEY=your_public_key
ALIPAY_NOTIFY_URL=http://your-domain.com/api/payment/alipay/callback
```

## 📊 服务端口

| 服务 | 端口 | 协议 | 说明 |
|------|------|------|------|
| API 服务 | 8080 | HTTP | RESTful API |
| WebSocket 服务 | 8082 | HTTP/WebSocket | 实时通信 |
| 管理后台 API | 8081 | HTTP | 管理后台 API |
| 管理后台前端 | 3000 | HTTP | Web 界面 |
| MySQL | 3306 | TCP | 数据库 |
| Redis | 6379 | TCP | 缓存 |

## 🔄 服务通信

### API 服务 → WebSocket 服务
```http
GET http://localhost:8082/ws/stats
```

### 客户端 → WebSocket 服务
```javascript
const ws = new WebSocket('ws://localhost:8082/ws');
```

### 管理后台前端 → 管理后台 API
```http
GET http://localhost:8081/api/admin/users
```

## 🛠️ 开发指南

### 添加新的 API 接口
1. 在 `api/handlers/` 中创建处理器函数
2. 在 `api/main.go` 中注册路由
3. 添加 Swagger 注释
4. 运行 `swag init` 更新文档

### 添加新的管理后台页面
1. 在 `admin/frontend/src/views/` 中创建页面组件
2. 在 `admin/frontend/src/api/` 中创建 API 接口文件
3. 在 `admin/frontend/src/router/` 中添加路由配置

### WebSocket 消息处理
1. 在 `websocket/message.go` 中定义消息类型
2. 在 `websocket/hub.go` 中处理消息逻辑
3. 更新前端 WebSocket 客户端

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

MIT License

## 📧 联系方式

- 项目主页: https://github.com/your-repo/mychat
- 问题反馈: https://github.com/your-repo/mychat/issues
- 架构设计: [SOA_ARCHITECTURE.md](SOA_ARCHITECTURE.md)

## 🎉 致谢

感谢所有为此项目做出贡献的开发者！

---

**注意**: 本项目采用 SOA 架构设计，每个服务都是独立的，可以单独部署和扩展。详细架构说明请参考 [SOA 架构文档](SOA_ARCHITECTURE.md)。
