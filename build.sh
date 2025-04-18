#!/bin/bash
if command -v go &> /dev/null; then
    echo "Go 环境已安装"
    go version
else
    echo "Go 环境未安装"
    exit 1
fi

go mod download
go mod verify
export CGO_ENABLED=0

PLATFORMS="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 linux/loong64 linux/mips64 linux/mips64le linux/ppc64 linux/ppc64le linux/riscv64 linux/s390x aix/ppc64 dragonfly/amd64 netbsd/amd64 netbsd/arm64 windows/amd64 windows/arm64"
OUT_DIR=./bin
MAIN_PATH=./cmd/BaPs/BaPs.go
NAME="BaPs"

  for platform in $PLATFORMS; do
      export GOOS=$(echo $platform | cut -d'/' -f1)
      export GOARCH=$(echo $platform | cut -d'/' -f2)
      OUTPUT_NAME=$NAME"_"$GOOS"_"$GOARCH
      if [ $GOOS = "windows" ]; then
        OUTPUT_NAME="$OUTPUT_NAME.exe"
      fi
      echo "Building $OUTPUT_NAME..."
      go build -ldflags="-s -w" -o $OUT_DIR/$OUTPUT_NAME $MAIN_PATH
    done