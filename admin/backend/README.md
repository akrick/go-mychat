# MyChat 管理后台后端服务

管理后台的 Go 后端服务，基于 Gin + GORM 实现。

## 目录结构

```
admin/backend/
├── main.go           # 主程序入口
├── go.mod            # Go 模块文件
├── go.sum            # 依赖锁定文件
├── database/         # 数据库相关
│   └── db.go
├── models/           # 数据模型
├── handlers/         # 处理器
├── middleware/       # 中间件
├── utils/            # 工具函数
├── websocket/        # WebSocket 相关
└── cache/            # 缓存相关
```

## 安装依赖

```bash
cd admin/backend
go mod tidy
```

## 运行服务

```bash
go run main.go
```

服务将启动在端口 :8081

## 编译

```bash
go build -o admin-backend.exe .
```

## API 接口

所有 API 接口路径为 `/api/admin/*`，需要 JWT 认证。

### 主要接口

- `POST /api/admin/login` - 管理员登录
- `GET /api/admin/users` - 获取用户列表
- `GET /api/admin/orders` - 获取订单列表
- `GET /api/admin/chat/sessions` - 获取聊天会话
- `GET /api/admin/statistics` - 获取统计数据

详细 API 文档请参考项目根目录的 ADMIN_README.md
