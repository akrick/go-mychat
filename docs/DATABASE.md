# MyChat 数据库文档

## 数据库信息

- 数据库名: `mychat`
- 字符集: `utf8mb4`
- 排序规则: `utf8mb4_unicode_ci`

## 数据表结构

### 1. users (用户表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| username | varchar(50) | 用户名 | UNIQUE |
| password | varchar(255) | 密码（加密） | |
| email | varchar(100) | 邮箱 | INDEX |
| phone | varchar(20) | 手机号 | INDEX |
| avatar | varchar(255) | 头像 | |
| nickname | varchar(50) | 昵称 | |
| status | tinyint | 状态：0-禁用，1-正常 | INDEX |
| is_admin | tinyint | 是否管理员 | |
| created_at | datetime | 创建时间 | INDEX |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE users (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  username varchar(50) NOT NULL,
  password varchar(255) NOT NULL,
  email varchar(100) DEFAULT NULL,
  phone varchar(20) DEFAULT NULL,
  avatar varchar(255) DEFAULT NULL,
  nickname varchar(50) DEFAULT NULL,
  status tinyint DEFAULT 1,
  is_admin tinyint DEFAULT 0,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY username (username),
  KEY idx_email (email),
  KEY idx_phone (phone),
  KEY idx_status (status),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2. counselors (咨询师表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| user_id | bigint unsigned | 关联用户ID | INDEX |
| name | varchar(50) | 姓名 | |
| avatar | varchar(255) | 头像 | |
| specialties | varchar(255) | 专业领域 | |
| description | text | 个人简介 | |
| price_per_minute | decimal(10,2) | 每分钟价格 | |
| status | tinyint | 状态：0-禁用，1-正常 | INDEX |
| rating | decimal(3,2) | 评分 | |
| review_count | int | 评价数量 | |
| created_at | datetime | 创建时间 | |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE counselors (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id bigint unsigned NOT NULL,
  name varchar(50) NOT NULL,
  avatar varchar(255) DEFAULT NULL,
  specialties varchar(255) DEFAULT NULL,
  description text,
  price_per_minute decimal(10,2) DEFAULT 0.00,
  status tinyint DEFAULT 1,
  rating decimal(3,2) DEFAULT 5.00,
  review_count int DEFAULT 0,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_user_id (user_id),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 3. counselor_accounts (咨询师账户表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| counselor_id | bigint unsigned | 咨询师ID | UNIQUE |
| balance | decimal(10,2) | 账户余额 | |
| total_income | decimal(10,2) | 累计收入 | |
| total_withdraw | decimal(10,2) | 累计提现 | |
| frozen_amount | decimal(10,2) | 冻结金额 | |
| created_at | datetime | 创建时间 | |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE counselor_accounts (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  counselor_id bigint unsigned NOT NULL,
  balance decimal(10,2) DEFAULT 0.00,
  total_income decimal(10,2) DEFAULT 0.00,
  total_withdraw decimal(10,2) DEFAULT 0.00,
  frozen_amount decimal(10,2) DEFAULT 0.00,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY counselor_id (counselor_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 4. orders (订单表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| order_no | varchar(32) | 订单号 | UNIQUE |
| user_id | bigint unsigned | 用户ID | INDEX |
| counselor_id | bigint unsigned | 咨询师ID | INDEX |
| type | varchar(20) | 类型：text/voice/video | |
| amount | decimal(10,2) | 订单金额 | |
| status | tinyint | 状态 | INDEX |
| created_at | datetime | 创建时间 | INDEX |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE orders (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  order_no varchar(32) NOT NULL,
  user_id bigint unsigned NOT NULL,
  counselor_id bigint unsigned NOT NULL,
  type varchar(20) NOT NULL,
  amount decimal(10,2) NOT NULL,
  status tinyint DEFAULT 0,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY order_no (order_no),
  KEY idx_user_id (user_id),
  KEY idx_counselor_id (counselor_id),
  KEY idx_status (status),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 5. payments (支付表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| payment_no | varchar(32) | 支付流水号 | UNIQUE |
| order_id | bigint unsigned | 订单ID | INDEX |
| user_id | bigint unsigned | 用户ID | |
| amount | decimal(10,2) | 支付金额 | |
| payment_method | varchar(20) | 支付方式 | |
| status | tinyint | 状态 | INDEX |
| transaction_id | varchar(64) | 第三方交易ID | |
| paid_at | datetime | 支付时间 | |
| created_at | datetime | 创建时间 | |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE payments (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  payment_no varchar(32) NOT NULL,
  order_id bigint unsigned NOT NULL,
  user_id bigint unsigned NOT NULL,
  amount decimal(10,2) NOT NULL,
  payment_method varchar(20) NOT NULL,
  status tinyint DEFAULT 0,
  transaction_id varchar(64) DEFAULT NULL,
  paid_at datetime DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY payment_no (payment_no),
  KEY idx_order_id (order_id),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 6. withdraw_records (提现记录表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| withdraw_no | varchar(32) | 提现单号 | UNIQUE |
| counselor_id | bigint unsigned | 咨询师ID | INDEX |
| amount | decimal(10,2) | 提现金额 | |
| status | tinyint | 状态 | INDEX |
| account_name | varchar(50) | 账户名称 | |
| account_number | varchar(50) | 账户号码 | |
| bank_name | varchar(50) | 银行名称 | |
| remark | varchar(255) | 备注 | |
| approved_at | datetime | 审核时间 | |
| transferred_at | datetime | 打款时间 | |
| created_at | datetime | 创建时间 | INDEX |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE withdraw_records (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  withdraw_no varchar(32) NOT NULL,
  counselor_id bigint unsigned NOT NULL,
  amount decimal(10,2) NOT NULL,
  status tinyint DEFAULT 0,
  account_name varchar(50) DEFAULT NULL,
  account_number varchar(50) DEFAULT NULL,
  bank_name varchar(50) DEFAULT NULL,
  remark varchar(255) DEFAULT NULL,
  approved_at datetime DEFAULT NULL,
  transferred_at datetime DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY withdraw_no (withdraw_no),
  KEY idx_counselor_id (counselor_id),
  KEY idx_status (status),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 7. chat_sessions (聊天会话表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| session_no | varchar(32) | 会话编号 | UNIQUE |
| user_id | bigint unsigned | 用户ID | INDEX |
| counselor_id | bigint unsigned | 咨询师ID | INDEX |
| order_id | bigint unsigned | 订单ID | |
| status | tinyint | 状态 | INDEX |
| started_at | datetime | 开始时间 | |
| ended_at | datetime | 结束时间 | |
| created_at | datetime | 创建时间 | |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE chat_sessions (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  session_no varchar(32) NOT NULL,
  user_id bigint unsigned NOT NULL,
  counselor_id bigint unsigned NOT NULL,
  order_id bigint unsigned DEFAULT NULL,
  status tinyint DEFAULT 0,
  started_at datetime DEFAULT NULL,
  ended_at datetime DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY session_no (session_no),
  KEY idx_user_id (user_id),
  KEY idx_counselor_id (counselor_id),
  KEY idx_status (status),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 8. chat_messages (聊天消息表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| session_id | bigint unsigned | 会话ID | INDEX |
| sender_id | bigint unsigned | 发送者ID | INDEX |
| sender_type | varchar(20) | 发送者类型 | |
| content | text | 消息内容 | |
| message_type | varchar(20) | 消息类型 | |
| created_at | datetime | 创建时间 | INDEX |

```sql
CREATE TABLE chat_messages (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  session_id bigint unsigned NOT NULL,
  sender_id bigint unsigned NOT NULL,
  sender_type varchar(20) NOT NULL,
  content text,
  message_type varchar(20) DEFAULT 'text',
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_session_id (session_id),
  KEY idx_sender_id (sender_id),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 9. chat_billing (聊天计费表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| session_id | bigint unsigned | 会话ID | INDEX |
| billing_no | varchar(32) | 计费单号 | UNIQUE |
| duration | int | 时长（秒） | |
| amount | decimal(10,2) | 金额 | |
| billed_at | datetime | 计费时间 | |
| created_at | datetime | 创建时间 | |

```sql
CREATE TABLE chat_billing (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  session_id bigint unsigned NOT NULL,
  billing_no varchar(32) NOT NULL,
  duration int NOT NULL,
  amount decimal(10,2) NOT NULL,
  billed_at datetime DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY billing_no (billing_no),
  KEY idx_session_id (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 10. reviews (评价表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| user_id | bigint unsigned | 用户ID | INDEX |
| counselor_id | bigint unsigned | 咨询师ID | INDEX |
| session_id | bigint unsigned | 会话ID | |
| rating | tinyint | 评分（1-5） | |
| content | text | 评价内容 | |
| reply | text | 回复内容 | |
| created_at | datetime | 创建时间 | INDEX |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE reviews (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id bigint unsigned NOT NULL,
  counselor_id bigint unsigned NOT NULL,
  session_id bigint unsigned DEFAULT NULL,
  rating tinyint NOT NULL,
  content text,
  reply text,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_user_id (user_id),
  KEY idx_counselor_id (counselor_id),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 11. system_logs (系统日志表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| user_id | bigint unsigned | 用户ID | INDEX |
| level | varchar(20) | 日志级别 | INDEX |
| module | varchar(50) | 模块名 | |
| action | varchar(100) | 操作 | |
| message | text | 日志消息 | |
| ip_address | varchar(50) | IP地址 | |
| user_agent | varchar(500) | 用户代理 | |
| created_at | datetime | 创建时间 | INDEX |

```sql
CREATE TABLE system_logs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id bigint unsigned DEFAULT NULL,
  level varchar(20) NOT NULL,
  module varchar(50) DEFAULT NULL,
  action varchar(100) DEFAULT NULL,
  message text,
  ip_address varchar(50) DEFAULT NULL,
  user_agent varchar(500) DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_user_id (user_id),
  KEY idx_level (level),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 12. system_configs (系统配置表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| key | varchar(100) | 配置键 | UNIQUE |
| value | text | 配置值 | |
| category | varchar(50) | 配置分类 | INDEX |
| label | varchar(100) | 配置标签 | |
| type | varchar(20) | 配置类型 | |
| is_system | tinyint | 是否系统配置 | |
| sort | int | 排序 | |
| created_at | datetime | 创建时间 | |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE system_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) NOT NULL,
  value text,
  category varchar(50) NOT NULL,
  label varchar(100) NOT NULL,
  type varchar(20) NOT NULL,
  is_system tinyint DEFAULT 0,
  sort int DEFAULT 0,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY `key` (`key`),
  KEY idx_category (category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 13. roles (角色表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| name | varchar(50) | 角色名称 | |
| description | varchar(255) | 角色描述 | |
| status | tinyint | 状态 | INDEX |
| created_at | datetime | 创建时间 | |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE roles (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(50) NOT NULL,
  description varchar(255) DEFAULT NULL,
  status tinyint DEFAULT 1,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 14. permissions (权限表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| name | varchar(100) | 权限名称 | |
| code | varchar(100) | 权限代码 | UNIQUE |
| type | varchar(20) | 权限类型 | |
| parent_id | bigint unsigned | 父权限ID | INDEX |
| sort | int | 排序 | |
| created_at | datetime | 创建时间 | |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE permissions (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(100) NOT NULL,
  code varchar(100) NOT NULL,
  type varchar(20) NOT NULL,
  parent_id bigint unsigned DEFAULT 0,
  sort int DEFAULT 0,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY code (code),
  KEY idx_parent_id (parent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 15. role_permissions (角色权限关联表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| role_id | bigint unsigned | 角色ID | INDEX |
| permission_id | bigint unsigned | 权限ID | INDEX |
| created_at | datetime | 创建时间 | |

```sql
CREATE TABLE role_permissions (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  role_id bigint unsigned NOT NULL,
  permission_id bigint unsigned NOT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_role_id (role_id),
  KEY idx_permission_id (permission_id),
  UNIQUE KEY role_permission (role_id, permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 16. user_roles (用户角色关联表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| user_id | bigint unsigned | 用户ID | INDEX |
| role_id | bigint unsigned | 角色ID | INDEX |
| created_at | datetime | 创建时间 | |

```sql
CREATE TABLE user_roles (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id bigint unsigned NOT NULL,
  role_id bigint unsigned NOT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_user_id (user_id),
  KEY idx_role_id (role_id),
  UNIQUE KEY user_role (user_id, role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 17. menus (菜单表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| name | varchar(50) | 菜单名称 | |
| path | varchar(200) | 路由路径 | |
| icon | varchar(50) | 图标 | |
| parent_id | bigint unsigned | 父菜单ID | INDEX |
| sort | int | 排序 | |
| status | tinyint | 状态 | INDEX |
| created_at | datetime | 创建时间 | |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE menus (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(50) NOT NULL,
  path varchar(200) DEFAULT NULL,
  icon varchar(50) DEFAULT NULL,
  parent_id bigint unsigned DEFAULT 0,
  sort int DEFAULT 0,
  status tinyint DEFAULT 1,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_parent_id (parent_id),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 18. lowcode_forms (低代码表单表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| name | varchar(100) | 表单名称 | |
| description | varchar(255) | 表单描述 | |
| schema | json | 表单结构 | |
| status | tinyint | 状态 | INDEX |
| created_at | datetime | 创建时间 | |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE lowcode_forms (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(100) NOT NULL,
  description varchar(255) DEFAULT NULL,
  schema json,
  status tinyint DEFAULT 1,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 19. lowcode_pages (低代码页面表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| name | varchar(100) | 页面名称 | |
| description | varchar(255) | 页面描述 | |
| schema | json | 页面结构 | |
| status | tinyint | 状态 | INDEX |
| created_at | datetime | 创建时间 | |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE lowcode_pages (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(100) NOT NULL,
  description varchar(255) DEFAULT NULL,
  schema json,
  status tinyint DEFAULT 1,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 20. lowcode_form_data (表单数据表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| form_id | bigint unsigned | 表单ID | INDEX |
| user_id | bigint unsigned | 提交用户ID | |
| data | json | 表单数据 | |
| created_at | datetime | 创建时间 | INDEX |

```sql
CREATE TABLE lowcode_form_data (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  form_id bigint unsigned NOT NULL,
  user_id bigint unsigned DEFAULT NULL,
  data json,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_form_id (form_id),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 21. notifications (通知表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| user_id | bigint unsigned | 用户ID | INDEX |
| type | varchar(20) | 通知类型 | INDEX |
| level | varchar(20) | 通知级别 | |
| title | varchar(255) | 通知标题 | |
| content | text | 通知内容 | |
| extra_data | text | 额外数据(JSON) | |
| is_read | tinyint | 是否已读 | INDEX |
| read_time | datetime | 阅读时间 | |
| created_at | datetime | 创建时间 | INDEX |

```sql
CREATE TABLE notifications (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id bigint unsigned NOT NULL,
  type varchar(20) NOT NULL,
  level varchar(20) DEFAULT 'info',
  title varchar(255) NOT NULL,
  content text,
  extra_data text,
  is_read tinyint DEFAULT 0,
  read_time datetime DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_user_id (user_id),
  KEY idx_type (type),
  KEY idx_is_read (is_read),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 22. files (文件表)

| 字段名 | 类型 | 说明 | 索引 |
|--------|------|------|------|
| id | bigint unsigned | 主键ID | PRIMARY |
| file_name | varchar(255) | 文件名 | |
| original_name | varchar(255) | 原始文件名 | |
| file_path | varchar(500) | 文件路径 | |
| file_size | bigint | 文件大小(字节) | |
| file_type | varchar(100) | 文件类型 | |
| mime_type | varchar(100) | MIME类型 | |
| md5 | varchar(32) | 文件MD5 | UNIQUE |
| storage_type | varchar(20) | 存储类型 | |
| bucket_name | varchar(100) | OSS桶名 | |
| uploaded_by | bigint unsigned | 上传人ID | INDEX |
| is_deleted | tinyint | 是否删除 | INDEX |
| created_at | datetime | 创建时间 | |
| updated_at | datetime | 更新时间 | |

```sql
CREATE TABLE files (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  file_name varchar(255) NOT NULL,
  original_name varchar(255) NOT NULL,
  file_path varchar(500) NOT NULL,
  file_size bigint NOT NULL,
  file_type varchar(100) NOT NULL,
  mime_type varchar(100) NOT NULL,
  md5 varchar(32) NOT NULL,
  storage_type varchar(20) DEFAULT 'local',
  bucket_name varchar(100) DEFAULT NULL,
  uploaded_by bigint unsigned NOT NULL,
  is_deleted tinyint DEFAULT 0,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY md5 (md5),
  KEY idx_uploaded_by (uploaded_by),
  KEY idx_is_deleted (is_deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 索引说明

### 主键索引
所有表都使用自增主键 `id` 作为主键。

### 唯一索引
- `users.username`: 用户名唯一
- `orders.order_no`: 订单号唯一
- `payments.payment_no`: 支付流水号唯一
- `withdraw_records.withdraw_no`: 提现单号唯一
- `chat_sessions.session_no`: 会话编号唯一
- `system_configs.key`: 配置键唯一
- `permissions.code`: 权限代码唯一
- `files.md5`: 文件MD5唯一

### 复合索引
- `role_permissions`: (role_id, permission_id)
- `user_roles`: (user_id, role_id)

### 普通索引
主要用于查询优化，包括：
- 用户相关: `email`, `phone`, `status`
- 订单相关: `user_id`, `counselor_id`, `status`, `created_at`
- 日志相关: `user_id`, `level`, `created_at`
- 配置相关: `category`
- 权限相关: `parent_id`, `status`

## 数据初始化

### 初始化脚本
```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE mychat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 导入数据结构
mysql -u root -p mychat < api/init_data.sql

# 导入RBAC数据
mysql -u root -p mychat < api/init_rbac.sql
```

### 默认数据
- 管理员账户: username=`admin`, password=`admin123`（需要加密）
- 36个系统配置
- 默认权限和角色

## 数据库优化建议

### 1. 定期清理
- 清理过期的聊天消息（保留90天）
- 清理已删除的文件
- 清理过期的系统日志

### 2. 分区表
对于数据量大的表（如 chat_messages），可以考虑按时间分区。

### 3. 读写分离
在高并发场景下，可以考虑使用主从复制，实现读写分离。

### 4. 索引优化
- 定期分析慢查询日志
- 根据实际查询情况调整索引
- 删除冗余索引

### 5. 数据备份
```bash
# 备份数据库
mysqldump -u root -p mychat > backup_$(date +%Y%m%d).sql

# 恢复数据库
mysql -u root -p mychat < backup_20260128.sql
```

## 字段类型说明

### 整型
- `tinyint`: 0-255，用于状态、标志位
- `int`: -2^31 到 2^31-1，用于数量、排序
- `bigint unsigned`: 0 到 2^64-1，用于主键ID

### 字符串
- `varchar(n)`: 可变长度字符串，n为最大长度
- `text`: 长文本，用于内容、描述

### 小数
- `decimal(10,2)`: 10位数字，其中2位小数

### 日期时间
- `datetime`: 日期时间，格式 `YYYY-MM-DD HH:MM:SS`

### JSON
- `json`: JSON格式数据，用于配置、表单结构等
