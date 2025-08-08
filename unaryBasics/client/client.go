package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"unaryBasics/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("🔗 Connecting to server...")

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server")
	}
	defer conn.Close()

	client := greetpb.NewGreetServiceClient(conn)

	doUnary(client)
}

func doUnary(c greetpb.GreetServiceClient) {

	req := &greetpb.GreetRequest{
		FirstName: "Mushahid",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.Greet(ctx, req)

	if err != nil {
		log.Fatalf("❌ Error calling Greet RPC: %v", err)
	}

	fmt.Printf("✅ Response: %v\n", res.GetResult())

}
