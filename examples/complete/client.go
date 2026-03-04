package main

import (
	"context"
	"fmt"

	userv1 "github.com/yeqown/protoc-gen-fieldmask/examples/proto"
)

// Client demonstrates FieldMask usage from the client's perspective
type Client struct {
	server *Server
}

func NewClient(server *Server) *Client {
	return &Client{server: server}
}

// ============================================================================
// Scenario 1: Partial Response (FILTER mode)
// Client specifies which fields to include in response
// ============================================================================

// Scenario1_MobileClient demonstrates a mobile client that only needs minimal data
func (c *Client) Scenario1_MobileClient() {
	fmt.Println("\n========== SCENARIO 1: FILTER Mode - Mobile Client ==========")
	fmt.Println("Client Goal: Only fetch user_id and username (minimal data)")
	fmt.Println("Expected: Only user_id and username are returned, other fields are empty")
	fmt.Println()

	req := &userv1.GetUserRequest{UserId: "user123"}

	// Client sets FieldMask to specify which fields to return
	fmt.Println(">>> Client: Setting FieldMask...")
	fm := req.FieldMask().Response()
	fm.MaskUserId()
	fm.MaskUsername()

	// Show what fields are marked
	fmt.Printf(">>> Client: FieldMask paths = %v\n", req.Fm.Paths)
	fmt.Printf(">>> Client: Masked user_id = %v\n", fm.MaskedUserId())
	fmt.Printf(">>> Client: Masked username = %v\n", fm.MaskedUsername())
	fmt.Printf(">>> Client: Masked email = %v\n", fm.MaskedEmail())
	fmt.Println()

	// Send request
	fmt.Println(">>> Client: Sending request to server...")
	resp, _ := c.server.GetUser(context.Background(), req)

	// Display response - show which fields are masked (returned) and which are not
	fmt.Println(">>> Client: Received response from server")
	fmt.Println()
	fmt.Println("Response Analysis:")
	fmt.Println("  ├─ Masked (Requested) fields:")
	displayField("  │  ├── user_id", resp.UserId, resp.UserId != "")
	displayField("  │  └── username", resp.Username, resp.Username != "")
	fmt.Println("  └─ Unmasked (Not requested) fields:")
	displayField("      ├── email", resp.Email, resp.Email != "")
	displayField("      ├── phone", resp.Phone, resp.Phone != "")
	displayField("      ├── bio", resp.Bio, resp.Bio != "")
	displayField("      ├── password_hash", resp.PasswordHash, resp.PasswordHash != "")
	displayField("      └── created_at", resp.CreatedAt, resp.CreatedAt != 0)
	fmt.Println()
}

// Scenario1_WebClient demonstrates a web client that needs most fields but excludes sensitive data
func (c *Client) Scenario1_WebClient() {
	fmt.Println("\n========== SCENARIO 1: FILTER Mode - Web Client ==========")
	fmt.Println("Client Goal: Fetch all fields except password_hash (sensitive data)")
	fmt.Println("Expected: All fields are returned except password_hash (empty)")
	fmt.Println()

	req := &userv1.GetUserRequest{UserId: "user456"}

	// Client sets FieldMask to include all fields except password_hash
	fmt.Println(">>> Client: Setting FieldMask...")
	fm := req.FieldMask().Response()
	fm.MaskUserId()
	fm.MaskUsername()
	fm.MaskEmail()
	fm.MaskPhone()
	fm.MaskAvatarUrl()
	fm.MaskBio()
	fm.MaskCreatedAt()
	fm.MaskUpdatedAt()
	// Note: password_hash and status are NOT marked - they will be excluded from response

	// Show what fields are marked
	fmt.Printf(">>> Client: FieldMask paths = %v\n", req.Fm.Paths)
	fmt.Printf(">>> Client: Total marked fields: %d\n", len(req.Fm.Paths))
	fmt.Println()

	// Send request
	fmt.Println(">>> Client: Sending request to server...")
	resp, _ := c.server.GetUser(context.Background(), req)

	// Display response
	fmt.Println(">>> Client: Received response from server")
	fmt.Println()
	fmt.Println("Response Analysis:")
	fmt.Println("  ├─ Masked (Requested) fields:")
	displayField("  │  ├── user_id", resp.UserId, resp.UserId != "")
	displayField("  │  ├── username", resp.Username, resp.Username != "")
	displayField("  │  ├── email", resp.Email, resp.Email != "")
	displayField("  │  ├── phone", resp.Phone, resp.Phone != "")
	displayField("  │  ├── avatar_url", resp.AvatarUrl, resp.AvatarUrl != "")
	displayField("  │  ├── bio", resp.Bio, resp.Bio != "")
	displayField("  │  └── created_at", resp.CreatedAt, resp.CreatedAt != 0)
	fmt.Println("  └─ Unmasked (Not requested - Sensitive/Not marked) fields:")
	displayField("      ├── password_hash", resp.PasswordHash, resp.PasswordHash != "")
	displayField("      └── status", resp.Status, resp.Status != "")
	fmt.Println()
	fmt.Println("  ✓ Sensitive password_hash was NOT returned (security)")
	fmt.Println()
}

