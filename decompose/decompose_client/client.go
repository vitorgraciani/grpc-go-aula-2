package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/vitorgraciani/grpc-go-aula-2/decompose/decomposepb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting client...")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error to connect to server: %v", err)
	}

	defer conn.Close()
	cc := decomposepb.NewDecomposeServiceClient(conn)
	decomposeStream(cc)
}

func decomposeStream(cc decomposepb.DecomposeServiceClient) {
	req := &decomposepb.DecomposeRequest{
		FirstNumer: 120,
	}

	stream, err := cc.PrimeNumberDecompose(context.Background(), req)
	if err != nil {
		log.Fatalf("error to decompose number %v", err)
	}
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		fmt.Printf("%d ", msg.GetResult())
	}
}
