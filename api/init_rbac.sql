-- MyChat RBAC权限管理系统初始化脚本
-- 包含角色、权限、用户角色、角色权限表

-- ============================================
-- 角色表 (roles)
-- ============================================
CREATE TABLE IF NOT EXISTS `roles` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` VARCHAR(50) NOT NULL COMMENT '角色名称',
  `code` VARCHAR(50) NOT NULL COMMENT '角色代码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '角色描述',
  `sort` INT DEFAULT 0 COMMENT '排序',
  `status` TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- ============================================
-- 权限表 (permissions)
-- ============================================
CREATE TABLE IF NOT EXISTS `permissions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '权限ID',
  `parent_id` INT UNSIGNED DEFAULT 0 COMMENT '父权限ID',
  `name` VARCHAR(50) NOT NULL COMMENT '权限名称',
  `code` VARCHAR(100) NOT NULL COMMENT '权限代码',
  `type` VARCHAR(20) DEFAULT 'menu' COMMENT '类型：menu-菜单，button-按钮，api-接口',
  `path` VARCHAR(255) DEFAULT NULL COMMENT '路由路径',
  `icon` VARCHAR(50) DEFAULT NULL COMMENT '图标',
  `component` VARCHAR(255) DEFAULT NULL COMMENT '组件路径',
  `sort` INT DEFAULT 0 COMMENT '排序',
  `status` TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- ============================================
-- 用户角色关联表 (user_roles)
-- ============================================
CREATE TABLE IF NOT EXISTS `user_roles` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
  `role_id` INT UNSIGNED NOT NULL COMMENT '角色ID',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- ============================================
