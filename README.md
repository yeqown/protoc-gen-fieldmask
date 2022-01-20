# protoc-gen-fieldmask

Generate FieldMask utility functions for protobuf

### Generated Preview

coding proto file：

```protobuf
syntax= "proto3";

import "google/protobuf/types/known/fieldmaskpb.proto";

message UserInfoRequest {
  string user_id = 1;
  google.protobuf.FieldMask field_mask = 2 [option = {message: "UserInfoResponse"}];
}

message Address {
  string country = 1;
  string province = 2;
}

message UserInfoResponse {
  string user_id = 1;
  string name = 2;
  string email = 3;
  Address address = 4;
}
```

generated `*.fm.go`：

```go
package api

type FieldMask_Mode uint8 

const (
	FieldMask_FILTER FieldMask_Mode = iota
	FieldMask_PRUNE
)

type UserInfoRequest struct {
	// ...
}

func (req *UserInfoRequest) Mask_UserId() *UserInfoRequest {return nil}
func (req *UserInfoRequest) Mask_Name() *UserInfoRequest {return nil}
func (req *UserInfoRequest) Mask_Email() *UserInfoRequest {return nil}
func (req *UserInfoRequest) Mask_Adress() *UserInfoRequest {return nil}
func (req *UserInfoRequest) Mask_Adress_Country() *UserInfoRequest {return nil}
func (req *UserInfoRequest) Mask_Adress_Province() *UserInfoRequest {return nil}
// FieldMask generated a message_FieldMask from UserInfoRequest
//
// mode：decide which mode will UserInfoResponse_FieldMask acts.
// prune：remove fields those not in UserInfoRequest.field_mask.
// filter：keep fields those in UserInfoRequest.field_mask.
func (req *UserInfoRequest) FieldMask(mode FieldMask_Mode) *UserInfoResponse_FieldMask {return nil}


type UserInfoResponse struct {}

// UserInfoResponse_FieldMask is a functions set to help FieldMask usage. 
type UserInfoResponse_FieldMask struct {}
func (fm *UserInfoResponse_FieldMask) Include_UserId() bool {return false}
func (fm *UserInfoResponse_FieldMask) Include_Name() bool {return false}
// ... more Include_xxx()

func (fm *UserInfoResponse_FieldMask) Do(msg *UserInfoResponse) *UserInfoResponse {return nil}
```