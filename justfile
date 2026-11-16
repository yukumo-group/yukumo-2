# Cross-platform command hub. All Go commands run in module mode so each
# module must resolve through its own go.mod + replace directives.
export GOWORK := "off"

default: check

# Local full check.
check: check-replaces lint vet test

# CI entry point: replace check, then lint before build, then tests/vet.
ci: check-replaces build test vet

# Build runs depguard + staticcheck first.
build: lint
    go -C cli build -o ../yukumo-script ./cmd/yukumo-script

run *ARGS:
    go -C cli run ./cmd/yukumo-script {{ARGS}}

install-tools:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

test:
    go -C core test ./...
    go -C runtime test ./...
    go -C core/api test ./...
    go -C cli test ./...

vet:
    go -C core vet ./...
    go -C runtime vet ./...
    go -C core/api vet ./...
    go -C cli vet ./...

lint:
    cd core && golangci-lint run
    cd runtime && golangci-lint run
    cd core/api && golangci-lint run
    cd cli && golangci-lint run

check-replaces:
    go run scripts/check-replaces.go
