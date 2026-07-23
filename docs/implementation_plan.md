# 沙盘经营系统需求驱动实施计划（首版）

> 更新日期：2026-07-23
> 适用方式：基于当前已确认的业务共识、Excel 规则底稿和原型方向，持续把沙盘经营系统首版需求拆成可执行任务，并同步更新状态。

## 计划规则

- 本计划只收录“首版业务需求任务”与“首版前置落地任务”。
- 每个任务必须包含：优先级、状态、当前产出（文档/原型/代码落点）、下一步。
- 状态定义：`已完成` / `部分完成` / `未开始` / `阻塞`。
- 需求来源以用户明确指令、Excel 规则底稿、现有 Excel 文件为准。
- 未经确认，不得擅自扩大首版范围。

## 任务总览（需求沉淀层当前基线）

- 总任务数：`82`
- 已完成：`80`
- 部分完成：`4`
- 未开始：`0`

## 开发执行层任务总览

- 总任务数：`76`
- 已完成：`75`
- 部分完成：`1`
- 未开始：`0`
- 阻塞：`0`
- 当前状态：`I13-02 ~ I13-07、I13-09 已完成；I13-08 订单专项造数命令已进入手动版实现，待继续完善；等待用户按订单验收清单复测`

---

## 1) 规则底座与需求基线

### 1.1 Excel 规则底稿整理

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 提取 Excel 工作簿结构、关键公式、跨年联动与汇总来源 | P0 | ??? | `game doc/Excel计算规则与跨表联动说明.md` | 后续若业务规则变动，同步标注与 Excel 的差异点 |
| 明确经营表、财报表、汇总表三层关系 | P0 | ??? | `game doc/Excel计算规则与跨表联动说明.md` 第 3~7 节 | 在数据模型设计阶段把这三层映射成实体与状态 |
| 确认财报手工输入边界（在制品/成品/材料/税率） | P0 | ??? | `game doc/Excel计算规则与跨表联动说明.md` 第 8 节 | 在正式页面中映射为绿色输入格和税率下拉 |
| 标记 Excel 中的疑点和待确认项 | P1 | ??? | `game doc/Excel计算规则与跨表联动说明.md` 第 10 节 | 继续在需求讨论中逐条消化 |
| 跟踪 Excel 单元格名称调整并冻结字段显示名称基线 | P0 | ??? | 已基于 `1组 最终版.xlsx` 冻结当前名称基线，并同步到字段映射文档 | 后续继续做全页标签收口与页面/接口文案核对 |

### 1.2 游戏规则与首版边界收敛

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 确认 `0年` 定位：非排名年，但参与后续数据传递 | P0 | ??? | 已在对话中确认；待写入正式需求章节 | 同步到页面导航说明和状态机设计 |
| 确认胜利规则：最终以所有者权益最高为准 | P0 | ??? | 用户确认 + 去年规则文档补充佐证 | 写入正式需求文件的“胜负与排名”章节 |
| 确认破产规则：现金流断裂或最终年度所有者权益为负 | P0 | ??? | 用户确认 | 写入系统状态与页面行为章节 |
| 确认首版只做 Excel 已承载规则 + 少量必要流程控制 | P0 | ??? | 用户确认 | 作为范围闸门，限制后续扩张 |
| 去年规则文档中的贷款上限、认证、小巨人、项目交付罚则等是否纳入首版 | P1 | ??? | 用户已明确：首版暂不纳入 | 作为“二期候选能力”记录即可 |

### 1.3 需求共识沉淀

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 沉淀《需求共识清单 v0.1》，显式区分“已确认 / 待确认 / 不纳入首版” | P0 | ??? | `docs/requirements_consensus_checklist.md` | 作为正式需求文档的前置基线继续使用 |
| 输出《正式需求文档（首版）》并收口首版范围、规则、状态与验收基线 | P0 | ??? | `docs/requirements_spec.md` | 进入技术选型、接口与数据模型设计阶段 |

---

## 2) 玩家端需求

### 2.1 玩家端页面方向

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 确认玩家端采用高保真 Excel 风格页面 | P0 | ??? | 用户确认；`excel-demo.html` | 以 `excel-combined-demo.html` 作为当前主方向继续微调 |
| 确认首页默认进入 `0年经营` | P0 | ??? | 用户确认；`excel-combined-demo.html` | 在正式需求中写为默认导航行为 |
| 确认经营页和财报页都采用整张 Excel 风格 | P0 | ??? | 用户确认；`excel-combined-demo.html` | 后续把财报具体位置继续向真实 Excel 靠近 |
| 确认右侧只保留轻量工作栏 | P1 | ??? | `excel-combined-demo.html` | 后续细化工作栏在经营页/财报页的差异化内容 |
| 输出正式前端页面结构清单（以讨论共识为准，不以 Demo 反推） | P0 | ??? | `docs/frontend_page_structure.md` | 后续前端工程按该清单落页面边界、路由与区块结构 |

### 2.2 玩家端经营页

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 经营页完整展示全部阶段，不按阶段拆成多页 | P0 | ??? | 用户确认；`excel-combined-demo.html` | 在正式需求中定义“同页锁定”机制 |
| 未开放阶段输入格可见但锁定不可编辑 | P0 | ??? | `excel-combined-demo.html` 已示意锁定态 | 后续细化锁定态颜色与只读提示 |
| 年初区随 Q1 一起提交，不单独做年初提交 | P0 | ??? | 用户确认 | 在提交流程章节写入 |
| Q1/Q2/Q3/Q4/年末分别提交 | P0 | ??? | 用户确认 | 在状态机和按钮行为里落地 |
| 经营页自动保存草稿，提交为单独动作 | P0 | ??? | 用户确认；采用每 5 分钟自动保存草稿的轻量方案 | 后续只需在正式需求中写清“草稿保存 vs 正式提交”的区别 |
| 经营页展示黄色“季末现金核对”区域，供玩家线下自行核对 | P0 | ??? | 用户确认 + 最新版 Excel 已体现 | 在正式需求、接口说明和前端 demo 中保持一致口径 |
| 玩家端通知区置于右侧工作栏顶部；普通通知由管理员发送；奖励/罚款由管理员按组/年下发，系统自动归属阶段，玩家只读 | P0 | ✅ | `I12-01 ~ I12-06` 已完成 | 按 `docs/adjustment_module_acceptance_checklist.md` 做业务场景复测 |

### 2.3 玩家端财报页

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 财报页仅在本年经营结束后开放 | P0 | ??? | 用户确认 | 写入状态流转与页面开放条件 |
| 财报页自动带出计算项，绿色单元格为手工项 | P0 | ??? | 用户确认；`excel-combined-demo.html` 已体现 | 后续与真实 Excel 位置继续贴齐 |
| 所得税税率采用下拉框（0.25 / 0.15 / 0） | P0 | ??? | 用户确认；`excel-combined-demo.html` 已体现 | 写入字段交互规范 |
| 财报提交时强校验“总资产 = 总负债 + 总负债权益” | P0 | ??? | 用户确认；demo 已体现校验区 | 在正式需求中定义失败提示文案 |

---

## 3) 阶段流程与提交流程

### 3.1 年份与阶段推进

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 季度由小组自行提交推进 | P0 | ??? | 用户确认 | 进入状态机设计 |
| 年度由管理员统一开放下一年 | P0 | ??? | 用户确认 | 进入管理员页需求 |
| `0年` 与正式年份流程一致，只是现场主持人会线下讲解 | P0 | ??? | 用户确认 | 在需求中明确“系统流程不特殊化” |
| 财报提交成功后，本年度完成 | P0 | ??? | 用户确认 | 与汇总更新时间联动写清楚 |
| 汇总页从 `1年` 开始显示，不显示 `0年` | P0 | ??? | 基于用户共识 + Excel 汇总结构 | 写入管理员汇总定义 |

### 3.2 提交校验规则

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 阶段提交时仅校验当前阶段范围内的手工必填项 | P0 | ??? | 用户确认 | 在字段与阶段校验章节落地 |
| `Q1 提交` 校验 `年初 + Q1` | P0 | ??? | 用户确认 | 写入明确示例 |
| `Q2/Q3/Q4` 提交只校验各自阶段，不重复校验已锁定阶段 | P0 | ??? | 用户确认 | 写入明确示例 |
| `年末提交` 只校验年末区手工项 | P0 | ??? | 用户确认 | 写入明确示例 |
| 经营阶段提交后返回系统计算的期末现金，供玩家自行线下核对 | P0 | ??? | 用户确认 + 最新版 Excel 的“核对季末现金”行 | 在接口响应和页面回显中统一“显示、不拦截”口径 |
| 系统不再要求额外输入现场现金，也不以现金一致性作为提交拦截条件 | P0 | ??? | 用户确认 | 同步移除旧“现金校验闸门”表述 |
| 默认值为空，`0` 必须显式填写 | P0 | ??? | 用户确认 | 在输入规则章节写清楚 |
| 首版不做重型校验清单，只做顶部错误提示 + 字段高亮 | P1 | ??? | 用户确认 | 在交互规范章节说明 |

---

## 4) 管理员端需求

### 4.1 管理员端范围

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 确认管理员端首版菜单：汇总 / 年度控制 / 初始基线 / 组数据 | P0 | ??? | 用户确认 | 进一步拆管理员页面功能点 |
| 汇总页年度汇总区展示正式年份收入/利润/权益与破产状态 | P0 | ??? | 用户确认（上午讨论） | 在正式需求中明确“无实时排名，只展示正式汇总结果” |
| 汇总页与最终排名区同页展示，最终排名区位于页面下方，未满足开放条件时显示状态提示 | P0 | ??? | 用户确认（上午讨论） | 在管理员前端信息架构与接口设计中保持同页结构 |
| 不做管理员进度大盘，不做导出 | P0 | ??? | 用户确认 | 作为首版范围闸门 |
| 初始基线采用共享模板单页录入，提交后按组扇出并锁定不可修改 | P0 | ??? | 用户确认（上午讨论） | 写入初始化流程和数据落库语义 |

### 4.2 年度控制与异常处理

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| `最终年份` 作为配置项 | P0 | ??? | 用户确认 | 在管理员“年度控制”中定义 |
| `当前开放年份` 作为运行时状态 | P0 | ??? | 用户确认 | 与“最终年份”分开建模 |
| 管理员手动开放下一年 | P0 | ??? | 用户确认 | 写入状态流转规则 |
| 管理员不做通用历史纠错引擎，正式版做“定向异常解锁” | P0 | ??? | 用户确认 | 区分经营页 / 财报页解锁，并补阶段级回退 |
| 经营页异常解锁支持按目标阶段回退，后续结果失效但保留草稿 | P0 | ??? | 用户确认 | 结合财报失效、汇总撤回与失效草稿规则写入状态机 |
| 异常解锁只允许在下一年尚未开放前执行 | P0 | ??? | 用户确认 | 写入限制条件 |
| 解锁必须填写原因并保留记录 | P1 | ??? | 用户确认 | 写入审计要求 |

---

## 5) 破产、状态与权限

### 5.1 破产后的系统行为

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 破产后后续所有年份和财报不可再填写 | P0 | ??? | 用户确认 | 写入状态机 |
| 破产后历史数据保留，不清空 | P0 | ??? | 用户确认 | 写入数据保留规则 |
| 管理员汇总页标记“已破产” | P0 | ??? | 用户确认；采用“经营状态列 + 红色已破产标签” | 在正式需求中固化展示规范 |
| 玩家仍可查看本组历史数据，但不可编辑 | P0 | ??? | 用户确认 | 写入权限规则 |

### 5.2 当前仍需细化的状态设计

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 定义“年份状态 + 阶段状态 + 财报状态 + 破产状态”的最小状态机 | P0 | ??? | `docs/minimal_state_machine.md` | 在正式需求文档中引用并细化接口/数据模型影响 |
| 定义“异常解锁后重新提交”的状态回收方式 | P0 | ??? | 用户已确认：覆盖生效、保留操作/提交日志；若已提交财报则财报失效、汇总撤回 | 在正式状态机和需求文档中固化 |

---

## 6) 原型与交互验证

### 6.1 已完成原型

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 游戏化概念页探索 | P2 | ??? | `demo.html` | 仅保留作早期方向记录 |
| Excel 风格经营页高保真探索 | P0 | ??? | `excel-demo.html` | 作为第一轮页面参考 |
| 右侧窄栏布局探索 | P1 | ??? | `excel-sidebar-demo.html` | 作为工作栏探索记录 |
| 主表 + 轻量右侧工作栏混合方案 | P0 | ??? | `excel-combined-demo.html` | 作为当前主方向继续微调 |

### 6.2 原型待收口项

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 经营页再进一步贴近真实 Excel 的排布细节 | P1 | ?? | `excel-combined-demo.html` | 继续调整局部位置、留白与色带 |
| 财报页再进一步贴近真实 Excel 的原始位置和说明区 | P1 | ?? | `excel-combined-demo.html` | 与原 Excel 对照微调 |
| 右侧工作栏在经营页/财报页的差异化内容 | P1 | ??? | 用户已确认：同布局、分内容 | 后续同步到原型与正式需求文档 |

---

## 7) 技术方案前置判断

### 7.1 当前判断

| 任务 | 优先级 | 状态 | 当前产出（文档/原型落点） | 下一步 |
|---|---|---|---|---|
| 判断系统不是“纯 Excel 网页化”，而是“Excel 风格 + 中等复杂度流程状态机” | P0 | ??? | 用户已确认 | 在正式需求文档中写清取舍理由 |
| 评估是否采用轻中型方案（页面轻，规则中等） | P0 | ??? | 用户已确认：采用轻中型方案，业务字段存储 + 页面映射 | 后续再细化技术栈与模块边界 |
| 输出《技术选型细化文档（Go 方向） v0.1》 | P0 | ??? | `docs/technical_selection.md` | 作为接口设计、数据库设计和 Go 实现的承接基线 |
| 输出《接口设计文档（正式版）》 | P0 | ??? | `docs/api_design.md` | 作为 Go Handler、权限和测试用例的直接依据 |
| 输出《数据库设计文档（正式版）》 | P0 | ??? | `docs/database_design.md` | 作为 MySQL DDL、迁移脚本和 Repository 设计依据 |
| 明确玩家端、管理员端、规则层三层解耦 | P1 | ??? | `docs/technical_selection.md` 已明确前后端分离、规则层独立、状态机显式建模 | 下一步在接口与表结构中落地 |

---

## 8) 本期明确不做

| 项目 | 原因 |
|---|---|
| 通用历史回退引擎 | 首版仅做定向异常解锁，不做任意快照回退 |
| 历史纠错引擎 | 首版过重，后续再议 |
| 导出功能 | 用户已明确首版不做 |
| 复杂排名图表与大屏 | 首版仅做基础汇总 |
| 去年规则文档中的认证/小巨人/贷款上限等完整电子化 | 当前不属于首版项目范围 |
| 桌游流程完全电子化 | 当前只做经营与财报系统，不做桌面流程全替代 |

---

## 9) 开发执行拆解（可直接领取）

### 9.1 拆解原则

- 本节任务用于承接 `docs/requirements_spec.md`、`docs/api_design.md`、`docs/database_design.md` 与 `docs/minimal_state_machine.md`，作为实际开发排期与领取任务的依据。
- 每个任务尽量做到“一项任务对应一块明确交付物”，避免出现“做完后仍不知道算不算完成”的模糊任务。
- 首版先保证主链路打通，再补体验增强；当前已以 `1组 最终版.xlsx` 冻结命名基线，后续命名收口任务必须显式回写同步范围，不再以“等待最终版”为由长期挂起。
- 任务状态以本文当前回写为准；后续开发过程中必须继续随做随更新，避免再次出现“文档进度落后于仓库实际状态”。

### 9.2 里程碑划分

| 里程碑 | 目标 | 说明 |
|---|---|---|
| M1 | 打通首版最小可运行主链 | 后端骨架、数据库、玩家经营、财报、管理员控制主链可跑通 |
| M2 | 补齐管理员运维能力与前端页面落地 | 汇总、组数据查看、回退与修正、Excel 风格页面联调完成 |
| M3 | 收口测试与 Excel 最终命名适配 | 测试补齐、字段命名映射、上线前回归 |

### 9.3 M1：首版最小可运行主链

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| M1-01 | 初始化 Go 单体项目骨架 | P0 | ??? | `docs/technical_selection.md` | `cmd/`、`internal/`、`configs/`、`migrations/`、`docs/README.md` | 项目可本地启动，具备基础目录、配置文件样例和启动入口 | 无 |
| M1-02 | 落统一响应体、错误码、日志与基础中间件 | P0 | ??? | M1-01、`docs/api_design.md` | 公共 `response`、`errors`、`middleware`、`logger` 模块 | 新增一个示例接口即可返回统一格式，错误码与日志链路可用 | 无 |
| M1-03 | 建立数据库迁移基线与首版核心表 DDL | P0 | ??? | M1-01、`docs/database_design.md` | 首版 migration 脚本、建表 SQL、回滚脚本 | 能一键初始化 `sg_*` 表结构，核心索引建立成功 | 无 |
| M1-04 | 建立账号、角色、组与基础种子数据 | P0 | ??? | M1-03 | 初始化数据脚本、管理员账号、可配置组账号种子、权限种子 | 本地环境可直接登录管理员与任一已初始化玩家组账号 | 无 |
| M1-05 | 实现年份状态、阶段状态、财报状态的状态机服务 | P0 | ??? | M1-02、M1-03、`docs/minimal_state_machine.md` | `state_machine` 服务与状态枚举常量 | 能正确判断“可编辑 / 只读 / 可提交 / 不可提交 / 可解锁” | 无 |
| M1-06 | 搭建规则层壳子与计算上下文对象 | P0 | ??? | M1-02、M1-03、`game doc/Excel计算规则与跨表联动说明.md` | `rules` 模块、计算输入输出结构、规则注册入口 | 后续经营计算、财报计算、汇总计算有统一调用入口 | 无 |
| M1-07 | 实现获取经营页视图接口 | P0 | ??? | M1-05、M1-06、`docs/api_design.md` 6.2.1 | `player-operating` 查询接口 | 可返回当前组某年的经营页视图、锁定态和字段值 | 无 |
| M1-08 | 实现保存经营页草稿接口 | P0 | ??? | M1-07 | `player-operating` 草稿保存接口 | 草稿可保存，不推进状态，不写正式提交日志 | 无 |
| M1-09 | 实现经营阶段提交接口 | P0 | ??? | M1-05、M1-06、M1-07、M1-08 | `submit-stage` 接口、阶段提交日志写入 | 完成当前阶段必填校验、期末现金计算回显、状态推进、破产判定与阶段锁定 | 无 |
| M1-10 | 实现获取财报页视图接口 | P0 | ??? | M1-05、M1-06、M1-09 | `player-report` 查询接口 | 仅在本年经营结束后开放，自动带出计算结果与手工项 | 无 |
| M1-11 | 实现保存财报草稿接口 | P0 | ??? | M1-10 | `player-report` 草稿保存接口 | 财报手工项可保存，不改年度完成态 | 无 |
| M1-12 | 实现财报提交接口与汇总快照写入 | P0 | ??? | M1-10、M1-11、M1-06 | 财报提交接口、`sg_group_report_submission`、`sg_group_summary_snapshot` 写入逻辑 | 完成平衡校验、财报锁定、年度完成、正式年份汇总入表 | 无 |

### 9.4 M2：管理员主链与前端落地

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| M2-01 | 实现游戏配置接口与最终年份管理 | P0 | ??? | M1-03、M1-05、`docs/api_design.md` 6.1/6.5 | `game-config`、`admin-control` 配置接口 | 管理员可读取与修改最终年份，系统能返回当前开放年份状态；当最终年份上调时，系统可自动补齐未来年份状态空间 | 无 |
| M2-02 | 实现初始基线读取与提交接口 | P0 | ??? | M1-03、M1-05、M1-06 | `initial-baseline` 接口 | 管理员可录入共享基线模板，提交后按组扇出并锁定不可修改 | 无 |
| M2-03 | 实现管理员开放下一年接口 | P0 | ??? | M1-05、M2-01、M2-02 | `open-next-year` 接口与状态推进逻辑 | 所有组在管理员动作后统一开放下一年 | 无 |
| M2-04 | 实现管理员汇总与最终排名接口 | P0 | ??? | M1-12、M1-04 | `admin-summary` 年度汇总/最终排名接口 | 前端可按正式年份顺序展示年度汇总区，并在满足条件后展示最终排名区 | 无 |
| M2-05 | 实现管理员组数据查看接口 | P1 | ??? | M1-07、M1-10 | `admin-group-data` 查询接口 | 管理员可按组、按年查看经营页与财报页 | 无 |
| M2-06 | 实现退回重提接口与审计日志 | P0 | ??? | M1-05、M1-09、M1-12、M2-03 | `unlock-year` 接口、`sg_admin_unlock_log`、`sg_admin_action_log`、回退日志与安全快照 | 满足“仅下一年未开放前可退回重提、原因可空并由系统默认记录、财报失效、汇总撤回” | 无 |
| M2-07 | 搭建玩家端经营页 Excel 风格静态壳子 | P0 | ??? | 现有 demo、`docs/requirements_spec.md` | 玩家经营页前端骨架、年份标签、阶段表格框架 | 当前正式前端页面已作为经营页视觉基线收口；整页 Excel 风格、年份标签、经营/财报切换、右侧轻量工作栏与“未来阶段可见但锁定”口径已落地 | 无 |
| M2-08 | 接入经营页查询、草稿、阶段提交链路 | P0 | ??? | M1-07、M1-08、M1-09、M2-07 | 经营页联调页面 | 已完成真实读取、草稿保存、错误提交拦截与成功提交流程闭环；并通过事务级集成测试验证 `Q1 -> Q2` 推进、阶段流水写入与可编辑范围切换 | 无 |
| M2-09 | 搭建玩家端财报页 Excel 风格静态壳子 | P0 | ?? | 现有 demo、`docs/requirements_spec.md` | 玩家财报页前端骨架 | 财报页已接入正式前端工程并具备 Excel 风格主表、绿色手工项与税率下拉；仍待继续按最终视觉口径验收 | 无 |
| M2-10 | 接入财报查询、草稿、提交链路 | P0 | ??? | M1-10、M1-11、M1-12、M2-09 | 财报页联调页面 | 财报查询、草稿、提交与平衡校验链路已接入正式前端工程，并通过事务级集成测试验证草稿持久化、正式年份提交、汇总快照写入与提交后只读 | 无 |
| M2-11 | 搭建管理员页面基础框架 | P1 | ??? | M2-04、现有 demo | 管理员导航、汇总页、年度控制页、初始基线页 | 管理员正式前端页面已落入 `frontend/` 工程，完成左侧菜单、汇总页、年度控制页、初始基线页与真实接口接入，并通过 `npm run build` | 无 |
| M2-12 | 搭建组数据查看前端并统一回退入口 | P1 | ??? | M2-05、M2-06、M2-11 | 组数据只读查看页、`回退与修正` 统一入口 | 管理员正式前端支持按组按年查看经营/财报只读视图、切换页面类型；退回重提统一在 `回退与修正` 页面发起，组数据页不再保留异常解锁按钮 | 无 |
| M2-13 | 实现最小登录认证接口与认证中间件 | P0 | ??? | M1-02、M1-04、`docs/api_design.md` 6.0 | `auth/login`、`auth/get-current-user`、`auth/logout`、认证中间件、密码校验 | 已补齐最小登录接口、Bearer 鉴权中间件、密码校验与登录态解析；业务接口默认按正式登录态识别 `ADMIN/GROUP`，未登录访问返回 `401`，并通过 `go test ./...` | 无 |
| M2-14 | 实现统一登录页与角色自动跳转 | P0 | ??? | M2-13、`docs/frontend_page_structure.md` | `/sandbox-game/login`、登录态存储、路由守卫、退出登录入口 | 已落统一登录页、登录态持久化、角色自动跳转与退出登录；玩家默认进入 `0年经营`，管理员默认进入 `汇总页`，并通过 `frontend/npm run build` | 无 |

