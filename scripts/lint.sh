#!/bin/bash

# ref: https://freshman.tech/linting-golang/
set -e

# https://github.com/mgechev/revive
revive -config revive-config.toml -formatter friendly ./...

# https://github.com/dominikh/go-tools
staticcheck ./...

# https://golangci-lint.run/usage/linters/
# https://golangci-lint.run/usage/configuration/
# https://golangci-lint.run/usage/quick-start/
golangci-lint run ./...

# https://github.com/kisielk/errcheck
errcheck ./...
