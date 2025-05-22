package main

import (
	"authentication/package/infra/api"
	"authentication/package/infra/redis/repo"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
)

func main() {
	mainRedisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_MAIN_ADR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Username: os.Getenv("REDIS_USERNAME"),
	})
	userRepositury := repo.NewUserRepository(mainRedisClient)

	handler := api.NewHandler(userRepositury, os.Getenv("SECRET"))
	http.HandleFunc("/authentication/login", handler.Login)
	http.HandleFunc("/authentication/registration", handler.Registration)
	http.HandleFunc("/authentication/token/refresh", handler.RefreshToken)
	http.ListenAndServe(":8082", nil)
}
