# MyChat 管理后台开发指南

## 开发环境搭建

### 前置要求

- Go 1.21+
- Node.js 16+
- MySQL 5.7+
- Redis 6.0+
- Git

### 环境配置

#### 1. 数据库配置

创建数据库并导入表结构：

```sql
CREATE DATABASE mychat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE mychat;

-- 导入表结构 SQL 文件
SOURCE /path/to/schema.sql;
```

#### 2. Redis 配置

确保 Redis 服务正常运行：

```bash
redis-server
```

#### 3. 环境变量

复制环境变量模板：

```bash
cd admin/backend
cp .env.example .env
```

编辑 `.env` 文件，配置数据库和 Redis 连接信息。

## 后端开发

### 项目结构

```
backend/
├── main.go              # 主程序入口
├── database/            # 数据库相关
│   └── db.go          # 数据库初始化
├── models/             # 数据模型
├── handlers/           # 业务处理器
├── middleware/         # 中间件
├── utils/              # 工具函数
├── cache/              # 缓存操作
└── websocket/          # WebSocket 相关
```

### 添加新的 API 接口

1. **创建 Handler**

在 `backend/handlers/` 创建新的处理器文件：

```go
package handlers

import (
    "akrick.com/mychat/admin/backend/database"
    "akrick.com/mychat/admin/backend/models"
    "github.com/gin-gonic/gin"
)

// ExampleHandler 示例处理器
func ExampleHandler(c *gin.Context) {
    // 业务逻辑
    c.JSON(200, gin.H{
        "code": 200,
        "msg":  "成功",
        "data": nil,
    })
}
```

2. **注册路由**

在 `backend/main.go` 中注册路由：

```go
// 公开路由
public := r.Group("/api")
{
    public.GET("/example", handlers.ExampleHandler)
}

// 需要认证的路由
auth := r.Group("/api")
auth.Use(authMiddleware)
{
    auth.GET("/admin/example", handlers.ExampleHandler)
}
```

3. **添加 Swagger 注释**

```go
// ExampleHandler godoc
// @Summary 示例接口
// @Description 这是一个示例接口
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int false "ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:成功"
// @Router /api/admin/example [get]
func ExampleHandler(c *gin.Context) {
    // ...
}
```

4. **更新 Swagger 文档**

```bash
cd backend
swag init
```

### 添加新的数据模型

1. **创建模型文件**

在 `backend/models/` 创建新的模型文件：

```go
package models

import (
    "time"
)

// Example 示例模型
type Example struct {
    ID        int        `json:"id" gorm:"primaryKey;autoIncrement"`
    Name      string     `json:"name" gorm:"size:100;not null;comment:名称"`
    Status    int        `json:"status" gorm:"not null;default:1;comment:状态"`
    CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (Example) TableName() string {
    return "examples"
}
```

2. **数据库迁移**

```bash
cd backend
go run main.go
```

程序启动时会自动执行 GORM 自动迁移。

## 前端开发

### 项目结构

```
frontend/src/
├── api/                # API 接口
├── router/             # 路由配置
├── stores/             # 状态管理
├── views/              # 页面组件
├── layout/             # 布局组件
├── utils/              # 工具函数
├── styles/             # 样式文件
├── main.js            # 主入口
└── App.vue            # 根组件
```

### 添加新的页面

1. **创建页面组件**

在 `frontend/src/views/` 创建新的页面组件：

```vue
<template>
  <div class="example-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>示例页面</span>
          <el-button type="primary">新增</el-button>
        </div>
      </template>
      <!-- 页面内容 -->
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

// 业务逻辑
</script>

<style scoped>
.example-container {
  padding: 20px;
}
</style>
```

2. **创建 API 接口**

在 `frontend/src/api/` 创建 API 接口文件：

```javascript
import request from '@/utils/request'

// 获取示例列表
export function getExampleList(params) {
  return request({
    url: '/admin/example',
    method: 'get',
    params
  })
}

// 创建示例
export function createExample(data) {
  return request({
    url: '/admin/example',
    method: 'post',
    data
  })
}
```

3. **配置路由**

在 `frontend/src/router/index.js` 添加路由：

```javascript
{
  path: 'example',
  name: 'Example',
  component: () => import('@/views/example/index.vue'),
  meta: { title: '示例页面', icon: 'Document' }
}
```

### 状态管理

使用 Pinia 进行状态管理：

```javascript
// stores/example.js
import { defineStore } from 'pinia'

export const useExampleStore = defineStore('example', {
  state: () => ({
    data: []
  }),
  
  getters: {
    getData: (state) => state.data
  },
  
  actions: {
    async loadData() {
      // 加载数据逻辑
    }
  }
})
```

### 样式规范

使用 SCSS 编写样式：

```scss
<style lang="scss" scoped>
.example-container {
  padding: 20px;
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>
```

## 常用命令

### 后端

```bash
# 安装依赖
go mod tidy

# 运行开发服务器
go run main.go

# 编译项目
go build -o admin-backend .

# 运行测试
go test ./...

# 更新 Swagger 文档
swag init

# 格式化代码
go fmt ./...
```

### 前端

```bash
# 安装依赖
npm install

# 运行开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview

# 代码检查
npm run lint
```

## 代码规范

### Go 代码规范

- 使用 `gofmt` 格式化代码
- 函数命名使用大驼峰
- 变量命名使用小驼峰
- 常量命名使用大写加下划线
- 导出函数/变量首字母大写
- 添加必要的注释

### Vue 代码规范

- 组件命名使用大驼峰
- 组件文件名使用小写加短横线
- 使用 Composition API
- 使用 `<script setup>` 语法
- 添加必要的注释

## 调试技巧

### 后端调试

1. 使用 `fmt.Println` 或 `log.Println` 打印日志
2. 使用 GoLand IDE 的调试功能
3. 查看 Gin 框架的请求日志

### 前端调试

1. 使用浏览器开发者工具
2. 使用 Vue DevTools 插件
3. 查看 Network 面板的请求信息
4. 使用 `console.log` 打印日志

## 性能优化

### 后端优化

- 使用 Redis 缓存热点数据
- 合理使用数据库索引
- 使用连接池
- 异步处理耗时操作

### 前端优化

- 路由懒加载
- 组件按需引入
- 使用虚拟滚动处理大量数据
- 图片懒加载

## 测试

### 单元测试

```bash
# 后端测试
cd backend
go test ./...

# 前端测试
cd frontend
npm run test
```

### 集成测试

编写集成测试用例，确保功能正常。

## 部署

### 后端部署

```bash
# 编译
go build -o admin-backend .

# 运行
./admin-backend
```

### 前端部署

```bash
# 构建
npm run build

# 部署到 Nginx
cp -r dist/* /path/to/nginx/html/
```

## 常见问题

### 1. CORS 错误

在后端添加 CORS 中间件：

```go
r.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    c.Next()
})
```

### 2. JWT Token 过期

检查 JWT 配置，确保 `JWT_EXPIRES` 设置合理。

### 3. 数据库连接失败

检查数据库连接配置，确保 MySQL 服务正常运行。

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 联系方式

- 项目主页: https://github.com/your-repo/mychat
- 问题反馈: https://github.com/your-repo/mychat/issues
