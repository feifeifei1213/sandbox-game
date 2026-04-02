-- 正式比赛环境初始化种子
-- 说明：
-- 1. 本脚本只初始化比赛配置与管理员账号。
-- 2. 不直接初始化玩家组、玩家账号、年份状态，也不提交初始基线。
-- 3. 正式比赛应由管理员首次登录后，在“赛前配置页”中配置小组数量并初始化比赛。

INSERT INTO sg_account (
    id, username, password_hash, role_type, group_id, status, last_login_time,
    creator, create_time, updater, update_time
)
VALUES
    (1, 'admin', CONCAT('{sha256}', SHA2('123456', 256)), 'ADMIN', NULL, 'ENABLED', NULL, 'seed-competition', NOW(), 'seed-competition', NOW())
ON DUPLICATE KEY UPDATE
    password_hash = VALUES(password_hash),
    role_type = VALUES(role_type),
    group_id = VALUES(group_id),
    status = VALUES(status),
    updater = 'seed-competition',
    update_time = NOW();

INSERT INTO sg_game_config (
    id, final_year, current_open_year, rule_version, template_version, initial_baseline_submitted,
    creator, create_time, updater, update_time
)
VALUES
    (1, 8, 0, 'excel-final-2026-03-25', 'excel-page-v1', 0, 'seed-competition', NOW(), 'seed-competition', NOW())
ON DUPLICATE KEY UPDATE
    final_year = VALUES(final_year),
    current_open_year = VALUES(current_open_year),
    rule_version = VALUES(rule_version),
    template_version = VALUES(template_version),
    initial_baseline_submitted = VALUES(initial_baseline_submitted),
    updater = 'seed-competition',
    update_time = NOW();
