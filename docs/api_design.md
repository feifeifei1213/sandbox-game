# 沙盘经营系统接口设计文档（正式版）

> 更新日期：2026-03-30  
> 适用方式：基于《正式需求文档（首版）》《最小状态机 v0.1》《技术选型细化文档（Go 方向） v0.1》，定义首版正式业务接口边界，作为后续 Go 后端开发、前端 API 客户端开发和接口测试的统一依据。  
> 文档定位：本文件定义接口域划分、路径风格、请求/响应结构、核心动作语义、关键错误状态和测试口径。  
> 说明：本文件已按 `1组 最终版.xlsx` 口径收口；接口域划分、动作语义和结构性字段已经冻结，若后续仅有页面标签细修，应优先更新映射文档，不直接改动接口结构。

## 1. 适用范围

本规范适用于沙盘经营系统全部首版业务接口：

- 首版最小登录接口
- 玩家端经营页接口
- 玩家端财报页接口
- 管理员汇总与控制接口
- 管理员组数据查看接口
- 审计日志查询接口

本文件不覆盖：

- 公司统一认证平台的内部实现
- 网关配置细节
- 非业务型静态资源接口

---

## 2. 基础约定

### 2.1 访问前缀

首版推荐统一前缀如下：

- API 根前缀：`/api/v1`
- 业务前缀：`/sandbox-game/*`
- 完整示例：`/api/v1/sandbox-game/player-operating/get-year-view`

说明：

- 若后续接入公司统一网关，可在网关层映射为 `/admin-api` 或其他统一前缀。
- 业务路径本身不应因为网关变化而重构。

### 2.2 认证与角色识别

首版正式版要求提供项目内最小登录闭环。

约束如下：

- 系统提供统一登录页与最小登录接口。
- 用户通过用户名和密码登录。
- 登录成功后，系统按账号角色自动识别玩家或管理员身份，不要求用户手动选择角色。
- 玩家登录成功后默认进入 `0年经营`；管理员登录成功后默认进入 `汇总页`。
- 业务接口默认要求调用方已完成认证。

业务侧最少需要识别以下身份信息：

- `userId`
- `username`
- `roleType`：`ADMIN` / `GROUP`
- `groupId`：管理员可为空；玩家必须具备所属组

推荐请求头：

- `Authorization: Bearer <accessToken>`

说明：

- 若后续复用公司统一认证，本文件中的业务接口无需整体改名。
- 首版即使采用项目内最小登录模块，业务接口层也只应依赖“已解析好的身份信息”，不把登录逻辑耦合到业务域。

### 2.3 权限命名

权限码统一建议：`sandbox-game:<domain>:<action>`

示例：

- `sandbox-game:player-operating:query`
- `sandbox-game:player-operating:submit-stage`
- `sandbox-game:player-report:submit`
- `sandbox-game:admin-control:open-next-year`
- `sandbox-game:admin-control:unlock-year`
- `sandbox-game:admin-summary:query`

### 2.4 URL 与方法规范

本项目采用 **动作式接口风格**，保持与参考规范一致。

约定：

- 查询详情：`GET /get-*` 或 `GET /.../get-*`
- 分页：`GET /page-*`
- 草稿保存：`PUT /save-draft`
- 正式提交：`POST /submit-*`
- 管理动作：`POST /open-next-year`、`POST /unlock-year` 等
- 配置更新：`PUT /update-*`

禁止：

- 同一业务域混用多套命名风格
- 用 `GET` 承载有副作用的业务动作
- 把复杂业务逻辑塞进路由层

---

## 3. 响应格式与错误语义

### 3.1 统一响应格式

成功响应统一采用：

```json
{
  "code": 0,
  "msg": "",
  "data": {}
}
```

说明：

- 成功：HTTP `2xx`，响应体 `code=0`
- 失败：HTTP `4xx/5xx`，响应体 `code` 与 HTTP 状态码保持一致
- 前端禁止通过 `msg` 文案做业务分支判断

### 3.2 分页结构

分页数据统一：

```json
{
  "list": [],
  "total": 0
}
```

### 3.3 错误状态码语义

本项目统一采用 REST 状态语义：

- `400`：参数错误
- `401`：未登录或令牌无效
- `403`：无权限
- `404`：资源不存在
- `409`：状态冲突或重复动作
- `422`：业务规则不满足
- `500`：系统异常

### 3.4 本项目重点使用边界

- `409` 适用于：
  - 当前阶段已提交，重复提交
  - 年份已开放，重复开放下一年
  - 年度已完成，再次提交财报
