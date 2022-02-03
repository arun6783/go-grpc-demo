package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/arun6783/go-grpc-demo/greet/greetpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	fmt.Println("Hello in client go")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect:%v\n\n", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//unaryApiReceive(c)

	//receiveServerStreaming(c)

	//sendClientStreaming(c)

	//sendBiDirectionalStream(c)

	dounaryWithDeadline(c, 1*time.Second)
	dounaryWithDeadline(c, 5*time.Second)
}

func dounaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {

	fmt.Printf("Greet with deadline client was called with duration %v\n", timeout)

	req := &greetpb.GreetWithDeadLineRequest{
		Greeting: &greetpb.Greeting{FirstName: "Rekha", LastName: "Bothiraj"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	res, err := c.GreetWithDeadLine(ctx, req)

	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit!! deadline was exceeded!!")
			} else {
				fmt.Printf("unexpected error occured %v\n", statusErr)
			}
		} else {
			log.Fatalf("Error occured when calling GreetWithDeadLine %v\n", err)
		}
	}

	fmt.Printf("Response from greet %v\n", res)
}

func sendBiDirectionalStream(c greetpb.GreetServiceClient) {
	fmt.Println("BiDirectionalStream client   was called")

	//we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream:%v\n\n", err)
	}

	waitc := make(chan struct{})
	//we send a bunch of MessageState (go routine)
	requests := getGreetEveryOneMessages()

	go func() {
		for _, req := range requests {
			log.Printf("sending message - %v\n\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}

		stream.CloseSend()
	}()
	//we receive a bunch of messages from the client time go routine

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				log.Printf("End of stream")
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving %v\n\n", err)
				break
			}

			log.Printf("Response from server : %v\n\n", response.GetResult())
		}
		close(waitc)
	}()
	//block until everything is done
	<-waitc
}

func sendClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("SendClientStreaming  was called")

	requests := getGreetMessages()

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error when sending client stream %v\n\n\n", err)
	}
	for _, req := range requests {
		fmt.Printf("sending request %v\n\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Microsecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error when receiving long stream %v\n\n", err)
	}
	fmt.Printf("Long stream response is %v\n\n", res)
}

func receiveServerStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("ReceiveServerStreaming  was called")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{FirstName: "Rekha", LastName: "Bothiraj"},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error occured when calling greet multiple %v\n\n\n", err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			log.Println("End of stream")
			break
		}
		if err != nil {
			log.Fatalf("ERror occured when receiving stream %v\n\n\n", err)
		}

		log.Printf("response from server is %v\n\n", msg.GetResult())
	}

}
func unaryApiReceive(c greetpb.GreetServiceClient) {

	fmt.Println("Greet  client was called")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{FirstName: "Rekha", LastName: "Bothiraj"},
	}

	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Error occured when calling greet %v\n\n\n", err)
	}

	fmt.Printf("Response from greet %v\n\n\n", res)
}

func getGreetMessages() []*greetpb.LongGreetRequest {
	return []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{FirstName: "John"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Tony"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Lucy"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Viktor"},
		},
	}
}

func getGreetEveryOneMessages() []*greetpb.GreetEveryoneRequest {
	return []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{FirstName: "John"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Tony"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Lucy"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Viktor"},
		},
	}
}
