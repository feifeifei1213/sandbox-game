# 沙盘经营系统接口设计文档（正式版）

> 更新日期：2026-05-20
> 适用方式：基于《正式需求文档（首版）》《最小状态机 v0.1》《技术选型细化文档（Go 方向） v0.1》，定义首版正式业务接口边界，作为后续 Go 后端开发、前端 API 客户端开发和接口测试的统一依据。  
> 文档定位：本文件定义接口域划分、路径风格、请求/响应结构、核心动作语义、关键错误状态和测试口径。  
> 说明：本文件已按 `1组 最终版.xlsx` 口径收口；接口域划分、动作语义和结构性字段已经冻结，若后续仅有页面标签细修，应优先更新映射文档，不直接改动接口结构。

## 1. 适用范围

本规范适用于沙盘经营系统全部首版业务接口：

- 首版最小登录接口
- 玩家端经营页接口
- 玩家端年度订单接口
- 玩家端财报页接口
- 管理员订单管理接口
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
- `sandbox-game:player-order:query`
- `sandbox-game:player-order:submit-market-investment`
- `sandbox-game:player-order:select-order`
- `sandbox-game:player-report:submit`
- `sandbox-game:admin-order:query`
- `sandbox-game:admin-order:update-config`
- `sandbox-game:admin-order:generate-pool`
- `sandbox-game:admin-order:control-bidding`
- `sandbox-game:admin-control:query-setup`
- `sandbox-game:admin-control:initialize-game`
- `sandbox-game:admin-control:open-next-year`
- `sandbox-game:admin-control:unlock-year`
- `sandbox-game:admin-dictionary:query`
- `sandbox-game:admin-dictionary:manage-scheme`
- `sandbox-game:admin-dictionary:update-current`
- `sandbox-game:admin-rollback:query`
- `sandbox-game:admin-rollback:create-snapshot`
- `sandbox-game:admin-rollback:restore-group-snapshot`
- `sandbox-game:admin-rollback:unlock-retry`
- `sandbox-game:admin-notice:query`
- `sandbox-game:admin-notice:send-general`
- `sandbox-game:admin-notice:send-adjustment`
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
  - 回退或提交过程中状态已变化，需要刷新后重试
- `422` 适用于：
  - 当前阶段必填项未完成
  - 财报平衡校验失败
  - 当前年份未开放
  - 当前阶段未到可提交时机
  - 破产组继续提交
  - 存在回退补提 / 待重提小组时尝试开放下一年

---

## 4. 业务接口域划分

首版正式接口域建议如下：

| 领域 | 说明 |
|---|---|
| `auth` | 登录、当前用户信息、退出登录 |
| `game-config` | 当前游戏配置、开放年份、规则版本信息 |
| `player-order` | 玩家年度订单页、当前订单模板完整市场投入、按轮次选择订单与交付状态查看 |
| `player-operating` | 玩家经营页读取、草稿保存、阶段提交 |
| `player-report` | 财报页读取、草稿保存、财报提交 |
| `admin-order` | 管理员多年订单数量控制台、年度市场开启、预览/确认订单池、标段释放与竞标控制 |
| `admin-summary` | 汇总页、最终排名 |
| `admin-control` | 最终年份设置、开放下一年、初始基线、异常解锁兼容旧接口 |
| `admin-dictionary` | 赛前业务显示字典方案、当前比赛字典快照、应用方案与修改日志 |
| `admin-rollback` | 管理员回退与修正、退回重提、快照列表、手动快照、单组快照恢复 |
| `admin-group-data` | 管理员查看任意组任一年经营/财报数据 |
| `audit-log` | 提交日志、解锁日志、管理员动作日志 |

说明：

- 首版不建议按“前端页面路径”来拆接口域。
- 首版也不建议把所有接口塞进单一 `game` 域，避免后续维护混乱。

### 4.1 沙盘版本包接口口径

- 沙盘版本包为系统内置配置，不提供管理员在线新增、编辑或删除版本包的接口。
- 管理员只能在比赛初始化前从服务端返回的版本包列表中选择本场比赛版本。
- 初始化完成后，版本包锁定；`game-config`、玩家经营页、玩家财报页、年度订单页等接口应返回当前版本标识和当前比赛字典版本，供前端选择字段模板并叠加业务显示字典。
- 当前已内置 `VIP_SERVICE_V1`，绑定贵宾服务版字段模板与通用公式/流程规则。
- 下一步新增 `PRODUCTION_V1`，绑定生产制造版经营页和财报页字段模板，公式规则和流程规则继续共用通用版本；订单字段模板暂时仍绑定 `VIP_ORDER_TEMPLATE_V1`，待生产版订单字段确认后再升级为生产版订单模板。
- 版本包负责字段结构、公式版本和流程规则；业务显示字典只负责显示名称，不改变稳定字段编码、payload、数据库字段、公式或流程。
- 业务显示字典方案必须绑定版本包；当前比赛初始化时复制一份字典快照，后续比赛中修改的是当前比赛字典快照。

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
- 市场编码：
  - `LOCAL`
  - `REGIONAL`
  - `NATIONAL`
  - `GLOBAL`
- 订单类型：
  - `AGENCY_INSPECTION`
  - `TWO_CABIN_VIP`
  - `BUSINESS_VIP`
  - `MEMBER_CUSTOM`
- 订单竞标状态：
  - `NOT_REQUIRED`
  - `NOT_GENERATED`
  - `PREVIEW_GENERATED`
  - `POOL_CONFIRMED`
  - `WAITING_INVESTMENT`
  - `INVESTMENT_READY`
  - `SEQUENCE_READY`
  - `SELECTING`
  - `COMPLETED`
  - `SKIPPED`
- 标段状态：
  - `MARKET_DISABLED`
  - `NO_ORDER_CONFIG`
  - `WAITING_RELEASE`
  - `SELECTING`
  - `COMPLETED`
  - `SKIPPED`
- 订单交付状态：
  - `SELECTED`
  - `DELIVERED`
  - `UNFINISHED`

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
| `defaultRoute` | 默认跳转地址；玩家为 `player/operating?yearNo=0`，管理员按初始化状态返回 `admin/setup` 或 `admin/summary` |

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
| `editionCode` | 当前沙盘版本包编码，例如 `VIP_SERVICE_V1` |
| `editionName` | 当前沙盘版本包显示名，例如 `贵宾服务版 V1` |
| `ruleVersion` | 当前规则版本 |
| `templateVersion` | 当前模板版本 |
| `operatingTemplateVersion` | 当前经营页字段模板版本 |
| `reportTemplateVersion` | 当前财报页字段模板版本 |
| `orderTemplateVersion` | 当前订单字段模板版本 |
| `processRuleVersion` | 当前流程规则版本 |
| `demoYearEnabled` | 是否启用 `0年` 引导年，首版固定为 `true` |

#### 6.1.1A 获取可选沙盘版本包

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/game-config/list-game-editions`
- 权限：管理员已登录可访问

返回字段建议：

| 字段 | 说明 |
|---|---|
| `editionCode` | 版本包编码 |
| `editionName` | 版本包显示名 |
| `description` | 版本说明 |
| `defaultEdition` | 是否默认版本 |
| `operatingTemplateVersion` | 经营页字段模板版本 |
| `reportTemplateVersion` | 财报页字段模板版本 |
| `orderTemplateVersion` | 订单字段模板版本 |
| `formulaVersion` | 公式规则版本 |
| `processRuleVersion` | 流程规则版本 |

规则：

- 当前至少返回 `VIP_SERVICE_V1`；新增生产制造版后应同时返回 `PRODUCTION_V1`。
- 本接口只返回系统内置版本包，不支持页面新增或修改版本包。
- 比赛已初始化后仍可查询，但不能再用于切换当前比赛版本。

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

正式年份经营页 `Q1` 提交前，服务端必须校验本组本年订单前置流程已完成。`0年` 不做订单前置校验。

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
- 若为正式年份 `Q1` 提交，本组本年订单前置流程必须已完成

成功结果：

- 记录阶段提交日志
- 推进到下一经营阶段，或进入 `REPORT_PENDING`
- 返回最新状态摘要
- 返回本次系统计算的期末现金，供玩家线下核对

重复提交规则：

- 同一阶段已正式提交后再次提交，返回 `409`

---

### 6.2A `player-order`

#### 6.2A.1 获取年度订单视图

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/player-order/get-year-view?yearNo=1`
- 权限：`sandbox-game:player-order:query`

返回字段建议：

| 字段 | 说明 |
|---|---|
| `groupId` | 当前玩家所属组 |
| `yearNo` | 年份 |
| `orderRequired` | 是否需要订单；`0年` 为 `false` |
| `markets[].marketCode` | 市场编码 |
| `markets[].marketName` | 市场名称 |
| `markets[].marketEnabled` | 当前年份该市场是否开启；未开启市场不生成订单、不抢单 |
| `markets[].marketInvestmentLimit` | 当前年份该市场投入上限；`null` 表示无上限，未开启市场前端显示为不可投 |
| `investmentStatus` | 本组当年当前订单模板全部市场投入提交状态 |
| `canSubmitInvestment` | 是否可提交本年市场投入 |
| `markets[].segments[].marketInvestment` | 本组该标段投入 |
| `markets[].segments[].investmentSubmitted` | 本组该标段投入是否已提交 |
| `markets[].isMarketLeader` | 本组是否为该市场本年市场龙头 |
| `markets[].segments[].orderType` | 标段订单类型 |
| `markets[].segments[].releaseSequenceNo` | 标段释放顺序 |
| `markets[].segments[].segmentStatus` | 标段状态 |
| `markets[].segments[].currentRoundNo` | 当前或最近处理轮次 |
| `markets[].segments[].nextRoundNo` | 下一待开启轮次；没有则为空 |
| `markets[].segments[].selfRoundStatus` | 本组当前轮次状态；无顺序记录时派生为“本轮无资格” |
| `markets[].segments[].selfSelectionSequenceNo` | 本组当前轮次顺序；无资格时为空 |
| `markets[].segments[].canSelectOrder` | 当前轮次是否轮到本组选择 |
| `markets[].segments[].ordersVisible` | 玩家端是否可以查看该标段订单明细；由后端按标段是否已释放进入过选单阶段计算 |
| `markets[].segments[].availableOrders` | 当前标段仍可选择订单列表；仅 `ordersVisible=true` 时返回明细，否则返回空数组 |
| `markets[].segments[].lockedOrders[]` | 已被选择订单的只读展示信息，玩家端只用于灰色不可选，不返回选中组；仅 `ordersVisible=true` 时返回明细，否则返回空数组 |
| `markets[].segments[].selectedOrders[]` | 本组该标段各轮已选订单，包含 `roundNo`；仅 `ordersVisible=true` 时返回明细，否则返回空数组 |
| `pollingIntervalSeconds` | 年度订单页建议自动轮询间隔，首版为 `3` |

