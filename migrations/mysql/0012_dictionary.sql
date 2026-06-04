SET @column_exists := (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sg_game_config' AND COLUMN_NAME = 'dictionary_revision'
);
SET @ddl := IF(@column_exists = 0, 'ALTER TABLE sg_game_config ADD COLUMN dictionary_revision INT NOT NULL DEFAULT 0 COMMENT ''当前比赛业务显示字典版本'' AFTER process_rule_version', 'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS sg_dictionary_scheme (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    edition_code VARCHAR(64) NOT NULL COMMENT '绑定沙盘版本包编码',
    scheme_name VARCHAR(64) NOT NULL COMMENT '字典方案名称',
    description VARCHAR(255) NOT NULL DEFAULT '' COMMENT '方案说明',
    built_in TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统内置方案',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_dictionary_scheme_name (edition_code, scheme_name),
    KEY idx_dictionary_scheme_edition (edition_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='业务显示字典方案表';

CREATE TABLE IF NOT EXISTS sg_dictionary_scheme_item (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    scheme_id BIGINT NOT NULL COMMENT '字典方案 ID',
    edition_code VARCHAR(64) NOT NULL COMMENT '绑定沙盘版本包编码',
    item_code VARCHAR(128) NOT NULL COMMENT '稳定显示项编码',
    item_category VARCHAR(32) NOT NULL COMMENT '显示项类别',
    default_name VARCHAR(128) NOT NULL COMMENT '版本默认显示名',
    display_name VARCHAR(128) NOT NULL COMMENT '方案显示名',
    display_order INT NOT NULL DEFAULT 0 COMMENT '展示排序',
    editable TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否允许管理员改显示名',
    related_payload JSON NULL COMMENT '关联前端字段或业务扩展信息',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_dictionary_scheme_item (scheme_id, item_code),
    KEY idx_dictionary_scheme_item_edition (edition_code),
    KEY idx_dictionary_scheme_item_category (item_category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='业务显示字典方案明细表';

CREATE TABLE IF NOT EXISTS sg_current_dictionary_item (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    edition_code VARCHAR(64) NOT NULL COMMENT '当前比赛沙盘版本包编码',
    item_code VARCHAR(128) NOT NULL COMMENT '稳定显示项编码',
    item_category VARCHAR(32) NOT NULL COMMENT '显示项类别',
    default_name VARCHAR(128) NOT NULL COMMENT '版本默认显示名',
    display_name VARCHAR(128) NOT NULL COMMENT '当前比赛显示名',
    display_order INT NOT NULL DEFAULT 0 COMMENT '展示排序',
    editable TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否允许管理员改显示名',
    related_payload JSON NULL COMMENT '关联前端字段或业务扩展信息',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_current_dictionary_item (item_code),
    KEY idx_current_dictionary_edition (edition_code),
    KEY idx_current_dictionary_category (item_category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='当前比赛业务显示字典快照表';

CREATE TABLE IF NOT EXISTS sg_dictionary_change_log (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    edition_code VARCHAR(64) NOT NULL COMMENT '沙盘版本包编码',
    change_type VARCHAR(32) NOT NULL COMMENT '变更类型',
    scheme_id BIGINT NULL COMMENT '关联方案 ID',
    reason VARCHAR(255) NOT NULL DEFAULT '' COMMENT '变更原因或说明',
    before_json JSON NOT NULL COMMENT '变更前快照',
    after_json JSON NOT NULL COMMENT '变更后快照',
    revision INT NOT NULL DEFAULT 0 COMMENT '变更后字典版本',
    operator_id BIGINT NOT NULL COMMENT '操作管理员 ID',
    operator_name VARCHAR(64) NOT NULL COMMENT '操作管理员名称',
    operate_time DATETIME NOT NULL COMMENT '操作时间',
    PRIMARY KEY (id),
    KEY idx_dictionary_log_edition (edition_code),
    KEY idx_dictionary_log_operate_time (operate_time),
    KEY idx_dictionary_log_revision (revision)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='业务显示字典修改日志表';
