hello:
	echo "hello world"

generate_protos:
	echo "about to generate protos files"
	protoc greet/greetpb/greet.proto  --go-grpc_out=. --go_out=.
	protoc calculator/calculatorpb/calculator.proto  --go-grpc_out=. --go_out=.

	