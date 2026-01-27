# MyChat SOA 架构说明

## 架构概述

MyChat 采用 SOA（面向服务架构）设计理念，将系统拆分为独立的、可部署的服务单元。每个服务负责特定的业务功能，通过标准的 HTTP/WebSocket 协议进行通信。

## 服务架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         客户端层                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │  用户 Web   │  │  用户 App   │  │  管理后台   │          │
│  │  (Vue.js)   │  │  (移动端)   │  │  (Vue 3)    │          │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘          │
└─────────┼────────────────┼────────────────┼───────────────────┘
          │                │                │
          ↓                ↓                ↓
┌─────────────────────────────────────────────────────────────────┐
│                        网关层 (Nginx)                         │
│                     反向代理 + 负载均衡                         │
└─────────────────────────────────────────────────────────────────┘
          │                │                │
          ↓                ↓                ↓
┌─────────────────────────────────────────────────────────────────┐
│                         应用服务层                             │
│                                                                  │
│  ┌────────────────┐  ┌────────────────┐  ┌───────────────┐  │
│  │  API 服务     │  │  WebSocket 服务│  │  管理后台 API  │  │
│  │  (:8080)      │  │  (:8082)      │  │  (:8081)       │  │
│  │               │  │               │  │               │  │
│  │  - 用户认证   │  │  - 实时聊天   │  │  - 用户管理   │  │
│  │  - 订单管理   │  │  - 消息推送   │  │  - 咨询师管理 │  │
│  │  - 支付结算   │  │  - 会话管理   │  │  - 订单管理   │  │
│  │  - 评价系统   │  │  - 在线状态   │  │  - 财务管理   │  │
│  │  - 通知系统   │  │               │  │  - 聊天管理   │  │
│  │  - 统计数据   │  │               │  │  - 权限管理   │  │
│  │               │  │               │  │  - 低代码     │  │
│  └───────┬────────┘  └───────┬────────┘  └───────┬───────┘  │
└──────────┼──────────────────┼──────────────────┼──────────────┘
           │                  │                  │
           └──────────────────┼──────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                         数据服务层                             │
│                                                                  │
│  ┌────────────────┐              ┌────────────────┐            │
│  │   MySQL        │              │    Redis       │            │
│  │   (:3306)      │              │    (:6379)     │            │
│  │                │              │                │            │
│  │  - 用户数据    │              │  - 缓存        │            │
│  │  - 订单数据    │              │  - 会话管理    │            │
│  │  - 聊天记录    │              │  - 在线状态    │            │
│  │  - 财务数据    │              │  - 消息队列    │            │
│  │                │              │                │            │
│  └────────────────┘              └────────────────┘            │
└─────────────────────────────────────────────────────────────────┘
```

## 服务详细说明

### 1. API 服务 (`api/`)

**端口**: 8080
**职责**: 核心业务逻辑处理

**功能模块**:
- 用户认证与授权
- 订单管理（创建、查询、更新、取消）
- 支付系统（微信支付、支付宝）
- 评价系统（创建、查询、回复）
- 通知系统（推送、已读标记）
- 统计数据（仪表盘、报表）
- 配置管理（支付配置）

**技术栈**:
- Gin Web Framework
- GORM (ORM)
- MySQL (数据库)
- Redis (缓存)
- JWT (认证)

**依赖**:
- MySQL 数据库
- Redis 缓存
- WebSocket 服务（用于获取在线状态等）

**目录结构**:
```
api/
├── main.go              # 服务入口
├── go.mod               # Go 模块
├── go.sum               # 依赖锁定
├── README.md            # 服务文档
├── handlers/            # 业务处理器
│   ├── auth.go         # 认证处理
│   ├── user.go         # 用户管理
│   ├── order.go        # 订单处理
│   ├── payment.go      # 支付处理
│   ├── review.go       # 评价处理
│   ├── notification.go # 通知处理
│   └── stats.go       # 统计处理
├── models/             # 数据模型
├── middleware/         # 中间件
├── utils/              # 工具函数
├── cache/              # 缓存操作
├── database/           # 数据库配置
└── websocket/          # WebSocket 客户端
```

---

### 2. WebSocket 服务 (`websocket/`)

**端口**: 8082
**职责**: 实时通信与消息推送

**功能模块**:
- 实时聊天（用户 ↔ 咨询师）
- 消息持久化
- 会话管理（创建、结束、超时）
- 在线状态管理
- 消息推送（点对点、广播）
- 聊天计费

**技术栈**:
- Gin Web Framework
- GORM (ORM)
- MySQL (数据库)
- Redis (会话管理)
- WebSocket (实时通信)

**依赖**:
- MySQL 数据库
- Redis 缓存

**目录结构**:
```
websocket/
├── main.go              # WebSocket 服务入口
├── go.mod               # Go 模块
├── README.md            # 服务文档
├── hub.go               # 连接管理
├── manager.go           # 会话管理
├── message.go           # 消息处理
├── stats.go             # 统计信息
├── models/              # 数据模型
├── database/           # 数据库配置
├── cache/              # 缓存操作
└── middleware/          # 中间件
```

**WebSocket 端点**:
- `ws://localhost:8082/ws` - 用户连接
- `ws://localhost:8082/ws/counselor/{id}` - 咨询师连接
- `GET /ws/stats` - 统计信息

