package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

func main() {
	// Подключение к серверу NATS
	natsConn, err := nats.Connect("nats://" + os.Getenv("NATS_USERNAME") + ":" + os.Getenv("NATS_PASSWORD") + "@nats:4222")

	if err != nil {
		log.Fatalf("Не удалось подключиться к NATS: %v", err)
	}
	defer natsConn.Close()

	log.Println("Успешное подключение к серверу NATS")

	// Подписка на все сообщения с использованием подстановочного знака ">"
	_, err = natsConn.Subscribe(">", func(msg *nats.Msg) {
		log.Printf("Получено сообщение из темы '%s': %s", msg.Subject, string(msg.Data))
	})
	if err != nil {
		log.Fatalf("Не удалось подписаться на все темы: %v", err)
	}

	log.Println("Подписка на все темы выполнена успешно")

	var forever chan struct{}

	<-forever
}
