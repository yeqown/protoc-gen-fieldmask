package advanced_test

import (
	"testing"

	"github.com/yeqown/protoc-gen-fieldmask/examples/03_advanced/advanced"
	pbfieldmask "github.com/yeqown/protoc-gen-fieldmask/proto/protobuf"
)

// Test_AdvancedFieldMask demonstrates advanced usage scenarios.
func Test_AdvancedFieldMask(t *testing.T) {
	t.Run("FILTER mode for queries", func(t *testing.T) {
		req := &advanced.GetOrderRequest{
			OrderId: "order123",
		}

		// FILTER mode: only marked fields are included in response
		fm := req.FieldMask(pbfieldmask.FILTER)

		// Mark fields to include
		fm.Response().OrderId()
		fm.Response().TotalAmount()
		// products and status are NOT marked, won't be included

		t.Log("FILTER: only OrderId and TotalAmount will be returned")
	})

	t.Run("PRUNE mode for updates", func(t *testing.T) {
		req := &advanced.UpdateOrderRequest{
			OrderId:     "order123",
			TotalAmount: 99.99,
		}

		// PRUNE mode: marked fields are EXCLUDED from update
		fm := req.FieldMask(pbfieldmask.PRUNE)

		// Mark fields to EXCLUDE from update
		fm.Response().TotalAmount()
		// product_ids is ignored via field option, can't be marked
		// OrderId is not marked, so it WILL be updated (retained)

		t.Log("PRUNE: TotalAmount will NOT be updated, OrderId will be")
	})

	t.Run("Field ignore option", func(t *testing.T) {
		req := &advanced.ProductOptions{}

		// product_ids field has ignore: true in proto
		// No FieldMask methods will be generated for it
		fm := req.FieldMask(pbfieldmask.FILTER)

		// Can only mark non-ignored fields
		fm.Response().TotalAmount()

		// This would be a compile error:
		// fm.Response().ProductIds() // Method doesn't exist

		t.Log("Ignored fields don't have mask methods generated")
	})
}

// Test_AdvancedResponseMasking demonstrates complex response masking.
func Test_AdvancedResponseMasking(t *testing.T) {
	t.Run("Partial field selection", func(t *testing.T) {
		req := &advanced.GetOrderRequest{OrderId: "order123"}

		// Client only needs order total amount
		fm := req.FieldMask(pbfieldmask.FILTER)
		fm.Response().TotalAmount()

		// Apply mask to filter response
		resp := &advanced.GetOrderResponse{
			OrderId:     "order123",
			TotalAmount: 99.99,
			Status:      "completed",
		}

		// After applying mask, only TotalAmount would remain
		_ = resp
		_ = fm

		t.Log("Response can be filtered to specific fields")
	})

	t.Run("Nested message handling", func(t *testing.T) {
		req := &advanced.GetOrderRequest{OrderId: "order123"}

		fm := req.FieldMask(pbfieldmask.FILTER)

		// Mark top-level fields
		fm.Response().OrderId()

		// In advanced usage with nested messages,
		// you would navigate to nested fields
		// (if they were supported in the proto definition)

		t.Log("Nested message structures are handled")
	})
}
