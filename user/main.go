package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/user"
	handlers "github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/handlers"
	models "github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/user/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	// Load Server certificate and its key
	serverCertificate, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		log.Fatalf("Failed to load server certificate and key. %s.", err)
	}

	// Load CA certificate
	trustedCertificate, err := ioutil.ReadFile("cacert.pem")
	if err != nil {
		log.Fatalf("Failed to load trusted certificate. %s.", err)
	}

	// Put the CA certificate to certificate pool
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(trustedCertificate) {
		log.Fatalf("Failed to append trusted certificate to certificate pool. %s.", err)
	}

	// create tls credentials
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCertificate},
		ClientCAs:    certPool,
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}
	cred := credentials.NewTLS(tlsConfig)

	// listen to incoming requests
	lis, err := net.Listen("tcp", cf.ServerEndpoint)
	if err != nil {
		log.Fatalf("Unable to connect server %v", err)
	}

	// server create
	// grpcServer := grpc.NewServer(grpc.Creds(cred), grpc.UnaryInterceptor(
	// 	interceptors.UnaryAuthInterceptor,
	// ))
	grpcServer := grpc.NewServer(grpc.Creds(cred))

	// register service
	user.RegisterUserServiceServer(grpcServer, &handlers.UserServer{DB: db})
	reflection.Register(grpcServer)

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	} else {
		fmt.Printf("Server started: %v", cf.ServerEndpoint)
	}

}
