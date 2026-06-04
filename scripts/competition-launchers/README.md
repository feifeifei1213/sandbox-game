# 正式版管理员启动脚本目录

本目录用于存放正式版现场启动脚本模板，默认面向以下部署结构：

- `D:\deploy\sandbox-game-shengchan`
- `D:\deploy\sandbox-game-guibing`
- 本机 MySQL 服务默认自动识别常见名称：`SandboxGameMySQL / MySQL / MySQL80 / MySQL57 / MariaDB`

目录内文件分为两类：

- 管理员直接双击使用的 `.bat`
- 给管理员阅读的 [安装部署手册-正式版.txt](E:\project\sand box game\scripts\competition-launchers\安装部署手册-正式版.txt) 与 [操作手册-正式版.txt](E:\project\sand box game\scripts\competition-launchers\操作手册-正式版.txt)

当前脚本职责边界：

- `启动制造业版.bat`：拉起制造业版整套环境
- `启动服务版.bat`：拉起服务版整套环境
- `首次初始化制造业版数据库.bat`：首次部署制造业版时初始化数据库
- `首次初始化服务版数据库.bat`：首次部署服务版时初始化数据库
- `重启制造业版服务.bat`：只重启制造业版，不清数据
- `重启服务版服务.bat`：只重启服务版，不清数据
- `赛前重置制造业版数据库.bat`：仅赛前可用，清空并重新初始化制造业版数据库
- `赛前重置服务版数据库.bat`：仅赛前可用，清空并重新初始化服务版数据库
- `停止全部服务.bat`：统一停止两套游戏前后端进程，不停止 MySQL
- `更新制造业版前端.bat`：制造业版前端覆盖更新脚本模板，需和同一个更新包里的 `frontend\dist\` 一起使用
- `更新服务版前端.bat`：服务版前端覆盖更新脚本模板，需和同一个更新包里的 `frontend\dist\` 一起使用

说明：

- 本目录当前是仓库内模板目录，后续交付时建议整体复制到 `D:\deploy\launcher\`
- 脚本默认兼容当前“后端直接提供前端页面”的正式包结构，因此只依赖各自包目录下的 `configs\competition.yaml`、`scripts\start-competition.ps1` 与 `frontend\dist\`
- 当前模板默认约定制造业版与服务版分别监听 `18080 / 28080`；对应正式包中的 `competition.yaml` 也需要同步配置成相同端口
- 当前目录中的两份手册是用于直接发给管理员的成品文案：安装部署手册负责“自己把环境装起来”，操作手册负责“比赛当天怎么用”
- 当前脚本会优先自动识别常见 MySQL 服务名；若现场机器使用自定义服务名，再手工修改这些 `.bat` 顶部的 `MYSQL_SERVICE` 变量
- 推荐赛前部署顺序：创建数据库 -> 双击“首次初始化数据库” -> 双击“启动对应版本”
- `更新制造业版前端.bat` 与 `更新服务版前端.bat` 不建议长期单独放到 `D:\deploy\launcher\` 里裸用；更适合和新的 `frontend\dist\` 一起打成单独“前端更新包”后再发给管理员
- 这两套更新脚本会自动在 `D:\deploy\_backup\shengchan` / `D:\deploy\_backup\guibing` 下创建带时间戳的前端备份，再覆盖正式目录中的 `frontend\dist`
