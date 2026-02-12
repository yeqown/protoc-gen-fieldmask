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
  FILTER = 0;  // 被 mask 的字段才有效，其他字段会被忽略
  PRUNE = 1;   // 被 mask 的字段会被忽略，其他字段会被返回
}

// ============================================================================
// Field-level options: 精细控制字段的 mask 行为
// ============================================================================

extend google.protobuf.FieldOptions {
  optional FieldOptions field = 1143;
}

// FieldOptions 定义单个字段的 mask 行为
message FieldOptions {
  // 是否忽略该字段的 mask 方法生成
  bool ignore = 1;
  
  // 是否为嵌套字段生成 mask 方法
  bool nested = 2;
}
```

## 使用示例

### CASE 1: 基本用法

```protobuf
syntax = "proto3";

package example;

import "google/protobuf/field_mask.proto";
import "protoc-gen-fieldmask/options.proto";

message UserInfoRequest {
  string user_id = 1;
  google.protobuf.FieldMask fm = 2;
}

message UserInfoResponse {
  string user_id = 1;
  string name = 2;
  string email = 3 [(protoc_gen_fieldmask.field) = {ignore: true}];
  Address address = 4 [(protoc_gen_fieldmask.field) = {nested: true}];
}

message Address {
  string country = 1;
  string province = 2 [(protoc_gen_fieldmask.field) = {ignore: true}]; // 不针对这个字段生成 Mask API
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

```protobuf
message UserListRequest {
  int32 page_size = 1;
  string page_token = 2;
  google.protobuf.FieldMask fm = 3;
}

service UserListService {
  // 复用 UserInfoResponse，无需 _FM_NEVER_USE hack
  rpc ListUsers(UserListRequest) returns (UserInfoResponse) {
    option (protoc_gen_fieldmask.rpc) = {
      field_name: "fm"
      mask_mode: FILTER
    };
  }
}
```

### CASE 3: 增量更新 (PRUNE 模式)

```protobuf
message UpdateUserRequest {
  string user_id = 1;
  string name = 2;
  string email = 3;
  google.protobuf.FieldMask update_mask = 4;
}

service UserUpdateService {
  rpc UpdateUser(UpdateUserRequest) returns (UserInfoResponse) {
    option (protoc_gen_fieldmask.rpc) = {
      field_name: "update_mask"
      mask_mode: PRUNE
    };
  }
}
```

## 预期生成代码

### API 设计概览

```go
// 统一入口，指定 mask 模式
fm := req.FieldMask(FILTER)  // 或 PRUNE

// 标记 request 自身字段
fm.Request().UserId()
fm.Request().Name()

// 标记 response 字段
fm.Response().Name()
fm.Response().Address().Country()

// 判断字段是否被标记
fm.Marked().UserId()           // 判断 request 的 user_id
fm.Marked().Response().Name()  // 判断 response 的 name

// 应用 mask
fm.Request().Apply(req)   // 应用到 request
fm.Response().Apply(resp) // 应用到 response
```

---

## 生成 API 详细说明

### 1. Request Message 生成的 API

Request message 包含 `google.protobuf.FieldMask` 字段，且在 RPC method 上配置了 `rpc` option。

#### 1.1 入口方法

```go
// UserInfoRequest 生成
func (m *UserInfoRequest) FieldMask() *UserInfoRequest_FieldMask
```

#### 1.2 Request 自身字段操作

```go
// 获取 Request 字段操作对象
func (fm *UserInfoRequest_FieldMask) Request() *UserInfoRequest_RequestMask

// 标记 request 的基础字段
func (rm *UserInfoRequest_RequestMask) UserId() *UserInfoRequest_RequestMask
func (rm *UserInfoRequest_RequestMask) Name() *UserInfoRequest_RequestMask

// 标记 request 的嵌套字段（如果字段标记了 nested: true）
// Address() 标记整个 address 字段
func (rm *UserInfoRequest_RequestMask) Address() *UserInfoRequest_RequestMask

// 进入 Address 的子字段操作
func (a *Address) FieldMask() *Address_FieldMask

// Address 的子字段
func (afm *Address_FieldMask) Country() *Address_FieldMask
func (afm *Address_FieldMask) Province() *Address_FieldMask

// 应用 mask 到 request
func (rm *UserInfoRequest_RequestMask) Apply(req *UserInfoRequest)
```

#### 1.3 Response 字段操作

```go
// 获取 Response 字段操作对象
func (fm *UserInfoRequest_FieldMask) Response() *UserInfoRequest_ResponseMask

// 标记 response 的基础字段
func (rm *UserInfoRequest_ResponseMask) UserId() *UserInfoRequest_ResponseMask
func (rm *UserInfoRequest_ResponseMask) Name() *UserInfoRequest_ResponseMask

// 标记 response 的嵌套字段（如果字段标记了 nested: true）
// Address() 标记整个 address 字段
func (rm *UserInfoRequest_ResponseMask) Address() *UserInfoRequest_ResponseMask

// 进入 Address 的子字段操作
func (rm *UserInfoRequest_ResponseMask) Address() *UserInfoRequest_ResponseMask_Address
func (a *UserInfoRequest_ResponseMask_Address) FieldMask() *Address_FieldMask

// Address 的子字段
func (afm *Address_FieldMask) Country() *Address_FieldMask
func (afm *Address_FieldMask) Province() *Address_FieldMask

// 应用 mask 到 response
func (rm *UserInfoRequest_ResponseMask) Apply(resp *UserInfoResponse)
```

#### 1.4 字段判断

```go
// 获取判断对象
func (fm *UserInfoRequest_FieldMask) Marked() *UserInfoRequest_Marked

// 判断 request 自身字段
func (m *UserInfoRequest_Marked) UserId() bool
func (m *UserInfoRequest_Marked) Name() bool
func (m *UserInfoRequest_Marked) Address() bool

// 判断 response 字段
func (m *UserInfoRequest_Marked) Response() *UserInfoRequest_Marked_Response
func (mr *UserInfoRequest_Marked_Response) UserId() bool
func (mr *UserInfoRequest_Marked_Response) Name() bool
func (mr *UserInfoRequest_Marked_Response) Address() bool
```

---

### 2. Response Message 生成的 API

Response message 本身不包含 FieldMask 字段，但被 RPC method 引用。

#### 2.1 独立的 Mask 操作（可选）

如果需要在 response 上直接操作（不通过 request），可以生成：

```go
// UserInfoResponse 生成独立的 mask 方法
func (m *UserInfoResponse) NewFieldMask(mode MaskMode) *UserInfoResponse_FieldMask

// 标记字段
func (fm *UserInfoResponse_FieldMask) UserId() *UserInfoResponse_FieldMask
func (fm *UserInfoResponse_FieldMask) Name() *UserInfoResponse_FieldMask
func (fm *UserInfoResponse_FieldMask) Address() *UserInfoResponse_FieldMask_Address

// 应用 mask
func (fm *UserInfoResponse_FieldMask) Apply(resp *UserInfoResponse)

// 判断字段
func (fm *UserInfoResponse_FieldMask) Marked() *UserInfoResponse_Marked
func (m *UserInfoResponse_Marked) UserId() bool
func (m *UserInfoResponse_Marked) Name() bool
```

**注意**: Response 的独立 API 是可选的，主要用于：
- 在没有 request 的情况下操作 response
- 跨多个 request 共享 response mask 逻辑

---

### 3. 普通 Message 生成的 API

普通 message（如 `Address`）如果被标记为 `nested: true`，会生成独立的 FieldMask 操作。

#### 3.1 作为嵌套字段使用

```go
// Address 作为 UserInfoResponse.address 的嵌套字段
// 在 Response mask 中生成

// 标记整个 address 字段
func (rm *UserInfoRequest_ResponseMask) Address() *UserInfoRequest_ResponseMask

// 获取 Address 的 FieldMask 操作对象
func (rm *UserInfoRequest_ResponseMask) Address() *UserInfoRequest_ResponseMask_Address
func (a *UserInfoRequest_ResponseMask_Address) FieldMask() *Address_FieldMask

// Address 的字段方法
func (afm *Address_FieldMask) Country() *Address_FieldMask
func (afm *Address_FieldMask) Province() *Address_FieldMask
```

**使用示例**:
```go
fm := req.FieldMask(FILTER)

// 标记整个 address
fm.Response().Address()

// 标记 address 的子字段
fm2 := fm.Response().Address().FieldMask()
fm2.Country()
fm2.Province()
```

#### 3.2 独立使用（可选）

如果 `Address` 需要独立的 mask 功能：

```go
// Address 生成独立的 mask 方法
func (m *Address) NewFieldMask(mode MaskMode) *Address_FieldMask

func (fm *Address_FieldMask) Country() *Address_FieldMask
func (fm *Address_FieldMask) Province() *Address_FieldMask

func (fm *Address_FieldMask) Apply(addr *Address)

func (fm *Address_FieldMask) Marked() *Address_Marked
func (m *Address_Marked) Country() bool
func (m *Address_Marked) Province() bool
```

---

### 4. 字段级别的控制

#### 4.1 ignore: true

```protobuf
message UserInfoResponse {
  string email = 3 [(protoc_gen_fieldmask.field) = {ignore: true}];
}
```

**效果**: 不生成 `Email()` 相关方法

#### 4.2 nested: true

```protobuf
message UserInfoResponse {
  Address address = 4 [(protoc_gen_fieldmask.field) = {nested: true}];
}
```

**效果**: 生成嵌套字段的 FieldMask 操作

```go
// 标记整个 address
func (rm *UserInfoRequest_ResponseMask) Address() *UserInfoRequest_ResponseMask

// 进入 Address 子字段操作
func (rm *UserInfoRequest_ResponseMask) Address() *UserInfoRequest_ResponseMask_Address
func (a *UserInfoRequest_ResponseMask_Address) FieldMask() *Address_FieldMask

// Address 的子字段
func (afm *Address_FieldMask) Country() *Address_FieldMask
func (afm *Address_FieldMask) Province() *Address_FieldMask
```

**使用**:
```go
fm.Response().Address()  // 标记整个 address

fm2 := fm.Response().Address().FieldMask()
fm2.Country()  // 标记 address.country
```

#### 4.3 默认行为（无 option）

```protobuf
message UserInfoResponse {
  string name = 2;  // 无 option
}
```

**效果**: 生成基础的 mask 方法

```go
func (rm *UserInfoRequest_ResponseMask) Name() *UserInfoRequest_ResponseMask
```

---

### 5. 完整示例对比

#### Proto 定义

```protobuf
message UserInfoRequest {
  string user_id = 1;
  google.protobuf.FieldMask fm = 2;
}

message UserInfoResponse {
  string user_id = 1;
  string name = 2;
  string email = 3 [(protoc_gen_fieldmask.field) = {ignore: true}];
  Address address = 4 [(protoc_gen_fieldmask.field) = {nested: true}];
}

message Address {
  string country = 1;
  string province = 2;
}

service UserService {
  rpc GetUserInfo(UserInfoRequest) returns (UserInfoResponse) {
    option (protoc_gen_fieldmask.rpc) = {
      field_name: "fm"
      mask_mode: FILTER
    };
  }
}
```

---

## 解决的问题

1. **配置简洁**: 只保留必要的配置项，去除冗余
2. **跨包引用**: 通过 method options 自动关联 request/response，无需 hack
3. **字段级控制**: 通过 `ignore` 和 `nested` 精细控制字段行为
4. **双模式支持**: FILTER 和 PRUNE 模式满足不同场景
5. **命名清晰**: `MaskOut_*` 表示输出字段，`Masked_*` 表示判断是否被 mask
