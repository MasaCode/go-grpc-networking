package main;

import (
	"context"
	"fmt"
	"io"
	"log"
	"google.golang.org/grpc"
	"go-grpc-networking/greet/pb"
)

func main () {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: $v", err)
	}
	defer conn.Close()

	c := greet.NewGreetServiceClient(conn)
	// fmt.Printf("Created client: %f", c)

	//doUnary(c)
	doServerStreaming(c)
}

func doServerStreaming (client greet.GreetServiceClient) {
	fmt.Println("Starting to do a server streaming RPC...")
	request := &greet.GreetManyTimesRequest{
		Greeting: &greet.Greeting{
			FirstName: "Masashi",
			LastName:  "Morita",
		},
	}
	stream, err := client.GreetManyTimes(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receiving stream from server: %v", err)
		}
		fmt.Printf("Server Response: %v\n", res.Result)
	}
}

func doUnary (client greet.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	request := &greet.GreetRequest{
		Greeting: &greet.Greeting{
			FirstName: "Stephane",
			LastName: "Bob",
		},
	}
	res, err := client.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}