### 9.5 M3：测试、收口与最终命名适配

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| M3-01 | 为状态机、经营校验、季末现金计算、财报平衡补单元测试 | P0 | ??? | M1-05、M1-06、M1-09、M1-12 | 规则层单元测试 | 已覆盖规则层、状态机、异常解锁与管理员控制链关键边界，并可执行 `go test ./...` | 无 |
| M3-02 | 为核心接口补集成测试 | P0 | ??? | M1、M2 全部主链任务 | API 集成测试 | 已补玩家经营/财报、管理员开年、异常解锁、年度汇总与最终排名场景，并补齐管理员 `.http` 用例 | 无 |
| M3-03 | 编写首版联调脚本与人工验收用例 | P0 | ?? | M2-08、M2-10、M2-11、M2-12 | 联调清单、验收步骤、演示账号说明 | 已补联调运行手册与 `.http` 冒烟脚本；仍待按演练库实际走完一次完整主链 | 无 |
| M3-04 | 建立 Excel 对账样例集 | P1 | ??? | M1-06、M1-09、M1-12 | 一组或多组样例输入与预期输出 | 已形成可回归的财报平衡样例与最终排名样例集 | 无 |
| M3-05 | 补字段显示名称映射文档 | P0 | ??? | `1组 最终版.xlsx`、`docs/excel_field_mapping.md` | 字段命名映射表、页面标签对照表、接口注释补充 | 已补齐关键页面标签对照，并明确最终版 Excel 与前端显示名映射 | 无 |
| M3-06 | 收口页面标签、接口注释与帮助文案 | P1 | ??? | M3-05 | 前后端文案更新、接口文档注释更新 | 已将管理端初始基线页、映射文档、接口注释与需求文档收口到最终版 Excel 命名 | 无 |

### 9.6 推荐开发顺序

1. 先做 `M1-01 ~ M1-06`
   - 先把工程底座、表结构、状态机、规则层壳子搭起来，否则后续接口都没有稳固落点。
2. 再做 `M1-07 ~ M1-12`
   - 这是玩家主链，决定这个系统最核心的“能不能用”。
3. 然后做 `M2-01 ~ M2-06`
   - 管理员控制链是游戏现场能否正常推进的关键。
4. 再做 `M2-07 ~ M2-14`
   - 页面联调与操作闭环完成后，才能拿去做现场可调试初版。
5. 最后做 `M3-01 ~ M3-06`
   - 测试和命名适配放最后收口，但单元测试建议边开发边补，不要真的全部压到最后一天。

### 9.7 当前可立即开工的最小任务包

如果要面向“正式版给真实玩家使用”的当前目标，在现有仓库进度基础上，下一批最值得继续推进的是下面这 5 个任务：

1. `I2-01` 补齐“通知与奖惩”接口 / 表结构正式设计
2. `I2-02` 建立通知与奖惩存储模型、迁移脚本与 Repository
3. `I2-03` 实现管理员端普通通知与奖惩下发接口
4. `I2-04` / `I2-05` 打通玩家端通知读取与奖惩计算联动
5. `I2-06` / `I2-07` / `I2-08` 完成管理端页面、玩家端展示与测试收口

## 10) 变更记录

| 日期 | 变更内容 |
|---|---|
| 2026-03-24 | 新建首版实施计划；纳入 Excel 规则底稿、玩家端/管理员端共识、破产规则、异常解锁规则、页面原型方向。 |
| 2026-03-24 | 新增 `docs/requirements_consensus_checklist.md` 任务落点；修正任务统计；把“待确认事项需显式写清”纳入当前基线。 |
| 2026-03-24 | 收口原待确认项；确认轻量自动保存、轻中型方案、异常解锁后的财报失效与汇总撤回规则；新增 Excel 命名调整外部依赖。 |
| 2026-03-24 | 新增 `docs/minimal_state_machine.md`；完成最小状态机沉淀，并将其从计划中的部分完成项转为已完成。 |
| 2026-03-24 | 新增 `docs/requirements_spec.md`；按参考规范沉淀正式需求文档，收口首版范围、规则、状态与验收基线。 |
| 2026-03-24 | 新增 `docs/technical_selection.md`；明确 Go 单体后端、MySQL、业务字段建模、自定义 Excel 风格前端与显式状态机的技术方向。 |
| 2026-03-24 | 新增 `docs/api_design.md` 与 `docs/database_design.md`；按正式版口径定义首版接口边界、动作语义、表清单、状态落库与审计日志设计。 |
| 2026-03-24 | 将原“本周执行建议”升级为“开发执行拆解”；新增 30 项可直接领取的开发任务、里程碑划分、最小可开工任务包与 Excel 最终命名阻塞标记。 |
| 2026-03-26 | 补充 `docs/api_design.md` 中 `admin-control` 六个接口的 Go DTO 草图、响应字段与错误常量建议；同步把“共享初始基线模板”“开放下一年阻塞原因摘要”“异常解锁最小硬校验”回写到正式文档。 |
| 2026-03-26 | 同步更新 `docs/database_design.md`：明确 `final_year >= current_open_year`、`sg_initial_baseline` 采用“页面共享模板 + 后端按组扇出落库”的实现口径，并补充异常解锁回收破产标记的落库要求。 |
| 2026-03-26 | 完成正式文档与 `AGENTS.md` 的术语收口：统一使用“初始基线”命名，并修正少量替换后不顺的表述，确保后续编码阶段命名一致。 |
| 2026-03-26 | 复核聊天共识后修正需求口径：经营页改为展示“季末现金核对”供玩家线下自核，系统不再要求额外录入现场现金，也不以现金一致性阻断提交；同步补入上午确认的管理端汇总页单页双区、共享初始基线页等要求。 |
| 2026-03-26 | 复核当前仓库代码后重写开发状态：确认 `M1-01 ~ M1-12` 与 `M2-04` 已落代码；将此前误记为已完成的 `M2-01/M2-02` 回退为未开始；将 `demo 网页/player-demo.html`、`demo 网页/admin-demo.html` 计入静态壳子部分完成。 |
| 2026-03-26 | 同步收口 6 份核心/关键文档旧口径：将 `1组（3.5）.xlsx` 统一切换为 `1组 最终版.xlsx`；移除“现金核对阻断提交”的旧规则，统一为“季末现金展示供玩家线下自核”；并把管理员汇总页同页双区、共享初始基线模板等共识同步回 `AGENTS.md`、需求文档、接口文档、状态机文档与技术选型文档。 |
| 2026-03-26 | 新增 `docs/frontend_page_structure.md`：按上午讨论共识沉淀首版正式前端页面结构清单，明确玩家端/管理端正式页面、必备区块、页面关系与明确不做项；并在 `docs/frontend_guide.md` 中将页面范围基线指向该清单。 |
| 2026-03-26 | 新增面向汇报的前端静态展示原型：`demo 网页/player-demo.html` 展示玩家端“Excel 整页 + 轻量工作栏”，`demo 网页/admin-demo.html` 展示管理端“汇总 / 年度控制 / 初始基线 / 组数据”信息架构。 |
| 2026-03-26 | 初始化 `E:\project\sand box game` 本地 Git 仓库并提交当前项目基线；补充 `.gitignore`，将文档、原型、Go 后端骨架与迁移脚本纳入首个可追溯版本。 |
| 2026-03-26 | 实现 `M2-01` 核心后端链路：补齐 `game-config/get-current`、`admin-control/get-config`、`admin-control/update-final-year`，新增管理员动作日志仓储与最终年份更新审计，并补充对应服务层测试。 |
| 2026-03-26 | 实现 `M2-02` 后端链路：补齐 `admin-control/get-initial-baseline`、`admin-control/submit-initial-baseline`，完成共享初始基线模板读取、按组扇出提交锁定、提交状态回写与管理员动作日志留痕。 |
| 2026-03-26 | 实现 `M2-03` 后端链路：补齐 `admin-control/open-next-year`，完成“仅开放下一自然年、全部未破产组完成当前年财报后才可开放、破产组后续年份保持锁定、当前开放年份更新与管理员动作日志留痕”的首版闭环，并通过 `go test ./...` 验证。 |
| 2026-03-26 | 实现 `M2-05` 管理员组数据查询链路：补齐 `admin-group-data/get-operating-view`、`admin-group-data/get-report-view`，实现管理员按组按年查看经营/财报只读视图，并复用现有玩家端计算与页面数据模型，同时通过 `go test ./...` 验证。 |
| 2026-03-26 | 补齐 `M2-01` 剩余链路：新增 `game-config/get-year-tabs`，按当前登录身份返回年份标签状态，冻结玩家端年份标签的“可进入 / 锁定 / 已完成 / 已破产只读”口径，并通过 `go test ./...` 验证。 |
| 2026-03-27 | 实现 `M2-06` 异常解锁后端闭环：补齐 `admin-control/unlock-year`、`sg_admin_unlock_log` 落库、财报提交失效、汇总撤回、破产恢复与管理员动作审计，并通过 `go test ./...` 验证。 |
| 2026-03-27 | 根据“管理员可动态设定游戏进行多少年”的新收口要求，补充正式需求、接口与前端结构文档：明确前端采用动态年份视图而非逐年生成独立页面，并确认 `update-final-year` 还需补齐“扩年初始化”能力；同步将 `M2-01` 从已完成调整为部分完成。 |
| 2026-03-27 | 完成 `M2-01` 扩年初始化闭环：`admin-control/update-final-year` 在上调最终年份时会自动补齐全部小组缺失的未来年份状态记录，响应中补充扩年范围信息，并通过 `go test ./...` 验证。 |
| 2026-03-27 | 启动正式前端工程：新增 `frontend/` Vue 3 + TypeScript + Vite + Pinia 项目骨架，完成玩家经营页首版页面、年份标签、右侧工作栏和经营页查询 / 草稿保存 / 阶段提交接口接入，并通过 `npm run build` 验证。 |
| 2026-03-27 | 完成经营页首轮功能验收：在线验证 `get-current`、`get-year-tabs`、`get-year-view` 读取正常，确认 `Q1_OPEN` 时仅开放 `YEAR_START + Q1`；真实验证 `save-draft` 可写且会刷新 `lastDraftSavedAt`，并验证未满足条件时 `submit-stage` 返回 `422` 且不会误推进阶段；同时再次通过 `go test ./...` 与 `npm run build`。 |
| 2026-03-28 | 完成 `M2-08` 成功提交闭环：新增 `internal/service/player_operating_command_service_test.go`，用事务级集成测试在不污染真实小组数据前提下验证 `Q1 -> Q2` 推进、阶段提交流水写入、经营页回读后仅开放 `Q2` 编辑，并再次通过 `go test ./...`。 |
| 2026-03-29 | 确认 `game doc/1组 最终版.xlsx` 为当前命名与公式冻结基线；清理 `docs/` 中残留的“待最终 Excel”旧口径，解除 `M3-05/M3-06` 的外部阻塞表述，并修正字段映射摘要中的旧显示名。 |
| 2026-03-29 | 完成 `M3-01 / M3-02 / M3-04 / M3-05 / M3-06`：新增管理员开年 / 异常解锁 / 年度汇总 / 最终排名事务级集成测试，补齐 `AdminControl.http`、`AdminSummary.http`、`docs/excel_reconciliation_samples.md`，并将管理端初始基线页与核心文档命名收口到最终版 Excel。 |
| 2026-03-28 | 推进 `M2-09` 财报页视觉收口：重做 `PlayerReportPage.vue`、`ReportSheet.vue`、`ReportSidebar.vue` 的正式页面壳子，使其更贴近 Excel 主表布局，并再次通过 `npm run build`。 |
| 2026-03-28 | 完成 `M2-10` 财报主链闭环：新增 `internal/service/player_report_command_service_test.go`，用事务级集成测试验证财报草稿保存、正式年份提交、汇总快照写入、提交后只读，并再次通过 `go test ./...`。 |
| 2026-03-28 | 完成 `M2-11` 管理员正式前端框架：新增管理员端路由、左侧菜单、汇总页、年度控制页、初始基线页与组数据路由壳子，接入管理员汇总 / 年度控制 / 初始基线真实接口，并通过 `frontend/npm run build`。 |
| 2026-03-28 | 完成 `M2-12` 组数据查看与异常解锁前端：新增 `admin-group-data/list-groups` 后端接口，管理员正式前端已支持按组按年查看经营/财报只读视图、切换页面类型、发起异常解锁并展示最近一次原因记录，同时再次通过 `go test ./...` 与 `frontend/npm run build`。 |
| 2026-03-31 | 新增“通知与奖惩”需求变更：允许管理员通过普通通知向全体或单组发送自由文本消息，并允许广播赛事事件信息；经营页中的奖励/罚款改为管理员按 `组 + 年 + 季` 下发、玩家只读、自动计入计算；同步回写 `docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/frontend_page_structure.md` 与本计划。 |

## 11) 当前开发进度回写

### 11.1 本轮完成情况（2026-03-27 ~ 2026-03-28）

- `M1-01 ~ M1-12`：已完成
  - 落点：`cmd/server/main.go`、`internal/app/router.go`、`internal/http/handler/player_operating_handler.go`、`internal/http/handler/player_report_handler.go`、`internal/service/player_operating_*`、`internal/service/player_report_*`、`internal/rules/*`、`internal/state/*`、`migrations/mysql/*`
  - 偏差说明：当前仓库实际进度已经超过“最小骨架”阶段，玩家端经营/财报主链、规则层、状态机、迁移脚本和种子数据均已落地；此前计划中的开发状态回写偏保守。
  - 下一步：优先承接管理员异常解锁，以及前端联调。
- `M2-04`：已完成
  - 落点：`internal/http/handler/admin_summary_handler.go`、`internal/service/admin_summary_query_service.go`、`internal/repository/summary_snapshot_repository.go`
  - 偏差说明：当前已完成“按年取正式汇总 + 最终排名”后端查询链路，但管理员前端仍未正式接入，页面层“单页双区”还停留在 Demo/需求确认阶段。
  - 下一步：在 `M2-11` 中把同页汇总区与最终排名区接到正式页面。
- `M2-01`：已完成
  - 落点：`internal/http/handler/game_config_handler.go`、`internal/service/game_config_query_service.go`、`internal/service/game_config_query_service_test.go`、`internal/http/handler/admin_control_handler.go`、`internal/service/admin_control_*`、`internal/repository/game_config_repository.go`、`internal/repository/group_year_state_repository.go`、`internal/repository/group_year_state_repository_test.go`、`internal/repository/admin_action_log_repository.go`、`internal/app/router.go`
  - 偏差说明：当前已补齐 `game-config/get-current`、`game-config/get-year-tabs`、`admin-control/get-config`、`admin-control/update-final-year`，并在管理员上调 `finalYear` 时自动为全部小组补齐缺失的未来年份状态记录；前端后续只需按接口返回的年份标签动态刷新页签上限，不再依赖数据库预置未来年份数据。
  - 下一步：在 `M2-11` 正式管理员页面中接入“最终年份修改后刷新年份标签上限”的交互联动。
- `M2-02`：已完成
  - 落点：`internal/http/handler/admin_control_handler.go`、`internal/http/dto/admin_control_dto.go`、`internal/service/admin_control_query_service.go`、`internal/service/admin_control_command_service.go`、`internal/repository/initial_baseline_repository.go`、`internal/repository/game_config_repository.go`、`internal/app/router.go`
  - 偏差说明：本轮已补齐 `admin-control/get-initial-baseline` 与 `admin-control/submit-initial-baseline`，实现“共享模板读取、提交后按组扇出、提交即锁定、重复提交拦截、管理员动作日志留痕”的首版闭环；当前实现直接复用已存在的 `sg_initial_baseline` 按组存储结构，不额外引入新表。
  - 下一步：作为年度控制前置能力继续复用，并在管理员页面接入初始基线查看/提交链路。
- `M2-03`：已完成
  - 落点：`internal/http/dto/admin_control_dto.go`、`internal/http/handler/admin_control_handler.go`、`internal/service/admin_control_command_service.go`、`internal/repository/game_config_repository.go`、`internal/app/router.go`、`internal/service/admin_control_command_service_test.go`
  - 偏差说明：本轮已补齐 `admin-control/open-next-year`，实现“目标年份必须为当前开放年份 + 1、不得超过最终年份、仅在全部未破产组完成当前年财报后才允许开放、破产组后续年份保持锁定、管理员动作日志留痕”的首版闭环；开放条件摘要继续复用查询层统一口径。
  - 下一步：继续在管理员页面接入“开放下一年”动作，并和年度控制区联调。
- `M2-05`：已完成
  - 落点：`internal/http/dto/admin_group_data_dto.go`、`internal/http/handler/admin_group_data_handler.go`、`internal/service/admin_group_data_query_service.go`、`internal/service/admin_group_data_query_service_test.go`、`internal/app/router.go`
  - 偏差说明：本轮先按任务完成标准落地 `admin-group-data/get-operating-view` 与 `admin-group-data/get-report-view` 两个核心只读接口，管理员现在可按组按年查看经营页与财报页；`page-stage-submissions` 与 `page-report-submissions` 尚未实现，保留在后续页面实际需要时补齐。
  - 下一步：在 `M2-12` 中把“查看组数据 + 异常解锁”串成完整管理操作流。
- `M2-06`：已完成
  - 落点：`internal/http/dto/admin_control_dto.go`、`internal/http/handler/admin_control_handler.go`、`internal/service/admin_control_command_service.go`、`internal/repository/admin_unlock_log_repository.go`、`internal/repository/group_repository.go`、`internal/repository/report_repository.go`、`internal/repository/summary_snapshot_repository.go`、`internal/app/router.go`、`internal/service/admin_control_command_service_test.go`
  - 偏差说明：本轮当时已补齐 `admin-control/unlock-year` 首版闭环，既有代码在命中破产年份时会恢复该组 `NORMAL` 状态；该历史实现已被 2026-07-20 的 `I12-01` 新口径取代，后续 `I12-02 ~ I12-06` 必须移除破产恢复路径，改为不可撤销破产。其余财报提交回收、汇总撤回与审计日志能力继续复用。
  - 下一步：在 `M2-12` 中把解锁弹窗、组数据页和操作反馈联调到正式前端。
- `M2-07`：已完成
  - 落点：`frontend/package.json`、`frontend/src/router/index.ts`、`frontend/src/stores/player-operating.ts`、`frontend/src/views/sandbox-game/player/operating/PlayerOperatingPage.vue`、`frontend/src/components/sandbox-game/player/OperatingSheet.vue`、`frontend/src/components/sandbox-game/player/OperatingSidebar.vue`
  - 偏差说明：当前正式前端中的经营页已经作为视觉基线收口，页面整体样式与结构已被确认为后续迭代基准；接下来不再以“重做页面外观”为重点，而是转向财报页与管理端的正式联调。
  - 下一步：进入 `M2-09` / `M2-10`，继续收口财报页页面与提交链路。
- `M2-08`：已完成
  - 落点：`frontend/src/api/sandbox-game/player-operating.ts`、`frontend/src/stores/player-operating.ts`、`frontend/src/views/sandbox-game/player/operating/PlayerOperatingPage.vue`、`internal/http/handler/player_operating_handler.go`、`internal/service/player_operating_*`、`internal/service/player_operating_command_service_test.go`
  - 偏差说明：当前已完成真实读取、草稿保存、阶段提交接口接入，并新增事务级集成测试，在不污染真实小组数据前提下验证“满足条件后成功从 `Q1 -> Q2` 推进”“阶段提交流水落库”“经营页回读后仅开放 `Q2` 编辑”；同时再次通过 `go test ./...`。
  - 下一步：转入 `M2-09` / `M2-10`，继续财报页正式联调与提交验收。
- `M2-09`：部分完成
  - 落点：`frontend/src/api/sandbox-game/player-report.ts`、`frontend/src/stores/player-report.ts`、`frontend/src/views/sandbox-game/player/report/PlayerReportPage.vue`、`frontend/src/components/sandbox-game/player/ReportSheet.vue`、`frontend/src/components/sandbox-game/player/ReportSidebar.vue`
  - 偏差说明：财报页已继续向高保真 Excel 收口，页头文案、状态区、主表列头 / 标题行 / 平衡校验区、右侧工作栏均已按正式页面方向调整，并再次通过 `npm run build`；当前仍保留“部分完成”，主要是还缺一轮用户视角的页面人工验收。
  - 下一步：继续做用户视角的页面确认，并把最终视觉验收并入 `M3-03` 的联调脚本与人工验收用例。
- `M2-10`：已完成
  - 落点：`frontend/src/api/sandbox-game/player-report.ts`、`frontend/src/stores/player-report.ts`、`frontend/src/views/sandbox-game/player/report/PlayerReportPage.vue`、`internal/http/handler/player_report_handler.go`、`internal/service/player_report_*`、`internal/service/player_report_command_service_test.go`
  - 偏差说明：本轮补齐了财报主链的事务级集成验证，在不污染真实小组数据前提下覆盖“草稿保存持久化”“正式年份财报提交”“年度完成态写回”“正式汇总快照写入”“提交后财报页只读回显”，并再次通过 `go test ./...`；浏览器侧最终视觉确认仍继续归入 `M2-09` / `M3-03` 的人工验收范围，不影响当前功能链闭环完成判定。
  - 下一步：转入 `M2-11` / `M2-12`，开始管理员端正式前端页面与解锁操作流接入。
- `M2-11`：已完成
  - 落点：`frontend/src/router/index.ts`、`frontend/src/views/sandbox-game/admin/AdminLayout.vue`、`frontend/src/views/sandbox-game/admin/summary/AdminSummaryPage.vue`、`frontend/src/views/sandbox-game/admin/control/AdminControlPage.vue`、`frontend/src/views/sandbox-game/admin/baseline/AdminBaselinePage.vue`、`frontend/src/components/sandbox-game/admin/AdminNav.vue`、`frontend/src/api/sandbox-game/admin-*.ts`、`frontend/src/stores/admin-*.ts`、`frontend/src/types/sandbox-game-admin.ts`
  - 偏差说明：本轮已将管理员 Demo 收口为正式 `frontend/` 工程页面，完成左侧菜单、汇总页、年度控制页、初始基线页与组数据路由壳子，并接入管理员汇总 / 年度控制 / 初始基线真实接口，同时通过 `npm run build`。
  - 下一步：转入 `M2-12` / `M3-03`，继续组数据操作流与管理员页面人工验收。
- `M2-12`：已完成
  - 落点：`frontend/src/api/sandbox-game/admin-group-data.ts`、`frontend/src/api/sandbox-game/admin-control.ts`、`frontend/src/stores/admin-group-data.ts`、`frontend/src/views/sandbox-game/admin/group-data/AdminGroupDataPage.vue`、`internal/http/handler/admin_group_data_handler.go`、`internal/service/admin_group_data_query_service.go`、`internal/app/router.go`
  - 偏差说明：本轮补齐了组数据页真实查询与异常解锁主链，新增管理员小组列表接口，前端已支持按组按年查看经营/财报只读视图、切换页面类型、发起异常解锁并展示最近一次原因记录，同时通过 `go test ./...` 与 `npm run build`；当前仍未单独提供“历史解锁日志列表查询”，若后续现场需要，可在管理端继续追加日志查询区。
  - 下一步：转入 `M3-03`，补管理员端联调脚本与人工验收用例。
- `M2-13`：已完成
  - 落点：`internal/http/handler/auth_handler.go`、`internal/service/auth_service.go`、`internal/http/middleware/auth_middleware.go`、`internal/app/router.go`、`internal/repository/account_repository.go`、`internal/service/auth_service_test.go`
  - 偏差说明：本轮在不引入重型账号系统的前提下，基于现有 `sg_account` 和种子账号补齐了最小登录认证闭环；默认认证模式切到项目内 `Bearer` 令牌，同时保留 `auth.mode=bypass` 作为开发态兜底，不影响正式版默认行为。
  - 下一步：转入 `M2-14`，补统一登录页、路由守卫与退出登录联调。
- `M2-14`：已完成
  - 落点：`frontend/src/views/sandbox-game/login/LoginPage.vue`、`frontend/src/stores/auth.ts`、`frontend/src/router/index.ts`、`frontend/src/api/http.ts`、`frontend/src/views/sandbox-game/admin/AdminLayout.vue`、`frontend/src/views/sandbox-game/player/operating/PlayerOperatingPage.vue`、`frontend/src/views/sandbox-game/player/report/PlayerReportPage.vue`
  - 偏差说明：本轮已将统一登录页、登录态持久化、角色自动跳转、页面守卫与退出登录接入正式前端工程，同时移除管理员接口对开发旁路头的依赖；玩家页 `preview=1` 仍仅保留为开发预览入口。
  - 下一步：回到 `M3-03`，按正式登录入口补联调手册并在演练库走完整主链。
