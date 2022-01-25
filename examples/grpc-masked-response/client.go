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
	maskedReq.MaskOut_Email(). // masking email field in response
					MaskOut_Address() // masking address field in response
	fmt.Printf("request: %+v\n", maskedReq.String())

	resp, err := client.GetUserInfo(context.Background(), maskedReq)
	if err != nil {
		panic(err)
	}

	// Output:
	// email:"yeqown@gmail.com" address:{country:"China" province:"Sichuan"}
	fmt.Printf("response: %+v\n", resp.String())
}
