# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`protoc-gen-fieldmask` is a protoc plugin that generates utility code for working with Google's FieldMask protobuf type. It helps developers avoid repetitive code when dealing with FieldMask messages in Go applications.

## Development Commands

### Build and Install
```bash
# Install the plugin locally
make install

# Build the plugin binary
go build -o bin/protoc-gen-fieldmask ./cmd/protoc-gen-fieldmask
```

### Testing
```bash
# Run all tests
make test
# or
go test -v ./... --count=1

# Run specific test file
go test -v ./module/fieldmask_test.go
```

### Code Generation
```bash
# Generate fieldmask protobuf definitions
make gen-fm-pb

# Generate code from proto files
protoc \
    -I. \
    -I$YOUR_PROTO_PATH \
    --go_out=paths=source_relative:. \
    --fieldmask_out=paths=source_relative,lang=go:. \
    your.proto
```

### Debugging
```bash
# Prepare debug data (requires protoc-gen-debug)
make prepare-debug
# or manually:
protoc \
    -I=./examples/pb \
    -I=./proto \
    --plugin=protoc-gen-debug=$(which protoc-gen-debug) \
    --debug_out="./internal/module/debugdata:." \
    ./examples/pb/user.proto
```

## Architecture

### Entry Point
- `cmd/protoc-gen-fieldmask/main.go`: Plugin entry point using protoc-gen-star framework
  - Registers the fieldmask module
  - Applies Go formatting as post-processor

### Core Generation Logic (`module/`)
- `fieldmask.go`: Main module implementation
  - Implements `pgs.Module` interface
  - Parses protobuf services and methods
  - Orchestrates template-based code generation
- `fm_message_out.go`: Data structures for fieldmask context
  - `fmMessagePair`: Associates RPC method with request/response messages
  - `checkMethodOptions()`: Validates method-level fieldmask options
  - `findFieldMaskField()`: Locates FieldMask field in request message
- `cache.go`: Message cache for cross-package message lookups

### Template System (`templates/`)
- `template.go`: Template registry factory for Go templates
- `go/`: Go-specific templates (embedded via go:embed)
  - `file.tpl`: Main file structure with package declaration
  - `request_mask.tpl`: Request field masking methods (`Mask_$field`)
  - `response_mask.tpl`: Response filtering methods (`FieldMask_Filter`, `FieldMask_Prune`)
  - `marked_checker.tpl`: Field checking utilities (`Masked_$field`)
- `shared/functions.go`: Template helper functions

### Protobuf Definitions
- `third_party/protoc_gen_fieldmask/option.proto`: Custom protobuf options for fieldmask configuration
- `protobuf/`: Generated protobuf Go files

## Key Concepts

### FieldMask Options
The plugin uses custom protobuf options defined in `third_party/protoc_gen_fieldmask/option.proto`:
- **Method-level options** (`(fieldmask.option.Option).in`/`.out`): Configure which fields generate mask methods
- **Field-level options**: Fine-grained control over individual field masking behavior

### Code Generation Flow
1. Plugin invoked via protoc with `--fieldmask_out` parameter
2. `main.go` initializes protoc-gen-star framework
3. `fieldmask.go` parses services/methods with fieldmask options
4. Template engine generates `.pb.fm.go` files
5. Go formatter post-processes generated code

### Masking Modes
- **FILTER**: Only masked fields are included in response
- **PRUNE**: Masked fields are removed from response

## Testing Strategy

- Unit tests in `module/` use `testify` framework
- Template tests in `templates/registry_test.go`
- Integration tests in `examples/` demonstrate real usage
- Debug tests use `protoc-gen-debug` output

## Important Notes

- Proto files with fieldmask options must specify the `lang=go` parameter
- The plugin requires `protoc-gen-star` framework for protobuf parsing
- Generated files use `.pb.fm.go` extension to distinguish from standard `.pb.go` files
- Cross-package message references require proper import path resolution
