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
  $envContent = $envContent.Replace('DATABASE_URL=', "DATABASE_URL=postgres://abclinic:${dbPassword}@localhost:55433/abclinic?sslmode=disable")
  [IO.File]::WriteAllText((Join-Path $projectRoot '.env'), $envContent, (New-Object Text.UTF8Encoding($false)))
  Write-Host 'Created .env with a random local PostgreSQL password. Do not commit it.'
 }
 $envFile = Join-Path $projectRoot '.env'
 $envContent = [IO.File]::ReadAllText($envFile)
 foreach ($key in @('JWT_ACCESS_SECRET','JWT_REFRESH_SECRET')) {
  $pattern = '(?m)^' + $key + '=([^\r\n]*)'
  $existing = [regex]::Match($envContent, $pattern)
  if ($existing.Success -and $existing.Groups[1].Value.Trim() -ne '') { continue }
  $secretBytes = New-Object byte[] 32
  $rng = [Security.Cryptography.RandomNumberGenerator]::Create()
  try { $rng.GetBytes($secretBytes) } finally { $rng.Dispose() }
  $secret = ([BitConverter]::ToString($secretBytes)).Replace('-', '').ToLowerInvariant()
  if ($existing.Success) { $envContent = [regex]::Replace($envContent, $pattern, "$key=$secret") }
  else { $envContent = $envContent.TrimEnd() + "`n$key=$secret`n" }
 }
 [IO.File]::WriteAllText($envFile, $envContent, (New-Object Text.UTF8Encoding($false)))
 & $goExe mod download
 if ($LASTEXITCODE -ne 0) { throw 'Dependency download failed' }
 Write-Host 'Ready: powershell -ExecutionPolicy Bypass -File .\scripts\run.ps1'
} finally { Pop-Location }
