# 订单历史年份与订单事实锁定保护开发文档

> 适用模块：订单管理、多年订单数量控制台、市场开启配置、标段释放顺序、订单池生成、回退与快照恢复
> 编写日期：2026-07-27
> 当前状态：首版开发已完成，后续进入页面人工验收

## 1. 背景

在使用 `seed-order-scenario -confirm-reset` 测试多轮订单功能时，出现了一个现象：

- 造数命令已经写入 `1年` 已选订单历史，用于测试 `2年` 市场龙头；
- 管理员在页面配置订单数量、市场开启、释放顺序并生成选单顺序后；
- 顺序表中 `上年该市场订单额` 全部显示为 `0`，`市场龙头` 全部显示为 `否`。

排查后判断：问题不是市场龙头公式本身，而是订单配置刷新动作可能重建了 `1年` 订单池，导致造数命令预置的 `1年 SELECTED 已选订单历史` 被删除。

当前系统已经有“订单池确认后锁定配置”的能力，但锁定条件主要依赖 `confirmed order batch`。这不足以覆盖以下情况：

- 历史年份已经完成，但没有完整的订单确认批次；
- 造数命令手工写入了已选订单历史，但没有写入 confirmed batch；
- 回退 / 快照恢复后，经营、财报、汇总需要失效或补提，但已选订单归属仍应保留；
- 订单池、选单记录、市场投入、标段状态可能因为造数、回退或异常修复出现不完全一致。

因此需要新增一套统一的“订单年份锁定”判断，防止历史年份或已产生订单事实的年份被订单配置动作覆盖。

## 2. 修复目标

本次修复目标可以概括为：

```text
凡是已经成为历史年份，或者已经产生订单事实的年份，都不能再被订单配置刷新订单池。
```

具体目标：

1. `yearNo < currentOpenYear` 的历史年份必须锁定订单配置和订单池重建。
2. 任何已经产生订单事实的年份必须锁定订单配置和订单池重建。
3. 保存多年订单数量控制台、市场开启配置、释放顺序、生成预览订单池前必须统一检查锁定状态。
4. 前端数量控制台和年度订单配置页必须展示锁定状态并禁用对应输入。
5. 回退 / 快照恢复不得释放已选订单归属、不得重排选单、不得重算市场龙头。
6. `seed-order-scenario` 进入 `2年` 后，`1年` 应被识别为历史年份并锁定，避免误删 `1年` 龙头历史。

## 3. 非目标

本次不做以下内容：

- 不修改市场龙头计算公式。
- 不释放已选订单。
- 不重排已经生成的选单顺序。
- 不把已选订单恢复为可选订单。
- 不做全局快照恢复。
- 不新增复杂人工改历史订单入口。
- 不改变回退 / 快照恢复“只影响经营、财报、汇总、交付有效性”的边界。

## 4. 统一订单年份锁定规则

新增统一的订单年份锁定判断，建议在订单服务层封装为类似：

```go
type OrderYearLockStatus struct {
    YearNo     int
    Locked     bool
    ReasonCode string
    ReasonText string
}
```

首版锁定条件如下，任一命中即锁定。

| 锁定条件 | 判断依据 | 推荐原因码 | 推荐展示文案 |
|---|---|---|---|
| 历史年份 | `yearNo < gameConfig.current_open_year` | `HISTORICAL_YEAR` | 该年份已成为历史年份 |
| 回退待重提 | `sg_group_year_state.rollback_pending = true` | `ROLLBACK_PENDING` | 该年份存在回退补提，订单事实保持锁定 |
| 订单池已确认 | `sg_order_generation_batch.batch_status = CONFIRMED` | `ORDER_POOL_CONFIRMED` | 订单池已确认 |
| 已有已选订单 | `sg_order_pool.pool_status = SELECTED` 或 `selected_group_id IS NOT NULL` | `SELECTED_ORDER_EXISTS` | 该年份已产生订单历史 |
| 已有小组选单记录 | `sg_group_order_selection` 存在该年记录 | `GROUP_SELECTION_EXISTS` | 该年份已有小组选单记录 |
| 玩家已提交市场投入 | `sg_group_market_bid` 存在该年记录 | `MARKET_BID_EXISTS` | 该年份已有市场投入 |
| 已生成选单顺序 | `sg_market_selection_order` 存在该年记录 | `SELECTION_SEQUENCE_EXISTS` | 该年份已生成选单顺序 |
| 竞标流程已推进 | `sg_market_bidding_state` 进入非纯配置状态 | `BIDDING_STARTED` | 该年份竞标流程已开始 |

