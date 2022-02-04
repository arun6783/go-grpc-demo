package main

import (
	"fmt"
	"log"

	"github.com/arun6783/go-grpc-demo/blog/blogpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	fmt.Println("Hello in Blog Client go")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error occured when dialing %v\n", err)
	}

	defer conn.Close()
	c := blogpb.NewBlogServiceClient(conn)

	//postBlog(c)
	//readBlog(c, "61fd6ef6318a346e9ae3ce84")

	//updateBlog(c)

	deleteBlog(c, "61fd6ef6318a346e9ae3ce84")
}

func deleteBlog(c blogpb.BlogServiceClient, blogId string) {

	fmt.Println("Delete blog  client was called")

	req := &blogpb.DeleteBlogRequest{BlogId: blogId}

	res, err := c.DeleteBlog(context.Background(), req)

	if err != nil {
		log.Fatalf("Error occured when deleting blog %v\n", err)
	}

	fmt.Printf("Successfully deleted blog %v\n", res)
}

func updateBlog(c blogpb.BlogServiceClient) {

	fmt.Println("Update blog  client was called")

	req := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{Id: "61fd6ef6318a346e9ae3ce84",
			AuthorId: "Arun", Title: "first blog post updated", Content: "This is the updated content!!!!!",
		},
	}

	res, err := c.UpdateBlog(context.Background(), req)
	if err != nil {
		fmt.Printf("Update blog error %v\n\n", err)
		st, ok := status.FromError(err)
		if ok {
			// Error was  a status error
			if st.Code() == codes.NotFound {
				fmt.Printf("cannot find blog . response from server - %v\n", st.Message())

			}
		} else {
			fmt.Printf("error converting update blog error to status error %v\n", err)
		}
	}
	if res != nil {
		fmt.Printf("Response from update blog!!! %v\n", res)
	}

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
