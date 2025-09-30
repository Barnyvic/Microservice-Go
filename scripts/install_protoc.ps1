# PowerShell script to install protoc on Windows

$ErrorActionPreference = "Stop"

Write-Host "Installing Protocol Buffers Compiler (protoc)..." -ForegroundColor Green

# Define protoc version and download URL
$PROTOC_VERSION = "25.1"
$PROTOC_ZIP = "protoc-$PROTOC_VERSION-win64.zip"
$PROTOC_URL = "https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/$PROTOC_ZIP"
$INSTALL_DIR = "$env:USERPROFILE\.protoc"
$TEMP_DIR = "$env:TEMP\protoc_install"

# Create directories
New-Item -ItemType Directory -Force -Path $TEMP_DIR | Out-Null
New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null

try {
    # Download protoc
    Write-Host "Downloading protoc $PROTOC_VERSION..." -ForegroundColor Yellow
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest -Uri $PROTOC_URL -OutFile "$TEMP_DIR\$PROTOC_ZIP"
    
    # Extract protoc
    Write-Host "Extracting protoc..." -ForegroundColor Yellow
    Expand-Archive -Path "$TEMP_DIR\$PROTOC_ZIP" -DestinationPath $INSTALL_DIR -Force
    
    # Add to PATH for current session
    $env:PATH = "$INSTALL_DIR\bin;$env:PATH"
    
    # Verify installation
    Write-Host "Verifying protoc installation..." -ForegroundColor Yellow
    & "$INSTALL_DIR\bin\protoc.exe" --version
    
    Write-Host "`nProtoc installed successfully!" -ForegroundColor Green
    Write-Host "Location: $INSTALL_DIR\bin" -ForegroundColor Cyan
    Write-Host "`nTo make this permanent, add the following to your PATH:" -ForegroundColor Yellow
    Write-Host "$INSTALL_DIR\bin" -ForegroundColor Cyan
    
} catch {
    Write-Host "Error installing protoc: $_" -ForegroundColor Red
    exit 1
} finally {
    # Cleanup
    Remove-Item -Path $TEMP_DIR -Recurse -Force -ErrorAction SilentlyContinue
}

Write-Host "`nGenerating Go protobuf files..." -ForegroundColor Green

# Change to project directory
$PROJECT_DIR = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
Set-Location $PROJECT_DIR

# Generate protobuf files
try {
    & "$INSTALL_DIR\bin\protoc.exe" `
        --go_out=. --go_opt=paths=source_relative `
        --go-grpc_out=. --go-grpc_opt=paths=source_relative `
        proto/product.proto proto/subscription.proto
    
    Write-Host "`nProtobuf files generated successfully!" -ForegroundColor Green
    Write-Host "Generated files:" -ForegroundColor Cyan
    Get-ChildItem -Path "proto" -Filter "*.pb.go" -Recurse | ForEach-Object { Write-Host "  - $($_.FullName)" -ForegroundColor Gray }
    
} catch {
    Write-Host "Error generating protobuf files: $_" -ForegroundColor Red
    exit 1
}

Write-Host "`nSetup complete! You can now build and run the application." -ForegroundColor Green

