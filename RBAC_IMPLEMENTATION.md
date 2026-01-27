# MyChat RBAC权限管理系统实现文档

## 完成时间
2026-01-27

## 系统架构

### RBAC模型
```
用户(User) ---- 多对多 ----> 角色(Role)
  |                                   |
  |                                   | 一对多
  |                                 权限(Permission)
  |                                   |
  └----------------- 一对多 ------------─┘
```

### 核心概念
1. **用户**: 系统中的实际用户
2. **角色**: 权限的集合，简化权限管理
3. **权限**: 最小权限单元，可以是菜单、按钮、API接口
4. **用户角色关联**: 用户与角色的多对多关系
5. **角色权限关联**: 角色与权限的多对多关系

## 数据库设计

### 1. 角色表 (roles)
```sql
CREATE TABLE `roles` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(50) NOT NULL COMMENT '角色名称',
  `code` VARCHAR(50) NOT NULL COMMENT '角色代码，如：admin、user',
  `description` VARCHAR(255) COMMENT '角色描述',
  `sort` INT DEFAULT 0 COMMENT '排序',
  `status` TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**字段说明**:
- `id`: 主键，自增
- `name`: 角色名称，如：超级管理员、管理员
- `code`: 角色代码，系统内部使用，如：super_admin、admin
- `description`: 角色描述信息
- `sort`: 排序字段，用于角色列表排序
- `status`: 状态，0-禁用，1-启用

### 2. 权限表 (permissions)
```sql
CREATE TABLE `permissions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_id` INT UNSIGNED DEFAULT 0 COMMENT '父权限ID，用于构建树形结构',
  `name` VARCHAR(50) NOT NULL COMMENT '权限名称',
  `code` VARCHAR(100) NOT NULL COMMENT '权限代码，唯一标识',
  `type` VARCHAR(20) DEFAULT 'menu' COMMENT '类型：menu-菜单，button-按钮，api-接口',
  `path` VARCHAR(255) COMMENT '路由路径',
  `icon` VARCHAR(50) COMMENT '图标名称',
  `component` VARCHAR(255) COMMENT '前端组件路径',
  `sort` INT DEFAULT 0 COMMENT '排序',
  `status` TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**字段说明**:
- `id`: 主键，自增
- `parent_id`: 父权限ID，用于构建树形结构（菜单树、权限树）
- `name`: 权限名称，如：用户管理
- `code`: 权限代码，格式：模块:操作，如：system:user:list
- `type`: 权限类型
  - `menu`: 菜单权限，对应前端路由
  - `button`: 按钮权限，对应页面操作按钮
  - `api`: 接口权限，对应后端API
- `path`: 路由路径，如：/system/user
- `icon`: 图标名称，如：User
- `component`: 前端组件路径，如：views/system/user/index
- `sort`: 排序字段
- `status`: 状态

**权限代码命名规范**:
```
模块:操作:子操作

示例：
- system:user:list - 系统模块：用户管理：查看列表
- system:user:create - 系统模块：用户管理：创建用户
- business:order:update - 业务模块：订单管理：更新订单
- finance:withdraw:approve - 财务模块：提现审核：审核通过
```

