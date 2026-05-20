# 沙盘经营系统 Excel 字段与页面映射文档（首版）

> 更新日期：2026-05-19
> 适用方式：基于当前 Excel 模板、正式需求与接口/数据库设计，先冻结“结构性映射”和“关键字段映射”，作为后续前端页面、接口字段、规则实现的统一对照表。  
> 文档定位：本文件不是完整逐格抄录 Excel，而是说明“哪些字段已经冻结、哪些映射当前仍未逐格穷尽”。

## 1. 文档目标

本文件用于回答四个问题：

1. Excel 中哪些区块已经有明确系统字段归属
2. 哪些字段属于玩家手工输入、管理员输入、系统计算
3. 哪些字段会进入接口、数据库、页面映射
4. 哪些字段仍待进一步补齐映射

---

## 2. 映射原则

### 2.1 三层命名分离

同一个字段至少区分三种命名：

| 层 | 命名示例 | 说明 |
|---|---|---|
| Excel 显示名 | 所得税税率 | 面向表格使用者 |
| 系统字段名 | `incomeTaxRate` | 面向前后端实现 |
| 数据库存储名 | `income_tax_rate` | 面向落库 |

### 2.2 首版冻结范围

首版当前先冻结：

- 页面区块映射
- 初始基线关键字段映射
- 财报手工项映射
- 财报自动结果关键字段映射
- 汇总指标映射
- 跨年承接关键字段映射
- 经营页 payload 区块锚点结构

首版暂不逐格穷尽：

- 全经营页逐格中文标签的完整清单
- 非业务型帮助文案与说明文案的排版细修
- 少量不影响规则实现的展示层标签补齐

补充约束：

- 当前中文显示名已按最新版 `1组 最终版.xlsx` 同步为 `材料`、`待折资产`、`生产线残值`、`厂房` 等口径。
- 系统字段名与数据库字段名保持稳定，不因 Excel 中文标签微调直接改名。

### 2.3 字段分类

| 类型 | 说明 |
|---|---|
| `ADMIN_INPUT` | 管理员录入 |
| `PLAYER_INPUT` | 玩家手工录入 |
| `SYSTEM_CALCULATED` | 后端规则层计算 |
| `SYSTEM_DERIVED` | 汇总或状态派生字段 |

---

## 3. 区块级映射

### 3.1 经营页区块

| 区块编码 | 页面区块 | 年份适用范围 | 阶段归属 | 数据类型 |
|---|---|---|---|---|
| `operating.beginning` | 年初区 | `0年 ~ 最终年` | `Q1` 提交时一起处理 | `PLAYER_INPUT + SYSTEM_CALCULATED` |
| `player.annualOrder` | 年度订单页 | `1年 ~ 最终年` | 经营页 `Q1` 提交前置流程 | `PLAYER_INPUT + SYSTEM_DERIVED` |
| `operating.q1` | Q1 区 | `0年 ~ 最终年` | `Q1` | `PLAYER_INPUT + SYSTEM_CALCULATED` |
| `operating.q2` | Q2 区 | `0年 ~ 最终年` | `Q2` | `PLAYER_INPUT + SYSTEM_CALCULATED` |
| `operating.q3` | Q3 区 | `0年 ~ 最终年` | `Q3` | `PLAYER_INPUT + SYSTEM_CALCULATED` |
| `operating.q4` | Q4 区 | `0年 ~ 最终年` | `Q4` | `PLAYER_INPUT + SYSTEM_CALCULATED` |
| `operating.yearEnd` | 年末区 | `0年 ~ 最终年` | `YEAR_END` | `PLAYER_INPUT + SYSTEM_CALCULATED` |
| `operating.extra` | 额外收入/罚款区 | `0年 ~ 最终年` | 随当期提交 | `PLAYER_INPUT` |
| `operating.derived` | 指标与汇总区 | `0年 ~ 最终年` | 只读展示 | `SYSTEM_CALCULATED` |

