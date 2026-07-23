# 测试与展示命令速查

> 适用场景：本地测试、给同事演示、切换到指定测试停点。  
> 注意：本文中的数据库重置 / 造数命令会清空 `configs/local.yaml` 指向的数据库，只能用于测试库。

## 1. 启动前后端

建议先准备好数据库状态，再启动前后端。如果已经启动服务，也可以执行造数命令后刷新页面并重新登录。

### 1.1 启动后端

新开一个 PowerShell 窗口：

```powershell
cd 'E:\project\sand box game'
$env:GOCACHE=(Resolve-Path .go-build-cache).Path
go run .\cmd\server\main.go -config .\configs\local.yaml
```

后端健康检查：

```text
http://127.0.0.1:8080/healthz
```

### 1.2 启动前端

再新开一个 PowerShell 窗口：

```powershell
cd 'E:\project\sand box game\frontend'
npm run dev
```

本机访问地址：

```text
http://127.0.0.1:5173/sandbox-game/login
```

如果给同事在局域网访问，使用 Vite 输出的 `Network` 地址中与你们同一网段的地址，例如：

```text
http://192.168.8.105:5173/sandbox-game/login
```

不要同时混用多个 `Network` 地址测试登录状态。

## 2. 状态一：回到赛前选择版本

用途：

- 测试管理员首次进入系统。
- 测试赛前配置页。
- 测试选择经营版本。
- 测试配置小组数量并初始化比赛。

执行命令：

```powershell
cd 'E:\project\sand box game'
$env:GOCACHE=(Resolve-Path .go-build-cache).Path
go run ./cmd/dbtool -config configs/local.yaml -action reset-competition
```

执行后状态：

- 只保留管理员账号。
- 只初始化游戏配置。
- 不创建玩家组。
- 不创建玩家账号。
- 不创建年份状态。
- 不提交初始基线。
- 不生成订单数据。

登录账号：

```text
admin / 123456
```

管理员登录后应进入：

```text
/sandbox-game/admin/setup
```

也就是“赛前配置 / 选择经营版本 / 配置小组数量”的页面。

## 3. 状态二：跳到订单开标测试停点

用途：

- 跳过真实 `0年 / 1年` 经营填报。
- 保留 `1年` 订单历史，用来测试 `2年` 市场龙头。
- 预置 `3` 个正常组和 `1` 个破产组，便于手工测试龙头过滤、破产过滤和多轮竞标。
- 从管理员端手工测试订单数量配置、市场开启、释放顺序、市场投入与多轮竞标。

执行命令：

```powershell
cd 'E:\project\sand box game'
$env:GOCACHE=(Resolve-Path .go-build-cache).Path
go run ./cmd/dbtool -config configs/local.yaml -action seed-order-scenario -confirm-reset
```

执行后状态：

- 初始化 `4` 个小组，其中第 `4` 组为破产组。
- 写入合法但非真实的 `0年 / 1年` 经营和财报数据。
- 写入 `1年` 各市场已选订单历史。
- 当前开放年份为 `2年`。
- `2年` 订单数量、市场开启、释放顺序、订单池和市场投入均由管理员在页面手工配置。

`2年` 市场龙头测试数据：

| 市场 | 第一组 | 第二组 | 第三组 | 第四组 | 2年龙头 |
|---|---:|---:|---:|---:|---|
| 本地市场 | 120 | 80 | 60 | 40 | 第一组 |
| 区域市场 | 70 | 130 | 90 | 60 | 第二组 |
| 全国市场 | 65 | 75 | 150 | 50 | 第三组 |
| 全球市场 | 140 | 100 | 110 | 9999 | 第一组 |

登录账号：

```text
admin / 123456
group01 / 123456
group02 / 123456
group03 / 123456
group04 / 123456
```

推荐测试顺序：

1. 管理员先配置多年订单数量控制台。
2. 管理员配置 `2年` 市场开启状态。
3. 管理员配置 `2年` 标段释放顺序。
4. 管理员生成并确认 `2年` 订单池。
5. 三个正常小组分别提交 `2年` 市场投入。
6. 管理员生成选单顺序。
7. 管理员释放标段。
8. 当前轮到的小组选择或放弃订单。
9. 验证市场龙头和上一年总额是否在顺序表中显示正确。
10. 验证已选订单对后续小组置灰锁定。

## 4. 状态三：跳到回退与修正测试停点

用途：

- 测试管理员端 `回退与修正` 页面。
- 测试自动快照列表、手动快照、退回重提、单组快照恢复。
- 测试回退后的 `待重提`、年度控制阻断、玩家页回退提示。
- 测试订单交付失效但已选订单归属不释放。

