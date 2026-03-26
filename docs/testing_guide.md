# 沙盘经营系统测试指南（正式版）

> 更新日期：2026-03-25  
> 适用方式：基于《正式需求文档（首版）》《接口设计文档（正式版）》《数据库设计文档（正式版）》《开发执行拆解》，定义首版测试分层、关键场景与质量门禁。  
> 文档定位：本文件回答“这个系统要怎么测、测到什么程度算可交付”。

## 1. 测试目标

首版测试至少要验证以下三件事：

1. 主链接口可用
2. 核心规则正确
3. 页面关键流程可回归

本项目不是普通 CRUD，测试重点必须放在：

- 状态推进
- 季末现金展示口径
- 财报平衡校验
- 异常解锁
- 破产判定
- 汇总与排名口径

---

## 2. 测试分层

### 2.1 API 调试/回归（`.http`）

统一规范：

- 所有 `.http` 用例统一放在：`tests/http/sandbox-game/`
- 新增接口或修改接口时，必须同步补齐 `.http` 用例
- 至少覆盖：成功 / 参数错误 / 业务错误 / 权限或状态冲突

推荐环境变量：

```http
@baseUrl = http://127.0.0.1:8080/api/v1
@adminToken = <admin-token>
@groupToken = <group-token>
```

### 2.2 后端单元测试

- 测试文件统一使用 Go 原生 `_test.go`
- 优先覆盖：
  - 状态机
  - 经营提交规则
  - 财报平衡校验
  - 跨年传递规则
  - 破产判定
  - 汇总快照生效/撤回

### 2.3 前端测试

- 工具：`Vitest`
- 至少覆盖：
  - 页面加载
  - 锁定态渲染
  - 提交前错误提示
  - 提交成功/失败提示
  - 管理员异常解锁交互

### 2.4 数据库结构冒烟（必做）

涉及数据库变更时，至少验证：

- 核心表存在
- 唯一约束存在
- 状态字段存在
- 汇总快照表存在
- 解锁日志表存在

---

## 3. 首版接口必测场景

每个变更接口至少覆盖以下四类场景：

1. 正常场景
2. 参数异常
3. 权限场景
4. 状态或规则场景

### 3.1 状态码断言要求

- 每个写接口至少断言一次 `2xx`
- 每个写接口至少断言一次 `400`
- 存在重复动作或状态冲突的接口，至少断言一次 `409`
- 存在规则校验的接口，至少断言一次 `422`
- 失败场景需同时断言 HTTP 状态码和响应体 `code`

---

## 4. 模块专项测试点

### 4.1 玩家经营页

必须覆盖：

- `Q1` 只校验 `年初 + Q1`
- `Q2/Q3/Q4` 只校验各自阶段
- `年末` 只校验年末区手工项
- 默认空值不允许直接提交
- 手工输入 `0` 后允许提交
- 经营页展示对应阶段的季末现金核对值
- 季末现金核对值不作为阻断提交条件
- 提交后上一阶段锁定

### 4.2 玩家财报页

必须覆盖：

- 年末提交成功后财报页开放
- 财报页绿色手工项可填，自动项只读
- 所得税税率仅允许 `0.25 / 0.15 / 0`
- 财报平衡校验通过后才允许提交
- 提交成功后本年完成
- 已完成年份财报不可重复提交

### 4.3 管理员控制

必须覆盖：

- 统一初始基线可录入并锁定，并同步应用到全部小组
- 最终年份可配置
- `finalYear < currentOpenYear` 时更新最终年份返回 `422`
- 传入与当前一致的 `finalYear` 按幂等更新返回成功
- 管理员可开放下一年
- 未提交统一初始基线时，开放下一年返回 `422`
- 当前开放年份仍有未破产组未完成财报时，开放下一年返回 `422`
- `targetYearNo` 与服务端当前状态不一致时，开放下一年返回 `409`
- 重复开放下一年返回冲突或规则错误
- 下一年已开放后，禁止异常解锁上一年
- 对已处于可编辑态的年份重复异常解锁返回 `409`
- 若被解锁年份曾触发破产，解锁成功后应恢复为临时可编辑态
- 异常解锁后财报失效、汇总撤回

### 4.4 汇总与排名

必须覆盖：

- 汇总从 `1年` 开始，不显示 `0年`
- 仅“已提交财报”的年份进入正式汇总
- 异常解锁后汇总撤回
- 最终排名按最终年份所有者权益排序
- 最终排名仅在所有未破产组完成最终年份财报后开放
- 破产组显示状态标识

### 4.5 破产与只读控制

必须覆盖：

- 经营提交节点触发现金流断裂后判定破产
- 最终年份所有者权益为负后判定破产
- 破产后后续年份不可编辑
- 破产后历史年份可查看

---

## 5. `.http` 用例建议清单

建议至少创建以下文件：

- `tests/http/sandbox-game/GameConfig.http`
- `tests/http/sandbox-game/PlayerOperating.http`
- `tests/http/sandbox-game/PlayerReport.http`
- `tests/http/sandbox-game/AdminControl.http`
- `tests/http/sandbox-game/AdminSummary.http`
- `tests/http/sandbox-game/AdminGroupData.http`
- `tests/http/sandbox-game/AuditLog.http`

### 5.1 示例

```http
### 获取当前配置
GET {{baseUrl}}/sandbox-game/game-config/get-current
Authorization: Bearer {{groupToken}}

### 提交 Q1
POST {{baseUrl}}/sandbox-game/player-operating/submit-stage
Content-Type: application/json
Authorization: Bearer {{groupToken}}

{
  "yearNo": 0,
  "stageStatus": "Q1_OPEN",
  "operatingPayload": {}
}
```

---

## 6. 数据库结构冒烟建议

涉及 SQL 变更时，建议至少执行：

```sql
SELECT DATABASE();
SHOW TABLES LIKE 'sg_group_year_state';
SHOW TABLES LIKE 'sg_group_stage_submission';
SHOW TABLES LIKE 'sg_group_report_submission';
SHOW TABLES LIKE 'sg_group_summary_snapshot';
SHOW TABLES LIKE 'sg_admin_unlock_log';

DESC sg_group_year_state;
DESC sg_group_stage_submission;
DESC sg_group_report_submission;
DESC sg_group_summary_snapshot;
```

判定标准：

- 核心表存在
- 状态字段存在且命名正确
- `group_id + year_no` 唯一约束存在
- 汇总快照与日志表存在

---

## 7. 人工验收主链

首版至少按以下链路完整走一遍：

1. 管理员录入统一初始基线
2. 玩家进入 `0年经营`
3. 玩家完成 `Q1 -> Q2 -> Q3 -> Q4 -> 年末` 提交
4. 玩家进入 `0年财报` 并提交
5. 管理员开放 `1年`
6. 玩家完成 `1年经营 + 1年财报`
7. 管理员查看汇总
8. 管理员对某组执行异常解锁
9. 被解锁组重新提交经营或财报
10. 汇总重新生效

---

## 8. 质量门禁建议

提交前至少满足：

1. 后端可编译
2. 前端类型检查通过
3. 涉及 SQL 变更时完成数据库结构冒烟
4. 受影响接口 `.http` 用例冒烟通过
5. 高风险规则具备自动化测试或人工回归记录
6. 文档同步更新
