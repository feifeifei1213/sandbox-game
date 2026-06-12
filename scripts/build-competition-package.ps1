param(
    [string]$ConfigPath = '.\configs\competition.yaml',
    [string]$OutputDir = '.\.runtime\competition-package'
)

$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path $PSScriptRoot -Parent
$resolvedConfigPath = if ([System.IO.Path]::IsPathRooted($ConfigPath)) { $ConfigPath } else { Join-Path $projectRoot $ConfigPath }
$resolvedOutputDir = if ([System.IO.Path]::IsPathRooted($OutputDir)) { $OutputDir } else { Join-Path $projectRoot $OutputDir }
$goCache = Join-Path $projectRoot '.gocache'
$goModCache = Join-Path $projectRoot '.cache\gomod'
$releaseName = 'sandbox-game-competition'
$releaseRoot = Join-Path $resolvedOutputDir $releaseName
$zipPath = Join-Path $resolvedOutputDir ($releaseName + '.zip')

New-Item -ItemType Directory -Force -Path $goCache | Out-Null
New-Item -ItemType Directory -Force -Path $goModCache | Out-Null
New-Item -ItemType Directory -Force -Path $resolvedOutputDir | Out-Null

$env:GOCACHE = $goCache
$env:GOMODCACHE = $goModCache
$env:GOTELEMETRY = 'off'

if (Test-Path $releaseRoot) {
    Remove-Item -LiteralPath $releaseRoot -Recurse -Force
}
if (Test-Path $zipPath) {
    Remove-Item -LiteralPath $zipPath -Force
}

New-Item -ItemType Directory -Force -Path (Join-Path $releaseRoot 'bin') | Out-Null
New-Item -ItemType Directory -Force -Path (Join-Path $releaseRoot 'frontend\dist') | Out-Null
New-Item -ItemType Directory -Force -Path (Join-Path $releaseRoot 'configs') | Out-Null
New-Item -ItemType Directory -Force -Path (Join-Path $releaseRoot 'scripts') | Out-Null
New-Item -ItemType Directory -Force -Path (Join-Path $releaseRoot 'migrations\mysql') | Out-Null
New-Item -ItemType Directory -Force -Path (Join-Path $releaseRoot 'docs') | Out-Null

Push-Location $projectRoot
try {
    Push-Location (Join-Path $projectRoot 'frontend')
    try {
        npm.cmd run build
        if ($LASTEXITCODE -ne 0) {
            throw "frontend build failed with exit code $LASTEXITCODE"
        }
    }
    finally {
        Pop-Location
    }

    go build -o (Join-Path $releaseRoot 'bin\sandbox-game-server.exe') .\cmd\server
    if ($LASTEXITCODE -ne 0) {
        throw "go build server failed with exit code $LASTEXITCODE"
    }
    go build -o (Join-Path $releaseRoot 'bin\dbtool.exe') .\cmd\dbtool
    if ($LASTEXITCODE -ne 0) {
        throw "go build dbtool failed with exit code $LASTEXITCODE"
    }

    Copy-Item -Path (Join-Path $projectRoot 'frontend\dist\*') -Destination (Join-Path $releaseRoot 'frontend\dist') -Recurse -Force
    Copy-Item -Path $resolvedConfigPath -Destination (Join-Path $releaseRoot 'configs\competition.yaml') -Force
    Copy-Item -Path (Join-Path $projectRoot 'scripts\init-competition.ps1') -Destination (Join-Path $releaseRoot 'scripts\init-competition.ps1') -Force
    Copy-Item -Path (Join-Path $projectRoot 'scripts\reset-competition.ps1') -Destination (Join-Path $releaseRoot 'scripts\reset-competition.ps1') -Force
    Copy-Item -Path (Join-Path $projectRoot 'scripts\start-competition.ps1') -Destination (Join-Path $releaseRoot 'scripts\start-competition.ps1') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0001_init.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0001_init.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0003_notice_adjustment.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0003_notice_adjustment.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0004_seed_competition_admin.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0004_seed_competition_admin.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0005_order_admin.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0005_order_admin.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0006_order_generation_refactor.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0006_order_generation_refactor.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0007_order_market_enable.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0007_order_market_enable.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0008_order_forecast_control.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0008_order_forecast_control.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0009_game_edition.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0009_game_edition.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0010_order_market_investment_limit.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0010_order_market_investment_limit.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0011_rollback_snapshot.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0011_rollback_snapshot.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0012_dictionary.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0012_dictionary.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0013_order_template_airport.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0013_order_template_airport.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'migrations\mysql\0014_order_unit_price_precision.sql') -Destination (Join-Path $releaseRoot 'migrations\mysql\0014_order_unit_price_precision.sql') -Force
    Copy-Item -Path (Join-Path $projectRoot 'docs\competition_launch_runbook.md') -Destination (Join-Path $releaseRoot 'docs\competition_launch_runbook.md') -Force
    Copy-Item -Path (Join-Path $projectRoot 'docs\competition_deploy_checklist.md') -Destination (Join-Path $releaseRoot 'docs\competition_deploy_checklist.md') -Force
    Copy-Item -Path (Join-Path $projectRoot 'docs\README.md') -Destination (Join-Path $releaseRoot 'README.md') -Force

    Compress-Archive -Path (Join-Path $releaseRoot '*') -DestinationPath $zipPath -Force
}
finally {
    Pop-Location
}

Write-Host "competition package ready: $releaseRoot"
Write-Host "competition zip ready: $zipPath"
