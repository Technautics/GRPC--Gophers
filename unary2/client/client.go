package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"unary2/protoc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Connecting to server....")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create new client connection
	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create a client connection....")
	}
	defer conn.Close()

	// Create client from generated code
	client := protoc.NewExampleClient(conn)

	doUnary(client)

}

func doUnary(c protoc.ExampleClient) {
	// Create request with a simple message
	req := &protoc.HelloRequest{
		Message: "Hey!, this is Mushahid.....",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Make the Unary RPC Call
	res, err := c.ServerReply(ctx, req)
	if err != nil {
		log.Fatalf("Error during displaying the message...")
	}

	// Print the reply from the server
	fmt.Println("Server Replied : ", res.GetReply())

}
