param(
    [string]$ConfigPath = '.\configs\competition.yaml',
    [string]$OutputDir = '',
    [string]$BackupName = ''
)

$ErrorActionPreference = 'Stop'

function Resolve-WorkspacePath {
    param(
        [string]$BasePath,
        [string]$InputPath
    )
    if ([string]::IsNullOrWhiteSpace($InputPath)) {
        return $InputPath
    }
    if ([System.IO.Path]::IsPathRooted($InputPath)) {
        return $InputPath
    }
    return Join-Path $BasePath $InputPath
}

function Read-CompetitionDsn {
    param([string]$Path)
    if (-not (Test-Path -LiteralPath $Path)) {
        throw "未找到配置文件：$Path"
    }

    $inMysqlSection = $false
    foreach ($line in Get-Content -LiteralPath $Path -Encoding UTF8) {
        if ($line -match '^\s*mysql\s*:\s*$') {
            $inMysqlSection = $true
            continue
        }
        if ($inMysqlSection -and $line -match '^\S') {
            $inMysqlSection = $false
        }
        if ($inMysqlSection -and $line -match '^\s*dsn\s*:\s*(.+?)\s*$') {
            $dsn = $matches[1].Trim()
            if (($dsn.StartsWith('"') -and $dsn.EndsWith('"')) -or ($dsn.StartsWith("'") -and $dsn.EndsWith("'"))) {
                $dsn = $dsn.Substring(1, $dsn.Length - 2)
            }
            if ([string]::IsNullOrWhiteSpace($dsn)) {
                throw "配置文件中的 mysql.dsn 为空：$Path"
            }
            return $dsn
        }
    }

    throw "配置文件中未找到 mysql.dsn：$Path"
}

function Parse-MySqlDsn {
    param([string]$Dsn)

    $marker = '@tcp('
    $markerIndex = $Dsn.IndexOf($marker)
    if ($markerIndex -lt 0) {
        throw "暂不支持的 MySQL DSN 格式：未找到 @tcp(...)"
    }

    $credentialPart = $Dsn.Substring(0, $markerIndex)
    $endpointStart = $markerIndex + $marker.Length
    $endpointEnd = $Dsn.IndexOf(')', $endpointStart)
    if ($endpointEnd -lt 0) {
        throw "暂不支持的 MySQL DSN 格式：tcp(...) 未闭合"
    }

    $endpoint = $Dsn.Substring($endpointStart, $endpointEnd - $endpointStart)
    $databasePart = $Dsn.Substring($endpointEnd + 1)
    if (-not $databasePart.StartsWith('/')) {
        throw "暂不支持的 MySQL DSN 格式：未找到数据库名"
    }

    $databasePart = $databasePart.Substring(1)
    $queryIndex = $databasePart.IndexOf('?')
    $database = if ($queryIndex -ge 0) { $databasePart.Substring(0, $queryIndex) } else { $databasePart }
    if ([string]::IsNullOrWhiteSpace($database)) {
        throw "MySQL DSN 中数据库名为空"
    }

    $credentialSeparator = $credentialPart.IndexOf(':')
    if ($credentialSeparator -ge 0) {
        $user = $credentialPart.Substring(0, $credentialSeparator)
        $password = $credentialPart.Substring($credentialSeparator + 1)
    }
    else {
        $user = $credentialPart
        $password = ''
    }
    if ([string]::IsNullOrWhiteSpace($user)) {
        throw "MySQL DSN 中用户名为空"
    }

    $mysqlHost = $endpoint
    $port = '3306'
    if ($endpoint.StartsWith('[')) {
        $closingBracket = $endpoint.IndexOf(']')
        if ($closingBracket -gt 0) {
            $mysqlHost = $endpoint.Substring(1, $closingBracket - 1)
            if ($endpoint.Length -gt ($closingBracket + 2) -and $endpoint.Substring($closingBracket + 1, 1) -eq ':') {
                $port = $endpoint.Substring($closingBracket + 2)
            }
        }
    }
    else {
        $lastColon = $endpoint.LastIndexOf(':')
        if ($lastColon -gt 0) {
            $mysqlHost = $endpoint.Substring(0, $lastColon)
            $port = $endpoint.Substring($lastColon + 1)
        }
    }

    [pscustomobject]@{
        User     = $user
        Password = $password
        Host     = $mysqlHost
        Port     = $port
        Database = $database
    }
}

