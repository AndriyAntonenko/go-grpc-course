package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/AndriyAntonenko/go-grpc-course/blog/blogpb"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	MongoUser string
	MongoPass string
	MongoUrl  string
	MongoPort string
	MongoDB   string
}

func initAppConfig() *Config {
	var err error
	var configInstance *Config

	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	configInstance = &Config{
		MongoUser: os.Getenv("MONGO_INITDB_DATABASE_USER"),
		MongoPass: os.Getenv("MONGO_INITDB_PASSWORD"),
		MongoUrl:  os.Getenv("MONGO_URL"),
		MongoPort: os.Getenv("MONGO_PORT"),
		MongoDB:   os.Getenv("MONGO_INITDB_DATABASE"),
	}

	return configInstance
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func connectMongo(config *Config) *mongo.Client {
	connectionString := "mongodb://" + config.MongoUser + ":" + config.MongoPass + "@" + config.MongoUrl + ":" + config.MongoPort + "/" + config.MongoDB
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatalf("error connect to mongo db %v\n", connectionString)
	}
	return mongoClient
}

var collection *mongo.Collection

type server struct{}

func (s server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("CreateBlog rpc is running")
	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}
	fmt.Println(data)

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}

	fmt.Println("Blog created successfully")
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func main() {
	config := initAppConfig()
	mongoClient := connectMongo(config)
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	collection = mongoClient.Database(config.MongoDB).Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to servie %v", err)
		}
	}()

	fmt.Println("Blog server is running!!!")

	// Wait for Ctrl+C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("closing the listener")
	lis.Close()
	fmt.Println("Closing mongodb connection")
	mongoClient.Disconnect(context.TODO())
	fmt.Println("End of program")
}
