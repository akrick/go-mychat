# MyChat 管理后台前端

管理后台的 Vue 3 前端应用，基于 Vue 3 + Element Plus + Vite 实现。

## 目录结构

```
admin/frontend/
├── index.html          # HTML 入口
├── package.json        # 依赖配置
├── vite.config.js      # Vite 配置
├── .gitignore          # Git 忽略文件
├── README.md           # 说明文档
└── src/
    ├── main.js         # 主入口文件
    ├── App.vue         # 根组件
    ├── api/            # API 接口
    ├── router/         # 路由配置
    ├── stores/         # 状态管理 (Pinia)
    ├── views/          # 页面组件
    ├── layout/         # 布局组件
    ├── utils/          # 工具函数
    └── styles/         # 样式文件
```

## 安装依赖

```bash
cd admin/frontend
npm install
```

## 开发环境运行

```bash
npm run dev
```

服务将启动在端口 :3000，API 请求代理到后端服务 :8081

## 生产环境构建

```bash
npm run build
```

构建产物在 `dist` 目录

## 主要功能页面

### 系统管理
- `/user` - 用户管理
- `/counselor` - 咨询师管理
- `/roles` - 角色管理
- `/permissions` - 权限管理

### 业务管理
- `/order` - 订单管理
- `/chat` - 聊天记录管理

### 财务管理
- `/withdraw` - 提现审核
- `/statistics` - 财务统计

### 低代码平台
- `/lowcode/forms` - 表单设计
- `/lowcode/pages` - 页面设计
- `/lowcode/data` - 数据管理

## 技术栈

- Vue 3
- Element Plus
- Vue Router
- Pinia
- Vite
- Axios

详细 API 文档请参考项目根目录的 ADMIN_README.md
