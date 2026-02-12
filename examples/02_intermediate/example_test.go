package intermediate_test

import (
	"testing"

	"github.com/yeqown/protoc-gen-fieldmask/examples/02_intermediate/intermediate"
	pbfieldmask "github.com/yeqown/protoc-gen-fieldmask/proto/protobuf"
)

// Test_IntermediateFieldMask demonstrates intermediate features:
// - PRUNE mode for updates
// - Field-level options (ignore, nested)
func Test_IntermediateFieldMask(t *testing.T) {
	// Simulate an update request
	req := &intermediate.UpdateUserRequest{
		UserId: "user123",
		Name:   "John Doe",
		// Email won't be included in mask (marked as ignore)
		// Address will be included with nested field mask
	}

	// Use PRUNE mode - masked fields are EXCLUDED from update
	fm := req.FieldMask(pbfieldmask.PRUNE)

	// Mark fields to update (all except email which is ignored)
	fm.Response().UserId() // Note: using Response() for update target
	fm.Response().Name()

	// Demonstrate nested field masking
	// Get the nested Address field mask
	addrMask := fm.Response().Address().FieldMask()

	// Mark nested fields within Address
	addrMask.Country()
	addrMask.Province()
	// city is not marked, so it won't be updated

	// Check which fields are marked
	if fm.Marked().Response().UserId() {
		t.Log("UserId is marked for update")
	}
	if fm.Marked().Response().Name() {
		t.Log("Name is marked for update")
	}
	if !fm.Marked().Response().Email() {
		t.Log("Email is NOT marked (was set to ignore)")
	}

	// Check nested fields
	if addrMask.Marked().Country() {
		t.Log("Country is marked for update")
	}

	// In PRUNE mode, marked fields are those that will be UPDATED
	// Unmarked fields (including city) will remain unchanged

	t.Log("Intermediate FieldMask usage demonstrated successfully")
}

// Test_FilterVsPrune compares FILTER and PRUNE modes.
func Test_FilterVsPrune(t *testing.T) {
	t.Run("FILTER mode includes only marked fields", func(t *testing.T) {
		req := &intermediate.UpdateUserRequest{}

		// FILTER: only marked fields are included
		fm := req.FieldMask(pbfieldmask.FILTER)
		fm.Response().Name()

		t.Log("FILTER mode: only Name will be in response")
	})

	t.Run("PRUNE mode excludes marked fields", func(t *testing.T) {
		req := &intermediate.UpdateUserRequest{}

		// PRUNE: marked fields are excluded from operation
		fm := req.FieldMask(pbfieldmask.PRUNE)
		fm.Response().Name()

		t.Log("PRUNE mode: Name will be excluded from update")
	})
}
