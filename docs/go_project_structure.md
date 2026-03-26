# 沙盘经营系统 Go 项目结构与首版开发拆分文档（正式版）

> 更新日期：2026-03-25  
> 适用方式：基于《技术选型细化文档（Go 方向） v0.1》《接口设计文档（正式版）》《数据库设计文档（正式版）》《Excel 字段与页面映射文档（首版）》，把首版 Go 后端工程结构、目录职责与首批任务落点进一步具体化，作为正式建仓和主链开发的直接依据。  
> 文档定位：本文件回答“Go 项目目录怎么搭、代码先落到哪里、每个任务应该改哪些文件”，不替代业务需求、接口定义和数据库设计。

## 1. 输入基线

本文件直接承接以下文档：

- `docs/requirements_spec.md`
- `docs/minimal_state_machine.md`
- `docs/technical_selection.md`
- `docs/api_design.md`
- `docs/database_design.md`
- `docs/backend_guide.md`
- `docs/excel_field_mapping.md`
- `docs/calculation_rule_spec.md`

说明：

- 若本文件与以上正式文档冲突，以需求、接口、数据库、规则文档为准。
- 本文件重点是“工程落点”和“任务拆分”，不是重新定义业务规则。

---

## 2. 当前推荐的实现冻结项

为保证“周末可调试初版”的推进速度，当前建议直接冻结以下实现基线：

- 语言：`Go 1.22+`
- Web 框架：`Gin`
- 数据库：`MySQL 8.x`
- 迁移工具：`golang-migrate`
- 参数校验：`go-playground/validator`
- 日志：`Zap`
- 配置管理：`Viper`
- 测试：`testing + testify + httptest`
- 项目形态：`单仓库、单体后端服务`

关于 ORM / 数据访问：

- 为了首版交付速度，本文默认推荐 `GORM` 作为首版落地方案。
- 这不是业务规则冻结项，而是工程效率选择。
- 如果团队已有更熟的 `sqlx/sqlc` 基座，可以替换 Repository 实现，但目录职责、接口边界与状态机分层不应改变。

---

## 3. 推荐仓库结构

推荐仓库结构如下：

```text
sandbox-game/
  cmd/
    server/
      main.go
  configs/
    local.example.yaml
    dev.example.yaml
  migrations/
    mysql/
      0001_init.sql
      0002_seed_data.sql
  docs/
  internal/
    app/
      bootstrap.go
      router.go
    config/
      config.go
    enum/
      role_type.go
      year_status.go
      stage_status.go
      report_status.go
      business_status.go
      error_code.go
    http/
      dto/
        common_result.go
        player_operating_dto.go
        player_report_dto.go
        admin_control_dto.go
        admin_summary_dto.go
      handler/
        game_config_handler.go
        player_operating_handler.go
        player_report_handler.go
        admin_control_handler.go
        admin_summary_handler.go
        admin_group_data_handler.go
        audit_log_handler.go
      middleware/
        auth_middleware.go
        logger_middleware.go
        recovery_middleware.go
        error_middleware.go
    model/
      entity/
        account.go
        group.go
        game_config.go
        group_year_state.go
        group_operating_draft.go
        group_stage_submission.go
        group_report.go
        group_report_submission.go
        group_summary_snapshot.go
        initial_baseline.go
        admin_unlock_log.go
        admin_action_log.go
      payload/
        operating_payload.go
        report_manual_payload.go
        report_computed_payload.go
        baseline_payload.go
        state_snapshot.go
    repository/
      account_repository.go
      group_year_state_repository.go
      operating_repository.go
      report_repository.go
      summary_repository.go
      admin_control_repository.go
      audit_log_repository.go
    service/
      game_config_service.go
      player_operating_query_service.go
      player_operating_command_service.go
      player_report_query_service.go
      player_report_command_service.go
      admin_control_service.go
      admin_summary_service.go
      admin_group_data_service.go
      audit_log_service.go
    rules/
      context/
        calculation_context.go
      operating/
        operating_calculator.go
        operating_validator.go
      report/
        report_calculator.go
        report_validator.go
      summary/
        summary_builder.go
      carryforward/
        carry_forward_builder.go
    state/
      state_machine.go
      transition_guard.go
      state_snapshot_builder.go
    assembler/
      player_operating_assembler.go
      player_report_assembler.go
      admin_summary_assembler.go
      admin_group_data_assembler.go
    infra/
      db/
        mysql.go
      logger/
        zap.go
      timeutil/
        clock.go
  tests/
    http/
    integration/
```

