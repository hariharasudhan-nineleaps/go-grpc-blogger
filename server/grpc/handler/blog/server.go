package blog

import (
	"context"

	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/proto/blog"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/proto/user"
)

type BlogServer struct {
	blog.UnimplementedBlogServiceServer
}

func (bs *BlogServer) CreateBlog(ctx context.Context, cbRequest *blog.CreateBlogRequest) (*blog.CreateBlogResponse, error) {
	return &blog.CreateBlogResponse{
		Id:          "1",
		Title:       "Test title",
		Description: "Test Description",
		Category:    0,
		Tags:        []string{"tag1", "tag2"},
		Author:      &user.User{Id: "1", Name: "Test Name", Email: "Test Email"},
	}, nil
}
