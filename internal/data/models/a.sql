CREATE TABLE `tb_company_cfg` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `third_company_id` varchar(20) NOT NULL COMMENT '三方租户id',
    `platform_ids` varchar(100) NOT NULL COMMENT '平台id, 用来区分多种数据源,多个用逗号分隔',
    `company_id` varchar(20) NOT NULL COMMENT '云文档租户id',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态,0-禁用,1-启用',
    `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE COMMENT '主键索引',
    UNIQUE KEY `uk_third_company_id` (`third_company_id`) USING BTREE COMMENT 'third_company_id唯一索引',
    UNIQUE KEY `uk_company_id` (`company_id`) USING BTREE COMMENT 'company_id唯一索引'
) ENGINE = InnoDB AUTO_INCREMENT = 5 DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '租户关系表' -- user
CREATE TABLE `tb_las_user` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `task_id` varchar(20) NOT NULL COMMENT '任务id',
    `third_company_id` varchar(20) NOT NULL COMMENT '租户id',
    `platform_id` varchar(60) NOT NULL COMMENT '平台id, 用来区分多种数据源，platform_id + uid 唯一',
    `uid` varchar(255) NOT NULL COMMENT '用户id',
    `def_did` varchar(255) DEFAULT NULL COMMENT '默认部门',
    `def_did_order` int DEFAULT '0' COMMENT '用户在默认部门下的排序',
    `account` varchar(255) NOT NULL COMMENT '登录名，对应account',
    `nick_name` varchar(255) NOT NULL COMMENT '用户昵称，对应nick_name',
    `password` varchar(255) DEFAULT NULL COMMENT '密码',
    `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
    `email` varchar(80) DEFAULT NULL COMMENT '邮箱',
    `gender` varchar(60) DEFAULT NULL COMMENT '用户性别',
    `title` varchar(255) DEFAULT NULL COMMENT '职称',
    `work_place` varchar(255) DEFAULT NULL COMMENT '办公地点',
    `leader` varchar(255) DEFAULT NULL COMMENT '上级主管ID',
    `employer` varchar(255) DEFAULT NULL COMMENT '员工工号',
    `employment_status` varchar(60) NOT NULL DEFAULT 'notactive' COMMENT '就职状态[active, notactive, disabled]',
    `employment_type` varchar(60) DEFAULT NULL COMMENT '就职类型[permanent, intern]',
    `phone` varchar(200) DEFAULT NULL COMMENT '手机号',
    `telephone` varchar(200) DEFAULT NULL COMMENT '座机号',
    `source` varchar(20) DEFAULT 'sync' COMMENT '来源, buildin/sync',
    `custom_fields` varchar(5000) DEFAULT NULL COMMENT '自定义字段，json数组',
    `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `check_type` tinyint NOT NULL DEFAULT '0' COMMENT '1-勾选 0-未勾选',
    PRIMARY KEY (`id`) USING BTREE COMMENT '主键索引',
    UNIQUE KEY `uk_task_uid` (`uid`, `task_id`, `platform_id`) USING BTREE COMMENT 'uid唯一索引',
    UNIQUE KEY `uk_task_company_name` (`account`, `task_id`, `third_company_id`) USING BTREE COMMENT 'name唯一索引',
    KEY `idx_task_company` (`task_id`, `third_company_id`) USING BTREE COMMENT 'task索引'
) ENGINE = InnoDB AUTO_INCREMENT = 1257 DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '三方用户表' -- department
CREATE TABLE `tb_las_department` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `did` varchar(255) NOT NULL COMMENT '部门id',
    `task_id` varchar(20) NOT NULL COMMENT '任务id',
    `third_company_id` varchar(20) NOT NULL COMMENT '租户id',
    `platform_id` varchar(60) NOT NULL COMMENT '平台id, 用来区分多种数据源，platform_id + did 唯一, 根部门例外',
    `pid` varchar(255) DEFAULT NULL COMMENT '父部门id',
    `name` varchar(255) NOT NULL COMMENT '部门名称',
    `order` int DEFAULT '0' COMMENT '排序',
    `source` varchar(20) DEFAULT 'sync' COMMENT '来源, buildin/sync',
    `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `check_type` tinyint NOT NULL DEFAULT '0' COMMENT '1-勾选 0-未勾选',
    `type` varchar(255) DEFAULT NULL COMMENT '类型, 仅支持小写字母和下划线组成',
    PRIMARY KEY (`id`) USING BTREE COMMENT '主键索引',
    UNIQUE KEY `uk_task_did` (
        `did`,
        `task_id`,
        `third_company_id`,
        `platform_id`
    ) USING BTREE COMMENT 'did唯一索引',
    KEY `idx_pid` (
        `pid`,
        `task_id`,
        `third_company_id`,
        `platform_id`
    ) USING BTREE COMMENT '父部门id索引',
    KEY `idx_task_company` (`task_id`, `third_company_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 5015 DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '三方部门表' -- relation
CREATE TABLE `tb_las_department_user` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `task_id` varchar(20) NOT NULL COMMENT '任务id',
    `third_company_id` varchar(20) NOT NULL COMMENT '租户id',
    `platform_id` varchar(60) NOT NULL COMMENT '平台id, 用来区分多种数据源，platform_id + id 唯一',
    `uid` varchar(255) NOT NULL COMMENT '用户id',
    `did` varchar(255) NOT NULL COMMENT '部门id',
    `order` int DEFAULT NULL COMMENT '用户在部门下的排序',
    `main` int DEFAULT '0' COMMENT '是否是主部门，1：是 0：不是',
    `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `check_type` tinyint NOT NULL DEFAULT '0' COMMENT '1-勾选 0-未勾选',
    PRIMARY KEY (`id`) USING BTREE COMMENT '主键索引',
    UNIQUE KEY `uid` (
        `uid`,
        `did`,
        `task_id`,
        `third_company_id`,
        `platform_id`
    ) COMMENT '唯一索引',
    KEY `idx_did` (`did`, `task_id`) USING BTREE COMMENT '部门id索引',
    KEY `idx_uid` (`uid`, `task_id`) USING BTREE COMMENT 'uid索引',
    KEY `idx_task_company` (`task_id`, `third_company_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 999 DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '三方部门用户关系表'

--- incre user

CREATE TABLE `tb_las_user_increment` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `third_company_id` varchar(20) NOT NULL COMMENT '租户id',
    `platform_id` varchar(60) NOT NULL COMMENT '平台id, 用来区分多种数据源，platform_id + uid 唯一',
    `uid` varchar(255) NOT NULL COMMENT '用户id',
    `def_did` varchar(255) DEFAULT NULL COMMENT '默认部门',
    `def_did_order` int DEFAULT '0' COMMENT '用户在默认部门下的排序',
    `account` varchar(255) NOT NULL COMMENT '登录名，对应account',
    `nick_name` varchar(255) NOT NULL COMMENT '用户昵称，对应nick_name',
    `password` varchar(255) DEFAULT NULL COMMENT '密码',
    `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
    `email` varchar(80) DEFAULT NULL COMMENT '邮箱',
    `gender` varchar(60) DEFAULT NULL COMMENT '用户性别',
    `title` varchar(255) DEFAULT NULL COMMENT '职称',
    `work_place` varchar(255) DEFAULT NULL COMMENT '办公地点',
    `leader` varchar(255) DEFAULT NULL COMMENT '上级主管ID',
    `employer` varchar(255) DEFAULT NULL COMMENT '员工工号',
    `employment_status` varchar(60) NOT NULL COMMENT '就职状态[active, notactive, disabled]',
    `employment_type` varchar(60) DEFAULT NULL COMMENT '就职类型[permanent, intern]',
    `phone` varchar(200) DEFAULT NULL COMMENT '手机号',
    `telephone` varchar(200) DEFAULT NULL COMMENT '座机号',
    `source` varchar(20) DEFAULT 'sync' COMMENT '来源, buildin/sync',
    `custom_fields` varchar(5000) DEFAULT NULL COMMENT '自定义字段，json数组',
    `sync_type` varchar(20) DEFAULT 'auto' COMMENT '同步方式，auto/manual',
    `update_type` varchar(20) NOT NULL COMMENT '修改类型, user_del/user_update/user_add',
    `status` int DEFAULT '0' COMMENT '0-默认状态，1-已同步 -1:同步失败',
    `msg` varchar(2000) DEFAULT NULL COMMENT '错误详情',
    `operator` varchar(100) NOT NULL DEFAULT '系统' COMMENT 'operator',
    `sync_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '增量数据变动时间',
    `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE COMMENT '主键索引',
    KEY `idx_mtime` (`mtime`) USING BTREE COMMENT 'mtime索引',
    KEY `idx_sync_time` (`sync_time`, `status`, `third_company_id`) USING BTREE COMMENT 'sync_time索引',
    KEY `idx_nick_name` (`nick_name`) USING BTREE COMMENT 'nick_name索引'
) ENGINE = InnoDB AUTO_INCREMENT = 55 DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '三方用户增量表'

-- incre department

CREATE TABLE `tb_las_department_increment` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `did` varchar(255) NOT NULL COMMENT '部门id',
    `third_company_id` varchar(20) NOT NULL COMMENT '租户id',
    `platform_id` varchar(60) NOT NULL COMMENT '平台id, 用来区分多种数据源，platform_id + did 唯一, 根部门例外',
    `pid` varchar(255) DEFAULT NULL COMMENT '父部门id',
    `name` varchar(255) NOT NULL COMMENT '部门名称',
    `order` int DEFAULT NULL COMMENT '排序',
    `source` varchar(20) DEFAULT 'sync' COMMENT '来源, buildin/sync',
    `sync_type` varchar(20) DEFAULT 'auto' COMMENT '同步方式，auto/manual',
    `update_type` varchar(20) NOT NULL COMMENT '修改类型, dept_del/dept_update/dept_add',
    `status` int DEFAULT '0' COMMENT '0-默认状态，1-已同步 -1:同步失败',
    `msg` varchar(2000) DEFAULT NULL COMMENT '错误详情',
    `operator` varchar(100) NOT NULL DEFAULT '系统' COMMENT 'operator',
    `sync_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '增量数据变动时间',
    `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `type` varchar(255) DEFAULT NULL COMMENT '类型, 仅支持小写字母和下划线组成',
    PRIMARY KEY (`id`) USING BTREE COMMENT '主键索引',
    KEY `idx_sync_time` (`sync_time`, `status`, `third_company_id`) USING BTREE COMMENT 'sync_time索引',
    KEY `idx_mtime` (`mtime`) USING BTREE COMMENT 'mtime索引',
    KEY `idx_name` (`name`) USING BTREE COMMENT 'name索引'
) ENGINE = InnoDB AUTO_INCREMENT = 83 DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '三方部门表'


