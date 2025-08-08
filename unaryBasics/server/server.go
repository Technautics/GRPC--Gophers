package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"unaryBasics/greetpb"

	"google.golang.org/grpc"
)

// Struct to implement greet service server interface
type server struct {
	greetpb.UnimplementedGreetServiceServer
}

// Implement the greet function.\
func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Received GreetRequest: %v\n", req)

	firstName := req.GetFirstName() // Extract name from request
	result := "Hello, " + firstName

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Starting gRPC Server...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to Listen %v", err)
	}

	s := grpc.NewServer() // Create a new grpc server

	greetpb.RegisterGreetServiceServer(s, &server{})
	// Register our service

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)

	}

}