---

## 4. 目录职责细化

### 4.1 `cmd/server`

- 只负责应用启动。
- 不承载业务逻辑。
- `main.go` 负责加载配置、初始化依赖、启动 HTTP 服务。

### 4.2 `internal/app`

- 负责应用装配。
- 推荐拆成：
  - `bootstrap.go`：创建 DB、Logger、Repository、Service
  - `router.go`：注册全部路由、中间件、健康检查接口

### 4.3 `internal/http`

职责边界：

- `handler`：接参、权限前置、调用 Service、回包
- `dto`：请求/响应对象
- `middleware`：认证、日志、错误处理、恢复

禁止事项：

- Handler 直接拼复杂 SQL
- Handler 直接做状态推进
- Handler 直接写公式计算

### 4.4 `internal/model`

建议拆成两层：

- `entity`：数据库实体
- `payload`：JSON 负载结构

这样做的原因：

- 经营页和财报页本身采用区块化 JSON 落库
- 但状态、汇总、日志等仍要保持强结构实体
- 可以避免把“数据库实体”和“Excel 风格表单负载”混在一个结构体里

### 4.5 `internal/repository`

- 一个 Repository 文件对应一个领域主表或一组强相关表。
- Repository 只负责读写，不负责判断业务是否允许。
- 涉及“当前有效版本”“最新草稿”“最新提交版本”的查询，应收敛在 Repository 中提供专门方法。

### 4.6 `internal/service`

建议按“查询服务 / 指令服务”拆分：

- 查询类：
  - `player_operating_query_service.go`
  - `player_report_query_service.go`
  - `admin_summary_service.go`
  - `admin_group_data_service.go`
- 写操作类：
  - `player_operating_command_service.go`
  - `player_report_command_service.go`
  - `admin_control_service.go`

这样做的价值：

- 读写链路更清晰
- 事务边界更集中
- 提交、解锁、开放下一年这类高风险操作更不容易被“顺手改坏”

### 4.7 `internal/rules`

规则层至少保持 4 个子域：

- `operating`
- `report`
- `summary`
- `carryforward`

推荐分工：

- `operating_calculator`：经营页计算、阶段派生结果、现金结果
- `operating_validator`：当前阶段必填校验、破产前置校验
- `report_calculator`：财报自动计算
- `report_validator`：平衡校验、税率合法性校验
- `summary_builder`：从财报结果生成汇总快照
- `carry_forward_builder`：从上年财报生成次年初始视图

### 4.8 `internal/state`

状态层必须独立存在，不应散落在 Service 中。

建议至少包含：

- `state_machine.go`：状态枚举与转移入口
- `transition_guard.go`：状态允许性判断
- `state_snapshot_builder.go`：提交前后状态摘要生成

### 4.9 `internal/assembler`

Assembler 用于解决两个问题：

- 数据库存储结构与前端视图结构不同
- 规则层输出与接口返回对象不同

典型场景：

- 经营页视图组装
- 财报页视图组装
- 管理员汇总页组装
- 管理员查看组数据视图组装

---

## 5. 首版核心模块与职责映射

| 模块 | 主要职责 | 直接依赖 |
|---|---|---|
| `game-config` | 当前配置、最终年份、当前开放年份 | `service + repository + handler` |
| `player-operating` | 经营页读取、草稿保存、阶段提交 | `service + rules.operating + state + assembler` |
| `player-report` | 财报页读取、草稿保存、财报提交 | `service + rules.report + rules.summary + state + assembler` |
| `admin-control` | 初始基线、开放下一年、异常解锁 | `service + state + repository + audit` |
| `admin-summary` | 汇总页、最终排名 | `service + assembler + repository` |
| `admin-group-data` | 管理员查看任意组任一年数据 | `service + assembler + repository` |
| `audit-log` | 审计日志查询 | `service + repository + handler` |