-- increment relation

CREATE TABLE `tb_las_department_user_increment` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `third_company_id` varchar(20) NOT NULL COMMENT '租户id',
    `platform_id` varchar(60) NOT NULL COMMENT '平台id, 用来区分多种数据源，platform_id + id 唯一',
    `uid` varchar(255) NOT NULL COMMENT '用户id',
    `did` varchar(255) NOT NULL COMMENT '默认部门id',
    `order` int DEFAULT NULL COMMENT '用户在部门下的排序',
    `main` int DEFAULT '0' COMMENT '是否是主部门，1：是 0：不是',
    `sync_type` varchar(20) DEFAULT 'auto' COMMENT '同步方式，auto/manual',
    `update_type` varchar(20) NOT NULL COMMENT '修改类型, user_dept_del/user_dept_update/user_dept_add',
    `status` int DEFAULT '0' COMMENT '0-默认状态，1-已同步 -1:同步失败',
    `msg` varchar(2000) DEFAULT NULL COMMENT '错误详情',
    `operator` varchar(100) NOT NULL DEFAULT '系统' COMMENT 'operator',
    `sync_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '增量数据变动时间',
    `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `dids` varchar(5000) DEFAULT NULL COMMENT 'dids, JSONArray, [{"did": 1, "order": 1}]',
    PRIMARY KEY (`id`) USING BTREE COMMENT '主键索引',
    KEY `idx_sync_time` (`sync_time`, `status`, `third_company_id`) USING BTREE COMMENT 'sync_time索引',
    KEY `idx_mtime` (`mtime`) USING BTREE COMMENT 'mtime索引'
) ENGINE = InnoDB AUTO_INCREMENT = 80 DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '三方部门用户关系表'