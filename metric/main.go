package main

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/activity"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/metric/handler"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/metric/utils"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// ctx
	ctx := context.Background()
	var wg sync.WaitGroup

	// load env
	cf, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Unable to connect server %v", err)
	}

	// redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cf.RedisEndpoint,
		Password: "",
		DB:       0,
	})

	// listen to incoming requests
	lis, err := net.Listen("tcp", cf.ServerEndpoint)
	if err != nil {
		log.Fatalf("Unable to connect server %v", err)
	}

	// grpc
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		grpcServer := grpc.NewServer()
		activity.RegisterActivityServiceServer(grpcServer, &handler.BlogActivityServerHandler{
			RDB: rdb,
		})
		reflection.Register(grpcServer)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			wg.Done()
		}
	}(&wg)

	// kafka
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		const (
			topic     = "blog_view"
			partition = 0
		)
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{cf.KafkaEndpoint},
			Topic:     topic,
			Partition: partition,
		})

		kh := handler.BlogActivityAsyncHandler{
			RDB: rdb,
		}

		for {
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				break
			}
			kh.HandleMessage(&msg)
		}

		if err := r.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}(&wg)

	wg.Wait()
}
