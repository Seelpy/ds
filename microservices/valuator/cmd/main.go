package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	redis "github.com/redis/go-redis/v9"
)

type User struct {
	Name  string
	Email string
}

func main() {
	// Инициализация клиента Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Обработчик HTTP-запросов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// Получение данных из Redis
		name, err := rdb.Get(ctx, "user:name").Result()
		if err != nil {
			log.Printf("Ошибка при получении имени: %v", err)
			name = "Гость"
		}

		email, err := rdb.Get(ctx, "user:email").Result()
		if err != nil {
			log.Printf("Ошибка при получении email: %v", err)
			email = "нет данных"
		}

		user := User{Name: name, Email: email}

		// Шаблон HTML
		tmpl := `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Профиль пользователя</title>
        </head>
        <body>
            <h1>Профиль пользователя</h1>
            <p>Имя: {{.Name}}</p>
            <p>Email: {{.Email}}</p>
        </body>
        </html>
        `

		// Парсинг и выполнение шаблона
		t, err := template.New("userProfile").Parse(tmpl)
		if err != nil {
			log.Printf("Ошибка при парсинге шаблона: %v", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, user)
		if err != nil {
			log.Printf("Ошибка при выполнении шаблона: %v", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Сервер запущен на http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
