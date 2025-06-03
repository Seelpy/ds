package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"protokey/pkg/app"
	"protokey/pkg/infra"
	"protokey/pkg/infra/transport"
	"syscall"
	"time"
)

func main() {
	snapshotFile := getEnv("SNAPSHOT_FILE", "data/snapshot.data")
	httpPort := getEnv("HTTP_PORT", "8080")
	snapshotDelay := getDurationEnv("SNAPSHOT_DELAY", time.Second)

	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	store := infra.NewStore()
	snapshotService, err := infra.NewFileSnapshotService(snapshotFile)
	if err != nil {
		log.Fatalf("Failed to create snapshot service: %v", err)
	}

	commandChan := make(chan app.Command)
	responseChan := make(chan app.Response)
	stopChan := make(chan struct{})

	engine, err := app.NewEngine(store, snapshotService, commandChan, responseChan, stopChan, snapshotDelay)
	if err != nil {
		log.Fatalf("Failed to create engine: %v", err)
	}

	engine.Run()

	service := app.NewProtoKeyService(commandChan, responseChan)

	handler := transport.NewHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/get", handler.GetValue).Methods("GET")
	router.HandleFunc("/set", handler.SetValue).Methods("POST")
	router.HandleFunc("/keys", handler.ListKeys).Methods("GET")

	server := &http.Server{
		Addr:         ":" + httpPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Server is shutting down...")

		close(stopChan)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v", err)
		}
		close(done)
	}()

	log.Printf("Server starting on port %s", httpPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	<-done
	log.Println("Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		duration, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("Invalid duration for %s: %v, using default", key, err)
			return defaultValue
		}
		return duration
	}
	return defaultValue
}