- `422` 适用于：
  - 当前阶段必填项未完成
  - 财报平衡校验失败
  - 当前年份未开放
  - 当前阶段未到可提交时机
  - 破产组继续提交

---

## 4. 业务接口域划分

首版正式接口域建议如下：

| 领域 | 说明 |
|---|---|
| `auth` | 登录、当前用户信息、退出登录 |
| `game-config` | 当前游戏配置、开放年份、规则版本信息 |
| `player-operating` | 玩家经营页读取、草稿保存、阶段提交 |
| `player-report` | 财报页读取、草稿保存、财报提交 |
| `admin-summary` | 汇总页、最终排名 |
| `admin-control` | 最终年份设置、开放下一年、初始基线、异常解锁 |
| `admin-group-data` | 管理员查看任意组任一年经营/财报数据 |
| `audit-log` | 提交日志、解锁日志、管理员动作日志 |

说明：

- 首版不建议按“前端页面路径”来拆接口域。
- 首版也不建议把所有接口塞进单一 `game` 域，避免后续维护混乱。

---

## 5. 核心状态与结构性字段

### 5.1 状态枚举

接口中涉及的结构性状态字段，正式版统一如下：

- 年份主状态：
  - `LOCKED`
  - `OPERATING`
  - `REPORT_PENDING`
  - `REPORTING`
  - `COMPLETED`
- 经营阶段状态：
  - `Q1_OPEN`
  - `Q2_OPEN`
  - `Q3_OPEN`
  - `Q4_OPEN`
  - `YEAR_END_OPEN`
- 财报状态：
  - `REPORT_LOCKED`
  - `REPORT_OPEN`
  - `REPORT_SUBMITTED`
- 经营状态：
  - `NORMAL`
  - `BANKRUPT`

### 5.2 已冻结的 Excel 命名映射边界

当前已以 `game doc/1组 最终版.xlsx` 作为命名冻结基线。

以下内容已经明确：

- `operatingPayload` / `reportManualPayload` 的系统字段名保持稳定
- Excel 中文显示名以 `docs/excel_field_mapping.md` 为准
- 页面标签若继续细修，应通过映射文档和前端文案收口完成

因此首版接口正式版采用：

- 结构性字段冻结
- 区块负载字段名稳定
- 显示层命名通过映射文档维护

---

## 6. 领域接口设计

### 6.0 `auth`（首版最小登录）

#### 6.0.1 登录

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/auth/login`
- 权限：匿名可访问

请求体建议：

```json
{
  "username": "group01",
  "password": "******"
}
```

业务规则：

- 账号必须存在且处于启用状态。
- 用户名与密码必须匹配。
- 不允许前端手动声明本次以“玩家”还是“管理员”身份登录。
- 登录成功后，服务端根据账号角色返回默认跳转地址。

返回字段建议：

| 字段 | 说明 |
|---|---|
| `accessToken` | 访问令牌 |
| `tokenType` | 固定返回 `Bearer` |
| `expiresIn` | 令牌有效期，单位秒 |
| `user.userId` | 用户 ID |
| `user.username` | 登录名 |
| `user.roleType` | `ADMIN` / `GROUP` |
| `user.groupId` | 玩家所属组，管理员为空 |
| `defaultRoute` | 默认跳转地址；玩家为 `player/operating?yearNo=0`，管理员为 `admin/summary` |

#### 6.0.2 获取当前登录用户

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/auth/get-current-user`
- 权限：已登录用户可访问

返回应至少包含：

- `userId`
- `username`
- `roleType`
- `groupId`
- `defaultRoute`

用途：

- 前端刷新后恢复当前登录用户信息
- 路由守卫判断当前是玩家还是管理员
- 顶部导航或退出登录区域展示用户名

#### 6.0.3 退出登录

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/auth/logout`
- 权限：已登录用户可访问

规则：

- 成功后清理当前登录态。
- 前端收到成功响应后返回统一登录页。
- 首版不要求实现复杂的多端会话管理。

### 6.1 `game-config`

#### 6.1.1 获取当前配置

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/game-config/get-current`
- 权限：已登录用户可访问

返回示例字段：

| 字段 | 说明 |
|---|---|
| `currentOpenYear` | 当前开放年份 |
| `finalYear` | 最终年份配置 |
| `ruleVersion` | 当前规则版本 |
| `templateVersion` | 当前模板版本 |
| `demoYearEnabled` | 是否启用 `0年` 引导年，首版固定为 `true` |

