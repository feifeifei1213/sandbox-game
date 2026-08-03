# 沙盘经营系统文档索引

> 更新日期：2026-04-02  
> 适用目录：`E:\project\sand box game`

## 1. 项目简介

本项目是一个面向企业培训与经营推演场景的沙盘经营系统。

当前首版重点覆盖：

- 玩家端经营页
- 玩家端财报页
- 管理员端汇总页
- 管理员端年度控制
- 管理员端初始基线
- 管理员端组数据查看与异常解锁
- 管理员端通知与奖惩
- 单组演练环境与全过程测试支持

系统定位不是把桌游流程完全电子化，而是先把 Excel 已承载的经营、财报、汇总与必要流程控制落到系统中。

---

## 2. 当前技术栈

- 后端：Go
- 前端：Vue 3 + TypeScript + Vite
- 数据库：MySQL
- 文档：统一收口在 `docs/`

---

## 3. 目录说明

项目中当前重点目录：

- `cmd/`：Go 启动入口与辅助工具
- `internal/`：后端业务代码
- `frontend/`：正式前端工程
- `configs/`：本地配置文件
- `migrations/`：数据库结构与种子数据
- `docs/`：需求、设计、联调、测试、运行说明
- `demo 网页/`：静态页面原型与演示稿
- `game doc/`：Excel、规则文档与原始业务资料

订单模块近期补充设计：

- `docs/order_market_auto_focus_design.md`：订单竞标阶段管理员端和玩家端按真实流程自动聚焦当前市场的开发设计。
- `docs/admin_order_pool_group_filter_design.md`：管理员端订单池按小组筛选查看已选订单的开发设计。

---

## 4. 如何快速启动

### 4.1 启动后端

```powershell
Set-Location 'E:\project\sand box game'
$env:GOCACHE='E:\project\sand box game\.gocache'
$env:GOMODCACHE='E:\project\sand box game\.cache\gomod'
go run .\cmd\server\main.go -config .\configs\local.yaml
```

### 4.2 启动前端

```powershell
Set-Location 'E:\project\sand box game\frontend'
npm run dev
```

### 4.3 常用地址

- 前端登录页：`http://127.0.0.1:5173/sandbox-game/login`
- 后端健康检查：`http://127.0.0.1:8080/healthz`

---

## 5. 常用联调账号

默认联调账号口径：

- 管理员：`admin / 123456`
- 玩家：`group01 / 123456`

说明：

- 多组演练环境下，玩家账号可能扩展为 `group01 ~ groupNN`
- 单组演练环境下，通常只保留 `admin` 与 `group01`

---

## 6. 推荐阅读顺序

如果你要理解这个项目，推荐按这个顺序看文档：

1. [requirements_spec.md](E:\project\sand box game\docs\requirements_spec.md)
2. [requirements_consensus_checklist.md](E:\project\sand box game\docs\requirements_consensus_checklist.md)
3. [technical_selection.md](E:\project\sand box game\docs\technical_selection.md)
4. [database_design.md](E:\project\sand box game\docs\database_design.md)
5. [api_design.md](E:\project\sand box game\docs\api_design.md)
6. [frontend_page_structure.md](E:\project\sand box game\docs\frontend_page_structure.md)
7. [implementation_plan.md](E:\project\sand box game\docs\implementation_plan.md)

---

## 7. 按用途查文档

### 7.1 看需求

- [requirements_spec.md](E:\project\sand box game\docs\requirements_spec.md)
- [requirements_consensus_checklist.md](E:\project\sand box game\docs\requirements_consensus_checklist.md)

### 7.2 看技术方案

- [technical_selection.md](E:\project\sand box game\docs\technical_selection.md)
- [go_project_structure.md](E:\project\sand box game\docs\go_project_structure.md)
- [backend_guide.md](E:\project\sand box game\docs\backend_guide.md)
- [frontend_guide.md](E:\project\sand box game\docs\frontend_guide.md)
- [player_operating_typography_spacing_design.md](E:\project\sand box game\docs\player_operating_typography_spacing_design.md)
- [input_navigation_design.md](E:\project\sand box game\docs\input_navigation_design.md)
- [admin_input_navigation_design.md](E:\project\sand box game\docs\admin_input_navigation_design.md)
- [keyboard_navigation_design.md](E:\project\sand box game\docs\keyboard_navigation_design.md)
- [order_multi_round_upgrade_design.md](E:\project\sand box game\docs\order_multi_round_upgrade_design.md)
- [admin_order_pool_group_filter_design.md](E:\project\sand box game\docs\admin_order_pool_group_filter_design.md)
- [order_template_airport_plan.md](E:\project\sand box game\docs\order_template_airport_plan.md)
- [airport_order_module_acceptance_checklist.md](E:\project\sand box game\docs\airport_order_module_acceptance_checklist.md)

