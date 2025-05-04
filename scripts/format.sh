#!/bin/bash

# ref: https://blog.devgenius.io/sort-go-imports-acb76224dfa7
# https://yolken.net/blog/cleaner-go-code-golines

set -e

# https://github.com/segmentio/golines
# # will do `gofmt` internally
golines -m 120 -w --ignore-generated .

# # https://pkg.go.dev/golang.org/x/tools/cmd/goimports
# goimports -l -w .

# https://github.com/incu6us/goimports-reviser
# https://github.com/incu6us/goimports-reviser/issues/118
# https://github.com/incu6us/goimports-reviser/issues/88
# https://github.com/incu6us/goimports-reviser/issues/104
gci write --skip-generated -s standard -s "prefix(github.com/THUSAAC-PSD/algorithmia-backend" -s default -s blank -s dot --custom-order  .

# https://golang.org/cmd/gofmt/
# gofmt -w .

# https://github.com/mvdan/gofumpt
# will do `gofmt` internally
gofumpt -l -w .
