$ErrorActionPreference = 'Stop'
. "$PSScriptRoot\environment.ps1"
$binDir = Join-Path $projectRoot '.tools\booking-app'
New-Item -ItemType Directory -Force -Path $binDir | Out-Null
$binary = Join-Path $binDir 'api.exe'
try {
 Push-Location $projectRoot
 & $goExe build -o $binary ./cmd/api
 if ($LASTEXITCODE -ne 0) { throw 'API build failed' }
} finally { Pop-Location }
$oldAddr = $env:SERVER_ADDR
$oldMode = $env:GIN_MODE
$child = $null
try {
 $env:SERVER_ADDR = '127.0.0.1:18081'
 $env:GIN_MODE = 'release'
 $child = Start-Process -FilePath $binary -WorkingDirectory $projectRoot -WindowStyle Hidden -PassThru -RedirectStandardOutput (Join-Path $binDir 'smoke.out.log') -RedirectStandardError (Join-Path $binDir 'smoke.err.log')
 $ready = $false
 for ($attempt = 0; $attempt -lt 30; $attempt++) {
  if ($child.HasExited) { throw 'API exited early; check .tools/booking-app/smoke.err.log' }
  try {
   $health = Invoke-RestMethod 'http://127.0.0.1:18081/healthz' -TimeoutSec 1
   if ($health.mode -eq 'scaffold') { $ready = $true; break }
  } catch { Start-Sleep -Milliseconds 100 }
 }
 if (-not $ready) { throw 'API startup timeout' }
 & "$PSScriptRoot\test-api.ps1" -BaseUrl 'http://127.0.0.1:18081'
} finally {
 if ($child -and -not $child.HasExited) { Stop-Process -Id $child.Id }
 $env:SERVER_ADDR = $oldAddr
 $env:GIN_MODE = $oldMode
}