规则：

- `0年` 返回 `orderRequired=false`，不进入市场选单。
- 未开启市场仍返回其 4 个订单类型投入项，但玩家端输入框禁用，提交时系统自动带 `0`；若绕过前端提交非 `0`，服务端返回 `422`。
- 玩家端应展示每个市场的 `marketInvestmentLimit`；`null` 显示为“无上限”，未开启市场显示为“未开启”。
- 市场投入阶段、等待其他小组提交投入、已生成选单顺序但标段尚未释放时，`ordersVisible=false`，玩家接口不得返回具体订单池明细、订单编号、金额、数量、单价、账期或交付面板所需的订单明细。
- 管理员释放标段并进入过选单阶段后，`ordersVisible=true`；前端按现有页面结构分别展示订单卡片、本组已选订单和交付面板，不新增合并后的统一大模块。
- 玩家不返回其他组已选订单明细。
- 已被选择的订单在玩家端显示为灰色不可选，但不返回被哪个小组选走。
- 只有当前释放到的标段才允许选择订单。
- 玩家端只展示本组状态和当前轮次，不返回其他小组投入、完整顺序或未来轮次参与名单。
- 首版通过自动轮询同步状态，不做 WebSocket。
- 轮询采用静默局部合并，不进入整页加载态，不改变页面滚动位置。

#### 6.2A.2 提交市场投入

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/player-order/submit-market-investments`
- 权限：`sandbox-game:player-order:submit-market-investment`

请求体建议：

```json
{
  "yearNo": 1,
  "investments": [
    {
      "marketCode": "LOCAL",
      "orderType": "AGENCY_INSPECTION",
      "marketInvestment": 100
    },
    {
      "marketCode": "LOCAL",
      "orderType": "TWO_CABIN_VIP",
      "marketInvestment": 0
    }
  ]
}
```

规则：

- 仅正式年份允许提交。
- `0年` 不走独立市场投入提交。
- 每次提交必须包含本年当前订单模板定义的全部 `市场 + 订单类型` 投入值。
- 投入金额必须大于等于 `0` 且必须为整数；空值和小数不允许提交。
- 若某市场未开启，该市场下 4 项 `marketInvestment` 必须全部为 `0`；前端应自动带 `0`，绕过前端提交非 `0` 返回 `422`。
- 若某市场配置了 `marketInvestmentLimit`，该市场下 4 项投入合计不得超过该上限；未配置上限时不做单市场上限校验。
- 提交后不可修改；重复提交返回 `409`。
- 任何小组某标段投入为 `0` 时均不参与该标段选单，市场龙头也不例外。
- 市场投入提交后回写经营页年初市场投入区域为只读展示，经营页 `Q1` 不再允许修改。

#### 6.2A.3 选择订单

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/player-order/select-order`
- 权限：`sandbox-game:player-order:select-order`

请求体建议：

```json
{
  "yearNo": 1,
  "marketCode": "LOCAL",
  "orderType": "AGENCY_INSPECTION",
  "orderId": 10001
}
```

规则：

- 仅在该标段处于 `SELECTING` 时允许选择。
- 必须轮到当前组。
- 当前组必须具有服务端当前轮次资格；客户端不提交权威 `roundNo`，后端从加锁后的标段状态读取当前轮次。
- 每组每标段每轮最多选择 `1` 个订单；贵宾/生产最多四轮，机场暂时一轮。
- 订单必须属于当前年份、市场和订单类型，且状态仍可选。
- 选择成功后订单锁定，不再对其他组可选。
- 玩家选择后不可撤销。

返回字段建议：

| 字段 | 说明 |
|---|---|
| `selectedOrderId` | 已选订单 ID |
| `marketCode` | 市场 |
| `orderType` | 订单类型 |
| `roundNo` | 实际选择轮次 |
| `selectionSequenceNo` | 本组顺序 |
| `nextGroupId` | 本轮下一顺位组；若本轮已结束则为空 |
| `segmentStatus` | 选择后的标段状态 |

#### 6.2A.4 放弃当前轮

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/player-order/pass-segment`
- 权限：`sandbox-game:player-order:select-order`

请求体建议：

```json
{
  "yearNo": 1,
  "marketCode": "LOCAL",
  "orderType": "AGENCY_INSPECTION"
}
```

规则：

- 仅当前轮到本组时允许放弃。
- 放弃只对当前轮生效，不影响后续有资格轮次或后续标段。
- 前端放弃前二次确认但不要求原因；放弃后本轮不能反悔，系统推进到本轮下一个有资格小组。
- 客户端不提交权威 `roundNo`，后端根据当前标段状态确定轮次。

#### 6.2A.5 交付已选订单

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/player-order/deliver-orders`
- 权限：`sandbox-game:player-order:deliver-order`

请求体建议：

```json
{
  "yearNo": 1,
  "stageCode": "Q1",
  "orderIds": [10001, 10002]
}
```

规则：

- 只能交付本组本年已选且未交付订单。
- 单个订单不能拆分到多个季度交付；一个季度可以交付多个完整订单。
- `stageCode` 必须等于当前经营季度。
- 服务端校验当前季度销售收入等于本次交付订单金额合计，校验通过后订单状态变为 `DELIVERED`。
- 年末仍未交付订单标记为 `UNFINISHED`，首版不自动处罚。

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
| `reportComputedPayload.reportTotalEquity` | 当前年份所有者权益；最佳 CEO 得分新增项的计算基数 |
| `reportComputedPayload.reportBestSalesDirectorScore` | 截至当前年份累计订单总额除以 `10` 后的最佳销售总监得分 |
| `reportComputedPayload.reportBestCeoScore` | 五项总监最终得分合计再加当前年份所有者权益的 `1/2` |
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

#### 6.3.4 总监得分计算与返回口径

- 本次公式调整不新增接口字段，不修改请求结构，也不需要数据库迁移；继续复用 `reportComputedPayload` 中现有字段。
- 最佳销售总监得分按以下业务公式返回：

```text
reportBestSalesDirectorScore_y =
  (orderTotal_0 + orderTotal_1 + ... + orderTotal_y) / 10
```

- 服务端按年递推时应实现为：

```text
previousReport.reportBestSalesDirectorScore + currentYear.orderTotal / 10
```

- 不得实现为 `(previousReport.reportBestSalesDirectorScore + currentYear.orderTotal) / 10`，避免已经缩放的历史得分被再次除以 `10`。
- 最佳 CEO 得分按以下公式返回：

```text
reportBestCeoScore =
  reportBestMarketDirectorScore
  + reportBestTechnologyDirectorScore
  + productionHumanScore
  + reportBestSalesDirectorScore
  + reportBestCfoScore
  + reportTotalEquity / 2
```

- `reportTotalEquity` 为当前年份财报权益；允许为负，负值会降低 CEO 得分。两个得分字段均不额外取整，允许小数。
- 新公式从 `0年` 起适用于生产制造版和贵宾服务版。获取历史财报视图时，服务端应按当前规则从 `0年` 递归重建总监得分后返回，不要求物理改写旧提交 JSON 或旧快照。
- 保存草稿与提交财报均由后端按当前公式生成计算结果；提交结果仍是正式权威口径。
- 前端财报实时预览可复用返回值中的 `reportBestSalesDirectorScore`，不得再次执行 `/ 10`；前端必须用当前预览得到的 `reportTotalEquity` 重新计算 CEO 的权益项，避免编辑态与提交态显示不一致。

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
| `list[].businessStatus` | 小组当前经营状态，来源于小组主数据当前 `businessStatus`，不使用某年汇总快照中的历史状态 |
| `list[].ranking` | 若为最终年份，可返回最终排名 |

规则：

- `0年` 不进入正式汇总
- 仅 `COMPLETED` 的正式年份计入汇总
- 被异常解锁且未重提财报的年份，不得出现在正式汇总口径中
- 年度收入、利润、权益来自对应年份正式汇总快照；`businessStatus` 用于页面行级状态展示，必须反映小组当前经营状态。若小组当前已破产，即使历史年度快照状态为 `NORMAL`，接口也返回 `BANKRUPT`。

#### 6.4.2 获取最终排名

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-summary/get-final-ranking`
- 权限：`sandbox-game:admin-summary:query`

规则：

- 仅在最终年份结果可用时返回有效排名
- 排名依据：最终年份 `equity` 倒序
- 最终排名区的 `businessStatus` 同样显示小组当前经营状态，不显示最终年份快照中的历史状态。

---

### 6.4A `admin-order`

#### 6.4A.0 获取市场预测

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/player-order/get-market-forecast`
- 权限：`sandbox-game:player-order:query`

