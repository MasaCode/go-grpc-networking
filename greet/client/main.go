package main;

import (
	"context"
	"fmt"
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
	doUnary(c)
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
