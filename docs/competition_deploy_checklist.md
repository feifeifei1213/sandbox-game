# 沙盘经营系统现场一键部署清单

> 更新日期：2026-08-03
> 适用对象：比赛现场技术支持、部署执行人、主持人  
> 文档定位：本文件不是解释“为什么这样部署”，而是给现场直接照着执行的清单。若要看完整背景和边界，请同时参考 `docs/competition_launch_runbook.md`。

## 1. 使用场景

本清单适用于以下目标：

- 在一台 `Windows` 比赛主机上部署首版正式比赛系统
- 使用 `Go + MySQL` 提供统一访问地址，由 Go 服务同时提供前端静态页面与业务 API
- 正式比赛库只保留 `admin + sg_game_config` 初始化结果
- 由管理员首次登录后，在系统内完成“赛前配置页”初始化小组数量
- 场景 B：管理员电脑已经跑过一次比赛时，先备份旧比赛数据库，再把新代码部署到新文件夹
- 新部署推荐采用单实例入口，生产制造版 / 贵宾服务版在网页赛前配置页选择，不再通过不同启动脚本区分

本清单默认你已经拿到了当前仓库。

如果要在开发机现场打包，需要：

- `PowerShell`
- `Go`
- `Node.js / npm`
- `MySQL`

如果比赛主机直接使用已经构建好的正式版上线包，则比赛主机只需要：

- `PowerShell`
- `MySQL`

---

## 2. 部署目标结果

部署完成后，现场应达到以下状态：

- 玩家和管理员统一访问一个地址，例如 `http://192.168.8.200:8080/sandbox-game/login`
- 前端页面与 `/api/` 都由同一个 Go 服务直接提供
- 正式比赛库已初始化，但尚未预建玩家组
- 管理员首次登录后进入 `赛前配置页`
- 管理员在页面中配置小组数量并点击“初始化比赛”后，才真正生成 `group01 ~ groupNN`
- 若是已有旧比赛数据的管理员电脑，旧代码目录和旧数据库备份文件已经保留，新比赛使用新部署目录和新数据库
- 管理员日常只需要双击 `启动沙盘系统.bat`，版本差异由系统内版本包处理

---

## 3. 赛前准备清单

比赛前一天，先逐项确认：

- 如果管理员电脑已有旧比赛数据，已先执行 `备份当前比赛数据库.bat`，并确认 `backup` 文件夹下生成 zip
- 新代码部署目录不会覆盖旧代码目录
- 比赛主机固定 IP 已确认，例如 `192.168.8.200`
- 比赛主机和所有参赛电脑在同一局域网
- MySQL 已安装并可用
- 如果要在比赛主机本机打包，主机还需安装 `Go` 和 `Node.js`
- 防火墙已放行对外访问端口，建议与 `competition.yaml` 中 `server.port` 保持一致
- 已准备一个正式比赛数据库，例如 `sandbox_game_competition`
- 已准备数据库账号，具备该库的建表和写入权限
- 已明确正式比赛不复用开发库、不复用演练库
- 已准备管理员账号发放口径，默认可用 `admin / 123456`

---

## 4. 第一步：修改正式版配置

在仓库根目录修改 `configs/competition.yaml`。

至少要改这两项：

```yaml
mysql:
  dsn: root:你的数据库密码@tcp(127.0.0.1:3306)/sandbox_game_competition?charset=utf8mb4&parseTime=True&loc=Local

auth:
  tokenSecret: 替换成一串随机密钥
```

检查项：

- `mysql.dsn` 指向正式比赛库
- `tokenSecret` 不再保留 `CHANGE_ME_TO_A_RANDOM_SECRET`
- `server.port` 就是正式对外访问端口，例如 `8080`
- `frontend.distDir` 保持 `frontend/dist` 即可，由后端直接读取打包后的前端静态文件

---

## 5. 第二步：构建正式版上线包

在开发机或比赛主机执行：