返回：

- 固定按 `1~3年 / 4~5年 / 6~8年` 三段返回市场预测。
- 每段包含四个市场、四类产品在各年份的预测金额/订单量摘要。
- 返回结构需能支撑玩家端按 Excel 竖状柱状图展示：阶段下按市场分图，横轴为年份，柱子区分四类订单。
- 返回管理员维护或导入的预测说明文字。
- 市场预测不因管理员配置的最终年份而裁剪；例如最终年份为 `5年` 时，仍可展示 `6~8年` 的预测趋势。

规则：

- 市场预测由多年订单数量控制台和订单生成公式链产生，玩家端只读展示。
- 市场预测用于玩家决策市场投入，不代表市场当年自动开启。

#### 6.4A.1 获取年度订单管理视图

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-order/get-generation-console?yearNo=1`
- 权限：`sandbox-game:admin-order:query`

返回：

- 返回当年从多年订单数量控制台带入的 `市场 + 订单类型` 订单卡片数量、标段释放顺序、订单池批次状态和风险提示。
- 返回当年 4 个市场开启状态；本地市场默认开启，区域市场、全国市场、全球市场默认关闭。
- 返回当年生成状态：`NOT_GENERATED / PREVIEW_GENERATED / POOL_CONFIRMED / SELECTING / COMPLETED`。
- 返回预览批次或正式批次摘要：批次 ID、公式版本、生成时间、确认时间。
- 返回市场预测快照摘要：预测版本、最近更新时间、是否已按当前控制台生成。
- 首版不返回均价、波动系数、最小/最大数量等复杂参数编辑项。

规则：

- 订单生成依据为 `道具-订单推算（服务企业）.xlsx` 的公式链。
- 多年订单数量控制台是订单数量主来源，当年订单管理只读取当前年份的控制台数量，不再维护另一套独立数量。
- 系统内置公式链，不把 Excel 文件作为运行时订单池上传结果。
- 后端生成订单时自动保证 `orderAmount = orderQuantity × unitPrice`；`orderQuantity` 为整数，`unitPrice` 可为小数，`orderAmount` 必须为整数，不新增管理员逐单调整金额或单价接口。
- “导入订单生成控制台参数/模板”可作为后续辅助能力，不是首版主链路。

#### 6.4A.1A 获取多年订单数量控制台

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-order/get-forecast-control`
- 权限：`sandbox-game:admin-order:query`

返回：

- 固定返回 `1年~8年 × 四个市场 × 四类产品` 的订单卡片数量控制台。
- 按 `1~3年 / 4~5年 / 6~8年` 三段组织，便于和 Excel 市场预测页一致。
- 返回每个市场、每个阶段的预测说明文字。
- 返回预测快照状态、最近保存时间、最近生成时间。

#### 6.4A.1C 更新多年订单数量控制台

- 方法：`PUT`
- 路径：`/api/v1/sandbox-game/admin-order/update-forecast-control`
- 权限：`sandbox-game:admin-order:update-config`

请求体建议：

```json
{
  "items": [
    {
      "yearNo": 1,
      "marketCode": "LOCAL",
      "orderType": "AGENCY_INSPECTION",
      "orderCount": 11
    }
  ],
  "narratives": [
    {
      "stageCode": "YEAR_1_3",
      "marketCode": "LOCAL",
      "content": "第一年：以代办过检为核心..."
    }
  ]
}
```

规则：

- `yearNo` 固定允许 `1~8`，不受管理员最终年份配置影响。
- `orderCount` 必须为非负整数，不设置固定业务最大值。
- 更新后应重新生成市场预测快照，或标记预测快照待生成。
- 若某年订单池已确认，修改多年控制台不得静默改变该年已确认订单池；后续是否允许重新生成需遵守订单池锁定规则。

#### 6.4A.1B 更新年度市场开启配置

- 方法：`PUT`
- 路径：`/api/v1/sandbox-game/admin-order/update-market-enabled-config`
- 权限：`sandbox-game:admin-order:update-config`

请求体建议：

```json
{
  "yearNo": 1,
  "markets": [
    { "marketCode": "LOCAL", "enabled": true, "marketInvestmentLimit": null },
    { "marketCode": "REGIONAL", "enabled": false, "marketInvestmentLimit": null },
    { "marketCode": "NATIONAL", "enabled": true, "marketInvestmentLimit": 80 },
    { "marketCode": "GLOBAL", "enabled": false, "marketInvestmentLimit": null }
  ]
}
```

规则：

- 只能在订单池确认前且本年尚无任何小组提交市场投入时更新；订单池确认后或已有任意小组提交市场投入后，本年市场开启状态和单市场投入上限锁定。
- `LOCAL` 默认开启，`REGIONAL / NATIONAL / GLOBAL` 默认关闭。
- 市场开启完全以管理员当年手动配置为准，不根据多年订单数量控制台中是否存在订单数量自动开启。
- 未开启市场下四个标段不生成订单池、不占用有效释放顺序、不进入选单。
- `marketInvestmentLimit` 为 `null` 表示无上限；非 `null` 时必须为非负整数，单位为 `M`。
- 关闭市场时上限不生效，前端应禁用上限输入。
- 若关闭市场时该市场已有未确认预览订单，应随重新生成预览流程覆盖或作废。

#### 6.4A.2 获取当年释放顺序配置

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-order/get-control-config?yearNo=1`
- 权限：`sandbox-game:admin-order:query`

返回：

- 按 `年份 + 市场 + 订单类型` 返回当年只读订单数量和标段释放顺序。
- 订单卡片数量来自多年订单数量控制台，在当前年度订单管理区只读展示，不在该接口中维护。
- 返回 `items[].marketEnabled` 与 `items[].marketInvestmentLimit`，用于前端将未开启市场行置灰、展示/编辑单市场投入上限；未开启市场的控制台数量只读保留，但不参与当年订单生成和释放顺序。
- 返回风险提示列表 `warnings[]`，用于提示订单数量不足、标段数量不足、某组可能没有可参与标段等情况。
- 首版不返回均价、波动系数、最小/最大数量等复杂参数。

#### 6.4A.3 更新当年释放顺序配置

- 方法：`PUT`
- 路径：`/api/v1/sandbox-game/admin-order/update-control-config`
- 权限：`sandbox-game:admin-order:update-config`

请求体建议：

```json
{
  "yearNo": 1,
  "items": [
    {
      "marketCode": "LOCAL",
      "orderType": "AGENCY_INSPECTION",
      "releaseSequenceNo": 1
    }
  ]
}
```

规则：

- 该接口只更新当年标段释放顺序，不更新订单数量。
- 订单数量必须通过多年订单数量控制台维护，并在本接口相关视图中只读带入。
- `releaseSequenceNo` 用于控制同一年内标段释放先后；同一年内不得重复。
- 标段释放顺序只对已开启且订单数量大于 `0` 的标段生效；未开启市场和订单数量为 `0` 的标段不占用有效释放顺序。
- 标段释放顺序在释放第一个有效标段前允许调整；释放第一个有效标段后不允许修改。
- 当年订单数量只做风险提示，不强制保底或自动补单；贵宾/生产多轮资格按固定规则生成。

#### 6.4A.4 生成预览订单池

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-order/generate-preview-pool`
- 权限：`sandbox-game:admin-order:generate-pool`

请求体建议：

```json
{
  "yearNo": 1
}
```

规则：

- 系统按多年订单数量控制台中的当年数量和内置 Excel 公式链生成预览订单池。
- 系统只读取当年、管理员已开启市场的数量；年度订单管理页中的数量为只读快照，不是另一套配置源。
- 系统只为已开启且订单数量大于 `0` 的标段生成订单。
- 每次生成预览时生成新的随机种子。
- 同一年度重复生成预览时，旧预览批次作废或覆盖。
- 订单池已确认后不允许重新生成预览。

#### 6.4A.4B 确认订单池

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-order/confirm-order-pool`
- 权限：`sandbox-game:admin-order:generate-pool`

请求体建议：

```json
{
  "yearNo": 1,
  "batchId": 1001
}
```

规则：

- 只能确认当前有效预览批次。
- 确认后保存正式批次、随机种子、公式版本、参数快照和订单明细。
- 确认后不允许修改影响该年订单池的控制台数量或重新生成订单池。

#### 6.4A.5 获取订单池

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-order/get-order-pool?yearNo=1&marketCode=ALL&orderType=ALL`
- 权限：`sandbox-game:admin-order:query`

用途：

- 管理员默认查看指定年份全部订单池，也可按 `marketCode`、`orderType` 筛选；筛选值省略或传 `ALL` 时表示全部。
- 返回订单业务编号或卡片编号字段，例如 `businessOrderNo` / `cardSequenceNo`，前端主列显示 `CARD-01` 等业务编号，不使用数据库自增 ID 作为主要展示编号。
- 管理员可查看订单市场、订单类型、金额、数量、单价、账期、当前状态与选中组。
- 返回的 `orderAmount` 应为整数金额，`unitPrice` 可为小数；两者与 `orderQuantity` 必须满足后端生成口径。

#### 6.4A.6 获取市场投入提交状态

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-order/get-investment-status?yearNo=1`
- 权限：`sandbox-game:admin-order:control-bidding`

返回：

- 各未破产小组是否已提交当年当前订单模板的全部市场投入。
- 已提交小组的提交时间。
- 未提交小组列表。

规则：

- 管理员只查看状态，不做代填。
- 未全部提交时，不允许开始开标。

#### 6.4A.7 生成标段选单顺序

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-order/generate-selection-sequence`
- 权限：`sandbox-game:admin-order:control-bidding`

请求体建议：

```json
{
  "yearNo": 1
}
```

规则：

