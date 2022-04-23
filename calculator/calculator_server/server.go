package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/vitorgraciani/grpc-go-aula-2/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
}

func (*server) Sum(ctx context.Context, req *calculatorpb.CalcRequest) (*calculatorpb.CalcResponse, error) {
	fmt.Printf("Server was invoked, req: %v \n", req)
	first := req.GetFirst()
	second := req.GetSecond()

	result := first + second

	resp := &calculatorpb.CalcResponse{
		Result: result,
	}
	return resp, nil
}

func (s *server) Max(stream calculatorpb.CalculatorService_MaxServer) error {
	log.Println("Receving request")
	var previous int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error to send on stream %v \n", err)
		}

		if previous <= req.Number {
			log.Printf("Previous: %d recv: %d", previous, req.Number)
			stream.SendMsg(&calculatorpb.MaxResponse{
				Result: req.Number,
			})

		}
		previous = req.Number
	}
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	log.Println("Receving request")
	var sum int32 = 0
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {

			stream.SendMsg(&calculatorpb.CalculatorResponse{
				Result: float64(sum) / float64(count),
			})
			break
		}
		if err != nil {
			log.Fatalf("Error to receive request on stream %v \n", err)
		}
		sum += req.GetElem()
		count += 1
	}
	return nil
}

func main() {

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Fail to list %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
