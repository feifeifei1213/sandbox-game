-- M1-04 基础种子数据
-- 说明：
-- 1. 本脚本用于本地开发和联调初始化，不直接替代正式活动导入流程。
-- 2. 默认口径：初始基线已提交、当前开放年份为 0 年，便于直接开始体验主链。
-- 3. 账号密码当前以 SHA2-256 形式存入 password_hash，明文默认口径统一为 123456，后续若接入正式认证方案可再替换算法。

INSERT INTO sg_group (
    id, group_no, group_code, group_name, business_status, bankrupt_year_no, bankrupt_reason,
    creator, create_time, updater, update_time
)
VALUES
    (1, 1, 'GROUP_01', '第一组', 'NORMAL', NULL, NULL, 'seed', NOW(), 'seed', NOW()),
    (2, 2, 'GROUP_02', '第二组', 'NORMAL', NULL, NULL, 'seed', NOW(), 'seed', NOW()),
    (3, 3, 'GROUP_03', '第三组', 'NORMAL', NULL, NULL, 'seed', NOW(), 'seed', NOW()),
    (4, 4, 'GROUP_04', '第四组', 'NORMAL', NULL, NULL, 'seed', NOW(), 'seed', NOW()),
    (5, 5, 'GROUP_05', '第五组', 'NORMAL', NULL, NULL, 'seed', NOW(), 'seed', NOW()),
    (6, 6, 'GROUP_06', '第六组', 'NORMAL', NULL, NULL, 'seed', NOW(), 'seed', NOW()),
    (7, 7, 'GROUP_07', '第七组', 'NORMAL', NULL, NULL, 'seed', NOW(), 'seed', NOW()),
    (8, 8, 'GROUP_08', '第八组', 'NORMAL', NULL, NULL, 'seed', NOW(), 'seed', NOW()),
    (9, 9, 'GROUP_09', '第九组', 'NORMAL', NULL, NULL, 'seed', NOW(), 'seed', NOW()),
    (10, 10, 'GROUP_10', '第十组', 'NORMAL', NULL, NULL, 'seed', NOW(), 'seed', NOW())
ON DUPLICATE KEY UPDATE
    group_name = VALUES(group_name),
    business_status = VALUES(business_status),
    bankrupt_year_no = VALUES(bankrupt_year_no),
    bankrupt_reason = VALUES(bankrupt_reason),
    updater = 'seed',
    update_time = NOW();

INSERT INTO sg_account (
    id, username, password_hash, role_type, group_id, status, last_login_time,
    creator, create_time, updater, update_time
)
VALUES
    (1, 'admin', CONCAT('{sha256}', SHA2('123456', 256)), 'ADMIN', NULL, 'ENABLED', NULL, 'seed', NOW(), 'seed', NOW()),
    (101, 'group01', CONCAT('{sha256}', SHA2('123456', 256)), 'GROUP', 1, 'ENABLED', NULL, 'seed', NOW(), 'seed', NOW()),
    (102, 'group02', CONCAT('{sha256}', SHA2('123456', 256)), 'GROUP', 2, 'ENABLED', NULL, 'seed', NOW(), 'seed', NOW()),
    (103, 'group03', CONCAT('{sha256}', SHA2('123456', 256)), 'GROUP', 3, 'ENABLED', NULL, 'seed', NOW(), 'seed', NOW()),
    (104, 'group04', CONCAT('{sha256}', SHA2('123456', 256)), 'GROUP', 4, 'ENABLED', NULL, 'seed', NOW(), 'seed', NOW()),
    (105, 'group05', CONCAT('{sha256}', SHA2('123456', 256)), 'GROUP', 5, 'ENABLED', NULL, 'seed', NOW(), 'seed', NOW()),
    (106, 'group06', CONCAT('{sha256}', SHA2('123456', 256)), 'GROUP', 6, 'ENABLED', NULL, 'seed', NOW(), 'seed', NOW()),
    (107, 'group07', CONCAT('{sha256}', SHA2('123456', 256)), 'GROUP', 7, 'ENABLED', NULL, 'seed', NOW(), 'seed', NOW()),
    (108, 'group08', CONCAT('{sha256}', SHA2('123456', 256)), 'GROUP', 8, 'ENABLED', NULL, 'seed', NOW(), 'seed', NOW()),
    (109, 'group09', CONCAT('{sha256}', SHA2('123456', 256)), 'GROUP', 9, 'ENABLED', NULL, 'seed', NOW(), 'seed', NOW()),
    (110, 'group10', CONCAT('{sha256}', SHA2('123456', 256)), 'GROUP', 10, 'ENABLED', NULL, 'seed', NOW(), 'seed', NOW())
