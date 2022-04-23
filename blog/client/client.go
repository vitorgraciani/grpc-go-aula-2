package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/vitorgraciani/grpc-go-aula-2/blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var addr string = "0.0.0.0:50051"

func main() {
	fmt.Println("Starting client...")
	cc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error to connect to server %v", err)
	}

	client := pb.NewBlogServiceClient(cc)

	id := createBlog(client)
	res := readBlog(client, id)
	log.Printf("res %v", res)
	//readBlog(client, "non-existing")
	//updateBlog(client, id)
	//listBlog(client)
	deleteBlog(client, id)
	readBlog(client, id)
}

func createBlog(c pb.BlogServiceClient) string {
	log.Println("createBlog was invoked")

	req := &pb.Blog{
		AuthorId: "Xerxes",
		Title:    "Xerxes the xerxes",
		Content:  "My first paper",
	}

	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Unexpected error %v\n", err)
	}

	log.Printf("Bloq was created %s\n", res.Id)
	return res.Id
}

func readBlog(c pb.BlogServiceClient, id string) *pb.Blog {
	log.Println("readBlog was invoked")

	res, err := c.ReadBlog(context.Background(), &pb.BlogID{Id: id})
	if err != nil {
		log.Printf("error while reading %v\n", err)
	}

	return res
}

func updateBlog(c pb.BlogServiceClient, id string) {
	log.Println("updateBlog was invoked")

	req := &pb.Blog{
		Id:       id,
		AuthorId: "Not you my friend",
		Title:    "New title",
		Content:  "New content",
	}

	if _, err := c.UpdateBlog(context.Background(), req); err != nil {
		log.Fatalf("error while updating blog %v\n", err)
	}
}

func listBlog(c pb.BlogServiceClient) {
	log.Println("listBlog was invoked")

	stream, err := c.ListBlogs(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("error while calling list blogs %v\n", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("omething happened %v\n", err)
		}
		log.Println(res)
	}
}

func deleteBlog(c pb.BlogServiceClient, id string) {
	log.Println("deleteBlog was invoked")

	_, err := c.DeleteBlog(context.Background(), &pb.BlogID{Id: id})
	if err != nil {
		log.Fatalf("error while deleting %v\n", err)
	}
	log.Println("Blog was deleted")
}
