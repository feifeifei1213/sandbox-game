CREATE TABLE IF NOT EXISTS sg_notice (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    target_scope VARCHAR(16) NOT NULL COMMENT '目标范围：ALL/GROUP',
    target_group_id BIGINT NULL COMMENT '目标小组 ID，ALL 时为空',
    content VARCHAR(1000) NOT NULL COMMENT '通知内容',
    pinned TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否置顶',
    published_at DATETIME NOT NULL COMMENT '发布时间',
    operator_id BIGINT NOT NULL COMMENT '操作管理员 ID',
    operator_name VARCHAR(64) NOT NULL COMMENT '操作管理员名称',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_target_scope_group (target_scope, target_group_id),
    KEY idx_published_at (published_at),
    KEY idx_pinned (pinned)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='普通通知表';

CREATE TABLE IF NOT EXISTS sg_group_adjustment (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    group_id BIGINT NOT NULL COMMENT '目标小组 ID',
    year_no INT NOT NULL COMMENT '目标年份',
    stage_code VARCHAR(16) NOT NULL COMMENT '目标季度：Q1/Q2/Q3/Q4',
    adjustment_type VARCHAR(16) NOT NULL COMMENT '类型：REWARD/PENALTY',
    amount DECIMAL(18,2) NOT NULL DEFAULT 0 COMMENT '金额',
    reason VARCHAR(500) NOT NULL COMMENT '奖惩原因',
    published_at DATETIME NOT NULL COMMENT '发布时间',
    operator_id BIGINT NOT NULL COMMENT '操作管理员 ID',
    operator_name VARCHAR(64) NOT NULL COMMENT '操作管理员名称',
    creator VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '创建人',
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '更新人',
    update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_group_year_stage (group_id, year_no, stage_code),
    KEY idx_published_at (published_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='按组按年按季奖惩记录表';