### 3.2 财报页区块

| 区块编码 | 页面区块 | 年份适用范围 | 数据类型 |
|---|---|---|---|
| `report.profit` | 损益表区 | `0年 ~ 最终年` | `SYSTEM_CALCULATED + PLAYER_INPUT(税率)` |
| `report.assets` | 资产侧 | `0年 ~ 最终年` | `SYSTEM_CALCULATED + PLAYER_INPUT(存货类)` |
| `report.liabilityEquity` | 负债权益侧 | `0年 ~ 最终年` | `SYSTEM_CALCULATED` |
| `report.validation` | 平衡校验区 | `0年 ~ 最终年` | `SYSTEM_DERIVED` |

### 3.3 管理员区块

| 区块编码 | 页面区块 | 数据类型 |
|---|---|---|
| `admin.initialBaseline` | 初始基线 | `ADMIN_INPUT` |
| `admin.orderControl` | 订单数量控制与订单池生成 | `ADMIN_INPUT + SYSTEM_DERIVED` |
| `admin.orderPool` | 订单池查看与市场竞标控制 | `ADMIN_INPUT + SYSTEM_DERIVED` |
| `admin.summary` | 汇总指标区 | `SYSTEM_DERIVED` |
| `admin.control` | 年度配置与推进 | `ADMIN_INPUT + SYSTEM_DERIVED` |
| `admin.unlock` | 异常解锁 | `ADMIN_INPUT + SYSTEM_DERIVED` |

---

## 4. 当前已冻结的关键字段映射

### 4.1 财报手工输入项

| 系统字段名 | 数据库存储名 | 当前 Excel 位置 | 当前显示名 | 类型 | 说明 |
|---|---|---|---|---|---|
| `workInProgress` | `work_in_progress` | `G15` | 在制品 | `PLAYER_INPUT` | 财报手工项 |
| `finishedGoods` | `finished_goods` | `G16` | 成品 | `PLAYER_INPUT` | 财报手工项 |
| `rawMaterials` | `raw_materials` | `G17` | 材料 | `PLAYER_INPUT` | 财报手工项 |
| `incomeTaxRate` | `income_tax_rate` | `C24` | 所得税税率 | `PLAYER_INPUT` | 下拉选项：`0.25/0.15/0` |

### 4.2 管理员初始基线关键字段

