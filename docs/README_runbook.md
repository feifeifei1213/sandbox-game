# 证书迁移运行手册（README_runbook）

基线日期：`2026-03-19`

## 1. 目标

本手册用于指导在本机执行证书迁移脚本，覆盖：

- `dry-run` 演练
- `sample` 小样本实跑
- `full` 全量迁移
- 失败重跑与定位
- B 侧历史飞行记录的“航司 -> 飞机型号”导出与 Excel 回填

## 2. 前置准备

1. 准备 Python 3.10+
2. 安装依赖

```powershell
python -m pip install -e .
```

3. 复制并填写本机配置

```powershell
Copy-Item config\config.local.example.yaml config\config.local.yaml
```

必须填写：

- `a_system.cookie`
- `a_system.certificate_query_filter`
- `b_system.bearer_token`

## 3. 命令入口

统一入口：

```powershell
python run_bootstrap.py --config config/config.local.yaml --project-root .
```

可选参数：

- `--run-mode dry-run|sample|full`：覆盖配置中的运行模式
- `--limit N`：限制处理条数
- `--resume-from-a-cert-id A_CERT_ID`：从指定 A 证书 ID 继续（含该 ID）

航司机型导出入口：

```powershell
python scripts/export_airline_aircraft_models.py --config config/config.local.yaml --airport-code ZBYN --page-size 200
```

可选参数：

- `--excel-path PATH`：覆盖默认航司字典 Excel，默认 `captures/excel/航司信息(2).xlsx`
- `--sheet-name 航司信息`：覆盖 Excel sheet 名
- `--output-dir PATH`：指定输出目录；不传时自动生成 `output/airline_aircraft_models_<airport>_<timestamp>/`
- `--max-pages N`：仅调试时限制拉取页数

## 4. 推荐执行顺序

### 4.1 Dry-run

```powershell
python run_bootstrap.py --config config/config.local.yaml --project-root . --run-mode dry-run
```

### 4.2 小样本（10 条）

```powershell
python run_bootstrap.py --config config/config.local.yaml --project-root . --run-mode sample --limit 10
```

### 4.3 全量

```powershell
python run_bootstrap.py --config config/config.local.yaml --project-root . --run-mode full
```

## 5. 产物说明

每次运行会在 `output/<run_id>/` 生成：

- `success.csv`
- `failed.csv`
- `skipped.csv`
- `updated.csv`
- `unmatched_person.csv`
- `run_summary.json`

说明：

- `skipped.csv` 中 `dry_run_dedup_exact_match` / `dedup_exact_match` 表示 B 侧证书主记录已完全一致，不会重复写主记录。
- `full` 模式下即使主记录命中 `skip`，仍会继续执行人员与机型绑定。

日志位于：

- `logs/<run_id>/migration.log`

航司机型导出会在 `output/airline_aircraft_models_<airport_code>_<timestamp>/` 生成：

- `summary.json`
- `airlines.json`
- `airlines.csv`
- `unmatched_flights.csv`
- `航司信息(2)_已回填机型.xlsx`

说明：

- `airlines.csv` / `airlines.json` 中会同时保留 `Excel原飞机型号` 与 `Excel回填飞机型号`
- `航司信息(2)_已回填机型.xlsx` 是副本，不会覆盖原始 `captures/excel/航司信息(2).xlsx`
- 回填值优先使用历史飞行记录聚合出的 `normalized_models`，当前分隔符为 ` / `
- 航司匹配优先级为：`airlineCompany(三字码)` -> `portIn/portOut 航班号前缀(二字码)`

## 6. 常见重跑方式

### 6.1 从失败证书继续

```powershell
python run_bootstrap.py --config config/config.local.yaml --project-root . --run-mode full --resume-from-a-cert-id <A_CERT_ID>
```

### 6.2 仅验证修复后的少量数据

```powershell
python run_bootstrap.py --config config/config.local.yaml --project-root . --run-mode sample --limit 5 --resume-from-a-cert-id <A_CERT_ID>
```

## 7. 故障排查

1. `B-side bearer_token still contains a placeholder`
- 原因：`config.local.yaml` 仍是占位符
- 处理：填真实 token 后重跑

2. `resume_from_a_cert_id not found`
- 原因：指定的 A 证书 ID 不在当前读取结果中
- 处理：确认筛选条件 `certificate_query_filter` 与 ID 是否一致

3. `failed.csv` 中 `MISSING_REMINDER_DATE`
- 原因：当前运行配置将 `execution.require_reminder_date=true` 时，A 侧提醒日期为空或格式非法
- 处理：若业务允许空值，则将 `execution.require_reminder_date=false`，脚本会省略 `reminderDate` 字段；若需严格校验，则先修正 A 侧源数据后重跑

4. `failed.csv` 中 `PERSON_NOT_FOUND_IN_B` / `PERSON_MULTI_MATCH`
- 原因：人员未命中、命名不一致或重名歧义
- 处理：优先确认是否存在已确认别名（如 `石磊（小） -> 石磊17`、`高鹏19 -> 高鹏1158`）；若仍不命中，再修正 B 人员档案或补齐部门辅助信息后重跑

5. `skipped.csv` 中 `dry_run_dedup_exact_match` / `dedup_exact_match`
- 原因：B 侧已存在同键同值证书主记录
- 处理：这是正常结果，不需要重复写主记录；只需确认 `b_cert_id` 是否符合预期

6. 航司机型导出结果中出现 `AMBIGUOUS_*`
- 原因：航司字典中同一二字码或三字码对应了多条记录，例如 `CF`、`CYZ`
- 处理：先检查 `summary.json` 中 `duplicate_iata_codes` / `duplicate_icao_codes`，必要时清洗 Excel 航司字典后重跑

7. 航司机型导出结果中出现 `CODE_NOT_FOUND`
- 原因：历史飞行记录里的航司代码未在 Excel 航司字典中维护，例如部分外航或特殊航班号
- 处理：先查看 `unmatched_flights.csv`，补齐字典后重跑

## 8. 安全与审计

- 不在代码仓库提交真实 Cookie/Token
- 失败快照会做敏感字段脱敏（如 `cookie` / `authorization` / `token`）
- 对外共享运行结果前，先复核 `failed.csv` 与日志内容



