package main

import (
	"log"
	"net"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/comment"
	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/user"
	handlers "github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/handlers"
	interceptors "github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/interceptors"
	models "github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		interceptors.UnaryAuthInterceptor,
	))

	// user grpc client
	conn, err := grpc.Dial("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Unable to connect user service %v", err)
	}
	defer conn.Close()

	userServiceClient := user.NewUserServiceClient(conn)

	// register service
	comment.RegisterCommentServiceServer(grpcServer, &handlers.CommentServer{DB: db, UserServiceClient: userServiceClient})
	reflection.Register(grpcServer)

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