#### 6.1.2 获取年份标签状态

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/game-config/get-year-tabs`
- 权限：已登录用户可访问

用途：

- 返回前端顶部年份标签状态
- 区分：可进入 / 锁定 / 已完成 / 已破产只读
- 前端应以本接口结果动态渲染年份按钮数量与可进入状态，不得把年份数量写死在页面中

---

### 6.2 `player-operating`

#### 6.2.1 获取经营页视图

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/player-operating/get-year-view?yearNo=0`
- 权限：`sandbox-game:player-operating:query`

返回应至少包含：

| 字段 | 说明 |
|---|---|
| `groupId` | 当前玩家所属组 |
| `yearNo` | 当前年份 |
| `yearStatus` | 年份主状态 |
| `stageStatus` | 当前经营阶段状态 |
| `reportStatus` | 当前财报状态 |
| `businessStatus` | 正常 / 已破产 |
| `hasInvalidDraft` | 当前经营页是否存在失效草稿 |
| `invalidScopes` | 已保留但失效的经营区域列表 |
| `hasRetainedReportDraft` | 本年是否仍保留失效的财报草稿 |
| `operatingPayload` | 当前经营页完整业务数据 |
| `editableScopes` | 当前允许编辑的区域列表 |
| `readonlyScopes` | 当前只读区域列表 |
| `stageSubmitHistory` | 各阶段已提交摘要 |
| `lastDraftSavedAt` | 最近草稿保存时间 |

#### 6.2.2 保存经营页草稿

- 方法：`PUT`
- 路径：`/api/v1/sandbox-game/player-operating/save-draft`
- 权限：`sandbox-game:player-operating:update`

请求体建议：

```json
{
  "yearNo": 1,
  "stageStatus": "Q2_OPEN",
  "operatingPayload": {},
  "clientSaveTime": "2026-03-24T10:00:00+08:00"
}
```

规则：

- 只保存草稿，不推进状态
- 只允许保存当前可编辑年份
- 保存策略：最后一次成功保存覆盖当前草稿

#### 6.2.3 提交经营阶段

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/player-operating/submit-stage`
- 权限：`sandbox-game:player-operating:submit-stage`

请求体建议：

```json
{
  "yearNo": 1,
  "stageCode": "Q2",
  "operatingPayload": {}
}
```

业务校验：

- 当前年份必须处于 `OPERATING`
- 当前阶段必须与 `stageCode` 匹配
- 当前阶段手工项完整
- 当前组未破产

成功结果：

- 记录阶段提交日志
- 推进到下一经营阶段，或进入 `REPORT_PENDING`
- 返回最新状态摘要
- 返回本次系统计算的期末现金，供玩家线下核对

重复提交规则：

- 同一阶段已正式提交后再次提交，返回 `409`

---

### 6.3 `player-report`

#### 6.3.1 获取财报页视图

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/player-report/get-view?yearNo=1`
- 权限：`sandbox-game:player-report:query`

返回应至少包含：

| 字段 | 说明 |
|---|---|
| `groupId` | 当前玩家所属组 |
| `yearNo` | 当前年份 |
| `yearStatus` | 年份主状态 |
| `reportStatus` | 财报状态 |
| `businessStatus` | 正常 / 已破产 |
| `hasInvalidDraft` | 当前财报页是否为失效草稿待重提状态 |
| `reportComputedPayload` | 自动计算结果 |
| `reportManualPayload` | 手工项当前值 |
| `manualFieldOptions` | 如税率下拉选项 |
| `lastDraftSavedAt` | 最近草稿保存时间 |

#### 6.3.2 保存财报草稿

- 方法：`PUT`
- 路径：`/api/v1/sandbox-game/player-report/save-draft`
- 权限：`sandbox-game:player-report:update`

请求体建议：

```json
{
  "yearNo": 1,
  "reportManualPayload": {}
}
```

规则：

- 仅保存当前财报手工项草稿
- 不推进年份状态

#### 6.3.3 提交财报

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/player-report/submit`
- 权限：`sandbox-game:player-report:submit`

请求体建议：

```json
{
  "yearNo": 1,
  "reportManualPayload": {}
}
```

业务校验：

- 当前年份必须处于 `REPORT_PENDING` 或 `REPORTING`
- 财报页已开放
- 手工项完整
- `总资产 = 总负债 + 总权益`
- 当前组未破产

成功结果：

- 财报状态变为 `REPORT_SUBMITTED`
- 年份主状态变为 `COMPLETED`
- 正式年份写入汇总快照

重复提交规则：

- 年份已完成后再次提交财报，返回 `409`

---

### 6.4 `admin-summary`

#### 6.4.1 获取年度汇总

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-summary/get-year-summary?yearNo=1`
- 权限：`sandbox-game:admin-summary:query`

