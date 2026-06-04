param(
    [string]$ConfigPath = '.\configs\competition.yaml'
)

$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path $PSScriptRoot -Parent
$resolvedConfigPath = if ([System.IO.Path]::IsPathRooted($ConfigPath)) { $ConfigPath } else { Join-Path $projectRoot $ConfigPath }
$goCache = Join-Path $projectRoot '.gocache'
$goModCache = Join-Path $projectRoot '.cache\gomod'
$dbToolExe = Join-Path $projectRoot 'bin\dbtool.exe'

New-Item -ItemType Directory -Force -Path $goCache | Out-Null
New-Item -ItemType Directory -Force -Path $goModCache | Out-Null

$env:GOCACHE = $goCache
$env:GOMODCACHE = $goModCache
$env:GOTELEMETRY = 'off'

Push-Location $projectRoot
try {
    if (Test-Path $dbToolExe) {
        & $dbToolExe -config $resolvedConfigPath -action init-competition
        if ($LASTEXITCODE -ne 0) {
            throw "dbtool init-competition failed with exit code $LASTEXITCODE"
        }
    }
    else {
        go run .\cmd\dbtool -config $resolvedConfigPath -action init-competition
        if ($LASTEXITCODE -ne 0) {
            throw "go run init-competition failed with exit code $LASTEXITCODE"
        }
    }
}
finally {
    Pop-Location
}
