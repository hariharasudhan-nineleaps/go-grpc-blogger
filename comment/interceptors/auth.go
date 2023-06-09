package interceptors

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Ignore login route from auth
	if info.FullMethod == "/user.UserService/Login" {
		log.Println("--> Skipping auth for method: ", info.FullMethod)
		return handler(ctx, req)
	}

	log.Println("--> Verifying auth for method: ", info.FullMethod)

	// fetch metadata based on protocol
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata not exists")
	}

	// check auth header
	values, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	// trim Bearer from string
	token, ok := strings.CutPrefix(values[0], "Bearer ")
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	// get cliams from token
	claims, err := utils.VerifyToken(token, "secret")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated -> Invalid token!")
	}

	// fetch userId and attach to context which will used in handlers.
	userId, ok := claims["id"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated -> Invalid token (userId)!")
	}

	// set context
	ctx = context.WithValue(ctx, "userId", userId.(string))
	ctx = context.WithValue(ctx, "userToken", token)

	fmt.Println("[Comment Service]===> AUTH - SUCCESS")
	return handler(ctx, req)
}
