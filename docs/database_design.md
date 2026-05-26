# 沙盘经营系统数据库设计文档（正式版）

> 更新日期：2026-05-20
> 适用方式：基于《正式需求文档（首版）》《最小状态机 v0.1》《技术选型细化文档（Go 方向） v0.1》，定义首版正式数据库设计方向，作为后续 MySQL 建表、迁移脚本、Repository 实现和状态机落库的统一依据。  
> 文档定位：本文件定义表清单、核心字段、关系、约束、索引与变更规则，不替代最终 SQL 脚本。  
> 说明：当前显示名称已按 `1组 最终版.xlsx` 冻结；数据库字段命名应坚持业务语义，不应直接跟随 Excel 中文标题逐格命名。

## 1. 数据库基线

- 数据库：`MySQL 8.x`
- 字符集：`utf8mb4`
- 排序规则：推荐 `utf8mb4_general_ci` 或团队统一规则
- 时区：统一使用 `Asia/Shanghai` 或存 UTC 后在应用层转换
- 主脚本建议：`sql/mysql/sandbox_game.sql`
- 升级脚本建议：`sql/mysql/sandbox-game-upgrade-*.sql`

---

## 2. 核心设计约定

### 2.1 主键与通用字段

- 主键：`BIGINT AUTO_INCREMENT`
- 时间字段：`DATETIME`
- 金额字段：优先 `DECIMAL(18,2)`
- JSON 字段：使用 `JSON`

### 2.2 审计字段约定

业务主表统一建议保留：

- `creator`
- `create_time`
- `updater`
- `update_time`

日志表除以上字段外，应增加业务动作专用字段：

- `operator_id`
- `operator_name`
- `operate_time`

### 2.3 删除策略

首版原则：

- 历史提交、汇总快照、异常解锁日志、管理员动作日志不允许物理删除
- 核心业务数据以“状态失效 / 快照保留”为主，不依赖删除回滚
- 若后续需要逻辑删除，仅用于配置型 / 辅助型表，不建议用于提交流水表

### 2.4 多租户策略

首版默认不主动引入 `tenant_id`，原因：

- 当前项目为单活动、单组织使用
- 用户规模小
- 业务复杂度主要在规则而非租户隔离

若后续接入公司统一平台，可再补充多租户字段。

---

## 3. 建模原则

### 3.1 不按 Excel 单元格建模

数据库不应直接以“单元格坐标”作为核心数据模型。

禁止：

- 用 `A1/B2/C3` 一类坐标做主业务字段
- 以单元格为粒度直接设计所有存储结构

### 3.2 采用“强状态 + 弱区块”的混合建模

首版推荐：

- 状态、关键结果、日志：强结构化关系型字段
- 经营页 / 财报页手工输入区：区块化 JSON
- 汇总结果：独立快照表

### 3.3 为什么这样建模

原因如下：

- 状态机、权限、汇总查询依赖强结构化字段
- 即使后续仍有少量页面文案细修，区块 JSON 仍比逐格字段更抗变化
- 关键结果不应每次都从原始区块全量重算

---

## 4. 表清单（首版正式版）

### 4.1 账户与组

#### 4.1.1 `sg_account`

用途：系统账号表。

说明：

- 若后续接入统一认证，本表可退化为本地账号映射表。
- 首版若项目独立运行，建议保留本表。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `username` | VARCHAR(64) | 登录名，唯一 |
| `password_hash` | VARCHAR(255) | 密码哈希 |
| `role_type` | VARCHAR(16) | `ADMIN` / `GROUP` |
| `group_id` | BIGINT NULL | 玩家账号对应组；管理员为空 |
| `status` | VARCHAR(16) | `ENABLED` / `DISABLED` |
| `last_login_time` | DATETIME NULL | 最近登录时间 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_username(username)`
- `idx_role_type(role_type)`
- `idx_group_id(group_id)`

#### 4.1.2 `sg_group`

用途：参赛组主表。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `group_no` | INT | 第几组 |
| `group_code` | VARCHAR(32) | 组编码，唯一 |
| `group_name` | VARCHAR(64) | 组名称 |
| `business_status` | VARCHAR(16) | `NORMAL` / `BANKRUPT` |
| `bankrupt_year_no` | INT NULL | 破产年份 |
| `bankrupt_reason` | VARCHAR(255) NULL | 破产原因 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_group_code(group_code)`
- `uk_group_no(group_no)`

说明：

