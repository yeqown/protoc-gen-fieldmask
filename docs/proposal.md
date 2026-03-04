# FieldMask Options V2 Proposal

## 概述

重新设计 protoc-gen-fieldmask 的 options，采用 **MethodOptions** 和 **FieldOptions** 结合的方式，提供更灵活和精细的 field mask 控制能力。

## Option 定义

### option.proto

```protobuf
syntax = "proto3";

package protoc_gen_fieldmask;

option go_package = "github.com/yeqown/protoc-gen-fieldmask/proto/fieldmask;fieldmask";

import "google/protobuf/descriptor.proto";

// ============================================================================
// Method-level options: 关联 RPC 的 request 和 response
// ============================================================================

extend google.protobuf.MethodOptions {
  optional MethodOptions rpc = 1142;
}

// MethodOptions 定义 RPC 方法级别的 field mask 配置
message MethodOptions {
  // request 中 field_mask 字段名称 (默认: "fm")
  string field_name = 1;
  
  // mask 模式: FILTER 或 PRUNE (默认: FILTER)
  MaskMode mode = 2;
}

enum MaskMode {
  FILTER = 0;  // 被 mask 的字段有效，其他字段无效（默认）
  PRUNE = 1;   // 被 mask 的字段无效，其他字段有效
}

// ============================================================================
// Field-level options: 精细控制字段的 mask 行为
// ============================================================================

extend google.protobuf.FieldOptions {
  optional FieldOptions field = 1142;
}

// FieldOptions 定义单个字段的 mask 行为
message FieldOptions {
  // 是否为该字段生成 mask 方法（默认: false）
  bool mask = 1;
  
  // 是否为嵌套字段生成 mask 方法，仅当对应的字段类型为 Message 时有效（默认: false）
  bool nested = 2;
}
```

## 使用示例

```bash
protoc \
        -I./path/to/proto \
        --go_out=paths=source_relative:./pb \
        --fieldmask_out=paths=source_relative,lang=go:./pb \
        ./pb/user.proto
```

### CASE 1: 基本用法

针对单个方法的 request 和 response 中的字段选项进行配置

```protobuf
// user.proto
syntax = "proto3";

package example;

import "google/protobuf/field_mask.proto";
import "protoc_gen_fieldmask/option.proto";

message UserInfoRequest {
  string user_id = 1;
  google.protobuf.FieldMask fm = 2;
}

message UserInfoResponse {
  string user_id = 1;
  string name = 2;
  string email = 3 [(protoc_gen_fieldmask.field) = {mask: true}];
  Address address = 4 [(protoc_gen_fieldmask.field) = {mask: true, nested: true}];
}

message Address {
  string country = 1;
  string province = 2 [(protoc_gen_fieldmask.field) = {mask: true}]; // 不针对这个字段生成 Mask API
}

service UserService {
  rpc GetUserInfo(UserInfoRequest) returns (UserInfoResponse) {
    option (protoc_gen_fieldmask.rpc) = {
      field_name: "fm"
      mode: FILTER
    };
  }
}
```

### CASE 2: 跨包引用

在 proto 中引用其他包的 message 时，pgfm 不会为其生成 mask 方法，因此 request 和 response 需要成对出现。

```protobuf
// order.proto
syntax = "proto3";

package example;

import "user.proto";

import "google/protobuf/field_mask.proto";
import "protoc_gen_fieldmask/option.proto";

message OrderInfoRequest {
  string order_id = 1;
  google.protobuf.FieldMask fm = 2;
}

message OrderInfoResponse {
  string order_id = 1;
  string user_id = 2;
  example.Address address = 3 [(protoc_gen_fieldmask.field) = {mask: true, nested: true}];
}

service OrderService {
  rpc GetUserAddress(example.UserInfoRequest) returns (example.UserInfoRequest) {
    option (protoc_gen_fieldmask.rpc) = {
      field_name: "fm"
      mask_mode: FILTER
    };
  }
  
  rpc GetOrderDetail(OrderInfoRequest) returns (OrderInfoResponse) {
	option (protoc_gen_fieldmask.rpc) = {
      field_name: "fm"
      mask_mode: FILTER
    };
  }
}
```

### CASE3: 增量更新

针对 Request 中只标记特定字段需要更新，其余字段则不需要。

```protobuf
// update.proto
syntax = "proto3";

import "google/protobuf/field_mask.proto";
import "protoc_gen_fieldmask/option.proto";

package example;

message UpdateUserInfoRequest {
  string user_id = 1;
  string name = 2 [(protoc_gen_fieldmask.field) = {mask: true}];
  string email = 3 [(protoc_gen_fieldmask.field) = {mask: true}];
}

message UpdateUserInfoResponse {}

service IncrementalUpdateService {
  rpc UpdateUserInfo(UpdateUserInfoRequest) returns (UpdateUserInfoResponse) {
    option (protoc_gen_fieldmask.rpc) = {
      field_name: "fm"
      mask_mode: FILTER
    };
  }
}
```

## 预期生成代码

### API 设计概览