// Scenario1_NoFieldMask demonstrates what happens when no FieldMask is provided
func (c *Client) Scenario1_NoFieldMask() {
	fmt.Println("\n========== SCENARIO 1: FILTER Mode - No FieldMask ==========")
	fmt.Println("Client Goal: Fetch all data (legacy client without FieldMask support)")
	fmt.Println("Expected: All fields are returned including password_hash")
	fmt.Println()

	req := &userv1.GetUserRequest{UserId: "user789"}

	// Client does NOT set FieldMask
	fmt.Println(">>> Client: NO FieldMask set (legacy client)")
	if req.Fm == nil {
		fmt.Println(">>> Client: FieldMask = nil")
	} else {
		fmt.Printf(">>> Client: FieldMask paths = %v\n", req.Fm.Paths)
	}
	fmt.Println()

	// Send request
	fmt.Println(">>> Client: Sending request to server...")
	resp, _ := c.server.GetUser(context.Background(), req)

	// Display response
	fmt.Println(">>> Client: Received response from server")
	fmt.Println()
	fmt.Println("Response Analysis:")
	fmt.Println("  All fields (including sensitive data):")
	displayField("  ├── user_id", resp.UserId, resp.UserId != "")
	displayField("  ├── username", resp.Username, resp.Username != "")
	displayField("  ├── email", resp.Email, resp.Email != "")
	displayField("  ├── password_hash", resp.PasswordHash, resp.PasswordHash != "")
	displayField("  └── created_at", resp.CreatedAt, resp.CreatedAt != 0)
	fmt.Println()
	fmt.Println("  ⚠ WARNING: password_hash was returned (no FieldMask to filter it)")
	fmt.Println()
}

// ============================================================================
// Scenario 2: Partial Update (PRUNE mode)
// Client specifies which fields to EXCLUDE from update
// ============================================================================

// Scenario2_UpdateProfileOnly demonstrates updating profile while excluding contact fields
func (c *Client) Scenario2_UpdateProfileOnly() {
	fmt.Println("\n========== SCENARIO 2: PRUNE Mode - Update Profile Only ==========")
	fmt.Println("Client Goal: Update username, avatar_url, bio")
	fmt.Println("             BUT EXCLUDE email and phone from update")
	fmt.Println("Expected: Only profile fields are updated, email/phone unchanged")
	fmt.Println()

	req := &userv1.UpdateUserRequest{
		UserId:    "user123",
		Username:  "new_username",
		Email:     "should_not_change@example.com", // Has value but should be excluded
		Phone:     "+9999999999",                   // Has value but should be excluded
		AvatarUrl: "https://example.com/new-avatar.jpg",
		Bio:       "Updated bio",
	}

	// Client sets FieldMask to mark which fields to UPDATE
	fmt.Println(">>> Client: Setting FieldMask to mark fields to UPDATE...")
	fmt.Println(">>> Client: Request data:")
	fmt.Printf("    username = %q (should UPDATE)\n", req.Username)
	fmt.Printf("    email = %q (should NOT update)\n", req.Email)
	fmt.Printf("    phone = %q (should NOT update)\n", req.Phone)
	fmt.Printf("    avatar_url = %q (should UPDATE)\n", req.AvatarUrl)
	fmt.Printf("    bio = %q (should UPDATE)\n", req.Bio)
	fmt.Println()

	fm := req.FieldMask()
	// Mark fields to UPDATE (not marking email/phone means don't update them)
	fm.Request().MaskUsername()
	fm.Request().MaskAvatarUrl()
	fm.Request().MaskBio()
	fm.Response().MaskUserId()
	fm.Response().MaskUsername()
	fm.Response().MaskEmail()

	fmt.Printf(">>> Client: FieldMask paths = %v\n", req.Fm.Paths)
	fmt.Println("  req.username = marked for UPDATE")
	fmt.Println("  req.avatar_url = marked for UPDATE")
	fmt.Println("  req.bio = marked for UPDATE")
	fmt.Println("  (email and phone are NOT marked - will not be updated)")
	fmt.Println()

	// Send request
	fmt.Println(">>> Client: Sending update request to server...")
	resp, _ := c.server.UpdateUser(context.Background(), req)

	// Display response
	fmt.Println(">>> Client: Received response from server")
	fmt.Println()
	fmt.Println("Response Analysis:")
	fmt.Println("  ├─ Updated (Marked) fields:")
	displayField("  │  ├── username", resp.Username, resp.Username == "new_username")
	fmt.Println("  │     Status: ✓ Updated")
	displayField("  │  ├── avatar_url", resp.AvatarUrl, resp.AvatarUrl == "https://example.com/new-avatar.jpg")
	fmt.Println("  │     Status: ✓ Updated")
	displayField("  │  └── bio", resp.Bio, resp.Bio == "Updated bio")
	fmt.Println("     Status: ✓ Updated")
	fmt.Println("  └─ Not Marked (Not Updated) fields:")
	displayField("      ├── email", resp.Email, resp.Email == "john@example.com")
	fmt.Println("         Status: ✓ NOT updated (not marked)")
	displayField("      └── phone", resp.Phone, resp.Phone == "+1234567890")
	fmt.Println("         Status: ✓ NOT updated (not marked)")
	fmt.Println()
}