- 若异常解锁回收了导致破产的那一年结果，应同步把 `business_status` 暂时恢复为 `NORMAL`，并清理对应破产标记，待重新提交后再重算。

---

### 4.2 游戏配置与规则版本

#### 4.2.1 `sg_game_config`

用途：全局配置表。

建议按单行配置或 `config_key/config_value` 形式二选一，首版更推荐单行表。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `final_year` | INT | 最终年份 |
| `current_open_year` | INT | 当前开放年份 |
| `edition_code` | VARCHAR(64) | 沙盘版本包编码，例如 `VIP_SERVICE_V1`、`PRODUCTION_V1` |
| `edition_name` | VARCHAR(64) | 沙盘版本包显示名，例如 `贵宾服务版 V1`、`生产制造版 V1` |
| `rule_version` | VARCHAR(32) | 规则版本 |
| `template_version` | VARCHAR(32) | 模板版本 |
| `operating_template_version` | VARCHAR(64) | 经营页字段模板版本 |
| `report_template_version` | VARCHAR(64) | 财报页字段模板版本 |
| `order_template_version` | VARCHAR(64) | 订单字段模板版本 |
| `process_rule_version` | VARCHAR(64) | 流程规则版本 |
| `initial_baseline_submitted` | TINYINT(1) | 初始基线是否已提交 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

说明：

- `edition_code` 用于锁定本场比赛采用的内置版本包，当前已支持 `VIP_SERVICE_V1`，下一步新增 `PRODUCTION_V1`。
- `rule_version` 和 `template_version` 用于追踪 Excel 规则基线与页面模板版本。`operating_template_version / report_template_version / order_template_version` 用于区分经营页、财报页和订单页字段模板。
- 当前生产版和贵宾版公式、流程一致，可共用 `COMMON_FORMULA_V1` 与 `COMMON_PROCESS_V1`。
- `PRODUCTION_V1` 首轮应写入 `PRODUCTION_OPERATING_TEMPLATE_V1`、`PRODUCTION_REPORT_TEMPLATE_V1`，订单模板暂写入 `VIP_ORDER_TEMPLATE_V1`；待生产版订单字段确认后再补充新的订单模板版本。
- 比赛初始化后版本包不允许切换；如需切换，应重新初始化新的比赛环境或比赛库。
- `final_year` 不能小于 `current_open_year`。

---

### 4.3 年份状态与经营数据

#### 4.3.1 `sg_group_year_state`

用途：`组 + 年` 的主状态表，是状态机落库核心。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `group_id` | BIGINT | 所属组 |
| `year_no` | INT | 年份 |
| `year_type` | VARCHAR(16) | `INITIAL` / `DEMO` / `FORMAL` |
| `year_status` | VARCHAR(32) | `LOCKED` / `OPERATING` / `REPORT_PENDING` / `REPORTING` / `COMPLETED` |
| `stage_status` | VARCHAR(32) | `Q1_OPEN` / `Q2_OPEN` / `Q3_OPEN` / `Q4_OPEN` / `YEAR_END_OPEN` |
| `report_status` | VARCHAR(32) | `REPORT_LOCKED` / `REPORT_OPEN` / `REPORT_SUBMITTED` |
| `summary_effective` | TINYINT(1) | 是否计入正式汇总 |
| `latest_stage_submit_version` | INT | 当前经营提交版本号 |
| `latest_report_submit_version` | INT | 当前财报提交版本号 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_group_year(group_id, year_no)`
- `idx_year_status(year_status)`
- `idx_group_id(group_id)`

说明：

- `sg_group_year_state` 不保存所有明细，但必须保存当前有效状态。
- 前端锁定区、汇总生效、异常解锁回收，都依赖本表。

#### 4.3.2 `sg_group_operating_draft`

用途：经营页当前草稿。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `group_id` | BIGINT | 所属组 |
| `year_no` | INT | 年份 |
| `stage_status` | VARCHAR(32) | 保存时对应阶段 |
| `operating_payload_json` | JSON | 经营页业务区块数据 |
| `last_auto_saved_at` | DATETIME | 最近自动保存时间 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_group_year_draft(group_id, year_no)`

说明：

- 草稿表只保留当前最新草稿。
- 草稿不等于正式结果。

#### 4.3.3 `sg_group_stage_submission`

