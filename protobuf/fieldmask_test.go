package protobuf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	common "github.com/yeqown/protoc-gen-fieldmask/third_party/test"
	"google.golang.org/protobuf/proto"
)

// TestSplitPaths tests the SplitPaths function
func TestSplitPaths(t *testing.T) {
	tests := []struct {
		name    string
		paths   []string
		wantReq map[string]struct{}
		wantRes map[string]struct{}
	}{
		{
			name:    "empty paths",
			paths:   []string{},
			wantReq: map[string]struct{}{},
			wantRes: map[string]struct{}{},
		},
		{
			name:  "only request paths",
			paths: []string{"req.user_id", "req.username"},
			wantReq: map[string]struct{}{
				"user_id":  {},
				"username": {},
			},
			wantRes: map[string]struct{}{},
		},
		{
			name:    "only response paths",
			paths:   []string{"res.email", "res.address"},
			wantReq: map[string]struct{}{},
			wantRes: map[string]struct{}{
				"email":   {},
				"address": {},
			},
		},
		{
			name:  "mixed request and response paths",
			paths: []string{"req.user_id", "res.email", "req.username", "res.address.city"},
			wantReq: map[string]struct{}{
				"user_id":  {},
				"username": {},
			},
			wantRes: map[string]struct{}{
				"email":        {},
				"address.city": {},
			},
		},
		{
			name:  "paths without prefix are ignored",
			paths: []string{"user_id", "req.user_id", "res.email", "invalid.prefix.field"},
			wantReq: map[string]struct{}{
				"user_id": {},
			},
			wantRes: map[string]struct{}{
				"email": {},
			},
		},
		{
			name: "nested paths",
			paths: []string{
				"req.address.street",
				"req.address.city",
				"res.address.geo.latitude",
				"res.address.geo.longitude",
			},
			wantReq: map[string]struct{}{
				"address.street": {},
				"address.city":   {},
			},
			wantRes: map[string]struct{}{
				"address.geo.latitude":  {},
				"address.geo.longitude": {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReq, gotRes := SplitPaths(tt.paths)

			// Check request paths
			assert.Equal(t, len(tt.wantReq), len(gotReq), "SplitPaths() req length mismatch")
			for path := range tt.wantReq {
				assert.Contains(t, gotReq, path, "SplitPaths() req missing path %q", path)
			}
			for path := range gotReq {
				assert.Contains(t, tt.wantReq, path, "SplitPaths() req extra path %q", path)
			}

			// Check response paths
			assert.Equal(t, len(tt.wantRes), len(gotRes), "SplitPaths() res length mismatch")
			for path := range tt.wantRes {
				assert.Contains(t, gotRes, path, "SplitPaths() res missing path %q", path)
			}
			for path := range gotRes {
				assert.Contains(t, tt.wantRes, path, "SplitPaths() res extra path %q", path)
			}
		})
	}
}

// TestFilter_CommonMessage tests Filter with CommonMessage
func TestFilter_CommonMessage(t *testing.T) {
	t.Run("filter keeps only masked fields", func(t *testing.T) {
		msg := &common.CommonMessage{
			Id:        "msg123",
			CreatedAt: "2024-01-01",
			UpdatedAt: "2024-01-02",
		}
		mask := map[string]struct{}{
			"id":         {},
			"created_at": {},
		}

		Filter(msg, mask)

		assert.Equal(t, "msg123", msg.Id, "expected Id 'msg123'")
		assert.Equal(t, "2024-01-01", msg.CreatedAt, "expected CreatedAt '2024-01-01'")
		assert.Empty(t, msg.UpdatedAt, "expected UpdatedAt to be empty")
	})

	t.Run("filter with empty mask clears all fields", func(t *testing.T) {
		msg := &common.CommonMessage{
			Id:        "msg123",
			CreatedAt: "2024-01-01",
			UpdatedAt: "2024-01-02",
		}
		Filter(msg, map[string]struct{}{})

		assert.Empty(t, msg.Id, "expected Id to be empty")
		assert.Empty(t, msg.CreatedAt, "expected CreatedAt to be empty")
		assert.Empty(t, msg.UpdatedAt, "expected UpdatedAt to be empty")
	})

	t.Run("filter with all fields mask keeps everything", func(t *testing.T) {
		msg := &common.CommonMessage{
			Id:        "msg123",
			CreatedAt: "2024-01-01",
			UpdatedAt: "2024-01-02",
		}
		mask := map[string]struct{}{
			"id":         {},
			"created_at": {},
			"updated_at": {},
		}
		Filter(msg, mask)

		assert.Equal(t, "msg123", msg.Id, "expected Id 'msg123'")
		assert.Equal(t, "2024-01-01", msg.CreatedAt, "expected CreatedAt '2024-01-01'")
		assert.Equal(t, "2024-01-02", msg.UpdatedAt, "expected UpdatedAt '2024-01-02'")
	})
}

// TestFilter_Address tests Filter with Address message (nested fields)
func TestFilter_Address(t *testing.T) {
	t.Run("filter with nested address fields", func(t *testing.T) {
		msg := &common.Address{
			Street:  "123 Main St",
			City:    "San Francisco",
			State:   "CA",
			Country: "USA",
			ZipCode: "94102",
			Geo: &common.GeoLocation{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		}
		mask := map[string]struct{}{
			"street":        {},
			"city":          {},
			"geo.latitude":  {},
			"geo.longitude": {},
		}

		Filter(msg, mask)

		assert.Equal(t, "123 Main St", msg.Street, "expected Street '123 Main St'")
		assert.Equal(t, "San Francisco", msg.City, "expected City 'San Francisco'")
		assert.Empty(t, msg.State, "expected State to be empty")
		assert.Empty(t, msg.Country, "expected Country to be empty")
		assert.Empty(t, msg.ZipCode, "expected ZipCode to be empty")
		assert.NotNil(t, msg.Geo, "expected Geo to not be nil")
		assert.Equal(t, 37.7749, msg.Geo.Latitude, "expected Latitude 37.7749")
		assert.Equal(t, -122.4194, msg.Geo.Longitude, "expected Longitude -122.4194")
	})
}

// TestFilter_ContactInfo tests Filter with ContactInfo
func TestFilter_ContactInfo(t *testing.T) {
	t.Run("filter ContactInfo", func(t *testing.T) {
		msg := &common.ContactInfo{
			Email:             "test@example.com",
			Phone:             "+1234567890",
			AlternativeEmails: []string{"alt1@example.com", "alt2@example.com"},
		}
		mask := map[string]struct{}{
			"email": {},
		}

		Filter(msg, mask)

		assert.Equal(t, "test@example.com", msg.Email, "expected Email 'test@example.com'")
		assert.Empty(t, msg.Phone, "expected Phone to be empty")
		assert.Nil(t, msg.AlternativeEmails, "expected AlternativeEmails to be nil")
	})
}

// TestPrune_CommonMessage tests Prune with CommonMessage
func TestPrune_CommonMessage(t *testing.T) {
	t.Run("prune removes only masked fields", func(t *testing.T) {
		msg := &common.CommonMessage{
			Id:        "msg123",
			CreatedAt: "2024-01-01",
			UpdatedAt: "2024-01-02",
		}
		mask := map[string]struct{}{
			"updated_at": {},
		}

		Prune(msg, mask)

		assert.Equal(t, "msg123", msg.Id, "expected Id 'msg123'")
		assert.Equal(t, "2024-01-01", msg.CreatedAt, "expected CreatedAt '2024-01-01'")
		assert.Empty(t, msg.UpdatedAt, "expected UpdatedAt to be empty (pruned)")
	})

	t.Run("prune with empty mask keeps everything", func(t *testing.T) {
		msg := &common.CommonMessage{
			Id:        "msg123",
			CreatedAt: "2024-01-01",
			UpdatedAt: "2024-01-02",
		}
		Prune(msg, map[string]struct{}{})

		assert.Equal(t, "msg123", msg.Id, "expected Id 'msg123'")
		assert.Equal(t, "2024-01-01", msg.CreatedAt, "expected CreatedAt '2024-01-01'")
		assert.Equal(t, "2024-01-02", msg.UpdatedAt, "expected UpdatedAt '2024-01-02'")
	})

	t.Run("prune with all fields mask clears everything", func(t *testing.T) {
		msg := &common.CommonMessage{
			Id:        "msg123",
			CreatedAt: "2024-01-01",
			UpdatedAt: "2024-01-02",
		}
		mask := map[string]struct{}{
			"id":         {},
			"created_at": {},
			"updated_at": {},
		}
		Prune(msg, mask)

		assert.Empty(t, msg.Id, "expected Id to be empty (pruned)")
		assert.Empty(t, msg.CreatedAt, "expected CreatedAt to be empty (pruned)")
		assert.Empty(t, msg.UpdatedAt, "expected UpdatedAt to be empty (pruned)")
	})
}

// TestPrune_Address tests Prune with Address message (nested fields)
func TestPrune_Address(t *testing.T) {
	t.Run("prune with nested address fields", func(t *testing.T) {
		msg := &common.Address{
			Street:  "123 Main St",
			City:    "San Francisco",
			State:   "CA",
			Country: "USA",
			ZipCode: "94102",
			Geo: &common.GeoLocation{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		}
		mask := map[string]struct{}{
			"state":        {},
			"country":      {},
			"zip_code":     {},
			"geo.latitude": {},
		}

		Prune(msg, mask)

		assert.Equal(t, "123 Main St", msg.Street, "expected Street '123 Main St'")
		assert.Equal(t, "San Francisco", msg.City, "expected City 'San Francisco'")
		assert.Empty(t, msg.State, "expected State to be empty (pruned)")
		assert.Empty(t, msg.Country, "expected Country to be empty (pruned)")
		assert.Empty(t, msg.ZipCode, "expected ZipCode to be empty (pruned)")
		assert.NotNil(t, msg.Geo, "expected Geo to not be nil")
		assert.Equal(t, float64(0), msg.Geo.Latitude, "expected Latitude to be 0 (pruned)")
		assert.Equal(t, -122.4194, msg.Geo.Longitude, "expected Longitude -122.4194")
	})
}

// TestFilter_NilMessage tests Filter with nil message
func TestFilter_NilMessage(t *testing.T) {
	t.Run("nil message", func(t *testing.T) {
		Filter(nil, map[string]struct{}{"id": {}})
		// Should not panic
	})
}

// TestPrune_NilMessage tests Prune with nil message
func TestPrune_NilMessage(t *testing.T) {
	t.Run("nil message", func(t *testing.T) {
		Prune(nil, map[string]struct{}{"id": {}})
		// Should not panic
	})
}

// TestFilter_EmptyMask tests Filter with empty mask
func TestFilter_EmptyMask(t *testing.T) {
	t.Run("empty mask", func(t *testing.T) {
		msg := &common.CommonMessage{
			Id:        "msg123",
			CreatedAt: "2024-01-01",
		}
		Filter(msg, map[string]struct{}{})

		assert.Empty(t, msg.Id, "expected Id to be empty with empty mask")
		assert.Empty(t, msg.CreatedAt, "expected CreatedAt to be empty with empty mask")
	})
}

// TestPrune_EmptyMask tests Prune with empty mask
func TestPrune_EmptyMask(t *testing.T) {
	t.Run("empty mask", func(t *testing.T) {
		msg := &common.CommonMessage{
			Id:        "msg123",
			CreatedAt: "2024-01-01",
		}
		Prune(msg, map[string]struct{}{})

		assert.Equal(t, "msg123", msg.Id, "expected Id 'msg123' with empty mask")
		assert.Equal(t, "2024-01-01", msg.CreatedAt, "expected CreatedAt '2024-01-01' with empty mask")
	})
}

// TestFilter_NestedFields tests Filter with deeply nested field paths
func TestFilter_NestedFields(t *testing.T) {
	t.Run("deeply nested geo location", func(t *testing.T) {
		msg := &common.Address{
			Street:  "123 Main St",
			City:    "San Francisco",
			State:   "CA",
			Country: "USA",
			ZipCode: "94102",
			Geo: &common.GeoLocation{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		}

		mask := map[string]struct{}{
			"geo.latitude":  {},
			"geo.longitude": {},
		}

		Filter(msg, mask)

		// All top-level fields should be cleared
		assert.Empty(t, msg.Street, "expected Street to be empty")
		assert.Empty(t, msg.City, "expected City to be empty")
		assert.Empty(t, msg.State, "expected State to be empty")
		assert.Empty(t, msg.Country, "expected Country to be empty")
		assert.Empty(t, msg.ZipCode, "expected ZipCode to be empty")

		// Only masked nested fields should remain
		assert.NotNil(t, msg.Geo, "expected Geo to not be nil")
		assert.Equal(t, 37.7749, msg.Geo.Latitude, "expected Latitude 37.7749")
		assert.Equal(t, -122.4194, msg.Geo.Longitude, "expected Longitude -122.4194")
	})
}

// TestPrune_NestedFields tests Prune with deeply nested field paths
func TestPrune_NestedFields(t *testing.T) {
	t.Run("deeply nested geo location", func(t *testing.T) {
		msg := &common.Address{
			Street:  "123 Main St",
			City:    "San Francisco",
			State:   "CA",
			Country: "USA",
			ZipCode: "94102",
			Geo: &common.GeoLocation{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		}

		mask := map[string]struct{}{
			"geo.latitude":  {},
			"geo.longitude": {},
		}

		Prune(msg, mask)

		// All top-level fields should remain (not in mask)
		assert.Equal(t, "123 Main St", msg.Street, "expected Street '123 Main St'")
		assert.Equal(t, "San Francisco", msg.City, "expected City 'San Francisco'")
		assert.Equal(t, "CA", msg.State, "expected State 'CA'")
		assert.Equal(t, "USA", msg.Country, "expected Country 'USA'")
		assert.Equal(t, "94102", msg.ZipCode, "expected ZipCode '94102'")

		// Only masked nested fields should be cleared
		assert.NotNil(t, msg.Geo, "expected Geo to not be nil")
		assert.Equal(t, float64(0), msg.Geo.Latitude, "expected Latitude to be 0 (pruned)")
		assert.Equal(t, float64(0), msg.Geo.Longitude, "expected Longitude to be 0 (pruned)")
	})
}

// TestFilter_TopLevelAndNestedFields tests Filter with both top-level and nested fields
func TestFilter_TopLevelAndNestedFields(t *testing.T) {
	t.Run("mix of top-level and nested fields", func(t *testing.T) {
		msg := &common.Address{
			Street:  "123 Main St",
			City:    "San Francisco",
			State:   "CA",
			Country: "USA",
			ZipCode: "94102",
			Geo: &common.GeoLocation{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		}

		mask := map[string]struct{}{
			"street":       {},
			"city":         {},
			"geo.latitude": {},
		}

		Filter(msg, mask)

		assert.Equal(t, "123 Main St", msg.Street, "expected Street '123 Main St'")
		assert.Equal(t, "San Francisco", msg.City, "expected City 'San Francisco'")
		assert.Empty(t, msg.State, "expected State to be empty")
		assert.Empty(t, msg.Country, "expected Country to be empty")
		assert.Empty(t, msg.ZipCode, "expected ZipCode to be empty")
		assert.NotNil(t, msg.Geo, "expected Geo to not be nil")
		assert.Equal(t, 37.7749, msg.Geo.Latitude, "expected Latitude 37.7749")
		assert.Equal(t, float64(0), msg.Geo.Longitude, "expected Longitude to be 0")
	})
}

// TestPrune_TopLevelAndNestedFields tests Prune with both top-level and nested fields
func TestPrune_TopLevelAndNestedFields(t *testing.T) {
	t.Run("mix of top-level and nested fields", func(t *testing.T) {
		msg := &common.Address{
			Street:  "123 Main St",
			City:    "San Francisco",
			State:   "CA",
			Country: "USA",
			ZipCode: "94102",
			Geo: &common.GeoLocation{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		}

		mask := map[string]struct{}{
			"state":        {},
			"country":      {},
			"geo.latitude": {},
		}

		Prune(msg, mask)

		assert.Equal(t, "123 Main St", msg.Street, "expected Street '123 Main St'")
		assert.Equal(t, "San Francisco", msg.City, "expected City 'San Francisco'")
		assert.Empty(t, msg.State, "expected State to be empty (pruned)")
		assert.Empty(t, msg.Country, "expected Country to be empty (pruned)")
		assert.Equal(t, "94102", msg.ZipCode, "expected ZipCode '94102'")
		assert.NotNil(t, msg.Geo, "expected Geo to not be nil")
		assert.Equal(t, float64(0), msg.Geo.Latitude, "expected Latitude to be 0 (pruned)")
		assert.Equal(t, -122.4194, msg.Geo.Longitude, "expected Longitude -122.4194")
	})
}

// TestFilter_PreservesOriginalMessage tests that Filter doesn't affect the original message
func TestFilter_PreservesOriginalMessage(t *testing.T) {
	t.Run("preserves original when using proto.Clone", func(t *testing.T) {
		original := &common.CommonMessage{
			Id:        "msg123",
			CreatedAt: "2024-01-01",
			UpdatedAt: "2024-01-02",
		}

		// Clone the message before filtering
		cloned := proto.Clone(original).(*common.CommonMessage)
		mask := map[string]struct{}{
			"id": {},
		}
		Filter(cloned, mask)

		// Original should be unchanged
		assert.Equal(t, "msg123", original.Id, "expected original Id 'msg123'")
		assert.Equal(t, "2024-01-01", original.CreatedAt, "expected original CreatedAt '2024-01-01'")
		assert.Equal(t, "2024-01-02", original.UpdatedAt, "expected original UpdatedAt '2024-01-02'")

		// Cloned should be filtered
		assert.Equal(t, "msg123", cloned.Id, "expected cloned Id 'msg123'")
		assert.Empty(t, cloned.CreatedAt, "expected cloned CreatedAt to be empty")
		assert.Empty(t, cloned.UpdatedAt, "expected cloned UpdatedAt to be empty")
	})
}

// TestMultiLevelNestedFields tests deeply nested field paths like "address.geo.latitude"
func TestMultiLevelNestedFields(t *testing.T) {
	t.Run("FILTER mode with multi-level nested fields", func(t *testing.T) {
		msg := &common.Address{
			Street:  "123 Main St",
			City:    "San Francisco",
			State:   "CA",
			Country: "USA",
			ZipCode: "94102",
			Geo: &common.GeoLocation{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		}

		// Mask only geo.latitude (multi-level nested)
		mask := map[string]struct{}{
			"geo.latitude": {},
		}

		Filter(msg, mask)

		// All top-level fields should be cleared
		assert.Empty(t, msg.Street, "expected Street to be empty")
		assert.Empty(t, msg.City, "expected City to be empty")
		assert.Empty(t, msg.State, "expected State to be empty")
		assert.Empty(t, msg.Country, "expected Country to be empty")
		assert.Empty(t, msg.ZipCode, "expected ZipCode to be empty")

		// Only geo.latitude should remain
		assert.NotNil(t, msg.Geo, "expected Geo to not be nil")
		assert.Equal(t, 37.7749, msg.Geo.Latitude, "expected Latitude 37.7749")
		assert.Equal(t, float64(0), msg.Geo.Longitude, "expected Longitude to be 0 (cleared)")
	})

	t.Run("PRUNE mode with multi-level nested fields", func(t *testing.T) {
		msg := &common.Address{
			Street:  "123 Main St",
			City:    "San Francisco",
			State:   "CA",
			Country: "USA",
			ZipCode: "94102",
			Geo: &common.GeoLocation{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		}

		// Prune only geo.latitude
		mask := map[string]struct{}{
			"geo.latitude": {},
		}

		Prune(msg, mask)

		// All top-level fields should remain (not in mask)
		assert.Equal(t, "123 Main St", msg.Street, "expected Street '123 Main St'")
		assert.Equal(t, "San Francisco", msg.City, "expected City 'San Francisco'")
		assert.Equal(t, "CA", msg.State, "expected State 'CA'")
		assert.Equal(t, "USA", msg.Country, "expected Country 'USA'")
		assert.Equal(t, "94102", msg.ZipCode, "expected ZipCode '94102'")

		// Only geo.latitude should be cleared
		assert.NotNil(t, msg.Geo, "expected Geo to not be nil")
		assert.Equal(t, float64(0), msg.Geo.Latitude, "expected Latitude to be 0 (pruned)")
		assert.Equal(t, -122.4194, msg.Geo.Longitude, "expected Longitude -122.4194")
	})
}
