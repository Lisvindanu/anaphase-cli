# Anaphase CLI Installation Script for Windows
# Run with: powershell -ExecutionPolicy Bypass -File install.ps1

Write-Host "🚀 Installing Anaphase CLI..." -ForegroundColor Cyan

# Check if Go is installed
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Go is not installed. Please install Go first:" -ForegroundColor Red
    Write-Host "   https://go.dev/dl/" -ForegroundColor Yellow
    exit 1
}

# Install the binary
Write-Host "📦 Installing anaphase binary..." -ForegroundColor Cyan
go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@latest

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Installation failed" -ForegroundColor Red
    exit 1
}

# Get GOPATH
$GOPATH = go env GOPATH
$GOBIN = "$GOPATH\bin"

# Check if GOBIN is in PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
$pathNeedsUpdate = $currentPath -notlike "*$GOBIN*"

if ($pathNeedsUpdate) {
    Write-Host ""
    Write-Host "⚠️  $GOBIN is not in your PATH" -ForegroundColor Yellow
    Write-Host ""

    $response = Read-Host "Would you like to add it to your PATH? (y/n)"

    if ($response -match '^[Yy]$') {
        # Add to user PATH
        $newPath = "$currentPath;$GOBIN"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")

        # Update current session
        $env:Path = "$env:Path;$GOBIN"

        Write-Host "✅ Added to PATH" -ForegroundColor Green
        Write-Host ""
        Write-Host "Please restart your terminal or run:" -ForegroundColor Yellow
        Write-Host "  `$env:Path = [System.Environment]::GetEnvironmentVariable('Path','User')" -ForegroundColor Cyan
    } else {
        Write-Host ""
        Write-Host "You can manually add it by:" -ForegroundColor Yellow
        Write-Host "  1. Open System Properties > Environment Variables" -ForegroundColor Cyan
        Write-Host "  2. Edit 'Path' under User variables" -ForegroundColor Cyan
        Write-Host "  3. Add: $GOBIN" -ForegroundColor Cyan
    }
} else {
    Write-Host "✅ PATH already configured" -ForegroundColor Green
}

Write-Host ""
Write-Host "🎉 Installation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Verify installation:" -ForegroundColor Cyan
Write-Host "  anaphase --version" -ForegroundColor White
Write-Host ""
Write-Host "Get started:" -ForegroundColor Cyan
Write-Host "  anaphase init my-project" -ForegroundColor White
Write-Host "  cd my-project" -ForegroundColor White
Write-Host "  anaphase gen domain --name user --prompt `"User with email and name`"" -ForegroundColor White
Write-Host ""
Write-Host "Documentation: https://anaphygon.my.id" -ForegroundColor Cyan
