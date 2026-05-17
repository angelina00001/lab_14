# Run: powershell -ExecutionPolicy Bypass -File .\scripts\run-tests.ps1
# Or: scripts\run-tests.cmd
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
$Root = Split-Path $PSScriptRoot -Parent

Write-Host "== go test ==" -ForegroundColor Cyan
Push-Location $Root
go test ./... -count=1
if ($LASTEXITCODE -ne 0) { Pop-Location; exit $LASTEXITCODE }
Pop-Location

Write-Host "== cargo test (no PyO3) ==" -ForegroundColor Cyan
Push-Location (Join-Path $Root "rust\sports-core")
cargo test --no-default-features -v
if ($LASTEXITCODE -ne 0) { Pop-Location; exit $LASTEXITCODE }
Pop-Location

Write-Host "== python venv + maturin + pytest ==" -ForegroundColor Cyan
$PythonDir = Join-Path $Root "python"
Push-Location $PythonDir

$Venv = Join-Path $PythonDir ".venv"
if (-not (Test-Path $Venv)) {
    # Default "python" (3.14 on your PC). Other version: py -3.12 -m venv .venv
    python -m venv $Venv
}

$py = Join-Path $Venv "Scripts\python.exe"
$pip = Join-Path $Venv "Scripts\pip.exe"
Write-Host "Python:" -NoNewline
& $py --version

$env:VIRTUAL_ENV = $Venv
$env:PATH = (Join-Path $Venv "Scripts") + ";" + $env:PATH

& $pip install -q -U pip
& $pip install -q -e ".[dev]"
& $pip install -q maturin
& $py -m maturin develop --release --features extension-module
if ($LASTEXITCODE -ne 0) { Pop-Location; exit $LASTEXITCODE }

& $py -m pytest -v
$code = $LASTEXITCODE
Pop-Location

if ($code -ne 0) { exit $code }
Write-Host ""
Write-Host "OK: all tests passed" -ForegroundColor Green
