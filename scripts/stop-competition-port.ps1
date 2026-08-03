param(
    [int]$Port = 8080
)

$ErrorActionPreference = 'Stop'

$connections = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue
$processIds = @($connections | Select-Object -ExpandProperty OwningProcess -Unique)

if ($processIds.Count -eq 0) {
    Write-Host "No process is listening on port $Port."
    exit 0
}

foreach ($processId in $processIds) {
    if ($processId -le 0) {
        continue
    }
    Stop-Process -Id $processId -Force -ErrorAction SilentlyContinue
    Write-Host "Stopped process $processId on port $Port."
}

exit 0
