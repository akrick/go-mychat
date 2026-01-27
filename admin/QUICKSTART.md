# MyChat 管理后台快速启动指南

## 前提条件

确保已安装:
- Go 1.21+
- Node.js 16+
- MySQL 5.7+
- Redis 6.0+

## 启动步骤

### 方式一: 使用批处理脚本 (推荐 Windows 用户)

#### 1. 启动后端服务

双击运行 `run-backend.bat` 文件,或在命令行中执行:

```cmd
cd d:\gospace\src\akrick.com\mychat\admin
run-backend.bat
```

后端服务将启动在: **http://localhost:8081**

#### 2. 启动前端服务

打开新的命令行窗口,双击运行 `run-frontend.bat` 文件,或执行:

```cmd
cd d:\gospace\src\akrick.com\mychat\admin
run-frontend.bat
```

前端服务将启动在: **http://localhost:3000**

### 方式二: 手动启动

#### 1. 启动后端服务

打开第一个命令行窗口:

```cmd
cd d:\gospace\src\akrick.com\mychat\admin\backend
go mod tidy
go run main.go
```

#### 2. 启动前端服务

打开第二个命令行窗口:

```cmd
cd d:\gospace\src\akrick.com\mychat\admin\frontend
npm install
npm run dev
```

### 方式三: 使用 PowerShell (推荐开发者)

#### 启动后端

```powershell
cd d:\gospace\src\akrick.com\mychat\admin\backend
go mod tidy
go run main.go
```

#### 启动前端

```powershell
cd d:\gospace\src\akrick.com\mychat\admin\frontend
npm install
npm run dev
```

## 验证服务是否正常

### 后端验证

在浏览器访问:
- http://localhost:8081

应该看到服务正常运行的提示

### 前端验证

在浏览器访问:
- http://localhost:3000

应该看到登录页面

## 默认登录账号

```
用户名: admin
密码: admin123
```

## 常见问题

### 1. 端口被占用

如果 8081 或 3000 端口被占用:

#### 修改后端端口

编辑 `backend/main.go` 文件,找到:
```go
log.Fatal(r.Run(":8081"))
```
改为其他端口,例如:
```go
log.Fatal(r.Run(":8881"))
```

#### 修改前端端口

编辑 `frontend/vite.config.js` 文件:
```javascript
server: {
    port: 3001,  // 改为其他端口
    ...
}
```

### 2. 数据库连接失败

检查:
- MySQL 服务是否启动
- `backend/database/db.go` 中的连接配置是否正确
- 数据库 `mychat` 是否已创建

### 3. Go 依赖下载失败

如果 `go mod tidy` 失败,尝试:

```cmd
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
```

### 4. npm 安装失败

如果 `npm install` 失败,尝试:

```cmd
npm config set registry https://registry.npmmirror.com
npm install
```

### 5. CORS 错误

前端请求后端时出现 CORS 错误,检查:
- 后端 CORS 中间件配置
- 前端 `vite.config.js` 中的 proxy 配置

## 服务架构

```
浏览器 → 前端 (localhost:3000)
            ↓
        代理转发
            ↓
        后端 (localhost:8081)
            ↓
    MySQL + Redis
```

## 停止服务

在对应的命令行窗口按 `Ctrl + C` 停止服务。

## 开发模式

开发模式下:
- 前端支持热重载 (修改代码自动刷新)
- 后端需要手动重启

## 生产部署

### 后端编译

```cmd
cd d:\gospace\src\akrick.com\mychat\admin\backend
go build -o admin-backend.exe .
```

运行编译后的程序:
```cmd
admin-backend.exe
```

### 前端构建

```cmd
cd d:\gospace\src\akrick.com\mychat\admin\frontend
npm run build
```

构建产物在 `dist` 目录,可以部署到 Nginx 等静态服务器。

## 技术支持

遇到问题?
1. 查看控制台错误信息
2. 检查日志文件
3. 参考 `DEVELOPMENT.md` 开发指南
4. 提交 Issue

---

**祝您使用愉快!** 🚀
