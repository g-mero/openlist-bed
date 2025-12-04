# Download and configure libvips (current session only)
# Usage: . .\setup-dev.ps1  (note the dot at the beginning)

Write-Host "Configuring libvips..." -ForegroundColor Green

$DEST_DIR = ".\third_party\libvips"

# Check if already downloaded locally
if (Test-Path $DEST_DIR) {
    Write-Host "OK libvips already exists in third_party" -ForegroundColor Green
    
    # Set environment variables for current session
    $vipsPath = Resolve-Path $DEST_DIR
    $env:PATH = "$vipsPath\bin;$env:PATH"
    $env:PKG_CONFIG_PATH = "$vipsPath\lib\pkgconfig"
    $env:CGO_ENABLED = "1"
    
    Write-Host "OK Environment variables set (current session)" -ForegroundColor Green
    return
}

# Download and install
$LIBVIPS_VERSION = "8.17.3"
# Use official Windows build - web version (smaller, sufficient for development)
$DOWNLOAD_URL = "https://github.com/libvips/build-win64-mxe/releases/download/v$LIBVIPS_VERSION/vips-dev-w64-all-$LIBVIPS_VERSION.zip"
$ARCHIVE_NAME = "vips-dev.zip"

if (-not (Test-Path "third_party")) {
    New-Item -ItemType Directory -Path "third_party" -Force | Out-Null
}

try {
    # Check if already downloaded
    if (-not (Test-Path $ARCHIVE_NAME)) {
        Write-Host "Downloading libvips $LIBVIPS_VERSION ..." -ForegroundColor Cyan
        Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $ARCHIVE_NAME -UseBasicParsing
        Write-Host "OK Download completed" -ForegroundColor Green
    } else {
        Write-Host "Using cached $ARCHIVE_NAME" -ForegroundColor Cyan
    }
    
    Write-Host "Extracting..." -ForegroundColor Cyan
    
    # Extract ZIP file (native Windows support)
    Expand-Archive -Path $ARCHIVE_NAME -DestinationPath "third_party" -Force
    
    # The extracted directory should be vips-dev-8.17
    $extractedDir = Get-ChildItem "third_party" -Directory | Where-Object { $_.Name -like "vips-dev-*" } | Select-Object -First 1
    
    if ($extractedDir) {
        if (Test-Path $DEST_DIR) {
            Remove-Item $DEST_DIR -Recurse -Force
        }
        Rename-Item $extractedDir.FullName "libvips"
        Write-Host "OK Installation completed ($ARCHIVE_NAME kept for reuse)" -ForegroundColor Green
    } else {
        throw "Failed to find extracted vips directory"
    }
    
    # Set environment variables for current session
    $vipsPath = Resolve-Path $DEST_DIR
    $env:PATH = "$vipsPath\bin;$env:PATH"
    $env:PKG_CONFIG_PATH = "$vipsPath\lib\pkgconfig"
    $env:CGO_ENABLED = "1"
    
    Write-Host ""
    Write-Host "OK Environment variables set (current session)" -ForegroundColor Green
    Write-Host "  PATH += $vipsPath\bin" -ForegroundColor Cyan
    Write-Host "  PKG_CONFIG_PATH = $vipsPath\lib\pkgconfig" -ForegroundColor Cyan
    Write-Host "  CGO_ENABLED = 1" -ForegroundColor Cyan
    Write-Host ""
} catch {
    Write-Host "ERROR: $_" -ForegroundColor Red
    return
}