用途：经营阶段正式提交记录。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `group_id` | BIGINT | 所属组 |
| `year_no` | INT | 年份 |
| `stage_code` | VARCHAR(16) | `Q1` / `Q2` / `Q3` / `Q4` / `YEAR_END` |
| `submit_version` | INT | 同阶段提交版本号 |
| `board_cash` | DECIMAL(18,2) | 玩家输入沙盘现金 |
| `system_cash` | DECIMAL(18,2) | 系统计算现金 |
| `cash_diff` | DECIMAL(18,2) | 差额 |
| `operating_payload_snapshot_json` | JSON | 提交时经营页快照 |
| `state_before_json` | JSON | 提交前状态摘要 |
| `state_after_json` | JSON | 提交后状态摘要 |
| `submitter_id` | BIGINT | 提交人 |
| `submit_time` | DATETIME | 提交时间 |

关键约束：

- `uk_group_year_stage_version(group_id, year_no, stage_code, submit_version)`
- `idx_group_year(group_id, year_no)`
- `idx_submit_time(submit_time)`

说明：

- 因为存在异常解锁后重提，同一阶段可能有多个版本。
- 当前有效版本由 `sg_group_year_state.latest_stage_submit_version` 结合阶段判断。

---

### 4.4 财报与汇总

#### 4.4.1 `sg_group_report`

用途：财报当前最新内容表。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `group_id` | BIGINT | 所属组 |
| `year_no` | INT | 年份 |
| `report_manual_payload_json` | JSON | 财报手工项 |
| `report_computed_payload_json` | JSON | 财报自动计算结果 |
| `balance_check_passed` | TINYINT(1) | 平衡校验是否通过 |
| `last_auto_saved_at` | DATETIME NULL | 最近草稿保存时间 |
| `submitted_at` | DATETIME NULL | 最近正式提交时间 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_group_year_report(group_id, year_no)`

#### 4.4.2 `sg_group_report_submission`

用途：财报正式提交历史。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `group_id` | BIGINT | 所属组 |
| `year_no` | INT | 年份 |
| `submit_version` | INT | 财报提交版本号 |
| `report_manual_snapshot_json` | JSON | 手工项快照 |
| `report_computed_snapshot_json` | JSON | 自动项快照 |
| `balance_check_passed` | TINYINT(1) | 平衡校验结果 |
| `state_before_json` | JSON | 提交前状态 |
| `state_after_json` | JSON | 提交后状态 |
| `submitter_id` | BIGINT | 提交人 |
| `submit_time` | DATETIME | 提交时间 |

关键约束：

- `uk_group_year_report_version(group_id, year_no, submit_version)`

#### 4.4.3 `sg_group_summary_snapshot`

用途：正式年份汇总快照表。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `group_id` | BIGINT | 所属组 |
| `year_no` | INT | 年份 |
| `revenue` | DECIMAL(18,2) | 收入 |
| `profit` | DECIMAL(18,2) | 利润 |
| `equity` | DECIMAL(18,2) | 权益 |
| `business_status` | VARCHAR(16) | `NORMAL` / `BANKRUPT` |
| `ranking_value` | DECIMAL(18,2) | 排名依据值，首版即 `equity` |
| `summary_effective` | TINYINT(1) | 当前快照是否生效 |
| `source_report_submit_version` | INT | 来源财报提交版本 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_group_year_summary(group_id, year_no)`
- `idx_year_no(year_no)`
- `idx_summary_effective(summary_effective)`

说明：

- 被异常解锁后，本表对应年份记录应改为 `summary_effective=0` 或更新为失效态。
- 汇总页查询正式口径时，只取 `summary_effective=1`。

---

### 4.5 通知与奖惩

#### 4.5.1 `sg_notice`

用途：管理员普通通知表。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `target_scope` | VARCHAR(16) | `ALL` / `GROUP` |
| `target_group_id` | BIGINT NULL | 目标组，发全体时为空 |
| `content` | VARCHAR(1000) | 通知内容 |
| `pinned` | TINYINT(1) | 是否置顶 |
| `published_at` | DATETIME | 发布时间 |
| `operator_id` | BIGINT | 操作管理员 |
| `operator_name` | VARCHAR(64) | 操作管理员名称 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `idx_target_scope_group(target_scope, target_group_id)`
- `idx_published_at(published_at)`
- `idx_pinned(pinned)`

说明：

- 普通通知只参与展示，不参与经营、财报、汇总计算。
- 同一组读取通知时，应同时看到 `ALL` 与自身 `GROUP` 通知。

#### 4.5.2 `sg_group_adjustment`