ON DUPLICATE KEY UPDATE
    password_hash = VALUES(password_hash),
    role_type = VALUES(role_type),
    group_id = VALUES(group_id),
    status = VALUES(status),
    updater = 'seed',
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
        1, 'seed', NOW(), 'seed', NOW()
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
    updater = 'seed',
    update_time = NOW();

INSERT INTO sg_initial_baseline (
    group_id, baseline_payload_json, submitted, submitter_id, submitted_at,
    creator, create_time, updater, update_time
)
SELECT
    g.id,
    JSON_OBJECT(
        'baselineSalesRevenue', 32,
        'baselineDirectCost', 15,
        'baselineComprehensiveCost', 13,
        'baselineDepreciation', 1,
        'baselineFinanceIncomeExpense', 2,
        'baselineExtraIncomeExpense', 2,
        'baselineIncomeTax', 1,
        'baselineFactoryAsset', 40,
        'baselineLineResidual', 3,
        'baselineDepreciableAsset', 0,
        'baselineCash', 36,
        'baselineReceivable', 0,
        'baselineWorkInProgress', 6,
        'baselineFinishedGoods', 4,
        'baselineRawMaterials', 1,
        'baselineShortTermLoan', 20,
        'baselineLongTermLoan', 0,
        'baselineShareCapital', 50,
        'baselineRetainedEarnings', 17
    ),
    1,
    1,
    NOW(),
    'seed',
    NOW(),
    'seed',
    NOW()
FROM sg_group g
ON DUPLICATE KEY UPDATE
    baseline_payload_json = VALUES(baseline_payload_json),
    submitted = VALUES(submitted),
    submitter_id = VALUES(submitter_id),
    submitted_at = VALUES(submitted_at),
    updater = 'seed',
    update_time = NOW();

INSERT INTO sg_group_year_state (
    group_id, year_no, year_type, year_status, stage_status, report_status,
    summary_effective, latest_stage_submit_version, latest_report_submit_version,
    creator, create_time, updater, update_time
)
SELECT
    g.id,
    y.year_no,
    CASE WHEN y.year_no = 0 THEN 'DEMO' ELSE 'FORMAL' END AS year_type,
    CASE WHEN y.year_no = 0 THEN 'OPERATING' ELSE 'LOCKED' END AS year_status,
    'Q1_OPEN' AS stage_status,
    'REPORT_LOCKED' AS report_status,
    0 AS summary_effective,
    0 AS latest_stage_submit_version,
    0 AS latest_report_submit_version,
    'seed' AS creator,
    NOW() AS create_time,
    'seed' AS updater,
    NOW() AS update_time
FROM sg_group g
CROSS JOIN (
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
WHERE 1 = 1
ON DUPLICATE KEY UPDATE
    year_type = VALUES(year_type),
    year_status = VALUES(year_status),
    stage_status = VALUES(stage_status),
    report_status = VALUES(report_status),
    summary_effective = VALUES(summary_effective),
    latest_stage_submit_version = VALUES(latest_stage_submit_version),
    latest_report_submit_version = VALUES(latest_report_submit_version),
    updater = 'seed',
    update_time = NOW();
