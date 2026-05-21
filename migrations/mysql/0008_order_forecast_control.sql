CREATE TABLE IF NOT EXISTS sg_order_forecast_control (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    year_no INT NOT NULL COMMENT '预测年份，固定 1~8',
    forecast_stage_code VARCHAR(32) NOT NULL COMMENT '预测阶段：YEAR_1_3/YEAR_4_5/YEAR_6_8',
    market_code VARCHAR(32) NOT NULL COMMENT '市场：LOCAL/REGIONAL/NATIONAL/GLOBAL',
    order_type VARCHAR(32) NOT NULL COMMENT '订单类型',
    order_count INT NOT NULL DEFAULT 0 COMMENT '订单卡片数量，范围 0~15',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_order_forecast_control (year_no, market_code, order_type),
    KEY idx_order_forecast_stage (forecast_stage_code, market_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='多年订单数量控制台表';

CREATE TABLE IF NOT EXISTS sg_order_market_forecast (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    forecast_stage_code VARCHAR(32) NOT NULL COMMENT '预测阶段',
    market_code VARCHAR(32) NOT NULL COMMENT '市场',
    forecast_data_json JSON NOT NULL COMMENT '该阶段市场预测数据',
    narrative TEXT NULL COMMENT '市场预测说明文字',
    formula_version VARCHAR(64) NOT NULL COMMENT '预测公式版本',
    random_seed VARCHAR(64) NOT NULL DEFAULT '' COMMENT '预测快照种子或版本摘要',
    control_snapshot_json JSON NOT NULL COMMENT '生成预测时的控制台快照',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_order_market_forecast (forecast_stage_code, market_code),
    KEY idx_order_market_forecast_stage (forecast_stage_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='订单市场预测快照表';

ALTER TABLE sg_order_generation_batch
    ADD COLUMN forecast_snapshot_json JSON NULL COMMENT '市场预测快照摘要' AFTER control_snapshot_json;

INSERT INTO sg_order_forecast_control (
    year_no,
    forecast_stage_code,
    market_code,
    order_type,
    order_count,
    creator,
    updater
)
SELECT
    year_no,
    CASE
        WHEN year_no BETWEEN 1 AND 3 THEN 'YEAR_1_3'
        WHEN year_no BETWEEN 4 AND 5 THEN 'YEAR_4_5'
        ELSE 'YEAR_6_8'
    END AS forecast_stage_code,
    market_code,
    order_type,
    order_count,
    'migration',
    'migration'
FROM sg_order_generation_config
WHERE year_no BETWEEN 1 AND 8
ON DUPLICATE KEY UPDATE
    order_count = VALUES(order_count),
    forecast_stage_code = VALUES(forecast_stage_code),
    updater = 'migration';
