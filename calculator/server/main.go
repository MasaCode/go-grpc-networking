package main;

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
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

func (*server) ComputeAverage(stream calculator.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Server was called with: a stream request \n")
	sum := int32(0)
	count := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculator.ComputeAverageResponse{
				Result: average,
			})
		}
		if err != nil {
			log.Fatalf("error while receiving client request: %v\n", err)
		}
		sum += req.GetNumber()
		count++
	}
}

func (*server) FindMaximum(stream calculator.CalculatorService_FindMaximumServer) error {
	fmt.Printf("Server was called with: a stream request \n")
	var currentMax int32 = 0
	for {
		req, err := stream.Recv()
		if  err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while fetching client request: %v\n", err)
			return err
		}

		// set currentMax only if the requested number is bigger than current
		if currentMax < req.GetNumber() {
			currentMax = req.GetNumber()
			sendErr := stream.Send(&calculator.FindMaximumResponse{
				Result: currentMax,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending response: %v\n", sendErr)
				return sendErr
			}
		}
	}
	return nil
}

func (*server) SquareRoot(ctx context.Context, req *calculator.SquareRootRequest) (*calculator.SquareRootResponse, error){
	fmt.Printf("Server was called with: %v\n", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &calculator.SquareRootResponse{
		Result: math.Sqrt(float64(number)),
	}, nil
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
