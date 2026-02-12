# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

protoc-gen-fieldmask is a protoc plugin that generates FieldMask utilities for Go. It reduces boilerplate when handling `google.protobuf.FieldMask` in gRPC applications, enabling two main use cases:
1. **Masking gRPC response fields** - filter or prune fields based on client request
2. **Incremental updates** - selectively update and process only requested fields

## Development Commands

### Build & Install
```bash
# Install the plugin (builds to $GOPATH/bin)
make install

# Or manually
go install ./
```

### Testing
```bash
# Run all tests
make test

# Run tests directly
go test -v ./... --count=1

# Run specific test
go test -v ./internal/module/... -run TestName
```

### Code Generation (Debugging)
```bash
# Generate protobuf files for the plugin's own proto definitions
make gen-fm-pb

# Generate example protobuf files
cd examples && make gen-pb

# Prepare debug data (requires protoc-gen-debug)
make prepare-debug
```

### Running the Plugin
```bash
# Basic usage
protoc \
  -I. \
  -I$PROTO_PATH \
  --go_out=paths=source_relative:. \
  --fieldmask_out=paths=source_relative,lang=go:. \
  file.proto
```

The `lang=go` parameter is **required**. Currently only Go is supported.

## Architecture

### Plugin Framework
Built on `protoc-gen-star` (pgs), following standard protoc plugin pattern:
- Entry point: `main.go` - registers the FieldMask module
- Core logic: `internal/module/fieldmask.go` - main module implementation
- Templates: `internal/templates/` - code generation via template files

### Execution Pipeline
The `Execute` method in `fieldmask.go` processes files in four phases:

1. **Parse Phase** (`parse()`): Scans protobuf files for FieldMask fields with custom options
	- Checks for `google.protobuf.FieldMask` fields with `(fieldmask.option.Option)` extension

2. **Locate Phase** (`locateMessage()`): Finds associated in/out messages
	- First checks current file for message
	- Then searches imported packages (with caching via `pkgMessageCache`)
	- Resolves Go package names and import paths

3. **Consummate Phase** (`consummate()`): Resolves import paths and package information
	- Deduplicates out message variables (fixes #8)
	- Handles package name conflicts with suffixing

4. **Generate Phase** (`generate()`): Produces code using templates
	- Loads language-specific templates from registry
	- Generates `.pb.fm.go` files alongside standard `.pb.go` files

### Template System
- `internal/templates/registry.go` - template registry supporting multiple languages
- `internal/templates/go/` - Go-specific templates:
	- `fm.in.tpl` - input message utilities (MaskIn_* methods)
	- `fm.out.tpl` - output message utilities (MaskOut_*, MaskedOut_* methods)
	- `message.tpl` - common message generation

### Custom Proto Extension
`proto/fieldmask/option.proto` defines the FieldMask options:
- `in.gen`: generate input field mask utilities
- `out.gen`: generate output field mask utilities
- `out.message`: specify associated output message (required for out utilities)

## Key Constraints & Limitations

1. **Language**: Currently only Go (multi-language support is TODO)
2. **Message location**: In and out messages must be in the same proto file
3. **FieldMask type**: Only supports `google.protobuf.FieldMask`
4. **Package resolution**: Go `go_package` option must be correctly specified

## Debugging

For debugging the plugin:
```bash
# 1. Install protoc-gen-debug
go install github.com/lyft/protoc-gen-star/protoc-gen-debug@latest

# 2. Prepare debug data
make prepare-debug

# 3. Run debug test
go test -v ./internal/module/... -run Test_ForDebug
```

The debug data is generated to `internal/module/debugdata/` and allows inspecting the parsed protobuf AST.

## Generated Code Patterns

For an input message with FieldMask field:
- `MaskIn_Field()` - adds field to mask (for incremental updates)
- `MaskOut_Field()` - adds field to output mask (for response filtering)
- `FieldMask_Filter()` - returns filter that keeps only masked fields
- `FieldMask_Prune()` - returns filter that removes masked fields
- `MaskedIn_Field()`, `MaskedOut_Field()` - check if field is masked

## Module Structure

```
internal/
├── module/              # Core plugin logic
│   ├── fieldmask.go    # Main module and pipeline
│   ├── fm_message_in.go   # Input message processing
│   ├── fm_message_out.go  # Output message processing
│   └── cache.go        # Message caching for imports
└── templates/          # Code generation templates
    ├── registry.go     # Template registry
    ├── shared/         # Shared template functions
    └── go/             # Go templates
```