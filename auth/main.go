package main

import (
	"log"
	"net"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/auth"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/auth/handlers"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/auth/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cf, cerr := utils.LoadConfig(".")
	if cerr != nil {
		log.Fatalf("Failed to Load env %v", cerr)
	}

	lis, nlerr := net.Listen("tcp", cf.ServerEndpoint)
	if nlerr != nil {
		log.Fatalf("Unable to listem %v", nlerr)
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, &handlers.AuthHandler{Config: cf})
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Unable to serve GRPC server %v", err)
	}
}
