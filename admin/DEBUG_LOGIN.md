# 调试登录问题

## ✅ 当前配置

### 数据库状态
```
用户名: admin
ID: 2
IsAdmin: 1 ✅
Status: 1 ✅
```

### 后端路由
```
POST /api/admin/login
```

### 登录逻辑

**临时特殊逻辑**：
- 当输入密码为 `admin123` 时：
  1. 生成正确的 bcrypt 哈希
  2. 更新到数据库
  3. 允许登录
- 其他密码：使用正常的 bcrypt 验证

## 🔍 调试步骤

### 1. 检查后端日志
后端会输出以下日志：
```
=== 登录请求开始 ===
接收到的请求体: {username: "admin", password: "xxx"}
用户名: admin, 密码: xxx
找到用户: ID=2, Username=admin, IsAdmin=true, Status=1
检测到 admin123，生成新的密码哈希
生成的哈希: $2a$10$...
密码哈希已更新到数据库
登录成功，生成的Token: xxx
=== 登录请求结束 ===
```

### 2. 测试登录
```
用户名: admin
密码: admin123
```

### 3. 预期结果
- ✅ 第一次登录成功（使用特殊逻辑）
- ✅ 数据库密码被更新为正确的 bcrypt 哈希
- ✅ 返回 token 和用户信息
- ✅ 前端跳转到首页

### 4. 后续登录
- 密码已更新为正确的哈希
- 使用正常的 bcrypt 验证
- 无需特殊逻辑

## 📝 常见问题排查

### 问题1: 401 Unauthorized
**可能原因**:
- 后端服务未启动
- 数据库连接失败
- 用户不存在或密码错误
- 请求格式不正确

**排查方法**:
1. 检查后端日志输出
2. 验证数据库连接：`mysql -uroot -p123456 -e "USE mychat; SELECT * FROM users;"`
3. 检查前端网络请求（F12 -> Network）

### 问题2: 后端无日志输出
**解决方法**:
- 确认后端服务正在运行
- 检查端口 8081 是否被占用
- 重新启动后端服务

### 问题3: 密码一直错误
**解决方法**:
- 第一次使用 `admin123` 登录会自动更新密码
- 如果失败，直接修改数据库：
  ```sql
  -- 查看当前密码
  SELECT username, password FROM users WHERE username='admin';

  -- 如果需要重置，可以手动更新
  UPDATE users SET password='任意值' WHERE username='admin';
  ```
  然后使用新密码登录（第一次会自动生成哈希）

## 🎯 推荐测试流程

1. 启动后端：`cd admin/backend && go run main.go`
2. 启动前端：`cd admin/frontend && npm run dev`
3. 访问：http://localhost:3000
4. 登录：`admin` / `admin123`
5. 查看后端日志确认流程
6. 退出后重新登录（使用正常密码验证）

## ⚠️ 注意事项

- 临时逻辑仅用于第一次登录
- 第一次登录成功后，密码会被正确哈希
- 后续登录使用标准密码验证
- 生产环境应移除临时逻辑，使用预置的密码哈希
