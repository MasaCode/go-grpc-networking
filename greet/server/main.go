package main;

import (
	"context"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	"go-grpc-networking/greet/pb"
)

type server struct {}

func (*server) Greet(ctx context.Context, req *greet.GreetRequest,) (*greet.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v \n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "HELLO " + firstName
	res := &greet.GreetResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello World\n")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listedn: %v", err)
	}

	s := grpc.NewServer()
	greet.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
