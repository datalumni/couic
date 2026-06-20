param([string]$Target = "linux")

switch ($Target) {
  "linux"   { $ext = "" }
  "windows" { $ext = ".exe" }
  "darwin"  { Write-Error "macOS builds must run natively (see AGENTS.md for commands)."; exit 1 }
  default   { Write-Error "Usage: $0 {linux|windows}"; exit 1 }
}

docker build -f scripts/Dockerfile --build-arg TARGET=$Target -t couic-builder .
docker create --name tmp couic-builder
docker cp tmp:/couic$ext ./couic$ext
docker rm tmp
Write-Host "Built: couic$ext"
