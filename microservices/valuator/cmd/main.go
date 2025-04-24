package main

import (
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"valuator/package/app/query"
	"valuator/package/app/service"
	"valuator/package/infra/api"
	infranats "valuator/package/infra/nats"
	infraredis "valuator/package/infra/redis"
	"valuator/package/infra/redis/repo"
)

func main() {
	mainRedisClient := redis.NewClient(&redis.Options{
		Addr: "redis-main:6379",
	})
	ruRedisClient := redis.NewClient(&redis.Options{
		Addr: "redis-ru:6379",
	})
	enRedisClient := redis.NewClient(&redis.Options{
		Addr: "redis-en:6379",
	})
	asiaRedisClient := redis.NewClient(&redis.Options{
		Addr: "redis-asia:6379",
	})
	log.Println("ASDDSASDDSA")

	natsConn, err := nats.Connect("http://nats:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsConn.Close()

	redisProvider := infraredis.NewShardRepository(mainRedisClient, ruRedisClient, enRedisClient, asiaRedisClient)
	textRepo := repo.NewTextRepository(redisProvider)
	natsDispatcher := infranats.NewNatsDispatcher(natsConn)
	textService := service.NewTextService(textRepo, redisProvider, natsDispatcher)
	textQueryService := query.NewTextQueryService(textRepo)

	handler := api.NewHandler(textService, textQueryService)

	http.HandleFunc("/valuator/create/form", handler.CreateForm)
	http.HandleFunc("/valuator/process", handler.ProcessText)
	http.HandleFunc("/valuator/delete", handler.Delete)
	http.HandleFunc("/valuator/list", handler.List)
	http.HandleFunc("/valuator/", handler.List)
	http.ListenAndServe(":8082", nil)
}
