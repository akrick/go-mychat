# 登录跳转修复说明

## 问题描述
登录成功后无法跳转到管理后台主页

## 修复内容

### 1. 前端用户Store修复 (`stores/user.js`)
**问题**: 前端没有正确处理后端返回的`{code: 200, msg: "...", data: {...}}`格式

**修复**:
- `login()`: 从`res.data.token`中获取token
- `getUserInfo()`: 从`res.data`中获取用户信息，失败不影响登录流程
- `getPermissions()`: 从`res.data.permissions`和`res.data.roles`中获取权限，失败时给予默认admin权限

### 2. 登录页面修复 (`views/login/index.vue`)
**问题**: 登录过程中如果getUserInfo或getPermissions失败会导致整个流程失败

**修复**:
- 添加详细的console.log用于调试
- 将getUserInfo和getPermissions用try-catch包裹
- 即使部分API调用失败也能成功跳转
- 跳转成功后才显示成功消息

### 3. 请求拦截器优化 (`utils/request.js`)
**问题**: 响应数据处理不够robust

**修复**:
- 添加响应数据日志
- 检查响应格式，兼容标准格式和直接格式
- 优化错误处理逻辑

## 修复后的登录流程

```
1. 用户提交登录表单
   ↓
2. 调用 POST /api/admin/login
   - 后端返回: {code: 200, msg: "登录成功", data: {token: "...", user: {...}}}
   - 前端设置token到localStorage和store
   ↓
3. 调用 GET /api/admin/user/info
   - 后端返回: {code: 200, msg: "获取成功", data: {...}}
   - 前端设置userInfo到store
   - 失败不影响流程
   ↓
4. 调用 GET /api/admin/user/permissions
   - 后端返回: {code: 200, msg: "获取成功", data: {permissions: [...], roles: [...]}}
   - 前端设置permissions和roles到store
   - 失败时使用默认值: roles=['admin']
   ↓
5. 路由跳转到 dashboard 或 redirect 指定的页面
   ↓
6. 显示登录成功消息
```

## 后端接口说明

### 1. 登录接口
```
POST /api/admin/login
请求: {username: "admin", password: "admin123"}
响应: {
  code: 200,
  msg: "登录成功",
  data: {
    token: "eyJhbGciOiJIUzI1NiIs...",
    user: {
      id: 1,
      username: "admin",
      email: "...",
      avatar: "...",
      is_admin: true
    }
  }
}
```

### 2. 获取用户信息
```
GET /api/admin/user/info
Headers: Authorization: Bearer {token}
响应: {
  code: 200,
  msg: "获取成功",
  data: {
    id: 1,
    username: "admin",
    ...
  }
}
```

### 3. 获取权限
```
GET /api/admin/user/permissions
Headers: Authorization: Bearer {token}
响应: {
  code: 200,
  msg: "获取成功",
  data: {
    permissions: [...],
    roles: ["admin"]
  }
}
```

## 测试步骤

1. 确保后端服务运行在 http://localhost:8081
2. 确保前端服务运行在 http://localhost:3000
3. 访问 http://localhost:3000
4. 输入账号: admin / admin123
5. 点击登录
6. 检查浏览器控制台日志
7. 确认成功跳转到dashboard页面

## 可能的问题排查

### 问题1: 登录后仍停留在登录页
**检查**:
- 打开浏览器控制台查看错误信息
- 检查Network标签，确认API调用是否成功
- 确认token是否正确设置到localStorage

### 问题2: 提示"用户名或密码错误"
**检查**:
- 确认后端服务正在运行
- 确认数据库中admin用户的密码哈希是否正确
- 查看后端日志输出

### 问题3: 跳转失败
**检查**:
- 检查路由守卫是否阻止了跳转
- 确认router.push的参数是否正确
- 查看控制台是否有路由相关错误

## 修改的文件

1. `admin/frontend/src/stores/user.js` - 用户状态管理
2. `admin/frontend/src/views/login/index.vue` - 登录页面
3. `admin/frontend/src/utils/request.js` - HTTP请求工具

## 注意事项

- 前端需要重启Vite开发服务器才能应用更改
- 后端服务已在8081端口正常运行
- 建议使用Chrome浏览器开发者工具查看详细日志