```powershell
Set-Location 'E:\project\sand box game'
.\scripts\build-competition-package.ps1 -ConfigPath .\configs\competition.yaml
```

执行完成后，重点查看两个产物：

- `.\.runtime\competition-package\sandbox-game-competition\`
- `.\.runtime\competition-package\sandbox-game-competition.zip`

上线包内应至少包含：

- `bin\sandbox-game-server.exe`
- `bin\dbtool.exe`
- `frontend\dist\`
- `configs\competition.yaml`
- `scripts\init-competition.ps1`
- `scripts\reset-competition.ps1`
- `scripts\start-competition.ps1`
- `scripts\backup-competition-database.ps1`
- `migrations\mysql\0017_order_delivery_revenue.sql`
- `首次初始化数据库.bat`
- `启动沙盘系统.bat`
- `重启沙盘系统.bat`
- `停止沙盘系统.bat`
- `备份当前比赛数据库.bat`
- `赛前重置当前比赛.bat`
- `安装部署手册-正式版.txt`
- `操作手册-正式版.txt`
- `docs\competition_launch_runbook.md`
- `docs\competition_deploy_checklist.md`
- `docs\competition_backup_single_deploy_plan.md`

如果比赛主机不能直接编译，也可以把整个 `sandbox-game-competition` 目录复制到比赛主机。

---

## 5.1 场景 B：已有旧比赛数据时的额外动作

如果管理员电脑已经做过一次比赛，先不要覆盖旧目录，也不要直接重置旧数据库。

推荐顺序：

1. 进入旧系统目录。
2. 双击 `备份当前比赛数据库.bat`。
3. 按提示输入备份文件名，例如 `第一场正式比赛数据`。
4. 确认旧系统目录下生成 `backup\第一场正式比赛数据.zip`。
5. 新建新的部署目录，例如 `E:\deploy\sandbox-game-competition`。
6. 将最新上线包解压到新目录。
7. 新比赛使用新的数据库，例如 `sandbox_game_competition`。

备份 zip 第一版应包含：

- 完整数据库 SQL；
- 结果摘要 CSV；
- 备份说明 TXT。

结果摘要用于快速查看每组排名、收入、利润、权益、CEO / CFO 等数据；严肃追溯仍以 SQL 恢复到临时数据库后的查询为准。

---

## 6. 第三步：初始化正式比赛库

在比赛主机进入上线包根目录，例如：

```powershell
Set-Location 'E:\deploy\sandbox-game-competition'
.\scripts\init-competition.ps1 -ConfigPath .\configs\competition.yaml
```

如果使用正式包根目录中的管理员双击入口，也可以直接双击：

```text
首次初始化数据库.bat
```

这一步只会初始化：

- 表结构与必要业务表
- `admin` 账号
- `sg_game_config`

这一步不会直接创建玩家组和玩家账号。这个行为是正确的，不是漏了。

执行后建议立刻确认：

- 数据库中已经有 `admin`
- `sg_group` 仍然是空的
- `sg_game_config` 已创建

---

## 7. 第四步：启动正式后端

在比赛主机进入上线包根目录执行：

```powershell
Set-Location 'E:\deploy\sandbox-game-competition'
.\scripts\start-competition.ps1 -ConfigPath .\configs\competition.yaml
```

如果使用正式包根目录中的管理员双击入口，也可以直接双击：

```text
启动沙盘系统.bat
```

后端默认监听：

- `0.0.0.0:8080`

可先本机检查：

```powershell
curl.exe -s -i "http://127.0.0.1:8080/healthz"
```

预期结果：

- 返回 `HTTP 200`

---

## 8. 第五步：现场联通验证

当前正式版不再需要单独配置 `Nginx`。

前端登录页与 `/api/` 接口都由 `sandbox-game-server.exe` 直接提供，因此只要后端服务启动成功，就可以直接访问：

- `http://127.0.0.1:8080/sandbox-game/login`

---