- 前置条件：所有未破产小组已提交当年当前订单模板的全部市场投入、订单池已确认、标段释放顺序已配置。
- 系统按每个 `市场 + 订单类型` 标段只计算一次基础顺序，并一次生成全部有效轮次顺序。
- 未开启市场进入 `MARKET_DISABLED`，订单数量为 `0` 的已开启标段进入 `NO_ORDER_CONFIG`，两者不进入选单顺序。
- 订单数量大于 `0` 但所有未破产小组该标段投入均为 `0` 的标段进入 `SKIPPED`。
- `1年` 按当前标段投入排序，投入相同随机。
- `2年` 起有资格的市场龙头优先；市场龙头当前标段投入为 `0` 时无资格。
- 市场龙头已破产时，本年按没有有效市场龙头处理。
- 其余普通小组按当前标段投入排序；投入相同时按上一年度该市场订单总额排序；仍相同则随机。
- 随机结果、市场龙头和排序依据必须保存，便于追溯。
- 贵宾/生产按当前标段投入 `1~2 / 3~5 / 6~8 / >=9` 分别生成 `1 / 2 / 3 / 4` 轮；机场固定生成一轮。
- 后续轮次只过滤基础顺序，不重新排序或随机。
- 返回每个标段的 `theoreticalMaxSelections`、`availableOrderCount` 和 `warnings[]`；订单池不足时包含“订单池可能提前选空”，但不阻止生成。

#### 6.4A.8 获取年度选单状态

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-order/get-selection-status?yearNo=1`
- 权限：`sandbox-game:admin-order:query`

返回：

- 年度订单竞标状态
- 各小组市场投入提交状态
- 各市场龙头
- 各标段释放顺序
- 当前释放标段
- 当前标段 `currentRoundNo / nextRoundNo / completionReason`
- 当前标段第一至第四轮完整顺序，每轮包含小组、标段投入、上一年市场订单金额、龙头、状态和所选订单
- 当前标段各轮状态：待本轮开始、待选择、当前选择、已选择、已放弃、管理员跳过、破产无资格
- 当前轮到的小组
- 已选订单与未交付状态
- 订单池剩余数量、理论最大选单机会和风险提示

#### 6.4A.9 释放下一个标段

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-order/release-next-segment`
- 权限：`sandbox-game:admin-order:control-bidding`

请求体建议：

```json
{
  "yearNo": 1
}
```

规则：

- 系统按管理员配置的 `releaseSequenceNo` 找到下一个未完成标段。
- 只能释放预设顺序中的下一个标段，不能跳序。
- 当前标段处于 `SELECTING` 或 `ROUND_READY` 时，不允许释放下一个标段。
- 释放第一个标段后，该年标段释放顺序锁定。
- 订单数量为 `0` 或所有未破产小组该标段投入均为 `0` 的标段进入 `SKIPPED`。
- 同一时间建议只存在一个 `SELECTING` 标段，避免玩家并行选单造成现场混乱。
- 释放成功后立即把第一轮首个有效小组设为当前选择，不增加开启第一轮接口或确认弹窗。

#### 6.4A.9A 开启下一轮

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-order/open-next-round`
- 权限：`sandbox-game:admin-order:control-bidding`

请求体建议：

```json
{
  "yearNo": 1,
  "marketCode": "LOCAL",
  "orderType": "AGENCY_INSPECTION"
}
```

规则：

- 只允许当前标段处于 `ROUND_READY` 时调用。
- 后端从加锁后的标段状态确定下一轮，客户端不提交权威 `roundNo`。
- 下一轮没有有效参与小组时自动跳过；后续全部轮次均无参与小组时直接完成标段。
- 订单池已经选空时禁止开启并保持标段完成。
- 重复请求只有第一次成功，不能重复开启或跨轮。
- 管理员端按钮统一显示“开启下一轮”，不增加确认弹窗。

#### 6.4A.10 管理员代跳过当前小组

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-order/admin-skip-current-group`
- 权限：`sandbox-game:admin-order:control-bidding`

请求体建议：

```json
{
  "yearNo": 1,
  "marketCode": "LOCAL",
  "orderType": "AGENCY_INSPECTION",
  "groupId": 1,
  "reason": "现场超时未操作"
}
```

规则：

- 仅允许跳过当前轮当前小组。
- 只对当前轮生效，不影响后续有资格轮次或后续标段。
- 管理员不能代玩家选择订单。
- `reason` 必填，并记录管理员动作日志。

---

### 6.5 `admin-control`

#### 6.5.1 获取赛前初始化状态

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-control/get-setup-status`
- 权限：`sandbox-game:admin-control:query-setup`

Go DTO 建议：

- 响应：`AdminControlSetupStatusResp`

返回字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `initialized` | `bool` | 比赛是否已初始化 |
| `groupCount` | `int` | 当前已初始化的小组数量 |
| `finalYear` | `int` | 当前最终年份配置 |
| `currentOpenYear` | `int` | 当前开放年份 |
| `editionCode` | `string` | 当前或默认沙盘版本包编码 |
| `editionName` | `string` | 当前或默认沙盘版本包显示名 |
| `availableEditions` | `array` | 可选沙盘版本包列表；当前包含 `VIP_SERVICE_V1`，新增生产制造版后包含 `PRODUCTION_V1` |
| `initialBaselineSubmitted` | `bool` | 初始基线是否已提交 |
| `dictionaryRevision` | `int` | 当前比赛字典快照版本；未初始化时为当前编辑草稿或默认版本 |
| `selectedDictionarySchemeId` | `int64/null` | 初始化前当前选中的字典方案；使用版本默认名称时为空 |
| `dictionarySchemeCount` | `int` | 当前版本包下可用自定义字典方案数量 |
| `defaultRoute` | `string` | 管理员当前默认跳转地址 |

说明：

- 首版建议直接以 `sg_group` 实际记录数推断 `initialized` 与 `groupCount`。
- 当 `initialized=false` 时，`defaultRoute` 应返回 `admin/setup`；当 `initialized=true` 时，应返回 `admin/summary`。
- 当 `initialized=false` 时，前端允许管理员选择 `availableEditions` 中的版本包；当 `initialized=true` 时，该版本包只读展示。
- 当 `initialized=false` 时，前端允许管理员选择或编辑与当前版本包绑定的字典方案；当 `initialized=true` 时，前端只读展示当前比赛字典快照，并允许管理员通过 `admin-dictionary` 接口解锁修改显示名称。

#### 6.5.2 初始化比赛

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-control/initialize-game`
- 权限：`sandbox-game:admin-control:initialize-game`

Go DTO 建议：

- 请求：`InitializeGameReq`
- 响应：`InitializeGameResp`

请求体：

```json
{
  "groupCount": 6,
  "editionCode": "VIP_SERVICE_V1",
  "dictionarySchemeId": 12,
  "dictionaryItems": [
    {
      "itemCode": "orderType.BUSINESS_VIP",
      "displayName": "商务接待"
    }
  ]
}
```

生产制造版初始化示例：

```json
{
  "groupCount": 6,
  "editionCode": "PRODUCTION_V1"
}
```

请求字段建议：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `groupCount` | `int` | 是 | 本场比赛要初始化的小组数量，首版建议限制在 `1 ~ 10` |
| `editionCode` | `string` | 是 | 沙盘版本包编码，当前支持 `VIP_SERVICE_V1`；新增生产制造版后支持 `PRODUCTION_V1` |
| `dictionarySchemeId` | `int64/null` | 否 | 本次初始化选择的字典方案；为空表示使用版本默认名称或前端传入的临时名称 |
| `dictionaryItems` | `array` | 否 | 本次初始化最终采用的业务显示名称；服务端以稳定编码保存为当前比赛字典快照 |

响应字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `initialized` | `bool` | 初始化后应返回 `true` |
| `groupCount` | `int` | 实际初始化的小组数量 |
| `editionCode` | `string` | 本场比赛锁定的沙盘版本包编码 |
| `editionName` | `string` | 本场比赛锁定的沙盘版本包显示名 |
| `dictionaryRevision` | `int` | 初始化后当前比赛字典快照版本 |
| `createdGroupCount` | `int` | 本次创建的小组主数据数量 |
| `createdAccountCount` | `int` | 本次创建的账号数量 |
| `createdYearStateCount` | `int` | 本次创建的年份状态数量 |

规则：

- 仅管理员可调用。
- 仅允许在比赛未初始化时调用；若已初始化，应返回 `409`。
- `editionCode` 必须属于系统内置版本包；非法版本返回 `422`。
- `dictionarySchemeId` 若不为空，必须属于当前 `editionCode` 绑定版本包；非法或跨版本方案返回 `422`。
- `dictionaryItems` 只允许传稳定编码与显示名；显示名不允许为空；服务端不得用中文显示名驱动公式或字段含义。
- 服务端应以单事务一次性创建 `sg_group`、`sg_account`、`sg_group_year_state`。
- 服务端应在同一初始化事务内写入当前比赛字典快照。
- 初始化成功后，本场比赛沙盘版本包锁定，不提供运行中切换接口。
- 初始化成功后，管理员账号保留 `admin`，玩家账号建议按 `group01 ~ groupNN` 自动生成。
- 初始化成功后，系统应处于“`0年` 已开放、正式年份已预置但锁定”的初始状态。

