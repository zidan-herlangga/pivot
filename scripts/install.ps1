# pivot - Multi-platform runtime version switcher
# Install script for Windows
# Usage: iwr -Uri https://raw.githubusercontent.com/zidan-herlangga/pivot/main/scripts/install.ps1 | iex

$Repo = "zidan-herlangga/pivot"
$BinDir = "$env:USERPROFILE\.pivot\bin"
$InstallDir = "$env:USERPROFILE\.pivot"

# Ensure directories
New-Item -ItemType Directory -Path $BinDir -Force | Out-Null
New-Item -ItemType Directory -Path "$InstallDir\runtimes" -Force | Out-Null
New-Item -ItemType Directory -Path "$InstallDir\projects" -Force | Out-Null
New-Item -ItemType Directory -Path "$InstallDir\profiles" -Force | Out-Null

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# GitHub API to get latest release
Write-Host "  Checking latest version..." -ForegroundColor Cyan
try {
    $Latest = (Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -ErrorAction Stop).tag_name
} catch {
    $Latest = "latest"
}

# Download binary
$Binary = "pivot-windows-$Arch.zip"
$Url = "https://github.com/$Repo/releases/download/$Latest/$Binary"
$ZipPath = "$env:TEMP\pivot.zip"

Write-Host "  Downloading $Binary ..." -ForegroundColor Cyan
try {
    Invoke-WebRequest -Uri $Url -OutFile $ZipPath -ErrorAction Stop
    Expand-Archive -Path $ZipPath -DestinationPath $BinDir -Force
    Remove-Item $ZipPath -Force
} catch {
    Write-Host "  Pre-built binary not found. Building from source..." -ForegroundColor Yellow
    if (Get-Command go -ErrorAction SilentlyContinue) {
        Push-Location "$env:TEMP"
        git clone --depth 1 "https://github.com/$Repo.git" 2>$null
        if (Test-Path "pivot") {
            Set-Location "pivot"
            go build -o "$BinDir\pivot.exe" .
        } else {
            Write-Host "  Could not clone repository." -ForegroundColor Red
            Write-Host "  Install Go from https://go.dev then run: go install github.com/$Repo@latest"
            exit 1
        }
        Pop-Location
    } else {
        Write-Host "  Could not download pre-built binary and Go is not installed." -ForegroundColor Red
        Write-Host "  Install Go from https://go.dev or download manually from:" -ForegroundColor Red
        Write-Host "    https://github.com/$Repo/releases"
        exit 1
    }
}

# Add to User PATH
$curPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($curPath -notmatch [regex]::Escape($BinDir)) {
    [Environment]::SetEnvironmentVariable('Path', "$BinDir;$curPath", 'User')
    $env:Path = "$BinDir;$env:Path"
    Write-Host "  Added $BinDir to User PATH." -ForegroundColor Green
}

Write-Host ""
Write-Host "  pivot installed successfully!" -ForegroundColor Green
Write-Host "  Run 'pivot' to start." -ForegroundColor White
Write-Host "  Restart your terminal or run: refreshenv" -ForegroundColor Gray
