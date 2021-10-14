package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AndriyAntonenko/go-grpc-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	fmt.Println("Calculator client is running")
	c := calculatorpb.NewSumServiceClient(cc)
	calculateSum(c)
}

func calculateSum(c calculatorpb.SumServiceClient) {
	req := &calculatorpb.SumRequest{
		Input: &calculatorpb.Sum{
			A: 10,
			B: 3,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	fmt.Printf("Response from Sum: %v", res.Result)
}
