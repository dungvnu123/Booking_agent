param([ValidateSet('api','worker')][string]$Target = 'api')
$ErrorActionPreference = 'Stop'
. "$PSScriptRoot\environment.ps1"
try {
 Push-Location $projectRoot
 & $goExe run "./cmd/$Target"
 if ($LASTEXITCODE -ne 0) { throw "ABClinic $Target stopped with an error" }
} finally { Pop-Location }
