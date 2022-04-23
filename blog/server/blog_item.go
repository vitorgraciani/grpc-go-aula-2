package main

import (
	"github.com/vitorgraciani/grpc-go-aula-2/blog/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

func documentToBlog(data *BlogItem) *proto.Blog {
	return &proto.Blog{
		ID:       data.ID.Hex(),
		AuthorID: data.AuthorID,
		Title:    data.Title,
		Content:  data.Content,
	}
}
