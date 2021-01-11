package main

import (
	"context"
	"fmt"
	"go-grpc-networking/greet/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
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

func (*server) GreetWithDeadline(ctx context.Context, req *greet.GreetWithDeadlineRequest) (*greet.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function was invoked with %v \n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "the client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()
	result := "HELLO " + firstName
	res := &greet.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
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
