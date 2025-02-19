package main

import (
	"net/http"
	"valuator/package/app/query"
	"valuator/package/app/service"
	"valuator/package/infra/api"
	"valuator/package/infra/redis/repo"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Инициализация клиента Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	textRepo := repo.NewTextRepository(rdb)
	ts := service.NewTextService(textRepo)
	ss := query.NewStatisticsQueryService(textRepo)
	handler := api.NewHandler(ts, ss)

	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/process", handler.ProcessText)
	http.ListenAndServe(":8082", nil)
}
