CREATE TABLE IF NOT EXISTS sg_order_market_config (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    year_no INT NOT NULL COMMENT '年份',
    market_code VARCHAR(32) NOT NULL COMMENT '市场：LOCAL/REGIONAL/NATIONAL/GLOBAL',
    market_enabled TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否开启该市场',
    config_status VARCHAR(32) NOT NULL DEFAULT 'DRAFT' COMMENT '配置状态：DRAFT/LOCKED',
    locked_batch_id BIGINT NULL COMMENT '锁定时对应确认订单池批次',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_order_market_config (year_no, market_code),
    KEY idx_order_market_enabled (year_no, market_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='年度订单市场开启配置表';

ALTER TABLE sg_order_pool
    ADD COLUMN business_order_no VARCHAR(64) NULL COMMENT '业务订单编号/页面展示编号' AFTER card_sequence_no;

ALTER TABLE sg_order_generation_config
    DROP INDEX uk_order_release_sequence,
    ADD KEY idx_order_release_sequence (year_no, release_sequence_no);

ALTER TABLE sg_market_bidding_state
    DROP INDEX uk_segment_release_sequence,
    ADD KEY idx_segment_release_sequence (year_no, release_sequence_no);
