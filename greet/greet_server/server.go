package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/vitorgraciani/grpc-go-aula-2/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Server was invoked, req: %v \n", req)
	name := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	resp := &greetpb.GreetResponse{
		Result: "Hello my name is " + name + " " + lastName,
	}
	return resp, nil
}

func main() {
	fmt.Println("Hello world")

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
