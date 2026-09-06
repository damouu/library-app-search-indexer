#!/bin/bash

set -e

go test \
  ./internal/application \
  ./internal/health \
  ./internal/kafka \
  ./internal/mapper \
  -coverprofile=coverage.out \
  -covermode=atomic

go tool cover -func=coverage.out