返回字段建议：

| 字段 | 说明 |
|---|---|
| `yearNo` | 正式年份 |
| `list[].groupId` | 组 ID |
| `list[].groupName` | 组名称 |
| `list[].revenue` | 收入 |
| `list[].profit` | 利润 |
| `list[].equity` | 权益 |
| `list[].businessStatus` | 经营状态 |
| `list[].ranking` | 若为最终年份，可返回最终排名 |

规则：

- `0年` 不进入正式汇总
- 仅 `COMPLETED` 的正式年份计入汇总
- 被异常解锁且未重提财报的年份，不得出现在正式汇总口径中

#### 6.4.2 获取最终排名

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-summary/get-final-ranking`
- 权限：`sandbox-game:admin-summary:query`

规则：

- 仅在最终年份结果可用时返回有效排名
- 排名依据：最终年份 `equity` 倒序

---

### 6.5 `admin-control`

#### 6.5.1 获取控制台配置

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-control/get-config`
- 权限：`sandbox-game:admin-control:query`

Go DTO 建议：

- 响应：`AdminControlConfigResp`

返回字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `finalYear` | `int` | 最终年份配置 |
| `currentOpenYear` | `int` | 当前开放年份 |
| `canOpenNextYear` | `bool` | 当前是否允许开放下一年 |
| `nextOpenableYear` | `int` | 若允许开放，下一次将开放的年份 |
| `openNextYearBlockedReason` | `string` | 不允许开放时的阻塞原因 |
| `ruleVersion` | `string` | 当前规则版本 |
| `templateVersion` | `string` | 当前模板版本 |
| `initialBaselineSubmitted` | `bool` | 初始基线是否已提交 |
| `initialBaselineSubmittedAt` | `string \| null` | 初始基线提交时间，RFC3339 |
| `initialBaselineSubmitterName` | `string \| null` | 初始基线提交人 |
| `latestAdminAction` | `object \| null` | 最近一次关键管理员动作摘要 |

说明：

- 本接口服务于管理员年度控制页顶部状态区。
- `openNextYearBlockedReason` 只返回原因摘要，不返回逐组明细。

#### 6.5.2 更新最终年份

- 方法：`PUT`
- 路径：`/api/v1/sandbox-game/admin-control/update-final-year`
- 权限：`sandbox-game:admin-control:update-final-year`

Go DTO 建议：

- 请求：`UpdateFinalYearReq`
- 响应：`UpdateFinalYearResp`

请求体：

```json
{
  "finalYear": 3
}
```

请求字段建议：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `finalYear` | `int` | 是 | 最终年份，必须大于等于 `currentOpenYear` |

响应字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `finalYear` | `int` | 更新后的最终年份 |
| `currentOpenYear` | `int` | 当前开放年份 |
| `initializedFromYear` | `int \| null` | 若本次上调扩年，自动补齐的起始年份 |
| `initializedToYear` | `int \| null` | 若本次上调扩年，自动补齐的结束年份 |
| `initializedYearCount` | `int` | 本次自动补齐的年份层数量，未扩年时为 `0` |
| `updatedAt` | `string` | 更新时间，RFC3339 |
| `updatedBy` | `string` | 操作管理员名称 |

规则：

- 仅管理员可执行。
- `finalYear` 不得小于 `currentOpenYear`。
- 若 `finalYear` 上调到当前已准备年份上限之外，服务端必须自动补齐全部小组缺失的未来年份状态记录。
- 自动补齐的未来年份默认应为：`FORMAL + LOCKED + Q1_OPEN + REPORT_LOCKED + summaryEffective=false + latest submit version=0`。
- 若 `finalYear` 下调且仍不小于 `currentOpenYear`，允许更新，但不物理删除已存在的未来年份数据。
- 更新成功后必须写入管理员动作日志。

