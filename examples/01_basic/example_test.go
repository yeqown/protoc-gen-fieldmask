package basic_test

import (
	"testing"

	"github.com/yeqown/protoc-gen-fieldmask/examples/01_basic/basic"
	pbfieldmask "github.com/yeqown/protoc-gen-fieldmask/proto/protobuf"
)

// Test_BasicFieldMask demonstrates the most basic usage of FieldMask.
func Test_BasicFieldMask(t *testing.T) {
	// Simulate a request from a client
	req := &basic.GetUserRequest{
		UserId: "user123",
	}

	// Client wants to include only user_id and name in response
	// The email field will be omitted
	fm := req.FieldMask(pbfieldmask.FILTER)

	// Mark fields to include in response
	fm.Response().UserId()
	fm.Response().Name()

	// Simulate response data
	resp := &basic.GetUserResponse{
		UserId: "user123",
		Name:   "John Doe",
		Email:  "john@example.com", // This will be filtered out
	}

	// Check which fields are marked
	if fm.Marked().Response().UserId() {
		t.Log("UserId is marked")
	}
	if fm.Marked().Response().Name() {
		t.Log("Name is marked")
	}
	if !fm.Marked().Response().Email() {
		t.Log("Email is NOT marked (will be filtered)")
	}

	// In a real implementation, you would apply the mask to the response
	// This filters out fields that are not marked
	// fm.Response().Apply(resp)

	t.Log("Basic FieldMask usage demonstrated successfully")
}