- `正式版最小登录闭环`：已完成
  - 落点：`internal/http/handler/auth_handler.go`、`internal/service/auth_service.go`、`frontend/src/views/sandbox-game/login/LoginPage.vue`、`frontend/src/stores/auth.ts`、`frontend/src/router/index.ts`
  - 偏差说明：当前已切换到“统一登录页 + Bearer 登录态 + 角色自动跳转”的正式口径，开发旁路仅保留为可配置兜底与玩家页预览模式，不再是正式入口。
  - 下一步：继续执行 `M3-03`，按正式登录入口完成一次完整联调演练。- `M3-03`：部分完成
  - 落点：`docs/integration_acceptance_runbook.md`、`docs/testing_guide.md`、`tests/http/sandbox-game/SandboxGame-Smoke.http`
  - 偏差说明：本轮已补齐首版联调运行手册、浏览器入口、开发态旁路身份说明与 `.http` 冒烟脚本，联调人员现在可以按统一路径检查玩家端、管理员端与核心只读接口；但“完整写链路现场走查”仍需在演练库按主持节奏再跑一次，暂不记为完全完成。
  - 下一步：在演练库按手册实际走完 `0年 -> 1年 -> 财报 -> 汇总 -> 开放下一年`，并将人工验收结果回填到本计划。
- `M3-01`：已完成
  - 落点：`internal/state/state_machine_test.go`、`internal/rules/operating/*_test.go`、`internal/rules/report/*_test.go`、`internal/rules/carryforward/*_test.go`、`internal/rules/summary/*_test.go`、`internal/service/admin_control_command_service_test.go`、`internal/service/admin_control_command_integration_test.go`
  - 偏差说明：本轮补齐了管理员开年与异常解锁的事务级场景，规则层、状态机和管理员控制链关键边界现在已有自动化回归承接；继续执行 `go test ./...` 可覆盖当前首版主链高风险规则。
  - 下一步：仅保留后续规则变更时随改随补，不再作为单独遗留项。
- `M3-02`：已完成
  - 落点：`internal/service/player_operating_command_service_test.go`、`internal/service/player_report_command_service_test.go`、`internal/service/admin_control_command_integration_test.go`、`internal/service/admin_summary_query_integration_test.go`、`tests/http/sandbox-game/SandboxGame-Smoke.http`、`tests/http/sandbox-game/AdminControl.http`、`tests/http/sandbox-game/AdminSummary.http`
  - 偏差说明：当前已形成“事务级集成测试 + `.http` 联调脚本”双层保障，覆盖玩家主链、管理员开年、异常解锁、年度汇总与最终排名，不再只停留在只读冒烟。
  - 下一步：后续若新增管理端接口，继续按同口径补对应 `.http` 与事务级测试。
- `M3-04`：已完成
  - 落点：`docs/excel_reconciliation_samples.md`
  - 偏差说明：当前先沉淀了财报平衡和最终排名两组高价值对账样例，优先保障首版最关键的结果口径；后续如 Excel 再新增重点公式，可继续扩样例集。
  - 下一步：如规则负责人要求逐格核对，再在本文件继续增补样例编号。
- `M3-05`：已完成
  - 落点：`docs/excel_field_mapping.md`
  - 偏差说明：本轮补齐了管理端初始基线、玩家财报页、经营页季末现金核对等关键标签映射，当前已经能支撑需求、接口和前端共同使用同一套名称基线。
  - 下一步：后续若新增页面标签，只需继续往映射表追加，不再单独开“命名收口”支线。
- `M3-06`：已完成
  - 落点：`frontend/src/views/sandbox-game/admin/baseline/AdminBaselinePage.vue`、`docs/api_design.md`、`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/frontend_page_structure.md`、`docs/calculation_rule_spec.md`、`docs/technical_selection.md`
  - 偏差说明：本轮已把“厂房 / 生产线残值 / 待折资产 / 材料 / 股东资本 / 利润留存”等关键显示名统一到最终版 Excel 口径，并清理掉接口与技术文档中的旧说法。
  - 下一步：后续只需在需求变更时同步更新，不再作为独立阻塞项。
- 文档口径收口：已完成
  - 落点：`AGENTS.md`、`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/api_design.md`、`docs/minimal_state_machine.md`、`docs/technical_selection.md`
  - 偏差说明：本轮主要修正文档之间的旧口径冲突，不改变当前已落地代码的业务边界；重点完成了 Excel 主依据文件名切换、现金规则去闸门化、管理员汇总页双区结构和共享初始基线模板的同步，并补齐了技术选型文档中的旧现金校验表述。
  - 下一步：继续检查其余辅助文档是否还残留旧“现金核对阻断提交”表述，并在后续开发中以这批已同步文档为准。
- 正式前端页面结构清单：已完成
  - 落点：`docs/frontend_page_structure.md`、`docs/frontend_guide.md`
  - 偏差说明：本轮页面清单以上午已确认的业务讨论结论为准，不以现有 Demo 倒推页面；因此它是正式页面边界基线，不是静态原型说明。
  - 下一步：前端工程落地时，优先按本清单拆页面与路由，再用 Demo 仅作视觉参考。
- 仓库初始化与基线提交：已完成
  - 落点：`.gitignore`、本地 Git 提交记录
  - 偏差说明：当前默认终端工作目录失效，因此初始化与提交均通过显式指定 `E:\project\sand box game` 路径执行；不影响仓库内容本身。
  - 下一步：后续开发继续在该仓库内按任务粒度提交，避免再混入旧目录 `E:\project\game` 的历史内容。
- 状态纠偏：已完成
  - 落点：本文档 `9.4` 与 `11.1`
  - 偏差说明：经复核，先前将 `M2-01/M2-02` 的状态记录得不够准确；当前已按仓库实际代码修正为 `M2-01` 已完成、`M2-02` 已完成，因此后续状态回写要继续以当前仓库为准。
  - 下一步：后续新增管理员控制接口时，按实际代码落点继续回写，不再用“曾讨论过/曾做过别处代码”代替当前仓库状态。
- 玩家页乱码修复与 Excel 标签回正：已完成
  - 落点：`frontend/src/components/sandbox-game/player/ReportSheet.vue`、`frontend/src/components/sandbox-game/player/OperatingSheet.vue`、`.editorconfig`
  - 偏差说明：本轮不再以旧提交版本为标签基线，而是直接对照 `game doc/1组 最终版.xlsx` 修复玩家财报页与经营页中的乱码文本，并补回财报页缺失的“非流动资产 / 流动资产 / 负债 / 总权益”等 Excel 标签；同时复核经营页第 `60/61` 行衍生指标绑定，修正了“待折资产增减 / 应收增减”值错位问题，并新增仓库级 UTF-8 编码约束，降低后续再次写入乱码的风险。
  - 下一步：后续若 Excel 显示名继续调整，优先以当前主依据工作簿为准同步玩家页标签，并在需要时继续补充 `docs/excel_field_mapping.md` 的展示映射。
- 经营页 6 个底部指标下沉后端规则层：已完成
  - 落点：`internal/rules/operating/operating_calculator.go`、`internal/service/operating_indicator_helper.go`、`internal/service/player_operating_query_service.go`、`internal/service/admin_group_data_query_service.go`、`frontend/src/components/sandbox-game/player/OperatingSheet.vue`、`frontend/src/utils/sandbox-game-operating-preview.ts`、`internal/rules/operating/operating_calculator_test.go`、`internal/service/operating_indicator_helper_test.go`
  - 偏差说明：本轮将 `市场回报比 / 研发投入强度 / 劳动生产率 / 净产收益率 / 总产收益 / 净利润率 / 毛率润率` 的展示口径统一收回后端；其中 `劳动生产率` 已按 Excel `O59` 固化在经营规则层，`净产收益率 / 总产收益 / 净利润率 / 毛率润率` 改为优先读取“当前年财报结果 / 财报预览”口径；若当前年尚无财报记录，服务端也会按“财报手工项空值视作 0”即时生成预览，确保经营页仍按 Excel 公式链显示，不再本地猜算错误值；同时修正了玩家经营页 store 的本地 preview 合并逻辑，避免后端已返回的财报口径指标被前端二次计算结果覆盖成 `--`；并已补上正式年份跨年联动回归，验证 `1年` 会正确承接 `0年财报` 再生成当前年财报预览与底部收益率指标；自动化回归改为“固定公式样例 + 服务层注入行为”双层校验，不依赖运行时读取 Excel 文件。
  - 下一步：后续若 Excel 再调整底部指标公式，只需同步更新规则层固定样例和文档口径，并继续保持经营页只读后端返回值。
## 12) 正式版迭代池（不计入当前 M1 ~ M3 完成统计）

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| I1-01 | 管理员异常解锁支持目标类型 `OPERATING / REPORT` 与 `targetStageCode` | P0 | 已完成 | `M2-06`、`docs/requirements_spec.md`、`docs/api_design.md`、`docs/minimal_state_machine.md` | `internal/http/dto/admin_control_dto.go`、`internal/http/handler/admin_control_handler.go`、`internal/service/admin_control_command_service.go`、`frontend/src/views/sandbox-game/admin/group-data/AdminGroupDataPage.vue`、`frontend/src/stores/admin-group-data.ts` | 管理员可明确选择“经营页 / 财报页”；当目标为 `OPERATING` 时可选择 `Q1 / Q2 / Q3 / Q4 / YEAR_END`，且服务端按目标阶段回退 | 后端已支持目标类型与目标阶段解锁，管理端弹窗已切到新请求结构，并通过 `go test ./...` 与 `npm run build` 验证 |
| I1-02 | 实现“有效结果 / 失效草稿”分离规则 | P0 | 已完成 | `I1-01` | `internal/assembler/draft_view_state.go`、`internal/repository/report_repository.go`、`internal/service/player_*_service.go`、`internal/service/admin_group_data_query_service.go` | 回退到 `Q2` 时，`Q3 / Q4 / 年末 / 财报` 原值保留但不再计入正式结果；重新提交后重新生效 | 已按“推断式失效草稿”落地：上一年承接仅读取有效财报，经营页可返回 `invalidScopes` / `hasRetainedReportDraft`，财报页可返回 `hasInvalidDraft`，并通过 `go test ./...` 验证 |
| I1-03 | 前端页面与交互收口：异常解锁弹窗、玩家页失效草稿展示、目标化文案 | P1 | 已完成 | `I1-01`、`I1-02` | `frontend/src/components/sandbox-game/player/*`、`frontend/src/views/sandbox-game/player/*`、`frontend/src/views/sandbox-game/admin/group-data/AdminGroupDataPage.vue`、`frontend/src/types/sandbox-game.ts` | 管理员端支持目标类型 + 阶段选择；玩家端能区分“已生效只读 / 当前可编辑 / 失效草稿待重提” | 玩家经营页、财报页与管理端组数据页已补齐失效草稿状态展示、单元格样式与提示文案，并通过 `npm run build` 验证 |
| I2-01 | 补齐“通知与奖惩”正式设计：接口、表结构、错误码与轮询口径 | P0 | 已完成 | `docs/requirements_spec.md`、`docs/frontend_page_structure.md` | `docs/api_design.md`、`docs/database_design.md`、`docs/testing_guide.md`、`tests/http/sandbox-game/AdminNotice.http` | 已补齐 `admin-notice` 接口域、请求/响应、锁定季度 `422` 语义、轮询口径与数据库表设计；后续开发不再需要边写边猜口径 | 已在本轮回写正式文档 |
| I2-02 | 建立通知与奖惩存储模型、迁移脚本与 Repository | P0 | 已完成 | `I2-01`、`M1-03` | `migrations/mysql/0003_notice_adjustment.sql`、Entity、Repository | 已落 `sg_notice` / `sg_group_adjustment` 两张表、迁移脚本与 Repository，并已将 `0003` 幂等执行到本地演练库 | 已完成 |
| I2-03 | 实现管理员端普通通知与奖惩下发接口 | P0 | 已完成 | `I2-01`、`I2-02` | Handler、DTO、Service、Router | 已提供 `get-records`、`send-general`、`send-adjustment` 三个接口；支持全体/单组通知、按组年季奖惩，以及锁定季度 `422` 拦截 | 已完成 |
| I2-04 | 实现玩家端通知读取能力与管理员记录查询能力 | P0 | 已完成 | `I2-01`、`I2-02` | Player Notice Board 装配、管理员记录查询接口 | 玩家经营页/财报页查询结果已附带 `noticeBoard`；管理员端可查看最近普通通知与奖惩记录；首版采用轮询/刷新口径，不做 WebSocket | 已完成 |
| I2-05 | 将奖惩接入经营计算、财报承接与汇总口径 | P0 | 已完成 | `I2-02`、`I2-03`、`M1-09`、`M1-12` | 规则层、查询层、装配层改造 | 奖惩已由后端统一叠加到经营页并继续进入财报承接；玩家提交时奖惩值由服务端覆盖，经营页仅保留只读展示 | 已完成 |
| I2-06 | 实现管理端“通知与奖惩页”前端页面 | P1 | 已完成 | `I2-03`、`I2-04`、`M2-11` | `/sandbox-game/admin/notices` 页面、表单区、记录区 | 管理员正式前端已可发送普通通知、下发奖惩并查看最近记录；页面结构已接入管理端导航 | 已完成 |
| I2-07 | 实现玩家端通知区与奖惩只读展示 | P0 | 已完成 | `I2-04`、`I2-05`、`M2-07`、`M2-09` | 玩家经营页/财报页右侧通知区、经营页奖惩只读展示 | 玩家经营页与财报页右侧工作栏顶部已展示通知区；经营页中的奖励/罚款已改为管理员下发的只读值 | 已完成 |
| I2-08 | 补通知与奖惩测试、联调脚本与验收样例 | P0 | 已完成 | `I2-03`、`I2-04`、`I2-05`、`I2-06`、`I2-07` | `internal/service/admin_notice_integration_test.go`、`tests/http/sandbox-game/AdminNotice.http` | 已补自动化测试、`.http` 联调脚本，并完成本地接口联调：普通通知、奖惩下发、玩家通知展示、锁定季度 `422` 拦截均已验证 | 已完成 |

| I3-01 | 明确单组演练库 / 单组演练模式的需求边界与文档口径 | P1 | 已完成 | `docs/requirements_spec.md`、`docs/testing_guide.md`、`docs/technical_selection.md` | 文档更新记录、环境边界说明 | 已明确“单组演练”仅用于开发 / 测试 / 联调 / 验收支持，不改变正式比赛规则；正式比赛库、多组演练库、单组演练库的定位已写入正式文档 | 已完成 |
| I3-02 | 落单组演练环境：独立配置、初始化数据与重置脚本 | P1 | 已完成 | `I3-01`、`M1-03`、`M1-04` | `configs/local-single.yaml`、`migrations/mysql/0002_seed_single_group.sql`、`scripts/init-single-rehearsal.ps1`、`scripts/reset-single-rehearsal.ps1`、`cmd/dbtool/main.go`、运行说明 | 能在独立环境中仅初始化 `1` 个管理员和 `1` 个玩家组，并支持重复执行 `0年 -> 最终年` 单链路推演，不污染正式比赛库；已实际完成单组库重置与数据核验（`groups=1`、`accounts=2`、`initialBaselineSubmitted=0`） | 已完成 |
| I3-03 | 明确“赛前初始化可配置小组数量”的正式需求与技术口径 | P0 | 已完成 | `docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/technical_selection.md` | 文档回写、边界说明 | 已明确仅支持赛前初始化可配置小组数量；比赛开始后不支持动态增减；运行期以 `sg_group` 实际数据为唯一依据 | 已完成 |
| I3-04 | 实现管理员端赛前配置页与比赛初始化接口 | P0 | 已完成 | `I3-03`、`M1-03`、`M2-11` | 管理员赛前配置页、`get-setup-status` / `initialize-game` 接口、批量生成服务 | 已完成后端 `get-setup-status / initialize-game` 接口、管理员动态默认路由与前端 `admin/setup` 页面；管理员可在比赛未初始化时配置 `groupCount` 并初始化比赛环境，系统以事务方式生成对应数量的 `sg_group`、`sg_account`、`sg_group_year_state` 数据；运行期以 `sg_group` 实际数据推断初始化状态与实际组数，并保留 `admin` 与 `group01 ~ groupNN` 命名规则 | 已完成 |
| I3-05 | 补管理员赛前初始化回归验证与运行说明 | P1 | 已完成 | `I3-04` | 测试文档、联调手册、验收清单 | 已补后端集成测试、管理员默认路由回归、本地独立运行配置 `configs/local-group-count-test.yaml`，并完成用户手工验收：管理员可在赛前配置页成功初始化 `6` 个小组，且管理员端可成功将最终年份配置为 `5`；本轮按用户确认作为完成标记 | 后续若需要更强边界回归，可继续补 `1 / 10` 组人工验收样例 |

