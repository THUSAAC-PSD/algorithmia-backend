#!/bin/bash

set -e

go get -u -t -v ./... && go mod tidy