#### 6.5.3 开放下一年

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-control/open-next-year`
- 权限：`sandbox-game:admin-control:open-next-year`

Go DTO 建议：

- 请求：`OpenNextYearReq`
- 响应：`OpenNextYearResp`

请求体：

```json
{
  "targetYearNo": 2
}
```

规则：

- 仅管理员可执行
- `targetYearNo` 必须等于 `currentOpenYear + 1`
- 不能超过 `finalYear`
- 仅当“当前开放年份下，全部未破产组已完成本年财报”时才允许开放
- 已破产组在新年份仍保持不可编辑
- 本接口默认依赖 `update-final-year` 已经补齐目标年份的状态空间；若目标年份状态记录缺失，应视为服务端初始化缺陷，而不是前端调用方式错误

响应字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `previousOpenYear` | `int` | 开放前年份 |
| `currentOpenYear` | `int` | 开放后年份 |
| `openedYearNo` | `int` | 本次实际开放年份 |
| `finalYear` | `int` | 当前最终年份 |
| `canOpenNextYear` | `bool` | 开放完成后是否还能继续开放下一年 |
| `openNextYearBlockedReason` | `string` | 若已到最终年份或后续被阻塞，返回摘要原因 |
| `latestAdminAction` | `object` | 本次动作摘要 |

#### 6.5.4 获取初始基线

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-control/get-initial-baseline`
- 权限：`sandbox-game:admin-control:query-initial-baseline`

说明：

- 首版页面语义采用“`1 份共享初始基线模板`”，不是按组分别录入。
- 后端可继续按组落库，但接口层只暴露一份共享模板。

Go DTO 建议：

- 响应：`InitialBaselineResp`

响应字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `submitted` | `bool` | 是否已提交 |
| `editable` | `bool` | 当前是否允许编辑 |
| `baselinePayload` | `object` | 初始基线业务数据 |
| `appliedGroupCount` | `int` | 已应用的小组数量 |
| `submitterName` | `string \| null` | 提交人 |
| `submittedAt` | `string \| null` | 提交时间，RFC3339 |

#### 6.5.5 提交初始基线

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-control/submit-initial-baseline`
- 权限：`sandbox-game:admin-control:submit-initial-baseline`

Go DTO 建议：

- 请求：`SubmitInitialBaselineReq`
- 响应：`SubmitInitialBaselineResp`

请求体建议：

```json
{
  "baselinePayload": {}
}
```

规则：

- 本次提交的是共享初始基线模板
- 提交成功后由后端应用到全部小组
- 提交后锁定
- 提交后不可修改
- 再次提交返回 `409`

说明：

- 业务要求上应在赛前完成该动作。
- 首版接口层只保留“未提交前可提交、已提交后不可再改”的最小硬约束，不额外电子化主持流程控制。

响应字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `submitted` | `bool` | 是否已提交 |
| `appliedGroupCount` | `int` | 本次应用到的小组数量 |
| `submittedAt` | `string` | 提交时间，RFC3339 |
| `submitterName` | `string` | 提交管理员名称 |

#### 6.5.6 异常解锁某组某年

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-control/unlock-year`
- 权限：`sandbox-game:admin-control:unlock-year`

Go DTO 建议：

- 请求：`UnlockYearReq`
- 响应：`UnlockYearResp`

请求体建议：

```json
{
  "groupId": 1,
  "yearNo": 2,
  "unlockTargetType": "OPERATING",
  "targetStageCode": "Q2",
  "reason": "现场核对后发现 Q2 数据需修正"
}
```

字段建议：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `groupId` | `int64` | 是 | 目标组 |
| `yearNo` | `int` | 是 | 目标年份 |
| `unlockTargetType` | `string` | 是 | `OPERATING` / `REPORT` |
| `targetStageCode` | `string` | 条件必填 | 当 `unlockTargetType = OPERATING` 时必填，允许值：`Q1 / Q2 / Q3 / Q4 / YEAR_END` |
| `reason` | `string` | 是 | 管理员填写的异常解锁原因 |

规则：

- 仅管理员可执行。
- 仅允许在 `下一年尚未开放前` 执行。
- 前端不做复杂可解锁预判，管理员点击后直接提交；服务端负责最终硬校验。
- 仅允许对“已正式提交”的目标阶段或财报结果执行异常解锁。
- 当 `unlockTargetType = OPERATING` 时，服务端不得仅因财报页仍处于可编辑状态，就返回“当前年份已处于可编辑状态”。
- 当 `unlockTargetType = REPORT` 时，不回退经营页已生效阶段。
- 经营页异常解锁后：
  - 早于目标阶段的经营结果继续有效并保持只读
  - 目标阶段恢复可编辑
  - 晚于目标阶段的经营结果失效但保留原值，标记为 `失效草稿`
  - 财报结果同步失效，财报手工值保留为 `失效草稿`
- 财报页异常解锁后：
  - 经营页已生效结果保持不变
  - 财报结果失效并恢复可编辑
  - 财报手工值保留为 `失效草稿`