| I4-01 | 明确线上订单与市场竞标正式需求口径 | P0 | 已完成 | 用户 2026-05-18/2026-05-19 确认、`道具-订单推算（服务企业）.xlsx` | `docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/minimal_state_machine.md`、`docs/calculation_rule_spec.md`、`docs/excel_field_mapping.md`、`docs/api_design.md`、`docs/database_design.md` | 已明确 `0年` 不需要订单；`1年` 起按 `市场 + 订单类型` 形成最多 `16` 个独立标段；系统按订单推算 Excel 公式链生成订单池，管理员每年开始前配置 `0~15` 的订单卡片数量，预览批次可重生成，确认后固化随机种子、公式版本、参数快照和订单明细；正式年份开标前玩家一次性提交 16 项市场投入，经营页只读带入；每标段一轮选单、每组最多一单；`1年` 按当前标段投入排序，`2年` 起市场龙头优先；非龙头同投入按上一年度该市场订单总额再随机；玩家可放弃，管理员可代跳过但不能代选；选单状态 `3` 秒自动轮询；交付绑定季度且销售收入需等于交付订单金额合计；未交付只记录状态，处罚暂不做 | 已完成文档沉淀；由于新确认改变了已实现链路，后续需执行 `I4-04R` 重构 |
| I4-02 | 实现管理员订单管理：Excel 上传、数量配置、标段释放顺序、订单池生成 | P0 | 已完成 | `I4-01`、`docs/api_design.md`、`docs/database_design.md` | `migrations/mysql/0005_order_admin.sql`、`internal/enum/order.go`、`internal/model/entity/order.go`、`internal/repository/order_repository.go`、`internal/service/order_excel_parser.go`、`internal/service/admin_order_service.go`、`internal/http/dto/admin_order_dto.go`、`internal/http/handler/admin_order_handler.go`、`internal/app/router.go`、`frontend/src/api/sandbox-game/admin-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/router/index.ts`、`frontend/src/components/sandbox-game/admin/AdminNav.vue`、`frontend/src/types/sandbox-game-admin.ts` | 历史实现口径：管理员可上传订单推算 Excel，按 `年份 + 市场 + 订单类型` 配置订单数量和释放顺序，生成固定订单池；后续已被 `I4-04R / I4-06 / I4-08` 逐步修正为“系统按公式生成、市场手动开启、多年订单数量控制台为唯一数量源”；本行仅保留为历史任务记录 | 已通过 `go test ./...` 与 `frontend/npm.cmd run build` 验证；后续以 `I4-08` 的多年订单数量控制台与市场预测链路作为最新开发口径 |
| I4-03 | 实现玩家年度订单页：市场投入、标段释放、选单顺序、选择/放弃订单 | P0 | 已完成 | `I4-02`、`docs/minimal_state_machine.md` | `internal/enum/order.go`、`internal/repository/order_repository.go`、`internal/service/player_order_service.go`、`internal/http/dto/player_order_dto.go`、`internal/http/handler/player_order_handler.go`、`internal/http/handler/admin_order_handler.go`、`internal/app/router.go`、`frontend/src/api/sandbox-game/player-order.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`frontend/src/api/sandbox-game/admin-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/components/sandbox-game/player/PageModeSwitch.vue`、`frontend/src/router/index.ts` | 玩家可在年度订单页提交市场投入、查看当前市场/标段顺序、按轮到本组时选择或放弃订单；已选订单立即锁定并在后续玩家侧置灰展示；页面按后端返回的 `pollingIntervalSeconds=3` 自动轮询；管理员订单页已补开放/关闭市场、释放下一个标段、查看市场选单状态、代跳过当前组控制能力；首版仍不做订单交付和经营页 Q1 前置拦截 | 已通过 `go test ./...` 与 `frontend/npm.cmd run build` 验证；下一步进入 `I4-04`，打通订单模块与经营页规则联动 |
| I4-04 | 打通订单模块与经营页规则联动 | P0 | 已完成 | `I4-03`、`docs/calculation_rule_spec.md` | `internal/service/order_operating_link_service.go`、`internal/repository/order_repository.go`、`internal/service/player_order_service.go`、`internal/service/player_operating_query_service.go`、`internal/service/player_operating_command_service.go`、`internal/service/player_report_query_service.go`、`internal/service/player_report_command_service.go`、`internal/service/admin_group_data_query_service.go`、`internal/http/dto/player_order_dto.go`、`internal/http/handler/player_order_handler.go`、`internal/http/handler/player_operating_handler.go`、`internal/app/router.go`、`frontend/src/api/sandbox-game/player-order.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/types/sandbox-game.ts`、`frontend/src/utils/sandbox-game-operating-preview.ts`、`frontend/src/components/sandbox-game/player/OperatingSheet.vue`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`docs/implementation_plan.md` | 正式年份经营页查询/保存/提交/财报计算统一读取订单模块汇总值；未完成订单前置时拦截 `Q1` 提交；已选订单金额汇总进入订单总额，市场投入汇总进入市场投入；订单页新增单订单交付与待交付订单批量交付，交付时校验当前季度销售收入等于本季度已交付订单金额合计；年末提交时将仍未交付订单标记 `UNFINISHED` 且不自动处罚；直接成本仍由玩家填写，账期只展示 | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build` 验证；下一步进入 `I4-05`，补订单模块测试、联调脚本与页面验收清单 |
| I4-04R | 重构订单生成与市场投入前置链路 | P0 | 已完成 | 用户 2026-05-19 新确认、`道具-订单推算（服务企业）.xlsx`、已完成的 `I4-02~I4-04` | `migrations/mysql/0006_order_generation_refactor.sql`、`internal/service/order_generation_engine.go`、`internal/service/admin_order_service.go`、`internal/service/player_order_service.go`、`internal/repository/order_repository.go`、`internal/model/entity/order.go`、`internal/http/dto/admin_order_dto.go`、`internal/http/dto/player_order_dto.go`、`internal/http/handler/admin_order_handler.go`、`internal/http/handler/player_order_handler.go`、`internal/app/router.go`、`frontend/src/api/sandbox-game/admin-order.ts`、`frontend/src/api/sandbox-game/player-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/types/sandbox-game-admin.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue` | 已将订单主链路从“上传 Excel / 解析固定订单池”调整为“系统生成订单池”：当时口径为管理员配置 16 个标段数量与释放顺序；后续 `I4-08` 已确认订单数量改由多年订单数量控制台统一维护，年度订单管理页只读带入当年数量并只保存市场开启与释放顺序；玩家端仍为开标前一次性提交 16 项市场投入；同一市场唯一市场龙头会在该市场 4 个产品标段共用 | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build` 验证；偏差说明：上传 Excel 功能保留为兼容接口，但不再是订单池主链路；首版订单生成公式按已提取的 Excel 控制台公式链落为 `ORDER_GEN_SERVICE_V1`，后续以 `I4-08` 接入多年控制台与市场预测 |
| I4-05 | 补订单模块测试、联调脚本与原型/页面验收 | P1 | 已完成 | `I4-04R` | `internal/service/order_generation_engine_test.go`、`internal/service/order_workflow_integration_test.go`、`tests/http/sandbox-game/AdminOrder.http`、`docs/order_module_acceptance_checklist.md` | 已补订单生成与排序单元测试、订单主链路集成测试、`.http` 联调脚本和管理员/玩家页面验收清单；覆盖 `0年` 无订单、订单生成预览/确认/随机种子固化、订单数量 `0~15` 校验、16 项市场投入必填、无投入标段跳过、投入并列按上年市场订单额再随机、市场龙头优先、破产龙头失效、标段释放顺序锁定、管理员不能跳序、每标段一轮选单、管理员代跳过、Q1 前置拦截、订单交付校验、未交付状态保留等场景 | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build` 验证；偏差说明：本轮以服务层自动化 + `.http` + 人工验收清单收口，不额外引入浏览器 E2E 或 WebSocket 实时验收 |
| I4-06 | 调整订单市场开启、订单池编号与管理员页操作链路 | P0 | 已完成 | 用户 2026-05-20 确认、`I4-04R`、`I4-05` | `migrations/mysql/0007_order_market_enable.sql`、`internal/enum/order.go`、`internal/model/entity/order.go`、`internal/repository/order_repository.go`、`internal/service/admin_order_service.go`、`internal/service/player_order_service.go`、`internal/service/order_generation_engine.go`、`internal/service/order_operating_link_service.go`、`internal/http/dto/admin_order_dto.go`、`internal/http/handler/admin_order_handler.go`、`internal/http/handler/player_order_handler.go`、`internal/app/router.go`、`frontend/src/api/sandbox-game/admin-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/types/sandbox-game-admin.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`frontend/src/components/sandbox-game/admin/AdminNav.vue` | 已完成市场开启配置表与接口：本地默认开启，区域/全国/全球默认关闭，订单池确认后锁定；未开启市场不生成订单、不进入选单，玩家仍需提交 16 项投入且未开启市场非 0 会被后端拦截；新增 `MARKET_DISABLED / NO_ORDER_CONFIG` 状态并纳入订单前置完成判断；订单池查询默认全量，支持按市场/类型筛选，页面主编号显示 `CARD-01` 等业务编号；管理员订单页已按“年份刷新、市场开启、标段配置、生成预览、确认订单池、生成顺序、释放标段、订单池查看”重排 | 已通过 `GOCACHE=<仓库绝对路径>/.go-build-cache go test ./...` 与 `frontend/npm.cmd run build` 验证；偏差说明：上传 Excel 兼容接口继续保留但不进入页面主链路 |
| I4-07 | 增加订单专项测试场景造数能力 | P0 | 已完成 | `I4-08`、用户 2026-05-21 确认 | `cmd/dbtool/main.go`、`cmd/dbtool/order_scenario.go`、`docs/order_module_acceptance_checklist.md`、`docs/implementation_plan.md` | 新增 `dbtool -action seed-order-scenario -confirm-reset` 测试造数命令：执行时会重建当前配置指向数据库并跑完整迁移，随后初始化 3 个小组、管理员与玩家账号、合法但非真实的 `0年 / 1年` 经营与平衡财报结果；同时造出 `1年` 各市场已选订单结果，用于形成 `2年` 市场龙头；最终停在 `2年已开放、1年订单历史已写入、2年订单配置未开始`，不预置 `2年` 多年订单数量控制台、市场开启状态、标段释放顺序或订单池，不提前提交市场投入、不关闭经营页与财报页校验、不改变正式比赛规则 | 已通过 `GOCACHE=.go-build-cache go test ./...` 验证；偏差说明：出于防误用考虑，造数命令必须显式追加 `--confirm-reset` 才会执行重置，且未自动运行真实数据库重置；2026-06-01 按用户测试诉求将停点前移，以便管理员端人工验证订单数量配置、市场开启、释放顺序、订单池生成与确认链路；下一步可按 `docs/order_module_acceptance_checklist.md` 从该停点做管理员端 / 玩家端人工验收 |
| I4-08 | 重构多年订单数量控制台与市场预测链路 | P0 | 已完成 | 用户 2026-05-21 确认、`道具-订单推算（服务企业）.xlsx`、`I4-06` | `migrations/mysql/0008_order_forecast_control.sql`、`internal/model/entity/order.go`、`internal/repository/order_repository.go`、`internal/service/admin_order_service.go`、`internal/http/dto/admin_order_dto.go`、`internal/http/handler/admin_order_handler.go`、`internal/http/handler/player_order_handler.go`、`internal/app/router.go`、`frontend/src/types/sandbox-game-admin.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/api/sandbox-game/admin-order.ts`、`frontend/src/api/sandbox-game/player-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`internal/service/order_generation_engine_test.go`、`internal/service/order_workflow_integration_test.go`、`internal/service/player_operating_command_service_test.go`、`tests/http/sandbox-game/AdminOrder.http`、`docs/implementation_plan.md` | 已落地新口径：新增多年订单数量控制台表与市场预测快照表，控制台固定覆盖 `1年~8年 × 四个市场 × 四类产品`；管理员端新增控制台维护区与三阶段预测展示/说明维护；年度订单管理页不再编辑订单数量，只读带入控制台当年数量并保存市场开启与释放顺序；订单生成从控制台读取当年数量，未开启市场仍不生成、不竞标；玩家端在市场投入前只读展示 `1~3年 / 4~5年 / 6~8年` 市场预测 | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build` 验证；偏差说明：上传 Excel 兼容接口继续保留但不作为控制台主链路，市场预测金额首版按已落地订单生成公式参数做确定性估算，后续若继续精确复刻 Excel 的全部随机展示口径需另立任务 |
| I4-09 | 修正订单页人工验收问题：静默轮询、Excel 竖状预测图、整数订单金额口径 | P0 | 已完成 | 用户 2026-05-22 测试反馈、`道具-订单推算（服务企业）.xlsx`、`I4-08` | `frontend/src/stores/player-order.ts`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`internal/service/order_generation_engine.go`、`internal/service/admin_order_service.go`、`internal/service/order_generation_engine_test.go`、`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/calculation_rule_spec.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/excel_field_mapping.md`、`docs/minimal_state_machine.md`、`docs/order_module_acceptance_checklist.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md` | 已完成三项修正：玩家端 `3` 秒轮询改为静默刷新，不触发整页加载态并保留滚动位置；玩家端与管理员端市场预测均改为贴近 Excel 的竖状柱状图，阶段下按市场分图，年份为横轴，四类订单为柱状系列；订单生成按公式目标值寻找最接近的整数金额，再反算两位小数单价，保证 `订单金额 = 数量 × 单价`，管理员端/玩家端金额显示为整数 | 待本轮验证完成后进入页面人工复测：重点检查订单页停留在订单池中等待轮询不跳顶、预测柱状图与 Excel 方向一致、重新生成订单池后金额不出现小数 |

### 9.7 I5：沙盘版本包与多版本字段模板

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| I5-01 | 明确沙盘版本包与贵宾服务版字段模板正式口径 | P0 | 已完成 | 用户 2026-05-26 确认、`guibing` 分支字段配置 | `docs/requirements_spec.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/excel_field_mapping.md`、`docs/implementation_plan.md` | 已明确字段模板、公式规则、流程规则拆分；首版先完整落地 `VIP_SERVICE_V1`，绑定贵宾经营页、贵宾财报页、贵宾订单页字段，公式和流程共用当前通用版本；管理员只在赛前初始化时选择内置版本包，初始化后锁定；不做在线字段编辑和在线公式编辑 | 已完成，后续生产版另立版本包 |
| I5-02 | 实现沙盘版本包底座与赛前版本选择 | P0 | 已完成 | `I5-01` | `internal/service/game_edition_registry.go`、`migrations/mysql/0009_game_edition.sql`、`sg_game_config` 版本字段、`get-current/list-game-editions/get-setup-status/initialize-game` 接口、管理员赛前初始化页版本选择 | 管理员可在未初始化时选择 `VIP_SERVICE_V1`；初始化后版本锁定并在配置接口只读展示；已有测试库默认补齐 `VIP_SERVICE_V1`，不影响现有数据 | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build` |
| I5-03 | 接入贵宾服务版经营页与财报页字段模板 | P0 | 已完成 | `I5-02`、`guibing` 分支 | `frontend/src/configs/sandbox-game-service-labels.ts`、`OperatingSheet.vue`、`ReportSheet.vue`、`ReportSidebar.vue`、`AdminBaselinePage.vue`、财报必填提示 | 当前主线经营页、财报页、初始基线页显示贵宾服务版字段；系统字段名、payload、公式计算不变；订单页继续使用贵宾订单字段 | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build` |
| I5-04 | 明确生产制造版版本包方案 | P0 | 已完成 | 用户 2026-05-26 新确认、`shengchan` 分支字段配置、已完成 `I5-01 ~ I5-03` | `docs/requirements_spec.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/excel_field_mapping.md`、`docs/implementation_plan.md` | 已明确新增 `PRODUCTION_V1`，经营页和财报页字段从 `shengchan` 分支提取，公式规则和流程规则继续共用当前通用版本；订单字段暂时继续绑定 `VIP_ORDER_TEMPLATE_V1`，待生产版订单字段来源确认后再新增生产版订单模板 | 已完成 |
| I5-05 | 实现 `PRODUCTION_V1` 内置版本包与版本列表返回 | P0 | 已完成 | `I5-04`、用户明确开始 | `internal/service/game_edition_registry.go`、`internal/service/admin_control_command_service_test.go`、`internal/service/admin_control_command_integration_test.go` | `list-game-editions` 同时返回 `VIP_SERVICE_V1` 与 `PRODUCTION_V1`；初始化接口允许选择 `PRODUCTION_V1` 并锁定版本字段；默认版本继续保持 `VIP_SERVICE_V1` | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build` |
| I5-06 | 接入生产制造版经营页与财报页字段模板 | P0 | 已完成 | `I5-05`、`shengchan` 分支字段配置 | `frontend/src/configs/sandbox-game-service-labels.ts`、`OperatingSheet.vue`、`ReportSheet.vue`、`ReportSidebar.vue`、`AdminBaselinePage.vue`、`PlayerOperatingPage.vue`、`PlayerReportPage.vue`、`player-report` store | 当前比赛为 `PRODUCTION_V1` 时，经营页、财报页、初始基线页显示生产制造字段；当前比赛为 `VIP_SERVICE_V1` 时继续显示贵宾服务字段；订单页仍使用贵宾订单字段 | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build` |
| I5-07 | 收口版本包前后端测试与页面验收 | P1 | 已完成 | `I5-05`、`I5-06` | 后端版本包测试、前端 build、手工验收步骤 | 验证赛前初始化页可选两个版本；初始化后版本锁定；两种版本下经营页/财报页/初始基线页显示正确；订单字段在生产版下仍按临时贵宾订单口径展示并在文档中标注 | 自动化验证已通过；页面人工验收待用户运行确认 |

### 9.8 I6：市场投入限制与全局手工数字整数化

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| I6-01 | 实现市场投入上限与未开启市场自动置零 | P0 | 已完成 | 用户 2026-05-27 确认、`I4-08`、`I4-09` | `migrations/mysql/0010_order_market_investment_limit.sql`、`cmd/dbtool/main.go`、`internal/model/entity/order.go`、`internal/repository/order_repository.go`、`internal/service/admin_order_service.go`、`internal/service/player_order_service.go`、`internal/http/dto/admin_order_dto.go`、`internal/http/handler/admin_order_handler.go`、`internal/http/handler/player_order_handler.go`、`internal/app/router.go`、`frontend/src/types/sandbox-game-admin.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`internal/service/order_generation_engine_test.go`、`internal/service/order_workflow_integration_test.go`、`internal/service/player_operating_command_service_test.go` | 未开启市场在玩家端禁用并自动按 `0` 提交；管理员可按 `年份 + 市场` 配置非负整数单市场投入上限，空值为无上限；已有任意小组提交本年投入后市场开启与上限锁定；玩家端展示“未开启 / 无上限 / 上限 XM”；提交时校验市场投入为非负整数且单市场 4 项合计不超过上限，错误提示精确到市场和字段 | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build`；本轮仅实现 I6-01，不扩大到 I6-02 全局手工数字整数校验 |
| I6-02 | 实现全局手工数字整数校验 | P0 | 已完成 | 用户 2026-05-27 确认、`I6-01` | `internal/service/manual_integer_validation.go`、`internal/service/manual_integer_validation_test.go`、`internal/service/player_operating_command_service.go`、`internal/service/player_report_command_service.go`、`internal/service/admin_control_command_service.go`、`internal/service/admin_notice_command_service.go`、`internal/http/handler/player_operating_handler.go`、`internal/http/handler/player_report_handler.go`、`internal/http/handler/admin_control_handler.go`、`internal/http/handler/admin_notice_handler.go`、`frontend/src/utils/manual-integer.ts`、`frontend/src/components/sandbox-game/player/OperatingSheet.vue`、`frontend/src/components/sandbox-game/player/ReportSheet.vue`、`frontend/src/components/sandbox-game/player/ReportSidebar.vue`、`frontend/src/stores/player-operating.ts`、`frontend/src/stores/player-report.ts`、`frontend/src/stores/admin-baseline.ts`、`frontend/src/stores/admin-notice.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`frontend/src/views/sandbox-game/player/report/PlayerReportPage.vue`、`frontend/src/views/sandbox-game/admin/baseline/AdminBaselinePage.vue`、`frontend/src/views/sandbox-game/admin/notice/AdminNoticePage.vue`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/views/sandbox-game/admin/control/AdminControlPage.vue`、`frontend/src/views/sandbox-game/admin/setup/AdminSetupPage.vue` | 所有用户手工填写的金额、数量、费用、成本、投入、上限、分数、最终年份、小组数量等数字均按整数口径校验；所得税税率下拉、系统计算结果、系统生成订单单价、历史/导入/公式结果展示不受限制；前端对小数输入高亮并阻止提交，后端在经营页、财报页、初始基线、奖惩金额等提交链路硬校验防绕过 | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build`；Go telemetry 仍因用户目录权限打印 token 告警但测试退出码为 0 |

### 9.9 I7：回退与修正、状态快照与单组恢复

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| I7-01 | 明确回退与快照正式规则口径 | P0 | 已完成 | 用户 2026-05-27 确认、`I1` 异常解锁、订单模块已落地口径 | `docs/requirements_spec.md`、`docs/minimal_state_machine.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/calculation_rule_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/implementation_plan.md` | 已明确管理员端统一 `回退与修正` 入口及单组阶段级快照恢复；2026-07-20 奖惩和破产子口径由 `I12-01` 更新为“普通退回保留奖惩、快照恢复奖惩状态、破产不可撤销”，其余订单、汇总和补提规则保持 | 已完成；奖惩与破产以 I12 最新规则为准 |
| I7-02 | 建立状态快照、回退日志与失效标记模型 | P0 | 已完成 | `I7-01` | `migrations/mysql/0011_rollback_snapshot.sql`、`cmd/dbtool/main.go`、`scripts/build-competition-package.ps1`、`internal/model/entity/state_snapshot.go`、`internal/model/entity/group_year_state.go`、`internal/model/entity/group_summary_snapshot.go`、`internal/model/entity/group_adjustment.go`、`internal/model/entity/admin_unlock_log.go`、`internal/model/entity/order.go`、`internal/enum/rollback.go`、`internal/repository/rollback_repository.go`、`internal/repository/rollback_stage.go`、相关仓储与测试兜底、`internal/service/rollback_errors.go` | 已新增状态快照主表、快照载荷表、回退日志表；补齐年份状态、汇总、奖惩、订单交付、异常解锁日志的回退/失效字段；`dbtool` 与正式打包脚本已纳入 `0011` 迁移；奖惩和订单交付新增默认生效与失效标记仓储方法；迁移和测试辅助均支持现有测试库补列 | 已通过 `GOCACHE=.go-build-cache go test ./...`；本轮只完成模型层，不实现自动快照生成、恢复服务或管理端页面 |
| I7-03 | 实现自动快照生成服务 | P0 | 已完成 | `I7-02` | `internal/service/rollback_snapshot_service.go`、自动触发点接入、`internal/model/entity/state_snapshot.go` | 经营阶段提交、财报提交、退回前、订单池确认、标段完成、开放下一年均能生成对应快照；快照载荷可用于审计和单组恢复 | 已通过 `go test ./...` 与前端构建 |
| I7-04 | 实现单组快照恢复服务 | P0 | 已完成 | `I7-02`、`I7-03` | `internal/service/admin_rollback_service.go`、`internal/http/handler/admin_rollback_handler.go`、`internal/http/dto/admin_rollback_dto.go`、`internal/app/router.go` | 管理员可选择单组快照恢复到目标阶段；系统生成安全快照，目标节点之后结果失效但保留；订单归属不释放；回退日志与待重提状态落库 | 已通过 `go test ./...` |
| I7-05 | 实现管理员端“回退与修正”页面 | P1 | 已完成 | `I7-04` | `frontend/src/views/sandbox-game/admin/AdminRollbackPage.vue`、`frontend/src/stores/admin-rollback.ts`、`frontend/src/api/sandbox-game/admin-rollback.ts`、管理员路由与导航 | 管理员可在统一页面执行退回重提、查询快照、创建手动快照、恢复单组快照，并看到风险提示与操作结果 | 已通过 `npm run build` |
| I7-06 | 接入年度控制阻断、玩家提示与汇总待重提展示 | P0 | 已完成 | `I7-04`、`I7-05` | `internal/service/admin_control_query_service.go`、`internal/service/admin_summary_query_service.go`、玩家经营/财报视图、汇总页 | 存在回退补提 / 待重提小组时管理员不能开放下一年；玩家端提示被回退到的年份阶段；汇总页对失效结果显示 `待重提` 且不计入排名 | 已通过 `go test ./...` 与 `npm run build` |
| I7-07 | 回归测试与人工验收清单 | P0 | 已完成 | `I7-02` ~ `I7-06` | `docs/rollback_module_acceptance_checklist.md`、服务层与前端构建验证 | 覆盖本年退回重提、跨年单组快照恢复、订单不释放、汇总待重提、开放下一年阻断；奖惩保留/快照恢复与破产不可撤销的新回归已由 `I12-06` 补齐 | 无 |

### 9.10 I8：赛前配置中心与业务显示字典

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| I8-01 | 明确赛前配置中心与业务显示字典正式口径 | P0 | 已完成 | 用户 2026-06-03 确认、`I5` 版本包、`I6` 整数输入、`I7` 回退模块 | `docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/excel_field_mapping.md`、`docs/implementation_plan.md` | 已明确版本包负责字段结构/公式/流程，业务显示字典只改显示名；赛前配置中心采用向导式流程，初始化后只读展示并支持解锁修改显示名称；字典方案可复用，当前比赛使用字典快照；字段名修改静默同步，不影响玩家草稿、滚动和输入状态；公式编辑、Excel 抽公式、流程规则编辑、新增/删除市场或订单类型均不纳入本轮 | 本轮只做文档沉淀，不进入代码实现 |
| I8-02 | 建立字典方案、当前比赛字典快照与修改日志模型 | P0 | 已完成 | `I8-01` | `migrations/mysql/0012_dictionary.sql`、`internal/model/entity/dictionary.go`、`internal/repository/dictionary_repository.go`、`internal/model/entity/game_config.go`、`internal/repository/game_config_repository.go`、`cmd/dbtool/main.go` | 已新增字典方案、方案明细、当前比赛字典快照、修改日志和 `sg_game_config.dictionary_revision`；支持当前比赛快照替换、方案明细替换、revision 更新与日志分页；未改变字段模板、公式或流程表结构 | 已通过 `go test ./...` |
| I8-03 | 实现业务显示字典后端接口 | P0 | 已完成 | `I8-02` | `internal/service/admin_dictionary_service.go`、`internal/service/dictionary_defaults.go`、`internal/http/handler/admin_dictionary_handler.go`、`internal/http/dto/admin_dictionary_dto.go`、`internal/app/router.go`、`internal/service/admin_control_command_service.go` | 已实现当前字典查询、方案列表、方案明细、保存/删除方案、更新当前比赛、应用方案、恢复默认、日志分页、revision 查询；初始化时写入当前比赛字典快照；跨版本方案应用会被拒绝 | 已通过 `go test ./...` |
| I8-04 | 改造赛前配置页面为向导式配置中心 | P0 | 已完成 | `I8-03`、现有 `admin/setup`、`admin/baseline` | `frontend/src/views/sandbox-game/admin/setup/AdminSetupPage.vue`、`frontend/src/api/sandbox-game/admin-control.ts`、`frontend/src/api/sandbox-game/admin-dictionary.ts`、`frontend/src/types/sandbox-game-admin.ts` | 赛前配置页已集中展示基础信息、版本包选择、字典方案编辑、初始基线入口和确认初始化；支持保存、更新、删除、选择字典方案；初始化时提交最终字典快照 | 已通过 `frontend/npm.cmd run build` |
| I8-05 | 实现初始化后当前比赛配置与字典维护页 | P1 | 已完成 | `I8-03`、`I8-04` | `frontend/src/views/sandbox-game/admin/setup/AdminSetupPage.vue`、`frontend/src/components/sandbox-game/admin/DictionaryEditor.vue`、`frontend/src/stores/dictionary.ts`、`docs/dictionary_module_acceptance_checklist.md` | 初始化后赛前配置页切为当前比赛配置视图，展示版本、公式、流程、字段模板、当前显示名称和修改日志；支持解锁修改当前比赛显示名称、应用方案、恢复默认、另存为方案；业务显示字典编辑器已改为分类展开，支持修改数量、空值数量、搜索、只看已修改、行级高亮和空值拦截；初始基线分类已从 6 个版本差异字段补齐为 20 个录入字段；最终年份仍在年度控制页维护 | 已通过 `frontend/npm.cmd run build` |
| I8-06 | 玩家端和管理员端页面接入当前比赛字典显示名 | P0 | 已完成 | `I8-03`、`I8-05` | `frontend/src/stores/dictionary.ts`、`frontend/src/configs/sandbox-game-service-labels.ts`、玩家经营/财报/订单页、管理员初始基线/订单/组数据页 | 玩家经营页、财报页、订单页、管理员初始基线页、订单管理页、组数据只读经营/财报视图已接入当前比赛字典显示名；市场和订单类型按字典显示；系统状态、按钮、菜单、错误提示和自由文本不做字典替换 | 本轮补齐 `frontend/src/views/sandbox-game/admin/group-data/AdminGroupDataPage.vue` 字典接入；已通过 `frontend/npm.cmd run build` |
| I8-07 | 实现字典 revision 静默同步 | P0 | 已完成 | `I8-03`、`I8-06` | `frontend/src/stores/dictionary.ts`、玩家经营/财报/订单页、管理员初始基线/订单/组数据页 | 前端字典 store 已按 `dictionaryRevision` 轻量检查并静默拉取新名称；页面只更新名称映射，不重新拉经营/财报/订单业务数据，不清空草稿、不改变输入状态、不滚动到顶部、不弹窗打断玩家 | 本轮补齐管理员组数据页静默同步；已通过 `frontend/npm.cmd run build` |
| I8-08 | 回归测试与人工验收清单 | P0 | 已完成 | `I8-02` ~ `I8-07` | `internal/service/admin_dictionary_service_integration_test.go`、`docs/dictionary_module_acceptance_checklist.md`、前端构建验证 | 已覆盖字典方案保存、读取、编辑、删除，初始化生成当前比赛快照，初始化后修改当前字典并自增 revision，应用方案、恢复默认、日志记录、跨版本拒绝；已补人工验收清单，覆盖玩家端静默同步和长名称展示 | 已通过 `go test ./...` 与 `frontend/npm.cmd run build` |

