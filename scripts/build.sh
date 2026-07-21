#!/bin/bash
set -e

APP=pokedexcli
VERSION="v1.0.0"
PKG=./cmd/pokedex

mkdir -p dist

echo "Building ${APP} ${VERSION}.."

GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=${VERSION}" -o dist/${APP}-darwin-arm64 $PKG
GOOS=linux   GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o dist/${APP}-linux-amd64 $PKG
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o dist/${APP}-windows-amd64.exe $PKG

echo "Done. Binaries in dist/"
ls -la dist/