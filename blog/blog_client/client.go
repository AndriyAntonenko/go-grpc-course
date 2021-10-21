package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AndriyAntonenko/go-grpc-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog client is running...")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	fmt.Println("Greet client is running")
	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Stephan",
		Title:    "My First Blog",
		Content:  "Conten of this blog",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Error during creating blog: %v\n", err)
		return
	}
	fmt.Printf("Blog has been created: %v\n", res)
}
