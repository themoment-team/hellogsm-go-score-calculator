#!/usr/bin/env pwsh
# build.ps1 - Windows PowerShell build script for hellogsm-go-score-calculator

Write-Host "Building hellogsm-go-score-calculator Lambda function..." -ForegroundColor Green

Write-Host "Running go mod tidy..." -ForegroundColor Yellow
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "go mod tidy failed!" -ForegroundColor Red
    exit 1
}

# 빌드 환경 변수 설정
Write-Host "Setting build environment..." -ForegroundColor Yellow
$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = "amd64"

Write-Host "Building binary..." -ForegroundColor Yellow
go build -o bootstrap main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!" -ForegroundColor Green
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Creating ZIP package..." -ForegroundColor Yellow
if (Test-Path "function.zip") {
    Remove-Item "function.zip" -Force
    Write-Host "Removed existing function.zip" -ForegroundColor Yellow
}

try {
    Compress-Archive -Path "bootstrap" -DestinationPath "function.zip" -CompressionLevel Optimal

    if (Test-Path "function.zip") {
        $size = (Get-Item "function.zip").Length
        $sizeKB = [math]::Round($size / 1KB, 2)
        Write-Host "ZIP package created successfully!" -ForegroundColor Green
        Write-Host "File size: $sizeKB KB" -ForegroundColor Cyan
        Write-Host "Ready to upload to AWS Lambda" -ForegroundColor Green
        Write-Host "`nBuild artifacts:" -ForegroundColor Cyan
        Get-ChildItem -Path "bootstrap", "function.zip" | Format-Table Name, Length, LastWriteTime -AutoSize
    } else {
        Write-Host "Failed to create ZIP package!" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Error creating ZIP package: $_" -ForegroundColor Red
    exit 1
}

Write-Host "`nBuild completed successfully!" -ForegroundColor Green
