package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

// Redis on and defer Close()
func RedisOn() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		// PoolSize: 10,
	})

	// defer client.Close()

	
	// ctx:= context.Background()
	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println("ping result: ", pong)
	}
	RC = client
}