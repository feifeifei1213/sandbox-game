-- M1-03 首版核心表 DDL
-- 说明：
-- 1. 本脚本对应 docs/database_design.md 的首版正式表结构。
-- 2. 首版优先保证状态、汇总、日志和核心 JSON 负载可落库。
-- 3. 当前先采用主键、唯一约束和索引约束，不额外引入外键级联，避免初始化和异常修复阶段过度耦合。

CREATE TABLE IF NOT EXISTS sg_group (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    group_no INT NOT NULL COMMENT '组序号',
    group_code VARCHAR(32) NOT NULL COMMENT '组编码',
    group_name VARCHAR(64) NOT NULL COMMENT '组名称',
    business_status VARCHAR(16) NOT NULL DEFAULT 'NORMAL' COMMENT '经营状态：NORMAL/BANKRUPT',
    bankrupt_year_no INT NULL COMMENT '破产年份',
    bankrupt_reason VARCHAR(255) NULL COMMENT '破产原因',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_group_no (group_no),
    UNIQUE KEY uk_group_code (group_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='参赛组主表';

CREATE TABLE IF NOT EXISTS sg_account (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    username VARCHAR(64) NOT NULL COMMENT '登录名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    role_type VARCHAR(16) NOT NULL COMMENT '角色类型：ADMIN/GROUP',
    group_id BIGINT NULL COMMENT '所属组ID，管理员为空',
    status VARCHAR(16) NOT NULL DEFAULT 'ENABLED' COMMENT '状态：ENABLED/DISABLED',
    last_login_time DATETIME NULL COMMENT '最近登录时间',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_username (username),
    KEY idx_role_type (role_type),
    KEY idx_group_id (group_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='系统账号表';

CREATE TABLE IF NOT EXISTS sg_game_config (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    final_year INT NOT NULL COMMENT '最终年份',
    current_open_year INT NOT NULL COMMENT '当前开放年份',
    rule_version VARCHAR(32) NOT NULL COMMENT '规则版本',
    template_version VARCHAR(32) NOT NULL COMMENT '模板版本',
    initial_baseline_submitted TINYINT(1) NOT NULL DEFAULT 0 COMMENT '初始基线是否已提交',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='全局配置表';

CREATE TABLE IF NOT EXISTS sg_group_year_state (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    group_id BIGINT NOT NULL COMMENT '所属组ID',
    year_no INT NOT NULL COMMENT '年份',
    year_type VARCHAR(16) NOT NULL COMMENT '年份类型：INITIAL/DEMO/FORMAL',
    year_status VARCHAR(32) NOT NULL COMMENT '年份主状态',
    stage_status VARCHAR(32) NOT NULL COMMENT '经营阶段状态',
    report_status VARCHAR(32) NOT NULL COMMENT '财报状态',
    summary_effective TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否计入正式汇总',
    latest_stage_submit_version INT NOT NULL DEFAULT 0 COMMENT '当前经营提交版本号',
    latest_report_submit_version INT NOT NULL DEFAULT 0 COMMENT '当前财报提交版本号',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_group_year (group_id, year_no),
    KEY idx_year_status (year_status),
    KEY idx_group_id (group_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='组年度主状态表';

CREATE TABLE IF NOT EXISTS sg_group_operating_draft (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    group_id BIGINT NOT NULL COMMENT '所属组ID',
    year_no INT NOT NULL COMMENT '年份',
    stage_status VARCHAR(32) NOT NULL COMMENT '保存时对应阶段',
    operating_payload_json JSON NOT NULL COMMENT '经营页业务区块数据',
    last_auto_saved_at DATETIME NULL COMMENT '最近自动保存时间',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_group_year_draft (group_id, year_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='经营页当前草稿表';

CREATE TABLE IF NOT EXISTS sg_group_stage_submission (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    group_id BIGINT NOT NULL COMMENT '所属组ID',
    year_no INT NOT NULL COMMENT '年份',
    stage_code VARCHAR(16) NOT NULL COMMENT '阶段编码：Q1/Q2/Q3/Q4/YEAR_END',
    submit_version INT NOT NULL COMMENT '同阶段提交版本号',
    period_end_cash DECIMAL(18,2) NOT NULL DEFAULT 0 COMMENT '提交时系统展示的期末现金',
    operating_payload_snapshot_json JSON NOT NULL COMMENT '经营页快照',
    state_before_json JSON NOT NULL COMMENT '提交前状态摘要',
    state_after_json JSON NOT NULL COMMENT '提交后状态摘要',
    submitter_id BIGINT NOT NULL COMMENT '提交人ID',
    submit_time DATETIME NOT NULL COMMENT '提交时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_group_year_stage_version (group_id, year_no, stage_code, submit_version),
    KEY idx_group_year (group_id, year_no),
    KEY idx_submit_time (submit_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='经营阶段提交流水表';

CREATE TABLE IF NOT EXISTS sg_group_report (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    group_id BIGINT NOT NULL COMMENT '所属组ID',
    year_no INT NOT NULL COMMENT '年份',
    report_manual_payload_json JSON NOT NULL COMMENT '财报手工项',
    report_computed_payload_json JSON NOT NULL COMMENT '财报自动计算结果',
    balance_check_passed TINYINT(1) NOT NULL DEFAULT 0 COMMENT '平衡校验是否通过',
    last_auto_saved_at DATETIME NULL COMMENT '最近自动保存时间',
    submitted_at DATETIME NULL COMMENT '最近正式提交时间',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_group_year_report (group_id, year_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='财报当前内容表';

CREATE TABLE IF NOT EXISTS sg_group_report_submission (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    group_id BIGINT NOT NULL COMMENT '所属组ID',
    year_no INT NOT NULL COMMENT '年份',
    submit_version INT NOT NULL COMMENT '财报提交版本号',
    report_manual_snapshot_json JSON NOT NULL COMMENT '财报手工项快照',
    report_computed_snapshot_json JSON NOT NULL COMMENT '财报自动项快照',
    balance_check_passed TINYINT(1) NOT NULL DEFAULT 0 COMMENT '平衡校验结果',
    state_before_json JSON NOT NULL COMMENT '提交前状态',
    state_after_json JSON NOT NULL COMMENT '提交后状态',
    submitter_id BIGINT NOT NULL COMMENT '提交人ID',
    submit_time DATETIME NOT NULL COMMENT '提交时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_group_year_report_version (group_id, year_no, submit_version)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='财报提交流水表';

CREATE TABLE IF NOT EXISTS sg_group_summary_snapshot (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    group_id BIGINT NOT NULL COMMENT '所属组ID',
    year_no INT NOT NULL COMMENT '年份',
    revenue DECIMAL(18,2) NOT NULL DEFAULT 0 COMMENT '收入',
    profit DECIMAL(18,2) NOT NULL DEFAULT 0 COMMENT '利润',
    equity DECIMAL(18,2) NOT NULL DEFAULT 0 COMMENT '权益',
    business_status VARCHAR(16) NOT NULL DEFAULT 'NORMAL' COMMENT '经营状态',
    ranking_value DECIMAL(18,2) NOT NULL DEFAULT 0 COMMENT '排名依据值',
    summary_effective TINYINT(1) NOT NULL DEFAULT 0 COMMENT '当前快照是否生效',
    source_report_submit_version INT NOT NULL DEFAULT 0 COMMENT '来源财报提交版本',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_group_year_summary (group_id, year_no),
    KEY idx_year_no (year_no),
    KEY idx_summary_effective (summary_effective)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='正式年份汇总快照表';

CREATE TABLE IF NOT EXISTS sg_initial_baseline (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    group_id BIGINT NOT NULL COMMENT '所属组ID',
    baseline_payload_json JSON NOT NULL COMMENT '初始基线数据',
    submitted TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已提交',
    submitter_id BIGINT NULL COMMENT '提交管理员ID',
    submitted_at DATETIME NULL COMMENT '提交时间',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_group_baseline (group_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='初始年财报基线表';

CREATE TABLE IF NOT EXISTS sg_admin_unlock_log (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    group_id BIGINT NOT NULL COMMENT '所属组ID',
    year_no INT NOT NULL COMMENT '年份',
    reason VARCHAR(500) NOT NULL COMMENT '解锁原因',
    state_before_json JSON NOT NULL COMMENT '解锁前状态',
    state_after_json JSON NOT NULL COMMENT '解锁后状态',
    operator_id BIGINT NOT NULL COMMENT '操作管理员ID',
    operator_name VARCHAR(64) NOT NULL COMMENT '操作管理员名称',
    operate_time DATETIME NOT NULL COMMENT '操作时间',
    PRIMARY KEY (id),
    KEY idx_group_year (group_id, year_no),
    KEY idx_operate_time (operate_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='异常解锁日志表';

CREATE TABLE IF NOT EXISTS sg_admin_action_log (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    action_code VARCHAR(64) NOT NULL COMMENT '管理动作编码',
    target_group_id BIGINT NULL COMMENT '目标组ID',
    target_year_no INT NULL COMMENT '目标年份',
    action_payload_json JSON NOT NULL COMMENT '动作参数',
    state_before_json JSON NOT NULL COMMENT '操作前状态',
    state_after_json JSON NOT NULL COMMENT '操作后状态',
    operator_id BIGINT NOT NULL COMMENT '操作管理员ID',
    operator_name VARCHAR(64) NOT NULL COMMENT '操作管理员名称',
    operate_time DATETIME NOT NULL COMMENT '操作时间',
    PRIMARY KEY (id),
    KEY idx_action_code (action_code),
    KEY idx_target_group_year (target_group_id, target_year_no),
    KEY idx_operate_time (operate_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='管理员关键动作日志表';
