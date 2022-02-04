package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/arun6783/go-grpc-demo/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	Collection *mongo.Collection
)

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := req.GetBlog()

	fmt.Printf("Update blog method was called with %v\n", blog)

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "InvalidArgument -  Error converting blogid")
	}

	filter := bson.D{{"_id", oid}}
	update := bson.D{{"$set", bson.D{{"author_id", data.AuthorID}, {"Title", "data.Title"}, {"content", data.Content}}}}

	res, err := Collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "InternalError occured when trying to update data in to db%v\n", err)
	}

	if res.MatchedCount == 0 {
		fmt.Println("No records found to update")
		return nil, status.Errorf(codes.NotFound, "NotFound cannot match item to update with given blog id %v\n", blog.GetId())
	}
	fmt.Printf("Documents matched: %v\n", res.MatchedCount)
	fmt.Printf("Documents updated: %v\n", res.ModifiedCount)

	return &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{Id: blog.GetId(), Title: blog.GetTitle(), Content: blog.GetContent(), AuthorId: blog.GetAuthorId()},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {

	blogId := req.GetBlogId()

	fmt.Printf("Read blog method was called with blogid %v\n", blogId)

	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "InvalidArgument -  Error converting blogid")
	}

	data := &blogItem{}

	filter := bson.D{{"_id", oid}}

	res := Collection.FindOne(context.Background(), filter)

	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot fund blog with specified id %v", oid)
	}
	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{Id: blogId, AuthorId: data.AuthorID, Title: data.Title, Content: data.Content},
	}, nil
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()

	fmt.Printf("Create blog method was called with %v\n", blog)

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	res, err := Collection.InsertOne(context.Background(), data)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "InternalError occured when trying to insert data in to db%v\n", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, status.Errorf(codes.Internal, "InternalError cannot convert to oid")
	}
	fmt.Printf("New blog post created")
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{Id: oid.Hex(), Title: blog.GetTitle(), Content: blog.GetContent(), AuthorId: blog.GetAuthorId()},
	}, nil
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Print("Hello world!! blog server started in .go\n")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://admin:pass@localhost:27017"))
	if err != nil {
		panic(err)
	}

	Collection = client.Database("mydb").Collection("blog")

	log.Println("connection to mongo successful")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting blog go server......")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	//wait for ctrl+c
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	fmt.Println("stopping server..")

	s.Stop()

	fmt.Println("Closing the listener")

	lis.Close()

	fmt.Println("closing mongo connecction")
	client.Disconnect(context.TODO())

	fmt.Println("End of program")

}
