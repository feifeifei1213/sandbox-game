param(
    [string]$ConfigPath = '.\configs\competition.yaml',
    [int]$Port = 8080,
    [string]$LogDir = ''
)

$ErrorActionPreference = 'Stop'

$projectRoot = Split-Path $PSScriptRoot -Parent
$resolvedConfigPath = if ([System.IO.Path]::IsPathRooted($ConfigPath)) {
    $ConfigPath
}
else {
    Join-Path $projectRoot $ConfigPath
}

if ([string]::IsNullOrWhiteSpace($LogDir)) {
    $LogDir = Join-Path $projectRoot 'logs'
}
elseif (-not [System.IO.Path]::IsPathRooted($LogDir)) {
    $LogDir = Join-Path $projectRoot $LogDir
}

New-Item -ItemType Directory -Force -Path $LogDir | Out-Null

$startScript = Join-Path $PSScriptRoot 'start-competition.ps1'
$stdoutPath = Join-Path $LogDir ("server-{0}.log" -f $Port)
$stderrPath = Join-Path $LogDir ("server-{0}.err.log" -f $Port)
$arguments = @(
    '-NoProfile',
    '-ExecutionPolicy',
    'Bypass',
    '-File',
    $startScript,
    '-ConfigPath',
    $resolvedConfigPath
)

Start-Process `
    -FilePath 'powershell.exe' `
    -ArgumentList $arguments `
    -WorkingDirectory $projectRoot `
    -WindowStyle Hidden `
    -RedirectStandardOutput $stdoutPath `
    -RedirectStandardError $stderrPath | Out-Null

Write-Host "Sandbox game server starting in background. Logs: $LogDir"
