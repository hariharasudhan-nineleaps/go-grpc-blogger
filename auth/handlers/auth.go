package handlers

import (
	"context"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/auth"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/auth/config"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/auth/utils"
)

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	Config *config.Config
}

func (ah *AuthHandler) GenerateToken(ctx context.Context, req *auth.GenerateTokenRequest) (*auth.GenerateTokenResponse, error) {
	claims := &jwt.MapClaims{
		"id":     req.UserId,
		"userId": req.UserId,
		"name":   req.Name,
		"email":  req.Email,
	}

	token, err := utils.GenerateToken(claims, ah.Config.JWTSecret)
	if err != nil {
		log.Fatalf("Unable to generate token %v", err)
	}

	return &auth.GenerateTokenResponse{
		AccessToken: token,
		OpaqueToken: token,
	}, nil
}

func (ah *AuthHandler) VerifyToken(ctx context.Context, req *auth.VerifyTokenRequest) (*auth.VerifyTokenResponse, error) {
	return &auth.VerifyTokenResponse{
		IsValid:     true,
		AccessToken: req.OpaqueToken,
	}, nil
}
