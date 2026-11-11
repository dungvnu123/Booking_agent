param([string]$BaseUrl = 'http://localhost:8081')
$ErrorActionPreference = 'Stop'
$base = $BaseUrl.TrimEnd('/')
$health = Invoke-RestMethod "$base/healthz"
if ($health.status -ne 'ok' -or $health.mode -ne 'scaffold') { throw 'Not the ABClinic scaffold server; restart the old demo server first.' }
$page = Invoke-WebRequest -UseBasicParsing "$base/swagger/index.html"
if ($page.StatusCode -ne 200) { throw 'Swagger page unavailable' }
$spec = Invoke-RestMethod "$base/api/v3/api-docs"
if (-not $spec.info.title.StartsWith('ABClinic Go')) { throw 'Wrong API document' }
$count = 0
foreach ($path in $spec.paths.PSObject.Properties) {
 foreach ($operation in $path.Value.PSObject.Properties) {
  if ($operation.Name -notin @('get','post','put','patch','delete') -or $operation.Value.'x-implementation-status' -ne 'not-implemented') { continue }
  $urlPath = [regex]::Replace($path.Name, '\{[^}]+\}', '1')
  try {
   Invoke-WebRequest -UseBasicParsing -Uri "$base$urlPath" -Method $operation.Name -TimeoutSec 5 | Out-Null
   throw "Expected 501 for $($operation.Name) $urlPath"
  } catch [Net.WebException] {
   if (-not $_.Exception.Response -or [int]$_.Exception.Response.StatusCode -ne 501) { throw }
  }
  $count++
 }
}
if ($count -ne 103) { throw "Expected 103 scaffold routes, found $count" }
foreach ($route in @('register','register-patient','login','refresh_token','logout')) {
 try {
  Invoke-WebRequest -UseBasicParsing -Uri "$base/api/auth/$route" -Method POST -TimeoutSec 5 | Out-Null
  throw "Expected 400 for empty auth request: $route"
 } catch [Net.WebException] {
  if (-not $_.Exception.Response -or [int]$_.Exception.Response.StatusCode -ne 400) { throw }
 }
}
Write-Host 'Health + Swagger + 103 scaffold routes (501) + 5 auth validation routes (400) OK. No business data was written.'
