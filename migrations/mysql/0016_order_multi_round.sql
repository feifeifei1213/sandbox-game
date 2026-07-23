SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_market_bidding_state' AND COLUMN_NAME = 'current_round_no'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_market_bidding_state ADD COLUMN current_round_no INT NOT NULL DEFAULT 0 COMMENT ''当前或最近处理轮次'' AFTER current_group_id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_market_bidding_state' AND COLUMN_NAME = 'completion_reason'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_market_bidding_state ADD COLUMN completion_reason VARCHAR(32) NULL COMMENT ''标段完成原因'' AFTER current_round_no', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_market_selection_order' AND COLUMN_NAME = 'round_no'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_market_selection_order ADD COLUMN round_no INT NOT NULL DEFAULT 1 COMMENT ''选单轮次'' AFTER order_type', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_market_selection_order' AND INDEX_NAME = 'uk_market_sequence'
);
SET @ddl := IF(@index_exists > 0, 'ALTER TABLE sg_market_selection_order DROP INDEX uk_market_sequence', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_market_selection_order' AND INDEX_NAME = 'uk_group_market_sequence'
);
SET @ddl := IF(@index_exists > 0, 'ALTER TABLE sg_market_selection_order DROP INDEX uk_group_market_sequence', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_market_selection_order' AND INDEX_NAME = 'uk_market_round_sequence'
);
SET @ddl := IF(@index_exists = 0, 'ALTER TABLE sg_market_selection_order ADD UNIQUE KEY uk_market_round_sequence (year_no, market_code, order_type, round_no, sequence_no)', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_market_selection_order' AND INDEX_NAME = 'uk_group_market_round'
);
SET @ddl := IF(@index_exists = 0, 'ALTER TABLE sg_market_selection_order ADD UNIQUE KEY uk_group_market_round (year_no, market_code, order_type, round_no, group_id)', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_order_selection' AND COLUMN_NAME = 'round_no'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_order_selection ADD COLUMN round_no INT NOT NULL DEFAULT 1 COMMENT ''选择轮次'' AFTER order_type', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_order_selection' AND COLUMN_NAME = 'selection_order_id'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_order_selection ADD COLUMN selection_order_id BIGINT NULL COMMENT ''选单顺序记录 ID'' AFTER round_no', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_order_selection' AND INDEX_NAME = 'uk_group_year_segment_selection'
);
SET @ddl := IF(@index_exists > 0, 'ALTER TABLE sg_group_order_selection DROP INDEX uk_group_year_segment_selection', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_order_selection' AND INDEX_NAME = 'uk_selection_order'
);
SET @ddl := IF(@index_exists = 0, 'ALTER TABLE sg_group_order_selection ADD UNIQUE KEY uk_selection_order (selection_order_id)', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
