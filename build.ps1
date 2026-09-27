$ErrorActionPreference = "Stop"

$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

$versionMatch = Select-String -Path ".\win32.go" -Pattern '^\s*version\s*=\s*"([^"]+)"'
if (-not $versionMatch) {
    throw "Could not determine MagzClicker version from win32.go"
}
$version = $versionMatch.Matches[0].Groups[1].Value
if ($version -notmatch '^\d+\.\d+\.\d+$') {
    throw "Unexpected MagzClicker version format: $version"
}

$resourceFile = ".\rsrc_windows_amd64.syso"

try {
    Remove-Item $resourceFile -Force -ErrorAction SilentlyContinue

    go run github.com/tc-hib/go-winres@v0.3.3 simply `
        --arch amd64 `
        --out rsrc `
        --icon assets/MagzClicker.ico `
        --manifest gui `
        --product-name "MagzClicker" `
        --file-description "MagzClicker - Auto Clicker for Windows" `
        --file-version "$version.0" `
        --product-version "$version.0" `
        --original-filename "MagzClicker.exe"

    go build -trimpath -ldflags="-s -w -H=windowsgui" -o MagzClicker.exe .
}
finally {
    Remove-Item $resourceFile -Force -ErrorAction SilentlyContinue
}

Write-Host ""
Write-Host "MagzClicker $version build complete."
Get-FileHash .\MagzClicker.exe -Algorithm SHA256
