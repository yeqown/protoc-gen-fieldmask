package protobuf

import (
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

const (
	// PathPrefixRequest is the prefix for request field paths
	PathPrefixRequest = "req."
	// PathPrefixResponse is the prefix for response field paths
	PathPrefixResponse = "res."
)

// SplitPaths splits the FieldMask paths into request and response maps.
// Paths with "req." prefix go to request map, "res." prefix go to response map.
// Returns (requestPaths, responsePaths)
func SplitPaths(paths []string) (map[string]struct{}, map[string]struct{}) {
	req := make(map[string]struct{})
	res := make(map[string]struct{})

	for _, path := range paths {
		if strings.HasPrefix(path, PathPrefixRequest) {
			// Remove prefix and store
			fieldPath := strings.TrimPrefix(path, PathPrefixRequest)
			req[fieldPath] = struct{}{}
		} else if strings.HasPrefix(path, PathPrefixResponse) {
			// Remove prefix and store
			fieldPath := strings.TrimPrefix(path, PathPrefixResponse)
			res[fieldPath] = struct{}{}
		}
		// Ignore paths without prefix
	}

	return req, res
}

// parseNestedMasks builds a map of top-level fields to their nested field masks
// Returns (topLevelFields, hasNested)
// topLevelFields[fieldName] = map of nested paths for that field
// hasNested[fieldName] = true if the field has nested masks
func parseNestedMasks(mask map[string]struct{}) (map[string]map[string]struct{}, map[string]bool) {
	topLevelFields := make(map[string]map[string]struct{})
	nestedFields := make(map[string]bool)

	for path := range mask {
		parts := strings.Split(path, ".")
		if len(parts) > 1 {
			// This is a nested path like "address.city"
			topLevel := parts[0]
			if topLevelFields[topLevel] == nil {
				topLevelFields[topLevel] = make(map[string]struct{})
			}
			// Store the nested path without the top-level prefix
			nestedPath := strings.Join(parts[1:], ".")
			topLevelFields[topLevel][nestedPath] = struct{}{}
			nestedFields[topLevel] = true
		}
	}

	return topLevelFields, nestedFields
}

// Filter keeps only the fields specified in the mask (FILTER mode)
// All other fields are cleared (set to zero value)
// Supports nested field paths like "address.city"
func Filter(message proto.Message, mask map[string]struct{}) {
	applyMask(message, mask, true)
}

// Prune removes the fields specified in the mask (PRUNE mode)
// Only fields in the mask are cleared, all others remain
// Supports nested field paths like "address.city"
func Prune(message proto.Message, mask map[string]struct{}) {
	applyMask(message, mask, false)
}

// applyMask applies the field mask to the message
// If isFilter is true, keeps only masked fields (FILTER mode)
// If isFilter is false, removes masked fields (PRUNE mode)
// Supports nested field paths like "address.city" and "address.geo.latitude"
// Nested expansion is only supported for single message fields, not repeated or map
func applyMask(message proto.Message, mask map[string]struct{}, isFilter bool) {
	if message == nil {
		return
	}
	// For FILTER mode, empty mask means clear all fields (keep nothing)
	// For PRUNE mode, empty mask means keep all fields (clear nothing)
	// So only return early for PRUNE mode with empty mask
	if !isFilter && len(mask) == 0 {
		return
	}

	refl := message.ProtoReflect()
	if !refl.IsValid() {
		return
	}

	topLevelFields, _ := parseNestedMasks(mask)

	refl.Range(func(field protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		fieldName := string(field.Name())

		// Check if this field has nested masks (e.g., "address.city", "address.geo.latitude")
		if nestedMasks, hasNested := topLevelFields[fieldName]; hasNested {
			// Only single message fields support nested expansion
			// For repeated, map, and primitive types, nested masks are ignored
			if field.Message() != nil && !field.IsMap() && !field.IsList() {
				// Recursively process nested message
				nestedMsg := value.Message()
				if nestedMsg.IsValid() {
					applyMask(nestedMsg.Interface(), nestedMasks, isFilter)
				}
				return true
			}
			// For other types with nested paths, ignore the nested part
			// Fall through to check if field itself is masked
		}

		// For fields without nested masks (or unsupported nested types),
		// check if the field itself is masked
		_, isMasked := mask[fieldName]
		var shouldClear bool
		if isFilter {
			shouldClear = !isMasked // FILTER: clear if NOT masked
		} else {
			shouldClear = isMasked // PRUNE: clear if IS masked
		}

		if shouldClear {
			refl.Clear(field)
		}
		return true
	})
}

// AddReqPathToFieldMask adds a request field path with "req." prefix to the FieldMask.
func AddReqPathToFieldMask(fm *fieldmaskpb.FieldMask, path string) {
	if fm == nil {
		return
	}
	fm.Paths = append(fm.Paths, PathPrefixRequest+path)
}

// AddResPathToFieldMask adds a response field path with "res." prefix to the FieldMask.
func AddResPathToFieldMask(fm *fieldmaskpb.FieldMask, path string) {
	if fm == nil {
		return
	}
	fm.Paths = append(fm.Paths, PathPrefixResponse+path)
}

// IsFieldMarked checks if a field is in the mask (marked)
func IsFieldMarked(mask map[string]struct{}, field string) bool {
	_, ok := mask[field]
	return ok
}
