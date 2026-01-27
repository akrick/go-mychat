# MyChat API 服务

MyChat 项目的主 API 服务，提供用户认证、订单管理、支付、评价等核心业务功能。

## 目录结构

```
api/
├── main.go              # 主程序入口
├── go.mod               # Go 模块文件
├── go.sum               # 依赖锁定文件
├── database/            # 数据库相关
├── models/              # 数据模型
├── handlers/            # HTTP 处理器
├── middleware/          # 中间件
├── utils/               # 工具函数
└── cache/               # 缓存操作
```

## 功能模块

### 1. 用户认证
- 用户注册
- 用户登录
- Token 刷新
- 个人信息管理

### 2. 管理后台接口
- 管理员登录
- 用户管理
- 角色权限管理
- 低代码平台
- 聊天记录查询
- 订单管理

### 3. 咨询师管理
- 咨询师列表
- 咨询师详情
- 咨询师创建/更新/删除

### 4. 订单管理
- 创建订单
- 订单查询
- 订单状态更新
- 订单取消

### 5. 支付系统
- 创建支付
- 支付状态查询
- 支付回调
- 支付退款
- 支持微信支付和支付宝

### 6. 评价系统
- 创建评价
- 查询评价
- 咨询师评分统计
- 评价回复

### 7. 统计数据
- 仪表盘统计
- 订单统计
- 咨询师排行榜

### 8. 配置管理
- 支付配置
- 配置测试

### 9. 通知系统
- 通知列表
- 标记已读
- 删除通知

### 10. 聊天相关
- 开始聊天会话
- 结束聊天会话
- 查询聊天记录
- 聊天账单查询
- 咨询师账户管理

## 安装依赖

```bash
cd api
go mod tidy
```

## 运行服务

```bash
go run main.go
```

服务将启动在端口 **:8080**

## API 文档

启动服务后访问：
- Swagger 文档: http://localhost:8080/swagger/index.html
- 健康检查: http://localhost:8080/health

## 主要接口

### 用户相关
- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录
- `GET /api/user/info` - 获取用户信息

### 管理员相关
- `POST /api/admin/login` - 管理员登录
- `GET /api/admin/users` - 获取用户列表
- `GET /api/admin/statistics` - 获取统计数据

### 咨询师相关
- `GET /api/counselor/list` - 咨询师列表
- `GET /api/counselor/:id` - 咨询师详情

### 订单相关
- `POST /api/order/create` - 创建订单
- `GET /api/order/list` - 订单列表

### 支付相关
- `POST /api/payment/create` - 创建支付
- `GET /api/payment/:id` - 支付状态

### 评价相关
- `POST /api/review/create` - 创建评价
- `GET /api/review/counselor/:id` - 咨询师评价

## 端口配置

| 服务 | 端口 | 说明 |
|------|------|------|
| API 服务 | 8080 | 主 API 服务 |
| WebSocket 服务 | 8082 | 独立的 WebSocket 服务（部署时需要同时启动）|
| 数据库 | 3306 | MySQL |
| Redis | 6379 | 缓存 |

## 依赖服务

本 API 服务依赖以下服务：

1. **MySQL 数据库** - 存储业务数据
2. **Redis** - 缓存和会话管理
3. **WebSocket 服务** (:8082) - 实时聊天功能

## 编译打包

```bash
go build -o mychat-api .
```

## 注意事项

1. WebSocket 功能已独立部署，相关接口会调用 WebSocket 服务
2. 启动前需要确保 MySQL 和 Redis 服务正常运行
3. 生产环境需要配置正确的数据库连接和 Redis 连接
4. 支付功能需要配置微信支付和支付宝的相关参数

## 相关服务

- **WebSocket 服务**: `../websocket/` - 独立的实时聊天服务
- **管理后台**: `../admin/backend/` - 管理后台 API
- **管理后台前端**: `../admin/frontend/` - 管理后台前端