function Find-MySqlTool {
    param([string]$ToolName)

    $command = Get-Command $ToolName -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($command) {
        return $command.Source
    }

    $candidatePaths = New-Object System.Collections.Generic.List[string]
    $programRoots = @($env:ProgramFiles, ${env:ProgramFiles(x86)}, 'C:\Program Files', 'C:\Program Files (x86)') |
        Where-Object { -not [string]::IsNullOrWhiteSpace($_) } |
        Select-Object -Unique

    foreach ($root in $programRoots) {
        $candidatePaths.Add((Join-Path $root "MySQL\MySQL Server 8.4\bin\$ToolName"))
        $candidatePaths.Add((Join-Path $root "MySQL\MySQL Server 8.0\bin\$ToolName"))
        $candidatePaths.Add((Join-Path $root "MySQL\MySQL Server 5.7\bin\$ToolName"))
        $candidatePaths.Add((Join-Path $root "MariaDB 11.0\bin\$ToolName"))
        $candidatePaths.Add((Join-Path $root "MariaDB 10.11\bin\$ToolName"))
        $candidatePaths.Add((Join-Path $root "MariaDB 10.6\bin\$ToolName"))
        $mysqlRoot = Join-Path $root 'MySQL'
        if (Test-Path -LiteralPath $mysqlRoot) {
            Get-ChildItem -LiteralPath $mysqlRoot -Directory -Filter 'MySQL Server*' -ErrorAction SilentlyContinue |
                ForEach-Object { $candidatePaths.Add((Join-Path $_.FullName "bin\$ToolName")) }
        }
        Get-ChildItem -LiteralPath $root -Directory -Filter 'MariaDB*' -ErrorAction SilentlyContinue |
            ForEach-Object { $candidatePaths.Add((Join-Path $_.FullName "bin\$ToolName")) }
    }

    foreach ($candidate in ($candidatePaths | Select-Object -Unique)) {
        if (Test-Path -LiteralPath $candidate) {
            return $candidate
        }
    }

    return $null
}

function Get-SafeBackupBaseName {
    param([string]$Name)

    $safe = $Name.Trim()
    if ($safe.EndsWith('.zip', [System.StringComparison]::OrdinalIgnoreCase)) {
        $safe = $safe.Substring(0, $safe.Length - 4)
    }
    foreach ($char in [System.IO.Path]::GetInvalidFileNameChars()) {
        $safe = $safe.Replace([string]$char, '_')
    }
    $safe = ($safe -replace '\s+', ' ').Trim(' ', '.')
    if ([string]::IsNullOrWhiteSpace($safe)) {
        $safe = '比赛数据备份'
    }
    return $safe
}

function Get-UniqueZipPath {
    param(
        [string]$Directory,
        [string]$BaseName
    )

    $zipPath = Join-Path $Directory ($BaseName + '.zip')
    if (-not (Test-Path -LiteralPath $zipPath)) {
        return $zipPath
    }

    $timestamp = Get-Date -Format 'yyyyMMdd_HHmmss'
    $zipPath = Join-Path $Directory ("{0}_{1}.zip" -f $BaseName, $timestamp)
    if (-not (Test-Path -LiteralPath $zipPath)) {
        return $zipPath
    }

    $suffix = [System.Guid]::NewGuid().ToString('N').Substring(0, 8)
    return (Join-Path $Directory ("{0}_{1}_{2}.zip" -f $BaseName, $timestamp, $suffix))
}

function New-DefaultsExtraFile {
    param(
        [object]$Connection,
        [string]$Directory
    )

    $path = Join-Path $Directory ('mysql-client-' + [System.Guid]::NewGuid().ToString('N') + '.cnf')
    $lines = @(
        '[client]',
        "user=$($Connection.User)",
        "password=$($Connection.Password)",
        "host=$($Connection.Host)",
        "port=$($Connection.Port)",
        'default-character-set=utf8mb4'
    )
    $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
    [System.IO.File]::WriteAllLines($path, $lines, $utf8NoBom)
    return $path
}

function Write-Utf8BomLines {
    param(
        [string]$Path,
        [string[]]$Lines
    )
    $utf8Bom = New-Object System.Text.UTF8Encoding($true)
    [System.IO.File]::WriteAllLines($Path, $Lines, $utf8Bom)
}

