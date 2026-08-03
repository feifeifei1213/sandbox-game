param(
    [string]$ConfigPath = '.\configs\competition.yaml',
    [int]$DefaultPort = 8080
)

$ErrorActionPreference = 'Stop'

function Resolve-ConfigPath {
    param([string]$Path)

    if ([System.IO.Path]::IsPathRooted($Path)) {
        return $Path
    }

    $projectRoot = Split-Path $PSScriptRoot -Parent
    return Join-Path $projectRoot $Path
}

$resolvedConfigPath = Resolve-ConfigPath -Path $ConfigPath
if (-not (Test-Path -LiteralPath $resolvedConfigPath)) {
    Write-Output $DefaultPort
    exit 0
}

$inServerSection = $false
foreach ($line in Get-Content -LiteralPath $resolvedConfigPath -Encoding UTF8) {
    if ($line -match '^\s*server\s*:\s*$') {
        $inServerSection = $true
        continue
    }
    if ($inServerSection -and $line -match '^\S') {
        $inServerSection = $false
    }
    if ($inServerSection -and $line -match '^\s*port\s*:\s*([0-9]+)\s*$') {
        Write-Output ([int]$matches[1])
        exit 0
    }
}

Write-Output $DefaultPort
