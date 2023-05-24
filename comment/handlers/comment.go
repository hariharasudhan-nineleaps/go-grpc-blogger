package handlers

import (
	"context"
	"fmt"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/comment"
	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/user"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/utils"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type CommentServer struct {
	comment.UnimplementedCommentServiceServer
	DB                *gorm.DB
	UserServiceClient user.UserServiceClient
}

func (cs *CommentServer) CreateComment(ctx context.Context, ccRequest *comment.CreateCommentRequest) (*comment.CreateCommentResponse, error) {
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, fmt.Errorf("Invalid userId")
	}

	commentToSave := &models.Comment{
		ID:       utils.ShortId(),
		Comment:  ccRequest.Comment,
		Entity:   ccRequest.Entity.String(),
		EntityID: ccRequest.EntityId,
		UserId:   userId,
	}
	result := cs.DB.Create(commentToSave)
	if result.Error != nil {
		return nil, fmt.Errorf("Unable to create comment %V", result.Error)
	}
	fmt.Printf("Comment with ID %v saved successfully", commentToSave.ID)

	return &comment.CreateCommentResponse{
		Id:       commentToSave.ID,
		Entity:   comment.CommentEntity(comment.CommentEntity_value[commentToSave.Entity]),
		EntityId: commentToSave.EntityID,
		Comment:  commentToSave.Comment,
	}, nil
}

func (cs *CommentServer) Comments(ctx context.Context, ccRequest *comment.CommentRequest) (*comment.CommentResponse, error) {

	userToken, ok := ctx.Value("userToken").(string)
	if !ok {
		return nil, fmt.Errorf("Invalid token")
	}

	// Fetch comments for blog
	var dbComments []models.Comment
	dbCommentsQuery := &models.Comment{
		Entity:   ccRequest.Entity.String(),
		EntityID: ccRequest.EntityId,
	}
	cs.DB.Where(dbCommentsQuery).Find(&dbComments)

	// Fetch users metadata from userID's
	var userIds []string
	for _, comment := range dbComments {
		userIds = append(userIds, comment.UserId)
	}

	md := metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %v", userToken),
	})
	ctx = metadata.NewOutgoingContext(ctx, md)
	usersResponse, err := cs.UserServiceClient.GetUsers(ctx, &user.GetUsersRequest{
		UserIds: userIds,
	})
	if err != nil {
		fmt.Print(err)
		return nil, fmt.Errorf("Error from userService %v", err)
	}
	usersMap := make(map[string]user.User, len(userIds))

	for _, resUser := range usersResponse.Users {
		usersMap[resUser.Id] = *resUser
	}

	var resComments []*comment.Comment
	for _, commentItem := range dbComments {
		userItem := usersMap[commentItem.UserId]
		fmt.Println("userItem", userItem)
		resComments = append(resComments, &comment.Comment{
			Id:        commentItem.ID,
			Comment:   commentItem.Comment,
			CreatedAt: timestamppb.New(commentItem.CreatedAt),
			User: &user.User{
				Id:    userItem.Id,
				Name:  userItem.Name,
				Email: userItem.Email,
			},
		})
	}
	fmt.Println("resComments")

	return &comment.CommentResponse{
		Metadata: &comment.CommentResponseMetadata{
			Total: uint32(len(resComments)),
		},
		Comments: resComments,
	}, nil
}