// Scenario2_UpdateContactOnly demonstrates updating contacts while excluding profile fields
func (c *Client) Scenario2_UpdateContactOnly() {
	fmt.Println("\n========== SCENARIO 2: PRUNE Mode - Update Contact Only ==========")
	fmt.Println("Client Goal: Update email and phone")
	fmt.Println("             BUT EXCLUDE profile fields (username, bio, avatar) from update")
	fmt.Println("Expected: Only contact fields are updated, profile fields unchanged")
	fmt.Println()

	req := &userv1.UpdateUserRequest{
		UserId:    "user456",
		Username:  "should_not_change", // Has value but should be excluded
		Email:     "newemail@example.com",
		Phone:     "+1111111111",
		AvatarUrl: "should_not_change", // Has value but should be excluded
		Bio:       "should_not_change", // Has value but should be excluded
	}

	// Client sets FieldMask to mark which fields to UPDATE
	fmt.Println(">>> Client: Setting FieldMask to mark fields to UPDATE...")
	fmt.Println(">>> Client: Request data:")
	fmt.Printf("    username = %q (should NOT update)\n", req.Username)
	fmt.Printf("    email = %q (should UPDATE)\n", req.Email)
	fmt.Printf("    phone = %q (should UPDATE)\n", req.Phone)
	fmt.Printf("    avatar_url = %q (should NOT update)\n", req.AvatarUrl)
	fmt.Printf("    bio = %q (should NOT update)\n", req.Bio)
	fmt.Println()

	fm := req.FieldMask()
	// Mark fields to UPDATE (not marking profile fields means don't update them)
	fm.Request().MaskEmail()
	fm.Request().MaskPhone()
	fm.Response().MaskUserId()
	fm.Response().MaskEmail()
	fm.Response().MaskPhone()

	fmt.Printf(">>> Client: FieldMask paths = %v\n", req.Fm.Paths)
	fmt.Println("  req.email = marked for UPDATE")
	fmt.Println("  req.phone = marked for UPDATE")
	fmt.Println("  (username, avatar_url, bio are NOT marked - will not be updated)")
	fmt.Println()

	// Send request
	fmt.Println(">>> Client: Sending update request to server...")
	resp, _ := c.server.UpdateUser(context.Background(), req)

	// Display response
	fmt.Println(">>> Client: Received response from server")
	fmt.Println()
	fmt.Println("Response Analysis:")
	fmt.Println("  ├─ Updated (Marked) fields:")
	displayField("  │  ├── email", resp.Email, resp.Email == "newemail@example.com")
	fmt.Println("  │     Status: ✓ Updated")
	displayField("  │  └── phone", resp.Phone, resp.Phone == "+1111111111")
	fmt.Println("     Status: ✓ Updated")
	fmt.Println("  └─ Not Marked (Not Updated) fields:")
	displayField("      ├── username", resp.Username, resp.Username == "johndoe")
	fmt.Println("         Status: ✓ NOT updated (not marked)")
	displayField("      ├── avatar_url", resp.AvatarUrl, resp.AvatarUrl == "https://example.com/avatar.jpg")
	fmt.Println("         Status: ✓ NOT updated (not marked)")
	displayField("      └── bio", resp.Bio, resp.Bio == "Software developer")
	fmt.Println("         Status: ✓ NOT updated (not marked)")
	fmt.Println()
}

// ============================================================================
// Helper functions
// ============================================================================