## 9. 第六步：现场联通验证补充

先在比赛主机浏览器打开：

- `http://127.0.0.1:8080/sandbox-game/login`

再在管理员电脑和至少一台玩家电脑打开：

- `http://比赛主机IP:8080/sandbox-game/login`

必须逐项确认：

1. 登录页可以打开
2. 管理员 `admin / 123456` 可以登录
3. 管理员首次登录后进入 `赛前配置页`
4. 管理员能填写小组数量，例如 `6`
5. 点击初始化比赛后，管理员进入正式管理页
6. 玩家账号 `group01 / 123456` 可以登录
7. 玩家默认进入 `0年经营页`
8. 玩家端页面刷新后仍可正常访问
9. 管理员端刷新后仍可正常访问
10. 赛前配置页中可选择本场比赛版本包，例如生产制造版或贵宾服务版

---

## 10. 比赛开始前最后动作

正式给玩家发网址前，按这个顺序做：

1. 管理员完成赛前初始化
2. 主持人确认小组数量正确
3. 技术支持确认 `http://比赛主机IP:8080/sandbox-game/login` 可从多台电脑访问
4. 导出或截图保存一份当前数据库备份记录
5. 再统一向玩家发放访问地址和账号

不要发这些地址：

- `http://127.0.0.1:5173`
- `http://localhost:5173`
- 开发者电脑临时 IP
- 健康检查地址或其他非登录页入口

---

## 11. 常见故障排查

### 11.1 登录页打不开

按顺序检查：

1. `sandbox-game-server.exe` 是否已启动
2. `frontend/dist` 是否已随正式包一并拷贝
3. `configs/competition.yaml` 中 `frontend.distDir` 是否仍指向 `frontend/dist`
4. 比赛主机防火墙是否放行 `server.port`
5. 玩家电脑和比赛主机是否在同一局域网

### 11.2 页面能打开，但登录失败或接口报错

按顺序检查：

1. 后端服务是否仍在运行
2. `http://127.0.0.1:8080/healthz` 是否正常
3. `configs/competition.yaml` 中数据库连接是否正确
4. `configs/competition.yaml` 中 `server.port` 是否与当前访问地址一致

### 11.3 管理员登录后没有进入赛前配置页

按顺序检查：

1. 正式比赛库是否错误地导入了旧的 10 组种子数据
2. `sg_group` 表中是否已经有记录
3. 是否错误执行了非正式比赛初始化脚本

说明：

- 正式比赛环境如果一开始就有 `sg_group`，系统会认为比赛已初始化
- 正式比赛库不能直接使用 `0002_seed_data.sql` 一类预建玩家组脚本

### 11.4 需要重新开赛前初始化

若比赛尚未正式开始，且确认要整库重置，可执行：

```powershell
Set-Location 'E:\deploy\sandbox-game-competition'
.\scripts\reset-competition.ps1 -ConfigPath .\configs\competition.yaml
.\scripts\init-competition.ps1 -ConfigPath .\configs\competition.yaml
```

注意：

- 这是重置正式比赛库的动作
- 只有在主持人与技术支持共同确认后才能执行

---

## 12. 最短执行路径

如果现场时间非常紧，最短路径就是这 `6` 步：

1. 改 `configs/competition.yaml`
2. 执行 `.\scripts\build-competition-package.ps1`
3. 把上线包复制到比赛主机
4. 执行 `.\scripts\init-competition.ps1`
5. 执行 `.\scripts\start-competition.ps1`
6. 管理员登录后先完成 `赛前配置页` 初始化，再给玩家发网址

---

## 13. 角色分工建议

- 技术支持：改配置、打包、启动后端、做连通性验证
- 主持人 / 管理员：登录系统、配置小组数量、初始化比赛、发账号
- 参赛玩家：只访问统一网址，不直接接触数据库、脚本和后端端口

---

## 14. 相关文档

- `docs/competition_launch_runbook.md`
- `docs/README.md`
