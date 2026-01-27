# MyChat 项目目录结构

## 完整目录树

```
mychat/
├── README.md                    # 项目主说明文档
├── ADMIN_README.md              # 管理后台使用说明
├── go.mod                       # Go 模块文件
├── go.sum                       # 依赖锁定文件
│
├── admin/                       # 管理后台目录
│   ├── README.md               # 管理后台说明
│   ├── STRUCTURE.md            # 管理后台结构说明
│   │
│   ├── backend/                # 管理后台后端服务
│   │   ├── main.go            # 主程序入口
│   │   ├── go.mod             # Go 模块文件
│   │   ├── go.sum             # 依赖锁定
│   │   ├── README.md          # 后端说明文档
│   │   │
│   │   ├── database/          # 数据库相关
│   │   │   └── db.go         # 数据库初始化
│   │   │
│   │   ├── models/            # 数据模型
│   │   │   ├── chat.go       # 聊天模型
│   │   │   ├── file.go       # 文件模型
│   │   │   ├── notification.go # 通知模型
│   │   │   ├── order.go      # 订单模型
│   │   │   ├── payment.go    # 支付模型
│   │   │   ├── rbac.go       # 权限模型
│   │   │   ├── review.go     # 评论模型
│   │   │   └── user.go       # 用户模型
│   │   │
│   │   ├── handlers/          # HTTP 处理器
│   │   │   ├── admin.go      # 管理员相关
│   │   │   ├── admin_chat.go # 聊天管理
│   │   │   ├── admin_order.go # 订单管理
│   │   │   ├── auth.go       # 认证相关
│   │   │   ├── chat.go       # 聊天处理
│   │   │   ├── config.go     # 配置处理
│   │   │   ├── counselor.go  # 咨询师管理
│   │   │   ├── lowcode.go    # 低代码平台
│   │   │   ├── notification.go # 通知处理
│   │   │   ├── order.go      # 订单处理
│   │   │   ├── order_validation.go # 订单验证
│   │   │   ├── payment.go    # 支付处理
│   │   │   ├── rbac.go       # 权限控制
│   │   │   ├── review.go     # 评论处理
│   │   │   ├── stats.go      # 统计数据
│   │   │   ├── upload.go     # 文件上传
│   │   │   ├── user.go       # 用户管理
│   │   │   └── websocket.go  # WebSocket 处理
│   │   │
│   │   ├── middleware/        # 中间件
│   │   │   └── auth.go      # JWT 认证
│   │   │
│   │   ├── utils/             # 工具函数
│   │   │   ├── alipay.go    # 支付宝支付
│   │   │   ├── jwt.go       # JWT 工具
│   │   │   ├── password.go  # 密码加密
│   │   │   ├── payment_helper.go # 支付辅助
│   │   │   ├── response.go  # 响应封装
│   │   │   └── wechat_pay.go # 微信支付
│   │   │
│   │   ├── websocket/         # WebSocket 相关
│   │   │   ├── hub.go       # WebSocket 集线器
│   │   │   ├── manager.go   # WebSocket 管理
│   │   │   ├── message.go   # 消息结构
│   │   │   └── stats.go     # 统计信息
│   │   │
│   │   └── cache/            # 缓存相关
│   │       ├── chat_cache.go # 聊天缓存
│   │       ├── context.go    # 上下文
│   │       ├── order_cache.go # 订单缓存
│   │       ├── payment_cache.go # 支付缓存
│   │       └── redis.go      # Redis 配置
│   │
│   └── frontend/             # 管理后台前端应用
│       ├── index.html        # HTML 入口
│       ├── package.json      # NPM 依赖
│       ├── vite.config.js    # Vite 配置
│       ├── .gitignore        # Git 忽略
│       └── README.md         # 前端说明
│       │
│       └── src/              # 源代码
│           ├── main.js       # 主入口
│           ├── App.vue       # 根组件
│           │
│           ├── api/          # API 接口
│           │   ├── admin.js
│           │   ├── adminChat.js
│           │   ├── adminCounselor.js
│           │   ├── adminOrder.js
│           │   ├── adminUser.js
│           │   ├── lowcode.js
│           │   ├── permission.js
│           │   ├── role.js
│           │   └── user.js
│           │
│           ├── router/       # 路由
│           ├── stores/       # 状态管理
│           ├── views/        # 页面组件
│           ├── layout/       # 布局
│           ├── utils/        # 工具
│           └── styles/       # 样式
│
├── api/                       # 主 API 服务
│   ├── main.go              # 主程序
│   ├── go.mod               # Go 模块
│   ├── go.sum               # 依赖锁定
│   ├── README.md            # API 说明
│   │
│   ├── database/            # 数据库
│   ├── models/              # 数据模型
│   ├── handlers/            # 处理器
│   ├── middleware/          # 中间件
│   ├── utils/               # 工具
│   ├── cache/               # 缓存
│   └── websocket/           # WebSocket 代码
│
├── websocket/                 # 独立 WebSocket 服务
│   ├── main.go              # 主程序
│   ├── go.mod               # Go 模块
│   ├── hub.go               # 集线器
│   ├── manager.go           # 会话管理
│   ├── message.go           # 消息
│   ├── stats.go             # 统计
│   └── README.md            # WebSocket 说明
│
├── database/                  # 共享数据库配置
├── models/                    # 共享数据模型
├── handlers/                  # 共享处理器
├── utils/                     # 共享工具
├── cache/                     # 共享缓存
├── websocket/                 # 共享 WebSocket 代码
├── middleware/                # 共享中间件
├── config/                    # 配置文件
├── tasks/                     # 定时任务
├── uploads/                   # 上传文件
├── cert/                      # 证书
├── docs/                      # API 文档
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
└── ...                        # 其他文件
```