- 若该年已计入汇总，则汇总必须立即失效。
- 若该年结果曾触发该组破产，且破产依据来自本次被失效的结果，则解锁成功后应临时恢复为 `NORMAL`，待重新提交后再重新判定。
- 必须记录异常解锁日志与管理员动作日志。

响应字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `groupId` | `int64` | 目标组 |
| `yearNo` | `int` | 目标年份 |
| `unlockTargetType` | `string` | 本次解锁目标类型 |
| `targetStageCode` | `string` | 本次回退阶段；当目标为 `REPORT` 时可为空 |
| `yearStatus` | `string` | 解锁后年份主状态 |
| `editableStageCode` | `string` | 当前重新开放可编辑的经营阶段；财报页解锁时可为空 |
| `reportStatus` | `string` | 解锁后财报状态 |
| `summaryEffective` | `bool` | 是否仍计入正式汇总 |
| `businessStatus` | `string` | 解锁后经营状态 |
| `unlockLogId` | `int64` | 解锁日志 ID |

#### 6.5.7 后端实现规则清单

1. DTO 与枚举
- `UnlockYearReq` 必须新增 `unlockTargetType`、`targetStageCode` 字段。
- 服务端应显式定义 `UnlockTargetType` 与 `StageCode` 枚举，避免字符串散落在业务代码中。

2. 服务端校验
- 校验目标组、目标年份存在。
- 校验下一年尚未开放。
- 校验 `reason` 非空。
- 校验 `unlockTargetType` 合法。
- 当目标为 `OPERATING` 时，校验 `targetStageCode` 合法且该阶段已正式提交。
- 当目标为 `REPORT` 时，校验财报结果已正式提交。

3. 经营页异常解锁执行规则
- 按目标阶段回退有效状态，而不是回退到“当前最新阶段”。
- 目标阶段之前的有效经营结果继续保留。
- 目标阶段之后的经营结果与财报结果转为 `失效草稿`。
- `失效草稿` 原值必须保留，不允许直接清空。

4. 财报页异常解锁执行规则
- 不改变经营页已生效阶段。
- 财报结果转为 `失效草稿` 并恢复为可编辑。
- 汇总结果同步失效。

5. 结果口径与关联处理
- 汇总、破产判定、下一年结转只读取 `有效结果`，不读取 `失效草稿`。
- 当异常解锁导致原破产依据失效时，应临时恢复 `NORMAL`。
- 重新提交成功后，再重新生成正式汇总与正式状态。

6. 审计与日志
- `sg_admin_unlock_log` 应补充或正式使用：目标类型、目标阶段、原因、解锁前状态、解锁后状态。
- `sg_admin_action_log` 继续记录操作人、操作时间与动作摘要。

7. 错误码与提示
- 错误提示必须目标化，不再只返回通用“当前年份已处于可编辑状态”。
- 推荐至少区分：经营页无需解锁、财报页无需解锁、目标阶段未提交、财报未提交、下一年已开放、目标类型非法、阶段非法。

8. 测试覆盖
- 单元测试覆盖：`OPERATING/Q1`、`OPERATING/Q2`、`OPERATING/YEAR_END`、`REPORT` 四类主场景。
- 集成测试覆盖：失效草稿保留、汇总失效、破产恢复、重新提交后重新生效。

---

### 6.6 `admin-group-data`

#### 6.6.1 查看任意组经营页视图

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-group-data/get-operating-view?groupId=1&yearNo=1`
- 权限：`sandbox-game:admin-group-data:query`

用途：

- 管理员查看指定组指定年的经营页数据与状态
- 返回结构可复用玩家端 `get-year-view` 的数据模型

#### 6.6.2 查看任意组财报页视图

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-group-data/get-report-view?groupId=1&yearNo=1`
- 权限：`sandbox-game:admin-group-data:query`

#### 6.6.3 分页查询阶段提交日志

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-group-data/page-stage-submissions`
- 权限：`sandbox-game:admin-group-data:query`

查询参数建议：

- `groupId`
- `yearNo`
- `stageCode`
- `pageNo`
- `pageSize`

#### 6.6.4 分页查询财报提交日志

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-group-data/page-report-submissions`
- 权限：`sandbox-game:admin-group-data:query`

---

### 6.7 `audit-log`

#### 6.7.1 分页查询异常解锁日志

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/audit-log/page-unlock-log`
- 权限：`sandbox-game:audit-log:query`

#### 6.7.2 分页查询管理员动作日志

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/audit-log/page-admin-action-log`
- 权限：`sandbox-game:audit-log:query`

用途：

- 查看开放下一年、更新最终年份、提交初始基线等动作

---

