param(
    [string]$Candidates = 'SandboxGameMySQL,MySQL,MySQL80,MySQL57,mysql,MariaDB,mariadb'
)

$ErrorActionPreference = 'Stop'

$service = $null
foreach ($name in $Candidates.Split(',')) {
    $trimmedName = $name.Trim()
    if ([string]::IsNullOrWhiteSpace($trimmedName)) {
        continue
    }
    $service = Get-Service -Name $trimmedName -ErrorAction SilentlyContinue
    if ($service) {
        break
    }
}

if (-not $service) {
    Write-Error "未找到可用的 MySQL / MariaDB 服务。已尝试：$Candidates"
    exit 2
}

if ($service.Status -ne 'Running') {
    Start-Service -Name $service.Name
    $service.WaitForStatus('Running', [TimeSpan]::FromSeconds(20))
}

Write-Host "MySQL service ready: $($service.Name)"
exit 0