function Invoke-DatabaseDump {
    param(
        [string]$DumpExe,
        [string]$DefaultsFile,
        [string]$Database,
        [string]$SqlPath
    )

    $dumpArgs = @(
        "--defaults-extra-file=$DefaultsFile",
        '--single-transaction',
        '--routines',
        '--triggers',
        '--events',
        '--default-character-set=utf8mb4',
        $Database
    )

    & $DumpExe @dumpArgs 1> $SqlPath
    if ($LASTEXITCODE -ne 0) {
        throw "mysqldump 执行失败，退出码：$LASTEXITCODE"
    }
    if (-not (Test-Path -LiteralPath $SqlPath)) {
        throw "mysqldump 未生成 SQL 文件：$SqlPath"
    }
    $sqlItem = Get-Item -LiteralPath $SqlPath
    if ($sqlItem.Length -le 0) {
        throw "mysqldump 生成的 SQL 文件为空：$SqlPath"
    }
}

function Convert-TsvLinesToCsvFile {
    param(
        [object[]]$TsvLines,
        [string]$CsvPath
    )

    if (-not $TsvLines -or $TsvLines.Count -eq 0) {
        throw "摘要查询没有返回任何内容"
    }

    $textLines = $TsvLines | ForEach-Object { [string]$_ }
    if ($textLines.Count -eq 1) {
        $headerLine = (($textLines[0] -split "`t") | ForEach-Object {
            '"' + ([string]$_).Replace('"', '""') + '"'
        }) -join ','
        Write-Utf8BomLines -Path $CsvPath -Lines @($headerLine)
        return
    }

    $rows = $textLines | ConvertFrom-Csv -Delimiter "`t"
    $csvLines = $rows | ConvertTo-Csv -NoTypeInformation
    Write-Utf8BomLines -Path $CsvPath -Lines $csvLines
}

