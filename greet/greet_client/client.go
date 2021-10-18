package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/AndriyAntonenko/go-grpc-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	fmt.Println("Greet client is running")
	c := greetpb.NewGreetServiceClient(cc)

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	doUnaryWithDeadline(c, 5)
	doUnaryWithDeadline(c, 2)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting do unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "andrii",
			LastName:  "antonenko",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	fmt.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting do server streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "andrii",
			LastName:  "antonenko",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimesRPC: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// we'he reached the end of stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		fmt.Printf("Response from GreatManyTimes %v\n", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting do client streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Andrii",
				LastName:  "Antonenko",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
				LastName:  "Black",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
				LastName:  "Twik",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet RPC: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("error while sending data to LongGreet RPC: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet RPC: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", response)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting do BiDi streaming RPC...")

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Andrii",
				LastName:  "Antonenko",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
				LastName:  "Black",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
				LastName:  "Twik",
			},
		},
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling GreetEveryone RPC: %v", err)
	}

	waitChanel := make(chan struct{})

	go func() {
		for _, req := range requests {
			fmt.Printf("sending request %v\n", req)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error while sending request: %v\n", err)
				close(waitChanel)
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving: %v\n", err)
			}
			fmt.Printf("received: %v\n", res.GetResult())
		}
		close(waitChanel)
	}()

	<-waitChanel
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, seconds int) {
	fmt.Println("Starting do unary with deadline RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "andrii",
			LastName:  "antonenko",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		resErr, ok := status.FromError(err)
		if !ok {
			log.Fatalf("error while calling GreetWithDeadline RPC: %v\n", err)
		}

		if resErr.Code() == codes.DeadlineExceeded {
			fmt.Println("Timeout was hit! Deadline was exceeded")
		} else {
			fmt.Printf("unexpected error: %v\n", err)
		}
		return
	}
	fmt.Printf("Response from GreetWithDeadline: %v\n", res.GetResult())
}