---

## 6. 首版最小文件落点

### 6.1 对应 `M1-01 ~ M1-06`

| Task ID | 需要落的核心文件/目录 | 说明 |
|---|---|---|
| `M1-01` | `cmd/server/main.go`、`internal/app/bootstrap.go`、`internal/app/router.go`、`configs/*.yaml` | 建起可启动骨架 |
| `M1-02` | `internal/http/dto/common_result.go`、`internal/http/middleware/*.go`、`internal/enum/error_code.go`、`internal/infra/logger/zap.go` | 统一响应体、错误码、中间件、日志 |
| `M1-03` | `migrations/mysql/0001_init.sql`、`internal/model/entity/*.go` | 首版表结构与实体定义 |
| `M1-04` | `migrations/mysql/0002_seed_data.sql` 或初始化脚本 | 10 个组、1 个管理员、默认配置 |
| `M1-05` | `internal/state/*.go`、`internal/enum/*status.go` | 状态机和状态判断 |
| `M1-06` | `internal/rules/context/*`、`internal/rules/operating/*`、`internal/rules/report/*`、`internal/rules/summary/*`、`internal/rules/carryforward/*` | 规则层壳子与计算上下文 |

### 6.2 对应 `M1-07 ~ M1-12`

| Task ID | 需要落的核心文件/目录 | 说明 |
|---|---|---|
| `M1-07` | `player_operating_handler.go`、`player_operating_query_service.go`、`player_operating_assembler.go` | 获取经营页视图 |
| `M1-08` | `player_operating_handler.go`、`player_operating_command_service.go`、`operating_repository.go` | 经营草稿保存 |
| `M1-09` | `player_operating_command_service.go`、`rules/operating/*`、`state/*` | 阶段提交、阶段锁定、破产判定 |
| `M1-10` | `player_report_handler.go`、`player_report_query_service.go`、`player_report_assembler.go` | 获取财报页视图 |
| `M1-11` | `player_report_command_service.go`、`report_repository.go` | 财报草稿保存 |
| `M1-12` | `player_report_command_service.go`、`rules/report/*`、`rules/summary/*`、`summary_repository.go` | 财报提交、汇总快照写入 |

### 6.3 对应 `M2-01 ~ M2-06`

| Task ID | 需要落的核心文件/目录 | 说明 |
|---|---|---|
| `M2-01` | `game_config_handler.go`、`game_config_service.go` | 最终年份与当前配置 |
| `M2-02` | `admin_control_handler.go`、`admin_control_service.go`、`baseline_payload.go` | 初始基线录入与锁定 |
| `M2-03` | `admin_control_service.go`、`state/*`、`audit_log_repository.go` | 管理员开放下一年 |
| `M2-04` | `admin_summary_handler.go`、`admin_summary_service.go`、`admin_summary_assembler.go` | 管理员汇总与排名 |
| `M2-05` | `admin_group_data_handler.go`、`admin_group_data_service.go`、`admin_group_data_assembler.go` | 查看组数据 |
| `M2-06` | `admin_control_service.go`、`audit_log_repository.go`、`state_snapshot_builder.go` | 异常解锁与日志 |

---

## 7. 请求链路建议

### 7.1 获取经营页

建议调用链：

`Handler -> QueryService -> Repository(状态/草稿/最新提交/基线) -> Rules/CarryForward -> Assembler -> Response`

重点：

- 当前年份若未开放，只返回锁定态视图
- 当前年份若无草稿，也要能用“基线 + 上年结果”组装初始视图

### 7.2 提交经营阶段

建议调用链：

`Handler -> CommandService -> StateGuard -> RulesValidator -> RulesCalculator -> Repository(Transaction) -> StateMachine -> AuditLog`

重点：