### 3. 用户角色关联表 (user_roles)
```sql
CREATE TABLE `user_roles` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
  `role_id` INT UNSIGNED NOT NULL COMMENT '角色ID',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**字段说明**:
- `id`: 主键，自增
- `user_id`: 用户ID，外键关联users表
- `role_id`: 角色ID，外键关联roles表
- `created_at`: 创建时间

**约束**:
- 联合唯一索引 `(user_id, role_id)`，一个用户不能重复分配同一个角色
- 一个用户可以有多个角色
- 一个角色可以分配给多个用户

### 4. 角色权限关联表 (role_permissions)
```sql
CREATE TABLE `role_permissions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `role_id` INT UNSIGNED NOT NULL COMMENT '角色ID',
  `permission_id` INT UNSIGNED NOT NULL COMMENT '权限ID',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**字段说明**:
- `id`: 主键，自增
- `role_id`: 角色ID，外键关联roles表
- `permission_id`: 权限ID，外键关联permissions表
- `created_at`: 创建时间

**约束**:
- 联合唯一索引 `(role_id, permission_id)`，一个角色不能重复分配同一个权限
- 一个角色可以有多个权限
- 一个权限可以分配给多个角色

## 后端实现

### 数据模型 (models/rbac.go)
```go
// Role 角色表
type Role struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:角色名称" json:"name"`
    Code        string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:角色代码" json:"code"`
    Description string    `gorm:"type:varchar(255);comment:描述" json:"description"`
    Sort        int       `gorm:"default:0;comment:排序" json:"sort"`
    Status      int       `gorm:"default:1;comment:状态:0-禁用,1-启用" json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Permission 权限表
type Permission struct {
    ID          uint         `gorm:"primaryKey" json:"id"`
    ParentID    uint         `gorm:"default:0;index;comment:父权限ID" json:"parent_id"`
    Name        string       `gorm:"type:varchar(50);not null;comment:权限名称" json:"name"`
    Code        string       `gorm:"type:varchar(100);uniqueIndex;not null;comment:权限代码" json:"code"`
    Type        string       `gorm:"type:varchar(20);default:menu;comment:类型:menu-菜单,button-按钮,api-接口" json:"type"`
    Path        string       `gorm:"type:varchar(255);comment:路由路径" json:"path"`
    Icon        string       `gorm:"type:varchar(50);comment:图标" json:"icon"`
    Component   string       `gorm:"type:varchar(255);comment:组件路径" json:"component"`
    Sort        int          `gorm:"default:0;comment:排序" json:"sort"`
    Status      int          `gorm:"default:1;comment:状态:0-禁用,1-启用" json:"status"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
    Children    []Permission `gorm:"-" json:"children,omitempty"`
}
```

### API接口实现

#### 1. 角色管理接口 (handlers/rbac.go)

**获取角色列表**
```go
GET /api/admin/roles?page=1&page_size=20&name=xxx
```

**创建角色**
```go
POST /api/admin/roles
Body: {
  "name": "角色名称",
  "code": "role_code",
  "description": "角色描述",
  "sort": 0,
  "status": 1
}
```

**更新角色**
```go
PUT /api/admin/roles/:id
Body: {
  "name": "角色名称",
  "description": "角色描述",
  "sort": 0,
  "status": 1
}
```

**删除角色**
```go
DELETE /api/admin/roles/:id
```

**获取角色权限**
```go
GET /api/admin/roles/:id/permissions
返回: [Permission对象数组]
```

**分配权限**
```go
POST /api/admin/roles/:id/permissions
Body: {
  "permission_ids": [1, 2, 3, 4, 5]
}
```

**获取角色用户列表**
```go
GET /api/admin/roles/:id/users?page=1&page_size=20
返回: {
  "users": [User对象数组],
  "total": 总数
}
```

#### 2. 权限管理接口 (handlers/rbac.go)

**获取权限树**
```go
GET /api/admin/permissions/tree
返回: [
  {
    "id": 1,
    "name": "系统管理",
    "code": "system",
    "type": "menu",
    "children": [...]
  }
]
```

**获取权限列表**
```go
GET /api/admin/permissions?page=1&page_size=20&name=xxx
```

**创建权限**
```go
POST /api/admin/permissions
Body: {
  "parent_id": 0,
  "name": "权限名称",
  "code": "system:user:list",
  "type": "menu",
  "path": "/system/user",
  "icon": "User",
  "component": "views/system/user/index",
  "sort": 1,
  "status": 1
}
```

**更新权限**
```go
PUT /api/admin/permissions/:id
Body: { 同创建 }
```

**删除权限**
```go
DELETE /api/admin/permissions/:id
```

#### 3. 菜单管理接口 (handlers/menu.go)

菜单本质上是类型为`menu`的权限，复用permissions表。

**获取菜单树**
```go
GET /api/admin/menus/tree
返回: 只包含type=menu的权限树
```

**获取菜单列表**
```go
GET /api/admin/menus
返回: 只包含type=menu的权限列表
```

**创建菜单**
```go
POST /api/admin/menus
Body: {
  "parent_id": 0,
  "name": "菜单名称",
  "code": "system:user:list",
  "type": "menu",
  "path": "/system/user",
  "icon": "User",
  "component": "views/system/user/index",
  "sort": 1,
  "status": 1
}
```

**更新菜单**
```go
PUT /api/admin/menus/:id
Body: { 同创建 }
```

**删除菜单**
```go
DELETE /api/admin/menus/:id
约束: 不能删除有子菜单的菜单
```

### 权限树构建算法

```go
func buildPermissionTree(permissions []Permission, parentID uint) []Permission {
    var tree []Permission
    for _, p := range permissions {
        if p.ParentID == parentID {
            children := buildPermissionTree(permissions, p.ID)
            if len(children) > 0 {
                p.Children = children
            }
            tree = append(tree, p)
        }
    }
    return tree
}
```

## 前端实现

### 1. 角色管理页面 (views/system/roles/index.vue)

**功能特性**:
- ✅ 角色列表展示（分页）
- ✅ 角色名称搜索
- ✅ 角色状态筛选
- ✅ 新增/编辑/删除角色
- ✅ 为角色分配权限（树形选择器）
- ✅ 查看角色下的用户
- ✅ 角色代码自动生成
- ✅ 角色排序
- ✅ 角色启用/禁用

**权限分配**:
- 树形权限展示
- 支持按类型筛选（菜单/按钮/接口）
- 支持全选/半选
- 展开/收起所有节点

### 2. 权限管理页面 (views/system/permissions/index.vue)

**功能特性**:
- ✅ 权限树形展示
- ✅ 权限类型显示（菜单/按钮/接口）
- ✅ 新增/编辑/删除权限
- ✅ 支持新增子权限
- ✅ 权限代码自动生成
- ✅ 权限排序
- ✅ 权限启用/禁用
- ✅ 图标选择器
- ✅ 展开/收起全部

**权限类型说明**:
- **菜单**: 对应前端路由，显示在导航栏
- **按钮**: 对应页面操作按钮，如新增、编辑、删除
- **接口**: 对应后端API接口，用于接口权限验证

### 3. 菜单管理页面 (views/system/menus/index.vue)

**功能特性**:
- ✅ 菜单树形展示
- ✅ 新增/编辑/删除菜单
- ✅ 支持新增子菜单
- ✅ 上级菜单选择
- ✅ 菜单图标选择
- ✅ 菜单排序
- ✅ 菜单启用/禁用

**菜单配置**:
- 菜单名称
- 上级菜单（支持树形选择）
- 路由路径
- 组件路径
- 图标
- 排序
- 状态

## 权限验证机制

### 1. JWT中间件 (middleware/auth.go)

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"code": 401, "msg": "未登录"})
            c.Abort()
            return
        }

        // 解析JWT token
        claims, err := ParseToken(token)
        if err != nil {
            c.JSON(401, gin.H{"code": 401, "msg": "token无效"})
            c.Abort()
            return
        }

        // 获取用户角色
        roles := getUserRoles(claims.UserID)
        // 获取角色权限
        permissions := getRolePermissions(roles)

        // 将用户信息存入context
        c.Set("user_id", claims.UserID)
        c.Set("roles", roles)
        c.Set("permissions", permissions)

        c.Next()
    }
}
```

