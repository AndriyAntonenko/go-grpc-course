package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/AndriyAntonenko/go-grpc-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	a := req.GetInput().GetA()
	b := req.GetInput().GetB()
	res := calculatorpb.SumResponse{
		Result: a + b,
	}
	return &res, nil
}

func (s *server) NumberDecomposition(req *calculatorpb.NumberDecompositionRequest, stream calculatorpb.SumService_NumberDecompositionServer) error {
	value := req.GetValue()
	var k int32 = 2
	for {
		if value <= 1 {
			break
		}

		if value%k == 0 {
			stream.Send(&calculatorpb.NumberDecompositionResponse{
				PrimeNumber: k,
			})
			value = value / k
			continue
		}

		k += 1
	}

	return nil
}

func (s *server) ComputeAverage(stream calculatorpb.SumService_ComputeAverageServer) error {
	var sum int32 = 0
	var count int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: float32(sum) / float32(count),
			})
		}
		if err != nil {
			log.Fatalf("error while receiving request %v", err)
		}

		sum += req.GetValue()
		count += 1
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterSumServiceServer(s, &server{})
	fmt.Println("Calculator server is running!!!")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to servie %v", err)
	}
}