| 系统字段名 | 数据库存储名 | 当前 Excel 位置 | 当前显示名 | 类型 | 说明 |
|---|---|---|---|---|---|
| `baselineSalesRevenue` | `baseline_sales_revenue` | `C4` | 销售收入 | `ADMIN_INPUT` | 初始基线损益表输入 |
| `baselineDirectCost` | `baseline_direct_cost` | `C5` | 直接成本 | `ADMIN_INPUT` | 初始基线损益表输入 |
| `baselineComprehensiveCost` | `baseline_comprehensive_cost` | `C10` | 综合费用 | `ADMIN_INPUT` | 初始基线损益表输入 |
| `baselineDepreciation` | `baseline_depreciation` | `C11` | 折旧 | `ADMIN_INPUT` | 初始基线损益表输入 |
| `baselineFinanceIncomeExpense` | `baseline_finance_income_expense` | `C16` | 财务收入/支出 | `ADMIN_INPUT` | 初始基线损益表输入 |
| `baselineExtraIncomeExpense` | `baseline_extra_income_expense` | `C17` | 额外收入/支出 | `ADMIN_INPUT` | 初始基线损益表输入 |
| `baselineIncomeTax` | `baseline_income_tax` | `C21` | 所得税 | `ADMIN_INPUT` | 初始基线税额直接录入 |
| `baselineWorkInConstruction` | `baseline_work_in_construction` | `G5` | 在建生产线 | `ADMIN_INPUT` | 初始基线非流动资产输入；制造业版显示名为“在建生产线”，服务版显示名为“在建贵宾厅” |
| `baselineFactoryAsset` | `baseline_factory_asset` | `G6` | 厂房 | `ADMIN_INPUT` | 初始基线资产输入 |
| `baselineLineResidual` | `baseline_line_residual` | `G7` | 生产线残值 | `ADMIN_INPUT` | 初始基线资产输入 |
| `baselineDepreciableAsset` | `baseline_depreciable_asset` | `G8` | 待折资产 | `ADMIN_INPUT` | 初始基线资产输入 |
| `baselineCash` | `baseline_cash` | `G13` | 现金 | `ADMIN_INPUT` | 作为 `0年经营` 的现金起点 |
| `baselineReceivable` | `baseline_receivable` | `G14` | 应收款 | `ADMIN_INPUT` | 初始基线流动资产输入 |
| `baselineWorkInProgress` | `baseline_work_in_progress` | `G15` | 在制品 | `ADMIN_INPUT` | 初始基线流动资产输入 |
| `baselineFinishedGoods` | `baseline_finished_goods` | `G16` | 成品 | `ADMIN_INPUT` | 初始基线流动资产输入 |
| `baselineRawMaterials` | `baseline_raw_materials` | `G17` | 材料 | `ADMIN_INPUT` | 初始基线流动资产输入 |
| `baselineShortTermLoan` | `baseline_short_term_loan` | `K7` | 短期负债 | `ADMIN_INPUT` | 初始基线负债输入 |
| `baselineLongTermLoan` | `baseline_long_term_loan` | `K8` | 长期负债 | `ADMIN_INPUT` | 初始基线负债输入 |
| `baselineShareCapital` | `baseline_share_capital` | `K15` | 股东资本 | `ADMIN_INPUT` | 初始基线权益输入 |
| `baselineRetainedEarnings` | `baseline_retained_earnings` | `K16` | 利润留存 | `ADMIN_INPUT` | 初始基线权益输入 |

### 4.3 财报自动计算关键字段

| 系统字段名 | 来源表与单元格 | 当前显示名 | 类型 | 主要用途 |
|---|---|---|---|---|
| `reportSalesRevenue` | 财报 `C4` | 销售收入 | `SYSTEM_CALCULATED` | 汇总收入、财报展示 |
| `reportDirectCost` | 财报 `C5` | 直接成本 | `SYSTEM_CALCULATED` | 毛利计算 |
| `reportGrossProfit` | 财报 `C7` | 毛利 | `SYSTEM_CALCULATED` | 损益表展示 |
| `reportComprehensiveCost` | 财报 `C10` | 综合费用 | `SYSTEM_CALCULATED` | 营业利润计算 |
| `reportDepreciation` | 财报 `C11` | 折旧 | `SYSTEM_CALCULATED` | 营业利润计算 |
| `reportOperatingProfit` | 财报 `C13` | 营业利润 | `SYSTEM_CALCULATED` | 损益表展示 |
| `reportFinanceIncomeExpense` | 财报 `C16` | 财务收入/支出 | `SYSTEM_CALCULATED` | 税前利润计算 |
| `reportExtraIncomeExpense` | 财报 `C17` | 额外收入/支出 | `SYSTEM_CALCULATED` | 税前利润计算 |
| `reportPreTaxProfit` | 财报 `C19` | 税前利润 | `SYSTEM_CALCULATED` | 所得税计算基数 |
| `reportIncomeTax` | 财报 `C21` | 所得税 | `SYSTEM_CALCULATED` | 税后利润与跨年承接 |
| `reportNetProfit` | 财报 `C22` | 年度净利润 | `SYSTEM_CALCULATED` | 汇总利润、权益更新 |
| `reportCash` | 财报 `G13` | 现金 | `SYSTEM_CALCULATED` | 跨年现金承接、资产侧展示 |
| `reportReceivable` | 财报 `G14` | 应收款 | `SYSTEM_CALCULATED` | 跨年应收承接 |
| `reportTotalAssets` | 财报 `G22` | 总资产 | `SYSTEM_CALCULATED` | 财报平衡校验 |
| `reportShortTermLiability` | 财报 `K7` | 短期负债 | `SYSTEM_CALCULATED` | 负债侧展示、跨年承接 |
| `reportLongTermLiability` | 财报 `K8` | 长期负债 | `SYSTEM_CALCULATED` | 负债侧展示、跨年承接 |
| `reportTotalEquity` | 财报 `K19` | 总股东权益 | `SYSTEM_CALCULATED` | 排名、破产判定 |
| `reportTotalLiabilityEquity` | 财报 `K22` | 总负债权益 | `SYSTEM_CALCULATED` | 财报平衡校验 |

