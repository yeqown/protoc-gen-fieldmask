# ============================================================================
# protoc-gen-fieldmask Makefile
# ============================================================================
# A protoc plugin that generates utility code for Google's FieldMask

.PHONY: help install build test clean gen-fm-pb gen-third-party gen-tests gen-examples gen-all prepare-debug ci all

# Default target
.DEFAULT_GOAL := help

# ============================================================================
# Development Commands
# ============================================================================

## install: Install the protoc-gen-fieldmask plugin locally
install:
	@echo "==> Installing protoc-gen-fieldmask plugin..."
	go install ./cmd/protoc-gen-fieldmask

## build: Build the plugin binary
build:
	@echo "==> Building protoc-gen-fieldmask plugin..."
	@mkdir -p bin
	go build -o bin/protoc-gen-fieldmask ./cmd/protoc-gen-fieldmask

## test: Run all tests with verbose output
test:
	@echo "==> Running tests..."
	go test -v ./... --count=1

## test-short: Run short tests only
test-short:
	@echo "==> Running short tests..."
	go test -v ./... --count=1 --short

## test-cover: Run tests with coverage report
test-cover:
	@echo "==> Running tests with coverage..."
	go test -v ./... --coverprofile=coverage.out --count=1
	go tool cover -html=coverage.out -o coverage.html
	@echo "==> Coverage report: coverage.html"

## clean: Remove generated files and build artifacts
clean:
	@echo "==> Cleaning generated files..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f tests/*.pb.fm.go
	rm -f examples/proto/*.pb.go examples/proto/*.pb.fm.go
	@echo "==> Clean complete"

# ============================================================================
# Code Generation Commands
# ============================================================================

## gen-fm-pb: Generate fieldmask protobuf definitions (protobuf/option.pb.go)
gen-fm-pb:
	@echo "==> Generating fieldmask protobuf definitions..."
	protoc \
		-I=./third_party \
		--go_opt=paths=import \
		--go_out=. \
		./third_party/protoc_gen_fieldmask/option.proto
	@if [ -d "github.com/yeqown/protoc-gen-fieldmask/protobuf" ]; then \
		mv github.com/yeqown/protoc-gen-fieldmask/protobuf/option.pb.go protobuf/option.pb.go; \
		rm -rf github.com; \
	fi
	@echo "==> Generated protobuf/option.pb.go"

## gen-third-party: Generate third_party test proto files
gen-third-party:
	@echo "==> Generating third_party test proto files..."
	protoc \
		-I=. \
		--go_out=paths=source_relative:. \
		./third_party/test/common.proto
	@echo "==> Generated third_party/test/common.pb.go"

## gen-tests: Generate integration test code from proto files
gen-tests:
	@echo "==> Generating integration test code..."
	protoc \
		-I=. \
		-I=./third_party \
		--go_out=paths=source_relative:. \
		--fieldmask_out=paths=source_relative,lang=go:. \
		./tests/example.proto
	@echo "==> Generated tests/example.pb.fm.go"

## gen-examples: Generate example code from proto files
gen-examples:
	@echo "==> Generating example code..."
	protoc \
		-I=. \
		-I=./examples \
		-I=./third_party \
		--go_out=paths=source_relative:. \
		--fieldmask_out=paths=source_relative,lang=go:. \
		./examples/proto/user.proto
	@echo "==> Generated examples/proto/user.pb.fm.go"

## gen-all: Generate all proto files (protobuf definitions, third_party, tests, and examples)
gen-all: gen-fm-pb gen-third-party gen-tests gen-examples
	@echo "==> All proto files generated"

## prepare-debug: Prepare debug data using protoc-gen-debug (requires protoc-gen-debug)
prepare-debug:
	@echo "==> Preparing debug data..."
	@if [ -z "$(shell which protoc-gen-debug)" ]; then \
		echo "Error: protoc-gen-debug not found. Install it from: github.com/chronos-tachyon/go-protoc-gen-debug"; \
		exit 1; \
	fi
	@mkdir -p internal/module/debugdata
	protoc \
		-I=./examples/pb \
		-I=./proto \
		--plugin=protoc-gen-debug=$(shell which protoc-gen-debug) \
		--debug_out="./internal/module/debugdata:." \
		./examples/pb/user.proto
	@echo "==> Debug data generated in internal/module/debugdata"

# ============================================================================
# CI/CD Commands
# ============================================================================

## ci: Run all CI checks
ci: test
	@echo "==> All CI checks passed"

## all: Run tests
all: test
	@echo "==> Build and test completed"

# ============================================================================
# Help
# ============================================================================

## help: Show this help message
help:
	@echo ""
	@echo "protoc-gen-fieldmask - A protoc plugin for Google's FieldMask"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Development Commands:"
	@echo "  install         Install the protoc-gen-fieldmask plugin locally"
	@echo "  build           Build the plugin binary to bin/"
	@echo "  test            Run all tests"
	@echo "  test-short      Run short tests only"
	@echo "  test-cover      Run tests with coverage report"
	@echo "  clean           Remove generated files and build artifacts"
	@echo ""
	@echo "Code Generation:"
	@echo "  gen-fm-pb       Generate fieldmask protobuf definitions"
	@echo "  gen-third-party Generate third_party test proto files"
	@echo "  gen-tests       Generate integration test code from proto files"
	@echo "  gen-examples    Generate example code from proto files"
	@echo "  gen-all         Generate all proto files"
	@echo "  prepare-debug   Prepare debug data using protoc-gen-debug"
	@echo ""
	@echo "CI/CD:"
	@echo "  ci              Run all CI checks"
	@echo "  all             Run tests"
	@echo ""