### 9.11 I9：订单模块可插拔与机场订单模板

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| I9-01 | 明确订单模块可插拔与机场订单模板正式口径 | P0 | 已完成 | 用户 2026-06-12 确认、`机场沙盘订单.xlsx`、`I5` 版本包、`I8` 业务显示字典 | `docs/order_template_airport_plan.md`、`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/calculation_rule_spec.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md` | 已明确订单模块下一阶段目标为模板可插拔；贵宾/生产继续绑定 `VIP_ORDER_TEMPLATE_V1`；新增 `AIRPORT_V1` 机场沙盘版和 `AIRPORT_ORDER_TEMPLATE_V1`；机场模板为国内/国际 × 窄体/宽体，每标段 `0~28`，国内默认开启、国际默认关闭；机场经营/财报待接入，首轮只支持订单模块独立联调 | 本轮只做文档沉淀，不进入代码实现 |
| I9-02 | 建立后端订单模板注册表与模板解析服务 | P0 | 已完成 | `I9-01` | `internal/service/order_template_registry.go`、`internal/service/game_edition_registry.go`、订单模板返回结构 | 后端已按当前比赛 `orderTemplateVersion` 返回市场、订单类型、最大订单数、默认开启状态、公式版本、是否支持交付和订单字段定义；未知模板返回明确错误；贵宾模板保持原四市场 × 四订单类型 | 无 |
| I9-03 | 扩展订单数据模型支持模板与扩展字段 | P0 | 已完成 | `I9-02` | `migrations/mysql/0013_order_template_airport.sql`、`internal/model/entity/order.go`、`internal/repository/order_repository.go` | 订单数量控制台、市场配置、预测、批次、订单池、竞标状态、投入、选单顺序和已选订单均补 `order_template_version`；订单池新增 `order_payload_json` 保存机场扩展字段；旧贵宾数据按默认模板兼容读取 | 无 |
| I9-04 | 将现有贵宾订单抽为 `VIP_ORDER_TEMPLATE_V1` 并保持兼容 | P0 | 已完成 | `I9-02`、`I9-03` | `internal/service/admin_order_service.go`、`internal/service/player_order_service.go`、`internal/service/order_generation_engine.go` | 贵宾服务版和生产制造版继续使用 `VIP_ORDER_TEMPLATE_V1`；仍为四市场 × 四订单类型、每标段 `0~15`；开标、选单、龙头、交付和经营页联动规则保持不变 | 无 |
| I9-05 | 实现 `AIRPORT_V1` 版本包与机场订单模板 | P0 | 已完成 | `I9-02`、`I9-04` | `internal/service/game_edition_registry.go`、`internal/service/order_template_registry.go`、`internal/service/dictionary_defaults.go` | 赛前可选择 `机场沙盘版 V1`；机场版绑定 `AIRPORT_ORDER_TEMPLATE_V1`；市场为国内 / 国际，订单类型为窄体 / 宽体；机场市场和订单类型已纳入业务显示字典 | 无 |
| I9-06 | 实现 `AIRPORT_ORDER_FORMULA_V1` 订单生成公式 | P0 | 已完成 | `I9-05` | `internal/service/order_generation_engine.go`、`internal/service/order_generation_engine_test.go` | 已按 `机场沙盘订单.xlsx` 口径生成架次、吞吐量、客座率、跑道要求、航线区域、单价、整数总收入和账期；订单金额按底层总收入四舍五入为整数；单价保留 4 位小数；每标段最多 28 张 | 无 |
| I9-07 | 管理员订单页模板化改造 | P0 | 已完成 | `I9-04`、`I9-06` | `frontend/src/types/sandbox-game-admin.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue` | 管理员订单页市场、订单类型、数量上限、默认开启状态、预测图、订单池字段均由当前订单模板驱动；机场版显示 2 市场 × 2 类型和扩展字段；贵宾/生产不回归 | 无 |
| I9-08 | 玩家订单页模板化与机场交付禁用 | P0 | 已完成 | `I9-07` | `frontend/src/types/sandbox-game-order.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue` | 玩家订单页按模板展示市场投入、预测图、订单卡字段和选单状态；机场版显示 4 项投入和机场订单字段；机场订单交付按钮置灰且后端拒绝交付；轮询仍保持静默不跳顶 | 无 |
| I9-09 | 机场经营页 / 财报页待接入占位 | P0 | 已完成 | `I9-05` | `frontend/src/views/sandbox-game/player/operating/PlayerOperatingPage.vue`、`frontend/src/views/sandbox-game/player/report/PlayerReportPage.vue` | 选择 `AIRPORT_V1` 后，玩家进入经营页 / 财报页显示待接入占位，不展示贵宾或生产字段，不允许保存或提交；订单模块仍可独立联调 | 无 |
| I9-10 | 机场订单专项造数命令与验收清单 | P1 | 已完成 | `I9-06`、`I9-08`、`I9-09` | `cmd/dbtool/airport_order_scenario.go`、`docs/test_demo_commands.md`、`docs/airport_order_module_acceptance_checklist.md` | 已支持 `seed-airport-order-scenario -confirm-reset` 停在机场版 `1年订单配置前`，支持 `seed-airport-order-leader-scenario -confirm-reset` 停在 `2年订单配置前` 并带 `1年` 市场龙头历史；验收清单覆盖机场订单和贵宾/生产回归 | 无 |

### 9.12 I10：第 9 步订单数量留痕

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| I10-01 | 冻结订单数量留痕与回退口径 | P0 | 已完成 | 用户 2026-07-17 确认、现有经营页与回退状态机 | `docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/api_design.md`、`docs/excel_field_mapping.md`、`docs/calculation_rule_spec.md`、`docs/minimal_state_machine.md`、`game doc/Excel计算规则与跨表联动说明.md` | 已明确生产版/贵宾版四类型 × 四季度、非负整数、无订单填 `0`、按年独立、不参与公式、随阶段级回退修改并保留版本 | 无 |
| I10-02 | 扩展后端经营负载与校验 | P0 | 已完成 | `I10-01` | `internal/model/payload/operating_payload.go`、`internal/rules/operating/operating_validator.go`、`internal/service/manual_integer_validation.go`、相关单元测试 | 新字段可保存、回显和进入完整快照；当前季度四项必填且为非负整数；旧 JSON 兼容；不进入计算 | 针对性 Go 测试已通过 |
| I10-03 | 实现生产版与贵宾版经营页录入 | P0 | 已完成 | `I10-02` | `frontend/src/components/sandbox-game/player/OperatingSheet.vue`、`frontend/src/types/sandbox-game.ts`、`frontend/src/stores/player-operating.ts`、版本化标签与字典默认项 | 两版本均显示四行季度录入，锁定、草稿、提交和管理员只读查看正常；机场版不受影响 | 前端生产构建已通过 |
| I10-04 | 回退回归、原型同步与验收 | P0 | 已完成 | `I10-02`、`I10-03` | `internal/model/payload/operating_payload_test.go`、经营计算/回退补提测试、`demo 网页/excel-combined-demo.html`、`demo 网页/player-operating-refined-demo.html`、`docs/implementation_plan.md` | 回退到目标季度后可修改并补提；旧提交/安全快照保留；旧 JSON 不伪补 `0`；数据不跨年、不影响公式；针对性 Go 测试与前端构建通过 | 全量 `go test ./...` 仅既有 `TestOrderWorkflowCoversGenerationSequenceSelectionDeliveryAndUnfinished` 市场龙头用例失败；排除该无关用例后全量通过 |

### 9.13 I11：总监得分公式调整

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| I11-01 | 冻结销售总监与 CEO 新公式及历史兼容口径 | P0 | 已完成 | 用户 2026-07-20 确认、`1组贵宾.xlsx`、`1组（服务）.xlsx`、现有财报递归重建机制 | `docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/calculation_rule_spec.md`、`docs/api_design.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md` | 已明确销售总监为累计订单总额整体除以 `10`，递推时仅将当年订单额除以 `10`；CEO 为五项总监最终得分合计加当前所有者权益的 `1/2`；新公式从 `0年` 追溯生效，负权益扣分，不额外取整；旧提交和快照不物理改写 | 无 |
| I11-02 | 修改后端权威公式并兼容历史重建、回退重提 | P0 | 已完成 | `I11-01` | `internal/rules/report/report_calculator.go`、`internal/rules/report/report_calculator_test.go`、`internal/service/player_report_command_service_test.go` | `0年` 与跨年销售得分均只缩放一次；CEO 使用当年 `reportTotalEquity / 2`；历史查询、后续承接、管理员只读查看和回退重提均按新公式计算；无需数据库迁移 | 已通过规则层与历史递归重建回归；旧财报 JSON 保持不改写 |
| I11-03 | 同步前端实时预览并完成公式回归 | P0 | 已完成 | `I11-02` | `frontend/src/types/sandbox-game.ts`、后端单元/集成测试、前端构建 | 前端不对后端销售总监得分二次除以 `10`，CEO 预览加入实时权益的一半；覆盖负权益、奇数权益小数、跨年递推、历史重建与旧数据不改写；前端构建和相关 Go 测试通过 | 完整 `go test ./...` 仍仅既有订单市场龙头用例失败；排除该无关用例后全量通过 |

### 9.14 I12：奖罚任意时段、年末重算与不可撤销破产

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| I12-01 | 冻结奖罚阶段归属、作废、静默同步与破产口径 | P0 | 已完成 | 用户 2026-07-20 确认、现有通知奖惩与回退链路 | `docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/minimal_state_machine.md`、`docs/api_design.md`、`docs/database_design.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md` | 已明确系统自动归属当前季度或 `YEAR_END`、财报草稿阶段可下发、正式提交后需先退回、折现费无年末字段、税前重算、下发/作废影响预览、3 秒静默局部同步、普通退回保留奖罚、快照恢复奖罚状态、破产快照不算正式财报且破产不可撤销 | 无 |
| I12-02 | 扩展奖罚事件存储、阶段解析与作废能力 | P0 | 已完成 | `I12-01` | `migrations/mysql/0015_adjustment_lifecycle.sql`、奖罚 Entity/Repository、`sg_group_adjustment_revision`、阶段解析与动作日志 | 已支持 `YEAR_END`、有效/已作废状态、完整作废审计和按组年单调 revision；阶段仅按服务端当前状态解析，忽略旧调用方传入阶段 | 无 |
| I12-03 | 实现影响预览、下发/作废权威重算与破产快照 | P0 | 已完成 | `I12-02`、经营/财报规则层、状态机 | `adjustment_impact_service.go`、`adjustment_bankruptcy_service.go`、preview/send/void 接口与破产快照载荷 | 下发与作废均可预览且正式执行再次权威重算；税后现金小于 `0` 时同事务标记不可撤销破产并生成 `ADJUSTMENT_BANKRUPTCY` 非正式只读快照 | 无 |
| I12-04 | 扩展经营/财报计算与玩家只读展示 | P0 | 已完成 | `I12-03` | operating/report calculator、玩家 DTO、`OperatingSheet.vue`、通知明细 | `YEAR_END` 只进入年度总计、年末现金和财报，不改变 Q1~Q4 季末现金；折现费年末为 `--`；奖罚汇总和有效/已作废/快照失效状态只读展示 | 无 |
| I12-05 | 实现管理员交互与玩家静默局部同步 | P0 | 已完成 | `I12-03`、`I12-04` | `AdminNoticePage.vue`、管理员 API/store、`PlayerAdjustmentSyncService`、玩家经营/财报局部同步 | 管理员支持下发/作废预览和破产警告；玩家页面可见时每 3 秒检查，隐藏暂停、恢复立即检查；只合并奖罚和派生结果，不调用整页刷新 | 无 |
| I12-06 | 补奖罚、回退、破产与用户体验专项回归 | P0 | 已完成 | `I12-02 ~ I12-05` | Go 专项测试、`tests/http/sandbox-game/AdminNotice.http`、`docs/adjustment_module_acceptance_checklist.md`、前端构建与页面验收 | 自动化覆盖阶段归属、正式提交锁定、revision、作废审计、普通退回保留、快照恢复、下发/作废破产、只读快照、破产冻结；浏览器验证预览弹窗、年末列、折现费 `--`、3 秒轮询不改变焦点/滚动/路由，控制台无错误 | 无；正式下发破产奖罚未在共享演练库点击确认，破产事务由集成测试覆盖 |

### 9.15 I13：订单模块多轮选单升级

> 最新口径以 `docs/order_multi_round_upgrade_design.md` 为主；本迭代覆盖贵宾服务版、生产制造版四轮资格和全部模板固定订单数量上限移除。机场版暂时保持单轮，造数命令最后单独升级。

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| I13-01 | 冻结多轮资格、龙头、订单池、回退与页面口径 | P0 | 已完成 | 用户 2026-07-21 多轮讨论确认、现有订单与回退模块 | `docs/order_multi_round_upgrade_design.md`、`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/minimal_state_machine.md`、`docs/database_design.md`、`docs/order_module_acceptance_checklist.md`、`docs/airport_order_module_acceptance_checklist.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md` | 已明确贵宾/生产固定 `3/6/9` 最多四轮、机场单轮、每市场每年一个龙头且投入 `0` 无资格、基础顺序只生成一次、固定订单池、后续轮手动开启、订单池选空结束、剩余订单过期、破产自动失去后续资格、回退不影响投入与竞标事实、全部模板取消固定数量上限 | 无；本任务只写文档，不修改代码或造数命令 |
| I13-02 | 扩展多轮数据模型并移除固定订单数量上限 | P0 | 已完成 | `I13-01` | `migrations/mysql/0016_order_multi_round.sql`、订单枚举/Entity/Repository、模板注册表、生成器、测试建表兜底、`cmd/dbtool/main.go`、打包脚本 | 已增加 `ROUND_READY/current_round_no/completion_reason/round_no/selection_order_id/UNSELECTED_EXPIRED`，调整轮次唯一索引；模板元数据以 `0` 表示无固定数量上限；生成器已验证可生成 `54` 张并拒绝负数 | 无；正式环境迁移执行留到统一部署流程 |
| I13-03 | 重构基础顺序、市场龙头与全部轮次一次生成 | P0 | 已完成 | `I13-02` | `internal/service/player_order_service.go`、`internal/service/order_template_registry.go`、订单仓储与现有订单工作流集成测试 | 已按当前标段投入生成资格轮次；同市场每年只计算一次龙头；零投入龙头被排除；同一基础顺序过滤生成全部有效轮次；机场模板仍返回单轮 | 无；更完整的 `0/1/3/6/9/10` 参数化测试并入 `I13-07` |
| I13-04 | 重构轮次推进、选单、放弃、跳过、选空和破产处理 | P0 | 已完成 | `I13-03` | `player_order_service.go`、`order_bankruptcy_service.go`、Repository、DTO/Handler/Router | 已实现第一轮自动开始、后续轮手动开启、当前轮选择/放弃/管理员跳过、空轮自动跳过、订单池选空结束、剩余池过期；经营提交/财报提交/管理员奖惩触发破产时只失效未完成轮次，当前组破产自动推进且已选订单保留 | 已补正常结束过期、放弃后续轮、破产完成原因和重复请求幂等测试；定序、释放和轮次状态查询使用事务锁 |
| I13-05 | 改造管理员端四轮顺序与控制区 | P0 | 已完成 | `I13-04` | 管理员订单类型、API、store、`AdminOrderPage.vue` | 已接入“开启下一轮”；按轮次纵向展示全部顺序；显示当前/下一轮、理论机会、剩余订单和订单池不足警告；无限数量输入不再被 `max=0` 限制 | 已通过浏览器验收管理员四轮纵向区、当前/下一轮状态和统一开启按钮；无资格轮次显示“本轮无参与小组” |
| I13-06 | 改造玩家端本组轮次状态与静默同步 | P0 | 已完成 | `I13-04` | 玩家订单类型、store、`PlayerOrderPage.vue` | 已展示本组当前轮状态和跨轮已选订单；放弃为“放弃本轮”并增加无原因确认；后端按服务端当前轮返回本组状态；沿用 3 秒静默轮询和滚动位置恢复 | 已通过浏览器验收玩家当前轮、订单卡和放弃入口；轮询后 URL、滚动位置、焦点保持不变，控制台无错误/警告 |
| I13-07 | 补多轮订单自动化、接口与人工验收 | P0 | 已完成 | `I13-02 ~ I13-06` | Go 单元/集成测试、`tests/http/sandbox-game/AdminOrder.http`、`docs/order_module_acceptance_checklist.md`、浏览器人工验收 | 已覆盖资格 `0/1/3/6/9/10`、基础顺序轮次过滤、机场单轮、放弃后续轮、管理员跳过原因、正常结束 `UNSELECTED_EXPIRED`、破产失效和已选历史保留、重复生成/重复开启幂等；补齐下一轮/放弃 HTTP 样例；`go test ./...`、前端构建和 `git diff --check` 通过；浏览器验收通过 | 无；I13-08 造数命令仍按用户确认延后 |
| I13-08 | 升级订单专项造数命令 | P1 | 部分完成 | `I13-07` | `cmd/dbtool/order_scenario.go`、`docs/test_demo_commands.md`、`docs/order_module_acceptance_checklist.md` | 已切到手动联调口径：保留 `1年` 龙头历史、预置 `3` 个正常组和 `1` 个破产组、停在 `2年已开放/订单数量与市场投入待手动配置`；后续再按需要补快速复现型停点 | 继续完善手动场景验证与必要的回归说明 |
| I13-09 | 收口玩家端订单明细可见性 | P0 | 已完成 | 用户 2026-07-23 确认“页面不改、展示分开，只控制展示时机” | `internal/service/player_order_service.go`、`internal/service/order_generation_engine_test.go`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`docs/order_multi_round_upgrade_design.md`、`docs/requirements_spec.md`、`docs/api_design.md`、`docs/order_module_acceptance_checklist.md`、`docs/implementation_plan.md` | 后端玩家订单视图已新增 `ordersVisible`，仅在 `SELECTING / ROUND_READY / COMPLETED` 返回订单池和本组已选订单明细；市场投入、等待提交、顺序已生成但标段未释放等状态返回空明细；前端保持原页面结构，订单卡片、本组已选订单和交付面板分开展示并统一受 `ordersVisible` 控制；已补服务层状态策略和裁剪测试 | 无；后续由用户按验收清单复测市场投入阶段和标段释放后的页面展示 |

### 12.1 本轮新增记录

