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
)

type server struct{}

// TODO: split code into files
type Config struct {
	MongoUser string
	MongoPass string
	MongoUrl  string
	MongoPort string
	MongoDB   string
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func InitAppConfig() *Config {
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

func connectMongo(config *Config) *mongo.Client {
	connectionString := "mongodb://" + config.MongoUrl + ":" + config.MongoPass + "@" + config.MongoUrl + ":" + config.MongoPort
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatalf("error connect to mongo db %v\n", connectionString)
	}
	return mongoClient
}

func getBlogCollection(mongoClient *mongo.Client, config *Config) *mongo.Collection {
	if mongoClient == nil {
		log.Fatalf("connection is nil")
	}
	return mongoClient.Database(config.MongoDB).Collection("blog")
}

func main() {
	config := InitAppConfig()
	mongoClient := connectMongo(config)
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	getBlogCollection(mongoClient, config)

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
