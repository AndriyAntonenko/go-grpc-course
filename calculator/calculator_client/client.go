package main

import (
	"context"
	"fmt"
	"io"
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
	// calculateSum(c)
	numberDecomposition(c)
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

func numberDecomposition(c calculatorpb.SumServiceClient) {
	req := &calculatorpb.NumberDecompositionRequest{
		Value: 44,
	}

	stream, err := c.NumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling NumberDecomposition RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		fmt.Printf("Prime number: %v\n", res.PrimeNumber)
	}
}
