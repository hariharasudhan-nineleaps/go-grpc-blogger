package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/activity"
	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/blog"
	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/user"
	handlers "github.com/hariharasudhan-nineleaps/go-grpc-blogger/blog/handlers"
	interceptors "github.com/hariharasudhan-nineleaps/go-grpc-blogger/blog/interceptors"
	models "github.com/hariharasudhan-nineleaps/go-grpc-blogger/blog/models"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/blog/utils"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	// kafka
	const (
		topic     = "blog_view"
		partition = 0
	)
	kafkaConn, err := kafka.DialLeader(context.Background(), "tcp", cf.KafkaEndpoint, topic, partition)
	if err != nil {
		log.Fatalf("Unable to connect with kafka %v", err)
	}

	// Load Client Certificates
	clientCertificate, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		log.Fatalf("Failed to load client certificate and key. %s.", err)
	}

	// Load CA Certificate
	trustedCertificate, err := ioutil.ReadFile("cacert.pem")
	if err != nil {
		log.Fatalf("Failed to load CA certificate. %s.", err)
	}

	// Put the CA certificate to certificate pool
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(trustedCertificate) {
		log.Fatalf("Failed to load CA certificate to pool. %s.", err)
	}

	// Create the TLS configuration
	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{clientCertificate},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}
	cred := credentials.NewTLS(&tlsConfig)

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
	uconn, uerr := grpc.Dial(cf.UserServiceEndpoint, grpc.WithTransportCredentials(cred))
	if uerr != nil {
		log.Fatalf("Unable to connect user service %v", uerr)
	}
	defer uconn.Close()
	userServiceClient := user.NewUserServiceClient(uconn)

	// activity grpc client
	aconn, aerr := grpc.Dial(cf.ActivityServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if aerr != nil {
		log.Fatalf("Unable to connect activity service %v", aerr)
	}
	defer aconn.Close()
	activityServiceClient := activity.NewActivityServiceClient(aconn)

	// register service
	blog.RegisterBlogServiceServer(grpcServer, &handlers.BlogServer{DB: db, UserServiceClient: userServiceClient, KafkaConn: kafkaConn, ActivityServiceClient: activityServiceClient})
	reflection.Register(grpcServer)

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
