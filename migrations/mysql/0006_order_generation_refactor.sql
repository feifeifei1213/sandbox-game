CREATE TABLE IF NOT EXISTS sg_order_generation_batch (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    year_no INT NOT NULL COMMENT '年份',
    batch_status VARCHAR(32) NOT NULL COMMENT '批次状态：PREVIEW/CONFIRMED/VOID',
    formula_version VARCHAR(32) NOT NULL COMMENT '订单生成公式版本',
    random_seed VARCHAR(64) NOT NULL COMMENT '随机种子',
    control_snapshot_json JSON NOT NULL COMMENT '订单数量与释放顺序快照',
    parameter_snapshot_json JSON NOT NULL COMMENT '公式参数快照',
    order_detail_json JSON NOT NULL COMMENT '生成订单明细快照',
    generated_order_count INT NOT NULL DEFAULT 0 COMMENT '生成订单数量',
    generated_by_id BIGINT NOT NULL COMMENT '生成管理员 ID',
    generated_by_name VARCHAR(64) NOT NULL COMMENT '生成管理员名称',
    generated_at DATETIME NOT NULL COMMENT '生成时间',
    confirmed_by_id BIGINT NULL COMMENT '确认管理员 ID',
    confirmed_by_name VARCHAR(64) NULL COMMENT '确认管理员名称',
    confirmed_at DATETIME NULL COMMENT '确认时间',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_order_generation_batch_year (year_no, batch_status),
    KEY idx_order_generation_batch_generated_at (generated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='订单生成预览与确认批次表';

ALTER TABLE sg_order_generation_config
    ADD COLUMN generation_batch_id BIGINT NULL COMMENT '当前确认订单池批次' AFTER release_sequence_no;

ALTER TABLE sg_order_pool
    ADD COLUMN segment_code VARCHAR(64) NULL COMMENT '标段编码' AFTER order_type,
    ADD COLUMN card_sequence_no INT NULL COMMENT '订单卡片序号' AFTER segment_code,
    ADD COLUMN generation_batch_id BIGINT NULL COMMENT '订单生成批次' AFTER selected_group_id,
    ADD COLUMN selected_at DATETIME NULL COMMENT '选中时间' AFTER selected_group_id,
    ADD COLUMN source_row_key VARCHAR(128) NULL COMMENT '公式链卡片标识' AFTER source_row_index;

ALTER TABLE sg_group_market_bid
    ADD COLUMN order_type VARCHAR(32) NOT NULL DEFAULT 'AGENCY_INSPECTION' COMMENT '订单类型' AFTER market_code,
    ADD COLUMN bid_status VARCHAR(32) NOT NULL DEFAULT 'SUBMITTED' COMMENT '投入状态：SUBMITTED/LOCKED' AFTER market_investment;

ALTER TABLE sg_group_market_bid
    DROP INDEX uk_group_year_market_bid;

ALTER TABLE sg_group_market_bid
    ADD UNIQUE KEY uk_group_year_segment_bid (group_id, year_no, market_code, order_type),
    ADD KEY idx_year_segment_bid (year_no, market_code, order_type);
