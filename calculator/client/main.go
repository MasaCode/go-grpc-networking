package main;

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"go-grpc-networking/calculator/pb"
	"io"
	"log"
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
	doServerStreaming(c)
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
