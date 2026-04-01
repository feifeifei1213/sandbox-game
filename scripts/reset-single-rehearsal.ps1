param(
    [string]$ConfigPath = '.\configs\local-single.yaml'
)

$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path $PSScriptRoot -Parent
$resolvedConfigPath = if ([System.IO.Path]::IsPathRooted($ConfigPath)) { $ConfigPath } else { Join-Path $projectRoot $ConfigPath }
$goCache = Join-Path $projectRoot '.gocache'
$goModCache = Join-Path $projectRoot '.cache\gomod'

New-Item -ItemType Directory -Force -Path $goCache | Out-Null
New-Item -ItemType Directory -Force -Path $goModCache | Out-Null

$env:GOCACHE = $goCache
$env:GOMODCACHE = $goModCache

Push-Location $projectRoot
try {
    go run .\cmd\dbtool\main.go -config $resolvedConfigPath -action reset-single
}
finally {
    Pop-Location
}
