#!/bin/bash
set -e

echo "Building hellogsm-go-score-calculator Lambda function..."

# Go 의존성 정리
go mod tidy

# Linux용 바이너리 빌드
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
else
    echo "Build failed!"
    exit 1
fi

echo "Creating ZIP package..."
zip -r function.zip bootstrap

if [ $? -eq 0 ]; then
    echo "ZIP package created successfully!"
    echo "File size: $(ls -lh function.zip | awk '{print $5}')"
    echo "Ready to upload to AWS Lambda"
else
    echo "Failed to create ZIP package!"
    exit 1
fi
