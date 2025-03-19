package main

import (
	"github.com/redis/go-redis/v9"
	"rankcalculator/package/infra/redis/provider"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	textProvider := provider.NewTextProvider(rdb)
}