## 服务端口配置

| 服务 | 端口 | 说明 |
|------|------|------|
| 主 API 服务 | 8080 | 用户认证、订单、支付等 |
| WebSocket 服务 | 8082 | 实时聊天服务 |
| 管理后台 API | 8081 | 管理后台专用 API |
| 管理后台前端 | 3000 | 管理后台界面 |
| MySQL | 3306 | 数据库 |
| Redis | 6379 | 缓存 |

## 启动命令

### 1. 主 API 服务
```bash
cd api
go run main.go
```

### 2. WebSocket 服务
```bash
cd websocket
go run main.go
```

### 3. 管理后台后端
```bash
cd admin/backend
go run main.go
```

### 4. 管理后台前端
```bash
cd admin/frontend
npm install
npm run dev
```

## 项目关系图

```
用户端
   ↓
主 API 服务 (:8080)
   ↓
数据库/Redis + WebSocket 服务 (:8082)

管理后台前端 (:3000)
   ↓
管理后台 API (:8081)
   ↓
数据库/Redis + WebSocket 服务 (:8082)
```

## 服务依赖

### 主 API 服务依赖
- MySQL 数据库
- Redis 缓存
- WebSocket 服务（用于聊天功能）

### WebSocket 服务依赖
- MySQL 数据库
- Redis 缓存

### 管理后台 API 依赖
- MySQL 数据库
- Redis 缓存
- WebSocket 服务（用于聊天记录查询）

### 管理后台前端依赖
- 管理后台 API (:8081)

## 共享代码说明

以下目录包含共享代码，api、admin/backend、websocket 都会使用：

- `database/` - 数据库连接配置
- `models/` - 数据模型定义
- `handlers/` - 共享的处理器
- `utils/` - 共享的工具函数
- `cache/` - 共享的缓存操作
- `websocket/` - 共享的 WebSocket 代码
- `middleware/` - 共享的中间件

## 编译打包

### 主 API 服务
```bash
cd api
go build -o mychat-api .
```

### WebSocket 服务
```bash
cd websocket
go build -o mychat-websocket .
```

### 管理后台后端
```bash
cd admin/backend
go build -o admin-backend .
```

### 管理后台前端
```bash
cd admin/frontend
npm run build
```

## 部署说明

### 环境准备
1. 安装 MySQL 5.7+
2. 安装 Redis 6.0+
3. 配置环境变量

### 服务部署
1. 部署 MySQL 数据库
2. 部署 Redis
3. 部署主 API 服务
4. 部署 WebSocket 服务
5. 部署管理后台 API
6. 部署管理后台前端
7. 配置 Nginx 反向代理

### 启动顺序
1. MySQL
2. Redis
3. 主 API 服务
4. WebSocket 服务
5. 管理后台 API
6. 管理后台前端

## 更多信息

- 项目说明: `../README.md`
- 管理后台使用: `ADMIN_README.md`
- API 服务说明: `api/README.md`
- WebSocket 服务说明: `websocket/README.md`
- 管理后台后端: `backend/README.md`
- 管理后台前端: `frontend/README.md`
