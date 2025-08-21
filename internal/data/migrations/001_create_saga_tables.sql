-- 创建 Saga 事务表（对应文档中的 saga_transactions）
CREATE TABLE IF NOT EXISTS `saga_transactions` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `transaction_id` varchar(64) NOT NULL COMMENT 'Saga 事务唯一标识',
    `name` varchar(255) NOT NULL COMMENT '事务名称',
    `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT 'Saga 状态',
    `current_step` varchar(64) DEFAULT NULL COMMENT '当前步骤',
    `progress` int NOT NULL DEFAULT '0' COMMENT '进度百分比',
    `start_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    `end_time` timestamp NULL DEFAULT NULL COMMENT '结束时间',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_transaction_id` (`transaction_id`),
    KEY `idx_status` (`status`),
    KEY `idx_start_time` (`start_time`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Saga 事务表';

-- 创建 Saga 步骤表（对应文档中的 saga_steps）
CREATE TABLE IF NOT EXISTS `saga_steps` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `step_id` varchar(64) NOT NULL COMMENT '步骤唯一标识',
    `transaction_id` varchar(64) NOT NULL COMMENT '关联的 Saga 事务ID',
    `step_name` varchar(255) NOT NULL COMMENT '步骤名称',
    `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '步骤状态',
    `action_data` json DEFAULT NULL COMMENT '操作数据（JSON格式）',
    `compensate_data` json DEFAULT NULL COMMENT '补偿数据（JSON格式）',
    `error_message` text COMMENT '错误信息',
    `retry_count` int NOT NULL DEFAULT '0' COMMENT '重试次数',
    `max_retries` int NOT NULL DEFAULT '3' COMMENT '最大重试次数',
    `start_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    `end_time` timestamp NULL DEFAULT NULL COMMENT '结束时间',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_step_id` (`step_id`),
    KEY `idx_transaction_id` (`transaction_id`),
    KEY `idx_status` (`status`),
    KEY `idx_start_time` (`start_time`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_saga_steps_transaction` FOREIGN KEY (`transaction_id`) REFERENCES `saga_transactions` (`transaction_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Saga 步骤表';

-- 创建 Saga 事件表（对应文档中的 saga_events）
CREATE TABLE IF NOT EXISTS `saga_events` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `transaction_id` varchar(64) NOT NULL COMMENT '关联的 Saga 事务ID',
    `step_id` varchar(64) DEFAULT NULL COMMENT '关联的步骤ID（可选）',
    `event_type` varchar(50) NOT NULL COMMENT '事件类型',
    `event_data` json DEFAULT NULL COMMENT '事件数据（JSON格式）',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_transaction_id` (`transaction_id`),
    KEY `idx_step_id` (`step_id`),
    KEY `idx_event_type` (`event_type`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_saga_events_transaction` FOREIGN KEY (`transaction_id`) REFERENCES `saga_transactions` (`transaction_id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_saga_events_step` FOREIGN KEY (`step_id`) REFERENCES `saga_steps` (`step_id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Saga 事件表';

-- 创建索引优化查询性能
CREATE INDEX IF NOT EXISTS `idx_saga_transactions_status_created` ON `saga_transactions` (`status`, `created_at`);
CREATE INDEX IF NOT EXISTS `idx_saga_steps_transaction_status` ON `saga_steps` (`transaction_id`, `status`);
CREATE INDEX IF NOT EXISTS `idx_saga_steps_transaction_start_time` ON `saga_steps` (`transaction_id`, `start_time`);
CREATE INDEX IF NOT EXISTS `idx_saga_events_transaction_created` ON `saga_events` (`transaction_id`, `created_at`);
CREATE INDEX IF NOT EXISTS `idx_saga_events_type_created` ON `saga_events` (`event_type`, `created_at`);