推荐判断优先级：

```text
历史年份
-> 回退待重提
-> 订单池已确认
-> 已选订单事实
-> 小组选单记录
-> 市场投入
-> 选单顺序
-> 标段状态已推进
```

优先级只影响页面展示原因，不影响锁定结果。

## 5. 需要接入锁定判断的后端入口

以下入口只要目标年份被锁定，就必须拒绝修改或重建订单池。

### 5.1 保存多年订单数量控制台

对应服务：

```text
AdminOrderCommandService.UpdateForecastControl
```

当前风险：

- 保存数量控制台会找出变化年份；
- 对变化年份自动调用 `generateOrderPreviewInTx(... Overwrite: true)`；
- 如果变化年份是历史年份或已有订单事实，就可能误删历史订单池。

修复要求：

```text
对每个 changedYear 调用统一订单年份锁定判断；
锁定则拒绝保存，并返回明确错误；
不得进入自动生成预览订单池逻辑。
```

### 5.2 保存市场开启配置

对应服务：

```text
AdminOrderCommandService.UpdateMarketConfig
```

当前风险：

该入口会删除当前年份订单池和标段状态，并重新生成预览。

修复要求：

```text
保存前检查目标年份锁定状态；
锁定则拒绝；
不允许删除订单池和标段状态。
```

### 5.3 保存标段释放顺序

对应服务：

```text
AdminOrderCommandService.UpdateControlConfig
```

当前风险：

该入口会删除当前年份订单池和标段状态，并重新生成预览。

修复要求：

```text
保存前检查目标年份锁定状态；
锁定则拒绝；
不允许删除订单池和标段状态。
```

### 5.4 生成 / 刷新预览订单池

对应服务：

```text
AdminOrderCommandService.GenerateOrderPool
generateOrderPreviewInTx
```

最关键的保护点是 `generateOrderPreviewInTx`，因为它会执行：

```go
poolRepo.DeleteByYear(ctx, cmd.YearNo)
stateRepo.DeleteByYear(ctx, cmd.YearNo)
```

修复要求：

```text
在真正 DeleteByYear 前调用统一锁定判断；
锁定则直接拒绝；
这是底层兜底保护，即使上层漏判，也不能删历史订单池。
```

### 5.5 确认订单池

对应服务：

```text
AdminOrderCommandService.ConfirmOrderPool
```

修复要求：

```text
如果该年份已经有订单事实、已生成选单、已开始竞标或已经是历史年份，则不能确认新的预览批次；
如果只是当前未锁定年份的正常 PREVIEW -> CONFIRMED，则允许确认。
```

## 6. 前端调整

### 6.1 多年订单数量控制台

当前数量控制台的年份锁定主要来自“订单池已确认”。

修复后，后端返回的 `yearLocks` 需要扩展为统一锁定结果：

```text
历史年份也锁
已有订单事实也锁
回退待重提年份也锁
订单池已确认仍然锁
```

前端继续复用现有逻辑：

```ts
isForecastYearLocked(yearNo)
forecastYearLockReason(yearNo)
```

预期效果：

```text
seed-order-scenario 后 currentOpenYear = 2；
1年数量列自动置灰；
管理员不能误改 1年订单数量；
1年已选订单历史不会被刷新订单池覆盖。
```

### 6.2 年度订单配置页

年度订单配置页的以下能力需要按统一锁定结果禁用：

- 市场开启配置保存；
- 市场投入上限编辑；
- 标段释放顺序保存；
- 生成 / 刷新预览订单池；
- 确认不合法预览批次。

