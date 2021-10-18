package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/AndriyAntonenko/go-grpc-course/calculator/calculatorpb"
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

	fmt.Println("Calculator client is running")
	c := calculatorpb.NewSumServiceClient(cc)
	// calculateSum(c)
	// numberDecomposition(c)
	// computeAverage(c)
	// findMaximum(c)
	squareRoot(c)
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

func computeAverage(c calculatorpb.SumServiceClient) {
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage RPC: %v", err)
	}

	requests := []*calculatorpb.ComputeAverageRequest{
		{Value: 1},
		{Value: 2},
		{Value: 3},
		{Value: 4},
	}

	for _, req := range requests {
		fmt.Printf("Sending req %v\n", req)
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("error while sending data to ComputeAverage RPC: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	result, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving data from ComputeAverage RPC: %v", err)
	}
	fmt.Printf("Average is: %v", result.GetResult())
}

func findMaximum(c calculatorpb.SumServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while calling FindMaximum RPC: %v", err)
	}

	requests := []*calculatorpb.FindMaximumRequest{
		{Value: 10},
		{Value: 5},
		{Value: 15},
		{Value: 20},
	}

	waitChanel := make(chan struct{})

	go func() {
		for _, req := range requests {
			fmt.Printf("Sending request: %v\n", req)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error while sending data: %v\n", err)
			}
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
				fmt.Printf("error while receiving data: %v\n", err)
				break
			}
			fmt.Printf("Maximum is: %v\n", res.GetMax())
		}
		close(waitChanel)
	}()

	<-waitChanel
}

func squareRoot(c calculatorpb.SumServiceClient) {
	// correct call
	doErrorCall(c, 10)

	// error call
	doErrorCall(c, -10)
}

func doErrorCall(c calculatorpb.SumServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: number,
	})

	if err != nil {
		resErr, ok := status.FromError(err)
		if !ok {
			log.Fatalf("big error calling SquareRoot RPC: %v\n", err)
		}

		// grpc error
		fmt.Println(resErr.Message())
		fmt.Println(resErr.Code())
		if resErr.Code() == codes.InvalidArgument {
			fmt.Println("we probably set negative number!")
			return
		}
	}
	fmt.Printf("result of square root of %v is %v\n", number, res.GetNumberRoot())
}
