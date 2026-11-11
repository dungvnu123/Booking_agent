$ErrorActionPreference = 'Stop'
. "$PSScriptRoot\environment.ps1"
try {
 Push-Location $projectRoot
 $goFiles = @(Get-ChildItem -LiteralPath cmd,internal,api,tests,migrations -Recurse -Filter '*.go' -File | Select-Object -ExpandProperty FullName)
 $unformatted = & (Join-Path (Split-Path -Parent $goExe) 'gofmt.exe') -l $goFiles
 if ($LASTEXITCODE -ne 0 -or $unformatted) { throw "Run gofmt on: $unformatted" }
 & $goExe test ./...
 if ($LASTEXITCODE -ne 0) { throw 'Tests failed' }
 & $goExe vet ./...
 if ($LASTEXITCODE -ne 0) { throw 'Vet failed' }
 & $goExe build ./cmd/api ./cmd/worker ./cmd/migrate
 if ($LASTEXITCODE -ne 0) { throw 'Build failed' }
} finally { Pop-Location }
