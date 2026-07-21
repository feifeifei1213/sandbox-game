SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_adjustment' AND COLUMN_NAME = 'voided_by_id'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_adjustment ADD COLUMN voided_by_id BIGINT NULL COMMENT ''作废操作管理员 ID'' AFTER effective', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_adjustment' AND COLUMN_NAME = 'voided_by_name'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_adjustment ADD COLUMN voided_by_name VARCHAR(64) NULL COMMENT ''作废操作管理员名称'' AFTER voided_by_id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_adjustment' AND COLUMN_NAME = 'void_reason'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_adjustment ADD COLUMN void_reason VARCHAR(500) NULL COMMENT ''作废原因'' AFTER voided_by_name', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_adjustment' AND COLUMN_NAME = 'voided_at'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_adjustment ADD COLUMN voided_at DATETIME NULL COMMENT ''作废时间'' AFTER void_reason', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_adjustment' AND INDEX_NAME = 'idx_group_year_effective'
);
SET @ddl := IF(@index_exists = 0, 'ALTER TABLE sg_group_adjustment ADD INDEX idx_group_year_effective (group_id, year_no, effective)', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS sg_group_adjustment_revision (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    group_id BIGINT NOT NULL COMMENT '目标小组',
    year_no INT NOT NULL COMMENT '目标年份',
    revision BIGINT NOT NULL DEFAULT 0 COMMENT '奖惩增量同步版本',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最近变更时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_adjustment_revision_group_year (group_id, year_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='小组年度奖惩增量同步版本';
