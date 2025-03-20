package main

import (
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"rankcalculator/package/app/calculator"
	"rankcalculator/package/app/command"
	"rankcalculator/package/app/service"
	"rankcalculator/package/infra/api"
	nats2 "rankcalculator/package/infra/nats"
	"rankcalculator/package/infra/redis/provider"
	"rankcalculator/package/infra/redis/repo"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	log.Println("ASDDSASDDSA")

	textProvider := provider.NewTextProvider(rdb)
	rankCalculator := calculator.NewRankCalculator(textProvider)
	rankRepository := repo.NewTextStatisticsRepository(rdb)
	rankService := service.NewStatisticsService(rankRepository, rankCalculator)

	natsConn, err := nats.Connect("http://nats:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsConn.Close()

	commandHandler := command.NewHandler(rankService)

	natsHandler := nats2.NewNATSHandler(natsConn, commandHandler)
	if err := natsHandler.Start(); err != nil {
		log.Fatalf("Failed to start NATS handler: %v", err)
	}
	log.Println("NATS handler is running...")

	handler := api.NewHandler(rankRepository)
	http.HandleFunc("/rankcalculator/statics", handler.Statistics)
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