用途：按 `组 + 年 + 季` 存储奖励 / 罚款记录。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `group_id` | BIGINT | 目标组 |
| `year_no` | INT | 目标年份 |
| `stage_code` | VARCHAR(16) | `Q1 / Q2 / Q3 / Q4` |
| `adjustment_type` | VARCHAR(16) | `REWARD` / `PENALTY` |
| `amount` | DECIMAL(18,2) | 金额 |
| `reason` | VARCHAR(500) | 奖惩原因 |
| `published_at` | DATETIME | 发布时间 |
| `operator_id` | BIGINT | 操作管理员 |
| `operator_name` | VARCHAR(64) | 操作管理员名称 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `idx_group_year_stage(group_id, year_no, stage_code)`
- `idx_published_at(published_at)`

说明：

- 首版允许同一季度存在多条奖惩记录，查询层按季度聚合展示。
- 奖惩只允许作用于尚未锁定的季度；锁定后如需修正，应走异常解锁。
- 经营页、财报页和汇总口径只读取当前有效年份状态对应的奖惩聚合结果。

---

### 4.6 管理动作与审计日志

#### 4.6.1 `sg_initial_baseline`

用途：管理员提交的初始基线数据。

说明：

- 页面录入语义为“`1 份共享初始基线模板`”。
- 后端提交时可将同一份模板扇出复制到全部小组，因此本表仍按 `group_id` 存储，便于后续按组追溯。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `group_id` | BIGINT | 所属组 |
| `baseline_payload_json` | JSON | 初始基线数据 |
| `submitted` | TINYINT(1) | 是否已提交 |
| `submitter_id` | BIGINT | 提交管理员 |
| `submitted_at` | DATETIME | 提交时间 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_group_baseline(group_id)`

#### 4.6.2 `sg_admin_unlock_log`

用途：异常解锁日志。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `group_id` | BIGINT | 所属组 |
| `year_no` | INT | 年份 |
| `reason` | VARCHAR(500) | 解锁原因 |
| `state_before_json` | JSON | 解锁前状态 |
| `state_after_json` | JSON | 解锁后状态 |
| `operator_id` | BIGINT | 操作管理员 |
| `operator_name` | VARCHAR(64) | 操作管理员名称 |
| `operate_time` | DATETIME | 操作时间 |

#### 4.6.3 `sg_admin_action_log`

用途：管理员关键动作日志。`action_code` 建议至少固定为：`UPDATE_FINAL_YEAR`、`OPEN_NEXT_YEAR`、`SUBMIT_INITIAL_BASELINE`。

动作至少覆盖：

- 更新最终年份
- 开放下一年
- 提交初始基线
- 管理员代跳过当前标段选单小组

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `action_code` | VARCHAR(64) | 管理动作编码 |
| `target_group_id` | BIGINT NULL | 目标组 |
| `target_year_no` | INT NULL | 目标年份 |
| `action_payload_json` | JSON | 动作参数 |
| `state_before_json` | JSON | 操作前状态 |
| `state_after_json` | JSON | 操作后状态 |
| `operator_id` | BIGINT | 操作管理员 |
| `operator_name` | VARCHAR(64) | 操作管理员名称 |
| `operate_time` | DATETIME | 操作时间 |

---

### 4.7 年度订单与市场竞标

#### 4.7.0 `sg_order_forecast_control`

用途：存储多年订单数量控制台，是市场预测和年度正式订单池的共同订单数量来源。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `year_no` | INT | 预测年份，固定 `1~8` |
| `market_code` | VARCHAR(32) | `LOCAL / REGIONAL / NATIONAL / GLOBAL` |
| `order_type` | VARCHAR(32) | `AGENCY_INSPECTION / TWO_CABIN_VIP / BUSINESS_VIP / MEMBER_CUSTOM` |
| `order_count` | INT | 订单卡片数量，范围 `0~15` |
| `forecast_stage_code` | VARCHAR(32) | `YEAR_1_3 / YEAR_4_5 / YEAR_6_8` |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_order_forecast_control(year_no, market_code, order_type)`
- `idx_order_forecast_stage(forecast_stage_code, market_code)`

说明：

- 该表固定覆盖 `1年~8年 × 四个市场 × 四类产品`，不随管理员配置的最终年份裁剪。
- 年度订单池生成时，只读取该年且管理员手动开启市场对应的控制台数量。
- 控制台有数量不代表市场自动开启；市场开启仍以 `sg_order_market_config` 为准。

