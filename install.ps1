# Install script for pray-cli (Windows)
# Usage: iwr -useb https://raw.githubusercontent.com/AbdElrahmaN31/pray-cli/main/install.ps1 | iex

$ErrorActionPreference = 'Stop'

Write-Host "Fetching latest version..." -ForegroundColor Yellow

# Get latest version
try {
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/AbdElrahmaN31/pray-cli/releases/latest"
    $version = $release.tag_name
    $versionNum = $version.TrimStart('v')
} catch {
    Write-Host "Failed to fetch latest version" -ForegroundColor Red
    exit 1
}

Write-Host "Latest version: $version" -ForegroundColor Green

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# Construct download URL
$filename = "pray-cli_${versionNum}_windows_${arch}.zip"
$url = "https://github.com/AbdElrahmaN31/pray-cli/releases/download/$version/$filename"

Write-Host "Downloading $filename..." -ForegroundColor Yellow

# Create temporary directory
$tmpDir = Join-Path $env:TEMP "pray-install-$(Get-Random)"
New-Item -ItemType Directory -Path $tmpDir | Out-Null

try {
    # Download
    $zipPath = Join-Path $tmpDir "pray.zip"
    Invoke-WebRequest -Uri $url -OutFile $zipPath

    # Extract
    Expand-Archive -Path $zipPath -DestinationPath $tmpDir -Force

    # Determine install location
    $installDir = "$env:LOCALAPPDATA\Programs\pray"

    # Create install directory if it doesn't exist
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir | Out-Null
    }

    # Copy binary
    $binaryPath = Join-Path $tmpDir "pray.exe"
    Copy-Item $binaryPath -Destination $installDir -Force

    Write-Host "Installed pray to $installDir" -ForegroundColor Green

    # Add to PATH if not already there
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$installDir*") {
        [Environment]::SetEnvironmentVariable(
            "Path",
            "$userPath;$installDir",
            "User"
        )
        Write-Host "Added $installDir to your PATH" -ForegroundColor Green
        Write-Host "Please restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
    }

    Write-Host ""
    Write-Host "Installation complete!" -ForegroundColor Green
    Write-Host "Run 'pray --help' to get started (restart your terminal first)" -ForegroundColor Cyan

} finally {
    # Cleanup
    Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
}
