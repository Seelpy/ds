package main

import (
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"rankcalculator/package/app/calculator"
	"rankcalculator/package/app/command"
	"rankcalculator/package/app/service"
	"rankcalculator/package/infra/api"
	"rankcalculator/package/infra/centrifugo"
	nats2 "rankcalculator/package/infra/nats"
	infraredis "rankcalculator/package/infra/redis"
	"rankcalculator/package/infra/redis/provider"
	"rankcalculator/package/infra/redis/repo"
	"rankcalculator/package/infra/redis/unique"
)

func main() {
	mainRedisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_MAIN_ADR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Username: os.Getenv("REDIS_USERNAME"),
	})
	ruRedisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_RU_ADR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Username: os.Getenv("REDIS_USERNAME"),
	})
	enRedisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_EN_ADR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Username: os.Getenv("REDIS_USERNAME"),
	})
	asiaRedisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ASIA_ADR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Username: os.Getenv("REDIS_USERNAME"),
	})

	log.Println("ASDDSASDDSA")

	redisProvider := infraredis.NewShardRepository(mainRedisClient, ruRedisClient, enRedisClient, asiaRedisClient)

	textProvider := provider.NewTextProvider(redisProvider)
	counter := unique.NewUniqueCounter(mainRedisClient)
	rankCalculator := calculator.NewRankCalculator(counter)
	rankRepository := repo.NewTextStatisticsRepository(redisProvider)
	rankService := service.NewStatisticsService(rankRepository, rankCalculator, counter, textProvider)

	centrifugoClient := centrifugo.NewCentrifugoClient()

	natsConn, err := nats.Connect("nats://" + os.Getenv("NATS_USERNAME") + ":" + os.Getenv("NATS_PASSWORD") + "@nats:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsConn.Close()

	commandHandler := command.NewHandler(rankService, centrifugoClient)

	natsHandler := nats2.NewNATSHandler(natsConn, commandHandler)
	if err := natsHandler.Start(); err != nil {
		log.Fatalf("Failed to start NATS handler: %v", err)
	}
	log.Println("NATS handler is running...")

	handler := api.NewHandler(rankRepository, os.Getenv("SECRET"))
	http.HandleFunc("/rankcalculator/statics", handler.AuthMiddleware(handler.Statistics))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
