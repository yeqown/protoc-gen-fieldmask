# protoc-gen-fieldmask Examples

This directory contains examples demonstrating the usage of `protoc-gen-fieldmask`, organized from basic to advanced.

## Directory Structure

```
examples/
├── 01_basic/          # Basic usage
├── 02_intermediate/    # Intermediate features
├── 03_advanced/        # Advanced scenarios
├── README.md          # This file
```

## 01_basic: Basic Usage

**Concept**: Getting started with FieldMask for simple response filtering.

**What you'll learn**:
- How to add FieldMask to your proto definitions
- Basic FILTER mode usage
- How to mark fields for inclusion in response
- How to check if fields are marked

**Key files**:
- `basic.proto` - Simple user service with FieldMask
- `example_test.go` - Demonstrates basic FieldMask usage

**Run**:
```bash
cd 01_basic
go test -v
```

## 02_intermediate: Intermediate Features

**Concept**: Learning FILTER vs PRUNE modes and field-level options.

**What you'll learn**:
- Difference between FILTER and PRUNE modes
- Field-level `ignore` option
- Field-level `nested` option for nested messages
- How to handle nested message field masking

**Key files**:
- `intermediate.proto` - User update service with advanced options
- `example_test.go` - Demonstrates intermediate features

**Run**:
```bash
cd 02_intermediate
go test -v
```

## 03_advanced: Advanced Scenarios

**Concept**: Real-world e-commerce scenarios with complex FieldMask usage.

**What you'll learn**:
- Multiple RPCs with different modes in same service
- Complex response filtering for performance
- Field-level options for fine-grained control
- Best practices for production use

**Key files**:
- `advanced.proto` - Order service with multiple operations
- `example_test.go` - Demonstrates advanced patterns

**Run**:
```bash
cd 03_advanced
go test -v
```

## Learning Path

1. **Start with `01_basic/`** - Understand the core concepts
2. **Move to `02_intermediate/`** - Learn modes and options
3. **Explore `03_advanced/`** - See production-ready patterns

## Proto File Generation

To generate Go code from proto files:

```bash
# From examples directory
protoc \
  -I. \
  -I../proto \
  --go_out=paths=source_relative:. \
  --fieldmask_out=paths=source_relative,lang=go:. \
  01_basic/basic.proto

# Or use make
make gen-pb
```

## Quick Reference

### Modes
- **FILTER**: Only masked fields are included (for queries)
- **PRUNE**: Masked fields are excluded (for updates)

### Field Options
- **ignore**: Don't generate mask methods for this field
- **nested**: Generate nested FieldMask for message fields

### API Pattern
```go
// Create FieldMask
fm := req.FieldMask(fieldmask.FILTER)

// Mark fields
fm.Response().Field1()
fm.Response().Field2()

// Check if marked
if fm.Marked().Response().Field1() {
    // Field is marked
}
```

## Notes

- Examples use Go 1.25+
- Requires protoc and protoc-gen-fieldmask installed
- Generated code is in `.pb.fm.go` files
- FieldMask paths use dot-notation: `"user_id"`, `"address.country"`