#### 4.7.0A `sg_order_market_forecast`

用途：存储市场预测展示快照和说明文字，供玩家端只读查看。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `forecast_stage_code` | VARCHAR(32) | `YEAR_1_3 / YEAR_4_5 / YEAR_6_8` |
| `market_code` | VARCHAR(32) | 市场 |
| `forecast_data_json` | JSON | 该阶段内各年份、四类产品的预测金额/订单量数据 |
| `narrative` | TEXT | 市场预测说明文字 |
| `formula_version` | VARCHAR(64) | 预测生成公式版本 |
| `random_seed` | VARCHAR(64) | 预测快照随机种子 |
| `control_snapshot_json` | JSON | 生成该预测时的多年控制台快照 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

说明：

- 玩家端市场预测固定按三段展示，不因最终年份变化隐藏 `6~8年`。
- 市场预测只是趋势展示，不代表当年市场自动开启，也不直接生成玩家可选订单。

#### 4.7.1 `sg_order_generation_batch`

用途：记录系统按订单推算 Excel 公式链生成订单池的预览/正式批次。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `year_no` | INT | 年份，正式年份从 `1` 开始 |
| `batch_status` | VARCHAR(32) | `PREVIEW / CONFIRMED / VOID` |
| `formula_version` | VARCHAR(64) | 订单生成公式版本 |
| `random_seed` | VARCHAR(64) | 随机种子，用于复盘 |
| `control_snapshot_json` | JSON | 订单数量控制台快照 |
| `forecast_snapshot_json` | JSON | 市场预测快照摘要 |
| `parameter_snapshot_json` | JSON | 均价、波动系数、最小/最大数量、账期范围等公式参数快照 |
| `confirmed_at` | DATETIME NULL | 确认时间 |
| `operator_id/operator_name/operate_time` | - | 生成/确认操作人信息 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

说明：

- 预览批次可被重新生成覆盖或置为 `VOID`。
- 管理员确认后，当前预览批次转为 `CONFIRMED`，订单池固定落库。
- 生成批次的订单数量来源应为多年订单数量控制台中该年、已开启市场的数量快照。
- “上传 Excel”不再作为订单池主链路；若后续保留导入能力，应另建“导入控制台参数/模板”的辅助记录，不替代本生成批次。

#### 4.7.2 `sg_order_generation_config`

用途：按 `年份 + 市场 + 订单类型` 存储当年标段释放顺序和生成批次关联；订单数量从多年订单数量控制台带入。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `year_no` | INT | 年份，正式年份从 `1` 开始 |
| `market_code` | VARCHAR(32) | `LOCAL / REGIONAL / NATIONAL / GLOBAL` |
| `order_type` | VARCHAR(32) | `AGENCY_INSPECTION / TWO_CABIN_VIP / BUSINESS_VIP / MEMBER_CUSTOM` |
| `order_count` | INT | 从多年订单数量控制台带入的该类订单生成数量快照 |
| `release_sequence_no` | INT | 标段释放顺序；释放单元为 `市场 + 订单类型` |
| `generation_batch_id` | BIGINT NULL | 当前确认订单池批次 |
| `config_status` | VARCHAR(32) | `DRAFT / CONFIRMED / LOCKED` |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_order_generation_config(year_no, market_code, order_type)`
- `uk_order_release_sequence(year_no, release_sequence_no)`
- `idx_year_market(year_no, market_code)`

说明：

- 多年订单数量控制台是订单数量主来源；本表中的 `order_count` 是当年生成/确认时的快照，不应作为另一套独立数量来源。
- `order_count` 范围为 `0 ~ 15`；`0` 表示不生成订单且该标段不进入开标。
- `release_sequence_no` 必须在开标前配置完成；释放第一个标段后不允许调整。
- 系统可基于 `order_count`、实际小组数、市场投入资格生成风险提示，但不做强制保底或多轮分配。
- 均价、波动系数、最小/最大订单数量、账期范围等复杂参数暂不进入首版配置页面，但保存到生成批次参数快照。

#### 4.7.2A `sg_order_market_config`

用途：按 `年份 + 市场` 存储当年市场开启状态。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `year_no` | INT | 年份，正式年份从 `1` 开始 |
| `market_code` | VARCHAR(32) | `LOCAL / REGIONAL / NATIONAL / GLOBAL` |
| `market_enabled` | TINYINT(1) | 是否开启该市场 |
| `config_status` | VARCHAR(32) | `DRAFT / LOCKED` |
| `locked_batch_id` | BIGINT NULL | 锁定时对应确认订单池批次 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_order_market_config(year_no, market_code)`
- `idx_order_market_enabled(year_no, market_enabled)`

