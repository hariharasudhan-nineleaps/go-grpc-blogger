package main

import (
	"log"
	"net"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/blog"
	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/user"
	handlers "github.com/hariharasudhan-nineleaps/go-grpc-blogger/blog/handlers"
	interceptors "github.com/hariharasudhan-nineleaps/go-grpc-blogger/blog/interceptors"
	models "github.com/hariharasudhan-nineleaps/go-grpc-blogger/blog/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/blog/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Blog{})

	// listen to incoming requests
	lis, err := net.Listen("tcp", cf.ServerEndpoint)
	if err != nil {
		log.Fatalf("Unable to connect server %v", err)
	}

	// grpc server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		interceptors.UnaryAuthInterceptor,
	))

	// user grpc client
	conn, err := grpc.Dial(cf.UserServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Unable to connect user service %v", err)
	}
	defer conn.Close()

	userServiceClient := user.NewUserServiceClient(conn)

	// register service
	blog.RegisterBlogServiceServer(grpcServer, &handlers.BlogServer{DB: db, UserServiceClient: userServiceClient})
	reflection.Register(grpcServer)

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