### 2. 权限检查函数

```go
func HasPermission(requiredPermission string, userPermissions []Permission) bool {
    for _, perm := range userPermissions {
        if perm.Code == requiredPermission {
            return true
        }
    }
    return false
}
```

### 3. 前端权限指令

```javascript
// main.js 注册权限指令
app.directive('permission', {
  mounted(el, binding) {
    const { value } = binding
    const permissions = useUserStore().permissions || []
    
    if (value && !permissions.includes(value)) {
      el.parentNode && el.parentNode.removeChild(el)
    }
  }
})

// 使用示例
<el-button v-permission="'system:user:create'">新增用户</el-button>
```

### 4. 前端路由守卫

```javascript
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  const requiredPermission = to.meta.permission

  if (requiredPermission && !userStore.hasPermission(requiredPermission)) {
    next({ path: '/' })
  } else {
    next()
  }
})
```

## 默认角色和权限

### 1. 预定义角色

| 角色 | 代码 | 说明 | 权限范围 |
|------|------|------|----------|
| 超级管理员 | super_admin | 拥有系统所有权限 | 全部权限 |
| 管理员 | admin | 拥有大部分管理权限 | 除超级权限外的所有权限 |
| 咨询师 | counselor | 咨询师角色 | 业务模块权限 |
| 普通用户 | user | 普通用户角色 | 基本查询权限 |

### 2. 权限模块划分

| 模块 | 权限代码前缀 | 说明 |
|------|--------------|------|
| 系统管理 | system: | 用户、角色、权限、菜单、日志等 |
| 业务管理 | business: | 订单、聊天等 |
| 财务管理 | finance: | 提现、统计、报表等 |
| 低代码平台 | lowcode: | 表单、页面、数据等 |

## 使用示例

### 1. 为用户分配角色

```sql
-- 用户ID=1，分配超级管理员角色
INSERT INTO user_roles (user_id, role_id) VALUES (1, 1);

-- 用户ID=2，分配管理员角色
INSERT INTO user_roles (user_id, role_id) VALUES (2, 2);
```

### 2. 为角色分配权限

```sql
-- 角色ID=2，分配权限ID列表
INSERT INTO role_permissions (role_id, permission_id) VALUES
(2, 1), (2, 2), (2, 3), (2, 4), (2, 5);
```

