// redis_connection.go
package databases

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

func ConnectToRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ping:", pong)

	return client
}
