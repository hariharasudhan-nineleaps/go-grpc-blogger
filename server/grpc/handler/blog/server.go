package blog

import (
	"context"
	"fmt"
	"strings"

	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/proto/blog"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/proto/user"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/utils"
	"gorm.io/gorm"
)

type BlogServer struct {
	blog.UnimplementedBlogServiceServer
	DB *gorm.DB
}

func (bs *BlogServer) CreateBlog(ctx context.Context, cbRequest *blog.CreateBlogRequest) (*blog.CreateBlogResponse, error) {
	userId, ok := ctx.Value("userId").(string)

	if !ok {
		return nil, fmt.Errorf("Invalid userId")
	}

	blogToSave := &models.Blog{
		ID:          utils.ShortId(),
		Title:       cbRequest.Title,
		Description: cbRequest.Description,
		Tags:        strings.Join(cbRequest.Tags, ","),
		AuthorId:    userId,
		Category:    cbRequest.Category.String(),
	}

	result := bs.DB.Create(&blogToSave)
	if result.Error != nil {
		return nil, result.Error
	}

	var author models.User
	author.ID = blogToSave.AuthorId

	authorResult := bs.DB.First(&author)
	if authorResult.Error != nil {
		return nil, authorResult.Error
	}

	return &blog.CreateBlogResponse{
		Id:          blogToSave.ID,
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

func (bs *BlogServer) GetUserBlogs(ctx context.Context, gubRequest *blog.UserBlogRequest) (*blog.UserBlogResponse, error) {
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, fmt.Errorf("Invalid userId")
	}

	fmt.Print(userId)

	return nil, fmt.Errorf("test")
}