func displayField(label string, value interface{}, hasValue bool) {
	valueStr := fmt.Sprintf("%v", value)
	if hasValue {
		fmt.Printf("%s = %s ✓\n", label, valueStr)
	} else {
		fmt.Printf("%s = %s (empty - masked out)\n", label, valueStr)
	}
}

// ============================================================================
// Server (for demonstration)
// ============================================================================

// Server represents the server that processes requests
type Server struct {
	users map[string]*UserData
}

type UserData struct {
	Id        string
	Username  string
	Email     string
	Phone     string
	AvatarUrl string
	Bio       string
	Status    string
	CreatedAt int64
	UpdatedAt int64
}

func NewServer() *Server {
	return &Server{
		users: make(map[string]*UserData),
	}
}

func (s *Server) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	// Get user from storage (simulated)
	user := s.getUser(req.UserId)
	fm := req.FieldMask()

	// Implementation method1: Build full response
	resp := &userv1.GetUserResponse{
		UserId:       user.Id,
		Username:     user.Username,
		Email:        user.Email,
		Phone:        user.Phone,
		AvatarUrl:    user.AvatarUrl,
		Bio:          user.Bio,
		Status:       user.Status,
		PasswordHash: "secret_hash_12345", // Sensitive data
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	// Implementation method2: Only executing the necessary querying
	// This is useful while the query API is a big "Aggregation"。
	//
	// resp := &userv1.GetUserResponse{
	// 	UserId:       user.Id,
	// 	Username:     user.Username,
	// 	Email:        user.Email,
	// 	Phone:        user.Phone,
	// 	AvatarUrl:    "",
	// 	Bio:          "",
	// 	Status:       user.Status,
	// 	PasswordHash: "secret_hash_12345", // Sensitive data
	// 	CreatedAt:    user.CreatedAt,
	// 	UpdatedAt:    user.UpdatedAt,
	// }
	// if fm.Response().MaskedAvatarUrl() {
	// 	resp.Bio = user.AvatarUrl
	// }
	// if fm.Response().MaskedBio() {
	// 	resp.Bio = user.Bio
	// }

	// Server applies FieldMask filter using generated code
	// Only call Apply() if Fm is not nil (nil means no filtering requested)
	fm.Response().Apply(resp)

	return resp, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	// Get existing user
	user := s.getUser(req.UserId)

	// Use generated FieldMask API to check which fields are marked for update
	// fm.Request().Masked*() returns true when field is marked (should be updated)
	fm := req.FieldMask()
	var markedFields []string
	if fm.Request().MaskedUsername() {
		markedFields = append(markedFields, "username")
	}
	if fm.Request().MaskedEmail() {
		markedFields = append(markedFields, "email")
	}
	if fm.Request().MaskedPhone() {
		markedFields = append(markedFields, "phone")
	}
	if fm.Request().MaskedAvatarUrl() {
		markedFields = append(markedFields, "avatar_url")
	}
	if fm.Request().MaskedBio() {
		markedFields = append(markedFields, "bio")
	}
	if len(markedFields) > 0 {
		fmt.Printf(">>> Server: Fields marked for update (using fm.Request().Masked*()): %v\n", markedFields)
	}

	// Update only fields that are marked (using generated API)
	if fm.Request().MaskedUsername() {
		user.Username = req.Username
	}
	if fm.Request().MaskedEmail() {
		user.Email = req.Email
	}
	if fm.Request().MaskedPhone() {
		user.Phone = req.Phone
	}
	if fm.Request().MaskedAvatarUrl() {
		user.AvatarUrl = req.AvatarUrl
	}
	if fm.Request().MaskedBio() {
		user.Bio = req.Bio
	}

	// Build response
	resp := &userv1.UpdateUserResponse{
		UserId:    user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		AvatarUrl: user.AvatarUrl,
		Bio:       user.Bio,
		Status:    user.Status,
		UpdatedAt: user.UpdatedAt,
	}

	// Apply response filter using generated code
	// Only call Apply() if Fm is not nil (nil means no filtering requested)
	if req.Fm != nil {
		fm.Response().Apply(resp)
	}

	return resp, nil
}

func (s *Server) getUser(userID string) *UserData {
	if user, ok := s.users[userID]; ok {
		return user
	}
	user := &UserData{
		Id:        userID,
		Username:  "johndoe",
		Email:     "john@example.com",
		Phone:     "+1234567890",
		AvatarUrl: "https://example.com/avatar.jpg",
		Bio:       "Software developer",
		Status:    "active",
		CreatedAt: 1609459200,
		UpdatedAt: 1609459200,
	}
	s.users[userID] = user
	return user
}