## 7. 关键请求对象定义

### 7.1 `OperatingPayload`

说明：

- 表示经营页业务区块数据
- 具体内部键名已按 `1组 最终版.xlsx` 和 `docs/excel_field_mapping.md` 冻结首版基线；后续若只调整页面显示名，不改 payload 结构
- 接口层将其视为结构化 JSON 对象

### 7.2 `ReportManualPayload`

说明：

- 表示财报页绿色手工项集合
- 当前至少覆盖：在制品、成品、材料、所得税税率
- 税率字段应限制在：`0.25 / 0.15 / 0`

### 7.3 `StateSnapshot`

建议字段：

- `yearStatus`
- `stageStatus`
- `reportStatus`
- `businessStatus`
- `summaryEffective`

### 7.4 管理端控制链 DTO 草图（Go）

建议集中放在：`internal/http/dto/admin_control_dto.go`

```go
type AdminControlConfigResp struct {
    FinalYear                    int                     `json:"finalYear"`
    CurrentOpenYear              int                     `json:"currentOpenYear"`
    CanOpenNextYear              bool                    `json:"canOpenNextYear"`
    NextOpenableYear             int                     `json:"nextOpenableYear"`
    OpenNextYearBlockedReason    string                  `json:"openNextYearBlockedReason"`
    RuleVersion                  string                  `json:"ruleVersion"`
    TemplateVersion              string                  `json:"templateVersion"`
    InitialBaselineSubmitted     bool                    `json:"initialBaselineSubmitted"`
    InitialBaselineSubmittedAt   *string                 `json:"initialBaselineSubmittedAt"`
    InitialBaselineSubmitterName *string                 `json:"initialBaselineSubmitterName"`
    LatestAdminAction            *AdminActionSummaryResp `json:"latestAdminAction"`
}

type UpdateFinalYearReq struct {
    FinalYear int `json:"finalYear" validate:"required,min=1"`
}

type UpdateFinalYearResp struct {
    FinalYear       int    `json:"finalYear"`
    CurrentOpenYear int    `json:"currentOpenYear"`
    UpdatedAt       string `json:"updatedAt"`
    UpdatedBy       string `json:"updatedBy"`
}

type OpenNextYearReq struct {
    TargetYearNo int `json:"targetYearNo" validate:"required,min=1"`
}

type OpenNextYearResp struct {
    PreviousOpenYear          int                     `json:"previousOpenYear"`
    CurrentOpenYear           int                     `json:"currentOpenYear"`
    OpenedYearNo              int                     `json:"openedYearNo"`
    FinalYear                 int                     `json:"finalYear"`
    CanOpenNextYear           bool                    `json:"canOpenNextYear"`
    OpenNextYearBlockedReason string                  `json:"openNextYearBlockedReason"`
    LatestAdminAction         *AdminActionSummaryResp `json:"latestAdminAction"`
}

type InitialBaselineResp struct {
    Submitted         bool           `json:"submitted"`
    Editable          bool           `json:"editable"`
    BaselinePayload   map[string]any `json:"baselinePayload"`
    AppliedGroupCount int            `json:"appliedGroupCount"`
    SubmitterName     *string        `json:"submitterName"`
    SubmittedAt       *string        `json:"submittedAt"`
}

type SubmitInitialBaselineReq struct {
    BaselinePayload map[string]any `json:"baselinePayload" validate:"required"`
}

type SubmitInitialBaselineResp struct {
    Submitted         bool   `json:"submitted"`
    AppliedGroupCount int    `json:"appliedGroupCount"`
    SubmitterName     string `json:"submitterName"`
    SubmittedAt       string `json:"submittedAt"`
}

type UnlockYearReq struct {
    GroupID int64  `json:"groupId" validate:"required,min=1"`
    YearNo  int    `json:"yearNo" validate:"required,min=0"`
    Reason  string `json:"reason" validate:"required,max=500"`
}

type UnlockYearResp struct {
    GroupID          int64  `json:"groupId"`
    YearNo           int    `json:"yearNo"`
    YearStatus       string `json:"yearStatus"`
    StageStatus      string `json:"stageStatus"`
    ReportStatus     string `json:"reportStatus"`
    SummaryEffective bool   `json:"summaryEffective"`
    BusinessStatus   string `json:"businessStatus"`
    UnlockLogID      int64  `json:"unlockLogId"`
}

type AdminActionSummaryResp struct {
    ActionCode   string `json:"actionCode"`
    ActionName   string `json:"actionName"`
    OperatorName string `json:"operatorName"`
    OperateTime  string `json:"operateTime"`
}
```