### 4.4 汇总指标字段

| 系统字段名 | 数据库存储名 | 来源表与单元格 | 当前显示名 | 类型 |
|---|---|---|---|---|
| `summaryRevenue` | `summary_revenue` | 财报 `C4` | 收入 | `SYSTEM_DERIVED` |
| `summaryProfit` | `summary_profit` | 财报 `C22` | 利润 | `SYSTEM_DERIVED` |
| `summaryEquity` | `summary_equity` | 财报 `K19` | 权益 | `SYSTEM_DERIVED` |

### 4.5 跨年承接关键字段

| 系统字段名 | 上年来源 | 当年去向 | 说明 |
|---|---|---|---|
| `previousIncomeTax` | 上年财报 `C21` | 当年经营 `M1` | 年初缴税基数 |
| `previousShortTermLoan` | 上年财报 `K7` | 当年经营贷款余额计算 | 短期负债承接 |
| `previousLongTermLoan` | 上年财报 `K8` | 当年经营贷款余额计算 | 长期负债承接 |
| `previousEquipmentResidual` | 上年财报 `G7` | 当年经营 `O49` | 生产线残值承接 |
| `previousDepreciableAsset` | 上年财报 `G8` | 当年经营 `O50` | 待折资产承接 |
| `previousCash` | 上年财报 `G13` | 当年经营 `C59` | 现金承接 |
| `previousReceivable` | 上年财报 `G14` | 当年财报资产侧 | 应收承接 |
| `shareholderCapital` | 上年财报 `K15` | 当年财报权益侧 | 股东资本承接 |
| `retainedEarnings` | 上年财报 `K16+K17` | 当年财报权益侧 | 利润留存承接 |

### 4.6 经营页季度现金核对展示

| 系统字段名 | 数据库存储名 | 当前 Excel 位置 | 当前显示名 | 类型 | 说明 |
|---|---|---|---|---|---|
| `quarterCashCheckLabel` | - | `B41` | 核对季末现金 | `SYSTEM_DERIVED` | 主表标签，仅展示 |
| `q1QuarterEndCashCheck` | - | `G41` | Q1 季末现金核对值 | `SYSTEM_CALCULATED` | 主表只读展示 |
| `q2QuarterEndCashCheck` | - | `I41` | Q2 季末现金核对值 | `SYSTEM_CALCULATED` | 主表只读展示 |
| `q3QuarterEndCashCheck` | - | `K41` | Q3 季末现金核对值 | `SYSTEM_CALCULATED` | 主表只读展示 |
| `q4QuarterEndCashCheck` | - | `M41` | Q4 季末现金核对值 | `SYSTEM_CALCULATED` | 主表只读展示 |
| `periodEndCashTitle` | - | `A59` | 期末现金 | `SYSTEM_DERIVED` | 主表标题位置 |

### 4.6.1 年度订单模块关键字段

说明：

- 订单推算来源为 `道具-订单推算（服务企业）.xlsx`。
- 系统按该 Excel 的公式链实现订单生成器，Excel 作为规则依据，不作为运行时订单池上传结果。
- 首版管理员先配置当年市场开启状态，再配置 `年份 + 市场 + 订单类型` 的订单卡片数量，系统生成预览订单池，确认后保存为固定业务数据。
- 本地市场默认开启，区域/全国/全球默认关闭；未开启市场不生成订单、不进入抢单，但玩家仍需在投入表中手动填写 `0`。
- 生成批次需保存随机种子、公式版本、控制台参数快照和公式参数快照，用于复盘与审计。