| 日期 | 记录 |
|---|---|
| 2026-07-23 | 任务状态：✅ 已完成；落点：`internal/repository/admin_action_log_repository.go`、`internal/service/admin_control_query_service_test.go`、`docs/implementation_plan.md`；偏差说明：用户发现后端启动后 `get-config` 周期性出现 `sg_admin_action_log` 的 `record not found` 日志。经定位，该接口查询“最新管理员动作”时允许日志表为空，上层已将 `gorm.ErrRecordNotFound` 视为无最新动作，但仓储层使用 GORM `First` 导致空结果被打印为报错噪音。本轮改为 `Limit(1).Find(&slice)` 查询，空结果仍按仓储契约返回 `gorm.ErrRecordNotFound`，但不再触发 GORM record-not-found 日志；新增测试覆盖空表查询不输出该噪音。验证：`go test ./internal/service -run TestAdminActionLogFindLatestEmptyDoesNotEmitRecordNotFound -count=1`、`go test ./internal/repository` 通过。下一步：用户重启后端后刷新管理员端，确认无最新管理员动作时不再持续打印该条 `record not found`。 |
| 2026-07-23 | 任务状态：✅ 已完成；落点：`docs/requirements_spec.md`、`docs/api_design.md`、`internal/repository/summary_snapshot_repository.go`、`internal/service/admin_summary_query_integration_test.go`；偏差说明：用户在 `seed-order-scenario` 停点发现第四组当前已破产，但管理员汇总页仍显示正常。经确认，汇总页最右侧 `经营状态` 应展示小组当前经营状态，而不是某一年汇总快照中的历史状态。本轮将年度汇总与最终排名查询中的状态来源改为 `sg_group.business_status`，收入、利润、权益仍来自正式汇总快照；新增测试覆盖“快照 NORMAL、当前小组 BANKRUPT 时接口返回 BANKRUPT”。验证：`go test ./internal/service`、`go test ./internal/repository`、`frontend/npm.cmd run build`、`git diff --check` 均通过。下一步：用户刷新管理员汇总页或重新运行订单造数后确认第四组显示红色 `已破产`。 |
| 2026-07-23 | 任务状态：✅ 已完成；落点：`internal/service/player_order_service.go`、`internal/service/order_generation_engine_test.go`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`docs/order_multi_round_upgrade_design.md`、`docs/implementation_plan.md`；偏差说明：按用户确认不调整玩家订单页结构，只给后端玩家视图新增 `ordersVisible` 并按服务端状态裁剪订单明细；`SEQUENCE_READY` 等未释放状态即使存在订单池和已选订单历史，也不返回订单金额、数量、单价、账期、订单编号和交付状态；`SELECTING / ROUND_READY / COMPLETED` 才返回明细；前端仍分别展示订单卡片、本组已选订单和交付面板。验证：`go test ./internal/service`、`frontend/npm.cmd run build`、`git diff --check` 均通过。下一步：用户按订单验收清单复测“市场投入阶段不可见订单”和“标段释放后可见订单/交付”。 |
| 2026-07-23 | 任务状态：🟡 部分完成；落点：`docs/order_multi_round_upgrade_design.md`、`docs/requirements_spec.md`、`docs/api_design.md`、`docs/order_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：本轮只按用户确认补文档，不修改代码和页面结构。已明确玩家填写市场投入时仍可看市场预测，但不能看到订单卡片、本组已选订单或交付面板；页面上三个区域继续分开展示，不新增统一大模块；后端玩家年度订单视图后续需返回 `ordersVisible`，不可见时不返回具体订单明细，不能只靠前端隐藏。下一步：如用户确认进入开发，则实施 `I13-09` 后端裁剪、前端展示条件和测试。 |
| 2026-07-22 | 任务状态：🟡 部分完成；落点：`cmd/dbtool/order_scenario.go`、`docs/test_demo_commands.md`、`docs/order_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：按用户确认将 `seed-order-scenario` 切为手动联调版，预置 `3` 个正常组和 `1` 个破产组，保留 `1年` 订单历史用于显示 `2年` 市场龙头，停在 `2年已开放、1年龙头历史已写入、2年订单数量/市场投入待手动配置`；验证：`go test ./cmd/dbtool` 通过；下一步：用户按手动流程在前端配置订单数量、市场开启和市场投入，继续观察多轮竞标。 |
| 2026-07-22 | 任务状态：✅ 已完成；落点：`internal/repository/order_repository.go`、`internal/service/player_order_service.go`、`internal/service/order_generation_engine_test.go`、`internal/service/order_workflow_integration_test.go`、`tests/http/sandbox-game/AdminOrder.http`、`docs/order_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：本轮按 `docs/order_multi_round_upgrade_design.md` 完成 `I13-04 ~ I13-07` 收口，新增生成顺序/释放标段的 `FOR UPDATE` 锁，防止并发请求重复重排或重复释放；补充 `0/1/3/6/9/10` 资格边界、机场单轮、基础顺序轮次过滤、放弃后续轮、正常结束剩余订单 `UNSELECTED_EXPIRED`、破产只失效未完成轮次且保留已选历史、重复生成/开启/放弃幂等测试；HTTP 样例增加“开启下一轮”和“放弃当前轮”；浏览器验证管理员四轮纵向展示、统一开启按钮、玩家当前轮/放弃入口及 3 秒静默轮询不跳顶；验证：`go test ./... -count=1`、`frontend/npm.cmd run build`、`git diff --check` 全部通过。`I13-08` 订单专项造数命令继续延后，未修改。下一步：用户可按 `docs/order_module_acceptance_checklist.md` 在独立演练库做完整业务场景复测。 |
| 2026-07-21 | 任务状态：🟡 部分完成；落点：`migrations/mysql/0016_order_multi_round.sql`、订单枚举/Entity/Repository、`internal/service/player_order_service.go`、`internal/service/order_bankruptcy_service.go`、订单 DTO/Handler/Router、管理员与玩家订单 types/API/store/page、订单测试基础设施；偏差说明：`I13-02 / I13-03` 已完成，`I13-04 ~ I13-07` 已完成主体编码但未宣告验收完成，`I13-08` 按用户要求仍未开始。当前实现包括无固定订单数量上限、固定 `3/6/9` 四轮资格、零投入龙头无资格、同市场单龙头、一次基础顺序过滤、多轮共用订单池、后续轮“开启下一轮”、选空结束、剩余订单过期、破产未完成轮次失效、管理员四轮纵向展示及玩家本组轮次展示。验证：`go test ./internal/service ./internal/repository ./internal/http/handler ./internal/app` 通过，`frontend/npm.cmd run build` 通过。下一步：先补 `I13-07` 专项测试（资格 `0/1/3/6/9/10`、放弃/跳过、过期、破产、机场单轮、并发幂等），再运行 `go test ./...`、`git diff --check` 和浏览器人工验收；验收通过后再开始 `I13-08` 订单专项造数命令。当前工作区尚未提交，续开发时保留无关的既有前端改动，不要回退。 |
| 2026-07-21 | 任务状态：✅ 已完成；落点：`docs/order_multi_round_upgrade_design.md`、`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/minimal_state_machine.md`、`docs/database_design.md`、`docs/order_module_acceptance_checklist.md`、`docs/airport_order_module_acceptance_checklist.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md`；偏差说明：本轮只完成 `I13-01` 文档收口，不修改代码、迁移或造数命令。已确认贵宾/生产按当前标段投入采用固定 `3/6/9` 最多四轮，机场暂时单轮；每市场每年一个龙头，龙头投入 `0` 无资格；每标段只生成一次基础顺序，后续轮按资格过滤；第一轮随标段释放自动开始，后续轮由管理员点击“开启下一轮”；固定订单池提前选空时立即结束，正常结束后的剩余订单只保留审计；全部模板取消 `15/28` 固定数量上限；回退/快照不修改订单投入和竞标事实；造数命令延后到功能开发完成后。下一步：从 `I13-02` 多轮数据模型和数量上限改造开始开发。 |
| 2026-07-21 | 任务状态：✅ 已完成；落点：`migrations/mysql/0015_adjustment_lifecycle.sql`、奖罚生命周期 Entity/Repository、`internal/service/adjustment_*`、`internal/service/player_adjustment_sync_service.go`、管理员通知与玩家经营/财报前端、`tests/http/sandbox-game/AdminNotice.http`、`docs/adjustment_module_acceptance_checklist.md`、`docs/testing_guide.md`；偏差说明：`I12-02 ~ I12-06` 已按冻结口径完成。实现包括服务端自动阶段归属、下发/作废前影响预览与执行时再次权威重算、税前奖罚重算、`YEAR_END` 年末列、作废审计、组年 revision、3 秒玩家局部同步、不可恢复的 `ADJUSTMENT_BANKRUPTCY` 破产快照、普通退回保留奖罚、快照恢复奖罚状态且破产不可撤销。验证：`go test ./... -skip '^TestOrderWorkflowCoversGenerationSequenceSelectionDeliveryAndUnfinished$' -count=1` 通过，`frontend/npm.cmd run build` 通过，`git diff --check` 通过；未排除时仍仅既有订单市场龙头用例失败，与 I12 无关。浏览器已验证管理员页无季度选择、Q1 自动归属、完整影响预览和破产警告，玩家经营页年末列/折现费 `--`，3 秒轮询后焦点、值、滚动、路由与年份保持不变，玩家财报通知区正常且控制台无错误。为避免污染共享演练库，未点击会永久破产的正式确认按钮；该事务由下发/作废破产集成测试覆盖。下一步：用户可按 `docs/adjustment_module_acceptance_checklist.md` 在独立演练库继续完整业务场景复测。 |
| 2026-07-20 | 任务状态：✅ 已完成；落点：`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/minimal_state_machine.md`、`docs/api_design.md`、`docs/database_design.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md`；偏差说明：本轮只完成 `I12-01` 文档收口，不修改代码、数据库迁移或原始 Excel。已确认管理员可在 Q1~Q4、年末经营和财报草稿阶段下发奖罚，系统自动归属当前阶段，财报阶段统一归属 `YEAR_END`；折现费用不增加年末字段；奖罚按税前收入/支出口径重算；下发和作废前均显示影响预览；玩家端采用 3 秒 revision 静默局部同步且不得刷新跳顶或覆盖未保存输入；普通退回保留奖罚，恢复快照按快照恢复有效状态；奖罚导致现金小于 0 时生成不计入正式财报、汇总和排名的破产快照，破产不可撤销且破产后奖罚冻结。下一步：从 `I12-02` 开始扩展存储、阶段解析和作废能力。 |
| 2026-07-20 | 任务状态：✅ 已完成；落点：`internal/rules/report/report_calculator.go`、`internal/rules/report/report_calculator_test.go`、`internal/service/player_report_command_service_test.go`、`frontend/src/types/sandbox-game.ts`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I11-02 / I11-03`。后端最佳销售总监得分改为“上一年已缩放得分 + 当年订单总额 / 10”，最佳 CEO 得分在五项总监最终得分合计基础上增加当前 `reportTotalEquity / 2`；前端财报实时预览直接复用后端销售总监得分，仅给 CEO 增加实时预览权益的一半。新增规则层回归覆盖 `0年`、跨年不重复缩放、负权益扣分与小数保留；新增服务层历史重建测试验证旧 `0年` 财报 JSON 不物理改写、`1年` 查询仍按 `100/10 + 140/10 = 24` 返回新销售得分。验证：针对性规则层与历史重建测试通过，`go test ./internal/rules/report ./internal/service -skip '^TestOrderWorkflowCoversGenerationSequenceSelectionDeliveryAndUnfinished$' -count=1` 通过，`go test ./... -skip '^TestOrderWorkflowCoversGenerationSequenceSelectionDeliveryAndUnfinished$' -count=1` 通过，`frontend/npm.cmd run build` 通过；未排除时仍仅既有订单市场龙头集成用例失败，失败内容与本次财报公式无关。下一步：页面人工验收 `0年 / 1年` 销售得分、负权益 CEO 扣分和财报手工项变化时的 CEO 实时预览；确认后可提交本次公式改动。 |
| 2026-07-20 | 任务状态：✅ 已完成；落点：`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/calculation_rule_spec.md`、`docs/api_design.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md`；偏差说明：本轮按用户要求只完成总监得分公式文档冻结，不修改代码、不改原始 Excel、不新增数据库字段。经只读核验 `1组贵宾.xlsx / 1组（服务）.xlsx`，原 `H27` 为从 `0年` 起累计订单总额，原 `H29` 为总监得分区合计；新规则明确销售总监得分为累计订单总额整体除以 `10`，跨年递推必须使用“上一年已缩放得分 + 当年订单总额 / 10”，CEO 在五项总监最终得分合计基础上增加当前 `K19 / reportTotalEquity` 的 `1/2`；从 `0年` 起追溯，负权益扣分，结果不额外取整，旧提交与快照不物理改写。下一步：进入 `I11-02 / I11-03` 修改后端权威公式、前端预览并补充跨年/历史/回退测试。 |
| 2026-07-17 | 任务状态：✅ 已完成；落点：`internal/model/payload/operating_payload.go`、`internal/rules/operating/operating_validator.go`、`internal/service/manual_integer_validation.go`、`internal/http/handler/player_operating_handler.go`、`frontend/src/components/sandbox-game/player/OperatingSheet.vue`、`frontend/src/types/sandbox-game.ts`、`frontend/src/stores/player-operating.ts`、版本化标签/字典、相关测试、规则文档与 HTML 原型；偏差说明：按用户确认将生产版“下新供应链订单”和贵宾版“下新服务订单”从静态提醒改为四类型 × 四季度的非负整数手工留痕，当前季度四项必填、无订单填 `0`、按年独立、不进入公式或正式订单链；现有阶段级退回重提/恢复快照可修改并保留旧提交、安全快照和补提版本，无需数据库迁移。验证：针对性 payload、经营校验、整数校验、计算不受影响和回退补提版本测试通过；`go test ./... -skip '^TestOrderWorkflowCoversGenerationSequenceSelectionDeliveryAndUnfinished$'` 全部通过；`frontend/npm.cmd run build` 通过；完整 `go test ./...` 中既有订单市场龙头集成用例仍失败且单独复跑一致，未修改该无关订单逻辑。下一步：在页面分别用 `PRODUCTION_V1 / VIP_SERVICE_V1` 做 Q1 填 `0`、负数拦截、季度锁定和退回 Q1 后修改的人工验收。 |
| 2026-06-12 | 任务状态：✅ 已完成；落点：`internal/service/order_template_registry.go`、`internal/service/game_edition_registry.go`、`internal/model/entity/order.go`、`migrations/mysql/0013_order_template_airport.sql`、`internal/service/order_generation_engine.go`、`internal/service/admin_order_service.go`、`internal/service/player_order_service.go`、`internal/service/dictionary_defaults.go`、`cmd/dbtool/airport_order_scenario.go`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/types/sandbox-game-admin.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/views/sandbox-game/player/operating/PlayerOperatingPage.vue`、`frontend/src/views/sandbox-game/player/report/PlayerReportPage.vue`、`docs/test_demo_commands.md`、`docs/airport_order_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I9-02 ~ I9-10` 首轮实现，订单模块已从固定贵宾订单升级为按当前比赛订单模板驱动；新增机场版 `AIRPORT_V1` 和 `AIRPORT_ORDER_TEMPLATE_V1`，机场订单支持数量配置、市场开关、订单池生成、市场投入、开标、选单、放弃、锁定和市场龙头联调；机场经营页和财报页仍按已确认边界只显示待接入占位，不做完整经营闭环。验证：已通过 `$env:GOCACHE='E:\project\sand box game\.go-build-cache'; go test ./...` 与 `frontend/npm.cmd run build`。下一步：按 `docs/airport_order_module_acceptance_checklist.md` 做页面人工验收，重点验证机场版 4 项投入、28 张上限、扩展字段展示、交付禁用和贵宾/生产不回归。 |
| 2026-06-12 | 任务状态：✅ 已完成；落点：`docs/order_template_airport_plan.md`、`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/calculation_rule_spec.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md`；偏差说明：本轮仅按用户要求完成“订单模块可插拔 + 机场订单模板”多轮讨论后的文档收口，不进入代码实现。已确认新增 `AIRPORT_V1` 机场沙盘版，订单模板为国内/国际 × 窄体/宽体，每标段 `0~28`，国内默认开启、国际默认关闭；机场订单按 Excel 公式生成架次、客座率、跑道、航线、单价和总收入，订单金额按 Excel 显示口径四舍五入为整数；机场经营页和财报页 Excel 尚未提供前只做占位，不套用贵宾/生产字段，订单模块可独立联调。下一步：用户确认后从 `I9-02` 后端订单模板注册表开始实现。 |
| 2026-06-04 | 任务状态：✅ 已完成；落点：`frontend/src/stores/admin-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`docs/implementation_plan.md`；偏差说明：用户人工测试发现玩家完成选单后，管理员端竞标控制区必须手动刷新才会更新标段状态并释放下一个标段。本轮仅优化管理员端状态刷新体验：竞标控制 tab 激活时每 3 秒静默拉取当前市场选单状态，切换 tab / 离开页面自动停止，不改变选单、跳过、释放和订单归属规则。下一步：用户在玩家选单后观察管理员端标段状态是否自动从“选单中”更新为“已完成 / 下一个可释放”。 |
| 2026-06-04 | 任务状态：✅ 已完成；落点：`internal/service/admin_order_service.go`、`internal/service/player_order_service.go`、`frontend/src/types/sandbox-game-admin.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`docs/implementation_plan.md`；偏差说明：本轮按用户确认调整管理员订单管理使用链路：数量控制台保存后自动刷新未确认年份的市场预测与预览订单池，已确认年份在数量控制台整列置灰并由后端拒绝变更；订单池 tab 去掉手动生成预览入口，改为展示预览/正式状态并保留确认订单池；生成选单顺序移动到竞标控制区，并在玩家市场投入未提交完整时提示具体未提交小组。下一步：运行后端测试与前端构建后，由用户按“保存数量 -> 查看预览订单池 -> 确认订单池 -> 玩家提交投入 -> 生成选单顺序”做页面人工复测。 |
| 2026-06-04 | 任务状态：✅ 已完成；落点：`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`docs/implementation_plan.md`；偏差说明：本轮按用户确认优化管理员订单管理页体验，将原先平铺区块改为顶部流程 tabs：数量控制、市场设置、标段顺序、订单池、竞标控制；保留顶部年度状态总览，未改变订单数量控制台、市场开启、标段释放、订单池生成确认、竞标控制、订单池查看等业务规则、接口或数据结构。下一步：用户在页面人工确认流程顺序与 tab 文案是否符合现场管理员使用习惯。 |
| 2026-06-04 | 任务状态：✅ 已完成；落点：`internal/service/dictionary_defaults.go`、`internal/service/admin_dictionary_service.go`、`internal/service/admin_dictionary_service_integration_test.go`、`frontend/src/configs/sandbox-game-service-labels.ts`、`frontend/src/views/sandbox-game/admin/baseline/AdminBaselinePage.vue`、`docs/implementation_plan.md`；偏差说明：本轮按用户确认补齐业务显示字典的初始基线范围，将 `BASELINE` 分类从原先仅覆盖 6 个版本差异名称扩展为 20 个初始基线录入字段，并对已有当前比赛字典快照 / 旧字典方案读取时自动补齐新增默认项；不改变初始基线 payload、字段编码、计算公式和提交流程。下一步：用户在赛前配置页和初始基线页人工确认 20 个字段均可修改并同步显示。 |
| 2026-06-04 | 任务状态：✅ 已完成；落点：`frontend/src/components/sandbox-game/admin/DictionaryEditor.vue`、`frontend/src/views/sandbox-game/admin/setup/AdminSetupPage.vue`、`docs/implementation_plan.md`；偏差说明：本轮按用户人工反馈只优化业务显示字典编辑体验，不改变字段编码、公式、payload、接口语义或玩家端静默同步规则；字典编辑器由平铺输入框改为分类展开，增加已修改 / 空值汇总、搜索、只看已修改、分类与行级高亮、撤销与恢复默认操作，并在保存当前比赛名称前展示按分类分组的变更预览；下一步：用户在页面上继续人工查看效果，必要时再微调列宽、颜色和文案。验证：已通过 `frontend/npm.cmd run build`。 |
| 2026-06-04 | 任务状态：✅ 已完成；落点：`migrations/mysql/0012_dictionary.sql`、`internal/service/admin_dictionary_service.go`、`internal/service/dictionary_defaults.go`、`internal/http/handler/admin_dictionary_handler.go`、`internal/http/dto/admin_dictionary_dto.go`、`internal/model/entity/dictionary.go`、`internal/repository/dictionary_repository.go`、`frontend/src/stores/dictionary.ts`、`frontend/src/api/sandbox-game/admin-dictionary.ts`、`frontend/src/views/sandbox-game/admin/setup/AdminSetupPage.vue`、`frontend/src/views/sandbox-game/admin/group-data/AdminGroupDataPage.vue`、`docs/dictionary_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：本轮核对发现 `I8-02 ~ I8-08` 主体代码已在仓库中落地但计划表仍停留在未开始；本轮补齐管理员组数据页只读经营/财报视图的当前比赛字典显示名与 revision 静默同步接入，并将 I8 状态按实际完成情况回写。验证：已通过 `$env:GOCACHE=(Resolve-Path '.go-build-cache').Path; go test ./...` 与 `frontend/npm.cmd run build`。下一步：按 `docs/dictionary_module_acceptance_checklist.md` 做页面人工联调验收。 |
| 2026-06-03 | 任务状态：✅ 已完成；落点：`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/excel_field_mapping.md`、`docs/implementation_plan.md`；偏差说明：本轮仅按用户要求把 `I8 赛前配置中心与业务显示字典` 写回文档，不进入代码实现。已确认：版本包负责字段结构、公式版本和流程规则；业务显示字典只改显示名，不改字段含义、公式、数据或流程；赛前配置集中承载小组数量、最终年份、版本包、字典方案、初始基线入口和确认初始化；初始化后可解锁修改当前比赛显示名称、应用同版本包字典方案、恢复默认并留日志；玩家端和管理员端静默同步名称变化，不影响草稿、输入状态或页面滚动。明确不做：管理员编辑公式、上传 Excel 自动抽公式、编辑流程规则、新增/删除市场或订单类型、调整字段结构、公式参数编辑、多场比赛历史归档。下一步：用户确认后可从 `I8-02` 建表与后端字典模型开始实现。 |
| 2026-03-31 | 正式确认异常解锁采用“方案 2”：管理员必须选择 `OPERATING / REPORT`；当目标为 `OPERATING` 时，还需选择 `Q1 / Q2 / Q3 / Q4 / YEAR_END`。 |
| 2026-03-31 | 正式确认经营页回退后的规则：后续经营阶段与财报结果失效，但原值不清空，统一保留为 `失效草稿`，并在重新提交后重新生效。 |
| 2026-03-31 | `I1-01` 已落地：后端异常解锁接口已支持 `OPERATING / REPORT` 与 `targetStageCode`，管理端组数据页已接入目标类型 / 阶段选择，并完成 `go test ./...` 与 `npm run build` 验证。 |
| 2026-03-31 | `I1-02` / `I1-03` 已落地：后端按当前状态推断“有效结果 / 失效草稿”，上一年承接仅读取有效财报；前端已在玩家经营页、财报页和管理端组数据页展示失效草稿状态，并完成 `go test ./...` 与 `npm run build` 验证。 |
| 2026-03-31 | 将“通知与奖惩”拆成独立迭代任务 `I2-01 ~ I2-08`，覆盖设计、存储、接口、计算联动、管理端页面、玩家端展示与测试收口，作为下一阶段正式版开发主线。 |
| 2026-03-31 | `I2-01 ~ I2-08` 已在当前仓库落地：后端新增 `admin-notice` 接口域、`sg_notice` / `sg_group_adjustment` 迁移与 Repository，前端完成管理端“通知与奖惩页”和玩家端通知区/奖惩只读展示，并补齐自动化测试与 `AdminNotice.http` 联调脚本。 |
| 2026-04-01 | 正式确认“单组演练库 / 单组演练模式”仅作为开发、测试、联调、验收与公式回归支持能力，不改变正式比赛的年度推进规则；相关边界已回写 `requirements / consensus / testing / technical / runbook / implementation` 文档。 |
| 2026-04-01 | 新增 `docs/competition_launch_runbook.md`，正式收口首版比赛主机部署、正式比赛库、统一访问地址、赛前检查、现场运行与备份恢复建议；明确正式比赛推荐采用单主机、单入口、非开发态部署方式。 |
| 2026-04-01 | `I3-02` 已落地：新增 `configs/local-single.yaml`、`migrations/mysql/0002_seed_single_group.sql`、`scripts/init-single-rehearsal.ps1`、`scripts/reset-single-rehearsal.ps1` 与 `cmd/dbtool/main.go`，并已完成单组演练库重置验证，当前可从“管理员提交初始基线”开始重复执行完整主链推演。 |
| 2026-04-01 | 正式确认需要支持“管理员端赛前配置小组数量并初始化比赛”能力，并收口边界：仅赛前初始化可配，比赛开始后不支持动态增减；运行期统一以 `sg_group` 实际数据为唯一依据；后续开发任务已拆为 `I3-03 ~ I3-05`。 |
| 2026-04-02 | 已补“小组数量赛前可配置”文档口径：`requirements / technical / api / implementation` 已明确管理员默认入口由初始化状态决定，建议通过 `get-setup-status / initialize-game` 落地赛前配置流程，运行期直接以 `sg_group` 实际数据推断 `initialized / groupCount`，初始化流程需单事务创建 `sg_group / sg_account / sg_group_year_state` 并禁止重复初始化。 |
| 2026-04-02 | 已完成经营页底部 6 个指标口径收口：`市场回报比 / 研发投入强度 / 劳动生产率` 改由经营规则层统一计算，其中 `劳动生产率` 按 Excel `O59` 实现；`净产收益率 / 总产收益 / 净利润率 / 毛率润率` 改为读取当前年财报结果或财报预览；若当前年还没有财报草稿，服务端也会按空手工项即时生成财报预览，确保经营页继续按 Excel 公式链显示；并补齐固定公式样例测试与服务层注入回归测试，当前已通过 `go test ./...` 与 `npm run build`。 |
| 2026-04-02 | 已修正玩家经营页的前端 `previewCalculation` 合并问题：正常模式下，后端已返回的 `净产收益率 / 总产收益 / 净利润率 / 毛率润率` 不再被前端本地 preview 结果覆盖丢失；当前已通过 `npm run build` 验证。 |
| 2026-04-02 | 已补正式年份跨年联动回归：当前规则实现仅区分 `0年` 与“正式年”两条路径，不按 `1~8年` 单独分叉；现已新增 `1年` 承接 `0年财报` 的服务层测试，验证经营页底部收益率指标在正式年份仍按上一年财报基线 + 当年经营数据公式链生成，并通过 `go test ./internal/service/...` 与 `go test ./internal/rules/...`。 |
| 2026-04-02 | `I3-04` 已落地：后端新增管理员赛前初始化状态查询与比赛初始化接口，动态根据 `sg_group` 实际数据返回管理员默认入口；前端新增 `admin/setup` 页面、导航切换与路由守卫，未初始化时仅允许进入赛前配置页；初始化成功后会一次性生成玩家组、玩家账号与 `0年 + 正式年` 主状态，并已通过 `go test ./...` 与 `frontend/npm run build` 验证。 |
| 2026-04-02 | 已补正式版上线包与文档收口：新增 `competition.yaml`、正式比赛专用 `admin + sg_game_config` 初始化 SQL、`init/reset/start/build-competition` 脚本与 `Nginx` 样例配置；同步修正文档，明确正式比赛库不能直接使用默认 10 组种子脚本，而应由管理员首登后在 `赛前配置页` 初始化比赛。 |
| 2026-04-02 | 已新增 `docs/competition_deploy_checklist.md` 作为现场一键部署清单，按“赛前准备 -> 改配置 -> 打包 -> 初始化正式库 -> 启动后端 -> 配置 Nginx -> 联通验证 -> 故障排查”收口正式比赛当天执行步骤；同时已把该文档加入 `docs/README.md` 索引，并补入正式版上线包构建脚本，确保打包后可随包交付给现场技术支持。 |
| 2026-04-02 | 已完成新版 Excel 经营页 / 财报页收口：经营页 `productionLineAdjustment` 区块已改成“建成生产线转固 / 生产线残值 / 待折资产”的两级结构，并补上 `短期贷款右移一格 / 产品进度更新 / 下新供应链订单 / 应收账款向左移动一格` 提醒；财报页新增“总监得分”区，新增 `企业认证得分 / 最佳生产人力总监得分 / 关账速度得分` 三个绿色手工项，以及 `最佳市场经营总监得分 / 最佳科技创新总监得分 / 最佳销售总监得分 / 最佳 CFO 得分 / 最佳 CEO 得分` 五个黄色计算项；其中 `最佳市场经营总监得分` 与 `最佳 CFO 得分` 已按“黄色基础分 + 绿色补充分”实现，`最佳 CEO 得分` 按各行最终得分汇总，且市场总监分跨年只累计黄色基础分，不滚入上一年绿色加分；后端规则层已支持 `0年 CFO = 0` 与“黄色基础分 + 绿色补充分”口径，并对旧报表 JSON 缺失新增得分字段的情况做了“查询/提交时尽力回填、失败不阻断主链”的兼容兜底；当前已通过 `go test ./...`（使用工作区内 `GOCACHE`）与 `frontend/npm run build` 验证。 |
| 2026-04-02 | 已修正财报“最佳 CFO（财务总监）得分”跨年口径：按正式 Excel 逐年复核后确认，`0年~8年财报` 的黄色公式始终为 `0.5 * (当年 K19 - 0年财报 K19)`，并不是“只取上一年权益差”；后端现已改为“上一年 CFO 黄色基础分 + 本年权益增量一半”的等价实现，并将玩家端查询、提交流程、管理端组数据查看在读取上年财报时统一改成“优先按当前规则重建总监得分后再覆盖旧值”，避免历史已提交错值继续传到后续年份；同时新增 `2年` CFO 累计回归测试，当前已通过 `go test ./...`。 |
| 2026-04-02 | 已确认“季末现金核对”属于系统新增核对能力，不要求与 Excel `G41 / I41 / K41 / M41` 原始季度公式逐格一致；后续公式审计时应将其视为有意差异，而不是规则实现错误。 |
| 2026-04-03 | `guibing` 服务版分支已开始按 `game doc/1组（服务）.xlsx` 收口展示层：玩家经营页、玩家财报页、财报右侧栏、玩家财报必填校验提示、管理员初始基线页统一切换为服务业口径；其中 `B19` 文案按用户确认保留为“服务进度更新”，不按制造业版“税前利润”显示。 |
| 2026-04-03 | 服务版分支已同步收口正式上线包默认交付物：`competition.yaml` 默认数据库名改为 `sandbox_game_service_competition`，正式版打包产物名改为 `sandbox-game-service-competition`，并同步更新 `Nginx` 样例路径、正式上线 Runbook 与现场部署清单，避免与制造业版正式包混用。 |
| 2026-04-03 | 已继续补充服务版正式部署文档口径：`competition_launch_runbook / competition_deploy_checklist / docs/README` 已显式写明“服务版正式环境必须使用独立数据库、独立配置、独立上线包目录”，并补充制造业版与服务版的建议库名 / 包名对照，避免现场误把两套游戏共用同一正式库。 |
| 2026-04-05 | 已新增本地手工验收配置 `configs/local-group-count-test.yaml`：该配置专用于“小组数量赛前初始化”功能联调，数据库默认指向 `sandbox_game_group_count_test`，保持 `admin + sg_game_config` 初始化口径，不预置玩家组，便于管理员亲手进入 `赛前配置页` 验证 `1 / 6 / 10` 组初始化链路。 |
| 2026-04-05 | 用户已完成“小组数量赛前初始化 + 最终年份配置”人工验收：在 `sandbox_game_group_count_test` 环境中，管理员可成功配置并初始化 `6` 个小组，且可成功将最终年份配置为 `5`；本轮按用户确认将该功能记为完成。 |
| 2026-04-05 | 已补服务版初始基线缺失字段：新增 `baselineWorkInConstruction / baselineWorkInConstruction`，管理员初始基线页补上“在建贵宾厅”；同时在 `0年` 经营/财报查询与提交链路中新增“仅当经营页未显式填写时，按初始基线预填在建贵宾厅”的默认值逻辑，并补上“缺失时回填、显式 `0` 保留”自动化测试，避免正式开赛前非零基线被静默丢失。 |
| 2026-04-05 | 已继续复核“在建贵宾厅”后续年份口径：对照 `1组 最终版.xlsx / 1组（服务）.xlsx` 确认，财报 `G5` 仅引用“当年经营 `O52`”，不属于“财报 -> 下一年经营”的自动承接项；因此本轮新增字段只影响“初始基线 -> 0年”，不应擅自扩展为正式年份自动滚动，并已补 `formal year does not carry forward previous workInConstruction` 回归测试与规则文档说明，防止后续误改。 |
| 2026-04-05 | 已同步校准 `shengchan` 制造业分支的管理员初始基线页口径：当前工作区原先混入了服务版标签引用，现已明确该分支应显示“在建生产线 / 厂房 / 生产线残值 / 在制品 / 成品 / 材料”，并保留同一套后端字段与 `0年` 默认值修复逻辑，确保制造业版与 `1组 最终版.xlsx` 一致。 |
| 2026-04-05 | 已新增正式版管理员启动脚本模板目录 `scripts/competition-launchers`：按“制造业版 / 服务版”分别提供启动、重启、赛前重置 `.bat`，并补 `停止全部服务.bat` 与《管理员操作手册-正式版》；脚本默认约定一台管理员电脑、本机 `SandboxGameMySQL` 服务、两套独立正式包目录 `E:\deploy\sandbox-game-shengchan / sandbox-game-guibing`，以及单端口访问口径（制造业 `18080`、服务版 `28080`），用于后续交付非技术管理员现场一键操作。 |
| 2026-04-05 | 已将正式版部署口径从“`Nginx + Go`”收口为“Go 直接提供前端静态页面 + API”：新增 `frontend.distDir` 配置项与 `internal/app/frontend_static.go`，服务端现在会直接返回 `frontend/dist` 中的页面与资源；正式版打包脚本已不再复制 `nginx` 目录，管理员启动脚本同步改成单端口模式（制造业 `18080`、服务版 `28080`），并已回写 `competition_launch_runbook / competition_deploy_checklist / docs/README`。 |
| 2026-04-05 | 已新增两套正式版配置模板：`configs/competition-shengchan.yaml` 与 `configs/competition-guibing.yaml`；分别固定制造业版 `18080 + sandbox_game_shengchan_prod`、服务版 `28080 + sandbox_game_guibing_prod` 的单端口数据库口径，后续正式打包时只需替换数据库密码与 `tokenSecret` 即可。 |
| 2026-04-05 | 已将管理员手册拆分为两份可直接交付文件：`scripts/competition-launchers/安装部署手册-正式版.txt` 负责指导管理员自行安装 MySQL、创建数据库、摆放正式包与首次试运行；`scripts/competition-launchers/操作手册-正式版.txt` 负责比赛当天日常启动、重启、停服与风险提醒；旧 `管理员操作手册-正式版.txt` 保留为跳转说明，避免管理员拿到旧文件名时走错手册。 |
| 2026-04-13 | 已收口正式版管理员启动脚本的 MySQL 服务名兼容性：`启动/赛前重置` 脚本不再写死 `SandboxGameMySQL`，改为优先自动识别 `SandboxGameMySQL / MySQL / MySQL80 / MySQL57 / MariaDB` 等常见服务名；若现场机器使用自定义服务名，再通过脚本顶部 `MYSQL_SERVICE` 手工指定，并已同步更新交付目录说明与安装部署手册。 |
| 2026-04-13 | 已补正式版管理员“首次初始化数据库”双击入口：新增 `首次初始化制造业版数据库.bat / 首次初始化服务版数据库.bat`，管理员首次部署时不再需要手工执行 PowerShell 命令；安装部署手册、启动器目录说明与交付目录已同步改成“创建数据库 -> 双击首次初始化 -> 双击启动”的鼠标操作路径。 |
| 2026-04-16 | 已按试用反馈继续收口 `shengchan` 制造业版经营页提醒样式：`请将短期贷款右移一格 / 产品进度更新 / 下新供应链订单 / 请将应收账款向左移动一格` 等黄色提醒行统一改为更醒目的高亮提醒条，仅调整前端字号、底色、左侧强调边与字重，不改任何业务口径与表格结构。 |
| 2026-04-16 | 已新增正式版“前端覆盖更新包”脚本模板：`更新制造业版前端.bat / 更新服务版前端.bat` 支持从任意临时解压目录读取同包内 `frontend\\dist`，自动创建 `D:\\deploy\\_backup\\...` 备份、覆盖正式版前端目录并尝试调用现场 `重启服务` 脚本；同时已同步补充启动器 `README` 与管理员《操作手册-正式版》中的前端更新说明，避免纯展示层改动时误重装正式包或误动数据库。 |
| 2026-05-18 | 任务状态：✅ 已完成；落点：`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/minimal_state_machine.md`、`docs/calculation_rule_spec.md`、`docs/excel_field_mapping.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/implementation_plan.md`；偏差说明：本轮只做订单模块规则、接口与数据模型文档收口，不进入代码实现；线上订单规则由“全年一次抢订单”调整为“按本地/区域/全国/全球四个市场分别竞标，每市场一轮选单”，并明确 `0年` 不需要订单；下一步：按 `I4-02 ~ I4-05` 依次做管理员订单管理、玩家年度订单页、经营页联动和测试验收。 |
| 2026-05-19 | 任务状态：✅ 已完成；落点：`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/minimal_state_machine.md`、`docs/calculation_rule_spec.md`、`docs/excel_field_mapping.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/implementation_plan.md`；偏差说明：用户进一步确认“每个市场每个产品要单独开标”，因此本轮将订单流程从“每市场一轮”修正为“每个 `市场 + 订单类型` 独立标段，每年最多 `16` 个标段”，并新增“管理员配置标段释放顺序”；排序规则和跳过规则仍按此前已确认口径，不引入解锁状态动态判定；下一步：开发时先做管理员端标段数量与释放顺序配置，再做玩家端按当前释放标段选单。 |
| 2026-05-19 | 任务状态：✅ 已完成；落点：`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/minimal_state_machine.md`、`docs/calculation_rule_spec.md`、`docs/excel_field_mapping.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/implementation_plan.md`；偏差说明：本轮继续补齐订单规则细节：不做保底/多轮分配，由管理员配置冗余订单并由系统提示风险；每个标段只走一轮，当前标段结束后由管理员手动释放下一标段；市场龙头按上一年度该市场订单总额确定且每市场唯一，龙头 `0` 投入仍可优先；玩家可放弃，管理员只可代跳过不可代选；订单状态采用 `3` 秒自动轮询；直接成本仍由玩家填写，账期只展示，交付需绑定当前季度并校验销售收入等于交付订单金额合计；未交付暂不自动处罚；下一步：进入 `I4-02` 管理员订单管理实现前，先按本轮文档做接口/表结构任务拆分。 |
| 2026-05-19 | 任务状态：✅ 已完成；落点：`migrations/mysql/0005_order_admin.sql`、`internal/enum/order.go`、`internal/model/entity/order.go`、`internal/repository/order_repository.go`、`internal/service/order_excel_parser.go`、`internal/service/admin_order_service.go`、`internal/http/dto/admin_order_dto.go`、`internal/http/handler/admin_order_handler.go`、`internal/app/router.go`、`frontend/src/api/sandbox-game/admin-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/router/index.ts`、`frontend/src/components/sandbox-game/admin/AdminNav.vue`、`frontend/src/types/sandbox-game-admin.ts`、`docs/implementation_plan.md`；偏差说明：本轮只完成 `I4-02` 管理员订单管理，不提前实现玩家抢单页、管理员释放下一标段、生成选单顺序、市场投入提交和订单交付联动；为保证“开标前可调整”，配置更新时会清理未开始的旧订单池与旧标段状态，已开放/已选单年份仍禁止覆盖；验证：已通过 `go test ./...` 与 `frontend/npm.cmd run build`；下一步：进入 `I4-03`，实现玩家年度订单页、标段释放、选单顺序与选择/放弃订单。 |
| 2026-05-19 | 任务状态：✅ 已完成；落点：`internal/enum/order.go`、`internal/repository/order_repository.go`、`internal/service/player_order_service.go`、`internal/http/dto/player_order_dto.go`、`internal/http/handler/player_order_handler.go`、`internal/http/handler/admin_order_handler.go`、`internal/app/router.go`、`frontend/src/api/sandbox-game/player-order.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`frontend/src/api/sandbox-game/admin-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/components/sandbox-game/player/PageModeSwitch.vue`、`frontend/src/views/sandbox-game/player/operating/PlayerOperatingPage.vue`、`frontend/src/views/sandbox-game/player/report/PlayerReportPage.vue`、`frontend/src/router/index.ts`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I4-03` 玩家年度订单页与管理员标段控制闭环，并补足后端 `player-order` 与管理员开标控制接口；仍按任务边界不做订单交付接口、不做经营页 Q1 订单前置拦截、不做 WebSocket 实时推送；验证：已通过 `go test ./...` 与 `frontend/npm.cmd run build`；下一步：进入 `I4-04`，将已选订单金额、市场投入汇总、订单交付状态与经营页规则联动。 |
| 2026-05-19 | 任务状态：✅ 已完成；落点：`internal/service/order_operating_link_service.go`、`internal/repository/order_repository.go`、`internal/service/player_order_service.go`、`internal/service/player_operating_query_service.go`、`internal/service/player_operating_command_service.go`、`internal/service/player_report_query_service.go`、`internal/service/player_report_command_service.go`、`internal/service/admin_group_data_query_service.go`、`internal/http/dto/player_order_dto.go`、`internal/http/handler/player_order_handler.go`、`internal/http/handler/player_operating_handler.go`、`internal/app/router.go`、`frontend/src/api/sandbox-game/player-order.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/types/sandbox-game.ts`、`frontend/src/utils/sandbox-game-operating-preview.ts`、`frontend/src/components/sandbox-game/player/OperatingSheet.vue`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I4-04` 订单模块与经营页规则联动，不提前展开 `I4-05` 的完整场景测试清单和 `.http` 联调脚本；经营页正式年份订单总额/市场投入由后端订单模块覆盖并在前端只读展示，旧草稿仍保留兼容读取；年度订单页除单订单交付外，补了待交付订单批量勾选交付，以匹配“同季度可交付多个完整订单”的已确认规则；验证：首次 `go test ./...` 因 Windows 用户目录 Go 缓存权限失败，改用仓库内 `.go-build-cache` 后已通过 `go test ./...`，并通过 `frontend/npm.cmd run build`；下一步：进入 `I4-05` 补订单模块测试、联调脚本与页面验收清单。 |
| 2026-05-19 | 任务状态：✅ 已完成；落点：`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/minimal_state_machine.md`、`docs/calculation_rule_spec.md`、`docs/excel_field_mapping.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/implementation_plan.md`；偏差说明：用户重新确认 `道具-订单推算（服务企业）.xlsx` 是订单生成器而不是固定订单池上传来源，因此本轮把订单主链路改为“系统按 Excel 公式链生成订单”：每年开始前生成当年订单；管理员只配置 `0~15` 的订单卡片数量；预览批次可反复重生成；确认后固化随机种子、公式版本、参数快照和订单明细；上传 Excel 不再作为订单池主链路，仅可作为后续导入控制台参数/模板的辅助能力；同时市场投入从经营页前移为正式年份开标前一次性提交 `16` 项，经营页只读带入，订单竞标结束后才允许进入或提交 `Q1`。下一步：新增 `I4-04R`，优先重构现有 I4-02~I4-04 代码以对齐新口径。 |
| 2026-05-20 | 任务状态：✅ 已完成；落点：`migrations/mysql/0006_order_generation_refactor.sql`、`internal/service/order_generation_engine.go`、`internal/service/admin_order_service.go`、`internal/service/player_order_service.go`、`internal/repository/order_repository.go`、`internal/model/entity/order.go`、`internal/http/dto/admin_order_dto.go`、`internal/http/dto/player_order_dto.go`、`internal/http/handler/admin_order_handler.go`、`internal/http/handler/player_order_handler.go`、`internal/app/router.go`、`frontend/src/api/sandbox-game/admin-order.ts`、`frontend/src/api/sandbox-game/player-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/types/sandbox-game-admin.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I4-04R` 代码重构，上传 Excel 保留为兼容接口但不作为订单池主链路；订单生成公式按 `ORDER_GEN_SERVICE_V1` 落地，管理员端生成预览/确认订单池/生成选单顺序，玩家端一次性提交 16 项市场投入；同一市场唯一市场龙头在该市场 4 个产品标段共用；验证：已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build`；下一步：进入 `I4-05` 补完整订单模块测试、`.http` 联调脚本与页面验收清单。 |
| 2026-05-20 | 任务状态：✅ 已完成；落点：`internal/service/order_generation_engine_test.go`、`internal/service/order_workflow_integration_test.go`、`internal/service/player_operating_command_service_test.go`、`tests/http/sandbox-game/AdminOrder.http`、`docs/order_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I4-05` 订单模块测试验收收口，以服务层自动化测试覆盖核心规则和事务链路，以 `.http` 脚本与人工验收清单覆盖页面联调路径，不额外引入浏览器 E2E 或 WebSocket 验收；为支持新演练库首次跑测试，补充订单相关测试建表兜底；验证：已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build`；下一步：可进入订单模块人工页面验收或按用户继续安排后续开发任务。 |
| 2026-05-20 | 任务状态：🟡 部分完成；落点：`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/minimal_state_machine.md`、`docs/calculation_rule_spec.md`、`docs/excel_field_mapping.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/order_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：用户人工查看管理员订单页后确认需要调整订单链路体验与规则细节：市场开启改为管理员手动控制，本地默认开启、区域/全国/全球默认关闭；未开启市场不生成订单、不抢单，但玩家必须手动填 `0`，填非 `0` 拦截；订单池默认查看全部订单，订单主编号使用 `CARD-01` 等业务编号而非数据库 ID；订单数量为 `0` 显示“未配置订单”，不再误显示“已跳过”；管理员页按真实操作顺序重排，保存配置放回标段配置区；下一步：待用户确认本轮文档修改后，进入 `I4-06` 代码实现。 |
| 2026-05-20 | 任务状态：✅ 已完成；落点：`migrations/mysql/0007_order_market_enable.sql`、`internal/enum/order.go`、`internal/model/entity/order.go`、`internal/repository/order_repository.go`、`internal/service/admin_order_service.go`、`internal/service/player_order_service.go`、`internal/service/order_generation_engine.go`、`internal/service/order_operating_link_service.go`、`internal/http/dto/admin_order_dto.go`、`internal/http/handler/admin_order_handler.go`、`internal/http/handler/player_order_handler.go`、`internal/app/router.go`、`frontend/src/api/sandbox-game/admin-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/types/sandbox-game-admin.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`frontend/src/components/sandbox-game/admin/AdminNav.vue`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I4-06` 代码实现，上传 Excel 兼容接口仍保留但未放入管理员订单页主操作链路；市场开启状态采用独立年度配置并在订单池确认后锁定，未开启市场与未配置订单分别显示 `市场未开启 / 未配置订单`；订单池全量查看默认生效，页面主编号改为 `CARD-01` 等业务编号；验证：已通过 `GOCACHE=<仓库绝对路径>/.go-build-cache go test ./...` 与 `frontend/npm.cmd run build`；下一步：进入页面人工验收，重点检查管理员订单页操作顺序、未开启市场投入 0 拦截、订单池全量查看和业务编号展示。 |
| 2026-05-21 | 任务状态：⏳ 未开始；落点：`docs/implementation_plan.md`；偏差说明：用户确认完整按真实游戏顺序测试会被三组真实经营数据缺失、财报平衡校验和 `0年` 对市场龙头无参考价值拖住，因此新增 `I4-07` 作为订单专项测试场景造数任务；该能力只用于测试库造出合法但非真实的数据，不关闭校验、不改变正式比赛规则；下一步：先讨论并设计造数口径，再实现 `dbtool` 场景初始化命令。 |
| 2026-05-21 | 任务状态：⏳ 未开始；落点：`game doc/Excel计算规则与跨表联动说明.md`、`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/minimal_state_machine.md`、`docs/calculation_rule_spec.md`、`docs/excel_field_mapping.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/order_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：用户与游戏管理员确认现场按 Excel 控制台订单数配置真实订单数量，因此新增 `I4-08`：多年订单数量控制台成为订单数量主来源，市场预测是该控制台的三阶段趋势展示，年度订单池从当年已开启市场的控制台数量生成；市场开启不做数量自动推断，完全以管理员手动开启为准；下一步：先确认数据模型迁移和页面改造方案，再进入代码实现。 |
| 2026-05-21 | 任务状态：✅ 已完成；落点：`migrations/mysql/0008_order_forecast_control.sql`、`internal/model/entity/order.go`、`internal/repository/order_repository.go`、`internal/service/admin_order_service.go`、`internal/http/dto/admin_order_dto.go`、`internal/http/handler/admin_order_handler.go`、`internal/http/handler/player_order_handler.go`、`internal/app/router.go`、`frontend/src/types/sandbox-game-admin.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/api/sandbox-game/admin-order.ts`、`frontend/src/api/sandbox-game/player-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`internal/service/order_generation_engine_test.go`、`internal/service/order_workflow_integration_test.go`、`internal/service/player_operating_command_service_test.go`、`tests/http/sandbox-game/AdminOrder.http`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I4-08` 代码实现，多年订单数量控制台成为订单数量唯一主来源，年度标段页只读带入当年数量并只保存释放顺序；玩家端新增市场预测展示；上传 Excel 保留为兼容接口但不进入主链路；验证：已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build`；下一步：可进入管理员订单页和玩家订单页人工验收，或继续 `I4-07` 测试场景造数能力。 |
| 2026-05-21 | 任务状态：✅ 已完成；落点：`cmd/dbtool/main.go`、`cmd/dbtool/order_scenario.go`、`docs/order_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I4-07` 订单专项测试场景造数能力，新增 `seed-order-scenario` 动作并要求 `--confirm-reset` 显式确认；原始造数以 `sg_order_forecast_control` 为数量源，通过正式订单服务生成并确认 `2年` 订单池，最终停在 `2年订单池已确认、市场投入尚未提交`；验证：已通过 `GOCACHE=.go-build-cache go test ./...`，Go telemetry 仍因用户目录权限输出 token 告警但测试退出码为 0；后续 2026-06-01 已按用户测试诉求将停点前移，当前停点以新记录为准。 |
| 2026-06-01 | 任务状态：✅ 已完成；落点：`cmd/dbtool/order_scenario.go`、`docs/order_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：用户希望测试管理员端从订单数量配置、市场开启到订单池生成确认的完整链路，因此本轮将 `seed-order-scenario` 停点前移：仍造好 3 个小组、`0年 / 1年` 合法经营财报和 `1年` 各市场已选订单历史，但不再预置 `2年` 多年订单数量控制台、市场开启、标段释放顺序、订单池生成与确认；最终停在 `2年已开放、1年订单历史已写入、2年订单配置未开始`；下一步：运行造数命令后，从管理员端手工配置订单数量、市场开启、释放顺序、生成并确认订单池，再进入玩家市场投入和开标验收。 |
| 2026-05-22 | 任务状态：🟡 部分完成；落点：`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/calculation_rule_spec.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/excel_field_mapping.md`、`docs/minimal_state_machine.md`、`docs/order_module_acceptance_checklist.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md`；偏差说明：用户在人工测试中发现玩家端订单页自动轮询会导致页面跳顶、市场预测展示不应为卡片而应贴近 Excel 竖状柱状图、订单金额不能出现小数且必须满足金额等于数量乘以单价。本轮只先完成文档口径，不进入代码实现；下一步：按 `I4-09` 修改玩家端静默轮询、市场预测图表、订单生成金额口径和前后端显示。 |
| 2026-05-22 | 任务状态：✅ 已完成；落点：`frontend/src/stores/player-order.ts`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`internal/service/order_generation_engine.go`、`internal/service/admin_order_service.go`、`internal/service/order_generation_engine_test.go`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I4-09` 代码实现，自动轮询改为静默加载并保留滚动位置；玩家端与管理员端市场预测由文字卡片/金额行补充为阶段-市场竖状柱状图；订单生成改为先按公式得到目标金额，再寻找最接近且可由两位小数单价乘回的整数金额，管理员端/玩家端金额显示为整数。下一步：跑 `go test ./...` 和前端 build，并由用户继续做页面人工复测。 |
| 2026-05-26 | 任务状态：✅ 已完成；落点：`frontend/src/stores/player-order.ts`、`docs/implementation_plan.md`；偏差说明：用户人工测试发现玩家端从 `2年` 切回 `1年` 时，市场投入输入框仍显示 `2年` 未提交草稿。本轮修正玩家端订单页市场投入草稿按年份隔离：切换年份时重置为目标年份后端返回值，同一年自动轮询/刷新时继续保留未提交草稿；验证：已通过 `frontend/npm.cmd run build`；下一步：用户继续按订单开标流程人工复测 `1年/2年` 切换、市场龙头和选单顺序。 |
| 2026-05-26 | 任务状态：✅ 已完成；落点：`internal/service/player_order_service.go`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`docs/implementation_plan.md`；偏差说明：用户人工测试发现管理员端在标段释放前看不到已生成的选单顺序，且顺序表中“市场投入”字段无法区分当前标段投入与上一年该市场订单额。本轮补充管理员选单状态接口返回 `previousMarketOrderAmount`，管理员端标段列表支持点击任一已生成顺序的标段查看顺序，顺序表改为展示“本标段投入 / 上年该市场订单额 / 市场龙头 / 状态 / 已选订单”，市场龙头摘要显示小组名与上年该市场订单额；验证：已通过 `go test ./...` 与 `frontend/npm.cmd run build`；下一步：用户继续验证释放前顺序查看、释放后当前组选择与市场龙头展示是否符合现场使用。 |
| 2026-05-26 | 任务状态：✅ 已完成；落点：`internal/service/game_edition_registry.go`、`internal/model/entity/game_config.go`、`internal/repository/game_config_repository.go`、`internal/service/admin_control_command_service.go`、`internal/service/admin_control_query_service.go`、`internal/service/game_config_query_service.go`、`internal/http/handler/game_config_handler.go`、`migrations/mysql/0009_game_edition.sql`、`cmd/dbtool/main.go`、`frontend/src/views/sandbox-game/admin/setup/AdminSetupPage.vue`、`frontend/src/configs/sandbox-game-service-labels.ts`、`frontend/src/components/sandbox-game/player/OperatingSheet.vue`、`frontend/src/components/sandbox-game/player/ReportSheet.vue`、`frontend/src/components/sandbox-game/player/ReportSidebar.vue`、`frontend/src/views/sandbox-game/admin/baseline/AdminBaselinePage.vue`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I5-01 ~ I5-03`，首版只内置 `VIP_SERVICE_V1`，不开放在线字段编辑或公式编辑；初始化接口要求管理员传入内置版本包编码并锁定到 `sg_game_config`，玩家经营页、财报页和管理员初始基线页切换为贵宾服务版显示名，系统字段名、payload 与公式规则保持不变；验证：已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build`。 |
| 2026-05-26 | 任务状态：✅ 已完成；落点：`docs/requirements_spec.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/excel_field_mapping.md`、`docs/implementation_plan.md`；偏差说明：用户要求先出方案和文档，不立即编码。本轮新增 `I5-04 ~ I5-07`：生产制造版作为独立 `PRODUCTION_V1` 版本包新增，经营页和财报页字段以 `shengchan` 分支生产制造口径为来源；订单字段暂时继续使用 `VIP_ORDER_TEMPLATE_V1` 贵宾订单口径，待生产版订单字段拿到后再新增生产版订单模板；公式规则与流程规则仍共用 `COMMON_FORMULA_V1 / COMMON_PROCESS_V1`。下一步：等用户明确说“开始”后，实现 `I5-05 ~ I5-07`。 |
| 2026-05-26 | 任务状态：✅ 已完成；落点：`internal/service/game_edition_registry.go`、`internal/service/admin_control_command_service_test.go`、`internal/service/admin_control_command_integration_test.go`、`frontend/src/configs/sandbox-game-service-labels.ts`、`frontend/src/components/sandbox-game/player/OperatingSheet.vue`、`frontend/src/components/sandbox-game/player/ReportSheet.vue`、`frontend/src/components/sandbox-game/player/ReportSidebar.vue`、`frontend/src/views/sandbox-game/player/operating/PlayerOperatingPage.vue`、`frontend/src/views/sandbox-game/player/report/PlayerReportPage.vue`、`frontend/src/stores/player-report.ts`、`frontend/src/views/sandbox-game/admin/baseline/AdminBaselinePage.vue`、`docs/implementation_plan.md`；偏差说明：本轮按用户指令实现 `I5-04 ~ I5-07`：新增 `PRODUCTION_V1` 内置版本包并允许赛前初始化选择；前端字段模板扩展为贵宾服务版和生产制造版两套，经营页、财报页、财报必填提示、财报税率说明与初始基线页均按当前比赛 `editionCode` 切换；生产制造版订单字段仍临时沿用 `VIP_ORDER_TEMPLATE_V1`，不改变订单页字段；验证：已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build`。下一步：用户运行页面人工验收生产制造版字段显示。 |
| 2026-05-27 | 任务状态：✅ 已完成；落点：`frontend/src/views/sandbox-game/admin/AdminLayout.vue`、`docs/implementation_plan.md`；偏差说明：用户验收生产制造版字段模板时，希望管理员端顶部状态条直接显示当前比赛版本。本轮在管理员端顶部状态 pill 中新增“版本：{editionName}”，优先读取当前游戏配置，未加载配置时使用赛前状态兜底，不改变赛前选择和初始化锁定逻辑；下一步：运行前端构建验证并由用户页面确认显示位置。 |
| 2026-05-27 | 任务状态：⏳ 未开始；落点：`docs/requirements_spec.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/excel_field_mapping.md`、`docs/calculation_rule_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/minimal_state_machine.md`、`docs/order_module_acceptance_checklist.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md`；偏差说明：本轮仅按用户要求完成新功能方案文档，不进入代码实现。新增 `I6-01 / I6-02`：未开启市场改为玩家端禁用并由系统自动按 `0` 提交；管理员可按 `年份 + 市场` 配置单市场投入上限，空值为无上限；市场投入和上限均必须为非负整数；已有任意小组提交本年市场投入后市场配置锁定；暂不做可用资金总额校验；同时确认全局“用户手工填写数字必须为整数”，但所得税税率下拉、系统计算结果、系统生成订单单价、历史/导入/公式结果展示除外。下一步：用户确认文档后，先实现 `I6-01`。 |
| 2026-05-27 | 任务状态：✅ 已完成；落点：`migrations/mysql/0010_order_market_investment_limit.sql`、`cmd/dbtool/main.go`、`internal/model/entity/order.go`、`internal/repository/order_repository.go`、`internal/service/admin_order_service.go`、`internal/service/player_order_service.go`、`internal/http/dto/admin_order_dto.go`、`internal/http/handler/admin_order_handler.go`、`internal/http/handler/player_order_handler.go`、`internal/app/router.go`、`frontend/src/types/sandbox-game-admin.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`internal/service/order_generation_engine_test.go`、`internal/service/order_workflow_integration_test.go`、`internal/service/player_operating_command_service_test.go`、`internal/service/admin_control_command_integration_test.go`、`docs/implementation_plan.md`；偏差说明：本轮按已确认范围完成 `I6-01`，未扩展到 `I6-02` 全局手工数字整数校验；未开启市场玩家端输入禁用并由提交 payload 自动置 `0`，后端继续硬校验非 `0` 绕过；管理员端新增单市场投入上限输入，订单池确认或已有任意小组提交本年市场投入后锁定；同时补充初始化集成测试的动作日志隔离，避免本地测试库历史记录影响全量测试；验证：已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build`，Go telemetry 仍因用户目录权限打印 token 告警但测试退出码为 0；下一步：进入 I6-02 前先再次确认全局整数输入覆盖范围与页面改动优先级。 |
| 2026-05-27 | 任务状态：✅ 已完成；落点：`internal/service/manual_integer_validation.go`、`internal/service/manual_integer_validation_test.go`、`internal/service/player_operating_command_service.go`、`internal/service/player_report_command_service.go`、`internal/service/admin_control_command_service.go`、`internal/service/admin_notice_command_service.go`、`internal/http/handler/player_operating_handler.go`、`internal/http/handler/player_report_handler.go`、`internal/http/handler/admin_control_handler.go`、`internal/http/handler/admin_notice_handler.go`、`frontend/src/utils/manual-integer.ts`、`frontend/src/types/sandbox-game.ts`、`frontend/src/components/sandbox-game/player/OperatingSheet.vue`、`frontend/src/components/sandbox-game/player/ReportSheet.vue`、`frontend/src/components/sandbox-game/player/ReportSidebar.vue`、`frontend/src/stores/player-operating.ts`、`frontend/src/stores/player-report.ts`、`frontend/src/stores/admin-baseline.ts`、`frontend/src/stores/admin-notice.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`frontend/src/views/sandbox-game/player/report/PlayerReportPage.vue`、`frontend/src/views/sandbox-game/admin/baseline/AdminBaselinePage.vue`、`frontend/src/views/sandbox-game/admin/notice/AdminNoticePage.vue`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/views/sandbox-game/admin/control/AdminControlPage.vue`、`frontend/src/views/sandbox-game/admin/setup/AdminSetupPage.vue`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I6-02` 全局手工数字整数校验；经营页与财报页为保留用户输入的小数字符串，将手工数字类型扩展为 `number | string | ''`，只对手工输入区生效，派生值、系统计算结果、订单单价、历史展示和所得税税率下拉不纳入整数限制；管理员最终年份和赛前小组数量补充前端整数提示，但后端仍沿用原有整数 DTO/业务校验。验证：已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build`，Go telemetry 仍因用户目录权限打印 token 告警但测试退出码为 0；下一步：用户可按经营页/财报页/初始基线/订单页/通知奖惩页面做小数输入人工验收。 |
| 2026-05-27 | 任务状态：✅ 已完成；落点：`docs/requirements_spec.md`、`docs/minimal_state_machine.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/calculation_rule_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/implementation_plan.md`；偏差说明：本轮仅按用户要求把“回退与修正”写回文档，不进入代码实现。新增 `I7-01 ~ I7-07`：管理员端统一 `回退与修正` 入口包含 `退回重提` 与 `恢复快照`；首版只做单组阶段级快照恢复，不做字段级回退和系统自动全局恢复；回退前生成安全快照，旧经营/财报数据保留为失效草稿，订单不释放、不重排、不重算市场龙头，奖惩和汇总按目标节点之后失效，普通通知不回退；存在回退补提 / 待重提小组时阻断开放下一年。下一步：用户确认后，从 `I7-02` 建表与回退日志模型开始实现。 |
| 2026-05-28 | 任务状态：✅ 已完成；落点：`migrations/mysql/0011_rollback_snapshot.sql`、`cmd/dbtool/main.go`、`scripts/build-competition-package.ps1`、`internal/model/entity/state_snapshot.go`、`internal/model/entity/group_year_state.go`、`internal/model/entity/group_summary_snapshot.go`、`internal/model/entity/group_adjustment.go`、`internal/model/entity/admin_unlock_log.go`、`internal/model/entity/order.go`、`internal/enum/rollback.go`、`internal/repository/rollback_repository.go`、`internal/repository/rollback_stage.go`、`internal/repository/group_year_state_repository.go`、`internal/repository/group_adjustment_repository.go`、`internal/repository/summary_snapshot_repository.go`、`internal/repository/order_repository.go`、`internal/service/rollback_errors.go`、`internal/service/admin_control_command_service.go`、`internal/service/admin_notice_command_service.go`、`internal/service/player_order_service.go`、`internal/service/player_operating_command_service_test.go`、`docs/implementation_plan.md`；偏差说明：本轮按 `I7-02` 只完成数据模型、仓储与迁移接入，不实现自动快照生成服务、单组快照恢复接口或管理员端页面；异常解锁日志先写入目标类型 / 目标阶段，`safety_snapshot_id` 等待 `I7-03/I7-04` 生成安全快照后接入；订单回退只新增交付状态失效标记，不释放已选订单归属。验证：已通过 `GOCACHE=.go-build-cache go test ./...`。下一步：进入 `I7-03` 自动快照生成服务。 |
| 2026-06-02 | 任务状态：✅ 已完成；落点：`docs/test_demo_commands.md`、`docs/README.md`、`docs/implementation_plan.md`；偏差说明：本轮按用户测试和展示需要新增命令速查文档，集中记录前后端启动命令、`reset-competition` 赛前选择版本停点、`seed-order-scenario -confirm-reset` 订单开标测试停点、多账号同时测试建议；未修改业务规则或代码；下一步：用户可按该文档在本地或同事电脑做演示和功能验收。 |
| 2026-06-02 | 任务状态：✅ 已完成；落点：`internal/service/rollback_snapshot_service.go`、`internal/service/admin_rollback_service.go`、`internal/http/handler/admin_rollback_handler.go`、`internal/http/dto/admin_rollback_dto.go`、`internal/app/router.go`、`internal/repository/*` 回退相关补充方法、`frontend/src/views/sandbox-game/admin/AdminRollbackPage.vue`、`frontend/src/stores/admin-rollback.ts`、`frontend/src/api/sandbox-game/admin-rollback.ts`、`frontend/src/router/index.ts`、`frontend/src/components/sandbox-game/admin/AdminNav.vue`、玩家经营/财报视图与汇总页、`docs/rollback_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I7-03 ~ I7-07`，把退回重提并入安全快照与回退日志链路，新增单组快照恢复接口和管理端页面；首版仍不做全局快照恢复、不释放订单归属、不重排选单、不重算市场龙头；跨年单组恢复会保留旧草稿和历史记录，并把目标年份之后的状态置为需重新推进。验证：已通过 `go test ./...` 与 `npm run build`。下一步：按 `docs/rollback_module_acceptance_checklist.md` 做人工联调验收。 |
| 2026-06-03 | 任务状态：✅ 已完成；落点：`cmd/dbtool/main.go`、`cmd/dbtool/rollback_scenario.go`、`docs/test_demo_commands.md`、`docs/implementation_plan.md`；偏差说明：用户确认回退模块测试数据必须是非零、可跑通、能触发真实快照动作的标准演练数据，不能只插全 0 占位数据或伪造快照。本轮新增 `seed-rollback-scenario -confirm-reset` 作为第三个测试停点：先按真实年度控制服务开放 `1年 / 2年` 生成全局快照，再通过订单服务生成并确认 `2年` 订单池、完成本地市场两个标段、让第一组完成订单交付，最后通过玩家经营 / 财报服务完成 `2年` 三组提交，生成单组 / 全局自动快照，用于验收退回重提、恢复快照、待重提阻断和订单交付失效。验证：已通过 `go test ./cmd/dbtool`，并实际执行 `go run ./cmd/dbtool -config configs/local.yaml -action seed-rollback-scenario -confirm-reset`，生成 `3` 个小组、`23` 条快照、`1` 条已选订单、`1` 条已交付订单，停在 `2年已完成、尚未开放3年`。下一步：按命令速查和回退验收清单做人工联调。 |
| 2026-06-04 | 任务状态：✅ 已完成；落点：`frontend/src/stores/admin-rollback.ts`、`frontend/src/views/sandbox-game/admin/AdminRollbackPage.vue`、`frontend/src/components/sandbox-game/common/YearTabs.vue`、`frontend/src/utils/sandbox-game-display.ts`、`internal/service/admin_control_command_service.go`、`internal/service/admin_control_command_service_test.go`、`internal/service/game_config_query_service.go`、`internal/service/game_config_query_service_test.go`、`docs/minimal_state_machine.md`、`docs/requirements_spec.md`、`docs/api_design.md`、`docs/calculation_rule_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/implementation_plan.md`；偏差说明：本轮优化退回重提 UX：移除管理员填写原因的必填限制，后端在原因为空时自动填入默认值 `管理员退回重提`；退回重提表单移除原因文本框；提交成功与失败均通过页面顶部 `pageMessage` 明确展示结果；新增 `ROLLBACK_PENDING` 年份 Tab 状态，玩家年份 Tab 在待重提时显示 `待重提` 标签并以橙黄色高亮；补充对应服务层测试与状态机/需求文档同步。验证：已通过 `go test ./internal/service/...` 与 `npm run build`（零报错）。下一步：重启后端后，在回退与修正页面做人工验收，重点验证退回重提无原因可提交、成功后小组年份 Tab 显示 `待重提`。 |
| 2026-06-12 | 任务状态：✅ 已完成；落点：`cmd/dbtool/airport_order_scenario.go`、`docs/test_demo_commands.md`、`docs/implementation_plan.md`；偏差说明：用户执行 `seed-airport-order-leader-scenario -confirm-reset` 时发现脚本在开放 `1年` 前被“初始基线尚未提交”规则拦截。本轮修正机场订单造数脚本：初始化 `AIRPORT_V1` 比赛后先通过正式 `SubmitInitialBaseline` 服务提交一份测试用初始基线，再标记 `0年` 完成并开放 `1年 / 2年`，避免绕过年度控制规则；同步在测试命令文档中说明机场停点会自动提交测试初始基线。验证：已通过 `go test ./cmd/dbtool`；未直接执行带 `--confirm-reset` 的造数命令，以免清空用户当前测试库状态。下一步：用户可重新运行机场造数命令验证停点。 |
| 2026-06-12 | 任务状态：✅ 已完成；落点：`migrations/mysql/0014_order_unit_price_precision.sql`、`cmd/dbtool/main.go`、`scripts/build-competition-package.ps1`、`docs/implementation_plan.md`；偏差说明：用户测试机场订单时发现国内窄体单价显示为 `0` 但订单金额正常。排查确认机场单价本身可能小于 `0.005`，原 `sg_order_pool.unit_price DECIMAL(18,2)` 会在入库时丢失精度；按用户确认，本轮只调整字段精度，不回填当前测试库旧数据。新增 `0014` 迁移将 `unit_price` 调整为 `DECIMAL(18,6)`，并接入 `dbtool` 与正式打包脚本。验证：已通过 `go test ./cmd/dbtool`。下一步：用户清库后重新运行机场造数并生成订单池，验证单价显示为小数。 |
| 2026-06-16 | 任务状态：✅ 已完成；落点：`cmd/dbtool/rollback_scenario.go`、`docs/test_demo_commands.md`、`docs/implementation_plan.md`；偏差说明：用户测试恢复快照时发现原 `seed-rollback-scenario` 没有每个小组 `0年 / 1年` 单组快照，且历史数据过于接近空值，无法验证跨年恢复和失效草稿保留。本轮将 `0年 / 1年` 历史年由直接写库改为通过玩家经营阶段提交和财报提交服务生成，`1年` 先生成并确认全 `0` 轻量订单前置以满足正式年份 `Q1` 提交检查且不预置抢单流程；每组 `0年 / 1年 / 2年` 均自检 `Q1 / Q2 / Q3 / Q4 / YEAR_END / REPORT` 自动快照，经营 payload 按年份/小组/阶段生成差异化非零值，财报手工项按试算平衡后拆分为可辨识的在制品/成品/材料值。验证：已通过 `GOCACHE=.go-build-cache go test ./cmd/dbtool`；未直接执行 `seed-rollback-scenario -confirm-reset`，避免清空用户当前测试库。下一步：用户可在测试库运行造数命令，重点验收 `group01 / 1年 / Q2`、`group02 / 0年 / REPORT` 等跨年恢复路径和失效草稿保留。 |
| 2026-06-16 | 任务状态：✅ 已完成；落点：`cmd/dbtool/rollback_scenario.go`、`docs/implementation_plan.md`；偏差说明：用户实际执行新版 `seed-rollback-scenario -confirm-reset` 时，`0年 group01` 财报试算出现 `gap=7.00` 正向缺口，原造数逻辑只按负向缺口生成库存补平，导致命令主动失败；继续验证时又发现 `1年 group03` 在 Q4 前现金流偏紧，年末提交被状态机拦截。本轮在差异化经营 payload 中增加每季度现金安全短贷，并在 Q4 增加等额“新增短贷 + 材料费”平衡缓冲，使现金流不断裂、财报缺口稳定转为可由库存手工项补平；不改变正式业务规则，仅修正测试造数数据。验证：已通过 `GOCACHE=.go-build-cache go test ./cmd/dbtool`，并实际执行 `go run ./cmd/dbtool -config configs/local.yaml -action seed-rollback-scenario -confirm-reset` 成功，摘要为 `snapshots=60 groupSnapshots=54 globalSnapshots=6 selections=1 delivered=1`。下一步：用户可直接在页面验收跨年恢复和失效草稿保留。 |
| 2026-06-16 | 任务状态：✅ 已完成；落点：`internal/service/player_operating_query_service.go`、`internal/service/draft_invalidation_integration_test.go`、`internal/model/entity/dictionary.go`、`docs/implementation_plan.md`；偏差说明：用户人工验收发现第二小组恢复到 `1年 Q1` 后，再打开 `2年经营` 会提示“获取经营页视图失败”。排查确认这是查询层缺少回退待重提兼容：跨年回退后 `1年` 有效财报被撤销，`2年经营` 不能再按正式承接链计算，但作为待重提年份应允许查看保留下来的失效草稿。本轮仅在经营页查询层增加 `rollbackPending` 兜底，缺少上一年有效财报时返回保留经营草稿、待重提状态和空派生值 / 空承接信息，不放松保存与提交链路；同时补齐字典实体字段长度标签，避免集成测试库 AutoMigrate 将索引字段改成长文本。验证：已通过 `GOCACHE=.go-build-cache go test ./internal/service/...` 与 `GOCACHE=.go-build-cache go test ./cmd/dbtool`。下一步：重启后端后，在页面复测 `group02` 恢复到 `1年 Q1` 后仍可查看 `2年经营` 保留草稿，但必须重新完成 `1年` 后才能继续提交后续年份。 |
| 2026-06-16 | 任务状态：✅ 已完成；落点：`internal/service/player_operating_query_service.go`、`internal/service/player_report_query_service.go`、`internal/service/admin_group_data_query_service.go`、`internal/service/draft_invalidation_integration_test.go`、`frontend/src/stores/player-operating.ts`、`frontend/src/stores/player-report.ts`、`frontend/src/views/sandbox-game/player/report/PlayerReportPage.vue`、`frontend/src/components/sandbox-game/player/OperatingSidebar.vue`、`frontend/src/components/sandbox-game/player/ReportSidebar.vue`、`docs/implementation_plan.md`；偏差说明：用户进一步确认恢复到 `1年 Q1` 后，`1年财报` 与 `2年财报` 也都无法查看。本轮将回退待重提场景拆成“正式计算链路”和“失效草稿查看链路”：财报页在流程未开放或上游承接断链时，如存在旧财报记录则只读返回保留的手工项和旧计算结果，并标记 `hasInvalidDraft/rollbackPending`；经营页在缺少新上年有效财报时只读返回旧经营草稿、空派生值和空承接信息；管理员组数据查询同步支持同样只读查看。前端停止在缺承接时做本地预览计算，并明确提示旧草稿仅供参考、不提交、不汇总、不承接；当上游重新完成并正式开放财报后，旧手工项仍可作为草稿核对并重新提交。验证：已通过 `GOCACHE=.go-build-cache go test ./internal/service/...` 与 `frontend/npm.cmd run build`。下一步：重启后端并刷新前端后，复测 `group02` 恢复到 `1年 Q1`，确认 `1年/2年` 经营与财报均可查看旧草稿，且只有重新完成上游链路后才能提交生效。 |
| 2026-06-16 | 任务状态：✅ 已完成；落点：`frontend/src/stores/player-operating.ts`、`internal/service/manual_integer_validation_test.go`、`docs/implementation_plan.md`；偏差说明：用户直接重新提交 `1年经营` 时，前端整数校验把订单联动写入经营页的 `marketCode` 字符串误判为小数。本轮不改变“用户手工数字必须为整数”的业务口径，只收窄经营页前端递归校验：仅校验数字和可解析为数字的字符串，非数字字符串视为编码/名称类元数据跳过；后端补充 `marketCode` 元数据不触发整数校验的回归测试。验证：已通过 `GOCACHE=.go-build-cache go test ./internal/service -run TestValidate` 与 `frontend/npm.cmd run build`。下一步：用户可复测回退到 `1年 Q1` 后不修改数据直接重新提交经营页。 |
| 2026-06-16 | 任务状态：✅ 已完成；落点：`internal/repository/operating_repository.go`、`internal/repository/report_repository.go`、`internal/service/player_operating_command_service.go`、`internal/service/player_report_command_service.go`、`internal/enum/rollback.go`、`internal/service/player_operating_command_service_test.go`、`internal/service/player_report_command_service_test.go`、`docs/implementation_plan.md`；偏差说明：用户复测回退后重新提交 `1年财报` 时，页面平衡差额为 `0` 但提交接口返回 `500`。排查确认风险点在回退后保留历史提交流水，而年度状态恢复到旧提交版本，重新提交可能撞上 `group/year/version` 唯一键。本轮不删除历史流水，改为经营阶段和财报提交均基于历史最大提交版本递增；当提交前处于 `rollback_pending` 时，自动快照 trigger 改为 `ROLLBACK_STAGE_RESUBMITTED / ROLLBACK_REPORT_RESUBMITTED`，描述显示“回退后重新提交...自动快照”，便于管理员区分首次提交与补提快照。验证：已通过 `GOCACHE=.go-build-cache go test ./internal/service -run "TestSubmit(PlayerReport|OperatingStage)|TestReportViewsForRollback|TestOperatingViewForRollback|TestFindEffective"` 与 `GOCACHE=.go-build-cache go test ./internal/service/...`。下一步：重启后端后复测第二小组从 `1年 Q1` 补提到 `1年财报`，确认接口不再 500，且回退与修正页面能看到补提生成的新自动快照。 |
| 2026-06-16 | 任务状态：✅ 已完成；落点：`frontend/src/views/sandbox-game/admin/AdminRollbackPage.vue`、`docs/api_design.md`、`docs/rollback_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：用户确认回退重提重新提交也需要产生快照，并且页面要能看出是重新提交产生。本轮在管理员 `回退与修正 -> 恢复快照` 列表新增快照说明列，详情中新增生成节点和快照说明；新增 trigger 前端中文映射 `ROLLBACK_STAGE_RESUBMITTED = 回退补提经营`、`ROLLBACK_REPORT_RESUBMITTED = 回退补提财报`，并同步接口文档与验收清单。验证：已通过 `frontend/npm.cmd run build`。下一步：重启前后端后，在回退待重提状态下分别重提经营阶段和财报，确认快照列表显示补提节点与说明。 |
| 2026-06-17 | 任务状态：✅ 已完成；落点：`frontend/src/views/sandbox-game/admin/group-data/AdminGroupDataPage.vue`、`frontend/src/stores/admin-group-data.ts`、`frontend/src/views/sandbox-game/admin/AdminRollbackPage.vue`、`frontend/src/components/sandbox-game/admin/AdminNav.vue`、`internal/service/admin_control_command_service.go`、`internal/http/handler/admin_control_handler.go`、`docs/implementation_plan.md`；偏差说明：本轮按用户确认统一功能使用位置，组数据页只保留按组 / 年 / 页面类型只读查看，删除原“异常解锁”按钮、侧栏、弹窗和最近解锁记录；退回重提继续统一在 `回退与修正` 页面发起，且无需填写原因，原因为空时仍由后端默认记录为 `管理员退回重提`；同步移除后端未使用的“解锁原因必填”错误口径残留；不改变恢复快照仍需填写原因的规则。验证：已通过 `npm.cmd run build` 与 `GOCACHE=E:\project\sand box game\.go-build-cache go test ./internal/service -run TestEnsureUnlock -count=1`。下一步：重启前后端后，在页面确认组数据页无退回入口、回退与修正页退回重提原因可空。 |
| 2026-06-18 | 任务状态：✅ 已完成；落点：`frontend/src/views/sandbox-game/player/operating/PlayerOperatingPage.vue`、`frontend/src/views/sandbox-game/player/report/PlayerReportPage.vue`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`frontend/src/components/sandbox-game/player/OperatingSidebar.vue`、`frontend/src/components/sandbox-game/player/ReportSidebar.vue`、`frontend/src/components/sandbox-game/player/ReportSheet.vue`、`docs/implementation_plan.md`；偏差说明：按用户要求清理玩家端页面中的解释性描述文案，删除经营页、财报页、订单页及侧栏中的说明句，仅保留必要标题、状态与交互控件，不改规则与接口；验证：已通过 `frontend/npm.cmd run build`；下一步：如后续还有同类描述文案，可继续按页面范围批量清理。 |





