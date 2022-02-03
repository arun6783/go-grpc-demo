package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/arun6783/go-grpc-demo/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	result := "Hello" + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greetmultiple function was invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "Hello" + firstName + strconv.Itoa(i)

		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)

		time.Sleep(1000 * time.Millisecond)
	}

	return nil

}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("Longgreet function was invoked with a streaming request\n")
	result := ""
	for {

		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("error when reading client stream %v\n", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + " !"

	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("Greeteveryone function was invoked with a streaming request\n")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client steram %v\n", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()

		result := "Hello " + firstName + " !"

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {

			log.Fatalf("Error occured when sending stream to client %v\n", sendErr)
			return sendErr
		}

	}

}

func main() {

	fmt.Print("Hello world!! inside greet server .go\n")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
