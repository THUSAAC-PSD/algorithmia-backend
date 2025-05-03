#!/bin/bash

set -e

# https://github.com/incu6us/goimports-reviser
go install -v github.com/incu6us/goimports-reviser/v3@latest

# https://github.com/daixiang0/gci
go install github.com/daixiang0/gci@latest

# https://pkg.go.dev/golang.org/x/tools/cmd/goimports
go install golang.org/x/tools/cmd/goimports@latest

# https://github.com/segmentio/golines
go install github.com/segmentio/golines@latest

# https://github.com/mvdan/gofumpt
go install mvdan.cc/gofumpt@latest

# https://github.com/dominikh/go-tools
go install honnef.co/go/tools/cmd/staticcheck@latest

go install github.com/mgechev/revive@latest

# https://golangci-lint.run/usage/install/#install-from-source
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

go install github.com/onsi/ginkgo/v2/ginkgo@latest