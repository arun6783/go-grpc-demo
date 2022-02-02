package main

import (
	"context"
	"fmt"
	"log"

	"github.com/arun6783/go-grpc-demo/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("hello from calculator client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error occured when dialing %v", err)
	}

	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	req := &calculatorpb.CalculatorRequest{
		Number_1:  210,
		Number_2:  20,
		Operation: calculatorpb.Operation_OPERATOR_ADD,
	}
	res, err := c.Calculate(context.Background(), req)

	if err != nil {
		log.Fatalf("error occured when calculating %v", err)
	}

	fmt.Print(fmt.Sprintf("response from calculator was %v", res))
}
