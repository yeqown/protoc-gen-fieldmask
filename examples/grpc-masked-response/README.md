## masked response

Client sends a request to the server, which returns a masked response with judgement that
depends on the field `request.FieldMask`.

```go
maskedReq := &pb.UserInfoRequest{
	UserId: "1",
}
maskedReq.MaskOut_Email(). // masking email field in response
          MaskOut_Address() // masking address field in response
fmt.Printf("request: %+v\n", maskedReq.String())
// output:
// request: user_id:"1" field_mask:{paths:"email" paths:"address"}
resp, err := client.GetUserInfo(context.Background(), maskedReq)
// output: 
// response: email:"yeqown@gmail.com" address:{country:"China" province:"Sichuan"}
```

### running example

```sh
$ go run server.go

# running client in another terminal session
$ go run client.go
$ request: user_id:"1" field_mask:{paths:"email" paths:"address"}
$ response: email:"yeqown@gmail.com" address:{country:"China" province:"Sichuan"}
```