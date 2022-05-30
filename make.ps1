param(
  [Parameter()]
  [switch]$test = $false,
  [switch]$run = $false,
  [switch]$docs = $false
)

$build = $true
# This is the windows build script in powershell
$folder = ".\deploy\bin"
Remove-Item -Path $folder -Force -Recurse
New-Item $folder -ItemType Directory


$version = (git describe --tags --always)
$IMPORT = "github.com/rest-api/internal/version"
$time = Get-Date -Format "yyyy-MM-ddTHH:mm:ss"
$GOFLAGS = "-s -w -X $IMPORT.Version=${version} -X ${IMPORT}.BuildTimeUTC=$time -X ${IMPORT}.AppName=rest-api"
$BLDDIR = "deploy\bin"


if ( $test ) {
  $build = $false
  go test -cover -race .\...
}

if ( $build ) {
  Write-Host "go build -ldflags '$GOFLAGS' -o $BLDDIR/rest-api.exe"
  go build -ldflags "$GOFLAGS" -o $BLDDIR/rest-api.exe
}

if ( $docs ) {
  .\deploy\bin\rest-api.exe -docs
}

if ( $run ) {
  .\deploy\bin\rest-api.exe -debug -pretty
}
