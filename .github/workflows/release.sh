#!/bin/bash

set -exuo pipefail

export GOOS="$1"
export GOARCH="$2"
export CGO_ENABLED=0

rm -rf _tmp
mkdir -p _tmp

go build -v -o _tmp/ ./observer ./scanners/...

mv _tmp/observer "_publish/observer-${GOOS}-${GOARCH}"
for i in _tmp/*; do
	mv "$i" "_publish/$(basename "$i")-scanner-${GOOS}-${GOARCH}"
done
