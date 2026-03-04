# FieldMask End-to-End Example

This example demonstrates the complete flow of using FieldMask:
1. Client sets FieldMask to specify desired fields
2. Server processes request based on FieldMask
3. Client receives response and analyzes which fields are masked/unmasked

## Running the Example

```bash
# From this directory
go run .

# Or specify all Go files
go run *.go
```

**Note:** Do NOT run `go run main.go` - this will only compile main.go and miss the client.go implementation.

## Scenarios Demonstrated

### Scenario 1: FILTER Mode (Partial Response)

In FILTER mode, only MARKED fields are returned. Unmarked fields are cleared.

#### 1a. Mobile Client
```
Request: user_id + username only
Response: Only user_id and username have values, others are empty
```

#### 1b. Web Client
```
Request: All fields except password_hash
Response: All marked fields have values, password_hash is empty (security)
```

#### 1c. No FieldMask (Legacy)
```
Request: No FieldMask
Response: All fields returned including sensitive password_hash
```

### Scenario 2: PRUNE Mode (Partial Update)

In PRUNE mode, MARKED fields are INCLUDED in update. Only marked fields are updated.

#### 2a. Update Profile Only
```
Request: Mark username/avatar_url/bio for update (email/phone NOT marked)
Result: Only profile fields updated, email/phone unchanged
```

#### 2b. Update Contact Only
```
Request: Mark email/phone for update (profile fields NOT marked)
Result: Only contact fields updated, profile fields unchanged
```

## Code Structure

- `main.go` - Entry point that runs all scenarios
- `client.go` - Client implementation with server mock

## Output Format

The example shows:
- `>>> Client:` - Client-side actions
- `>>> Server:` - Server-side processing
- `Response Analysis:` - Shows which fields are masked (✓) or empty
