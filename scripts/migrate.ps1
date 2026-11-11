$ErrorActionPreference = 'Stop'
. "$PSScriptRoot\environment.ps1"
try {
 Push-Location $projectRoot
 & $goExe run ./cmd/migrate
 if ($LASTEXITCODE -ne 0) { throw 'Database migration failed' }
} finally { Pop-Location }
