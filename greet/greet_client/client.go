package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vitorgraciani/grpc-go-aula-2/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Starting client...")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error to connect to server: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	//fmt.Printf("created client %f", c)

	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Piruleibe",
			LastName:  "Maluco",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error to call greet service %v", err)
	}
	fmt.Printf("Greet service response: %s \n", res.Result)
}
