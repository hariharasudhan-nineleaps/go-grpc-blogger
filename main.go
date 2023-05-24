package main

import (
	"log"
	"net"

	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/models"
	authServerHandler "github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/handler/auth"
	blogServerHandler "github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/handler/blog"
	commentServerHandler "github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/handler/comment"
	interceptor "github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/interceptor"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/proto/auth"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/proto/blog"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/proto/comment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// auth server

func main() {

	// dependencies
	db, err := gorm.Open(sqlite.Open("blogger.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(
		&models.User{},
		&models.Blog{},
		&models.Comment{},
	)

	// listen to incoming requests
	lis, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalf("Unable to connect server %v", err)
	}

	// server create
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		interceptor.UnaryAuthInterceptor,
	))

	// register service
	auth.RegisterAuthServiceServer(grpcServer, &authServerHandler.AuthServer{DB: db})
	blog.RegisterBlogServiceServer(grpcServer, &blogServerHandler.BlogServer{DB: db})
	comment.RegisterCommentServiceServer(grpcServer, &commentServerHandler.CommentServer{DB: db})
	reflection.Register(grpcServer)

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