| 系统字段名 | 数据库存储名 | 来源/页面 | 当前显示名 | 类型 | 说明 |
|---|---|---|---|---|---|
| `orderYearNo` | `year_no` | 订单池 | 年份 | `SYSTEM_DERIVED` | 订单所属年份，`1年 ~ 最终年` |
| `marketCode` | `market_code` | 订单池/年度订单页 | 市场 | `ADMIN_INPUT + PLAYER_INPUT` | `LOCAL / REGIONAL / NATIONAL / GLOBAL` |
| `marketEnabled` | `market_enabled` | 管理员订单管理/年度订单页 | 市场开启 | `ADMIN_INPUT + SYSTEM_DERIVED` | 本地市场默认开启，其他市场默认关闭；未开启市场投入必须为 `0` |
| `orderType` | `order_type` | 订单池/年度订单页 | 订单类型 | `ADMIN_INPUT + SYSTEM_DERIVED` | `代办过检 / 两舱贵宾 / 商务贵宾 / 会员定制` |
| `segmentCode` | `segment_code` | 订单池/年度订单页 | 标段 | `SYSTEM_DERIVED` | `市场 + 订单类型` 的组合编码 |
| `releaseSequenceNo` | `release_sequence_no` | 管理员订单管理 | 标段释放顺序 | `ADMIN_INPUT` | 管理员配置同一年内标段开标先后 |
| `generationBatchId` | `generation_batch_id` | 订单池 | 订单生成批次 | `SYSTEM_DERIVED` | 系统按 Excel 公式链生成订单池后固化的批次 ID |
| `formulaVersion` | `formula_version` | 订单生成批次 | 公式版本 | `SYSTEM_DERIVED` | 订单生成规则版本 |
| `randomSeed` | `random_seed` | 订单生成批次 | 随机种子 | `SYSTEM_DERIVED` | 后台复盘字段，首版页面可不展示 |
| `cardSequenceNo` | `card_sequence_no` | 订单池 | 订单卡片序号 | `SYSTEM_DERIVED` | 每个标段 `1 ~ 15` |
| `businessOrderNo` | `business_order_no` | 订单池 | 订单编号 | `SYSTEM_DERIVED` | 页面主编号，如 `CARD-01`；不使用数据库自增 ID 作为主要展示编号 |
| `segmentStatus` | `segment_status` | 年度订单页 | 标段状态 | `SYSTEM_DERIVED` | `MARKET_DISABLED / NO_ORDER_CONFIG / WAITING_RELEASE / SELECTING / COMPLETED / SKIPPED` |
| `selectionSequenceNo` | `sequence_no` | 年度订单页 | 选单顺序 | `SYSTEM_DERIVED` | 玩家和管理员均可查看完整顺序，但玩家端不展示排序依据 |
| `selectionStatus` | `selection_status` | 年度订单页 | 小组标段状态 | `SYSTEM_DERIVED` | `INELIGIBLE / WAITING / CURRENT / SELECTED / PASSED / ADMIN_SKIPPED` |
| `isMarketLeader` | `is_market_leader` | 年度订单页 | 是否市场龙头 | `SYSTEM_DERIVED` | 市场龙头只在对应市场生效，本年 `0` 投入仍可优先选单 |
| `orderAmount` | `order_amount` | 订单池 | 订单金额 | `SYSTEM_DERIVED` | 系统按 Excel 公式链生成后写入 |
| `orderQuantity` | `order_quantity` | 订单池 | 数量 | `SYSTEM_DERIVED` | 单张订单卡的随机服务数量，不是控制台订单卡片数量 |
| `unitPrice` | `unit_price` | 订单池 | 单价 | `SYSTEM_DERIVED` | 订单卡片单价字段 |
| `accountTerm` | `account_term` | 订单池 | 账期 | `SYSTEM_DERIVED` | 订单卡片账期字段；首版只展示，不参与应收账款自动计算 |
| `orderPoolStatus` | `status` | 订单池 | 订单状态 | `SYSTEM_DERIVED` | `AVAILABLE / SELECTED / VOID` 等后续实现枚举 |
| `marketInvestment` | `market_investment` | 年度订单页 | 市场投入 | `PLAYER_INPUT` | 玩家按 `市场 + 订单类型` 提交 16 项投入，提交后不可修改；未开启市场必须手动填写 `0` |
| `selectedOrderId` | `order_id` | 年度订单页 | 已选订单 | `PLAYER_INPUT + SYSTEM_DERIVED` | 玩家轮到本组且当前标段释放时选择的订单 |
| `selectedOrderAmount` | `selected_order_amount` | 年度订单页/经营页 | 已选订单金额 | `SYSTEM_DERIVED` | 汇总进入经营页订单总额 |
| `selectedOrderTotal` | - | 经营页 | 订单总额 | `SYSTEM_DERIVED` | 正式年份由本组全部标段已选订单金额汇总 |
| `marketInvestmentTotal` | - | 经营页 | 市场投入 | `SYSTEM_DERIVED` | 正式年份由本组 16 项标段投入汇总 |
| `orderDeliveryStatus` | `delivery_status` | 年度订单页/经营页 | 交付状态 | `SYSTEM_DERIVED` | `SELECTED / DELIVERED / UNFINISHED` |
| `deliveredStageCode` | `delivered_stage_code` | 年度订单页/经营页 | 交付季度 | `SYSTEM_DERIVED` | 玩家执行交付操作时绑定当前经营季度 |
| `pollingIntervalSeconds` | - | 年度订单页 | 自动轮询间隔 | `SYSTEM_DERIVED` | 首版固定 `3` 秒，不做 WebSocket |