#### 6.5.3 获取控制台配置

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
| `editionCode` | `string` | 当前沙盘版本包编码 |
| `editionName` | `string` | 当前沙盘版本包显示名 |
| `ruleVersion` | `string` | 当前规则版本 |
| `templateVersion` | `string` | 当前模板版本 |
| `operatingTemplateVersion` | `string` | 经营页字段模板版本 |
| `reportTemplateVersion` | `string` | 财报页字段模板版本 |
| `orderTemplateVersion` | `string` | 订单字段模板版本 |
| `processRuleVersion` | `string` | 流程规则版本 |
| `initialBaselineSubmitted` | `bool` | 初始基线是否已提交 |
| `initialBaselineSubmittedAt` | `string \| null` | 初始基线提交时间，RFC3339 |
| `initialBaselineSubmitterName` | `string \| null` | 初始基线提交人 |
| `hasRollbackPending` | `bool` | 是否存在回退补提 / 待重提小组 |
| `rollbackBlockedGroupCount` | `int` | 阻断开放下一年的回退补提小组数量 |
| `latestAdminAction` | `object \| null` | 最近一次关键管理员动作摘要 |

说明：

- 本接口服务于管理员年度控制页顶部状态区。
- `openNextYearBlockedReason` 只返回原因摘要，不返回逐组明细。

#### 6.5.4 更新最终年份

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

#### 6.5.5 开放下一年

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
- 不存在处于 `回退补提 / 待重提` 状态的未破产小组时才允许开放
- 已破产组在新年份仍保持不可编辑
- 开放成功后生成全局状态快照，用于审计和后续全局恢复扩展
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
| `hasRollbackPending` | `bool` | 开放后是否仍存在回退补提 / 待重提小组 |
| `rollbackBlockedGroupCount` | `int` | 开放后阻断后续开放的小组数量 |
| `latestAdminAction` | `object` | 本次动作摘要 |

#### 6.5.6 获取初始基线

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

#### 6.5.7 提交初始基线

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

### 6.5A `admin-dictionary`

#### 6.5A.1 获取版本包默认字典与当前显示名称

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-dictionary/get-current?editionCode=VIP_SERVICE_V1`
- 权限：登录态均可调用；初始化前主要供管理员赛前配置页使用，初始化后供玩家端和管理员端静默同步使用

用途：

- 初始化前：返回所选版本包的默认字典和前端编辑草稿所需字段。
- 初始化后：返回当前比赛字典快照、`dictionaryRevision` 和是否允许解锁修改。

返回字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| `initialized` | `bool` | 比赛是否已初始化 |
| `editionCode` | `string` | 版本包编码 |
| `editionName` | `string` | 版本包名称 |
| `dictionaryRevision` | `int` | 当前比赛字典快照版本；未初始化时可为默认版本 |
| `canUpdateCurrent` | `bool` | 是否允许修改当前比赛字典快照 |
| `items` | `array` | 字典项列表 |

字典项字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| `itemCode` | `string` | 稳定编码，例如 `market.LOCAL`、`orderType.BUSINESS_VIP`、`report.rawMaterials` |
| `itemCategory` | `string` | `MARKET` / `ORDER_TYPE` / `OPERATING` / `REPORT` / `BASELINE` |
| `defaultName` | `string` | 当前版本包默认名称 |
| `displayName` | `string` | 当前实际显示名称 |
| `displayOrder` | `int` | 展示排序 |
| `editable` | `bool` | 是否允许管理员修改 |
| `relatedPayload` | `object/null` | 前端映射辅助信息，不参与公式计算 |

规则：

- 前端默认不展示 `itemCode`，但接口必须返回稳定编码。
- 字典项只驱动显示名称，不驱动公式、字段含义或提交 payload。
- 系统状态、按钮、菜单、错误提示、Q1/Q2/Q3/Q4 等流程文案不通过字典替换。

#### 6.5A.2 查询字典方案列表

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-dictionary/list-schemes?editionCode=VIP_SERVICE_V1`
- 权限：管理员

返回字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| `editionCode` | `string` | 绑定版本包 |
| `list` | `array` | 字典方案列表 |

方案摘要字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `int64` | 字典方案 ID；`0` 表示版本默认名称 |
| `editionCode` | `string` | 绑定版本包 |
| `schemeName` | `string` | 字典方案名称 |
| `description` | `string` | 方案说明 |
| `builtIn` | `bool` | 是否内置项 |
| `itemCount` | `int` | 字典项数量 |
| `updatedAt` | `string` | 最近更新时间 |
| `updatedBy` | `string` | 最近更新人 |

规则：

- 只返回与当前版本包绑定的方案。
- 删除或编辑方案不影响当前比赛字典快照。

#### 6.5A.3 获取字典方案详情

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-dictionary/get-scheme-detail?id=12`
- 权限：管理员

返回字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `int64` | 字典方案 ID |
| `editionCode` | `string` | 绑定版本包 |
| `schemeName` | `string` | 方案名称 |
| `description` | `string` | 方案说明 |
| `builtIn` | `bool` | 是否内置项 |
| `items` | `array` | 方案字典项，字段同 `get-current.items` |
| `updatedAt` | `string` | 最近更新时间 |
| `updatedBy` | `string` | 最近更新人 |

#### 6.5A.4 保存字典方案

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-dictionary/save-scheme`
- 权限：管理员

请求字段：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `schemeId` | `int64/null` | 否 | 为空时新建，不为空时编辑 |
| `editionCode` | `string` | 是 | 绑定版本包 |
| `schemeName` | `string` | 是 | 方案名称 |
| `description` | `string` | 否 | 方案说明 |
| `items` | `array` | 是 | 方案字典项，仅包含 `itemCode / displayName` |

规则：

- 字典方案必须绑定版本包。
- 显示名不允许为空。
- 保存方案默认只影响以后新比赛，不自动影响当前比赛。
- 编辑已使用过的方案也不影响当前比赛字典快照，除非管理员后续执行“应用方案到当前比赛”。

#### 6.5A.5 删除字典方案

- 方法：`DELETE`
- 路径：`/api/v1/sandbox-game/admin-dictionary/delete-scheme?id=12`
- 权限：管理员

规则：

- 可删除已使用过的方案，因为当前系统不做多场比赛历史归档。
- 删除方案不影响当前比赛字典快照。

#### 6.5A.6 更新当前比赛显示名称

- 方法：`PUT`
- 路径：`/api/v1/sandbox-game/admin-dictionary/update-current`
- 权限：管理员

请求字段：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `items` | `array` | 是 | 当前比赛字典项，仅包含 `itemCode / displayName` |
| `reason` | `string` | 否 | 修改原因或备注 |

规则：

- 仅比赛已初始化后可调用。
- 只修改当前比赛字典快照。
- 更新成功后 `dictionaryRevision` 自增，并写入字段名称修改日志。
- 玩家端和管理员端通过轻量 revision 检查静默同步显示名称。
- 首版不要求前端提交 `baseRevision`；服务端在事务内读取当前配置并自增 `dictionaryRevision`。

#### 6.5A.7 应用字典方案到当前比赛

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-dictionary/apply-scheme-to-current`
- 权限：管理员

请求字段：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `schemeId` | `int64` | 是 | 要应用的字典方案 |

规则：

- 方案必须与当前比赛版本包绑定。
- 采用整体覆盖当前比赛字典快照方式。
- 前端确认弹窗展示变更数量和前若干项变更明细。
- 不改变公式、数据、订单、流程，只改变页面显示名称。
- 更新成功后写日志并返回新的 `dictionaryRevision`。

#### 6.5A.8 恢复当前版本默认名称

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-dictionary/restore-current-default`
- 权限：管理员

规则：

- 仅恢复当前比赛字典快照为当前版本包默认名称。
- 前端确认弹窗展示变更数量和前若干项变更明细。
- 更新成功后立即影响当前比赛显示名称，写入日志。

#### 6.5A.9 查询字段名称修改日志

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-dictionary/page-change-logs`
- 权限：管理员

返回列表字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `int64` | 日志 ID |
| `editionCode` | `string` | 沙盘版本包编码 |
| `changeType` | `string` | `INITIALIZE` / `UPDATE_CURRENT` / `APPLY_SCHEME` / `RESTORE_DEFAULT` / `SAVE_SCHEME` / `DELETE_SCHEME` |
| `schemeId` | `int64/null` | 关联方案 ID |
| `reason` | `string` | 修改原因或备注 |
| `revision` | `int` | 变更后当前比赛字典版本；方案维护日志可为 `0` |
| `operatorName` | `string` | 操作人 |
| `operateTime` | `string` | 操作时间 |
| `changedSummary` | `string` | 变更摘要 |

#### 6.5A.10 字典 revision 轻量检查

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-dictionary/get-revision`
- 权限：管理员或玩家登录态均可调用

规则：

- 返回当前比赛 `dictionaryRevision`。
- 玩家端页面可按较低频率轻量检查；发现版本变化后只刷新名称映射，不重新加载经营/财报/订单业务数据。
- 不允许出现刷新后页面跳到顶部、草稿丢失或输入焦点被强制打断。

#### 6.5.8 异常解锁某组某年

