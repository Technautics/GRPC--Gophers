package main

import (
	"clientStream/student"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("📡 Connecting to server......")

	// Connect to grpc server at localhost:50052
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌Could not connect.. %v", err)
	}
	defer conn.Close()

	// Create a client for student service
	client := student.NewStudentServiceClient(conn)

	// Start the client stream to send messages.
	stream, err := client.SubmitAssignments(context.Background())
	if err != nil {
		log.Fatalf("Error while starting the stream....")
	}

	// Prepare a list of messages(assignments) to send....
	assignments := []student.Assignment{
		{Title: "Maths", Marks: 85},
		{Title: "Chemistry", Marks: 89},
		{Title: "Physics", Marks: 90},
	}

	// Loop over messages(assignments) and send each to server.
	for i := range assignments {
		fmt.Printf("📤 Sending -> %s : %d \n", assignments[i].Title, assignments[i].Marks)
		if err := stream.Send(&assignments[i]); err != nil {
			log.Fatalf("Error sending the messages. %v", err)
		}
		time.Sleep(time.Second) // simmulate 1s delay between sending messages(optional)
	}

	// All assignments are send. Now close the stream and send the final grade..
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("❌Error receiving the final grade. %v", err)
	}

	// Print the final grade and remark from the server.
	fmt.Printf("🎓Final Grade: %.2f - %s", reply.GetAverage(), reply.GetRemarks())
}
