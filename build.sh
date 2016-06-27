#!/bin/bash
set -euo pipefail

version=$(git describe --always)

build() {
    name="smartos-pxe-$GOOS-$GOARCH-$version"
    rm -rf "$name"

    mkdir -p "$name/bin"
    go build -i -v -ldflags "-w -s -X main.version=$version" -o "$name/bin/smartos-pxe"

    cp -r data "$name/data"
    echo "$name" > "$name/data/BUILD"
    tar zcvf "$name.tar.gz" "$name"

    rm -rf  "$name"
}

GOOS=linux GOARCH=arm build
GOOS=linux GOARCH=amd64 build
GOOS=darwin GOARCH=amd64 build
GOOS=solaris GOARCH=amd64 build
GOOS=freebsd GOARCH=amd64 build