### 4.6.2 订单模块与经营页映射

| 订单模块字段 | 经营页字段/区块 | 映射规则 |
|---|---|---|
| `marketInvestment` 按 16 个标段汇总 | `beginning.marketBid.marketInvestmentTotal` / 市场投入 | `1年 ~ 最终年` 使用订单模块汇总并只读带入；`0年` 不适用 |
| `selectedOrderAmount` 按全部标段汇总 | `beginning.marketBid.selectedOrderTotal` / 订单总额 | `1年 ~ 最终年` 使用订单模块汇总；剩余未选订单不计入 |
| `selectedOrderAmount` 按交付季度汇总 | `quarter.deliverySettlement.salesRevenue` / 交货销售额 | 玩家交付一个或多个完整订单后，该季度销售收入必须匹配交付订单金额合计 |
| `orderDeliveryStatus` | 年度订单页订单状态 | 经营页交付后回写，年末未交付为 `UNFINISHED` |
| `accountTerm` | 订单卡片展示 | 首版仅展示，不自动生成或移动经营页应收账款 |

### 4.7 `OperatingPayload` 建议结构

说明：

- 本节冻结的是“后端 payload 结构”和“页面区块锚点”，不是要求把全部单元格逐格建成数据库字段。
- 块内仍然会同时存在玩家输入格、公式格和只读展示格；提交时应以区块语义处理，而不是以单元格坐标直接驱动业务。
- 当前结构已经足够支撑 `save-draft`、`get-year-view`、`submit-stage` 的接口实现。

