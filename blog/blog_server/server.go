package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/arun6783/go-grpc-demo/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
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

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Print("Hello world!! blog server started in .go\n")

	// client, merr := mongo.NewClient("mongodb://admin:pass@localhost:27017")
	// if merr != nil {
	// 	log.Fatalf("Error occured when connecting to mongo server%v\n", merr)
	// }

	// clientErr := client.Connect(context.TODO())
	// if clientErr != nil {
	// 	log.Fatalf("Mongo Error occured when connecting to todo context%v\n", clientErr)

	// }

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
