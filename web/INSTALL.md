# Web 前端项目说明

## 项目简介
基于 Vue 3 + Element Plus 的心理咨询平台前端项目。

## 技术栈
- Vue 3
- Vue Router 4
- Pinia
- Element Plus
- Axios
- Vite
- Day.js

## 安装依赖
```bash
npm install
```

## 开发运行
```bash
npm run dev
```

访问 http://localhost:3000

## 构建生产
```bash
npm run build
```

## 功能特性
- 用户注册/登录
- 咨询师列表/详情
- 订单创建/管理
- 咨询聊天（含实时倒计时）
- 个人中心
- 响应式设计

## API 接口
后端 API 运行在 `http://localhost:8081`

主要接口：
- POST /api/register - 用户注册
- POST /api/login - 用户登录
- GET /api/counselor/list - 获取咨询师列表
- GET /api/counselor/:id - 获取咨询师详情
- POST /api/order/create - 创建订单
- GET /api/order/list - 获取订单列表
- POST /api/chat/session/:orderId/start - 开始聊天会话
- POST /api/chat/session/:sessionId/message - 发送消息
- GET /api/chat/session/:sessionId/messages - 获取消息列表
- POST /api/chat/session/:sessionId/end - 结束会话