建议后端配置查询返回类似：

```json
{
  "orderYearLock": {
    "locked": true,
    "reasonCode": "HISTORICAL_YEAR",
    "reasonText": "该年份已成为历史年份"
  }
}
```

若暂不扩展返回结构，也可以先复用已有 `canUpdateConfig / canGeneratePreview / canConfirmPool`，但错误提示会不够精确。

## 7. 回退与快照恢复边界

本次修复必须保持现有回退边界：

回退 / 快照恢复可以影响：

```text
经营草稿与提交有效性
财报草稿与提交有效性
汇总快照有效性
奖罚有效性
订单交付有效性
```

回退 / 快照恢复不允许影响：

```text
已选订单归属
订单池 SELECTED 状态
小组选单记录
选单顺序
市场龙头历史依据
```

当前代码中的 `InvalidateDeliveryAfterTarget` 只应让交付有效性失效，不应释放订单归属。

## 8. 与年度推进规则的关系

年度控制已经要求：

```text
当前开放年所有未破产小组 YearStatus = COMPLETED；
且不存在 rollback_pending；
才允许开放下一年。
```

因此：

```text
如果 currentOpenYear = 2，
则 1年在业务上必须被视为历史年份；
1年订单配置和订单池重建必须锁定。
```

回退后可能出现：

```text
currentOpenYear = 2；
某小组 1年 rollback_pending = true。
```

这只代表该小组经营 / 财报需要补提，不代表 `1年` 订单配置重新可编辑。订单事实仍然保持锁定。

## 9. 错误提示建议

建议新增或复用统一错误：

```go
ErrAdminOrderYearLocked
```

推荐前端提示：

```text
该年份订单数据已锁定，不能修改订单配置或重新生成订单池。
```

若可带锁定原因，则使用更具体的提示：

- 该年份已成为历史年份，订单配置不可修改。
- 该年份已产生订单历史，不能重新生成订单池。
- 该年份已有市场投入，不能修改订单配置。
- 该年份存在回退补提，订单事实保持锁定。

## 10. 测试要求

### 10.1 历史年份不能刷新订单池

场景：

```text
currentOpenYear = 2
1年存在订单池或订单历史
修改 1年订单数量
```

期望：

```text
接口返回锁定错误
1年订单池不被删除
1年 SELECTED 历史仍存在
```

### 10.2 已有 SELECTED 订单不能生成预览池

场景：

```text
1年 sg_order_pool 存在 SELECTED 订单
调用 GenerateOrderPool(yearNo=1, overwrite=true)
```

期望：

```text
接口返回锁定错误
SELECTED 订单仍存在
selected_group_id 不被清空
```

### 10.3 已有市场投入不能改订单配置

场景：

```text
某年已有 sg_group_market_bid
管理员保存市场开启配置或释放顺序
```

期望：

```text
接口返回锁定错误
订单池和标段状态不被删除
```

### 10.4 回退不释放已选订单

场景：

```text
2年已有已选订单
恢复到 1年 Q1 快照
```

期望：

```text
sg_group_order_selection 仍存在
sg_order_pool.selected_group_id 仍保留
只按目标节点之后失效 delivery_effective
不重排选单
不重算龙头
```

### 10.5 造数命令回归

场景：

```text
运行 seed-order-scenario -confirm-reset
currentOpenYear = 2
尝试修改 1年订单数量
```

期望：

```text
1年被识别为历史年份并锁定
保存被拒绝
1年龙头历史保留
2年生成选单顺序后，本地市场第一组仍为龙头
```

## 11. 验收标准

本功能完成后，需要满足：

1. 历史年份订单数量列置灰，不能编辑。
2. 历史年份市场开启配置和释放顺序不能保存。
3. 已有订单事实的年份不能重新生成预览订单池。
4. `generateOrderPreviewInTx` 不会删除锁定年份订单池。
5. 回退 / 快照恢复不释放已选订单归属。
6. `seed-order-scenario` 测试 `2年` 龙头时，`1年` 历史不会被误删。
7. 市场龙头仍按上一年已固化订单事实计算，不改公式。
