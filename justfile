# Cross-platform command hub. All Go commands run in module mode so each
# module must resolve through its own go.mod + replace directives.
export GOWORK := "off"

default: check

# Local full check.
check: check-replaces check-cli-imports lint build-modules vet test

# CI entry point: replace check, then lint before build, then tests/vet.
ci: check-replaces check-cli-imports build build-modules test vet

# Build runs depguard + staticcheck first.
build: lint
    go -C cli build -o ../yukumo-script ./cmd/yukumo-script

run *ARGS:
    go -C cli run ./cmd/yukumo-script {{ARGS}}

build-modules:
    go -C core/port build ./...
    go -C core build ./...
    go -C runtime build ./...
    go -C core/api build ./...
    go -C cli build ./...

install-tools:
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

test:
    go -C core/port test ./...
    go -C core test ./...
    go -C runtime test ./...
    go -C core/api test ./...
    go -C cli test ./...

vet:
    go -C core/port vet ./...
    go -C core vet ./...
    go -C runtime vet ./...
    go -C core/api vet ./...
    go -C cli vet ./...

lint:
    cd core/port && golangci-lint run
    cd core && golangci-lint run
    cd runtime && golangci-lint run
    cd core/api && golangci-lint run
    cd cli && golangci-lint run

check-replaces:
    go run scripts/check-replaces.go

check-cli-imports:
    go run scripts/check-cli-imports.go

# wasm host smoke: copy the toolchain wasm_exec.js, build GOOS=js, check the
# dependency closure, then run the Node harness.
wasm-prepare:
    go run scripts/prepare-wasm.go

wasm-smoke: wasm-prepare
    GOOS=js GOARCH=wasm go -C wasmtest build -o .build/smoke.wasm ./...
    go run scripts/check-wasm-deps.go
    node wasmtest/harness.cjs

wasm-lint:
    cd wasmtest && GOOS=js GOARCH=wasm golangci-lint run
