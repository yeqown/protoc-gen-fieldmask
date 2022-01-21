# protoc-gen-fieldmask

Generate FieldMask utilities for protobuf, support [Go](https://golang.org), maybe more programing languages
later.

### Installation

```sh
go install github.com/yeqwown/protoc-gen-fieldmask@latest
```

### Get Started

```sh
protoc \
	-I. \
	-I$YOUR_PROTO_PATH \
	--go_out=paths=source_relative:. \
	--fieldmask_out=paths=source_relative,lang=go:. \
	example.proto
```

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

func (req *UserInfoRequest) Mask_UserId() *UserInfoRequest {return nil}
func (req *UserInfoRequest) Mask_Name() *UserInfoRequest {return nil}
func (req *UserInfoRequest) Mask_Email() *UserInfoRequest {return nil}
func (req *UserInfoRequest) Mask_Adress() *UserInfoRequest {return nil}
func (req *UserInfoRequest) Mask_Adress_Country() *UserInfoRequest {return nil}
func (req *UserInfoRequest) Mask_Adress_Province() *UserInfoRequest {return nil}
// FieldMaskWithMode generated a message_FieldMask from UserInfoRequest
//
// mode：decide which mode will UserInfoResponse_FieldMask acts.
// prune：remove fields those not in UserInfoRequest.field_mask.
// filter：keep fields those in UserInfoRequest.field_mask.
func (req *UserInfoRequest) FieldMaskWithMode(mode FieldMask_Mode) *UserInfoResponse_FieldMask {return nil}
func (req *UserInfoRequest) FieldMask_Filter() *UserInfoResponse_FieldMask {return nil}
func (req *UserInfoRequest) FieldMask_Prune() *UserInfoResponse_FieldMask {return nil}


type UserInfoResponse struct {}

// UserInfoResponse_FieldMask is a functions set to help FieldMask usage. 
type UserInfoResponse_FieldMask struct {
	maskedMap map[string]struct{}
}
func (fm *UserInfoResponse_FieldMask) Masked_UserId() bool {return false}
func (fm *UserInfoResponse_FieldMask) Masked_Name() bool {return false}
// ... more Include_xxx()

func (fm *UserInfoResponse_FieldMask) Mask(msg *UserInfoResponse) *UserInfoResponse {return nil}
```

### How to debug

- prepare a `debugdata`
- install `protoc-gen-debug`: `go install github.com/lyft/protoc-gen-star/protoc-gen-debug@latest`
- compile target proto file with `protoc-gen-debug`: 
	
    ```sh
    protoc \
        -I=./examples/normal \
        -I=./proto \
        --plugin=protoc-gen-debug=$(which protoc-gen-debug) \
        --debug_out="./debugdata,lang=go:./debugdata" \
        ./examples/normal/user.proto
    ```
- debug [Test_ForDebug](./internal/module/fieldmask_test.go#L46) test suite 
in [internal/module/fieldmask_test.go](./internal/module/fieldmask_test.go)