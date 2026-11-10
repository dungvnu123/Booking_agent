$ErrorActionPreference = 'Stop'
. "$PSScriptRoot\environment.ps1"
try {
 Push-Location $projectRoot
 & $goExe version
 if ($LASTEXITCODE -ne 0) { throw 'Go is unavailable' }
 if (-not (Test-Path -LiteralPath '.env')) {
  $random = New-Object byte[] 24
  $rng = [Security.Cryptography.RandomNumberGenerator]::Create()
  try { $rng.GetBytes($random) } finally { $rng.Dispose() }
  $dbPassword = ([BitConverter]::ToString($random)).Replace('-', '').ToLowerInvariant()
  $envContent = [IO.File]::ReadAllText((Join-Path $projectRoot '.env.example'))
  $envContent = $envContent.Replace('POSTGRES_PASSWORD=', "POSTGRES_PASSWORD=$dbPassword")
  $envContent = $envContent.Replace('DATABASE_URL=', "DATABASE_URL=postgres://abclinic:${dbPassword}@localhost:5433/abclinic?sslmode=disable")
  [IO.File]::WriteAllText((Join-Path $projectRoot '.env'), $envContent, (New-Object Text.UTF8Encoding($false)))
  Write-Host 'Created .env with a random local PostgreSQL password. Do not commit it.'
 }
 & $goExe mod download
 if ($LASTEXITCODE -ne 0) { throw 'Dependency download failed' }
 Write-Host 'Ready: powershell -ExecutionPolicy Bypass -File .\scripts\run.ps1'
} finally { Pop-Location }
