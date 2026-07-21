# I12 奖罚任意时段与不可撤销破产验收清单

## 1. 数据库与接口

- [ ] 执行 `0015_adjustment_lifecycle.sql` 后，`sg_group_adjustment` 存在作废人、作废原因和作废时间字段。
- [ ] `sg_group_adjustment` 存在 `idx_group_year_effective(group_id, year_no, effective)`。
- [ ] `sg_group_adjustment_revision` 存在，且 `(group_id, year_no)` 为唯一键。
- [ ] `preview-adjustment`、`send-adjustment`、`void-adjustment` 和 `get-adjustment-sync` 均可访问。
- [ ] 下发请求不包含 `stageCode`；响应中的 `resolvedStageCode` 或 `stageCode` 由服务端决定。

## 2. 阶段自动归属

| 目标组状态 | 预期归属 | 是否允许下发/作废 |
|---|---|---|
| `Q1_OPEN` | `Q1` | 允许 |
| `Q2_OPEN` | `Q2` | 允许 |
| `Q3_OPEN` | `Q3` | 允许 |
| `Q4_OPEN` | `Q4` | 允许 |
| `YEAR_END_OPEN` | `YEAR_END` | 允许 |
| `REPORT_PENDING + REPORT_OPEN` | `YEAR_END` | 允许 |
| `REPORTING + REPORT_OPEN` | `YEAR_END` | 允许 |
| 财报正式提交或年份完成 | 无 | 禁止，先退回 |
| 年份未开放 | 无 | 禁止 |
| 小组已破产 | 无 | 永久禁止 |

## 3. 影响预览与权威执行

- [ ] 创建预览显示归属阶段、调整前现金、调整后现金、税前利润、所得税、净利润和所有者权益。
- [ ] 作废预览按移除目标奖罚后的有效事件集合计算。
- [ ] 预览显示本次所依据的经营/财报草稿保存时间。
- [ ] 正式下发或作废时再次读取最新状态和草稿，不直接复用旧预览结果。
- [ ] 奖励按税前收入、罚款按税前支出进入所得税、净利润、现金、资产、权益和得分计算。
- [ ] 金额必须是大于 `0` 的整数；空原因、负数、`0` 和小数均被拒绝。

## 4. 年末与玩家只读展示

- [ ] 玩家经营页奖罚区显示 `Q1 / Q2 / Q3 / Q4 / 年末 / 总计`。
- [ ] 财报阶段下发的奖励或罚款进入年末列和年度总计。
- [ ] `YEAR_END` 奖罚改变年末现金和财报结果，但不反向改变 Q1 至 Q4 季末现金。
- [ ] 折现费用的年末位置显示 `--`，不保存年末值，也不计入年度总计。
- [ ] 玩家不能编辑奖励和罚款，只能查看汇总及有效/已作废/快照失效明细。

## 5. 作废与 revision

- [ ] 作废不编辑或物理删除原奖罚记录。
- [ ] 作废后记录 `voided_by_id`、`voided_by_name`、`void_reason` 和 `voided_at`。
- [ ] 首次下发后 revision 为 `1`，作废后为 `2`，快照恢复造成有效状态变化后继续递增。
- [ ] 玩家携带相同 revision 检查时返回 `notModified=true`。
- [ ] 玩家携带旧 revision 检查时返回最新奖罚聚合、经营派生结果、财报预览、通知区和平衡/破产状态。

## 6. 破产快照

- [ ] 下发罚款导致所得税后现金小于 `0` 时，同一事务内将小组标记为 `BANKRUPT`。
- [ ] 作废奖励导致所得税后现金小于 `0` 时，同样进入 `BANKRUPT`。
- [ ] 生成 `GROUP + AUTO + ADJUSTMENT_BANKRUPTCY` 快照，metadata 包含奖罚 ID、操作类型、财务影响、经营结果、财报预览、破产时间和原因。
- [ ] 破产快照显示为只读审计快照，恢复按钮禁用，后端恢复接口也拒绝。
- [ ] 破产快照不创建正式财报提交，不创建有效汇总，不进入排名。
- [ ] 破产后下发和作废均被拒绝；普通退回和恢复快照均不能撤销破产。

## 7. 回退语义

- [ ] 普通退回重提保留现有有效奖罚，不改变奖罚 revision。
- [ ] 恢复可恢复快照时，奖罚 `effective` 按快照时点恢复，原记录和作废/回退审计字段不物理删除。
- [ ] 快照恢复确实改变奖罚有效状态时，对每个受影响年份分别递增 revision。
- [ ] 已确认破产的小组恢复早期快照后仍保持 `BANKRUPT` 及原破产年份/原因。

## 8. 静默同步体验

- [ ] 经营页和财报页可见时每 `3` 秒检查 revision。
- [ ] 页面隐藏后停止请求，重新可见后立即检查一次。
- [ ] revision 变化时只更新奖罚只读区、经营派生结果、财报计算结果和通知区。
- [ ] 玩家正在编辑的手工输入不被覆盖，输入焦点和光标位置不改变。
- [ ] 页面滚动位置、当前路由和当前年份不改变，不触发整页 `loadView()`。
- [ ] 右下角提示约 `5` 秒后自动消失，不插入会改变页面高度的顶部消息条。

## 9. 自动化验证

```powershell
$env:GOCACHE='E:\project\sand box game\.go-build-cache'
go test ./internal/rules/operating ./internal/rules/report ./internal/service -skip '^TestOrderWorkflowCoversGenerationSequenceSelectionDeliveryAndUnfinished$' -count=1
go test ./... -skip '^TestOrderWorkflowCoversGenerationSequenceSelectionDeliveryAndUnfinished$' -count=1

Set-Location frontend
npm.cmd run build
```

已知例外：未排除时，仓库既有 `TestOrderWorkflowCoversGenerationSequenceSelectionDeliveryAndUnfinished` 市场龙头用例可能失败；该问题与 I12 奖罚实现无关，需单独跟踪。
