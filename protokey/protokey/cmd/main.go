package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"protokey/pkg/app"
	"protokey/pkg/infra/inmemory"
	"protokey/pkg/infra/transport"
	"time"
)

func setupRoutes(handler *transport.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/get", handler.GetValue).Methods("GET")
	router.HandleFunc("/set", handler.SetValue).Methods("POST")
	router.HandleFunc("/keys", handler.ListKeys).Methods("GET")

	return router
}

func main() {
	commandLogFile := getEnv("COMMAND_LOG_FILE", "data/command.log")
	snapshotFile := getEnv("SNAPSHOT_FILE", "data/snapshot.json")

	commandChan := make(chan app.Command)
	responseChan := make(chan app.Response)

	storeConfig := inmemory.Config{
		CommandLogFile:   commandLogFile,
		SnapshotFile:     snapshotFile,
		SnapshotInterval: 5 * time.Minute,
		FlushInterval:    1 * time.Second,
	}

	store := inmemory.NewStore(storeConfig, commandChan, responseChan)
	defer store.Close()

	svc := app.NewProtoKeyService(commandChan, responseChan)

	httpHandler := transport.NewHandler(svc)
	router := setupRoutes(httpHandler)

	log.Println("Server is listening on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
