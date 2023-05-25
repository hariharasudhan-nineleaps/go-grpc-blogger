package main

import (
	"log"
	"net"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/user"
	handlers "github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/handlers"
	interceptors "github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/interceptors"
	models "github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// auth server

func main() {

	// load env
	cf, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Unable to connect server %v", err)
	}

	// dependencies
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{})

	// listen to incoming requests
	lis, err := net.Listen("tcp", cf.ServerEndpoint)
	if err != nil {
		log.Fatalf("Unable to connect server %v", err)
	}

	// server create
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		interceptors.UnaryAuthInterceptor,
	))

	// register service
	user.RegisterUserServiceServer(grpcServer, &handlers.UserServer{DB: db})
	reflection.Register(grpcServer)

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
