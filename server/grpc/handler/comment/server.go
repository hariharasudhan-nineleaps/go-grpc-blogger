package comment

import (
	"context"
	"fmt"

	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/proto/comment"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/proto/user"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type CommentServer struct {
	comment.UnimplementedCommentServiceServer
	DB *gorm.DB
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
	_, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, fmt.Errorf("Invalid userId")
	}

	// Fetch comments for blog
	var dbComments []*models.Comment
	dbCommentsQuery := &models.Comment{
		Entity:   ccRequest.Entity.String(),
		EntityID: ccRequest.EntityId,
	}
	cs.DB.Where(&dbCommentsQuery).Find(dbComments)

	// Fetch userIds to get user metadata
	var userIds []string
	var dbUsers []*models.User
	var dbUsersMap map[string]*models.User
	for _, comment := range dbComments {
		userIds = append(userIds, comment.UserId)
	}
	cs.DB.Find(&dbUsers, userIds)
	for _, dbUser := range dbUsers {
		dbUsersMap[dbUser.ID] = dbUser
	}

	var resComments []*comment.Comment
	for _, commentItem := range dbComments {
		userItem := dbUsersMap[commentItem.UserId]
		resComments = append(resComments, &comment.Comment{
			Id:        commentItem.ID,
			Comment:   commentItem.Comment,
			CreatedAt: timestamppb.New(commentItem.CreatedAt),
			User: &user.User{
				Id:    userItem.ID,
				Name:  userItem.Name,
				Email: userItem.Email,
			},
		})
	}

	return &comment.CommentResponse{
		Metadata: &comment.CommentResponseMetadata{
			Total: uint32(len(resComments)),
		},
		Comments: resComments,
	}, nil
}
