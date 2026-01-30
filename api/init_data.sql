-- MyChat 数据库初始化脚本
-- 创建时间: 2026-01-26

-- ============================================
-- 1. 用户相关表
-- ============================================

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    avatar VARCHAR(255),
    status INT DEFAULT 1 COMMENT '1-正常,0-禁用',
    is_admin BOOLEAN DEFAULT FALSE COMMENT '是否管理员',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 咨询师表
CREATE TABLE IF NOT EXISTS counselors (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    title VARCHAR(50) COMMENT '职称',
    avatar VARCHAR(255) COMMENT '头像',
    bio TEXT COMMENT '个人简介',
    specialty VARCHAR(255) COMMENT '擅长领域',
    price DECIMAL(10,2) NOT NULL COMMENT '单价(元/分钟)',
    years_exp INT COMMENT '从业年限',
    rating DECIMAL(3,2) DEFAULT 5.00 COMMENT '评分',
    status INT NOT NULL DEFAULT 1 COMMENT '状态:1-启用,0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='咨询师表';

-- 咨询师账户表
CREATE TABLE IF NOT EXISTS counselor_accounts (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    counselor_id INT UNSIGNED NOT NULL UNIQUE,
    total_income DECIMAL(10,2) DEFAULT 0.00 COMMENT '总收入',
    withdrawn DECIMAL(10,2) DEFAULT 0.00 COMMENT '已提现',
    balance DECIMAL(10,2) DEFAULT 0.00 COMMENT '可用余额',
    frozen_amount DECIMAL(10,2) DEFAULT 0.00 COMMENT '冻结金额',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_counselor_id (counselor_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='咨询师账户表';

-- 咨询师统计表
CREATE TABLE IF NOT EXISTS counselor_statistics (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    counselor_id INT UNSIGNED NOT NULL UNIQUE,
    total_orders INT NOT NULL DEFAULT 0 COMMENT '总订单数',
    completed_orders INT NOT NULL DEFAULT 0 COMMENT '已完成订单数',
    cancelled_orders INT NOT NULL DEFAULT 0 COMMENT '已取消订单数',
    total_duration INT NOT NULL DEFAULT 0 COMMENT '总咨询时长(分钟)',
    total_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '总金额',
    review_count INT NOT NULL DEFAULT 0 COMMENT '评价数量',
    avg_rating DECIMAL(3,2) NOT NULL DEFAULT 0.00 COMMENT '平均评分',
    sum_rating INT NOT NULL DEFAULT 0 COMMENT '总评分',
    last_order_time TIMESTAMP NULL COMMENT '最后订单时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_counselor_id (counselor_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='咨询师统计表';

-- ============================================
-- 2. 订单相关表
-- ============================================

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(32) NOT NULL UNIQUE,
    user_id INT UNSIGNED NOT NULL,
    counselor_id INT UNSIGNED NOT NULL,
    duration INT NOT NULL COMMENT '咨询时长(分钟)',
    amount DECIMAL(10,2) NOT NULL,
    status INT NOT NULL DEFAULT 0 COMMENT '0-待支付,1-已支付,2-已完成,3-已取消,4-已退款',
    schedule_time TIMESTAMP NOT NULL COMMENT '预约时间',
    notes TEXT COMMENT '备注',
    pay_time TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_order_no (order_no),
    INDEX idx_user_id (user_id),
    INDEX idx_counselor_id (counselor_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 支付记录表
CREATE TABLE IF NOT EXISTS payments (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    payment_no VARCHAR(32) NOT NULL UNIQUE,
    order_id INT UNSIGNED NOT NULL,
    order_no VARCHAR(32) NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    payment_method VARCHAR(20) NOT NULL COMMENT '支付方式:wechat/alipay',
    trade_type VARCHAR(20) COMMENT '交易类型',
    transaction_id VARCHAR(64) UNIQUE COMMENT '第三方支付交易号',
    amount DECIMAL(10,2) NOT NULL,
    status INT NOT NULL DEFAULT 0 COMMENT '支付状态:0-待支付,1-已支付,2-失败,3-已退款,4-已取消',
    pay_time TIMESTAMP NULL,
    notify_time TIMESTAMP NULL,
    notify_data TEXT COMMENT '支付回调原始数据',
    failure_reason VARCHAR(255) COMMENT '失败原因',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_payment_no (payment_no),
    INDEX idx_order_id (order_id),
    INDEX idx_order_no (order_no),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_transaction_id (transaction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付记录表';

-- 支付配置表
CREATE TABLE IF NOT EXISTS payment_configs (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    payment_method VARCHAR(20) NOT NULL UNIQUE COMMENT '支付方式',
    app_id VARCHAR(64) COMMENT '应用ID',
    mch_id VARCHAR(64) COMMENT '商户号',
    api_secret VARCHAR(128) COMMENT 'API密钥',
    api_cert_path VARCHAR(255) COMMENT '证书路径',
    api_key_path VARCHAR(255) COMMENT '密钥路径',
    notify_url VARCHAR(255) COMMENT '回调地址',
    private_key_path VARCHAR(255) COMMENT '私钥路径',
    public_key_path VARCHAR(255) COMMENT '公钥路径',
    is_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    is_sandbox BOOLEAN DEFAULT FALSE COMMENT '是否沙箱环境',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_payment_method (payment_method)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付配置表';

-- ============================================
-- 3. 聊天相关表
-- ============================================

-- 聊天会话表
CREATE TABLE IF NOT EXISTS chat_sessions (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id INT UNSIGNED NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    counselor_id INT UNSIGNED NOT NULL,
    status INT NOT NULL DEFAULT 0 COMMENT '状态:0-待开始,1-进行中,2-已结束,3-已超时',
    start_time TIMESTAMP NULL,
    end_time TIMESTAMP NULL,
    duration INT COMMENT '实际时长(秒)',
    price DECIMAL(10,2) COMMENT '单价(元/分钟)',
    total_amount DECIMAL(10,2) COMMENT '总金额(元)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_order_id (order_id),
    INDEX idx_user_id (user_id),
    INDEX idx_counselor_id (counselor_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天会话表';

-- 聊天消息表
CREATE TABLE IF NOT EXISTS chat_messages (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    session_id INT UNSIGNED NOT NULL,
    sender_id INT UNSIGNED NOT NULL,
    sender_type VARCHAR(20) NOT NULL COMMENT '发送者类型:user/counselor',
    content_type VARCHAR(20) DEFAULT 'text' COMMENT '内容类型:text/image/file',
    content TEXT COMMENT '消息内容',
    file_url VARCHAR(255) COMMENT '文件URL',
    is_read BOOLEAN DEFAULT FALSE COMMENT '是否已读',
    read_time TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_session_id (session_id),
    INDEX idx_sender_id (sender_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天消息表';

-- 聊天计费记录表
CREATE TABLE IF NOT EXISTS chat_billings (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    session_id INT UNSIGNED NOT NULL UNIQUE,
    order_id INT UNSIGNED NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    counselor_id INT UNSIGNED NOT NULL,
    duration INT NOT NULL COMMENT '计费时长(秒)',
    price_per_minute DECIMAL(10,2) NOT NULL COMMENT '单价(元/分钟)',
    total_amount DECIMAL(10,2) NOT NULL COMMENT '总金额',
    platform_fee DECIMAL(10,2) NOT NULL COMMENT '平台费用(30%)',
    counselor_fee DECIMAL(10,2) NOT NULL COMMENT '咨询师收入(70%)',
    status INT NOT NULL DEFAULT 0 COMMENT '状态:0-待结算,1-已结算',
    settled_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_session_id (session_id),
    INDEX idx_order_id (order_id),
    INDEX idx_user_id (user_id),
    INDEX idx_counselor_id (counselor_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天计费记录表';

-- ============================================
-- 4. 评价相关表
-- ============================================

-- 评价表
CREATE TABLE IF NOT EXISTS reviews (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id INT UNSIGNED NOT NULL UNIQUE,
    order_no VARCHAR(32) NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    counselor_id INT UNSIGNED NOT NULL,
    rating INT NOT NULL COMMENT '评分(1-5)',
    service_rating INT NOT NULL DEFAULT 0 COMMENT '服务评分',
    professionalism INT NOT NULL DEFAULT 0 COMMENT '专业度评分',
    effectiveness INT NOT NULL DEFAULT 0 COMMENT '有效性评分',
    content TEXT COMMENT '评价内容',
    is_anonymous BOOLEAN DEFAULT FALSE COMMENT '是否匿名',
    status INT NOT NULL DEFAULT 1 COMMENT '状态:1-显示,0-隐藏',
    reply_content TEXT COMMENT '咨询师回复',
    reply_time TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_order_id (order_id),
    INDEX idx_order_no (order_no),
    INDEX idx_user_id (user_id),
    INDEX idx_counselor_id (counselor_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评价表';

-- ============================================
-- 5. 财务相关表
-- ============================================

-- 提现记录表
CREATE TABLE IF NOT EXISTS withdraw_records (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    counselor_id INT UNSIGNED NOT NULL,
    amount DECIMAL(10,2) NOT NULL COMMENT '提现金额',
    status INT NOT NULL DEFAULT 0 COMMENT '状态:0-待审核,1-已通过,2-已拒绝,3-已打款',
    bank_name VARCHAR(50) COMMENT '开户行',
    bank_account VARCHAR(50) COMMENT '银行账号',
    account_name VARCHAR(50) COMMENT '账户名',
    rejected_reason VARCHAR(255) COMMENT '拒绝原因',
    audited_at TIMESTAMP NULL,
    transferred_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_counselor_id (counselor_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='提现记录表';

-- ============================================
-- 6. RBAC权限相关表
-- ============================================

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255) COMMENT '描述',
    sort INT DEFAULT 0 COMMENT '排序',
    status INT DEFAULT 1 COMMENT '状态:0-禁用,1-启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_name (name),
    UNIQUE INDEX idx_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    parent_id INT UNSIGNED DEFAULT 0,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    type VARCHAR(20) DEFAULT 'menu' COMMENT '类型:menu-菜单,button-按钮,api-接口',
    path VARCHAR(255) COMMENT '路由路径',
    icon VARCHAR(50) COMMENT '图标',
    component VARCHAR(255) COMMENT '组件路径',
    sort INT DEFAULT 0 COMMENT '排序',
    status INT DEFAULT 1 COMMENT '状态:0-禁用,1-启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id),
    UNIQUE INDEX idx_code (code),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL,
    role_id INT UNSIGNED NOT NULL,
    UNIQUE INDEX idx_user_role (user_id, role_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    role_id INT UNSIGNED NOT NULL,
    permission_id INT UNSIGNED NOT NULL,
    UNIQUE INDEX idx_role_permission (role_id, permission_id),
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- ============================================
-- 7. 系统管理表
-- ============================================

-- 系统日志表
CREATE TABLE IF NOT EXISTS sys_logs (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL,
    username VARCHAR(50) NOT NULL COMMENT '操作用户名',
    module VARCHAR(50) NOT NULL COMMENT '操作模块',
    action VARCHAR(100) NOT NULL COMMENT '操作动作',
    method VARCHAR(10) NOT NULL COMMENT '请求方法',
    ip VARCHAR(50) COMMENT 'IP地址',
    url VARCHAR(255) COMMENT '请求URL',
    params TEXT COMMENT '请求参数',
    result TEXT COMMENT '返回结果',
    status INT DEFAULT 1 COMMENT '状态:1-成功,0-失败',
    error_msg TEXT COMMENT '错误信息',
    duration INT COMMENT '执行时长(毫秒)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_module (module),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统日志表';

-- 在线用户表
CREATE TABLE IF NOT EXISTS online_users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL UNIQUE,
    token VARCHAR(255) NOT NULL COMMENT 'Token',
    ip VARCHAR(50) COMMENT 'IP地址',
    user_agent VARCHAR(500) COMMENT '浏览器UA',
    login_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活动时间',
    UNIQUE INDEX idx_user_id (user_id),
    INDEX idx_updated_at (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='在线用户表';

-- 系统配置表
CREATE TABLE IF NOT EXISTS sys_configs (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    config_key VARCHAR(100) NOT NULL UNIQUE,
    config_name VARCHAR(100) NOT NULL COMMENT '配置名称',
    config_type VARCHAR(20) DEFAULT 'string' COMMENT '配置类型:string/number/boolean/json',
    config_val TEXT NOT NULL COMMENT '配置值',
    is_public BOOLEAN DEFAULT FALSE COMMENT '是否公开',
    remark VARCHAR(500) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_config_key (config_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- ============================================
-- 8. 低代码平台表
-- ============================================

-- 表单设计表
CREATE TABLE IF NOT EXISTS form_designs (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(500) COMMENT '表单描述',
    form_schema LONGTEXT NOT NULL COMMENT '表单配置JSON',
    is_published BOOLEAN DEFAULT FALSE COMMENT '是否发布',
    created_by INT UNSIGNED COMMENT '创建人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_code (code),
    INDEX idx_created_by (created_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='表单设计表';

-- 表单数据表
CREATE TABLE IF NOT EXISTS form_data (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    form_id INT UNSIGNED NOT NULL,
    submit_by INT UNSIGNED NOT NULL COMMENT '提交人ID',
    data LONGTEXT NOT NULL COMMENT '表单数据JSON',
    ip VARCHAR(50) COMMENT '提交IP',
    user_agent VARCHAR(500) COMMENT '浏览器UA',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_form_id (form_id),
    INDEX idx_submit_by (submit_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='表单数据表';

-- 页面设计表
CREATE TABLE IF NOT EXISTS page_designs (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    path VARCHAR(255) NOT NULL UNIQUE COMMENT '页面路径',
    description VARCHAR(500) COMMENT '页面描述',
    page_config LONGTEXT NOT NULL COMMENT '页面配置JSON',
    is_published BOOLEAN DEFAULT FALSE COMMENT '是否发布',
    created_by INT UNSIGNED COMMENT '创建人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_code (code),
    UNIQUE INDEX idx_path (path),
    INDEX idx_created_by (created_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='页面设计表';

-- ============================================
-- 9. 文件和通知表
-- ============================================

-- 文件表
CREATE TABLE IF NOT EXISTS files (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL COMMENT '文件名',
    original_name VARCHAR(255) NOT NULL COMMENT '原始文件名',
    file_path VARCHAR(500) NOT NULL COMMENT '文件路径',
    file_size BIGINT NOT NULL COMMENT '文件大小(字节)',
    file_type VARCHAR(100) NOT NULL COMMENT '文件类型',
    mime_type VARCHAR(100) NOT NULL COMMENT 'MIME类型',
    md5 VARCHAR(32) UNIQUE COMMENT '文件MD5',
    storage_type VARCHAR(20) DEFAULT 'local' COMMENT '存储类型:local/oss/qiniu',
    bucket_name VARCHAR(100) COMMENT 'OSS桶名',
    uploaded_by INT UNSIGNED NOT NULL COMMENT '上传人ID',
    is_deleted BOOLEAN DEFAULT FALSE COMMENT '是否删除',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_md5 (md5),
    INDEX idx_uploaded_by (uploaded_by),
    INDEX idx_is_deleted (is_deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件表';

-- 通知表
CREATE TABLE IF NOT EXISTS notifications (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL,
    title VARCHAR(100) NOT NULL COMMENT '通知标题',
    content TEXT NOT NULL COMMENT '通知内容',
    type VARCHAR(20) NOT NULL COMMENT '通知类型:order/chat/system',
    is_read BOOLEAN DEFAULT FALSE COMMENT '是否已读',
    read_time TIMESTAMP NULL,
    link_url VARCHAR(255) COMMENT '跳转链接',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_is_read (is_read),
    INDEX idx_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知表';

-- ============================================
-- 初始化数据
-- ============================================

-- 插入默认管理员账号 (密码: admin123)
INSERT INTO users (username, password, email, is_admin, status) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdJ17fqwlK8s', 'admin@mychat.com', TRUE, 1)
ON DUPLICATE KEY UPDATE username = VALUES(username);

-- 插入默认角色
INSERT INTO roles (name, code, description, status) VALUES
('超级管理员', 'super_admin', '拥有所有权限', 1),
('管理员', 'admin', '拥有大部分权限', 1),
('运营人员', 'operator', '运营管理权限', 1),
('咨询师', 'counselor', '咨询师权限', 1),
('用户', 'user', '普通用户权限', 1)
ON DUPLICATE KEY UPDATE code = VALUES(code);

-- 插入默认系统配置
INSERT INTO sys_configs (config_key, config_name, config_type, config_val, is_public, remark) VALUES
('site_name', '网站名称', 'string', 'MyChat', TRUE, '网站显示名称'),
('site_logo', '网站Logo', 'string', '', TRUE, '网站Logo URL'),
('platform_rate', '平台抽成比例', 'number', '0.3', FALSE, '平台抽成比例(0-1)'),
('counselor_rate', '咨询师分成比例', 'number', '0.7', FALSE, '咨询师分成比例(0-1)'),
('min_withdraw', '最低提现金额', 'number', '100', FALSE, '最低提现金额(元)'),
('max_withdraw', '最高提现金额', 'number', '10000', FALSE, '最高提现金额(元)'),
('withdraw_fee', '提现手续费', 'number', '0', FALSE, '提现手续费比例'),
('system_notice', '系统公告', 'string', '', TRUE, '系统公告内容')
ON DUPLICATE KEY UPDATE config_key = VALUES(config_key);

-- 插入测试咨询师数据
INSERT INTO counselors (name, title, avatar, bio, specialty, price, years_exp, rating, status) VALUES
('张明', '国家二级心理咨询师', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=300&h=300&fit=crop&crop=face', '拥有10年心理咨询经验，擅长情绪管理、人际关系、婚姻家庭咨询，已帮助超过2000名来访者走出困境。', '情绪管理,人际关系,婚姻家庭', 2.50, 10, 4.90, 1),
('李雪', '高级心理咨询师', 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=300&h=300&fit=crop&crop=face', '专注于青少年心理咨询，擅长学习压力、青春期困惑、亲子关系等领域，具有丰富的临床经验。', '青少年心理,学习压力,亲子关系', 3.00, 8, 4.85, 1),
('王芳', '婚姻家庭咨询师', 'https://images.unsplash.com/photo-1580489944761-15a19d654956?w=300&h=300&fit=crop&crop=face', '专业从事婚姻家庭咨询15年，擅长夫妻关系改善、离婚危机干预、家庭矛盾调解，帮助众多家庭重获幸福。', '婚姻关系,家庭矛盾,情感咨询', 3.50, 15, 4.95, 1),
('刘伟', '职业规划咨询师', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=300&h=300&fit=crop&crop=face', '职业规划专家，曾任职于多家知名企业HR总监，擅长职业规划、面试辅导、职场人际关系处理。', '职业规划,面试辅导,职场咨询', 4.00, 12, 4.88, 1),
('陈静', '认知行为治疗师', 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=300&h=300&fit=crop&crop=face', '认知行为治疗(CBT)专业认证咨询师，擅长焦虑症、抑郁症、强迫症等心理疾病的康复指导。', '焦虑抑郁,强迫症,心理疾病', 5.00, 9, 4.92, 1),
('赵敏', '儿童心理发展专家', 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=300&h=300&fit=crop&crop=face', '儿童心理学博士，专注于儿童心理发展、学习障碍、多动症、自闭症等问题的早期干预和治疗。', '儿童心理,学习障碍,多动症', 4.50, 7, 4.87, 1),
('孙涛', '心理危机干预师', 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=300&h=300&fit=crop&crop=face', '心理危机干预专家，处理过大量突发性创伤事件后的心理救援工作，包括自然灾害、事故创伤等。', '危机干预,创伤后应激,创伤修复', 3.80, 11, 4.91, 1),
('周琳', '艺术治疗师', 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=300&h=300&fit=crop&crop=face', '艺术治疗专业认证咨询师，运用绘画、音乐、舞蹈等艺术形式帮助来访者表达情感、释放压力。', '艺术治疗,压力管理,情感表达', 3.20, 6, 4.84, 1)
ON DUPLICATE KEY UPDATE name = VALUES(name);
