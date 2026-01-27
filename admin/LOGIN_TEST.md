# 管理员登录测试

## ✅ 已更新 admin 用户密码

### 用户信息
- **用户名**: `admin`
- **密码**: `admin123`
- **角色**: 管理员 (is_admin = 1)
- **状态**: 正常 (status = 1)

### 密码哈希信息
- **加密方式**: bcrypt
- **Cost**: 10
- **哈希值**: `$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy`

### 测试步骤

1. **启动后端服务**
   ```bash
   cd admin/backend
   go run main.go
   ```

2. **启动前端服务**
   ```bash
   cd admin/frontend
   npm run dev
   ```

3. **访问登录页面**
   - 前端地址: http://localhost:3000
   - 后端地址: http://localhost:8081

4. **输入登录信息**
   - 用户名: `admin`
   - 密码: `admin123`

5. **预期结果**
   - ✅ 登录成功
   - ✅ 跳转到数据看板
   - ✅ 显示管理员权限

### 已修复的问题

1. ✅ 修复了前端响应拦截器的错误处理
2. ✅ 更新了 admin 用户的密码为 `admin123`
3. ✅ 恢复了密码验证功能（移除了临时跳过代码）
4. ✅ 确保 is_admin = 1（管理员权限）

### 注意事项

- 如果密码验证失败，请检查：
  - 数据库中的密码哈希是否正确
  - 前端发送的密码是否正确
  - 后端的密码验证逻辑是否启用

- 默认情况下，系统会自动创建 admin/admin123 账号（在 `database/db.go` 中的 `createDefaultAdmin()` 函数）
