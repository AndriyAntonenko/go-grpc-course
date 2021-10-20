package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/AndriyAntonenko/go-grpc-course/blog/blogpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
	fmt.Println("End of program")
}
