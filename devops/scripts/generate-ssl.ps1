param([string]$Domain = "yago.loc")

$ErrorActionPreference = "Stop"

$CertDir = Join-Path $PSScriptRoot "..\nginx\ssl"
$Mkcert  = Join-Path $PSScriptRoot "..\nginx\mkcert.exe"

New-Item -ItemType Directory -Force -Path $CertDir | Out-Null

if (-not (Test-Path $Mkcert)) {
    Write-Host "❌ mkcert.exe not found: $Mkcert" -ForegroundColor Red
    exit 1
}

& $Mkcert -install | Out-Null

$ErrorActionPreference = "SilentlyContinue"

& $Mkcert `
    -cert-file (Join-Path $CertDir "$Domain.crt") `
    -key-file (Join-Path $CertDir "$Domain.key") `
    $Domain `
    "*.$Domain" `
    "admin.$Domain" `
    "storage.$Domain" `
    "api.$Domain" | Out-Null
$ErrorActionPreference = "Stop"

if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}