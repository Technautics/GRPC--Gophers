package main

import (
	"bidirectionalStream/stockpb"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	stockpb.UnimplementedStockServiceServer
}

// StreamPrices handles the bidirectional streaming logic
func (s *server) StreamPrices(stream stockpb.StockService_StreamPricesServer) error {
	log.Println("📡 Client connected for stock updates")

	// Store subscribed stock symbols from client
	symbols := make(map[string]bool)
	mu := sync.Mutex{} // Mutex to safely update symbols map across goroutines

	// Signal channel for when client stops sending data
	done := make(chan struct{})

	// ================================
	// Goroutine 1: Receiving from Client
	// ================================
	go func() {
		for {
			// Receive a StockRequest from client
			req, err := stream.Recv()
			if err == io.EOF {
				// Client closed the sending stream
				log.Println("❌ Client stopped sending symbols")
				close(done)
				return
			}
			if err != nil {
				// Some error occurred while receiving data
				log.Printf("❌ Error receiving symbol: %v", err)
				close(done)
				return
			}

			// Add new symbol to subscription list
			mu.Lock()
			symbols[req.Symbol] = true
			mu.Unlock()
			log.Printf("✅ Subscribed to: %s", req.Symbol)
		}
	}()

	// Ticker to send updates every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// ================================
	// Goroutine 2: Sending to Client
	// ================================
	for {
		select {
		case <-done:
			// Stop sending when client is done
			return nil
		case <-ticker.C:
			// Iterate through subscribed symbols and send updates
			mu.Lock()
			for symbol := range symbols {
				price := rand.Float64()*1000 + 100 // Generate mock price
				res := &stockpb.StockPrice{
					Symbol: symbol,
					Price:  price,
					Time:   time.Now().Format(time.RFC3339),
				}

				// Send stock update to client
				if err := stream.Send(res); err != nil {
					log.Printf("❌ Error sending update: %v", err)
					mu.Unlock()
					return err
				}
				log.Printf("📤 Sent %s update: $%.2f", symbol, price)
			}
			mu.Unlock()
		}
	}
}

func main() {
	// Start listening on TCP port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}

	// Create new gRPC server
	grpcServer := grpc.NewServer()

	// Register our StockService with gRPC server
	stockpb.RegisterStockServiceServer(grpcServer, &server{})

	log.Println("🚀 Stock Price Server listening on port 50051")

	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}
