#!/bin/bash

set -exuo pipefail

root=$PWD

directories=$(find "$root/" -type f -name '*.go' -exec dirname {} \; | sort -u)

for dir in $directories; do
    golangci-lint -c .golangci.yml run "$dir"
done