执行命令：

```powershell
cd 'E:\project\sand box game'
$env:GOCACHE=(Resolve-Path .go-build-cache).Path
go run ./cmd/dbtool -config configs/local.yaml -action seed-rollback-scenario -confirm-reset
```

执行后状态：

- 初始化 `3` 个小组，账号仍为 `group01 / group02 / group03`。
- 通过玩家经营 / 财报正式服务完成 `0年 / 1年 / 2年` 三组提交，经营与财报数据均为非零、可通过系统校验的演练数据。
- `1年` 会先通过订单服务生成并确认一份全 `0` 轻量订单前置，用于满足正式年份 `Q1` 提交前置检查，不预置 `1年` 抢单流程。
- 通过订单服务生成并确认 `2年` 订单池，生成全局自动快照。
- 通过开标服务完成至少一个标段，生成标段完成全局自动快照。
- 每个小组的 `0年 / 1年 / 2年` 均生成 `Q1 / Q2 / Q3 / Q4 / YEAR_END / REPORT` 单组自动快照。
- 通过年度控制服务依次开放到 `1年 / 2年`，生成开放下一年全局自动快照。
- 本地市场会完成 `代办过检`、`两舱贵宾` 两个标段；第一组有 `1` 条已选并已交付订单。
- 最终停在 `2年已完成、尚未开放3年`，每个小组都有可用于跨年恢复的单组历史快照，且有订单交付记录可测试回退失效。
- 命令跑通后的控制台摘要通常显示 `snapshots=60 groupSnapshots=54 globalSnapshots=6 selections=1 delivered=1`；其中 `groupSnapshots=54` 是硬性自检基线。

验收重点：

1. 管理员进入 `回退与修正 -> 恢复快照`，应能看到 `GROUP + AUTO` 和 `GLOBAL + AUTO` 快照。
2. 筛选单组快照并查看详情，确认载荷中包含小组、年份、阶段、经营 / 财报 / 汇总 / 订单选择等摘要。
3. 推荐筛选并恢复 `group01 / 1年 / Q2` 或 `group02 / 0年 / REPORT`，验证从 `2年已完成` 跨年恢复到历史节点。
4. 在 `退回重提` 中选择小组、年份和目标阶段，提交后汇总页对应年份显示 `待重提`。
5. 存在 `待重提` 小组时，年度控制页应阻断开放下一年。
6. 玩家重新进入经营页或财报页，应看到回退提示，原经营 / 财报数据保留为可修改草稿，且不是空白或全 `0`。
7. 恢复单组快照时必须填写原因并输入 `确认恢复`，恢复前系统会自动生成安全快照。
8. 回退后已选订单仍归属原小组，但目标节点之后的订单交付状态失效。

说明：

- 这个停点不是全 0 占位数据。
- `0年 / 1年 / 2年` 单组快照不是手工伪造，而是由经营阶段提交、财报提交等真实业务动作触发。
- 全局快照由订单池确认、标段完成、开放下一年等真实业务动作触发。
- 该命令会清空 `configs/local.yaml` 指向的数据库，只能用于测试库。

## 5. 状态四：机场订单模块测试停点

用途：

- 测试 `机场沙盘版 V1` 的订单模块独立联调。
- 从管理员端手工配置国内 / 国际、窄体 / 宽体的订单数量、市场开启、标段顺序、订单池和开标流程。
- 不测试机场经营页和财报页；当前这两个页面只显示待接入占位。

### 5.1 停在 1 年订单配置前

执行命令：

```powershell
cd 'E:\project\sand box game'
$env:GOCACHE=(Resolve-Path .go-build-cache).Path
go run ./cmd/dbtool -config configs/local.yaml -action seed-airport-order-scenario -confirm-reset
```

执行后状态：

- 初始化 `3` 个小组。
- 当前比赛版本为 `AIRPORT_V1`。
- 自动提交一份测试用初始基线，用于通过正式年度开放规则。
- `0年` 已完成，并已开放 `1年`。
- `1年` 订单数量、市场开启、释放顺序、订单池均未配置。
- 国内市场默认开启，国际市场默认关闭，但管理员仍可手动调整。

### 5.2 停在 2 年订单配置前，带 1 年市场龙头历史

执行命令：

```powershell
cd 'E:\project\sand box game'
$env:GOCACHE=(Resolve-Path .go-build-cache).Path
go run ./cmd/dbtool -config configs/local.yaml -action seed-airport-order-leader-scenario -confirm-reset
```

