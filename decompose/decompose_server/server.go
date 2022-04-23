package main

import (
	"log"
	"net"

	"github.com/vitorgraciani/grpc-go-aula-2/decompose/decomposepb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) PrimeNumberDecompose(in *decomposepb.DecomposeRequest, stream decomposepb.DecomposeService_PrimeNumberDecomposeServer) error {
	number := in.GetFirstNumer()
	var prime int32 = 2
	for number > 1 {
		if number%prime == 0 {
			stream.Send(&decomposepb.DecomposeResponse{
				Result: prime,
			})
			number = number / prime
		} else {
			prime += 1
		}
	}

	return nil
}

func main() {

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Fail to list %v", err)
	}

	s := grpc.NewServer()
	decomposepb.RegisterDecomposeServiceServer(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