---

### 3. 管理后台 API (`admin/backend/`)

**端口**: 8081
**职责**: 管理后台业务逻辑

**功能模块**:
- 用户管理（CRUD、密码重置）
- 咨询师管理（CRUD、状态控制）
- 订单管理（列表、统计、状态更新）
- 聊天记录管理（查询、搜索、统计）
- 财务管理（提现审核、统计报表）
- 权限管理（RBAC、角色、权限）
- 低代码平台（表单设计、页面设计）

**技术栈**:
- Gin Web Framework
- GORM (ORM)
- MySQL (数据库)
- Redis (缓存)
- JWT (认证)

**依赖**:
- MySQL 数据库
- Redis 缓存
- WebSocket 服务（用于查询聊天记录）

**目录结构**:
```
admin/backend/
├── main.go              # 服务入口
├── go.mod               # Go 模块
├── go.sum               # 依赖锁定
├── README.md            # 服务文档
├── handlers/            # 管理后台处理器
│   ├── admin.go        # 管理员功能
│   ├── admin_chat.go  # 聊天记录管理
│   ├── admin_order.go # 订单管理
│   ├── user.go        # 用户管理
│   ├── counselor.go   # 咨询师管理
│   ├── rbac.go       # 权限管理
│   ├── lowcode.go     # 低代码平台
│   └── ...
├── models/             # 数据模型
├── middleware/         # 中间件
├── utils/              # 工具函数
├── cache/              # 缓存操作
├── database/           # 数据库配置
└── websocket/          # WebSocket 客户端
```

---

### 4. 管理后台前端 (`admin/frontend/`)

**端口**: 3000
**职责**: 管理后台用户界面

**技术栈**:
- Vue 3
- Element Plus
- Vue Router
- Pinia
- Vite

**目录结构**:
```
admin/frontend/
├── index.html           # HTML 入口
├── package.json         # 依赖配置
├── vite.config.js       # 构建配置
├── README.md           # 前端文档
└── src/
    ├── main.js         # 主入口
    ├── App.vue         # 根组件
    ├── api/            # API 接口
    ├── router/         # 路由配置
    ├── stores/         # 状态管理
    ├── views/          # 页面组件
    ├── layout/         # 布局组件
    ├── utils/          # 工具函数
    └── styles/         # 样式文件
```

---

## 服务间通信

### API 服务 → WebSocket 服务

```javascript
// API 服务调用 WebSocket 服务获取在线状态
GET http://websocket:8082/ws/counselors/online
```

### 管理后台 API → WebSocket 服务

```javascript
// 管理后台查询聊天统计
GET http://websocket:8082/ws/stats
```

### 客户端 → 服务

```javascript
// 用户 API 调用
POST http://api:8080/api/login
GET http://api:8080/api/order/list

// WebSocket 连接
ws://websocket:8082/ws

// 管理后台 API 调用
GET http://admin-backend:8081/api/admin/users
```

---

## 数据一致性

### 数据库共享
所有服务共享同一个 MySQL 数据库，确保数据一致性。

### Redis 共享
- API 服务使用 Redis 作为缓存
- WebSocket 服务使用 Redis 管理会话和在线状态
- 管理后台 API 使用 Redis 缓存管理数据

