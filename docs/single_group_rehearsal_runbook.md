# 单组演练环境使用说明

> 更新日期：2026-04-01  
> 适用范围：`E:\project\sand box game` 当前仓库中的单组演练环境。  
> 文档定位：本文件专门说明“如何使用单组演练库反复跑完整游戏主链”。

## 1. 这份文档解决什么问题

单组演练环境的目的不是改变正式比赛规则，而是提供一套独立的测试环境，让我们可以：

- 只用 `1` 个管理员账号和 `1` 个玩家组账号
- 从“管理员提交初始基线”开始，完整推演 `0年 -> 最终年`
- 反复重置环境，重新测试
- 不污染多组演练库和正式比赛库

如果你下午要做“全过程推演测试”，优先使用这套环境。

---

## 2. 单组演练环境包含什么

当前单组演练环境由以下文件组成：

- 配置文件：`configs/local-single.yaml`
- 单组种子数据：`migrations/mysql/0002_seed_single_group.sql`
- 初始化脚本：`scripts/init-single-rehearsal.ps1`
- 重置脚本：`scripts/reset-single-rehearsal.ps1`
- 数据库工具：`cmd/dbtool/main.go`

数据库名：

- `sandbox_game_single_rehearsal`

默认账号：

- 管理员：`admin / 123456`
- 玩家：`group01 / 123456`

默认初始状态：

- 只有 `1` 个小组
- `0年` 已存在
- `1~8年` 状态已初始化
- `初始基线未提交`

这意味着：

- 你登录后，需要先由管理员提交初始基线
- 然后才能开始完整测试 `0年` 经营与财报

---

## 3. `init` 和 `reset` 的区别

### 3.1 `init-single-rehearsal.ps1`

用途：

- 首次初始化单组演练库
- 如果数据库还没有建好，可以用它补齐
- 如果表结构和基础数据缺失，可以再次执行

特点：

- 不主动清空已有数据
- 更适合“首次准备环境”

执行命令：

```powershell
Set-Location 'E:\project\sand box game'
powershell -ExecutionPolicy Bypass -File '.\scripts\init-single-rehearsal.ps1'
```

### 3.2 `reset-single-rehearsal.ps1`

用途：

- 删除并重建单组演练库
- 清空之前的推演结果，重新开始

特点：

- 会把单组演练库恢复到干净初始态
- 更适合“我要重新测一遍”

执行命令：

```powershell
Set-Location 'E:\project\sand box game'
powershell -ExecutionPolicy Bypass -File '.\scripts\reset-single-rehearsal.ps1'
```

建议：

- 第一次使用时，可直接先跑一次 `reset`
- 后面每次想重新推演，也优先用 `reset`

---

## 4. 使用前准备

请先确认：

1. MySQL 服务已经启动
2. 当前本机 MySQL root 密码仍是 `123456`
3. 本机 Go 环境已经可用
4. 前端依赖已经安装过，`frontend/node_modules` 可正常使用

如果 MySQL 账号、密码或端口变了，需要先同步修改：

- `configs/local-single.yaml`

当前默认连接：

```yaml
mysql:
  dsn: root:123456@tcp(127.0.0.1:3306)/sandbox_game_single_rehearsal?charset=utf8mb4&parseTime=True&loc=Local
```

---

## 5. 最推荐的完整使用流程

这套流程最适合下午做全过程推演测试。

### 第 1 步：重置单组演练库

```powershell
Set-Location 'E:\project\sand box game'
powershell -ExecutionPolicy Bypass -File '.\scripts\reset-single-rehearsal.ps1'
```

成功后，数据库会被重建为干净状态。

### 第 2 步：启动后端

```powershell
Set-Location 'E:\project\sand box game'
$env:GOCACHE='E:\project\sand box game\.gocache'
$env:GOMODCACHE='E:\project\sand box game\.cache\gomod'
go run .\cmd\server\main.go -config .\configs\local-single.yaml
```

启动成功后，访问：

- 后端健康检查：`http://127.0.0.1:8080/healthz`

### 第 3 步：启动前端

```powershell
Set-Location 'E:\project\sand box game\frontend'
npm run dev
```

启动成功后，前端默认地址：

- `http://127.0.0.1:5173`

