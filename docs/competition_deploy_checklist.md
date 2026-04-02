# 沙盘经营系统现场一键部署清单

> 更新日期：2026-04-02  
> 适用对象：比赛现场技术支持、部署执行人、主持人  
> 文档定位：本文件不是解释“为什么这样部署”，而是给现场直接照着执行的清单。若要看完整背景和边界，请同时参考 `docs/competition_launch_runbook.md`。

## 1. 使用场景

本清单适用于以下目标：

- 在一台 `Windows` 比赛主机上部署首版正式比赛系统
- 使用 `Nginx + Go + MySQL` 提供统一访问地址
- 正式比赛库只保留 `admin + sg_game_config` 初始化结果
- 由管理员首次登录后，在系统内完成“赛前配置页”初始化小组数量

本清单默认你已经拿到了当前仓库，且比赛主机可以运行：

- `PowerShell`
- `Go`
- `Node.js / npm`
- `MySQL`
- `Nginx`

---

## 2. 部署目标结果

部署完成后，现场应达到以下状态：

- 玩家和管理员统一访问一个地址，例如 `http://192.168.8.200`
- 前端页面由 `Nginx` 提供
- `/api/` 请求由 `Nginx` 反向代理到本机 `127.0.0.1:8080`
- 正式比赛库已初始化，但尚未预建玩家组
- 管理员首次登录后进入 `赛前配置页`
- 管理员在页面中配置小组数量并点击“初始化比赛”后，才真正生成 `group01 ~ groupNN`

---

## 3. 赛前准备清单

比赛前一天，先逐项确认：

- 比赛主机固定 IP 已确认，例如 `192.168.8.200`
- 比赛主机和所有参赛电脑在同一局域网
- MySQL 已安装并可用
- Nginx 已安装并可用
- 比赛主机已安装 `Go` 和 `Node.js`
- 防火墙已放行对外访问端口，建议 `80`
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
- `server.port` 保持 `8080` 即可，供 `Nginx` 反向代理

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
- `nginx\sandbox-game.competition.conf`
- `docs\competition_launch_runbook.md`
- `docs\competition_deploy_checklist.md`

如果比赛主机不能直接编译，也可以把整个 `sandbox-game-competition` 目录复制到比赛主机。

---

## 6. 第三步：初始化正式比赛库

在比赛主机进入上线包根目录，例如：

```powershell
Set-Location 'E:\deploy\sandbox-game-competition'
.\scripts\init-competition.ps1 -ConfigPath .\configs\competition.yaml
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

后端默认监听：

- `0.0.0.0:8080`

可先本机检查：

```powershell
curl.exe -s -i "http://127.0.0.1:8080/healthz"
```

预期结果：

- 返回 `HTTP 200`

---

## 8. 第五步：配置 Nginx

参考样例文件 `scripts/nginx/sandbox-game.competition.conf`。

你至少要改一项：

- 把 `root` 改成比赛主机上 `frontend/dist` 的真实绝对路径

示例：

```nginx
server {
    listen       80;
    server_name  _;

    root   E:/deploy/sandbox-game-competition/frontend/dist;
    index  index.html;

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
    }

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

应用配置后，重载 `Nginx`。

如果你使用的是 Windows 下的 `Nginx`，常见命令形态如下：

```powershell
cd 'E:\nginx'
.\nginx.exe -s reload
```

---

## 9. 第六步：现场联通验证

先在比赛主机浏览器打开：

- `http://127.0.0.1`

再在管理员电脑和至少一台玩家电脑打开：

- `http://比赛主机IP`

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

---

## 10. 比赛开始前最后动作

正式给玩家发网址前，按这个顺序做：

1. 管理员完成赛前初始化
2. 主持人确认小组数量正确
3. 技术支持确认 `http://比赛主机IP` 可从多台电脑访问
4. 导出或截图保存一份当前数据库备份记录
5. 再统一向玩家发放访问地址和账号

不要发这些地址：

- `http://127.0.0.1:5173`
- `http://localhost:5173`
- 开发者电脑临时 IP
- 后端裸地址 `http://比赛主机IP:8080`

---

## 11. 常见故障排查

### 11.1 登录页打不开

按顺序检查：

1. `Nginx` 是否已启动
2. `Nginx` 配置中的 `root` 是否指向正确的 `frontend/dist`
3. 比赛主机防火墙是否放行 `80`
4. 玩家电脑和比赛主机是否在同一局域网

### 11.2 页面能打开，但登录失败或接口报错

按顺序检查：

1. 后端服务是否仍在运行
2. `http://127.0.0.1:8080/healthz` 是否正常
3. `configs/competition.yaml` 中数据库连接是否正确
4. `Nginx` 的 `/api/` 是否确实代理到 `127.0.0.1:8080`

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
5. 执行 `.\scripts\start-competition.ps1` 并配置 `Nginx`
6. 管理员登录后先完成 `赛前配置页` 初始化，再给玩家发网址

---

## 13. 角色分工建议

- 技术支持：改配置、打包、启动后端、配 `Nginx`、做连通性验证
- 主持人 / 管理员：登录系统、配置小组数量、初始化比赛、发账号
- 参赛玩家：只访问统一网址，不直接接触数据库、脚本和后端端口

---

## 14. 相关文档

- `docs/competition_launch_runbook.md`
- `docs/README.md`
