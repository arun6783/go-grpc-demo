package main

import (
	"fmt"
	"log"

	"github.com/arun6783/go-grpc-demo/blog/blogpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello in Blog Client go")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error occured when dialing %v\n", err)
	}

	defer conn.Close()
	c := blogpb.NewBlogServiceClient(conn)

	postBlog(c)

}
func postBlog(c blogpb.BlogServiceClient) {

	fmt.Println("Create new blog  client was called")

	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Arun", Title: "Welcome to GO-gRPC", Content: "Welcome to go grpc learning!!!!!",
		},
	}

	res, err := c.CreateBlog(context.Background(), req)

	if err != nil {
		log.Fatalf("Error occured when creating new blog %v\n\n\n", err)
	}

	fmt.Printf("Response from Create blog %v\n\n\n", res)
}
