package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"serverStream/news"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("📡 Connecting to server......")

	// Connect to gRPC server on localhost 50051
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌Could not connect.. %v", err)
	}
	defer conn.Close()

	// Create a client for news service
	client := news.NewNewsServiceClient(conn)

	// Make a request to the server
	req := &news.NewsRequest{
		Category: "technology",
	}

	// Call server streaming RPC
	stream, err := client.GetNewsStream(context.Background(), req)
	if err != nil {
		log.Fatalf("❌ Error while calling GetNewsStream: %v", err)
	}

	// Receive the streamed responses one by one
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("📭 End of stream.")
			break
		}
		if err != nil {
			log.Fatalf("❌ Error receiving: %v", err)
		}

		log.Printf("📰 Headline received: %s", res.Headline)
		time.Sleep(500 * time.Millisecond)
	}

}