说明：

- 本地市场默认 `market_enabled=1`。
- 区域市场、全国市场、全球市场默认 `market_enabled=0`。
- 订单池确认前可修改市场开启状态；确认后该年市场开启状态锁定。
- 市场开启完全以管理员当年手动配置为准，不因多年订单数量控制台中存在订单数量而自动开启。
- 未开启市场下四个标段不生成订单池、不占用有效释放顺序、不进入选单。
- 玩家仍需提交未开启市场对应的 4 项投入，且必须为 `0`；服务端负责拦截非 `0`。

#### 4.7.3 `sg_order_pool`

用途：固定订单池。管理员生成后，玩家从中选择订单。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `year_no` | INT | 年份 |
| `market_code` | VARCHAR(32) | 市场 |
| `order_type` | VARCHAR(32) | 订单类型 |
| `segment_code` | VARCHAR(64) | 标段编码 |
| `card_sequence_no` | INT | 订单卡片序号，`1 ~ 15` |
| `business_order_no` | VARCHAR(64) | 业务订单编号/页面展示编号，例如 `CARD-01` |
| `order_amount` | DECIMAL(18,2) | 订单金额，业务口径为整数金额 |
| `order_quantity` | DECIMAL(18,2) | 数量，订单生成口径为整数 |
| `unit_price` | DECIMAL(18,2) | 单价，允许小数 |
| `account_term` | VARCHAR(64) | 账期 |
| `pool_status` | VARCHAR(32) | `AVAILABLE / SELECTED / VOID` |
| `selected_group_id` | BIGINT NULL | 选中小组 |
| `selected_at` | DATETIME NULL | 选中时间 |
| `generation_batch_id` | BIGINT | 订单生成批次 |
| `source_row_key` | VARCHAR(128) NULL | Excel 公式链对应卡片标识 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `idx_order_pool_year_market(year_no, market_code)`
- `idx_order_pool_status(pool_status)`
- `idx_order_pool_selected_group(selected_group_id)`

说明：

- 订单池生成后保存固定值，不在运行时持续依赖 Excel 随机公式。
- 生成时由后端自动保证 `order_amount = order_quantity × unit_price`，且 `order_amount` 为整数；管理员不维护逐单金额或单价。
- 订单池来自系统按 Excel 公式链生成后的批次，不来自上传固定订单明细。
- 预览阶段可覆盖；确认后不允许覆盖或重新生成。
- 管理端和玩家端订单列表应优先展示 `business_order_no` 或由 `card_sequence_no` 格式化出的 `CARD-01`，不应把数据库自增 `id` 作为业务订单编号展示给管理员。

#### 4.7.4 `sg_market_bidding_state`

用途：记录每个 `年份 + 市场 + 订单类型` 标段的竞标状态。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `year_no` | INT | 年份 |
| `market_code` | VARCHAR(32) | 市场 |
| `order_type` | VARCHAR(32) | 订单类型 |
| `segment_code` | VARCHAR(64) | 标段编码，建议由市场和订单类型组合 |
| `release_sequence_no` | INT | 标段释放顺序 |
| `segment_status` | VARCHAR(32) | `MARKET_DISABLED / NO_ORDER_CONFIG / WAITING_INVESTMENT / SEQUENCE_READY / WAITING_RELEASE / SELECTING / COMPLETED / SKIPPED` |
| `leader_group_id` | BIGINT NULL | 本年该市场优先的市场龙头 |
| `leader_rule_json` | JSON NULL | 市场龙头计算依据 |
| `random_seed` | VARCHAR(64) NULL | 随机排序种子或结果摘要 |
| `current_group_id` | BIGINT NULL | 当前轮到的小组 |
| `investment_ready_at` | DATETIME NULL | 当年所有未破产组市场投入提交完成时间 |
| `released_at` | DATETIME NULL | 标段释放时间 |
| `completed_at` | DATETIME NULL | 标段选单完成时间 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_market_bidding_state(year_no, market_code, order_type)`
- `uk_segment_release_sequence(year_no, release_sequence_no)`
- `idx_segment_status(segment_status)`

#### 4.7.5 `sg_group_market_bid`

用途：记录每组在每年每个 `市场 + 订单类型` 提交的市场投入。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `group_id` | BIGINT | 小组 |
| `year_no` | INT | 年份 |
| `market_code` | VARCHAR(32) | 市场 |
| `order_type` | VARCHAR(32) | 订单类型 |
| `market_investment` | DECIMAL(18,2) | 市场投入 |
| `bid_status` | VARCHAR(32) | `SUBMITTED / LOCKED` |
| `submitted_at` | DATETIME | 提交时间 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_group_year_segment_bid(group_id, year_no, market_code, order_type)`
- `idx_year_segment_bid(year_no, market_code, order_type)`

