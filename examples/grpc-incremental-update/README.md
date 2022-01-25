## incremental update

Client masked it's fields to let the server known those fields are changed or not, 
so server can only update the fields by easily calling `request.MaskedIn_xxx`.

```go
maskedReq := &pb.UserInfoRequest{
	UserId: "1",
}

// output (server-side):
// userServer.UpdateUserInfo is called: user_id:"1"

// currently, server in filter mode, so masked fields means fields are changed and needed to be updated.
// uncomment the following line to see the effect of filter mode.
//
// maskedReq.MaskIn_UserId() // masking user_id field in request

// output (server-side):
// userServer.UpdateUserInfo is called: user_id:"1"  field_mask:{paths:"user_id"}
// userId want to be updated.

resp, err := client.UpdateUserInfo(context.Background(), maskedReq)
```

### running example

```sh
# running server
go run examples/grpc-masked-response/server.go
# output:
# userServer.UpdateUserInfo is called: user_id:"1"
#
# with uncommenting the following line in client.go:
# maskedReq.MaskIn_UserId() // masking user_id field in request.
# output:
# userServer.UpdateUserInfo is called: user_id:"1"  field_mask:{paths:"user_id"}
# userId want to be updated.

# running client in another terminal session
go run client.go
```