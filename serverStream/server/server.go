package main

import (
	"fmt"
	"log"
	"net"
	"serverStream/news"
	"time"

	"google.golang.org/grpc"
)

// Server struct implementing generated server interface
type server struct {
	news.UnimplementedNewsServiceServer
}

// Implementing the server stream RPC
func (s *server) GetNewsStream(req *news.NewsRequest, stream news.NewsService_GetNewsStreamServer) error {
	fmt.Printf("📥 Received request for category: %s", req.Category)

	// Simulated headlines for demo purposes
	newsList := []string{
		"🚀 SpaceX launches another rocket",
		"📉 Market hits all-time low",
		"🏆 Local team wins championship",
		"📱 New smartphone released today",
		"🌧️ Heavy rains expected this week",
	}

	// Simulate streaming headlines one by one.
	for _, headline := range newsList {
		res := &news.NewsResponse{
			Headline: headline,
		}

		// Send response to client stream
		if err := stream.Send(res); err != nil {
			log.Fatalf("❌ Failed to send: %v", err)
			return err
		}

		log.Printf("✅ Sent: %s", headline)
		time.Sleep(time.Second) // Simulate delay
	}

	log.Println("✅ All news sent.")
	return nil
}

func main() {
	fmt.Println("🚀Starting gRPC Server......")

	// Listening for tcp connections on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌Failed to listen....")
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register our server
	news.RegisterNewsServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while listening to server %v", err)
	}

}