### 数据同步策略
1. 主从复制：数据库采用主从复制，读写分离
2. 缓存更新：采用缓存穿透和缓存雪崩保护
3. 消息队列：关键操作通过 Redis Pub/Sub 通知其他服务

---

## 部署架构

### 开发环境

```
单机部署
├── API 服务 (:8080)
├── WebSocket 服务 (:8082)
├── 管理后台 API (:8081)
├── 管理后台前端 (:3000)
├── MySQL (:3306)
└── Redis (:6379)
```

### 生产环境

```
┌──────────────────────────────────────────────────────────┐
│                    Nginx 负载均衡                        │
└──────────────────────────────────────────────────────────┘
         │                    │                    │
         ↓                    ↓                    ↓
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│  API 服务    │    │ WebSocket 服务│    │ 管理后台 API │
│  集群       │    │  集群        │    │  集群        │
│  (3 实例)   │    │  (2 实例)    │    │  (2 实例)    │
└───────────────┘    └───────────────┘    └───────────────┘
         │                    │                    │
         └────────────────────┼────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────┐
│                      数据中心                            │
│  ┌────────────┐          ┌────────────┐               │
│  │  MySQL     │          │  Redis     │               │
│  │  主从集群  │          │  集群      │               │
│  └────────────┘          └────────────┘               │
└──────────────────────────────────────────────────────────┘

管理后台前端部署到 CDN 或独立服务器
```

---

## 服务独立原则

### 1. 独立部署
- 每个服务可以独立部署、升级、扩展
- 服务之间通过 HTTP/WebSocket 协议通信
- 使用 Docker 容器化部署

### 2. 独立开发
- 每个服务有独立的代码仓库
- 独立的依赖管理（go.mod）
- 独立的测试环境

### 3. 独立扩展
- 根据负载独立扩展各服务
- WebSocket 服务需要更多实例时可以单独扩展
- API 服务和 WebSocket 服务可以独立水平扩展

### 4. 独立监控
- 每个服务独立的监控指标
- 独立的日志收集
- 独立的告警配置

---

## 快速开始

### 1. 启动所有服务

```bash
# 终端 1: API 服务
cd api
go mod tidy
go run main.go

# 终端 2: WebSocket 服务
cd websocket
go mod tidy
go run main.go

# 终端 3: 管理后台 API
cd admin/backend
go mod tidy
go run main.go

# 终端 4: 管理后台前端
cd admin/frontend
npm install
npm run dev
```

### 2. 访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| API 文档 | http://localhost:8080/swagger/index.html | Swagger 文档 |
| WebSocket | ws://localhost:8082/ws | WebSocket 连接 |
| 管理后台 | http://localhost:3000 | 管理后台界面 |

---

## Docker 部署

### Docker Compose

```yaml
version: '3.8'

services:
  api:
    build: ./api
    ports:
      - "8080:8080"
    depends_on:
      - mysql
      - redis

  websocket:
    build: ./websocket
    ports:
      - "8082:8082"
    depends_on:
      - mysql
      - redis

  admin-backend:
    build: ./admin/backend
    ports:
      - "8081:8081"
    depends_on:
      - mysql
      - redis

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: mychat
    ports:
      - "3306:3306"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
```

---

## 文档索引

| 文档 | 路径 | 说明 |
|------|------|------|
| SOA 架构 | `SOA_ARCHITECTURE.md` | 本文档 |
| 项目说明 | `README.md` | 项目总体介绍 |
| 项目结构 | `PROJECT_STRUCTURE.md` | 目录结构说明 |
| 管理后台使用 | `ADMIN_README.md` | 管理后台指南 |
| API 服务 | `api/README.md` | API 服务文档 |
| WebSocket 服务 | `websocket/README.md` | WebSocket 文档 |
| 管理后台后端 | `admin/backend/README.md` | 管理后端文档 |
| 管理后台前端 | `admin/frontend/README.md` | 管理前端文档 |

---

## 设计原则

1. **单一职责**: 每个服务只负责特定的业务领域
2. **高内聚低耦合**: 服务内部高度内聚，服务之间低耦合
3. **独立部署**: 服务可以独立部署和扩展
4. **故障隔离**: 单个服务故障不影响其他服务
5. **可观测性**: 每个服务提供监控和日志
6. **接口标准化**: 使用 RESTful API 和 WebSocket 标准接口
