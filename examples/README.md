# FieldMask 使用示例

本目录展示了 protoc-gen-fieldmask 的两种主要使用场景：

## 目录结构

```
examples/
├── proto/
│   └── user.proto      # 用户服务定义
├── server/
│   └── server.go       # 服务器端实现
├── client/
│   └── client.go       # 客户端实现示例
└── README.md           # 本文档
```

---

## 场景 1: 部分响应 (FILTER 模式)

### 用途

客户端指定需要返回哪些字段，服务器只返回被标记的字段。适用于：

- **减少响应体积** - 移动端只需要少量字段
- **保护敏感数据** - 密码哈希等字段永不返回
- **按需加载** - 不同客户端需要不同字段组合

### 实现原理

**服务器端**：
```go
func (s *Server) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
    // 1. 获取完整用户数据
    resp := &GetUserResponse{
        UserId:       user.Id,
        Username:     user.Name,
        Email:        user.Email,
        PasswordHash: "secret", // 敏感数据
        // ... 其他字段
    }

    // 2. 解析 FieldMask
    _, respMask := protobuf.SplitPaths(req.GetFm().GetPaths())

    // 3. 应用过滤器：只保留被标记的字段
    protobuf.Filter(resp, respMask)

    return resp, nil
}
```

**客户端（移动端 - 只需要 ID 和用户名）**：
```go
req := &GetUserRequest{UserId: "user123"}
fm := req.FieldMask()

// 只标记需要的字段
fm.Response().MaskUserId()
fm.Response().MaskUsername()

// 调用 RPC - 只会返回 user_id 和 username
resp, _ := client.GetUser(ctx, req)
```

**客户端（Web 端 - 需要除密码外的所有字段）**：
```go
req := &GetUserRequest{UserId: "user123"}
fm := req.FieldMask()

// 标记所有需要的字段（不标记 password_hash）
fm.Response().MaskUserId()
fm.Response().MaskUsername()
fm.Response().MaskEmail()
fm.Response().MaskPhone()
// ... 其他非敏感字段

// 调用 RPC - 返回所有被标记的字段，password_hash 不会被返回
resp, _ := client.GetUser(ctx, req)
```

### 生成的 FieldMask 路径

```
req.fm.paths = ["res.user_id", "res.username"]  // 移动端示例
```

---

## 场景 2: 部分更新 (PRUNE 模式)

### 用途

客户端指定**排除**哪些字段的更新。被标记的字段**不会**被更新，未被标记的字段会被更新。适用于：

- **增量更新** - 只更新修改过的字段
- **防止误更新** - 某些字段不应被意外修改

### 实现原理

**服务器端**：
```go
func (s *Server) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error) {
    user, _ := s.getUserFromStore(req.UserId)

    // 解析排除列表
    reqMask, _ := protobuf.SplitPaths(req.GetFm().GetPaths())

    // 只更新未被标记（未被排除）的字段
    if !protobuf.IsFieldMarked(reqMask, "username") && req.Username != "" {
        user.Username = req.Username  // 会更新
    }
    if !protobuf.IsFieldMarked(reqMask, "email") && req.Email != "" {
        user.Email = req.Email  // 会更新
    }
    // ... 其他字段

    return resp, nil
}
```

**客户端（更新个人资料 - 排除联系方式）**：
```go
req := &UpdateUserRequest{
    UserId:    "user123",
    Username:  "new_name",
    Email:     "should_not_change@example.com", // 不会被更新
    Phone:     "+9999999999",                   // 不会被更新
    AvatarUrl: "https://example.com/new.jpg",
    Bio:       "Updated bio",
}
fm := req.FieldMask()

// 标记要排除的字段（这些字段不会被更新）
fm.Request().MaskEmail()
fm.Request().MaskPhone()

// 调用 RPC - 只有 username、avatar_url、bio 会被更新
resp, _ := client.UpdateUser(ctx, req)
```

**客户端（只更新联系方式 - 排除个人资料）**：
```go
req := &UpdateUserRequest{
    UserId:    "user123",
    Username:  "should_not_change", // 不会被更新
    Email:     "new@example.com",
    Phone:     "+1111111111",
    AvatarUrl: "should_not_change", // 不会被更新
    Bio:       "should_not_change", // 不会被更新
}
fm := req.FieldMask()

// 标记要排除的字段
fm.Request().MaskUsername()
fm.Request().MaskAvatarUrl()
fm.Request().MaskBio()

// 调用 RPC - 只有 email、phone 会被更新
resp, _ := client.UpdateUser(ctx, req)
```

### 生成的 FieldMask 路径

```
req.fm.paths = ["req.email", "req.phone"]  // 排除这些字段的更新
```

---

## 与标准 PATCH 的区别

### 标准 PATCH (典型实现)
```protobuf
// 标准方式：标记要更新的字段
message UpdateRequest {
    string user_id = 1;
    string username = 2;
    string email = 3;
    google.protobuf.FieldMask update_mask = 4;  // 指定要更新的字段
}

// 服务器逻辑
if updateMask.HasPath("username") {
    user.Username = req.Username  // 只更新在 mask 中的字段
}
```

### PRUNE 模式 (本插件)
```protobuf
// PRUNE 方式：标记要排除的字段
message UpdateRequest {
    string user_id = 1;
    string username = 2;
    string email = 3;
    google.protobuf.FieldMask fm = 4;  // 指定要排除的字段
}

// 服务器逻辑
if !IsFieldMarked(reqMask, "username") {
    user.Username = req.Username  // 更新不在 mask 中的字段
}
```

---

## 代码生成

生成示例代码：

```bash
cd examples
protoc \
    -I=. \
    -I=../third_party \
    --go_out=paths=source_relative:. \
    --fieldmask_out=paths=source_relative,lang=go:. \
    ./proto/user.proto
```

这会生成：
- `proto/user.pb.go` - 标准 protobuf 生成的代码
- `proto/user.pb.fm.go` - FieldMask 辅助代码

---

## 运行示例

```bash
# 生成代码
make gen-examples

# 运行测试（需要先实现真实的服务器/客户端）
go run server/main.go
go run client/main.go
```

---

## 生成的 API

### GetUserRequest
```go
req := &GetUserRequest{UserId: "user123"}
fm := req.FieldMask()

// 标记响应字段
fm.Response().MaskUserId()
fm.Response().MaskUsername()
// ... 其他字段

// 检查字段是否被标记
if fm.Response().MaskedUserId() {
    // user_id 被标记
}
```

### UpdateUserRequest
```go
req := &UpdateUserRequest{UserId: "user123", ...}
fm := req.FieldMask()

// 标记请求字段（排除这些字段）
fm.Request().MaskEmail()
fm.Request().MaskPhone()

// 检查字段是否被标记
if fm.Request().MaskedEmail() {
    // email 被标记（将被排除）
}
```
