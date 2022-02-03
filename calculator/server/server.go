package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/arun6783/go-grpc-demo/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("Calculate SquareRoot was called with req:%v", req)

	number := req.GetNumber()

	if number < 0 {
		return nil,
			status.Errorf(
				codes.InvalidArgument,
				fmt.Sprintf("Received a negative number:%v", number),
			)
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}
func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {

	fmt.Printf("Calculate server was called with req:%v", req)
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

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {

	fmt.Printf("FindMaximum server was called with stream")
	var prev int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {

			log.Fatalf("Error occured when receiving stream%v\n", err)
			return err
		}

		current := req.GetNumber_1()

		if current > prev {
			sendErr := stream.Send(&calculatorpb.FindMaximumReponse{
				Result: current,
			})

			if sendErr != nil {
				log.Fatalf("Error occured when sending stream %v\n", sendErr)
				return sendErr
			}
		}

		prev = current
	}
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
