package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/vitorgraciani/grpc-go-aula-2/blog/model"
	pb "github.com/vitorgraciani/grpc-go-aula-2/blog/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var collection *mongo.Collection
var addr string = "0.0.0.0:50051"

type server struct {
	pb.BlogServiceServer
}

func (s *server) CreateBlog(ctx context.Context, in *pb.Blog) (*pb.BlogID, error) {
	log.Printf("CreateBlog was invoked %v", in)

	data := model.BlogItem{
		AuthorID: in.AuthorId,
		Title:    in.Title,
		Content:  in.Content,
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v\n", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			"Cannot convertto OID",
		)
	}

	return &pb.BlogID{
		Id: oid.Hex(),
	}, nil
}

func (s *server) ReadBlog(ctx context.Context, in *pb.BlogID) (*pb.Blog, error) {
	log.Printf("ReadBlog was invoked %v", in)

	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}

	data := &model.BlogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(ctx, filter)

	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"Cannot find blog with id provided",
		)
	}

	return model.DocumentToBlog(data), nil
}

func (s *server) UpdateBlog(ctx context.Context, in *pb.Blog) (*emptypb.Empty, error) {
	log.Printf("UpdateBlog was invoked %v", in)

	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}
	data := &model.BlogItem{
		AuthorID: in.AuthorId,
		Title:    in.Title,
		Content:  in.Content,
	}
	res, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": data},
	)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Could not update",
		)
	}

	if res.MatchedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			"Cannot findblog with id",
		)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) DeleteBlog(ctx context.Context, in *pb.BlogID) (*emptypb.Empty, error) {
	log.Printf("DeleteBlog was invoked %v\n", in)

	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}

	res, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete object in mongoDB %v", err),
		)
	}
	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			"Blog was not found",
		)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) ListBlogs(in *emptypb.Empty, stream pb.BlogService_ListBlogsServer) error {
	log.Println("ListBlogs was invoked")

	cur, err := collection.Find(context.Background(), primitive.D{{}})

	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error %v", err),
		)
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := &model.BlogItem{}
		if err := cur.Decode(data); err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MONGODB %v", err),
			)
		}
		stream.Send(model.DocumentToBlog(data))

		if err := cur.Err(); err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Unknown internal error %v", err),
			)
		}
	}
	return nil
}

func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		log.Fatalf("Fail to connect mongo %v", err)
	}

	if err := client.Connect(context.Background()); err != nil {
		log.Fatal(err)
	}

	collection = client.Database("blogdb").Collection("blog")

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Fail to list %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBlogServiceServer(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
