package main

import (
	"github.com/redis/go-redis/v9"
	"net/http"
	"valuator/package/app/query"
	"valuator/package/app/service"
	"valuator/package/infra/api"
	"valuator/package/infra/redis/repo"
	"valuator/package/infra/redis/unique"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	textRepo := repo.NewTextRepository(rdb)
	uniqueCounter := unique.NewUniqueCounter(rdb)
	textService := service.NewTextService(textRepo, uniqueCounter)
	statisticsQueryService := query.NewStatisticsQueryService(textRepo, uniqueCounter)
	textQueryService := query.NewTextQueryService(textRepo)

	handler := api.NewHandler(textService, statisticsQueryService, textQueryService)

	http.HandleFunc("/create/form", handler.CreateForm)
	http.HandleFunc("/process", handler.ProcessText)
	http.HandleFunc("/statics", handler.Statistics)
	http.HandleFunc("/delete", handler.Delete)
	http.HandleFunc("/list", handler.List)
	http.HandleFunc("/", handler.List)
	http.ListenAndServe(":8082", nil)
}
