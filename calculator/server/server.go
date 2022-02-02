package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/arun6783/go-grpc-demo/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	operator := req.GetOperation()
	if operator == calculatorpb.Operation_OPERATOR_UNKNOWN {
		return nil, errors.New("invalid operator")

	}
	var res float64
	if operator == calculatorpb.Operation_OPERATOR_ADD {
		res = float64(req.GetNumber_1() + req.GetNumber_2())
	} else {
		res = float64(req.GetNumber_1() - req.GetNumber_2())
	}
	return &calculatorpb.CalculatorResponse{
		Result: float64(res),
	}, nil
}

func main() {

	fmt.Println("hello from calculator server")

	listner, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Error occured when creating new server %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(listner); err != nil {
		log.Fatalf("Error occured when starting to serve %v", err)
	}
}
