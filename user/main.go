package main

import (
	"log"
	"net"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/auth"
	interceptor "github.com/hariharasudhan-nineleaps/go-grpc-blogger/server/grpc/interceptor"
	handlers "github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/handlers"
	models "github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/models"
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
	db.AutoMigrate(&models.User{})

	// listen to incoming requests
	lis, err := net.Listen("tcp", "localhost:3001")
	if err != nil {
		log.Fatalf("Unable to connect server %v", err)
	}

	// server create
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		interceptor.UnaryAuthInterceptor,
	))

	// register service
	auth.RegisterAuthServiceServer(grpcServer, &handlers.AuthServer{DB: db})
	reflection.Register(grpcServer)

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
