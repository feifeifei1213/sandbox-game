# 沙盘经营系统需求驱动实施计划（首版）

> 更新日期：2026-05-22
> 适用方式：基于当前已确认的业务共识、Excel 规则底稿和原型方向，持续把沙盘经营系统首版需求拆成可执行任务，并同步更新状态。

## 计划规则

- 本计划只收录“首版业务需求任务”与“首版前置落地任务”。
- 每个任务必须包含：优先级、状态、当前产出（文档/原型/代码落点）、下一步。
- 状态定义：`已完成` / `部分完成` / `未开始` / `阻塞`。
- 需求来源以用户明确指令、Excel 规则底稿、现有 Excel 文件为准。
- 未经确认，不得擅自扩大首版范围。

## 任务总览（需求沉淀层当前基线）

- 总任务数：`71`
- 已完成：`69`
- 部分完成：`4`
- 未开始：`0`

## 开发执行层任务总览

- 总任务数：`39`
- 已完成：`35`
- 部分完成：`2`
- 未开始：`2`
- 阻塞：`0`
- 当前状态：`?? 进行中`

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
| 玩家端通知区置于右侧工作栏顶部；普通通知由管理员发送；奖励/罚款改为管理员按组/年/季下发，玩家只读 | P0 | ??? | 用户确认（2026-03-31）；后端、前端、迁移与测试已落地 | 后续只做细节迭代，不再作为首版阻塞项 |

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
| M2 | 补齐管理员运维能力与前端页面落地 | 汇总、组数据查看、异常解锁、Excel 风格页面联调完成 |
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
| M2-06 | 实现异常解锁接口与审计日志 | P0 | ??? | M1-05、M1-09、M1-12、M2-03 | `unlock-year` 接口、`sg_admin_unlock_log`、`sg_admin_action_log` | 满足“仅下一年未开放前可解锁、解锁需写原因、财报失效、汇总撤回” | 无 |
| M2-07 | 搭建玩家端经营页 Excel 风格静态壳子 | P0 | ??? | 现有 demo、`docs/requirements_spec.md` | 玩家经营页前端骨架、年份标签、阶段表格框架 | 当前正式前端页面已作为经营页视觉基线收口；整页 Excel 风格、年份标签、经营/财报切换、右侧轻量工作栏与“未来阶段可见但锁定”口径已落地 | 无 |
| M2-08 | 接入经营页查询、草稿、阶段提交链路 | P0 | ??? | M1-07、M1-08、M1-09、M2-07 | 经营页联调页面 | 已完成真实读取、草稿保存、错误提交拦截与成功提交流程闭环；并通过事务级集成测试验证 `Q1 -> Q2` 推进、阶段流水写入与可编辑范围切换 | 无 |
| M2-09 | 搭建玩家端财报页 Excel 风格静态壳子 | P0 | ?? | 现有 demo、`docs/requirements_spec.md` | 玩家财报页前端骨架 | 财报页已接入正式前端工程并具备 Excel 风格主表、绿色手工项与税率下拉；仍待继续按最终视觉口径验收 | 无 |
| M2-10 | 接入财报查询、草稿、提交链路 | P0 | ??? | M1-10、M1-11、M1-12、M2-09 | 财报页联调页面 | 财报查询、草稿、提交与平衡校验链路已接入正式前端工程，并通过事务级集成测试验证草稿持久化、正式年份提交、汇总快照写入与提交后只读 | 无 |
| M2-11 | 搭建管理员页面基础框架 | P1 | ??? | M2-04、现有 demo | 管理员导航、汇总页、年度控制页、初始基线页 | 管理员正式前端页面已落入 `frontend/` 工程，完成左侧菜单、汇总页、年度控制页、初始基线页与真实接口接入，并通过 `npm run build` | 无 |
| M2-12 | 搭建组数据查看与异常解锁前端 | P1 | ??? | M2-05、M2-06、M2-11 | 组数据查看页、解锁弹窗、日志提示 | 管理员正式前端已支持按组按年查看经营/财报只读视图、切换页面类型、发起异常解锁并展示最近一次原因记录，同时通过 `go test ./...` 与 `npm run build` | 无 |
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
  - 偏差说明：本轮已补齐 `admin-control/unlock-year` 首版闭环，服务端会校验“下一年是否已开放、原因是否为空、当前年份是否仍处于可编辑态”；解锁时会回收财报提交状态、撤回汇总快照，并在命中破产年份时恢复该组 `NORMAL` 状态，同时写入 `sg_admin_unlock_log` 与 `sg_admin_action_log`。当前仍不支持“回滚到某个历史季度快照”的精细回退，只支持最小结果回收后由玩家重新提交。
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
| I4-07 | 增加订单专项测试场景造数能力 | P0 | 已完成 | `I4-08`、用户 2026-05-21 确认 | `cmd/dbtool/main.go`、`cmd/dbtool/order_scenario.go`、`docs/order_module_acceptance_checklist.md`、`docs/implementation_plan.md` | 新增 `dbtool -action seed-order-scenario -confirm-reset` 测试造数命令：执行时会重建当前配置指向数据库并跑完整迁移，随后初始化 3 个小组、管理员与玩家账号、合法但非真实的 `0年 / 1年` 经营与平衡财报结果；按 I4-08 新链路写入 `1~8年` 多年订单数量控制台并生成市场预测快照，再配置、生成并确认 `2年` 订单池；同时造出 `1年` 各市场已选订单结果，用于形成 `2年` 市场龙头；最终停在 `2年订单池已确认、市场投入尚未提交`，不提前提交市场投入、不关闭经营页与财报页校验、不改变正式比赛规则 | 已通过 `GOCACHE=.go-build-cache go test ./...` 验证；偏差说明：出于防误用考虑，造数命令必须显式追加 `--confirm-reset` 才会执行重置，且未自动运行真实数据库重置；下一步可按 `docs/order_module_acceptance_checklist.md` 从该停点做管理员端 / 玩家端人工验收 |
| I4-08 | 重构多年订单数量控制台与市场预测链路 | P0 | 已完成 | 用户 2026-05-21 确认、`道具-订单推算（服务企业）.xlsx`、`I4-06` | `migrations/mysql/0008_order_forecast_control.sql`、`internal/model/entity/order.go`、`internal/repository/order_repository.go`、`internal/service/admin_order_service.go`、`internal/http/dto/admin_order_dto.go`、`internal/http/handler/admin_order_handler.go`、`internal/http/handler/player_order_handler.go`、`internal/app/router.go`、`frontend/src/types/sandbox-game-admin.ts`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/api/sandbox-game/admin-order.ts`、`frontend/src/api/sandbox-game/player-order.ts`、`frontend/src/stores/admin-order.ts`、`frontend/src/stores/player-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`internal/service/order_generation_engine_test.go`、`internal/service/order_workflow_integration_test.go`、`internal/service/player_operating_command_service_test.go`、`tests/http/sandbox-game/AdminOrder.http`、`docs/implementation_plan.md` | 已落地新口径：新增多年订单数量控制台表与市场预测快照表，控制台固定覆盖 `1年~8年 × 四个市场 × 四类产品`；管理员端新增控制台维护区与三阶段预测展示/说明维护；年度订单管理页不再编辑订单数量，只读带入控制台当年数量并保存市场开启与释放顺序；订单生成从控制台读取当年数量，未开启市场仍不生成、不竞标；玩家端在市场投入前只读展示 `1~3年 / 4~5年 / 6~8年` 市场预测 | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build` 验证；偏差说明：上传 Excel 兼容接口继续保留但不作为控制台主链路，市场预测金额首版按已落地订单生成公式参数做确定性估算，后续若继续精确复刻 Excel 的全部随机展示口径需另立任务 |
| I4-09 | 修正订单页人工验收问题：静默轮询、Excel 竖状预测图、整数订单金额口径 | P0 | 已完成 | 用户 2026-05-22 测试反馈、`道具-订单推算（服务企业）.xlsx`、`I4-08` | `frontend/src/stores/player-order.ts`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`internal/service/order_generation_engine.go`、`internal/service/admin_order_service.go`、`internal/service/order_generation_engine_test.go`、`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/calculation_rule_spec.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/excel_field_mapping.md`、`docs/minimal_state_machine.md`、`docs/order_module_acceptance_checklist.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md` | 已完成三项修正：玩家端 `3` 秒轮询改为静默刷新，不触发整页加载态并保留滚动位置；玩家端与管理员端市场预测均改为贴近 Excel 的竖状柱状图，阶段下按市场分图，年份为横轴，四类订单为柱状系列；订单生成按公式目标值寻找最接近的整数金额，再反算两位小数单价，保证 `订单金额 = 数量 × 单价`，管理员端/玩家端金额显示为整数 | 待本轮验证完成后进入页面人工复测：重点检查订单页停留在订单池中等待轮询不跳顶、预测柱状图与 Excel 方向一致、重新生成订单池后金额不出现小数 |

### 9.7 I5：沙盘版本包与贵宾服务版字段模板

| Task ID | 任务 | 优先级 | 状态 | 依赖 | 建议交付物 | 完成标准 | 阻塞情况 |
|---|---|---|---|---|---|---|---|
| I5-01 | 明确沙盘版本包与贵宾服务版字段模板正式口径 | P0 | 已完成 | 用户 2026-05-26 确认、`guibing` 分支字段配置 | `docs/requirements_spec.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/excel_field_mapping.md`、`docs/implementation_plan.md` | 已明确字段模板、公式规则、流程规则拆分；首版先完整落地 `VIP_SERVICE_V1`，绑定贵宾经营页、贵宾财报页、贵宾订单页字段，公式和流程共用当前通用版本；管理员只在赛前初始化时选择内置版本包，初始化后锁定；不做在线字段编辑和在线公式编辑 | 已完成，后续生产版另立版本包 |
| I5-02 | 实现沙盘版本包底座与赛前版本选择 | P0 | 已完成 | `I5-01` | `internal/service/game_edition_registry.go`、`migrations/mysql/0009_game_edition.sql`、`sg_game_config` 版本字段、`get-current/list-game-editions/get-setup-status/initialize-game` 接口、管理员赛前初始化页版本选择 | 管理员可在未初始化时选择 `VIP_SERVICE_V1`；初始化后版本锁定并在配置接口只读展示；已有测试库默认补齐 `VIP_SERVICE_V1`，不影响现有数据 | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build` |
| I5-03 | 接入贵宾服务版经营页与财报页字段模板 | P0 | 已完成 | `I5-02`、`guibing` 分支 | `frontend/src/configs/sandbox-game-service-labels.ts`、`OperatingSheet.vue`、`ReportSheet.vue`、`ReportSidebar.vue`、`AdminBaselinePage.vue`、财报必填提示 | 当前主线经营页、财报页、初始基线页显示贵宾服务版字段；系统字段名、payload、公式计算不变；订单页继续使用贵宾订单字段 | 已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build` |

### 12.1 本轮新增记录

| 日期 | 记录 |
|---|---|
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
| 2026-05-21 | 任务状态：✅ 已完成；落点：`cmd/dbtool/main.go`、`cmd/dbtool/order_scenario.go`、`docs/order_module_acceptance_checklist.md`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I4-07` 订单专项测试场景造数能力，新增 `seed-order-scenario` 动作并要求 `--confirm-reset` 显式确认；造数以 `sg_order_forecast_control` 为数量源，通过正式订单服务生成并确认 `2年` 订单池，最终停在 `2年订单池已确认、市场投入尚未提交`；验证：已通过 `GOCACHE=.go-build-cache go test ./...`，Go telemetry 仍因用户目录权限输出 token 告警但测试退出码为 0；下一步：可用该停点按验收清单做订单模块人工联调。 |
| 2026-05-22 | 任务状态：🟡 部分完成；落点：`docs/requirements_spec.md`、`docs/requirements_consensus_checklist.md`、`docs/calculation_rule_spec.md`、`docs/api_design.md`、`docs/database_design.md`、`docs/excel_field_mapping.md`、`docs/minimal_state_machine.md`、`docs/order_module_acceptance_checklist.md`、`game doc/Excel计算规则与跨表联动说明.md`、`docs/implementation_plan.md`；偏差说明：用户在人工测试中发现玩家端订单页自动轮询会导致页面跳顶、市场预测展示不应为卡片而应贴近 Excel 竖状柱状图、订单金额不能出现小数且必须满足金额等于数量乘以单价。本轮只先完成文档口径，不进入代码实现；下一步：按 `I4-09` 修改玩家端静默轮询、市场预测图表、订单生成金额口径和前后端显示。 |
| 2026-05-22 | 任务状态：✅ 已完成；落点：`frontend/src/stores/player-order.ts`、`frontend/src/views/sandbox-game/player/order/PlayerOrderPage.vue`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`internal/service/order_generation_engine.go`、`internal/service/admin_order_service.go`、`internal/service/order_generation_engine_test.go`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I4-09` 代码实现，自动轮询改为静默加载并保留滚动位置；玩家端与管理员端市场预测由文字卡片/金额行补充为阶段-市场竖状柱状图；订单生成改为先按公式得到目标金额，再寻找最接近且可由两位小数单价乘回的整数金额，管理员端/玩家端金额显示为整数。下一步：跑 `go test ./...` 和前端 build，并由用户继续做页面人工复测。 |
| 2026-05-26 | 任务状态：✅ 已完成；落点：`frontend/src/stores/player-order.ts`、`docs/implementation_plan.md`；偏差说明：用户人工测试发现玩家端从 `2年` 切回 `1年` 时，市场投入输入框仍显示 `2年` 未提交草稿。本轮修正玩家端订单页市场投入草稿按年份隔离：切换年份时重置为目标年份后端返回值，同一年自动轮询/刷新时继续保留未提交草稿；验证：已通过 `frontend/npm.cmd run build`；下一步：用户继续按订单开标流程人工复测 `1年/2年` 切换、市场龙头和选单顺序。 |
| 2026-05-26 | 任务状态：✅ 已完成；落点：`internal/service/player_order_service.go`、`frontend/src/types/sandbox-game-order.ts`、`frontend/src/views/sandbox-game/admin/order/AdminOrderPage.vue`、`docs/implementation_plan.md`；偏差说明：用户人工测试发现管理员端在标段释放前看不到已生成的选单顺序，且顺序表中“市场投入”字段无法区分当前标段投入与上一年该市场订单额。本轮补充管理员选单状态接口返回 `previousMarketOrderAmount`，管理员端标段列表支持点击任一已生成顺序的标段查看顺序，顺序表改为展示“本标段投入 / 上年该市场订单额 / 市场龙头 / 状态 / 已选订单”，市场龙头摘要显示小组名与上年该市场订单额；验证：已通过 `go test ./...` 与 `frontend/npm.cmd run build`；下一步：用户继续验证释放前顺序查看、释放后当前组选择与市场龙头展示是否符合现场使用。 |
| 2026-05-26 | 任务状态：✅ 已完成；落点：`internal/service/game_edition_registry.go`、`internal/model/entity/game_config.go`、`internal/repository/game_config_repository.go`、`internal/service/admin_control_command_service.go`、`internal/service/admin_control_query_service.go`、`internal/service/game_config_query_service.go`、`internal/http/handler/game_config_handler.go`、`migrations/mysql/0009_game_edition.sql`、`cmd/dbtool/main.go`、`frontend/src/views/sandbox-game/admin/setup/AdminSetupPage.vue`、`frontend/src/configs/sandbox-game-service-labels.ts`、`frontend/src/components/sandbox-game/player/OperatingSheet.vue`、`frontend/src/components/sandbox-game/player/ReportSheet.vue`、`frontend/src/components/sandbox-game/player/ReportSidebar.vue`、`frontend/src/views/sandbox-game/admin/baseline/AdminBaselinePage.vue`、`docs/implementation_plan.md`；偏差说明：本轮完成 `I5-01 ~ I5-03`，首版只内置 `VIP_SERVICE_V1`，不开放在线字段编辑或公式编辑；初始化接口要求管理员传入内置版本包编码并锁定到 `sg_game_config`，玩家经营页、财报页和管理员初始基线页切换为贵宾服务版显示名，系统字段名、payload 与公式规则保持不变；验证：已通过 `GOCACHE=.go-build-cache go test ./...` 与 `frontend/npm.cmd run build`。 |





