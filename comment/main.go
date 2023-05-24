package main

import (
	"log"
	"net"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/comment"
	handlers "github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/handlers"
	models "github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/models"
	interceptor "github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// auth server

func main() {

	// dependencies
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Comment{})

	// listen to incoming requests
	lis, err := net.Listen("tcp", "localhost:3003")
	if err != nil {
		log.Fatalf("Unable to connect server %v", err)
	}

	// server create
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		interceptor.UnaryAuthInterceptor,
	))

	// register service
	comment.RegisterCommentServiceServer(grpcServer, &handlers.CommentServer{DB: db})
	reflection.Register(grpcServer)

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