```go
// 统一 FieldMask 入口
fm := req.FieldMask()  

// 标记 request 自身字段, 适用于增量更新
fm.Request().MaskUserId()
fm.Request().MaskName()

// 标记 response 字段，适用于仅获取特定字段
fm.Response().MaskName()
fm.Response().MaskAddress()
fm.Response().MaskAddress_Country()

// 判断字段是否被标记
fm.Request().MaskedUserId()           // 判断 request 的 user_id
fm.Response().MaskedName()  // 判断 response 的 name
fm.Response().MaskedAddress()  // 判断 response 的 address

resp := new(UserInfoResponse)
fm.Response().Apply(resp)
```

## 举例说明

以如下的 proto 定义为例：

```protobuf
// user.proto
syntax = "proto3";

package example;

import "google/protobuf/field_mask.proto";
import "protoc_gen_fieldmask/option.proto";

message UserInfoRequest {
  string user_id = 1;
  google.protobuf.FieldMask fm = 2;
}

message UserInfoResponse {
  string user_id = 1;
  string name = 2;
  string email = 3 [(protoc_gen_fieldmask.field) = {mask: true}];
  Address address = 4 [(protoc_gen_fieldmask.field) = {mask: true, nested: true}];
}

message Address {
  string country = 1;
  string province = 2 [(protoc_gen_fieldmask.field) = {mask: true}]; // 不针对这个字段生成 Mask API
}

service UserService {
  rpc GetUserInfo(UserInfoRequest) returns (UserInfoResponse) {
    option (protoc_gen_fieldmask.rpc) = {
      field_name: "fm"
      mode: FILTER
    };
  }
}
```

最终应该生成如下的代码：`user.pb.fm.go`

```go
// Code generated by protoc-gen-fieldmask. DO NOT EDIT.
// versions:
//  protoc-gen-fieldmask v0.4.1
//  source: user.proto
package example

import (
	pgfmpb "github.com/yeqown/protoc-gen-fieldmask/protobuf"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

func (x *UserInfoRequest) FieldMask() *GetUserInfo_FieldMask {
	// 1. 根据配置 mode 设定 mask mode
	// 2. 从 request 中获取 FieldMask, 根据配置的 field_name 进行解析
	return newGetUserInfo_FieldMask(pgfmpb.MaskMode_FILTER, x.GetFm())
}

type GetUserInfo_FieldMask struct {
	mode pgfmpb.MaskMode
	mask *fieldmaskpb.FieldMask

	req     map[string]struct{}
	res     map[string]struct{}
}

func newGetUserInfo_FieldMask(mode pgfmpb.MaskMode, fieldmask *fieldmaskpb.FieldMask) *GetUserInfo_FieldMask {
	fm := &GetUserInfo_FieldMask{
		mode: mode,
		mask: fieldmask,
		req:  nil,
		res:  nil,
	}

	// 根据 fieldmask 来填充 req 和 res
	if fieldmask == nil || len(fieldmask.GetPaths()) == 0 {
		return fm
	}
	
	fm.req, fm.res = pgfmpb.SplitPaths(fieldmask)

	return fm
}

func (fm *GetUserInfo_FieldMask) Request() *GetUserInfo_Req {
	if fm.req == nil {
		fm.req = make(map[string]struct{}, 4)
	}
	
	return &GetUserInfo_Req{fm: fm}
}

func (fm *GetUserInfo_FieldMask) Response() *GetUserInfo_Res {
	if fm.res == nil {
		fm.res = make(map[string]struct{}, 4)
	}
	
	return &GetUserInfo_Res{fm: fm}
}

func (fm *GetUserInfo_FieldMask) setReqFieldPath(path string) { fm.req[path] = struct{}{} }
func (fm *GetUserInfo_FieldMask) setResFieldPath(path string) { fm.res[path] = struct{}{} }

func (fm *GetUserInfo_FieldMask) hasReqFieldPath(path string) bool {
	if fm == nil || len(fm.req) == 0 {
		return false
	}

	_, exists := fm.req[path]
	return exists
}

func (fm *GetUserInfo_FieldMask) hasResFieldPath(path string) bool {
	if fm == nil || len(fm.res) == 0 {
		return false
	}

	_, exists := fm.res[path]
	return exists
}

type GetUserInfo_Req struct{ fm *GetUserInfo_FieldMask }

func (r *GetUserInfo_Req) MaskUserId() { r.fm.setReqFieldPath("user_id") }

func (r *GetUserInfo_Req) MaskedUserId() bool { return r.fm.hasReqFieldPath("user_id") }

type GetUserInfo_Res struct{ fm *GetUserInfo_FieldMask }

func (r *GetUserInfo_Res) MaskEmail()           { r.fm.setResFieldPath("email") }
func (r *GetUserInfo_Res) MaskAddress()         { r.fm.setResFieldPath("address") }
func (r *GetUserInfo_Res) MaskAddress_Country() { r.fm.setResFieldPath("address.country") }

func (r *GetUserInfo_Res) MaskedEmail() bool           { return r.fm.hasResFieldPath("email") }
func (r *GetUserInfo_Res) MaskedAddress() bool         { return r.fm.hasResFieldPath("address") }
func (r *GetUserInfo_Res) MaskedAddress_Country() bool { return r.fm.hasResFieldPath("address.country") }

// 根据 mode 和 fieldmask 来应用 mask，使用 proto 反射来实现
func (r *GetUserInfo_Res) Apply(resp proto.Message) {
	if r.fm.mode == pgfmpb.MaskMode_PRUNE {
		pgfmpb.Prune(resp, r.fm.res)
		return
	}

	pgfmpb.Filter(resp, r.fm.res)
}
```