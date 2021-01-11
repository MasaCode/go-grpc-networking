package main;

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"go-grpc-networking/calculator/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello, I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error while connecting to tcp: $v", err)
	}
	defer conn.Close()

	c := calculator.NewCalculatorServiceClient(conn)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDirectionalStreaming(c)

	doErrorUnary(c)
}

func doErrorUnary (client calculator.CalculatorServiceClient) {
	doErrorCall(client, 10)
	doErrorCall(client, -10)
}

func doErrorCall (client calculator.CalculatorServiceClient, number int32) {
	res, err := client.SquareRoot(context.Background(), &calculator.SquareRootRequest{Number: number,})
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			fmt.Println(resErr.Message())
			fmt.Println(resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
			return
		} else {
			log.Fatalf("Big Error callling SquareRoot: %v\n", err)
			return
		}
	}

	fmt.Printf("Result of Square root of %v is %v\n", number, res)
}

func doBiDirectionalStreaming (client calculator.CalculatorServiceClient) {
	stream, err := client.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while opening stream: %v\n", err)
	}

	waitc := make(chan struct{})

	go func () {
		numbers := []int32{2, 5, 100, 3, 4, 200}
		for _, number := range numbers {
			fmt.Printf("Sending Request: %v\n", number)
			err := stream.Send(&calculator.FindMaximumRequest{
				Number: number,
			})
			if err != nil {
				log.Fatalf("error while sending request: %v\n", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("error while closing connection: %v\n", err)
		}
	}()

	go func () {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving response: %v\n", err)
			}
			fmt.Printf("Server Response: %v\n", res.Result)
		}
		close(waitc)
	}()

	<-waitc
}

func doClientStreaming (client calculator.CalculatorServiceClient) {
	requests := []*calculator.ComputeAverageRequest{
		&calculator.ComputeAverageRequest{
			Number: 1,
		},
		&calculator.ComputeAverageRequest{
			Number: 2,
		},
		&calculator.ComputeAverageRequest{
			Number: 3,
		},
		&calculator.ComputeAverageRequest{
			Number: 4,
		},
	}
	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while invoking ComputeAverage RPC method: %v\n", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending Request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving server response: %v\n", err)
	}
	fmt.Printf("\n\nServer Response: %v\n", res.Result)
}

func doServerStreaming (client calculator.CalculatorServiceClient) {
	request := &calculator.PrimeNumberDecompositionRequest{
		Number: 120,
	}
	stream, err := client.PrimeNumberDecomposition(context.Background(), request)
	if err != nil {
		log.Fatalf("error while invoking rpc function: %v\n", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receiving response from server: $v\n", err)
		}
		fmt.Printf("Server Response: %v\n", res.Result)
	}
}

func doUnary (client calculator.CalculatorServiceClient) {
	request := &calculator.SumRequest{
		FirstNumber: 40,
		SecondNumber: 5,
	}
	result, err := client.Sum(context.Background(), request)
	if err != nil {
		log.Fatalf("error while invoking rpc function: %v", err)
	}
	fmt.Printf("Result: %v\n", result)
}
