$projectRoot = Split-Path -Parent $PSScriptRoot
$goCommand = Get-Command go -ErrorAction SilentlyContinue
$goExe = if (Test-Path -LiteralPath 'C:\Program Files\Go\bin\go.exe') { 'C:\Program Files\Go\bin\go.exe' } elseif ($goCommand) { $goCommand.Source } else { throw 'Install Go 1.27+ from https://go.dev/dl/ first.' }
$env:GOROOT = Split-Path -Parent (Split-Path -Parent $goExe)
if (-not (Test-Path -LiteralPath (Join-Path $env:GOROOT 'src\fmt'))) { throw 'Go installation is incomplete.' }
$env:PATH = "$(Split-Path -Parent $goExe);$env:PATH"
$env:GOPATH = Join-Path $projectRoot '.tools\gopath'
$env:GOCACHE = Join-Path $projectRoot '.tools\cache'
$env:GOTOOLCHAIN = 'local'
$env:GOWORK = Join-Path $projectRoot 'go.work'
if ($env:GOFLAGS -notmatch '(^|\s)-buildvcs=') { $env:GOFLAGS = "$env:GOFLAGS -buildvcs=false".Trim() }