| payloadPath | Excel 区块锚点 | 建议数据结构 | 提交节点 | 当前冻结程度 | 说明 |
|---|---|---|---|---|---|
| `beginning.taxAndPlanning` | `A1:O2` | `object` | `Q1` | `结构冻结` | 包含年初缴税、年度计划收入/综合费用等年初信息；与 Q1 一起提交 |
| `beginning.marketBid` | `B3:O8` | `array<object>` | `Q1` | `结构冻结，标签可细修` | 按市场维度组织，含各产品总价、订单总额、市场投入 |
| `quarter.shortTermLoan` | `B9:P12` | `quarterMap<object>` | `Q1/Q2/Q3/Q4` | `结构冻结` | 包含到期还贷、付利息、新增贷款 |
| `quarter.materialPayment` | `A13:O18` | `quarterMap<object>` | `Q1/Q2/Q3/Q4` | `结构冻结` | 按产品维度记录材料费 |
| `quarter.productionLineAdjustment` | `B19:O28` | `quarterMap<object>` | `Q1/Q2/Q3/Q4` | `结构冻结` | 包含变更产品、拆除、出售、新安装、转固相关输入 |
| `quarter.humanResource` | `B29:O30` | `quarterMap<object>` | `Q1/Q2/Q3/Q4` | `结构冻结` | 招聘、辞退、待岗、培训等费用输入 |
| `quarter.salaryAndProduction` | `B31:O32` | `quarterMap<object>` | `Q1/Q2/Q3/Q4` | `结构冻结` | 生产线投料与工资相关输入 |
| `quarter.researchAndManagement` | `B33:O34` | `quarterMap<object>` | `Q1/Q2/Q3/Q4` | `结构冻结` | 技术研发、质量环境健康投入 |
| `quarter.receivableUpdate` | `B35:O37` | `quarterMap<object>` | `Q1/Q2/Q3/Q4` | `结构冻结` | 供应链订单、应收账款更新与回款记录 |
| `quarter.deliverySettlement` | `B38:O40` | `quarterMap<object>` | `Q1/Q2/Q3/Q4` | `结构冻结` | 包含交货销售额、交货成本、管理人员费用 |
| `yearEnd.longTermLoan` | `B42:P44` | `object` | `YEAR_END` | `结构冻结` | 长期贷款账期更新、还款与新贷款 |
| `yearEnd.assetAdjustment` | `B45:O53` | `object` | `YEAR_END` | `结构冻结，细项可补充` | 包含维护费、厂房、生产线残值/折旧、新市场培育 |
| `extra.incomeAndPenalty` | `A54:O58` | `quarterMap<object>` | `随当期提交` | `结构冻结` | 包含折现费用、额外支出及罚款、额外收入及奖励 |

补充建议：

- `quarterMap<object>` 推荐统一按 `q1/q2/q3/q4` 四个键组织，避免前后端对列坐标有隐式依赖。
- 若某区块是“按产品/区域/费用类型”的二维输入，推荐以显式对象数组存储，而不是 `G14/I14/K14/M14` 这类坐标命名。
- 经营页中的 `periodEndCash`、`quarterCashCheck`、`planRevenue`、`comprehensiveCostTotal` 等结果值应由规则层计算后回填视图，不建议前端自行算。

### 4.7.1 页面关键标签对照（当前已落前端）

| 页面 | 系统字段 / 区块 | Excel 当前显示名 | 当前前端显示名 | 说明 |
|---|---|---|---|---|
| 管理端初始基线 | `baselineFactoryAsset` | 厂房 | 厂房 | 已按最终版 Excel 收口 |
| 管理端初始基线 | `baselineLineResidual` | 生产线残值 | 生产线残值 | 已按最终版 Excel 收口 |
| 管理端初始基线 | `baselineDepreciableAsset` | 待折资产 | 待折资产 | 已按最终版 Excel 收口 |
| 管理端初始基线 | `baselineWorkInConstruction` | 在建生产线 | 在建生产线 | 制造业版新增补齐，作为初始基线资产项；服务版对应显示名为“在建贵宾厅” |
| 管理端初始基线 | `baselineReceivable` | 应收款 | 应收款 | 已按最终版 Excel 收口 |
| 管理端初始基线 | `baselineRawMaterials` | 材料 | 材料 | 已按最终版 Excel 收口 |
| 管理端初始基线 | `baselineShortTermLoan` | 短期负债 | 短期负债 | 已按最终版 Excel 收口 |
| 管理端初始基线 | `baselineLongTermLoan` | 长期负债 | 长期负债 | 已按最终版 Excel 收口 |
| 管理端初始基线 | `baselineShareCapital` | 股东资本 | 股东资本 | 已按最终版 Excel 收口 |
| 管理端初始基线 | `baselineRetainedEarnings` | 利润留存 | 利润留存 | 已按最终版 Excel 收口 |
| 玩家财报页 | `rawMaterials` | 材料 | 材料 | 绿色手工项 |
| 玩家财报页 | `incomeTaxRate` | 所得税税率 | 所得税税率 | 下拉值 `0.25 / 0.15 / 0` |
| 玩家经营页 | `quarterCashCheckLabel` | 核对季末现金 | 核对季末现金 | 只读展示标签 |
### 4.8 与接口/数据库的直接落点

