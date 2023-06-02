package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/activity"
	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/blog"
	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/user"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/blog/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/blog/utils"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type BlogServer struct {
	blog.UnimplementedBlogServiceServer
	DB                    *gorm.DB
	KafkaConn             *kafka.Conn
	UserServiceClient     user.UserServiceClient
	ActivityServiceClient activity.ActivityServiceClient
}

type BlogViewPayload struct {
	BlogId string `json:"blogId"`
	UserId string `json:"userId"`
}

func (bs *BlogServer) CreateBlog(ctx context.Context, cbRequest *blog.CreateBlogRequest) (*blog.CreateBlogResponse, error) {
	userId, userIdOk := ctx.Value("userId").(string)
	userToken, userTokenOk := ctx.Value("userToken").(string)
	if !userIdOk || !userTokenOk {
		return nil, fmt.Errorf("Invalid userId or userToken")
	}

	md := metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %v", userToken),
	})

	ctx = metadata.NewOutgoingContext(ctx, md)
	resUser, err := bs.UserServiceClient.GetUser(ctx, &user.GetUserRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, fmt.Errorf("User service call failed %v", err)
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

	return &blog.CreateBlogResponse{
		Id:          blogToSave.ID,
		Title:       blogToSave.Title,
		Description: blogToSave.Description,
		Category:    cbRequest.Category,
		Tags:        strings.Split(blogToSave.Tags, ","),
		Author:      resUser,
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

func (bs *BlogServer) ViewBlog(ctx context.Context, vR *blog.ViewBlogRequest) (*blog.ViewBlogResponse, error) {
	// check user in ctx
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, fmt.Errorf("Invalid userId")
	}

	// check blog with ID exists
	var blogItem models.Blog
	blogItem.ID = vR.BlogId
	res := bs.DB.First(&blogItem)
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("Invalid blog Id %v", vR.BlogId)
	}

	// write to kafka
	blogViewJson, bJsonerr := json.Marshal(BlogViewPayload{
		BlogId: vR.BlogId,
		UserId: userId,
	})
	if bJsonerr != nil {
		log.Fatalf("Unable to create json %v", bJsonerr)
	}

	_, fWriteMsgerr := bs.KafkaConn.WriteMessages(kafka.Message{
		Headers: []kafka.Header{{Key: "type", Value: []byte("BLOG_VIEW")}},
		Value:   blogViewJson,
	})
	if fWriteMsgerr != nil {
		log.Fatalf("Unable to send message to kafka %v", fWriteMsgerr)
	}

	return &blog.ViewBlogResponse{}, nil
}

func (bs *BlogServer) MostViewedBlogs(ctx context.Context, mvBReq *blog.MostViewedBlogsRequest) (*blog.MostViewedBlogsResponse, error) {
	_, userIdOk := ctx.Value("userId").(string)
	userToken, userTokenOk := ctx.Value("userToken").(string)
	if !userIdOk || !userTokenOk {
		return nil, fmt.Errorf("Invalid userId or userToken")
	}

	md := metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %v", userToken),
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	ares, aerr := bs.ActivityServiceClient.MostViewedBlogIds(ctx, &activity.MostViewedBlogIdsRequest{})
	if aerr != nil {
		log.Fatalf("Activity service call failed %v", aerr)
	}
	blogIds := ares.BlogIds
	var blogs []models.Blog
	bs.DB.Find(&blogs, blogIds)

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

	return &blog.MostViewedBlogsResponse{
		Metadata: &blog.MostViewedBlogsResponseMetadata{
			Total: uint32(len(resBlogs)),
		},
		Blogs: resBlogs,
	}, nil
}
