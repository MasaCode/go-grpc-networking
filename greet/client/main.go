package main;

import (
	"context"
	"fmt"
	"io"
	"log"
	"google.golang.org/grpc"
	"go-grpc-networking/greet/pb"
	"time"
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
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDirectionalStreaming(c)
}

func doBiDirectionalStreaming (client greet.GreetServiceClient) {
	fmt.Println("Starting to do a di-directional streaming RPC...")
	stream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v\n", err)
	}

	waitc := make(chan struct{})

	// send a bunch of message to the client (go routine)
	go func() {
		requests := []*greet.GreetEveryoneRequest{
			&greet.GreetEveryoneRequest{
				Greeting: &greet.Greeting{
					FirstName: "Masashi",
					LastName:  "Morita",
				},
			},
			&greet.GreetEveryoneRequest{
				Greeting: &greet.Greeting{
					FirstName: "Bob",
					LastName:  "Barker",
				},
			},
			&greet.GreetEveryoneRequest{
				Greeting: &greet.Greeting{
					FirstName: "John",
					LastName:  "Doe",
				},
			},
			&greet.GreetEveryoneRequest{
				Greeting: &greet.Greeting{
					FirstName: "Carl",
					LastName:  "Steel",
				},
			},
		}
		for _, req := range requests {
			fmt.Printf("Sending Request: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive a bunch of message from the server (go routine)
	go func() {
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
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func doClientStreaming (client greet.GreetServiceClient) {
	fmt.Println("Starting to do a client streaming RPC...")

	requests := []*greet.LongGreetRequest{
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Masashi",
				LastName:  "Morita",
			},
		},
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Bob",
				LastName:  "Barker",
			},
		},
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Carl",
				LastName:  "Steel",
			},
		},
	}

	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v\n", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receing server response: %v\n", err)
	}
	fmt.Printf("\n\nServer Response: %v\n", res.Result)
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
