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
	//readBlog(c, "61fd6ef6318a346e9ae3ce84")

}

func readBlog(c blogpb.BlogServiceClient, blogId string) {

	fmt.Println("Read blog  client was called")

	req := &blogpb.ReadBlogRequest{BlogId: blogId}

	res, err := c.ReadBlog(context.Background(), req)

	if err != nil {
		log.Fatalf("Error occured when reading new blog %v\n", err)
	}
	fmt.Printf("Response from read blog %v\n", res)
}

func postBlog(c blogpb.BlogServiceClient) {

	fmt.Println("Create new blog  client was called")

	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Arun", Title: "Second Post GO-gRPC", Content: "This is the second post go grpc learning!!!!!",
		},
	}

	res, err := c.CreateBlog(context.Background(), req)

	if err != nil {
		log.Fatalf("Error occured when creating new blog %v\n", err)
	}

	fmt.Printf("Response from Create blog!!! %v\n now going to read the blog\n\n", res)

	readBlog(c, res.GetBlog().Id)
}
