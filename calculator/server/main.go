package main;

import (
	"context"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	"go-grpc-networking/calculator/pb"
)

type server struct {}
func (*server) Sum (ctx context.Context, req *calculator.SumRequest) (*calculator.SumResponse, error) {
	a := req.GetNumbers().GetA()
	b := req.GetNumbers().GetB()
	result := a + b
	res := &calculator.SumResponse{
		Result: result,
	}
	return res, nil
}

func main () {
	fmt.Println("Hello world, this is server.")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("error listing tcp connection: %v\n", err)
	}

	s := grpc.NewServer()
	calculator.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error occurred while serving server: %v\n", err)
	}
}
