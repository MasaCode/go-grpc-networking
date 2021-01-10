package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"google.golang.org/grpc"
	"go-grpc-networking/greet/pb"
	"strconv"
	"time"
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

func (*server) GreetManyTimes(req *greet.GreetManyTimesRequest, stream greet.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v \n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "HELLO " + firstName + " number " + strconv.Itoa(i)
		res := &greet.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greet.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request \n")
	result := "HELLO "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greet.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v\n", err)
		}
		firstName := req.GetGreeting().FirstName
		result += firstName + "! "
	}
}

func (*server) GreetEveryone(stream greet.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with a streaming request \n")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while receving client request: %v\n", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		sendErr := stream.Send(&greet.GreetEveryoneResponse{
			Result: "HELLO " + firstName,
		})
		if sendErr != nil {
			log.Fatalf("error while sending response: %v\n", sendErr)
			return sendErr
		}
	}
	return nil
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
