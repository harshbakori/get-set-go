package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func initRedis(address string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: address,
	})
}

func isUniqueRequest(id int) bool {
	val, err := rdb.SetNX(ctx, fmt.Sprintf("request_id_%d", id), 1, 1*time.Minute).Result()
	if err != nil {
		log.Printf("Failed to check/set unique ID in Redis: %v", err)
		return false
	}
	return val
}

func logUniqueRequests() {
	for {
		time.Sleep(1 * time.Minute)

		count, err := rdb.DBSize(ctx).Result()
		if err != nil {
			log.Printf("Failed to get unique request count from Redis: %v", err)
			continue
		}
		log.Printf("Unique request count in the last minute: %d", count)

		sendToKafka(count)

		rdb.FlushDB(ctx)
	}
}
