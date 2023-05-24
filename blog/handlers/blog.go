package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/blog"
	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/user"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	var blogs []models.Blog
	bs.DB.Where(&models.Blog{AuthorId: userId}).Find(&blogs)

	var resBlogs []*blog.UserBlog

	for _, blogItem := range blogs {
		resBlogs = append(resBlogs, &blog.UserBlog{
			Id:          blogItem.ID,
			Title:       blogItem.Title,
			Description: blogItem.Description,
			Category:    blog.BlogCategory(blog.BlogCategory_value[blogItem.Category]),
			Tags:        strings.Split(blogItem.Tags, ","),
			CreatedAt:   timestamppb.New(blogItem.CreatedAt),
		})
	}

	return &blog.UserBlogResponse{
		Metadata: &blog.UserBlogResponseMetadata{Total: uint32(len(blogs))},
		Blogs:    resBlogs,
	}, nil
}

func (bs *BlogServer) GetUserBlog(ctx context.Context, gubRequest *blog.GetUserBlogRequest) (*blog.UserBlog, error) {
	_, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, fmt.Errorf("Invalid userId")
	}

	var blogItem models.Blog
	blogId := gubRequest.BlogId
	blogItem.ID = blogId
	res := bs.DB.First(&blogItem)

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("Invalid blog Id %v", blogId)
	}

	return &blog.UserBlog{
		Id:          blogItem.ID,
		Title:       blogItem.Title,
		Description: blogItem.Description,
		Category:    blog.BlogCategory(blog.BlogCategory_value[blogItem.Category]),
		Tags:        strings.Split(blogItem.Tags, ","),
		CreatedAt:   timestamppb.New(blogItem.CreatedAt),
	}, nil
}
