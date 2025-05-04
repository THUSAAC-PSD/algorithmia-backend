#!/bin/bash

# https://blog.devgenius.io/go-golang-testing-tools-tips-to-step-up-your-game-4ed165a5b3b5
# https://github.com/testcontainers/testcontainers-go/pull/1394
# https://github.com/testcontainers/testcontainers-go/issues/1359

# Might be useful in the future

set -e

readonly type="$1"

go test -tags="$type" -timeout=30m  -count=1 -p=1 -parallel=1 ./...
