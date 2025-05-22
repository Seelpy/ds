package main

import (
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"valuator/package/app/query"
	"valuator/package/app/service"
	"valuator/package/infra/api"
	infranats "valuator/package/infra/nats"
	infraredis "valuator/package/infra/redis"
	"valuator/package/infra/redis/repo"
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

	natsConn, err := nats.Connect("http://nats:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsConn.Close()

	redisProvider := infraredis.NewShardRepository(mainRedisClient, ruRedisClient, enRedisClient, asiaRedisClient)
	textRepo := repo.NewTextRepository(redisProvider)
	natsDispatcher := infranats.NewNatsDispatcher(natsConn)
	textService := service.NewTextService(textRepo, natsDispatcher)
	textQueryService := query.NewTextQueryService(textRepo)

	handler := api.NewHandler(textService, textQueryService, "secret")

	http.HandleFunc("/valuator/create/form", handler.AuthMiddleware(handler.CreateForm))
	http.HandleFunc("/valuator/login/form", handler.Login)
	http.HandleFunc("/valuator/process", handler.AuthMiddleware(handler.ProcessText))
	http.HandleFunc("/valuator/delete", handler.AuthMiddleware(handler.Delete))
	http.HandleFunc("/valuator/list", handler.AuthMiddleware(handler.List))
	http.HandleFunc("/valuator/", handler.AuthMiddleware(handler.List))
	http.ListenAndServe(":8082", nil)
}
