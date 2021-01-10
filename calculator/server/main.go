package main;

import (
	"context"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	"go-grpc-networking/calculator/pb"
	"time"
)

type server struct {}
func (*server) Sum (ctx context.Context, req *calculator.SumRequest) (*calculator.SumResponse, error) {
	fmt.Printf("Server was called with: %v\n", req)
	a := req.FirstNumber
	b := req.SecondNumber
	result := a + b
	res := &calculator.SumResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculator.PrimeNumberDecompositionRequest, stream calculator.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Server was called with: %v\n", req)

	number := req.GetNumber()
	var factor int32 = 2
	for number > 1 {
		if (number % factor) == 0 {
			number /=  factor
			stream.Send(&calculator.PrimeNumberDecompositionResponse{
				Result: factor,
			})
			time.Sleep(1000 * time.Millisecond)
		} else {
			factor += 1
		}
	}
	return nil
}

func main () {
	fmt.Println("Hello world, this is server.")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("error listing tcp connection: %v\n", err)
	}

	s := grpc.NewServer()
	calculator.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error occurred while serving server: %v\n", err)
	}
}
