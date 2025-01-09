package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

const (
	streamName = "dataStream"
	groupName  = "consumerGroup"
	consumerID = "consumer1"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis_service:6379",
	})
	defer rdb.Close()

	ctx := context.Background()

	err := setupStream(ctx, rdb)
	if err != nil {
		log.Fatalf("Failed to set up stream: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go producer(ctx, rdb, &wg)
	go consumer(ctx, rdb, &wg)

	wg.Wait()
}

func setupStream(ctx context.Context, rdb *redis.Client) error {
	err := rdb.XGroupCreateMkStream(ctx, streamName, groupName, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}
	return nil
}

func producer(ctx context.Context, rdb *redis.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 10; i++ {
		_, err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{"value": i},
		}).Result()
		if err != nil {
			log.Printf("Failed to produce message: %v", err)
			return
		}
	}

	log.Println("Producer finished sending messages")
}

func consumer(ctx context.Context, rdb *redis.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		messages, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerID,
			Streams:  []string{streamName, ">"},
			Count:    10,
			Block:    0,
		}).Result()

		if err != nil {
			log.Printf("Failed to consume messages: %v", err)
			return
		}

		for _, stream := range messages {
			for _, message := range stream.Messages {
				val := message.Values["value"].(string)
				intVal, _ := strconv.Atoi(val)
				fmt.Println(intVal)

				_, err := rdb.XAck(ctx, streamName, groupName, message.ID).Result()
				if err != nil {
					log.Printf("Failed to acknowledge message: %v", err)
				}

				if intVal == 9 {
					log.Println("Consumer finished processing messages")
					return
				}
			}
		}
	}
}