> 产品入口说明：异常解锁保留为 `回退与修正` 页面中的 `退回重提` 模式。该接口可继续沿用 `admin-control/unlock-year` 的旧路径，也可在实现 `admin-rollback` 时增加等价代理入口；无论路径如何，业务规则与日志口径必须一致。

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-control/unlock-year`
- 权限：`sandbox-game:admin-control:unlock-year`
- 新页面代理路径：`/api/v1/sandbox-game/admin-rollback/unlock-retry`
- 新页面代理权限：`sandbox-game:admin-rollback:unlock-retry`

Go DTO 建议：

- 请求：`UnlockYearReq`
- 响应：`UnlockYearResp`

请求体建议：

```json
{
  "groupId": 1,
  "yearNo": 2,
  "unlockTargetType": "OPERATING",
  "targetStageCode": "Q2"
}
```

字段建议：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `groupId` | `int64` | 是 | 目标组 |
| `yearNo` | `int` | 是 | 目标年份 |
| `unlockTargetType` | `string` | 是 | `OPERATING` / `REPORT` |
| `targetStageCode` | `string` | 条件必填 | 当 `unlockTargetType = OPERATING` 时必填，允许值：`Q1 / Q2 / Q3 / Q4 / YEAR_END` |
| `reason` | `string` | 否 | 退回说明；为空时服务端自动记录为 `管理员退回重提` |

规则：

- 仅管理员可执行。
- 轻量退回重提仍优先用于下一年尚未开放前的本年修正；若需要跨年恢复，使用 `admin-rollback` 的单组快照恢复接口。
- 前端不做复杂可解锁预判，管理员点击后直接提交；服务端负责最终硬校验。
- 仅允许对“已正式提交”的目标阶段或财报结果执行异常解锁。
- 执行前必须自动创建目标组的回退前安全快照。
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
- 已破产小组不得执行异常解锁以恢复 `NORMAL`；破产判定不可撤销。
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
| `safetySnapshotId` | `int64` | 本次执行前自动创建的安全快照 ID |

#### 6.5.9 后端实现规则清单

1. DTO 与枚举
- `UnlockYearReq` 必须新增 `unlockTargetType`、`targetStageCode` 字段。
- 服务端应显式定义 `UnlockTargetType` 与 `StageCode` 枚举，避免字符串散落在业务代码中。

2. 服务端校验
- 校验目标组、目标年份存在。
- 校验下一年尚未开放。
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
- 异常解锁不得清除破产标记；目标组已破产时应拒绝解锁请求。
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

### 6.5B `admin-rollback`

#### 6.5B.1 分页查询状态快照

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-rollback/page-snapshots`
- 权限：`sandbox-game:admin-rollback:query`

查询参数建议：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `snapshotScope` | `string` | 否 | `GROUP` / `GLOBAL` |
| `snapshotType` | `string` | 否 | `AUTO` / `MANUAL` / `SAFETY` |
| `groupId` | `int64` | 否 | 小组快照筛选条件 |
| `yearNo` | `int` | 否 | 年份筛选条件 |
| `stageCode` | `string` | 否 | 阶段筛选条件 |
| `pageNo` | `int` | 是 | 页码 |
| `pageSize` | `int` | 是 | 每页条数 |

返回列表字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `int64` | 快照 ID |
| `snapshotScope` | `string` | `GROUP` / `GLOBAL` |
| `snapshotType` | `string` | `AUTO` / `MANUAL` / `SAFETY` |
| `triggerCode` | `string` | 生成节点，例如 `STAGE_SUBMITTED`、`REPORT_SUBMITTED`、`ROLLBACK_STAGE_RESUBMITTED`、`ROLLBACK_REPORT_RESUBMITTED`、`ORDER_POOL_CONFIRMED`、`SEGMENT_COMPLETED`、`OPEN_NEXT_YEAR`、`BEFORE_ROLLBACK` |
| `groupId` | `int64 \| null` | 小组快照所属组；全局快照为空 |
| `groupName` | `string \| null` | 小组名称 |
| `yearNo` | `int \| null` | 快照对应年份 |
| `stageCode` | `string \| null` | 快照对应阶段 |
| `description` | `string` | 管理员说明或系统说明 |
| `createdByName` | `string` | 创建人 |
| `createdAt` | `string` | 创建时间，RFC3339 |

说明：

- 列表默认按 `createdAt` 倒序。
- 首版全局快照只用于审计和后续扩展，不提供全局恢复按钮。

#### 6.5B.2 查看快照详情

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-rollback/get-snapshot-detail`
- 权限：`sandbox-game:admin-rollback:query`

查询参数：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `snapshotId` | `int64` | 是 | 快照 ID |

返回字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `snapshot` | `object` | 快照摘要 |
| `stateSummary` | `object` | 快照中可读状态摘要 |
| `payloadPreview` | `object` | 经营、财报、订单、汇总等业务区块摘要；首版不要求展示完整 JSON |

#### 6.5B.3 手动创建状态快照

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-rollback/create-snapshot`
- 权限：`sandbox-game:admin-rollback:create-snapshot`

请求体建议：

```json
{
  "snapshotScope": "GROUP",
  "groupId": 1,
  "yearNo": 3,
  "stageCode": "Q2",
  "description": "管理员现场确认前手动留档"
}
```

字段建议：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `snapshotScope` | `string` | 是 | `GROUP` / `GLOBAL` |
| `groupId` | `int64` | 条件必填 | `GROUP` 快照必填 |
| `yearNo` | `int` | 否 | 目标年份；为空时按当前状态生成摘要 |
| `stageCode` | `string` | 否 | 目标阶段 |
| `description` | `string` | 是 | 管理员填写说明 |

规则：

- `GROUP` 快照只覆盖目标小组相关状态、经营、财报、订单归属、奖惩、汇总摘要。
- `GLOBAL` 快照覆盖全局配置、订单池、标段状态、全部小组当前关键状态摘要。
- 手动快照不改变任何业务状态，只写快照与管理员动作日志。

#### 6.5B.4 恢复单组快照

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-rollback/restore-group-snapshot`
- 权限：`sandbox-game:admin-rollback:restore-group-snapshot`

请求体建议：

```json
{
  "snapshotId": 1001,
  "reason": "第三组 3年 Q2 起数据录入错误，需要回到该节点重新提交",
  "confirmText": "确认恢复"
}
```

字段建议：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `snapshotId` | `int64` | 是 | 目标单组快照 ID |
| `reason` | `string` | 是 | 回退原因 |
| `confirmText` | `string` | 是 | 二次确认文本，防止误操作 |

响应字段建议：

| 字段 | 类型 | 说明 |
|---|---|---|
| `rollbackLogId` | `int64` | 回退日志 ID |
| `safetySnapshotId` | `int64` | 回退前自动创建的安全快照 ID |
| `groupId` | `int64` | 被回退小组 |
| `targetYearNo` | `int` | 恢复到的年份 |
| `targetStageCode` | `string` | 恢复到的阶段 |
| `yearStatus` | `string` | 恢复后的年份状态 |
| `stageStatus` | `string` | 恢复后的经营阶段状态 |
| `reportStatus` | `string` | 恢复后的财报状态 |
| `businessStatus` | `string` | 恢复后的经营状态 |
| `hasRollbackPending` | `bool` | 是否进入回退补提 / 待重提 |

规则：

- 首版只允许恢复 `GROUP` 快照，不允许恢复 `GLOBAL` 快照。
- 回退粒度是阶段级，不支持字段级回退。
- 支持单组跨年度回退，例如从 `5年` 回到 `3年 Q2`。
- 回退不改变全局当前开放年份。
- 执行前必须自动创建目标组的 `SAFETY` 快照。
- 回退后，目标节点之后的经营、财报与汇总生效结果进入失效态，但原值保留为 `失效草稿` 或历史记录；奖惩作为管理员独立业务事件，在普通退回重提时继续有效，恢复快照时才按快照时点恢复有效状态。
- 普通通知不回退。
- 已确认破产不可通过退回重提或恢复快照撤销，已破产小组不得执行回退恢复为 `NORMAL`。
- 单组快照回退不释放已选订单、不重排历史选单顺序、不重算市场龙头、不重新生成已确认订单池。
- 订单交付状态在目标节点之后的记录失效，玩家重新推进到对应年份后再确认交付。
- 已选的未来年份订单保持归属，但处于“失效但不释放”口径，待该组重新推进到对应年份时继续作为订单来源。
- 如果该组尚未参与某个已生成顺序的标段，系统只提示风险，现场节奏由管理员控制。
- 回退执行期间短暂锁定目标组写操作；玩家并发提交时返回 `409`，提示“状态已变化，请刷新页面”。
- 存在任意未破产小组处于回退补提 / 待重提状态时，管理员不能开放下一年。

#### 6.5B.5 后端自动快照节点

系统应在以下节点自动生成快照：

- 单组经营阶段提交成功后生成 `GROUP + AUTO` 快照。
- 单组财报提交成功后生成 `GROUP + AUTO` 快照。
- 管理员执行退回重提或恢复快照前生成 `GROUP + SAFETY` 快照。
- 管理员确认订单池后生成 `GLOBAL + AUTO` 快照。
- 每个标段完成后生成 `GLOBAL + AUTO` 快照。
- 管理员开放下一年后生成 `GLOBAL + AUTO` 快照。

快照在本场比赛内不自动清理。

### 6.6 `admin-notice`

#### 6.6.1 查看通知与奖惩最近记录

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-notice/get-records?limit=20`
- 权限：`sandbox-game:admin-notice:query`

Go DTO 建议：

- 响应：`AdminNoticeRecordsResp`

查询参数建议：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `limit` | `int` | 否 | 最近记录条数，首版建议默认 `20`，并限制最大值 |

规则：

- 返回最近普通通知列表与最近奖惩列表。
- 采用轻量轮询口径，不做 WebSocket；奖惩响应同时返回单调递增的 `adjustmentRevision`，供经营页和财报页判断是否需要局部同步。
- 管理端记录区不允许编辑或物理删除历史记录；对未破产小组的有效奖惩提供“作废”入口，并显示 `EFFECTIVE / VOIDED / SNAPSHOT_INACTIVE` 状态。

