package main

import (
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	// Подключение к серверу NATS
	natsConn, err := nats.Connect("nats://nats:4222")
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

	// Бесконечный цикл для удержания программы
	select {}
}
