package main

import (
	"bidirectionalStream/stockpb"
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create StockService client
	client := stockpb.NewStockServiceClient(conn)

	// Create a bidirectional stream
	stream, err := client.StreamPrices(context.Background())
	if err != nil {
		log.Fatalf("❌ Failed to create stream: %v", err)
	}

	// =========================================
	// Goroutine 1: Receiving stock updates from server
	// =========================================
	go func() {
		for {
			// Receive StockPrice from server
			update, err := stream.Recv()
			if err == io.EOF {
				// Server closed stream
				log.Println("📴 Server closed the stream")
				return
			}
			if err != nil {
				// Some error occurred
				log.Fatalf("❌ Error receiving: %v", err)
			}

			// Print stock update to console
			log.Printf("💹 %s -> $%.2f at %s",
				update.Symbol, update.Price, update.Time)
		}
	}()

	// =========================================
	// Main Loop: Sending stock subscriptions to server
	// =========================================
	scanner := bufio.NewScanner(os.Stdin)
	log.Println("📥 Enter stock symbols to subscribe (e.g., AAPL). Type 'exit' to quit:")

	for scanner.Scan() {
		symbol := scanner.Text()
		if symbol == "exit" {
			break
		}

		// Send StockRequest to server
		req := &stockpb.StockRequest{Symbol: symbol}
		if err := stream.Send(req); err != nil {
			log.Fatalf("❌ Failed to send symbol: %v", err)
		}

		// Delay to avoid overwhelming server (optional)
		time.Sleep(500 * time.Millisecond)
	}

	// Close sending side of the stream when done
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("❌ Failed to close stream: %v", err)
	}
}