| 文档对象 | 对应映射章节 | 说明 |
|---|---|---|
| `player-report.save-draft / submit` 的 `reportManualPayload` | `4.1` | 财报手工项应直接使用本节字段名 |
| `sg_initial_baseline.baseline_payload_json` | `4.2` | 管理员初始基线建议直接采用本节字段名 |
| `report_computed_payload_json` | `4.3` | 财报自动计算结果建议至少覆盖本节字段 |
| `sg_group_summary_snapshot` | `4.4` | 汇总表正式字段直接对应收入、利润、权益 |
| `carry_forward_rules` | `4.5` | 跨年承接必须使用本节语义，不直接写死年份表名 |
| `player-operating` 的 `operatingPayload` | `4.7` | 经营页负载按区块化结构落库和回显 |
| `player-order` 的市场投入与选单结果 | `4.6.1 / 4.6.2` | 正式年份订单前置流程、市场投入汇总与订单总额来源 |
| `admin-order` 的订单数量控制与订单池 | `4.6.1` | 管理员订单管理页面和订单池落库字段 |

---

## 5. 状态与页面行为映射

| 结构性状态 | 页面含义 | 前端行为 |
|---|---|---|
| `LOCKED` | 年份未开放 | 标签可见但不可编辑 |
| `OPERATING + Q1_OPEN` | 当前处于 Q1 | 年初区、Q1 区可编辑 |
| `OPERATING + Q2_OPEN` | 当前处于 Q2 | Q2 区可编辑，其余历史区只读 |
| `REPORT_PENDING/REPORTING` | 财报开放 | 财报页可进入 |
| `COMPLETED` | 本年完成 | 本年全部只读 |
| `BANKRUPT` | 已破产 | 后续年份不可编辑 |

---

## 6. 当前尚待补齐的内容

以下内容当前不建议硬编码为最终版：

- 经营页全部中文标签清单
- 经营页逐格显示名称
- 少量仍可能调整的表头名称
- 帮助文案和提示语中的 Excel 原词

如果后续要进一步贴近 Excel 原始版式，可以继续补一张“逐格展示映射表”，字段建议如下：

| 字段 | 说明 |
|---|---|
| `pageCode` | 页面编码 |
| `blockCode` | 区块编码 |
| `fieldCode` | 系统字段编码 |
| `excelSheetName` | Excel 页名 |
| `excelCell` | Excel 单元格 |
| `currentDisplayName` | 当前显示名 |
| `uiDisplayName` | 页面显示名 |
| `inputType` | 输入类型 |
| `ownerRole` | 玩家/管理员/系统 |
| `status` | 已冻结/待细化 |

---

## 7. 变更控制原则

- 若只是标签名称变化，优先修改本映射文档和前端展示文案。
- 若变化影响字段语义、计算口径、跨年承接关系，必须同步更新：
  - `docs/requirements_spec.md`
  - `docs/api_design.md`
  - `docs/database_design.md`
  - `docs/calculation_rule_spec.md`
- 禁止因为 Excel 某个中文标题变化，就直接重命名数据库字段。


