ALTER TABLE sg_game_config
    ADD COLUMN edition_code VARCHAR(64) NOT NULL DEFAULT 'VIP_SERVICE_V1' COMMENT '沙盘版本包编码' AFTER current_open_year,
    ADD COLUMN edition_name VARCHAR(64) NOT NULL DEFAULT '贵宾服务版 V1' COMMENT '沙盘版本包显示名' AFTER edition_code,
    MODIFY COLUMN rule_version VARCHAR(64) NOT NULL DEFAULT 'COMMON_FORMULA_V1' COMMENT '公式规则版本',
    MODIFY COLUMN template_version VARCHAR(64) NOT NULL DEFAULT 'VIP_SERVICE_V1' COMMENT '总模板版本',
    ADD COLUMN operating_template_version VARCHAR(64) NOT NULL DEFAULT 'VIP_OPERATING_TEMPLATE_V1' COMMENT '经营页字段模板版本' AFTER template_version,
    ADD COLUMN report_template_version VARCHAR(64) NOT NULL DEFAULT 'VIP_REPORT_TEMPLATE_V1' COMMENT '财报页字段模板版本' AFTER operating_template_version,
    ADD COLUMN order_template_version VARCHAR(64) NOT NULL DEFAULT 'VIP_ORDER_TEMPLATE_V1' COMMENT '订单字段模板版本' AFTER report_template_version,
    ADD COLUMN process_rule_version VARCHAR(64) NOT NULL DEFAULT 'COMMON_PROCESS_V1' COMMENT '流程规则版本' AFTER order_template_version;

UPDATE sg_game_config
SET
    edition_code = 'VIP_SERVICE_V1',
    edition_name = '贵宾服务版 V1',
    rule_version = 'COMMON_FORMULA_V1',
    template_version = 'VIP_SERVICE_V1',
    operating_template_version = 'VIP_OPERATING_TEMPLATE_V1',
    report_template_version = 'VIP_REPORT_TEMPLATE_V1',
    order_template_version = 'VIP_ORDER_TEMPLATE_V1',
    process_rule_version = 'COMMON_PROCESS_V1';
