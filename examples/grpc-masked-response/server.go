package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"examples/pb"
)

type userServer struct{}

func (u userServer) GetUserInfo(
	ctx context.Context,
	request *pb.UserInfoRequest,
) (*pb.UserInfoResponse, error) {
	// FieldMask_Filter means that the fields are included in the response,
	// otherwise they are omitted.
	// FieldMask_Prune means masked fields are not included in the response,
	// otherwise they are included.
	fm := request.FieldMask_Filter()
	// fm2 := request.FieldMask_Prune()

	fmt.Printf("userServer.GetUserInfo is called: %v\n", request.String())

	resp := &pb.UserInfoResponse{
		UserId:  "69781",
		Name:    "yeqown",
		Email:   "yeqown@gmail.com",
		Address: nil,
	}

	// judge if the field masked or not, avoid unnecessary call or calculation.
	if fm.MaskedOut_Address() {
		// filter more, so the address field should be included.
		resp.Address = &pb.Address{
			Country:  "China",
			Province: "Sichuan",
		}
	}

	// filter the field masked out.
	_ = fm.Mask(resp)
	return resp, nil
}

func (u userServer) UpdateUserInfo(ctx context.Context, request *pb.UserInfoRequest) (*pb.NonEmpty, error) {
	// FieldMask_Filter means that the fields are expected to update,
	// otherwise they are ignored.
	// FieldMask_Prune means masked fields are ignored to update,
	// otherwise they are expected to update.
	fm := request.FieldMask_Filter()
	// fm2 := request.FieldMask_Prune()

	fmt.Printf("userServer.UpdateUserInfo is called: %v\n", request.String())

	if fm.MaskedIn_UserId() {
		// userId want to be updated, so you should use request.UserId to update.
		fmt.Println("userId want to be updated.")
	}

	return new(pb.NonEmpty), nil
}

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	server := new(userServer)
	pb.RegisterUserInfoServer(s, server)

	fmt.Println("server is running..., :8080")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