### 3. 检查用户权限

```go
// 在Handler中检查权限
if !HasPermission("system:user:create", userPermissions) {
    c.JSON(403, gin.H{"code": 403, "msg": "权限不足"})
    return
}
```

### 4. 前端控制按钮显示

```vue
<template>
  <div>
    <el-button v-permission="'system:user:create'">新增用户</el-button>
    <el-button v-permission="'system:user:update'">编辑用户</el-button>
    <el-button v-permission="'system:user:delete'">删除用户</el-button>
  </div>
</template>
```

## 初始化数据库

执行初始化SQL脚本：

```bash
mysql -u root -p mychat < api/init_rbac.sql
```

该脚本会：
1. ✅ 创建4张RBAC表
2. ✅ 插入4个预定义角色
3. ✅ 插入完整的权限树（50+权限）
4. ✅ 构建权限树结构
5. ✅ 为超级管理员分配所有权限
6. ✅ 为管理员分配大部分权限
7. ✅ 为普通用户分配基本权限

## API接口清单

### 角色管理
- `GET /api/admin/roles` - 获取角色列表
- `POST /api/admin/roles` - 创建角色
- `PUT /api/admin/roles/:id` - 更新角色
- `DELETE /api/admin/roles/:id` - 删除角色
- `GET /api/admin/roles/:id/permissions` - 获取角色权限
- `POST /api/admin/roles/:id/permissions` - 分配权限
- `GET /api/admin/roles/:id/users` - 获取角色用户

### 权限管理
- `GET /api/admin/permissions/tree` - 获取权限树
- `GET /api/admin/permissions` - 获取权限列表
- `POST /api/admin/permissions` - 创建权限
- `PUT /api/admin/permissions/:id` - 更新权限
- `DELETE /api/admin/permissions/:id` - 删除权限

### 菜单管理
- `GET /api/admin/menus/tree` - 获取菜单树
- `GET /api/admin/menus` - 获取菜单列表
- `POST /api/admin/menus` - 创建菜单
- `PUT /api/admin/menus/:id` - 更新菜单
- `DELETE /api/admin/menus/:id` - 删除菜单

## 测试建议

### 1. 角色管理测试
- [ ] 创建新角色
- [ ] 编辑角色信息
- [ ] 为角色分配权限
- [ ] 查看角色下的用户
- [ ] 禁用/启用角色
- [ ] 删除角色

### 2. 权限管理测试
- [ ] 创建顶级权限（菜单）
- [ ] 创建子权限（按钮）
- [ ] 编辑权限信息
- [ ] 禁用/启用权限
- [ ] 删除权限（检查级联）

### 3. 菜单管理测试
- [ ] 创建一级菜单
- [ ] 创建二级菜单
- [ ] 设置菜单图标
- [ ] 配置路由路径
- [ ] 配置组件路径
- [ ] 调整菜单排序

### 4. 权限验证测试
- [ ] 为用户分配角色
- [ ] 验证菜单显示
- [ ] 验证按钮显示
- [ ] 验证接口访问权限
- [ ] 测试无权限时的行为

## 注意事项

1. **权限代码唯一性**: 权限代码在系统中必须唯一，建议遵循命名规范
2. **删除限制**: 删除权限时，需要检查是否有子权限，避免破坏树结构
3. **角色引用**: 删除角色前，检查是否有用户使用该角色
4. **缓存优化**: 权限数据建议缓存，减少数据库查询
5. **前端同步**: 前端路由需要与权限菜单保持一致
6. **权限更新**: 角色权限变化后，需要用户重新登录生效
7. **超级管理员**: 超级管理员拥有所有权限，可以跳过部分权限检查

## 扩展建议

1. **数据权限**: 可以扩展数据行级权限，控制用户能看到的数据范围
2. **权限组**: 可以按模块分组权限，便于批量分配
3. **临时权限**: 支持为用户临时授予某个权限，到期自动失效
4. **审批流程**: 重要权限变更需要审批
5. **权限审计**: 记录所有权限变更日志
6. **权限依赖**: 设置权限之间的依赖关系，自动关联依赖权限

## 总结

MyChat RBAC权限管理系统实现了：
- ✅ 完整的角色管理功能
- ✅ 完整的权限管理功能
- ✅ 完整的菜单管理功能
- ✅ 灵活的权限分配机制
- ✅ 树形权限结构
- ✅ 前后端权限验证
- ✅ 完善的初始化数据
- ✅ 详细的API接口文档

系统具备企业级权限管理能力，可满足复杂的业务需求！