### 7.3 看数据库与接口

- [database_design.md](E:\project\sand box game\docs\database_design.md)
- [api_design.md](E:\project\sand box game\docs\api_design.md)
- [minimal_state_machine.md](E:\project\sand box game\docs\minimal_state_machine.md)

### 7.4 看测试与联调

- [testing_guide.md](E:\project\sand box game\docs\testing_guide.md)
- [integration_acceptance_runbook.md](E:\project\sand box game\docs\integration_acceptance_runbook.md)
- [single_group_rehearsal_runbook.md](E:\project\sand box game\docs\single_group_rehearsal_runbook.md)
- [test_demo_commands.md](E:\project\sand box game\docs\test_demo_commands.md)
- [order_module_acceptance_checklist.md](E:\project\sand box game\docs\order_module_acceptance_checklist.md)
- [airport_order_module_acceptance_checklist.md](E:\project\sand box game\docs\airport_order_module_acceptance_checklist.md)

### 7.5 看比赛上线与现场使用

- [competition_launch_runbook.md](E:\project\sand box game\docs\competition_launch_runbook.md)
- [competition_deploy_checklist.md](E:\project\sand box game\docs\competition_deploy_checklist.md)
- [competition_backup_single_deploy_plan.md](E:\project\sand box game\docs\competition_backup_single_deploy_plan.md)

### 7.5.1 正式版上线包落点

- `configs/competition.yaml`
- `migrations/mysql/0004_seed_competition_admin.sql`
- `scripts/init-competition.ps1`
- `scripts/reset-competition.ps1`
- `scripts/start-competition.ps1`
- `scripts/backup-competition-database.ps1`
- `scripts/build-competition-package.ps1`
- `首次初始化数据库.bat / 启动沙盘系统.bat / 重启沙盘系统.bat / 停止沙盘系统.bat / 备份当前比赛数据库.bat / 赛前重置当前比赛.bat`
- `安装部署手册-正式版.txt / 操作手册-正式版.txt`
- `docs/competition_deploy_checklist.md`
- `docs/competition_backup_single_deploy_plan.md`

### 7.6 看 Excel 规则与对账样例

- [calculation_rule_spec.md](E:\project\sand box game\docs\calculation_rule_spec.md)
- [excel_field_mapping.md](E:\project\sand box game\docs\excel_field_mapping.md)
- [excel_reconciliation_samples.md](E:\project\sand box game\docs\excel_reconciliation_samples.md)

---

## 8. 当前重要说明

- 正式比赛规则与测试/演练环境要区分看待。
- 单组演练脚本和单组演练库是测试工具，不是正式比赛主流程入口。
- 正式产品方向已经收口为：管理员端赛前配置比赛，再录入初始基线，再开始比赛。
- 正式比赛数据库初始化不应直接使用 `0002_seed_data.sql`，而应使用正式比赛专用初始化脚本，仅保留 `admin + sg_game_config`，再由管理员首登后在 `赛前配置页` 初始化比赛。
- 当前正式版推荐采用 `Go + MySQL` 结构：Go 服务直接提供前端静态资源与业务 API，不再依赖单独的 `Nginx` 进程。
- 文档若发生冲突，以最新需求文档、共识清单与实施计划为准。

## 9. 正式版上线最短路径

推荐按以下顺序执行：

1. 如果管理员电脑已有旧比赛数据，先执行 `备份当前比赛数据库.bat`
2. 新代码解压到新部署目录，不覆盖旧目录
3. 修改 `configs/competition.yaml`
4. 执行 `scripts/build-competition-package.ps1`
5. 在比赛主机双击 `首次初始化数据库.bat`
6. 双击 `启动沙盘系统.bat`
7. 直接访问 `http://比赛主机IP:server.port/sandbox-game/login`
8. 管理员登录后先在 `赛前配置页` 选择版本包并初始化比赛，再开始正式比赛

现场如果需要逐条照着执行，请直接使用：

- [competition_deploy_checklist.md](E:\project\sand box game\docs\competition_deploy_checklist.md)

---

## 10. 文档维护约定

后续如新增页面、接口、规则或初始化能力，至少同步更新：

- `requirements_spec.md`
- `requirements_consensus_checklist.md`
- `implementation_plan.md`
- 对应专题文档

这样可以保证需求、设计、开发、测试口径一致。
