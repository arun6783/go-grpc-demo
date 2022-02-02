package main

import (
	"fmt"
	"io"
	"log"

	"github.com/arun6783/go-grpc-demo/greet/greetpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello in client go")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect:%v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//doGreet(c)

	doGreetMultiple(c)
}

func doGreetMultiple(c greetpb.GreetServiceClient) {

	fmt.Println("Greet multiple client was called")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{FirstName: "Rekha", LastName: "Bothiraj"},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error occured when calling greet multiple %v\n", err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			log.Println("End of stream")
			break
		}
		if err != nil {
			log.Fatalf("ERror occured when receiving stream %v\n", err)
		}

		log.Printf("response from server is %v", msg.GetResult())
	}

}
func doGreet(c greetpb.GreetServiceClient) {

	fmt.Println("Greet  client was called")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{FirstName: "Rekha", LastName: "Bothiraj"},
	}

	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Error occured when calling greet %v\n", err)
	}

	fmt.Printf("Response from greet %v\n", res)
}
