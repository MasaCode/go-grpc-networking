package main;

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"go-grpc-networking/calculator/pb"
	"log"
)

func main() {
	fmt.Println("Hello, I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error while connecting to tcp: $v", err)
	}
	defer conn.Close()

	c := calculator.NewSumServiceClient(conn)
	doUnary(c)
}

func doUnary (client calculator.SumServiceClient) {
	request := &calculator.SumRequest{
		Numbers: &calculator.SumNumbers{
			A: 4,
			B: 5,
		},
	}
	result, err := client.Sum(context.Background(), request)
	if err != nil {
		log.Fatalf("error while invoking rpc function: %v", err)
	}
	fmt.Printf("Result: %v\n", result)
}
