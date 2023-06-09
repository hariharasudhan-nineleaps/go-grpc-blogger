package handlers

import (
	context "context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/auth"
	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/user"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/utils"
	"gorm.io/gorm"
)

type UserServer struct {
	DB                *gorm.DB
	AuthServiceClient auth.AuthServiceClient
	user.UnimplementedUserServiceServer
}

func getNameFromEmail(email string) string {
	split := strings.Split(email, "@")
	return split[0]
}

func buildUser(email string, plainPassword string) models.User {
	password, err := utils.HashPassword(plainPassword)
	if err != nil {
		log.Fatalf("Password hash failed %v", err)
	}

	return models.User{
		ID:       utils.ShortId(),
		Email:    email,
		Name:     getNameFromEmail(email),
		Password: password,
	}
}

func (a *UserServer) Login(ctx context.Context, authRequest *user.AuthRequest) (*user.AuthResponse, error) {

	// check user already exists
	var dbUser models.User
	dbUser.Email = authRequest.Email
	res := a.DB.First(&dbUser)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Printf("No existing record with email creating new record")
		user := buildUser(authRequest.Email, authRequest.Password)
		result := a.DB.Create(&user)
		if result.Error != nil {
			log.Fatalf("Error saving record %v", result.Error)
		}

		dbUser = user
	} else if res.Error != nil {
		log.Fatalf("Error saving record %v", res.Error)
	} else {
		log.Printf("record with email %v found", authRequest.Email)
	}

	token, err := a.AuthServiceClient.GenerateToken(context.Background(), &auth.GenerateTokenRequest{
		Email:  dbUser.Email,
		UserId: dbUser.Email,
		Name:   dbUser.Name,
	})
	if err != nil {
		log.Fatalf("Token generation failed %v", err)
	}

	return &user.AuthResponse{
		Id:          dbUser.ID,
		AccessToken: token.OpaqueToken,
		Name:        dbUser.Name,
		Email:       dbUser.Email,
	}, nil
}

func (a *UserServer) GetUsers(ctx context.Context, authRequest *user.GetUsersRequest) (*user.GetUsersResponse, error) {
	// check user already exists
	var dbUsers []models.User
	a.DB.Find(&dbUsers, authRequest.UserIds)

	// map db users to proto-buf user
	var resUsers []*user.User
	for _, dbUser := range dbUsers {
		resUsers = append(resUsers, &user.User{
			Id:    dbUser.ID,
			Name:  dbUser.Name,
			Email: dbUser.Email,
		})
	}

	return &user.GetUsersResponse{
		Users: resUsers,
	}, nil
}

func (a *UserServer) GetUser(ctx context.Context, authRequest *user.GetUserRequest) (*user.User, error) {
	var dbUser models.User
	dbUser.ID = authRequest.UserId
	result := a.DB.First(&dbUser)

	if result.Error != nil {
		return nil, fmt.Errorf("User with ID %v not exists", authRequest.UserId)
	}

	return &user.User{
		Id:    dbUser.ID,
		Name:  dbUser.Name,
		Email: dbUser.Email,
	}, nil
}
