-- I3-02 单组演练环境种子数据
-- 说明：
-- 1. 本脚本仅用于单组演练库，不替代正式比赛或多组演练库初始化。
-- 2. 默认只初始化 1 个管理员账号和 1 个玩家组账号，便于重复执行 0年 -> 最终年 的完整推演。
-- 3. 初始基线默认未提交，方便从“管理员提交初始基线”开始走完整主链。

INSERT INTO sg_group (
    id, group_no, group_code, group_name, business_status, bankrupt_year_no, bankrupt_reason,
    creator, create_time, updater, update_time
)
VALUES
    (1, 1, 'GROUP_01', '第一组', 'NORMAL', NULL, NULL, 'seed-single', NOW(), 'seed-single', NOW())
ON DUPLICATE KEY UPDATE
    group_name = VALUES(group_name),
    business_status = VALUES(business_status),
    bankrupt_year_no = VALUES(bankrupt_year_no),
    bankrupt_reason = VALUES(bankrupt_reason),
    updater = 'seed-single',
    update_time = NOW();

INSERT INTO sg_account (
    id, username, password_hash, role_type, group_id, status, last_login_time,
    creator, create_time, updater, update_time
)
VALUES
    (1, 'admin', CONCAT('{sha256}', SHA2('123456', 256)), 'ADMIN', NULL, 'ENABLED', NULL, 'seed-single', NOW(), 'seed-single', NOW()),
    (101, 'group01', CONCAT('{sha256}', SHA2('123456', 256)), 'GROUP', 1, 'ENABLED', NULL, 'seed-single', NOW(), 'seed-single', NOW())
ON DUPLICATE KEY UPDATE
    password_hash = VALUES(password_hash),
    role_type = VALUES(role_type),
    group_id = VALUES(group_id),
    status = VALUES(status),
    updater = 'seed-single',
    update_time = NOW();

INSERT INTO sg_game_config (
    id, final_year, current_open_year, edition_code, edition_name, rule_version, template_version,
    operating_template_version, report_template_version, order_template_version, process_rule_version,
    initial_baseline_submitted,
    creator, create_time, updater, update_time
)
VALUES
    (
        1, 8, 0, 'VIP_SERVICE_V1', '贵宾服务版 V1', 'COMMON_FORMULA_V1', 'VIP_SERVICE_V1',
        'VIP_OPERATING_TEMPLATE_V1', 'VIP_REPORT_TEMPLATE_V1', 'VIP_ORDER_TEMPLATE_V1', 'COMMON_PROCESS_V1',
        0, 'seed-single', NOW(), 'seed-single', NOW()
    )
ON DUPLICATE KEY UPDATE
    final_year = VALUES(final_year),
    current_open_year = VALUES(current_open_year),
    edition_code = VALUES(edition_code),
    edition_name = VALUES(edition_name),
    rule_version = VALUES(rule_version),
    template_version = VALUES(template_version),
    operating_template_version = VALUES(operating_template_version),
    report_template_version = VALUES(report_template_version),
    order_template_version = VALUES(order_template_version),
    process_rule_version = VALUES(process_rule_version),
    initial_baseline_submitted = VALUES(initial_baseline_submitted),
    updater = 'seed-single',
    update_time = NOW();

INSERT INTO sg_group_year_state (
    group_id, year_no, year_type, year_status, stage_status, report_status,
    summary_effective, latest_stage_submit_version, latest_report_submit_version,
    creator, create_time, updater, update_time
)
SELECT
    1,
    y.year_no,
    CASE WHEN y.year_no = 0 THEN 'DEMO' ELSE 'FORMAL' END AS year_type,
    CASE WHEN y.year_no = 0 THEN 'OPERATING' ELSE 'LOCKED' END AS year_status,
    'Q1_OPEN' AS stage_status,
    'REPORT_LOCKED' AS report_status,
    0 AS summary_effective,
    0 AS latest_stage_submit_version,
    0 AS latest_report_submit_version,
    'seed-single' AS creator,
    NOW() AS create_time,
    'seed-single' AS updater,
    NOW() AS update_time
FROM (
    SELECT 0 AS year_no UNION ALL
    SELECT 1 UNION ALL
    SELECT 2 UNION ALL
    SELECT 3 UNION ALL
    SELECT 4 UNION ALL
    SELECT 5 UNION ALL
    SELECT 6 UNION ALL
    SELECT 7 UNION ALL
    SELECT 8
) y
ON DUPLICATE KEY UPDATE
    year_type = VALUES(year_type),
    year_status = VALUES(year_status),
    stage_status = VALUES(stage_status),
    report_status = VALUES(report_status),
    summary_effective = VALUES(summary_effective),
    latest_stage_submit_version = VALUES(latest_stage_submit_version),
    latest_report_submit_version = VALUES(latest_report_submit_version),
    updater = 'seed-single',
    update_time = NOW();
