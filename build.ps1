$ErrorActionPreference = "Stop"

$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

python tools/make_icon.py assets/MagzClicker.ico
go build -trimpath -ldflags="-s -w -H=windowsgui" -o MagzClicker.exe .
python tools/embed_icon.py MagzClicker.exe assets/MagzClicker.ico

Write-Host ""
Write-Host "MagzClicker 1.4.3 build complete."
Get-FileHash .\MagzClicker.exe -Algorithm SHA256
