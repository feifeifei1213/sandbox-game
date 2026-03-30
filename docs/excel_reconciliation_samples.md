# 沙盘经营系统 Excel 对账样例集（首版）

> 更新日期：2026-03-29  
> 适用方式：基于 `game doc/1组 最终版.xlsx`、当前规则实现与现有集成测试夹具，沉淀一组“可反复回归”的固定样例，用来核对系统计算口径是否仍与 Excel 一致。  
> 文档定位：本文件不是完整测试报告，而是首版上线前的“对账样本基线”。

## 1. 使用目标

本样例集用于解决三个问题：

1. 当 Excel 公式或字段名微调后，如何快速判断系统是否仍对齐
2. 当后端规则层重构后，如何用固定输入回归关键结果
3. 当主持人或产品问“这条口径有没有参照样本”时，能拿出一份具体样例

---

## 2. 使用约定

- 样例以“输入条件 + 关键预期输出”组织，不追求逐格穷尽。
- 每条样例都应能被自动化测试或 `.http` 联调脚本复用。
- 若 `1组 最终版.xlsx` 公式发生变化，必须同步更新本文件与对应测试。

---

## 3. 样例清单

### 3.1 `REC-RP-001` 正式年财报平衡样例

用途：

- 对账财报页核心平衡口径
- 回归 `report submit -> summary snapshot` 主链

适用前提：

- 当前年份：`1年`
- 年份类型：`FORMAL`
- 上一年财报已存在，并作为本年承接基数
- 本年经营草稿为空或不影响该样例关注口径

输入条件：

| 类别 | 字段 | 值 |
|---|---|---|
| 上年财报承接 | `reportFactoryAsset` | `40` |
| 上年财报承接 | `reportLineResidual` | `3` |
| 上年财报承接 | `reportDepreciableAsset` | `0` |
| 上年财报承接 | `reportCash` | `35` |
| 上年财报承接 | `reportShortTermLiability` | `20` |
| 上年财报承接 | `reportShareCapital` | `50` |
| 上年财报承接 | `reportRetainedEarnings` | `19` |
| 财报手工项 | `workInProgress` | `6` |
| 财报手工项 | `finishedGoods` | `4` |
| 财报手工项 | `rawMaterials` | `1` |
| 财报手工项 | `incomeTaxRate` | `0` |

关键预期输出：

| 输出字段 | 期望值 |
|---|---|
| `reportTotalAssets` | `89` |
| `reportTotalLiability` | `20` |
| `reportTotalEquity` | `69` |
| `reportTotalLiabilityEquity` | `89` |
| `balanceGap` | `0` |
| `summaryRevenue` | `0` |
| `summaryProfit` | `0` |
| `summaryEquity` | `69` |
| `summaryEffective` | `true` |

当前落地依据：

- 自动化测试：`internal/service/player_report_command_service_test.go`
- 自动化测试：`internal/service/admin_control_command_integration_test.go`

### 3.2 `REC-ADM-001` 最终排名样例

用途：

- 对账管理员汇总页的“年度汇总区 + 最终排名区”
- 验证破产组不会阻塞最终排名开放

适用前提：

- 最终年份：`2年`
- `第1组`、`第3组` 已完成最终年份财报
- `第2组` 已破产，不再阻塞最终排名

输入条件：

| 小组 | 状态 | `equity` | `revenue` | `profit` |
|---|---|---|---|---|
| 第1组 | `NORMAL + COMPLETED` | `68` | `80` | `10` |
| 第2组 | `BANKRUPT + LOCKED` | 不参与有效汇总 | 不参与有效汇总 | 不参与有效汇总 |
| 第3组 | `NORMAL + COMPLETED` | `82` | `92` | `14` |

关键预期输出：

| 输出场景 | 期望 |
|---|---|
| 年度汇总区列表顺序 | 按 `groupNo` 正序展示有效快照 |
| 年度汇总区排名列 | 第1组显示 `2`，第3组显示 `1` |
| 最终排名区列表顺序 | 按 `equity` 倒序展示 |
| 最终排名第 1 名 | 第3组 |
| 最终排名第 2 名 | 第1组 |
| 破产组阻塞性 | 不阻塞最终排名开放 |

当前落地依据：

- 自动化测试：`internal/service/admin_summary_query_integration_test.go`

---

## 4. 维护规则

- 新增或调整样例时，优先补自动化测试，再回写本文件。
- 若某个样例只剩文档、没有自动化测试承接，不应视为稳定对账样例。
- 若页面标签调整但输出口径不变，本文件只需更新显示名，不应更改样例编号。
