package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"unary2/protoc"

	"google.golang.org/grpc"
)

type server struct {
	protoc.UnimplementedExampleServer
}

func (s *server) ServerReply(ctx context.Context, req *protoc.HelloRequest) (*protoc.HelloResponse, error) {
	fmt.Printf("Received Client Request : %v", req.GetMessage())

	reply := "Server says : Got your message" + req.GetMessage()

	result := &protoc.HelloResponse{
		Reply: reply,
	}

	return result, nil

}

func main() {
	fmt.Println("Starting the server")

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()

	protoc.RegisterExampleServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)

	}

}
