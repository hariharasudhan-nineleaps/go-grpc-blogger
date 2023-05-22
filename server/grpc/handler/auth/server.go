package auth

import (
	context "context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/proto/auth"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/utils"
	"gorm.io/gorm"
)

type AuthServer struct {
	DB *gorm.DB
	auth.UnimplementedAuthServiceServer
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
		Email:    email,
		Name:     getNameFromEmail(email),
		Password: password,
	}
}

func (a *AuthServer) Login(ctx context.Context, authRequest *auth.AuthRequest) (*auth.AuthResponse, error) {

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

	token, err := utils.GenerateToken(&jwt.MapClaims{
		"id":    dbUser.ID,
		"name":  dbUser.Name,
		"email": dbUser.Email,
	}, "secret")
	if err != nil {
		log.Fatalf("Token generation failed %v", err)
	}

	return &auth.AuthResponse{
		Id:          fmt.Sprintf("%v", dbUser.ID),
		AccessToken: token,
		Name:        dbUser.Name,
		Email:       dbUser.Email,
	}, nil
}
