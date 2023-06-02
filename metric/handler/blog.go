package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hariharasudhan-nineleaps/blogger-proto/grpc/proto/activity"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type BlogActivityAsyncHandler struct {
	RDB *redis.Client
}

type BlogActivityServerHandler struct {
	activity.UnimplementedActivityServiceServer
	RDB *redis.Client
}

type BlogViewPayload struct {
	BlogId string `json:"blogId"`
	UserId string `json:"userId"`
}

func (bh *BlogActivityAsyncHandler) HandleMessage(msg *kafka.Message) {
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

func (bh *BlogActivityAsyncHandler) handleBlogView(msg *kafka.Message) {
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

func (bh *BlogActivityServerHandler) MostViewedBlogIds(ctx context.Context, mvRequest *activity.MostViewedBlogIdsRequest) (*activity.MostViewedBlogIdsResponse, error) {
	values, err := bh.RDB.ZRevRangeByScoreWithScores(context.Background(), "blog_views", &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  10,
	}).Result()

	var blogIds []string
	for _, blogId := range values {
		blogIds = append(blogIds, blogId.Member.(string))
	}

	if err != nil {
		log.Fatalf("Unable to fetch top blogs %v", err)
	}

	return &activity.MostViewedBlogIdsResponse{
		BlogIds: blogIds,
	}, nil
}
