package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/comment"
	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/user"
	handlers "github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/handlers"
	interceptors "github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/interceptors"
	models "github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
	db.AutoMigrate(&models.Comment{})

	// listen to incoming requests
	lis, err := net.Listen("tcp", cf.ServerEndpoint)
	if err != nil {
		log.Fatalf("Unable to connect server %v", err)
	}

	// Load Client Certificate and key
	clientCertificate, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		log.Fatalf("Failed to load client certificate and key. %s.", err)
	}

	// Load CA certificate
	trustedCertificate, err := ioutil.ReadFile("cacert.pem")
	if err != nil {
		log.Fatalf("Failed to load CA certificate. %s.", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(trustedCertificate) {
		log.Fatalf("Failed to load CA certificate pool. %s.", err)
	}

	// tls config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCertificate},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}
	cred := credentials.NewTLS(tlsConfig)

	// server create
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		interceptors.UnaryAuthInterceptor,
	))

	// user grpc client
	conn, err := grpc.Dial(cf.UserServiceEndpoint, grpc.WithTransportCredentials(cred))
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
