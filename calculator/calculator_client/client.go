package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vitorgraciani/grpc-go-aula-2/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error to connect to server %v", err)
	}

	client := calculatorpb.NewCalculatorServiceClient(cc)
	//unary(client)
	//computeAverage(client)
	max(client)
}

func unary(client calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.CalcRequest{
		First:  10,
		Second: 3,
	}
	resp, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error to sum values %v", err)
	}
	log.Printf("server response: %d", resp.Result)
}

func computeAverage(client calculatorpb.CalculatorServiceClient) {
	reqs := []*calculatorpb.CalculatorRequest{
		{Elem: 1},
		{Elem: 2},
		{Elem: 3},
		{Elem: 4},
	}

	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error to send request on stream %v \n ", err)
	}
	for _, req := range reqs {
		log.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error whie receiving response from stream %v\n", err)
	}
	log.Printf("Compute average result is %f", res.GetResult())
}

func max(client calculatorpb.CalculatorServiceClient) {
	reqs := []*calculatorpb.MaxRequest{
		{Number: 1},
		{Number: 5},
		{Number: 3},
		{Number: 6},
		{Number: 2},
		{Number: 20},
	}
	stream, err := client.Max(context.Background())
	if err != nil {
		log.Fatalf("Error to create stream")
	}

	ch := make(chan struct{})

	go func() {
		for _, req := range reqs {
			log.Printf("Sending %v", req)
			stream.Send(req)
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
				log.Fatalf("Error while reading stream %v \n", err)
			}
			log.Printf("Result is: %d \n", res.Result)
		}
		close(ch)
	}()

	<-ch
}