#### 6.6.2 发送普通通知

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-notice/send-general`
- 权限：`sandbox-game:admin-notice:send-general`

Go DTO 建议：

- 请求：`SendAdminGeneralNoticeReq`
- 响应：`SendAdminGeneralNoticeResp`

请求体建议：

```json
{
  "targetScope": "GROUP",
  "targetGroupId": 3,
  "content": "第三组请核对本年经营数据。",
  "pinned": true
}
```

字段建议：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `targetScope` | `string` | 是 | `ALL` / `GROUP` |
| `targetGroupId` | `int64` | 条件必填 | 当 `targetScope = GROUP` 时必填 |
| `content` | `string` | 是 | 自由文本通知内容 |
| `pinned` | `bool` | 否 | 是否置顶 |

规则：

- 普通通知支持发全体或单组。
- 普通通知仅负责展示，不参与经营、财报、汇总计算。
- 允许发送赛事播报类消息，例如 `第一小组已破产`。

#### 6.6.3 下发奖惩

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-notice/send-adjustment`
- 权限：`sandbox-game:admin-notice:send-adjustment`

Go DTO 建议：

- 请求：`SendAdminAdjustmentReq`
- 响应：`SendAdminAdjustmentResp`

请求体建议：

```json
{
  "groupId": 3,
  "yearNo": 0,
  "adjustmentType": "REWARD",
  "amount": 88,
  "reason": "主持人现场奖励"
}
```

字段建议：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `groupId` | `int64` | 是 | 目标组 |
| `yearNo` | `int` | 是 | 目标年份 |
| `adjustmentType` | `string` | 是 | `REWARD` / `PENALTY` |
| `amount` | `decimal` | 是 | 必须大于 `0` |
| `reason` | `string` | 是 | 奖惩原因 |

规则：

- `stageCode` 不再由管理员提交，服务端必须根据目标组当前状态自动决定 `Q1 / Q2 / Q3 / Q4 / YEAR_END`；财报填写和财报草稿阶段统一归属 `YEAR_END`。
- 奖惩属于计算型业务事件，必须进入经营页计算，并继续进入财报承接口径。
- 玩家端经营页中的 `额外收入 / 奖励` 与 `额外支出 / 罚款` 改为只读展示，不允许玩家自行录入。
- 若目标年份尚未开放，返回 `422`。
- 服务端按目标组当前运行状态判断是否允许下发，不得根据历史季度提交流水判断锁定。
- Q1、Q2、Q3、Q4、年末经营、财报填写和财报草稿阶段均允许下发。
- 财报正式提交或年份完成时返回 `422`，提示先退回到可编辑状态。
- 若目标组已破产，返回 `422`。
- 奖惩按税前收入 / 支出口径参与权威重算；响应应返回系统归属阶段、调整前后关键财务影响、最新 `adjustmentRevision` 和破产结果。
- 若所得税后现金小于 `0`，服务端在同一事务内生成破产快照并将小组标记为不可撤销的破产状态。

#### 6.6.4 预览奖惩影响

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-notice/preview-adjustment`
- 权限：`sandbox-game:admin-notice:send-adjustment`

创建奖惩预览请求：

```json
{
  "operation": "CREATE",
  "groupId": 3,
  "yearNo": 1,
  "adjustmentType": "PENALTY",
  "amount": 10,
  "reason": "现场讨论决定"
}
```

作废奖惩预览请求：

```json
{
  "operation": "VOID",
  "adjustmentId": 123
}
```

- `operation = CREATE` 时，创建奖惩所需字段必填。
- `operation = VOID` 时，仅需 `adjustmentId`；服务端从原事件读取组、年份、阶段、类型和金额。
- 财报填写阶段的预览和破产判定以服务端最新已保存财报草稿为依据；玩家浏览器内尚未保存的手工输入不作为管理员动作的权威数据源。

响应至少包含：

| 字段 | 说明 |
|---|---|
| `resolvedStageCode` | 服务端根据当前状态解析出的阶段 |
| `cashBefore / cashAfter` | 调整前后所得税后现金 |
| `preTaxProfitAfter / incomeTaxAfter / netProfitAfter` | 调整后的核心损益结果 |
| `totalEquityAfter` | 调整后的所有者权益 |
| `willBankrupt` | 是否将立即触发破产 |
| `calculationBasisSavedAt` | 本次影响预览所依据的经营 / 财报草稿保存时间 |

- 下发接口必须再次按最新状态权威重算，不能直接相信预览值；若状态已改变，应返回最新影响或明确要求管理员重新确认。
- 管理员确认按钮在 `willBankrupt=true` 时使用“确认下发并执行破产判定”。

#### 6.6.5 作废奖惩

- 方法：`POST`
- 路径：`/api/v1/sandbox-game/admin-notice/void-adjustment`
- 权限：`sandbox-game:admin-notice:void-adjustment`

请求建议：

```json
{
  "adjustmentId": 123,
  "reason": "管理员确认原奖惩录入有误"
}
```

- 不得编辑或物理删除原记录；服务端记录作废人、作废时间、作废原因并递增 `adjustmentRevision`。
- 仅允许作废未破产小组的有效奖惩；财报正式提交或年份完成时必须先退回，已破产时永久禁止。
- 作废前使用 `operation = VOID` 调用影响预览；正式作废时服务端再次权威重算。
- 作废导致所得税后现金小于 `0` 时，同样生成破产快照并判定不可撤销破产。

#### 6.6.6 玩家奖惩增量同步

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/player-notice/get-adjustment-sync?yearNo=1&knownRevision=12`
- 权限：玩家本人小组

- 页面可见时每 `3` 秒调用；页面隐藏时暂停，重新可见时立即调用。
- revision 未变化时返回轻量 `notModified=true`；变化时只返回有效奖惩聚合、必要明细、派生经营结果、财报预览、平衡状态和破产状态，不返回整页载荷。
- 客户端必须局部合并，不得覆盖未保存手工输入，不得改变焦点、光标、滚动位置、路由或年份。
- 玩家正式提交时后端仍须重新读取最新有效奖惩并权威重算，不使用 `knownRevision` 作为强制提交冲突闸门。

---

### 6.7 `admin-group-data`

#### 6.7.1 查看任意组经营页视图

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-group-data/get-operating-view?groupId=1&yearNo=1`
- 权限：`sandbox-game:admin-group-data:query`

用途：

- 管理员查看指定组指定年的经营页数据与状态
- 返回结构可复用玩家端 `get-year-view` 的数据模型

#### 6.7.2 查看任意组财报页视图

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-group-data/get-report-view?groupId=1&yearNo=1`
- 权限：`sandbox-game:admin-group-data:query`

#### 6.7.3 分页查询阶段提交日志

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-group-data/page-stage-submissions`
- 权限：`sandbox-game:admin-group-data:query`

查询参数建议：

- `groupId`
- `yearNo`
- `stageCode`
- `pageNo`
- `pageSize`

#### 6.7.4 分页查询财报提交日志

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/admin-group-data/page-report-submissions`
- 权限：`sandbox-game:admin-group-data:query`

---

### 6.8 `audit-log`

#### 6.8.1 分页查询异常解锁日志

- 方法：`GET`
- 路径：`/api/v1/sandbox-game/audit-log/page-unlock-log`
- 权限：`sandbox-game:audit-log:query`

#### 6.8.2 分页查询管理员动作日志

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
- `quarter.supplyChainOrderRecord` 为第 9 步订单数量留痕，按 `q1/q2/q3/q4` 组织；每季度固定包含 `basicProduct / standardProduct / precisionProduct / intelligentProduct` 四个稳定字段。
- 生产制造版和贵宾服务版共用上述字段编码，仅通过版本化业务显示字典切换产品/服务名称。
- 四项值为非负整数且当前季度必填，无订单时显式提交 `0`；该字段不跨年带入，不进入计算和正式订单模块。
- 旧经营 JSON 缺少该字段时按空对象兼容读取，不批量回填为 `0`；旧年份回退后重新提交目标季度时，按新规则补齐四项。

示例：

```json
{
  "quarter": {
    "supplyChainOrderRecord": {
      "q1": {
        "basicProduct": 3,
        "standardProduct": 2,
        "precisionProduct": 0,
        "intelligentProduct": 0
      }
    }
  }
}
```

### 7.2 `ReportManualPayload`

说明：

- 表示财报页绿色手工项集合
- 当前至少覆盖：在制品、成品、材料、所得税税率、企业认证得分、最佳生产/服务人力总监得分、关账速度得分
- 税率字段应限制在：`0.25 / 0.15 / 0`
- `reportBestSalesDirectorScore`、`reportBestCeoScore`、`reportTotalEquity` 均属于 `ReportComputedPayload`，不得由客户端通过本对象直接提交覆盖

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
    EditionCode                  string                  `json:"editionCode"`
    EditionName                  string                  `json:"editionName"`
    RuleVersion                  string                  `json:"ruleVersion"`
    TemplateVersion              string                  `json:"templateVersion"`
    OperatingTemplateVersion     string                  `json:"operatingTemplateVersion"`
    ReportTemplateVersion        string                  `json:"reportTemplateVersion"`
    OrderTemplateVersion         string                  `json:"orderTemplateVersion"`
    ProcessRuleVersion           string                  `json:"processRuleVersion"`
    InitialBaselineSubmitted     bool                    `json:"initialBaselineSubmitted"`
    InitialBaselineSubmittedAt   *string                 `json:"initialBaselineSubmittedAt"`
    InitialBaselineSubmitterName *string                 `json:"initialBaselineSubmitterName"`
    LatestAdminAction            *AdminActionSummaryResp `json:"latestAdminAction"`
}

type InitializeGameReq struct {
    GroupCount  int    `json:"groupCount" validate:"required,min=1,max=10"`
    EditionCode string `json:"editionCode" validate:"required"`
}