function New-ResultSummaryCsv {
    param(
        [string]$MysqlExe,
        [string]$DefaultsFile,
        [string]$Database,
        [string]$CsvPath
    )

    $summarySql = @"
SET @final_year := COALESCE((SELECT final_year FROM sg_game_config ORDER BY id LIMIT 1), 0);
SET @rank := 0;
SELECT
  g.group_no AS `小组编号`,
  g.group_name AS `小组名称`,
  g.business_status AS `当前经营状态`,
  COALESCE(fs.year_no, ls.year_no, '') AS `最终有效年份`,
  COALESCE(fs.revenue, ls.revenue, '') AS `销售收入`,
  COALESCE(fs.profit, ls.profit, '') AS `净利润`,
  COALESCE(fs.equity, ls.equity, '') AS `所有者权益`,
  COALESCE(CAST(r.rank_no AS CHAR), '') AS `排名`,
  COALESCE(JSON_UNQUOTE(JSON_EXTRACT(COALESCE(frs.report_computed_snapshot_json, lrs.report_computed_snapshot_json, cr.report_computed_payload_json), '$.reportBestCeoScore')), '') AS `最佳CEO得分`,
  COALESCE(JSON_UNQUOTE(JSON_EXTRACT(COALESCE(frs.report_computed_snapshot_json, lrs.report_computed_snapshot_json, cr.report_computed_payload_json), '$.reportBestCfoScore')), '') AS `最佳CFO得分`,
  COALESCE(JSON_UNQUOTE(JSON_EXTRACT(COALESCE(frs.report_computed_snapshot_json, lrs.report_computed_snapshot_json, cr.report_computed_payload_json), '$.reportBestSalesDirectorScore')), '') AS `最佳销售总监得分`,
  COALESCE(JSON_UNQUOTE(JSON_EXTRACT(COALESCE(frs.report_computed_snapshot_json, lrs.report_computed_snapshot_json, cr.report_computed_payload_json), '$.reportBestMarketDirectorScore')), '') AS `最佳市场经营总监得分`,
  COALESCE(JSON_UNQUOTE(JSON_EXTRACT(COALESCE(frs.report_computed_snapshot_json, lrs.report_computed_snapshot_json, cr.report_computed_payload_json), '$.reportBestTechnologyDirectorScore')), '') AS `最佳科技创新总监得分`,
  COALESCE(JSON_UNQUOTE(JSON_EXTRACT(COALESCE(frs.report_manual_snapshot_json, lrs.report_manual_snapshot_json, cr.report_manual_payload_json), '$.productionHumanScore')), '') AS `最佳生产或服务人力总监得分`,
  COALESCE(JSON_UNQUOTE(JSON_EXTRACT(COALESCE(frs.report_manual_snapshot_json, lrs.report_manual_snapshot_json, cr.report_manual_payload_json), '$.enterpriseCertificationScore')), '') AS `企业认证得分`,
  CASE
    WHEN fs.id IS NOT NULL THEN '最终年有效汇总'
    WHEN ls.id IS NOT NULL THEN '未找到最终年有效汇总，使用最后有效汇总'
    WHEN cr.id IS NOT NULL THEN '未找到有效汇总，使用当前财报记录'
    ELSE '未找到有效汇总或财报'
  END AS `摘要说明`
FROM sg_group g
LEFT JOIN sg_group_summary_snapshot fs
  ON fs.group_id = g.id
 AND fs.year_no = @final_year
 AND fs.summary_effective = 1
LEFT JOIN (
  SELECT ss.*
  FROM sg_group_summary_snapshot ss
  JOIN (
    SELECT group_id, MAX(year_no) AS year_no
    FROM sg_group_summary_snapshot
    WHERE summary_effective = 1
    GROUP BY group_id
  ) latest
    ON latest.group_id = ss.group_id
   AND latest.year_no = ss.year_no
  WHERE ss.summary_effective = 1
) ls
  ON ls.group_id = g.id
LEFT JOIN (
  SELECT ranked.group_id, ranked.rank_no
  FROM (
    SELECT ordered.group_id, (@rank := @rank + 1) AS rank_no
    FROM (
      SELECT s2.group_id
      FROM sg_group_summary_snapshot s2
      JOIN sg_group g2 ON g2.id = s2.group_id
      WHERE s2.year_no = @final_year
        AND s2.summary_effective = 1
      ORDER BY s2.ranking_value DESC, g2.group_no ASC
    ) ordered
  ) ranked
) r
  ON r.group_id = g.id
LEFT JOIN sg_group_report_submission frs
  ON frs.group_id = g.id
 AND frs.year_no = fs.year_no
 AND frs.submit_version = fs.source_report_submit_version
LEFT JOIN sg_group_report_submission lrs
  ON lrs.group_id = g.id
 AND lrs.year_no = ls.year_no
 AND lrs.submit_version = ls.source_report_submit_version
LEFT JOIN sg_group_report cr
  ON cr.group_id = g.id
 AND cr.year_no = COALESCE(fs.year_no, ls.year_no)
ORDER BY g.group_no ASC;
"@

    $mysqlArgs = @(
        "--defaults-extra-file=$DefaultsFile",
        '--batch',
        '--raw',
        '--default-character-set=utf8mb4',
        "--execute=$summarySql",
        $Database
    )

    $summaryOutput = & $MysqlExe @mysqlArgs 2>&1
    if ($LASTEXITCODE -ne 0) {
        $message = ($summaryOutput | ForEach-Object { [string]$_ }) -join "`n"
        throw "结果摘要查询失败，退出码：$LASTEXITCODE。$message"
    }

    Convert-TsvLinesToCsvFile -TsvLines $summaryOutput -CsvPath $CsvPath
}

function Get-GitCommit {
    param([string]$Root)
    $git = Get-Command git -ErrorAction SilentlyContinue | Select-Object -First 1
    if (-not $git) {
        return 'unknown'
    }
    $commit = & $git.Source -C $Root rev-parse --short HEAD 2>$null
    if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($commit)) {
        return 'unknown'
    }
    return ([string]$commit).Trim()
}

$projectRoot = Split-Path $PSScriptRoot -Parent
$resolvedConfigPath = Resolve-WorkspacePath -BasePath $projectRoot -InputPath $ConfigPath
$resolvedOutputDir = if ([string]::IsNullOrWhiteSpace($OutputDir)) {
    Join-Path $projectRoot 'backup'
}
else {
    Resolve-WorkspacePath -BasePath $projectRoot -InputPath $OutputDir
}

if ([string]::IsNullOrWhiteSpace($BackupName)) {
    $BackupName = Read-Host '请输入备份文件名'
}
$backupBaseName = Get-SafeBackupBaseName -Name $BackupName

