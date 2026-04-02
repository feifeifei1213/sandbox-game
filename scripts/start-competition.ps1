param(
    [string]$ConfigPath = '.\configs\competition.yaml'
)

$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path $PSScriptRoot -Parent
$resolvedConfigPath = if ([System.IO.Path]::IsPathRooted($ConfigPath)) { $ConfigPath } else { Join-Path $projectRoot $ConfigPath }
$goCache = Join-Path $projectRoot '.gocache'
$goModCache = Join-Path $projectRoot '.cache\gomod'
$serverExe = Join-Path $projectRoot 'bin\sandbox-game-server.exe'

New-Item -ItemType Directory -Force -Path $goCache | Out-Null
New-Item -ItemType Directory -Force -Path $goModCache | Out-Null

$env:GOCACHE = $goCache
$env:GOMODCACHE = $goModCache
$env:GOTELEMETRY = 'off'

Push-Location $projectRoot
try {
    if (Test-Path $serverExe) {
        & $serverExe -config $resolvedConfigPath
        if ($LASTEXITCODE -ne 0) {
            throw "sandbox-game-server exited with code $LASTEXITCODE"
        }
    }
    else {
        go run .\cmd\server\main.go -config $resolvedConfigPath
        if ($LASTEXITCODE -ne 0) {
            throw "go run competition server failed with exit code $LASTEXITCODE"
        }
    }
}
finally {
    Pop-Location
}
