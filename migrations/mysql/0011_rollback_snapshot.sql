CREATE TABLE IF NOT EXISTS sg_state_snapshot (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    snapshot_scope VARCHAR(16) NOT NULL COMMENT '快照范围：GROUP/GLOBAL',
    snapshot_type VARCHAR(16) NOT NULL COMMENT '快照类型：AUTO/MANUAL/SAFETY',
    trigger_code VARCHAR(64) NOT NULL COMMENT '触发节点',
    target_group_id BIGINT NULL COMMENT '单组快照所属组，全局快照为空',
    target_year_no INT NULL COMMENT '快照对应年份',
    target_stage_code VARCHAR(32) NULL COMMENT '快照对应阶段',
    target_report_status VARCHAR(32) NULL COMMENT '快照对应财报状态',
    description VARCHAR(500) NULL COMMENT '管理员说明或系统说明',
    payload_hash VARCHAR(128) NOT NULL COMMENT '快照载荷哈希',
    created_by_id BIGINT NOT NULL COMMENT '创建人 ID',
    created_by_name VARCHAR(64) NOT NULL COMMENT '创建人名称',
    created_at DATETIME NOT NULL COMMENT '创建时间',
    PRIMARY KEY (id),
    KEY idx_snapshot_scope_type (snapshot_scope, snapshot_type),
    KEY idx_snapshot_group_year (target_group_id, target_year_no),
    KEY idx_snapshot_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='状态快照主表';

CREATE TABLE IF NOT EXISTS sg_state_snapshot_payload (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    snapshot_id BIGINT NOT NULL COMMENT '快照主表 ID',
    payload_version VARCHAR(32) NOT NULL COMMENT '快照载荷结构版本',
    payload_json JSON NOT NULL COMMENT '快照业务载荷',
    payload_size INT NOT NULL DEFAULT 0 COMMENT '载荷大小',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_snapshot_payload (snapshot_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='状态快照载荷表';

CREATE TABLE IF NOT EXISTS sg_rollback_log (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    rollback_type VARCHAR(32) NOT NULL COMMENT '回退类型：UNLOCK_RETRY/GROUP_SNAPSHOT_RESTORE',
    target_group_id BIGINT NOT NULL COMMENT '被修正小组 ID',
    target_year_no INT NOT NULL COMMENT '目标年份',
    target_stage_code VARCHAR(32) NULL COMMENT '目标阶段',
    snapshot_id BIGINT NULL COMMENT '恢复来源快照 ID',
    safety_snapshot_id BIGINT NOT NULL COMMENT '回退前安全快照 ID',
    reason VARCHAR(500) NOT NULL COMMENT '管理员填写原因',
    state_before_json JSON NOT NULL COMMENT '执行前状态摘要',
    state_after_json JSON NOT NULL COMMENT '执行后状态摘要',
    operator_id BIGINT NOT NULL COMMENT '操作管理员 ID',
    operator_name VARCHAR(64) NOT NULL COMMENT '操作管理员名称',
    operate_time DATETIME NOT NULL COMMENT '操作时间',
    PRIMARY KEY (id),
    KEY idx_rollback_group_year (target_group_id, target_year_no),
    KEY idx_rollback_safety_snapshot (safety_snapshot_id),
    KEY idx_rollback_operate_time (operate_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='回退与修正日志表';

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_year_state' AND COLUMN_NAME = 'rollback_pending'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_year_state ADD COLUMN rollback_pending TINYINT(1) NOT NULL DEFAULT 0 COMMENT ''是否处于回退补提或待重提'' AFTER latest_report_submit_version', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_year_state' AND COLUMN_NAME = 'rollback_target_year_no'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_year_state ADD COLUMN rollback_target_year_no INT NULL COMMENT ''最近一次回退目标年份'' AFTER rollback_pending', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_year_state' AND COLUMN_NAME = 'rollback_target_stage_code'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_year_state ADD COLUMN rollback_target_stage_code VARCHAR(32) NULL COMMENT ''最近一次回退目标阶段'' AFTER rollback_target_year_no', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_year_state' AND COLUMN_NAME = 'rollback_log_id'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_year_state ADD COLUMN rollback_log_id BIGINT NULL COMMENT ''最近一次回退日志 ID'' AFTER rollback_target_stage_code', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_year_state' AND INDEX_NAME = 'idx_rollback_pending'
);
SET @ddl := IF(@index_exists = 0, 'ALTER TABLE sg_group_year_state ADD INDEX idx_rollback_pending (rollback_pending)', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_summary_snapshot' AND COLUMN_NAME = 'invalidated_by_rollback_id'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_summary_snapshot ADD COLUMN invalidated_by_rollback_id BIGINT NULL COMMENT ''被回退置为失效时对应回退日志'' AFTER source_report_submit_version', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_summary_snapshot' AND COLUMN_NAME = 'invalidated_at'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_summary_snapshot ADD COLUMN invalidated_at DATETIME NULL COMMENT ''失效时间'' AFTER invalidated_by_rollback_id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_adjustment' AND COLUMN_NAME = 'effective'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_adjustment ADD COLUMN effective TINYINT(1) NOT NULL DEFAULT 1 COMMENT ''是否仍参与经营财报汇总计算'' AFTER reason', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_adjustment' AND COLUMN_NAME = 'invalidated_by_rollback_id'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_adjustment ADD COLUMN invalidated_by_rollback_id BIGINT NULL COMMENT ''被回退置为失效时对应回退日志'' AFTER effective', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_adjustment' AND COLUMN_NAME = 'invalid_reason'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_adjustment ADD COLUMN invalid_reason VARCHAR(255) NULL COMMENT ''失效原因'' AFTER invalidated_by_rollback_id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_adjustment' AND COLUMN_NAME = 'invalidated_at'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_adjustment ADD COLUMN invalidated_at DATETIME NULL COMMENT ''失效时间'' AFTER invalid_reason', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_adjustment' AND INDEX_NAME = 'idx_adjustment_effective'
);
SET @ddl := IF(@index_exists = 0, 'ALTER TABLE sg_group_adjustment ADD INDEX idx_adjustment_effective (effective)', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_admin_unlock_log' AND COLUMN_NAME = 'unlock_target_type'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_admin_unlock_log ADD COLUMN unlock_target_type VARCHAR(32) NULL COMMENT ''解锁目标类型'' AFTER reason', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_admin_unlock_log' AND COLUMN_NAME = 'target_stage_code'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_admin_unlock_log ADD COLUMN target_stage_code VARCHAR(32) NULL COMMENT ''经营页回退目标阶段'' AFTER unlock_target_type', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_admin_unlock_log' AND COLUMN_NAME = 'safety_snapshot_id'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_admin_unlock_log ADD COLUMN safety_snapshot_id BIGINT NULL COMMENT ''执行前安全快照 ID'' AFTER target_stage_code', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_order_selection' AND COLUMN_NAME = 'delivery_effective'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_order_selection ADD COLUMN delivery_effective TINYINT(1) NOT NULL DEFAULT 1 COMMENT ''交付状态是否仍有效'' AFTER delivered_at', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_order_selection' AND COLUMN_NAME = 'invalidated_by_rollback_id'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_order_selection ADD COLUMN invalidated_by_rollback_id BIGINT NULL COMMENT ''交付状态被回退置为失效时对应回退日志'' AFTER delivery_effective', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_order_selection' AND COLUMN_NAME = 'invalidated_at'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_order_selection ADD COLUMN invalidated_at DATETIME NULL COMMENT ''交付状态失效时间'' AFTER invalidated_by_rollback_id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_order_selection' AND INDEX_NAME = 'idx_order_delivery_effective'
);
SET @ddl := IF(@index_exists = 0, 'ALTER TABLE sg_group_order_selection ADD INDEX idx_order_delivery_effective (delivery_effective)', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
