package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type BlogAsyncHandler struct {
	RDB *redis.Client
}

type BlogViewPayload struct {
	BlogId string `json:"blogId"`
	UserId string `json:"userId"`
}

func (bh *BlogAsyncHandler) HandleMessage(msg *kafka.Message) {
	var msgType string
	for _, header := range msg.Headers {
		if header.Key == "type" {
			msgType = string(header.Value)
		}
	}

	switch msgType {
	case "BLOG_VIEW":
		fmt.Printf("Handling msg %v \n", string(msg.Value))
		bh.handleBlogView(msg)
	default:
		log.Fatalf("Invalid message type.")
	}
}

func (bh *BlogAsyncHandler) handleBlogView(msg *kafka.Message) {
	blogViewPayload := &BlogViewPayload{}
	uerr := json.Unmarshal(msg.Value, &blogViewPayload)
	if uerr != nil {
		log.Fatalf("Unable to parse json %v", uerr)
	}

	_, rerr := bh.RDB.ZIncrBy(context.Background(), "blog_views", 1, blogViewPayload.BlogId).Result()
	if rerr != nil {
		log.Fatalf("Unable to write to redis %v", rerr)
	}
}