New-Item -ItemType Directory -Force -Path $resolvedOutputDir | Out-Null
$zipPath = Get-UniqueZipPath -Directory $resolvedOutputDir -BaseName $backupBaseName
$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ('sandbox-game-backup-' + [System.Guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Force -Path $tempDir | Out-Null

$defaultsFile = $null
$summaryStatus = '未生成'
$summaryError = ''

try {
    [Console]::OutputEncoding = New-Object System.Text.UTF8Encoding($false)
    $OutputEncoding = New-Object System.Text.UTF8Encoding($false)

    $dsn = Read-CompetitionDsn -Path $resolvedConfigPath
    $connection = Parse-MySqlDsn -Dsn $dsn

    $dumpExe = Find-MySqlTool -ToolName 'mysqldump.exe'
    if (-not $dumpExe) {
        $dumpExe = Find-MySqlTool -ToolName 'mysqldump'
    }
    if (-not $dumpExe) {
        throw "未找到 mysqldump.exe。请确认 MySQL / MariaDB 客户端已安装，并已加入 PATH，或安装在常见 Program Files 目录下。"
    }

    $mysqlExe = Find-MySqlTool -ToolName 'mysql.exe'
    if (-not $mysqlExe) {
        $mysqlExe = Find-MySqlTool -ToolName 'mysql'
    }

    $defaultsFile = New-DefaultsExtraFile -Connection $connection -Directory ([System.IO.Path]::GetTempPath())
    $sqlFileName = $backupBaseName + '.sql'
    $summaryFileName = $backupBaseName + '_结果摘要.csv'
    $infoFileName = $backupBaseName + '_备份说明.txt'
    $sqlPath = Join-Path $tempDir $sqlFileName
    $summaryPath = Join-Path $tempDir $summaryFileName
    $infoPath = Join-Path $tempDir $infoFileName

    Write-Host "正在备份数据库：$($connection.Database)"
    Invoke-DatabaseDump -DumpExe $dumpExe -DefaultsFile $defaultsFile -Database $connection.Database -SqlPath $sqlPath

    if ($mysqlExe) {
        try {
            Write-Host '正在生成结果摘要 CSV...'
            New-ResultSummaryCsv -MysqlExe $mysqlExe -DefaultsFile $defaultsFile -Database $connection.Database -CsvPath $summaryPath
            $summaryStatus = '成功'
        }
        catch {
            $summaryStatus = '失败'
            $summaryError = $_.Exception.Message
            Write-Warning "数据库已完整备份，但结果摘要生成失败：$summaryError"
        }
    }
    else {
        $summaryStatus = '失败'
        $summaryError = '未找到 mysql.exe，无法生成结果摘要 CSV。完整 SQL 备份不受影响。'
        Write-Warning $summaryError
    }

    $commit = Get-GitCommit -Root $projectRoot
    $infoLines = @(
        '沙盘经营系统数据库备份说明',
        '',
        "备份时间：$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')",
        "部署目录：$projectRoot",
        "配置文件：$resolvedConfigPath",
        "数据库主机：$($connection.Host)",
        "数据库端口：$($connection.Port)",
        "数据库名称：$($connection.Database)",
        "代码提交：$commit",
        "SQL 文件：$sqlFileName",
        "结果摘要：$summaryStatus",
        "结果摘要文件：$(if ($summaryStatus -eq '成功') { $summaryFileName } else { '未生成' })",
        "结果摘要失败原因：$(if ([string]::IsNullOrWhiteSpace($summaryError)) { '无' } else { $summaryError })",
        '',
        '说明：',
        '1. SQL 文件是完整数据库备份，包含表结构和全部比赛数据。',
        '2. 结果摘要 CSV 仅用于快速查看最终排名、收入、利润、权益和总监得分。',
        '3. 如需严肃追溯订单、回退、快照、奖罚等细节，应将 SQL 恢复到临时数据库后查询。',
        '4. 本备份脚本不会清空数据库，也不会修改比赛数据。'
    )
    Write-Utf8BomLines -Path $infoPath -Lines $infoLines

    if ($defaultsFile -and (Test-Path -LiteralPath $defaultsFile)) {
        Remove-Item -LiteralPath $defaultsFile -Force -ErrorAction SilentlyContinue
        $defaultsFile = $null
    }

    Compress-Archive -Path (Join-Path $tempDir '*') -DestinationPath $zipPath -Force

    Write-Host ''
    Write-Host "备份成功：$zipPath"
    if ($summaryStatus -ne '成功') {
        Write-Host '提示：完整数据库 SQL 已备份成功，但结果摘要 CSV 未生成，请查看备份说明。'
    }
}
finally {
    if ($defaultsFile -and (Test-Path -LiteralPath $defaultsFile)) {
        Remove-Item -LiteralPath $defaultsFile -Force -ErrorAction SilentlyContinue
    }
    if (Test-Path -LiteralPath $tempDir) {
        Remove-Item -LiteralPath $tempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}
