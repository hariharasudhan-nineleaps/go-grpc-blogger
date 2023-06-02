package main

import (
	"context"
	"log"
	"sync"

	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/metric/handler"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/metric/utils"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
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

	// kafka
	const (
		topic     = "blog_view"
		partition = 0
	)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{cf.KafkaEndpoint},
		Topic:     topic,
		Partition: partition,
	})

	// listen kafka events
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		kh := handler.BlogAsyncHandler{
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
