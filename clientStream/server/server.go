package main

import (
	"clientStream/student"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

// struct to implement student service server interface
type server struct {
	student.UnimplementedStudentServiceServer
}

// Submit Assignments handles client streaming rpc.
// It receives multiple assignments from the client and calculates the average.
// And returns the final grade, once all assignments are received.
func (s *server) SubmitAssignments(stream student.StudentService_SubmitAssignmentsServer) error {
	var total int32 // to calculate the sum of marks
	var count int32 // to count how many assignments are received

	fmt.Println(" 📩 Receiving assignments from the client......")

	// Receive ssignments in a loop until the client finishes sending them.
	for {
		req, err := stream.Recv() // receive one message(assignmment) one at a time.

		// Client finished seending all assignemnts
		if err == io.EOF {
			// Calculate average marks
			avg := float32(total) / float32(count)

			// Decide remarks based on average
			var remark string
			if avg >= 90 {
				remark = "Outstanding"
			} else if avg >= 60 {
				remark = "Average"
			} else {
				remark = "Needs Improvement"
			}

			// Send back the final grade and close the stream
			return stream.SendAndClose(&student.FinalGrade{
				Average: avg,
				Remarks: remark,
			})
		}

		// If there's an error streaming not EOF then log the error
		if err != nil {
			log.Fatalf("❌ Error receiving the assignment %v", err)
		}

		// If there's no error then...
		fmt.Printf("✅ Received %s, (%d marks)", req.GetTitle(), req.GetMarks())

		// Accumulate total and count for average calculation
		total += req.GetMarks()
		count++
	}
}

func main() {
	fmt.Println("🚀Starting gRPC Server......")

	// Listen for tcp conncetions on port 50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("❌Failed to listen....")
	}

	// Create a new grpc server
	s := grpc.NewServer()

	// Register our server
	student.RegisterStudentServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while listening to server %v", err)
	}
}
