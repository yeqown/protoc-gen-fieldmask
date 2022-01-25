package main

import (
	"context"
	"fmt"

	"examples/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		panic(err)
	}

	client := pb.NewUserInfoClient(cc)
	defer cc.Close()

	maskedReq := &pb.UserInfoRequest{
		UserId: "1",
	}

	// output:
	// userServer.UpdateUserInfo is called: user_id:"1"

	// currently, server in filter mode, so masked fields means fields are changed and needed to be updated.
	// uncomment the following line to see the effect of filter mode.
	//
	// maskedReq.MaskIn_UserId() // masking user_id field in request

	// output:
	// userServer.UpdateUserInfo is called: user_id:"1"  field_mask:{paths:"user_id"}
	// userId want to be updated.

	fmt.Printf("request: %+v\n", maskedReq.String())

	resp, err := client.UpdateUserInfo(context.Background(), maskedReq)
	if err != nil {
		panic(err)
	}
	_ = resp
}
