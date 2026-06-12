SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_order_generation_config' AND COLUMN_NAME = 'order_template_version'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_order_generation_config ADD COLUMN order_template_version VARCHAR(64) NOT NULL DEFAULT ''VIP_ORDER_TEMPLATE_V1'' COMMENT ''订单模板版本'' AFTER id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_order_forecast_control' AND COLUMN_NAME = 'order_template_version'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_order_forecast_control ADD COLUMN order_template_version VARCHAR(64) NOT NULL DEFAULT ''VIP_ORDER_TEMPLATE_V1'' COMMENT ''订单模板版本'' AFTER id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_order_market_forecast' AND COLUMN_NAME = 'order_template_version'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_order_market_forecast ADD COLUMN order_template_version VARCHAR(64) NOT NULL DEFAULT ''VIP_ORDER_TEMPLATE_V1'' COMMENT ''订单模板版本'' AFTER id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_order_market_config' AND COLUMN_NAME = 'order_template_version'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_order_market_config ADD COLUMN order_template_version VARCHAR(64) NOT NULL DEFAULT ''VIP_ORDER_TEMPLATE_V1'' COMMENT ''订单模板版本'' AFTER id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_order_generation_batch' AND COLUMN_NAME = 'order_template_version'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_order_generation_batch ADD COLUMN order_template_version VARCHAR(64) NOT NULL DEFAULT ''VIP_ORDER_TEMPLATE_V1'' COMMENT ''订单模板版本'' AFTER id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_order_pool' AND COLUMN_NAME = 'order_template_version'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_order_pool ADD COLUMN order_template_version VARCHAR(64) NOT NULL DEFAULT ''VIP_ORDER_TEMPLATE_V1'' COMMENT ''订单模板版本'' AFTER id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_order_pool' AND COLUMN_NAME = 'order_payload_json'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_order_pool ADD COLUMN order_payload_json JSON NULL COMMENT ''订单模板专属扩展字段'' AFTER source_row_key', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_market_bidding_state' AND COLUMN_NAME = 'order_template_version'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_market_bidding_state ADD COLUMN order_template_version VARCHAR(64) NOT NULL DEFAULT ''VIP_ORDER_TEMPLATE_V1'' COMMENT ''订单模板版本'' AFTER id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_market_bid' AND COLUMN_NAME = 'order_template_version'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_market_bid ADD COLUMN order_template_version VARCHAR(64) NOT NULL DEFAULT ''VIP_ORDER_TEMPLATE_V1'' COMMENT ''订单模板版本'' AFTER id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_market_selection_order' AND COLUMN_NAME = 'order_template_version'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_market_selection_order ADD COLUMN order_template_version VARCHAR(64) NOT NULL DEFAULT ''VIP_ORDER_TEMPLATE_V1'' COMMENT ''订单模板版本'' AFTER id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_group_order_selection' AND COLUMN_NAME = 'order_template_version'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_group_order_selection ADD COLUMN order_template_version VARCHAR(64) NOT NULL DEFAULT ''VIP_ORDER_TEMPLATE_V1'' COMMENT ''订单模板版本'' AFTER id', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
