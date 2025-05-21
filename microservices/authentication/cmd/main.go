package main

import (
	"authentication/package/infra/api"
	"authentication/package/infra/redis/repo"
	"github.com/redis/go-redis/v9"
	"net/http"
)

func main() {
	mainRedisClient := redis.NewClient(&redis.Options{
		Addr: "redis-main:6379",
	})
	userRepositury := repo.NewUserRepository(mainRedisClient)

	handler := api.NewHandler(userRepositury, "secret")
	http.HandleFunc("/authentication/login", handler.Login)
	http.HandleFunc("/authentication/registration", handler.Registration)
	http.HandleFunc("/authentication/token/refresh", handler.RefreshToken)
	http.ListenAndServe(":8082", nil)
}
