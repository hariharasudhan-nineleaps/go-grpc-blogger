package main

import (
	"fmt"
	"log"
	"net"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/auth"
	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/user"
	handlers "github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/handlers"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/interceptors"
	models "github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// auth server

func main() {
	// load env
	cf, cerr := utils.LoadConfig(".")
	if cerr != nil {
		log.Fatalf("Unable to connect server %v", cerr)
	}

	// dependencies
	db, dberr := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if dberr != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{})

	// listen to incoming requests
	lis, lerr := net.Listen("tcp", cf.ServerEndpoint)
	if lerr != nil {
		log.Fatalf("Unable to listen server %v", lerr)
	}

	// auth grpc client
	aconn, aerr := grpc.Dial(cf.AuthServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if aerr != nil {
		log.Fatalf("Unable to dial auth server %v", aerr)
	}
	authClient := auth.NewAuthServiceClient(aconn)

	// server create
	userServiceInterceptors := interceptors.Interceptor{
		AuthServiceClient: authClient,
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(userServiceInterceptors.UnaryAuthInterceptor))
	user.RegisterUserServiceServer(grpcServer, &handlers.UserServer{DB: db, AuthServiceClient: authClient})
	reflection.Register(grpcServer)

	fmt.Print("hahahhahah", cf.ServerEndpoint)

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