说明：

- 提交后不可修改。
- 每组每年必须提交 16 条投入记录；`market_investment=0` 可保存，但普通小组不参与该标段选单。
- 市场投入从经营页前移到年度订单页提交，经营页只读带入汇总值。

#### 4.7.6 `sg_market_selection_order`

用途：记录系统生成的市场选单顺序。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `year_no` | INT | 年份 |
| `market_code` | VARCHAR(32) | 市场 |
| `order_type` | VARCHAR(32) | 订单类型 |
| `sequence_no` | INT | 选单顺序 |
| `group_id` | BIGINT | 小组 |
| `market_investment` | DECIMAL(18,2) | 排序时市场投入 |
| `previous_market_order_amount` | DECIMAL(18,2) | 上一年该市场订单总额，用于同投入时排序 |
| `is_market_leader` | TINYINT(1) | 是否市场龙头优先 |
| `rank_basis_json` | JSON | 排序依据与随机结果 |
| `selection_status` | VARCHAR(32) | `INELIGIBLE / WAITING / CURRENT / SELECTED / PASSED / ADMIN_SKIPPED` |
| `selected_order_id` | BIGINT NULL | 该轮选择的订单 |
| `selected_at` | DATETIME NULL | 选择时间 |
| `skipped_by_admin_id` | BIGINT NULL | 管理员代跳过操作人 |
| `skipped_reason` | VARCHAR(255) NULL | 管理员代跳过原因 |
| `skipped_at` | DATETIME NULL | 管理员代跳过时间 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_market_sequence(year_no, market_code, order_type, sequence_no)`
- `uk_group_market_sequence(year_no, market_code, order_type, group_id)`

说明：

- 玩家主动放弃记录为 `PASSED`；管理员代跳过记录为 `ADMIN_SKIPPED`，并同步写入管理员动作日志。
- 排序依据用于追溯，不在玩家端展示。

#### 4.7.7 `sg_group_order_selection`

用途：记录小组最终选择订单与交付状态。

关键字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | BIGINT | 主键 |
| `group_id` | BIGINT | 小组 |
| `year_no` | INT | 年份 |
| `market_code` | VARCHAR(32) | 市场 |
| `order_type` | VARCHAR(32) | 订单类型 |
| `order_id` | BIGINT | 订单 ID |
| `selection_status` | VARCHAR(32) | `SELECTED` |
| `delivery_status` | VARCHAR(32) | `SELECTED / DELIVERED / UNFINISHED` |
| `delivered_stage_code` | VARCHAR(16) NULL | `Q1 / Q2 / Q3 / Q4` |
| `delivered_at` | DATETIME NULL | 交付确认时间 |
| `selected_at` | DATETIME | 选择时间 |
| `creator/create_time/updater/update_time` | - | 审计字段 |

关键约束：

- `uk_group_year_segment_selection(group_id, year_no, market_code, order_type)`
- `uk_order_selected(order_id)`
- `idx_group_year_order_selection(group_id, year_no)`

说明：

- 每组每年每标段最多一个选择记录。
- 一个订单只能被一个小组选择。
- 账期字段来自订单池，首版只展示，不驱动经营页应收账款自动计算。
- 年末仍未交付时更新为 `UNFINISHED`，但首版不阻断财报提交。

---

## 5. 关系说明（逻辑）

- `sg_group` 1:N `sg_group_year_state`
- `sg_group` 1:1 `sg_initial_baseline`
- `sg_group` 1:N `sg_group_operating_draft`
- `sg_group` 1:N `sg_group_stage_submission`
- `sg_group` 1:N `sg_group_report`
- `sg_group` 1:N `sg_group_report_submission`
- `sg_group` 1:N `sg_group_summary_snapshot`
- `sg_group` 1:N `sg_notice`（按目标范围读取）
- `sg_group` 1:N `sg_group_adjustment`
- `sg_group` 1:N `sg_group_market_bid`
- `sg_group` 1:N `sg_market_selection_order`
- `sg_group` 1:N `sg_group_order_selection`
- `sg_group` 1:N `sg_admin_unlock_log`
- `sg_account` N:1 `sg_group`（玩家账号场景）
- `sg_game_config` 为单实例全局配置表，并锁定本场比赛沙盘版本包与字段/公式/流程版本
- `sg_order_generation_batch` 1:N `sg_order_generation_config`
- `sg_order_generation_batch` 1:N `sg_order_market_config`（通过锁定批次追溯）
- `sg_order_generation_batch` 1:N `sg_order_pool`
- `sg_order_pool` 1:0/1 `sg_group_order_selection`
- `sg_market_bidding_state` 1:N `sg_market_selection_order`

---

## 6. 索引与性能建议

### 6.1 必备索引

建议至少保留以下索引：

- `sg_group_year_state`：`uk_group_year`、`idx_year_status`
- `sg_group_stage_submission`：`uk_group_year_stage_version`、`idx_group_year`
- `sg_group_report_submission`：`uk_group_year_report_version`
- `sg_group_summary_snapshot`：`uk_group_year_summary`、`idx_year_no`
- `sg_notice`：`idx_target_scope_group`、`idx_published_at`、`idx_pinned`
- `sg_group_adjustment`：`idx_group_year_stage`、`idx_published_at`
- `sg_order_generation_config`：`uk_order_generation_config`、`idx_year_market`
- `sg_order_market_config`：`uk_order_market_config`、`idx_order_market_enabled`
- `sg_order_pool`：`idx_order_pool_year_market`、`idx_order_pool_status`
- `sg_market_bidding_state`：`uk_market_bidding_state`、`uk_segment_release_sequence`、`idx_segment_status`
- `sg_group_market_bid`：`uk_group_year_market_bid`、`idx_year_market_bid`
- `sg_market_selection_order`：`uk_market_sequence`、`uk_group_market_sequence`
- `sg_group_order_selection`：`uk_group_year_segment_selection`、`uk_order_selected`
- `sg_admin_unlock_log`：`idx_group_year`（可加）
- `sg_admin_action_log`：`idx_operate_time`

### 6.2 查询建议

- 汇总页查询优先走 `sg_group_summary_snapshot`
- 页面加载优先走“当前状态表 + 当前草稿/当前财报表”
- 历史查看、日志查看走提交流水表
- 年度订单页优先按 `sg_order_market_config + sg_market_bidding_state + sg_group_market_bid + sg_market_selection_order + sg_order_pool + sg_group_order_selection` 组合查询，并按市场开启状态与 `release_sequence_no` 展示标段释放进度
- 经营页正式年份订单总额与市场投入优先读取已锁定订单选择与市场投入汇总，不从经营页草稿反推

### 6.3 不建议的性能做法

- 不建议每次打开汇总页都全量即时重算全部年份
- 不建议用一个超大 JSON 表承载所有业务数据
- 不建议完全依赖前端缓存来判断状态
- 不建议用 Excel 文件作为运行时订单随机引擎；订单池应在管理员确认后固定落库

---

## 7. 变更规范

### 7.1 结构变更要求

数据库结构变更必须同步：

- 升级脚本
- 本文档
- 接口文档
- 对应 Repository / ORM 模型

### 7.2 变更脚本命名建议

- `sandbox-game-upgrade-1.sql`
- `sandbox-game-upgrade-2.sql`
- ...

### 7.3 幂等要求

升级脚本应尽量支持幂等：

- `IF NOT EXISTS`
- `ADD COLUMN IF NOT EXISTS`
- 索引存在性校验

---

## 8. 当前仍待后续补齐但不影响本版成立的内容

以下内容后续可补，但不影响数据库设计正式版成立：

- `operating_payload_json` 内部键名与 Excel 映射清单
- `report_manual_payload_json` 的字段级映射说明
- 沙盘版本包若后续从纯内置常量升级为数据库可查询配置，可补充 `sg_game_edition` / `sg_game_template_field` 等配置表；首版不要求落独立配置表
- 最终 SQL DDL 文件
- 数据初始化脚本（10 个组 + 1 个管理员）
- 订单模块最终 SQL DDL 与迁移脚本
- 订单生成批次样例与字段级校验规则
