package pkg

import (
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// MaskMode represents the field mask operation mode
type MaskMode int

const (
	// FILTER mode keeps only the fields specified in the mask
	FILTER MaskMode = 0
	// PRUNE mode removes the fields specified in the mask
	PRUNE MaskMode = 1
)

// FieldMaskUtil provides utilities for working with FieldMask
type FieldMaskUtil struct {
	mode MaskMode
	mask *fieldmaskpb.FieldMask
}

// NewFieldMaskUtil creates a new FieldMaskUtil
func NewFieldMaskUtil(mode MaskMode, mask *fieldmaskpb.FieldMask) *FieldMaskUtil {
	return &FieldMaskUtil{
		mode: mode,
		mask: mask,
	}
}

// GetMode returns the mask mode
func (u *FieldMaskUtil) GetMode() MaskMode {
	return u.mode
}

// GetMask returns the underlying field mask
func (u *FieldMaskUtil) GetMask() *fieldmaskpb.FieldMask {
	return u.mask
}

// IsFieldMasked checks if a field is masked
func (u *FieldMaskUtil) IsFieldMasked(fieldPath string) bool {
	if u.mask == nil {
		return false
	}

	for _, path := range u.mask.Paths {
		if path == fieldPath {
			return true
		}
	}

	return false
}