说明：

- `Payload` 只承载业务字段，不暴露 Excel 单元格坐标。
- 时间字段统一返回 RFC3339 字符串。
- `Reason` 建议服务端先做 `trim` 再校验长度。

### 7.5 管理端控制链错误常量建议（Go）

建议集中放在：`internal/common/errors/admin_control_errors.go`

| Go 常量名 | 默认 HTTP | 触发场景 | 默认提示语建议 |
|---|---|---|---|
| `ErrAdminControlConfigNotFound` | `404` | 配置表未初始化 | `游戏配置不存在` |
| `ErrAdminControlFinalYearTooSmall` | `422` | `finalYear < currentOpenYear` | `最终年份不能小于当前开放年份` |
| `ErrAdminControlTargetYearMismatch` | `409` | `targetYearNo != currentOpenYear + 1` | `开放年份与当前状态不一致` |
| `ErrAdminControlFinalYearReached` | `409` | 已到最终年份仍尝试开放 | `已达到最终年份，无法继续开放` |
| `ErrAdminControlOpenNextYearBlocked` | `422` | 未满足开放下一年条件 | `当前仍有未完成财报的小组` |
| `ErrAdminControlInitialBaselineSubmitted` | `409` | 初始基线重复提交 | `初始基线已提交，不能重复提交` |
| `ErrAdminControlInitialBaselineInvalid` | `422` | 初始基线载荷缺失或结构非法 | `初始基线数据不完整` |
| `ErrAdminControlUnlockTargetNotFound` | `404` | 目标组或目标年份不存在 | `未找到需要解锁的目标数据` |
| `ErrAdminControlUnlockNextYearOpened` | `409` | 下一年已开放后仍尝试解锁 | `下一年已开放，不能再解锁本年` |
| `ErrAdminControlUnlockTargetTypeRequired` | `422` | 未选择解锁目标类型 | `请选择解锁目标` |
| `ErrAdminControlUnlockTargetTypeInvalid` | `422` | 解锁目标类型非法 | `异常解锁目标类型不合法` |
| `ErrAdminControlUnlockStageRequired` | `422` | 经营页异常解锁未选择阶段 | `请选择要回退的经营阶段` |
| `ErrAdminControlUnlockStageInvalid` | `422` | 回退阶段非法 | `经营回退阶段不合法` |
| `ErrAdminControlUnlockTargetNotSubmitted` | `409` | 目标阶段或财报尚未正式提交 | `目标尚未正式提交，不能解锁` |
| `ErrAdminControlOperatingAlreadyEditable` | `409` | 经营页当前无需解锁 | `经营页当前无需解锁` |
| `ErrAdminControlReportAlreadyEditable` | `409` | 财报页当前无需解锁 | `财报页当前无需解锁` |
| `ErrAdminControlUnlockReasonRequired` | `422` | 未填写解锁原因 | `请填写异常解锁原因` |

---

## 8. 重复提交与并发处理规则

### 8.1 草稿保存

- 采用最后写入覆盖当前草稿的策略
- 不要求乐观锁拦截首版自动保存

### 8.2 正式提交

- 阶段提交、财报提交必须以当前服务端状态为准
- 服务端发现状态已变化时，返回 `409`
- 服务端禁止以客户端传来的状态字段直接替代服务端状态

### 8.3 管理动作

- `open-next-year`
- `submit-initial-baseline`
- `unlock-year`

以上动作必须记录：

- 操作人
- 操作时间
- 操作前状态
- 操作后状态
- 关键原因 / 备注

---

## 9. 测试口径要求

按正式版接口设计，后续 `.http` / API 测试至少覆盖：

- 正常场景：
  - 获取经营页
  - 保存草稿
  - 阶段提交成功
  - 财报提交成功
  - 管理员开放下一年成功
- 参数异常：
  - 缺少 `yearNo`
  - 非法 `stageCode`
  - 税率非法
- 权限异常：
  - 玩家访问管理员接口
  - 管理员访问无权限域
- 状态迁移异常：
  - 当前阶段必填项未完成
  - 财报平衡失败
  - 已完成年份重复提交财报
  - 下一年已开放后再解锁上一年

---

## 10. 当前待后续补齐但不影响本版成立的内容

以下内容后续可补充，不影响本接口设计正式版成立：

- `OperatingPayload` 的最终字段映射清单
- `ReportManualPayload` 的最终字段命名说明
- 详细接口示例 JSON 样例库
- `.http` 用例文件
- 异常解锁日志的高级筛选条件与导出策略