### 第 4 步：进入统一登录页

浏览器打开：

- `http://127.0.0.1:5173/sandbox-game/login`

### 第 5 步：管理员先提交初始基线

管理员账号登录：

- 用户名：`admin`
- 密码：`123456`

操作顺序：

1. 进入“初始基线”页
2. 填写初始基线数据
3. 提交并锁定

说明：

- 单组演练库默认就是“初始基线未提交”
- 这是故意保留的，用来测试正式主链第一步

### 第 6 步：玩家完成 `0年` 经营与财报

玩家账号登录：

- 用户名：`group01`
- 密码：`123456`

操作顺序：

1. 进入 `0年经营`
2. 依次提交 `Q1 -> Q2 -> Q3 -> Q4 -> 年末`
3. 完成后进入 `0年财报`
4. 提交 `0年财报`

说明：

- 当前系统里“季末核对现金”是展示型实时试算，不阻断提交
- 玩家需要先看页面内实时算出来的值，再自己核对

### 第 7 步：管理员开放 `1年`

管理员登录后：

1. 进入“年度控制”页
2. 开放 `1年`

因为单组库里只有 `1` 个组，所以不会卡在“其他组未完成”的限制上。

### 第 8 步：继续往后推演

后续就按同样方式循环：

1. 玩家完成本年经营
2. 玩家完成本年财报
3. 管理员开放下一年
4. 管理员查看汇总和排名区域

---

## 6. 如果想重新测试怎么办

最简单的方法就是再次执行：

```powershell
Set-Location 'E:\project\sand box game'
powershell -ExecutionPolicy Bypass -File '.\scripts\reset-single-rehearsal.ps1'
```

然后重新启动后端，再重新登录测试。

这会把单组演练库恢复成最初状态：

- 只有 `admin`
- 只有 `group01`
- 初始基线未提交
- 所有经营/财报/汇总结果被清空

---

## 7. 当前推荐的测试顺序

建议按下面顺序测：

1. 管理员登录
2. 提交初始基线
3. 玩家完成 `0年经营`
4. 玩家完成 `0年财报`
5. 管理员开放 `1年`
6. 玩家完成 `1年经营`
7. 玩家完成 `1年财报`
8. 管理员查看汇总页
9. 管理员测试异常解锁
10. 玩家重新提交被解锁的年度数据
11. 再次查看汇总是否恢复

---

## 8. 常见问题

### 8.1 为什么不用多组演练库？

因为多组演练库更适合测：

- 多组同时推进
- 汇总页联动
- 管理员控制多个组的状态
- “全部组完成才能开放下一年”的真实限制

而单组演练库更适合：

- 快速反复跑完整主链
- 排查某个公式或某个流程问题
- 做演示前自测

### 8.2 单组演练是不是放宽了正式规则？

不是。

正式规则没有改。

只是因为这个库里天然只有 `1` 个组，所以它自己就满足“全部未破产组已完成”的开放条件。

### 8.3 为什么默认不提交初始基线？

因为这样更适合全过程推演。

如果默认已经提交，你就测不到“管理员提交初始基线”这一步了。

### 8.4 如果我只想补齐数据库，不想删历史数据怎么办？

用：

```powershell
powershell -ExecutionPolicy Bypass -File '.\scripts\init-single-rehearsal.ps1'
```

### 8.5 如果我想彻底重来怎么办？

用：

```powershell
powershell -ExecutionPolicy Bypass -File '.\scripts\reset-single-rehearsal.ps1'
```

---

## 9. 一页版命令清单

### 重置数据库

```powershell
Set-Location 'E:\project\sand box game'
powershell -ExecutionPolicy Bypass -File '.\scripts\reset-single-rehearsal.ps1'
```

### 启动后端

```powershell
Set-Location 'E:\project\sand box game'
$env:GOCACHE='E:\project\sand box game\.gocache'
$env:GOMODCACHE='E:\project\sand box game\.cache\gomod'
go run .\cmd\server\main.go -config .\configs\local-single.yaml
```

### 启动前端

```powershell
Set-Location 'E:\project\sand box game\frontend'
npm run dev
```

### 登录页

- `http://127.0.0.1:5173/sandbox-game/login`

### 账号

- `admin / 123456`
- `group01 / 123456`