执行后状态：

- 初始化 `3` 个小组。
- 当前比赛版本为 `AIRPORT_V1`。
- 自动提交一份测试用初始基线，用于通过正式年度开放规则。
- 写入 `1年` 国内 / 国际市场已选订单历史。
- 当前开放年份为 `2年`。
- `2年` 订单数量、市场开启、释放顺序、订单池均未配置。

机场版推荐测试顺序：

1. 管理员进入订单管理，确认数量控制台只显示国内市场 / 国际市场 × 窄体 / 宽体。
2. 配置 `1年` 或 `2年` 各标段订单数量，数量必须为非负整数且无固定业务上限；机场版仍按单轮流程测试。
3. 配置市场开启状态，验证国内默认开启、国际默认关闭。
4. 保存标段释放顺序，生成并确认订单池。
5. 玩家端提交市场投入，确认只需提交 `4` 项。
6. 管理员生成选单顺序并释放标段。
7. 玩家按顺序选择或放弃订单，验证已选订单对后续小组置灰。
8. 机场订单交付按钮应置灰或提示暂不支持交付。
9. 玩家进入经营页 / 财报页，应看到机场版待接入占位。

## 6. 多账号同时测试建议

同一个浏览器的多个标签页会共用登录态，不能同时登录多个账号。

推荐做法：

```text
Chrome：admin
Edge：group01
Firefox：group02
Chrome 无痕窗口：group03
```

如果两台电脑一起测，也可以拆开：

```text
你的电脑：admin + group01
同事电脑：group02 + group03
```

关键是每个账号使用独立浏览器环境。

## 7. 停点区别

| 命令 | 停点 | 适合测试 |
|---|---|---|
| `reset-competition` | 赛前未初始化，只能管理员登录 | 选择经营版本、小组数量初始化 |
| `seed-order-scenario -confirm-reset` | 已到 `2年`，但订单配置未开始 | 订单数量、市场开启、订单池生成、开标抢单 |
| `seed-rollback-scenario -confirm-reset` | `0年 / 1年 / 2年` 均已有非零经营 / 财报 / 自动快照历史 | 跨年快照恢复、失效草稿保留、退回重提、待重提阻断、订单交付失效 |
| `seed-airport-order-scenario -confirm-reset` | 机场版 `1年` 订单配置前 | 机场版订单数量、市场开启、订单池生成、开标抢单 |
| `seed-airport-order-leader-scenario -confirm-reset` | 机场版 `2年` 订单配置前，带 `1年` 订单历史 | 机场版国内 / 国际市场龙头和排序 |

## 8. 最简命令汇总

### 8.1 启动后端

```powershell
cd 'E:\project\sand box game'
$env:GOCACHE=(Resolve-Path .go-build-cache).Path
go run .\cmd\server\main.go -config .\configs\local.yaml
```

### 8.2 启动前端

```powershell
cd 'E:\project\sand box game\frontend'
npm run dev
```

### 8.3 回到最开始：选择经营版本

```powershell
cd 'E:\project\sand box game'
$env:GOCACHE=(Resolve-Path .go-build-cache).Path
go run ./cmd/dbtool -config configs/local.yaml -action reset-competition
```

### 8.4 跳到订单开标测试

```powershell
cd 'E:\project\sand box game'
$env:GOCACHE=(Resolve-Path .go-build-cache).Path
go run ./cmd/dbtool -config configs/local.yaml -action seed-order-scenario -confirm-reset
```

### 8.5 跳到回退与修正测试

```powershell
cd 'E:\project\sand box game'
$env:GOCACHE=(Resolve-Path .go-build-cache).Path
go run ./cmd/dbtool -config configs/local.yaml -action seed-rollback-scenario -confirm-reset
```

### 8.6 跳到机场订单模块测试

```powershell
cd 'E:\project\sand box game'
$env:GOCACHE=(Resolve-Path .go-build-cache).Path
go run ./cmd/dbtool -config configs/local.yaml -action seed-airport-order-scenario -confirm-reset
```

### 8.7 跳到机场订单龙头测试

```powershell
cd 'E:\project\sand box game'
$env:GOCACHE=(Resolve-Path .go-build-cache).Path
go run ./cmd/dbtool -config configs/local.yaml -action seed-airport-order-leader-scenario -confirm-reset
```

### 8.8 登录地址和账号

```text
本机地址：http://127.0.0.1:5173/sandbox-game/login
局域网地址：http://你的电脑IP:5173/sandbox-game/login

admin / 123456
group01 / 123456
group02 / 123456
group03 / 123456
```