- 只校验当前阶段应校验的区块
- 不重复校验已锁定历史区
- 破产判定要在事务内与状态写入保持一致

### 7.3 提交财报

建议调用链：

`Handler -> CommandService -> StateGuard -> ReportValidator -> ReportCalculator -> SummaryBuilder -> Repository(Transaction) -> StateMachine`

重点：

- 财报提交必须生成当前计算结果快照
- 正式年份提交成功后，才写正式汇总快照
- `0年` 提交成功后允许承接到 `1年`，但不计排名

### 7.4 开放下一年 / 异常解锁

建议调用链：

`Handler -> AdminControlService -> StateGuard -> Repository(Transaction) -> StateMachine -> Summary撤回/状态恢复 -> AdminActionLog`

重点：

- `开放下一年` 是全组统一动作
- `异常解锁` 只允许在下一年未开放前执行
- 若已生成汇总快照，解锁时必须撤回汇总有效性

---

## 8. 配置、迁移与初始化建议

### 8.1 配置文件

建议最少提供：

- `configs/local.example.yaml`
- `configs/dev.example.yaml`

配置项建议：

- `server.port`
- `server.readTimeout`
- `server.writeTimeout`
- `mysql.dsn`
- `mysql.maxOpenConns`
- `mysql.maxIdleConns`
- `log.level`
- `auth.mode`

### 8.2 迁移脚本

脚本命名建议：

- `0001_init.sql`
- `0002_seed_data.sql`
- `0003_add_unlock_log.sql`

要求：

- 能独立初始化本地库
- 变更后同步更新 `docs/database_design.md`

### 8.3 种子数据

本地开发至少要有：

- 1 个管理员账号
- 10 个组账号
- 10 个组主数据
- 默认 `final_year`
- `0年 ~ 最终年` 的年份状态初始化逻辑

---

## 9. 周末可调试初版的建议落地顺序

若目标是尽快做出“能登录、能看经营页、能提交流程、能看汇总”的版本，建议按下面顺序推进：

1. 先完成工程骨架和迁移
2. 再完成状态机与规则层壳子
3. 然后打通玩家经营主链
4. 再打通财报提交与汇总快照
5. 最后补管理员开放下一年与异常解锁

不建议一开始就做：

- 前端逐像素还原
- OpenAPI 全量自动生成
- 全量 Excel 逐格映射
- 泛化的公式引擎
- 复杂权限平台对接

---

## 10. 当前阶段的最小建仓清单

如果现在就进入正式编码，第一批应该先创建这些文件或目录：

- `cmd/server/main.go`
- `internal/app/bootstrap.go`
- `internal/app/router.go`
- `internal/config/config.go`
- `internal/http/dto/common_result.go`
- `internal/http/middleware/error_middleware.go`
- `internal/http/middleware/logger_middleware.go`
- `internal/enum/error_code.go`
- `internal/enum/year_status.go`
- `internal/enum/stage_status.go`
- `internal/enum/report_status.go`
- `internal/enum/business_status.go`
- `internal/model/entity/group_year_state.go`
- `internal/model/payload/operating_payload.go`
- `internal/model/payload/report_manual_payload.go`
- `internal/model/payload/baseline_payload.go`
- `internal/repository/group_year_state_repository.go`
- `internal/rules/context/calculation_context.go`
- `internal/rules/operating/operating_validator.go`
- `internal/rules/report/report_validator.go`
- `internal/state/state_machine.go`
- `migrations/mysql/0001_init.sql`
- `migrations/mysql/0002_seed_data.sql`

这些文件建起来以后，`M1-01 ~ M1-06` 基本就有了明确落点，不会再停留在“知道要做，但不知道代码该从哪儿开始”的状态。

---

## 11. 变更控制原则

- 若后续只调整 Excel 中文标签或局部位置，本文件不需要推翻，只需同步 `docs/excel_field_mapping.md`。
- 若后续修改接口域划分、主表结构、状态机主语义，本文件必须同步更新。
- 若团队最终决定不用 `GORM`，只需调整 Repository 实现方案，不应改动 Handler / Service / Rules / State 的职责边界。
