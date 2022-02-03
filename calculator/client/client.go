package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/arun6783/go-grpc-demo/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("hello from calculator client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error occured when dialing %v\n", err)
	}

	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)

	//doCalculation(c)

	doBiDirectionalStream(c)
}

func doBiDirectionalStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Print("do BiDirectionalStream client method was called")
	//we create a stream by invoking the client

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error occured when creating a stream for FindMaximum%v\n", err)
	}

	requests := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest{
			Number_1: 1,
		},
		&calculatorpb.FindMaximumRequest{
			Number_1: 5,
		},
		&calculatorpb.FindMaximumRequest{
			Number_1: 3,
		},
		&calculatorpb.FindMaximumRequest{
			Number_1: 6,
		},
		&calculatorpb.FindMaximumRequest{
			Number_1: 2,
		},
		&calculatorpb.FindMaximumRequest{
			Number_1: 20,
		},
	}
	waitc := make(chan struct{})

	//send messages to server
	go func() {
		for _, req := range requests {

			log.Printf("sending message to server %v", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	//receive messages from server
	go func() {

		for {

			response, err := stream.Recv()

			if err == io.EOF {
				log.Print("end of stream")
				break
			}
			if err != nil {
				log.Fatalf("Error occured when receiving the stream %v\n", err)
				break
			}

			log.Printf("Response from server : %v\n", response.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}

func doCalculation(c calculatorpb.CalculatorServiceClient) {
	fmt.Print("do calculation client method was called")
	req := &calculatorpb.CalculatorRequest{
		Number_1:  210,
		Number_2:  20,
		Operation: calculatorpb.Operation_OPERATOR_ADD,
	}
	res, err := c.Calculate(context.Background(), req)

	if err != nil {
		log.Fatalf("error occured when calculating %v\n", err)
	}

	fmt.Printf("response from calculator was %v\n", res)
}
