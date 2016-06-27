#!/bin/bash
set -euo pipefail

version=$(git describe --always)

build-tar() {
    name="smartos-pxe-$GOOS-$GOARCH-$version"
    rm -rf "$name"

    mkdir -p "$name/bin"
    go build -i -v -ldflags "-w -s -X main.version=$version" -o "$name/bin/smartos-pxe"

    cp -r data "$name/data"
    echo "$name" > "$name/data/BUILD"
    tar zcvf "$name.tar.gz" "$name"

    rm -rf  "$name"
}

build-zip() {
    name="smartos-pxe-$GOOS-$GOARCH-$version"
    rm -rf "$name"

    mkdir -p "$name/bin"
    go build -i -v -ldflags "-w -s -X main.version=$version" -o "$name/bin/smartos-pxe.exe"

    cp -r data "$name/data"
    echo "$name" > "$name/data/BUILD"
    zip -r "$name.zip" "$name"

    rm -rf  "$name"
}

GOOS=linux GOARCH=arm build-tar
GOOS=linux GOARCH=amd64 build-tar
GOOS=darwin GOARCH=amd64 build-tar
GOOS=solaris GOARCH=amd64 build-tar
GOOS=freebsd GOARCH=amd64 build-tar
GOOS=windows GOARCH=amd64 build-zip