#!/bin/bash

# Facebook Comment Scraper - Cross-Platform Build Script
# Generates executables for Windows, Linux, and macOS

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Build information
APP_NAME="facebook-comment-scraper"
VERSION="1.0.0"
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${BLUE}🚀 Building Facebook Comment Scraper${NC}"
echo -e "${BLUE}=====================================${NC}"
echo -e "Version: ${VERSION}"
echo -e "Build Time: ${BUILD_TIME}"
echo -e "Git Commit: ${GIT_COMMIT}"
echo ""

# Create builds directory
mkdir -p builds
cd builds

# Build flags
LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -s -w"

echo -e "${YELLOW}📦 Building for Windows (amd64)...${NC}"
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o "${APP_NAME}-windows-amd64.exe" ../main.go
echo -e "${GREEN}✅ Windows executable created: ${APP_NAME}-windows-amd64.exe${NC}"

echo -e "${YELLOW}📦 Building for Linux (amd64)...${NC}"
GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o "${APP_NAME}-linux-amd64" ../main.go
echo -e "${GREEN}✅ Linux executable created: ${APP_NAME}-linux-amd64${NC}"

echo -e "${YELLOW}📦 Building for macOS (amd64)...${NC}"
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o "${APP_NAME}-macos-amd64" ../main.go
echo -e "${GREEN}✅ macOS (Intel) executable created: ${APP_NAME}-macos-amd64${NC}"

echo -e "${YELLOW}📦 Building for macOS (arm64 - Apple Silicon)...${NC}"
GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o "${APP_NAME}-macos-arm64" ../main.go
echo -e "${GREEN}✅ macOS (Apple Silicon) executable created: ${APP_NAME}-macos-arm64${NC}"

echo ""
echo -e "${GREEN}🎉 All builds completed successfully!${NC}"
echo ""
echo -e "${BLUE}📁 Generated executables:${NC}"
ls -la *.exe *.macos-* *.linux-* 2>/dev/null || ls -la

echo ""
echo -e "${YELLOW}💡 Usage Instructions:${NC}"
echo -e "• Windows: Run ${APP_NAME}-windows-amd64.exe"
echo -e "• Linux: Run ./${APP_NAME}-linux-amd64"
echo -e "• macOS (Intel): Run ./${APP_NAME}-macos-amd64"
echo -e "• macOS (Apple Silicon): Run ./${APP_NAME}-macos-arm64"
echo ""
echo -e "${BLUE}🔒 Security Note:${NC}"
echo -e "These executables are clean and should pass antivirus scans."
echo -e "If flagged, it's likely a false positive due to the scraping nature."