type InitializeGameResp struct {
    Initialized           bool   `json:"initialized"`
    GroupCount            int    `json:"groupCount"`
    EditionCode           string `json:"editionCode"`
    EditionName           string `json:"editionName"`
    CreatedGroupCount     int    `json:"createdGroupCount"`
    CreatedAccountCount   int    `json:"createdAccountCount"`
    CreatedYearStateCount int    `json:"createdYearStateCount"`
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
| `ErrAdminControlAlreadyInitialized` | `409` | 比赛已初始化后重复初始化 | `比赛已初始化，不能重复执行初始化` |
| `ErrAdminControlInitializeInvalid` | `422` | 初始化请求缺失或 `groupCount` 非法 | `初始化参数不合法` |
| `ErrAdminControlEditionRequired` | `422` | 初始化比赛未选择沙盘版本包 | `请选择沙盘版本` |
| `ErrAdminControlEditionInvalid` | `422` | 初始化比赛选择了不存在或未启用的版本包 | `沙盘版本不合法` |
| `ErrAdminControlEditionLocked` | `409` | 比赛初始化后尝试切换版本包 | `比赛已初始化，不能切换沙盘版本` |
| `ErrAdminControlFinalYearTooSmall` | `422` | `finalYear < currentOpenYear` | `最终年份不能小于当前开放年份` |
| `ErrAdminControlTargetYearMismatch` | `409` | `targetYearNo != currentOpenYear + 1` | `开放年份与当前状态不一致` |
| `ErrAdminControlFinalYearReached` | `409` | 已到最终年份仍尝试开放 | `已达到最终年份，无法继续开放` |
| `ErrAdminControlOpenNextYearBlocked` | `422` | 未满足开放下一年条件 | `当前仍有未完成财报的小组` |
| `ErrAdminControlRollbackPending` | `422` | 存在回退补提 / 待重提小组时尝试开放下一年 | `仍有小组处于回退补提中，不能开放下一年` |
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
| `ErrRollbackSnapshotNotFound` | `404` | 快照不存在 | `未找到状态快照` |
| `ErrRollbackSnapshotScopeUnsupported` | `422` | 首版尝试恢复全局快照 | `首版暂不支持恢复全局快照` |
| `ErrRollbackTargetInvalid` | `422` | 快照目标无法作为回退节点 | `快照目标状态不支持恢复` |
| `ErrRollbackReasonRequired` | `422` | 未填写回退原因 | `请填写回退原因` |
| `ErrRollbackConfirmRequired` | `422` | 未完成二次确认 | `请完成回退确认` |
| `ErrRollbackStateChanged` | `409` | 回退执行期间目标组状态被并发改变 | `状态已变化，请刷新页面` |
| `ErrRollbackGroupWriteLocked` | `409` | 目标组正在回退中仍提交写操作 | `小组状态正在修正，请刷新后重试` |

---

### 7.6 年度订单 DTO 草图（Go）

建议集中放在：`internal/http/dto/order_dto.go`

```go
type MarketCode string
type OrderType string

type SubmitMarketInvestmentReq struct {
    YearNo           int     `json:"yearNo" validate:"required,min=1"`
    MarketCode       string  `json:"marketCode" validate:"required"`
    MarketInvestment float64 `json:"marketInvestment" validate:"min=0"`
}

type SelectOrderReq struct {
    YearNo     int    `json:"yearNo" validate:"required,min=1"`
    MarketCode string `json:"marketCode" validate:"required"`
    OrderType  string `json:"orderType" validate:"required"`
    OrderID    int64  `json:"orderId" validate:"required,min=1"`
}

type PassOrderSegmentReq struct {
    YearNo     int    `json:"yearNo" validate:"required,min=1"`
    MarketCode string `json:"marketCode" validate:"required"`
    OrderType  string `json:"orderType" validate:"required"`
}

type DeliverOrdersReq struct {
    YearNo    int     `json:"yearNo" validate:"required,min=1"`
    StageCode string  `json:"stageCode" validate:"required"`
    OrderIDs  []int64 `json:"orderIds" validate:"required,min=1"`
}

type AdminSkipCurrentGroupReq struct {
    YearNo     int    `json:"yearNo" validate:"required,min=1"`
    MarketCode string `json:"marketCode" validate:"required"`
    OrderType  string `json:"orderType" validate:"required"`
    GroupID    int64  `json:"groupId" validate:"required,min=1"`
    Reason     string `json:"reason" validate:"required"`
}

type UpdateOrderControlConfigReq struct {
    YearNo int                    `json:"yearNo" validate:"required,min=1"`
    Items  []OrderControlItemReq  `json:"items" validate:"required"`
}

type OrderControlItemReq struct {
    MarketCode        string `json:"marketCode" validate:"required"`
    OrderType         string `json:"orderType" validate:"required"`
    OrderCount        int    `json:"orderCount" validate:"min=0"`
    ReleaseSequenceNo int    `json:"releaseSequenceNo" validate:"required,min=1"`
}

type UpdateOrderMarketEnabledReq struct {
    YearNo  int                         `json:"yearNo" validate:"required,min=1"`
    Markets []OrderMarketEnabledItemReq `json:"markets" validate:"required"`
}

type OrderMarketEnabledItemReq struct {
    MarketCode string `json:"marketCode" validate:"required"`
    Enabled    bool   `json:"enabled"`
}
```

说明：

- 金额类型在正式实现时建议统一使用项目内 `decimal` 类型，不用 `float64` 作为最终落库口径。
- 上述 DTO 仅表达接口结构草图，最终字段类型应与现有项目金额类型保持一致。

### 7.7 年度订单错误常量建议

| Go 常量名 | 默认 HTTP | 触发场景 | 默认提示语建议 |
|---|---|---|---|
| `ErrOrderNotRequiredForDemoYear` | `422` | `0年` 试图提交订单动作 | `0年不需要订单` |
| `ErrOrderMarketDisabledInvestmentMustBeZero` | `422` | 未开启市场提交非零投入 | `该市场未开启，系统自动按0提交，不允许填写非0投入` |
| `ErrOrderMarketConfigLocked` | `409` | 订单池确认后或已有小组提交市场投入后修改市场配置 | `市场配置已锁定，不能修改开启状态或投入上限` |
| `ErrOrderPoolNotGenerated` | `422` | 订单池未生成或未确认就尝试开标/选单 | `订单池尚未确认` |
| `ErrOrderPoolNotConfirmed` | `422` | 未确认预览批次就生成选单顺序 | `请先确认订单池` |
| `ErrOrderCountInvalid` | `422` | 订单卡片数量为负数或非整数 | `订单数量必须为非负整数` |
| `ErrOrderInvestmentIncomplete` | `422` | 未提交当前订单模板的完整市场投入 | `请先提交本年全部市场投入` |
| `ErrOrderInvestmentNotInteger` | `422` | 市场投入或单市场投入上限包含小数 | `市场投入和投入上限必须为整数` |
| `ErrOrderInvestmentLimitExceeded` | `422` | 某市场 4 项投入合计超过单市场上限 | `市场投入超过单市场上限` |
| `ErrOrderInvestmentSubmitted` | `409` | 重复提交市场投入 | `本年市场投入已提交，不能修改` |
| `ErrOrderMarketNoInvestment` | `422` | 当前组无标段投入却尝试选单 | `本组未投入该标段，不能选择订单` |
| `ErrOrderSelectionNotReady` | `422` | 选单顺序未生成 | `该标段尚未进入选单阶段` |
| `ErrOrderSelectionNotTurn` | `409` | 未轮到当前组 | `当前还未轮到本组选择订单` |
| `ErrOrderSegmentNotReleased` | `422` | 标段尚未释放 | `该标段尚未释放` |
| `ErrOrderReleaseSequenceDuplicated` | `422` | 同一年标段释放顺序重复 | `标段释放顺序重复` |
| `ErrOrderReleaseSequenceLocked` | `409` | 开标后修改释放顺序 | `标段释放顺序已锁定，不能修改` |
| `ErrOrderReleaseSequenceOutOfOrder` | `409` | 管理员跳序释放标段 | `请按预设顺序释放下一个标段` |
| `ErrOrderAlreadySelectedInRound` | `409` | 同标段同轮重复选单 | `本组已在本轮选择订单` |
| `ErrOrderRoundAlreadyPassed` | `409` | 放弃后再次操作当前轮 | `本组已放弃本轮，不能再次选择` |
| `ErrOrderRoundNotReady` | `409` | 标段不处于下一轮待开启状态 | `当前没有可开启的下一轮` |
| `ErrOrderRoundAlreadyOpened` | `409` | 重复开启已经开始或结束的轮次 | `该轮次已经开启，不能重复操作` |
| `ErrOrderPoolExhausted` | `409` | 订单池已选空后继续选择或开启下一轮 | `订单池已选空，当前标段已经结束` |
| `ErrOrderUnavailable` | `409` | 订单已被选择或不可选 | `该订单已不可选择` |
| `ErrOrderPoolLocked` | `409` | 订单池确认后仍尝试修改数量或重新生成 | `订单池已确认，不能重新生成` |
| `ErrOrderDeliveryRevenueMismatch` | `422` | 季度销售收入与交付订单金额合计不一致 | `本季度销售收入必须等于交付订单金额合计` |
| `ErrOrderAdminCannotSelect` | `403` | 管理员尝试代选订单 | `首版不支持管理员代选订单` |
| `ErrOrderPrerequisiteIncomplete` | `422` | 正式年份 Q1 进入/提交时订单前置未完成 | `本年订单竞标尚未结束` |

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
  - 正式年份订单前置未完成时提交 Q1
  - 财报平衡失败
  - 已完成年份重复提交财报
  - 下一年已开放后再解锁上一年

---

## 10. 当前待后续补齐但不影响本版成立的内容

以下内容后续可补充，不影响本接口设计正式版成立：

- `OperatingPayload` 的最终字段映射清单
- `ReportManualPayload` 的最终字段命名说明
- 订单生成批次与订单卡片字段的完整样例 JSON
- 详细接口示例 JSON 样例库
- `.http` 用例文件
- 异常解锁日志的高级筛选条件与导出策略