-- 角色权限关联表 (role_permissions)
-- ============================================
CREATE TABLE IF NOT EXISTS `role_permissions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `role_id` INT UNSIGNED NOT NULL COMMENT '角色ID',
  `permission_id` INT UNSIGNED NOT NULL COMMENT '权限ID',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- ============================================
-- 初始化默认角色
-- ============================================

-- 超级管理员
INSERT INTO `roles` (`name`, `code`, `description`, `sort`, `status`) VALUES
('超级管理员', 'super_admin', '拥有系统所有权限', 1, 1),
('管理员', 'admin', '拥有大部分管理权限', 2, 1),
('咨询师', 'counselor', '咨询师角色', 3, 1),
('普通用户', 'user', '普通用户角色', 4, 1)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- ============================================
-- 初始化权限数据
-- ============================================

-- 系统管理模块
INSERT INTO `permissions` (`name`, `code`, `type`, `path`, `icon`, `component`, `sort`, `status`) VALUES
-- 系统管理 (一级菜单)
('系统管理', 'system', 'menu', '/system', 'Setting', NULL, 1, 1),
-- 用户管理 (二级菜单)
('用户管理', 'system:user:list', 'menu', '/system/user', 'User', 'views/system/user/index', 11, 1),
('新增用户', 'system:user:create', 'button', '', '', NULL, 111, 1),
('编辑用户', 'system:user:update', 'button', '', '', NULL, 112, 1),
('删除用户', 'system:user:delete', 'button', '', '', NULL, 113, 1),
('重置密码', 'system:user:reset', 'button', '', '', NULL, 114, 1),
-- 角色管理 (二级菜单)
('角色管理', 'system:role:list', 'menu', '/system/roles', 'Avatar', 'views/system/roles/index', 12, 1),
('新增角色', 'system:role:create', 'button', '', '', NULL, 121, 1),
('编辑角色', 'system:role:update', 'button', '', '', NULL, 122, 1),
('删除角色', 'system:role:delete', 'button', '', '', NULL, 123, 1),
('分配权限', 'system:role:assign', 'button', '', '', NULL, 124, 1),
-- 权限管理 (二级菜单)
('权限管理', 'system:permission:list', 'menu', '/system/permissions', 'Key', 'views/system/permissions/index', 13, 1),
('新增权限', 'system:permission:create', 'button', '', '', NULL, 131, 1),
('编辑权限', 'system:permission:update', 'button', '', '', NULL, 132, 1),
('删除权限', 'system:permission:delete', 'button', '', '', NULL, 133, 1),
-- 菜单管理 (二级菜单)
('菜单管理', 'system:menu:list', 'menu', '/system/menus', 'Menu', 'views/system/menus/index', 14, 1),
('新增菜单', 'system:menu:create', 'button', '', '', NULL, 141, 1),
('编辑菜单', 'system:menu:update', 'button', '', '', NULL, 142, 1),
('删除菜单', 'system:menu:delete', 'button', '', '', NULL, 143, 1),
-- 咨询师管理 (二级菜单)
('咨询师管理', 'system:counselor:list', 'menu', '/system/counselor', 'UserFilled', 'views/system/counselor/index', 15, 1),
('新增咨询师', 'system:counselor:create', 'button', '', '', NULL, 151, 1),
('编辑咨询师', 'system:counselor:update', 'button', '', '', NULL, 152, 1),
('删除咨询师', 'system:counselor:delete', 'button', '', '', NULL, 153, 1),
-- 系统日志 (二级菜单)
('系统日志', 'system:log:list', 'menu', '/system/logs', 'Document', 'views/system/logs/index', 16, 1),
-- 在线用户 (二级菜单)
('在线用户', 'system:online:list', 'menu', '/system/online', 'User', 'views/system/online/index', 17, 1),
-- 系统配置 (二级菜单)
('系统配置', 'system:config:manage', 'menu', '/system/config', 'Setting', 'views/system/config/index', 18, 1)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- 业务管理模块
INSERT INTO `permissions` (`name`, `code`, `type`, `path`, `icon`, `component`, `sort`, `status`) VALUES
-- 业务管理 (一级菜单)
('业务管理', 'business', 'menu', '/business', 'Grid', NULL, 2, 1),
-- 订单管理 (二级菜单)
('订单管理', 'business:order:list', 'menu', '/business/order', 'Document', 'views/business/order/index', 21, 1),
('查看订单', 'business:order:view', 'button', '', '', NULL, 211, 1),
('更新订单状态', 'business:order:update', 'button', '', '', NULL, 212, 1),
('导出订单', 'business:order:export', 'button', '', '', NULL, 213, 1),
-- 聊天记录 (二级菜单)
('聊天记录', 'business:chat:list', 'menu', '/business/chat', 'ChatDotRound', 'views/business/chat/index', 22, 1),
('查看会话', 'business:chat:view', 'button', '', '', NULL, 221, 1),
('查看消息', 'business:chat:message:view', 'button', '', '', NULL, 222, 1),
('搜索消息', 'business:chat:search', 'button', '', '', NULL, 223, 1)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- 财务管理模块
INSERT INTO `permissions` (`name`, `code`, `type`, `path`, `icon`, `component`, `sort`, `status`) VALUES
-- 财务管理 (一级菜单)
('财务管理', 'finance', 'menu', '/finance', 'Wallet', NULL, 3, 1),
-- 提现审核 (二级菜单)
('提现审核', 'finance:withdraw:list', 'menu', '/finance/withdraw', 'Wallet', 'views/finance/withdraw/index', 31, 1),
('审核通过', 'finance:withdraw:approve', 'button', '', '', NULL, 311, 1),
('审核拒绝', 'finance:withdraw:reject', 'button', '', '', NULL, 312, 1),
-- 财务统计 (二级菜单)
('财务统计', 'finance:statistics:view', 'menu', '/finance/statistics', 'DataLine', 'views/finance/statistics/index', 32, 1),
-- 财务报表 (二级菜单)
('财务报表', 'finance:reports:view', 'menu', '/finance/reports', 'Document', 'views/finance/reports/index', 33, 1)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- 低代码平台模块
INSERT INTO `permissions` (`name`, `code`, `type`, `path`, `icon`, `component`, `sort`, `status`) VALUES
-- 低代码平台 (一级菜单)
('低代码平台', 'lowcode', 'menu', '/lowcode', 'Edit', NULL, 4, 1),
-- 表单设计 (二级菜单)
('表单设计', 'lowcode:form:design', 'menu', '/lowcode/forms', 'Edit', 'views/lowcode/forms/index', 41, 1),
('页面设计', 'lowcode:page:design', 'menu', '/lowcode/pages', 'Grid', 'views/lowcode/pages/index', 42, 1),
('数据管理', 'lowcode:data:manage', 'menu', '/lowcode/data', 'Database', 'views/lowcode/data/index', 43, 1)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- ============================================
-- 初始化菜单树结构（设置父ID）
-- ============================================

-- 设置系统管理子菜单的parent_id
UPDATE permissions SET parent_id = (SELECT id FROM (SELECT id FROM permissions WHERE code = 'system') AS tmp) WHERE code IN ('system:user:list', 'system:role:list', 'system:permission:list', 'system:menu:list', 'system:counselor:list', 'system:log:list', 'system:online:list', 'system:config:manage');

-- 设置用户管理子按钮的parent_id
UPDATE permissions SET parent_id = (SELECT id FROM (SELECT id FROM permissions WHERE code = 'system:user:list') AS tmp) WHERE code IN ('system:user:create', 'system:user:update', 'system:user:delete', 'system:user:reset');

-- 设置角色管理子按钮的parent_id
UPDATE permissions SET parent_id = (SELECT id FROM (SELECT id FROM permissions WHERE code = 'system:role:list') AS tmp) WHERE code IN ('system:role:create', 'system:role:update', 'system:role:delete', 'system:role:assign');

-- 设置权限管理子按钮的parent_id
UPDATE permissions SET parent_id = (SELECT id FROM (SELECT id FROM permissions WHERE code = 'system:permission:list') AS tmp) WHERE code IN ('system:permission:create', 'system:permission:update', 'system:permission:delete');

-- 设置菜单管理子按钮的parent_id
UPDATE permissions SET parent_id = (SELECT id FROM (SELECT id FROM permissions WHERE code = 'system:menu:list') AS tmp) WHERE code IN ('system:menu:create', 'system:menu:update', 'system:menu:delete');

-- 设置咨询师管理子按钮的parent_id
UPDATE permissions SET parent_id = (SELECT id FROM (SELECT id FROM permissions WHERE code = 'system:counselor:list') AS tmp) WHERE code IN ('system:counselor:create', 'system:counselor:update', 'system:counselor:delete');

-- 设置业务管理子菜单的parent_id
UPDATE permissions SET parent_id = (SELECT id FROM (SELECT id FROM permissions WHERE code = 'business') AS tmp) WHERE code IN ('business:order:list', 'business:chat:list');

-- 设置订单管理子按钮的parent_id
UPDATE permissions SET parent_id = (SELECT id FROM (SELECT id FROM permissions WHERE code = 'business:order:list') AS tmp) WHERE code IN ('business:order:view', 'business:order:update', 'business:order:export');

-- 设置聊天记录子按钮的parent_id
UPDATE permissions SET parent_id = (SELECT id FROM (SELECT id FROM permissions WHERE code = 'business:chat:list') AS tmp) WHERE code IN ('business:chat:view', 'business:chat:message:view', 'business:chat:search');

-- 设置财务管理子菜单的parent_id
UPDATE permissions SET parent_id = (SELECT id FROM (SELECT id FROM permissions WHERE code = 'finance') AS tmp) WHERE code IN ('finance:withdraw:list', 'finance:statistics:view', 'finance:reports:view');

-- 设置提现审核子按钮的parent_id
UPDATE permissions SET parent_id = (SELECT id FROM (SELECT id FROM permissions WHERE code = 'finance:withdraw:list') AS tmp) WHERE code IN ('finance:withdraw:approve', 'finance:withdraw:reject');

-- 设置低代码平台子菜单的parent_id
UPDATE permissions SET parent_id = (SELECT id FROM (SELECT id FROM permissions WHERE code = 'lowcode') AS tmp) WHERE code IN ('lowcode:form:design', 'lowcode:page:design', 'lowcode:data:manage');

-- ============================================
-- 为超级管理员分配所有权限
-- ============================================
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions WHERE status = 1
ON DUPLICATE KEY UPDATE role_id = role_id;

-- ============================================
-- 为管理员分配大部分权限
-- ============================================
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE status = 1 AND code NOT LIKE 'super:%'
ON DUPLICATE KEY UPDATE role_id = role_id;

-- ============================================
-- 为普通用户分配基本权限
-- ============================================
INSERT INTO role_permissions (role_id, permission_id)
SELECT 4, id FROM permissions WHERE status = 1 AND code IN ('system:user:list', 'business:order:list', 'business:chat:list')
ON DUPLICATE KEY UPDATE role_id = role_id;

-- ============================================
-- 为管理员用户分配超级管理员角色
-- ============================================
-- 注意：这里需要实际的用户ID，暂时跳过
-- INSERT INTO user_roles (user_id, role_id) VALUES (1, 1);

-- ============================================
-- 完成
-- ============================================
