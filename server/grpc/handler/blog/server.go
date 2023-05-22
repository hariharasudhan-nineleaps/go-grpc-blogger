package blog

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/proto/blog"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/proto/user"
	"gorm.io/gorm"
)

type BlogServer struct {
	blog.UnimplementedBlogServiceServer
	DB *gorm.DB
}

func (bs *BlogServer) CreateBlog(ctx context.Context, cbRequest *blog.CreateBlogRequest) (*blog.CreateBlogResponse, error) {
	userId := ctx.Value("userId")
	if userId == nil {
		return nil, fmt.Errorf("Permission denied.")
	}

	blogToSave := &models.Blog{
		Title:       cbRequest.Title,
		Description: cbRequest.Description,
		Tags:        strings.Join(cbRequest.Tags, ","),
		AuthorId:    fmt.Sprintf("%v", userId),
		Category:    cbRequest.Category.String(),
	}

	result := bs.DB.Create(&blogToSave)
	if result.Error != nil {
		return nil, result.Error
	}

	authorId, _ := strconv.ParseInt(blogToSave.AuthorId, 0, 8)
	var author models.User
	author.ID = uint(authorId)

	authorResult := bs.DB.First(&author)
	if authorResult.Error != nil {
		return nil, authorResult.Error
	}

	return &blog.CreateBlogResponse{
		Id:          fmt.Sprint(blogToSave.ID),
		Title:       blogToSave.Title,
		Description: blogToSave.Description,
		Category:    cbRequest.Category,
		Tags:        strings.Split(blogToSave.Tags, ","),
		Author: &user.User{
			Id:    blogToSave.AuthorId,
			Name:  author.Name,
			Email: author.Email,
		},
	}, nil
}
