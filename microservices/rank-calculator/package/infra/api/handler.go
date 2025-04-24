package api

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"html/template"
	"log"
	"net/http"
	"rankcalculator/package/app/model"
	"rankcalculator/package/app/notification"
	"time"
)

func NewHandler(rankRepo model.ReadOnlyTextStatisticsRepository) *Handler {
	return &Handler{
		rankRepo: rankRepo,
	}
}

type Handler struct {
	rankRepo model.ReadOnlyTextStatisticsRepository
}

func (h *Handler) Statistics(w http.ResponseWriter, r *http.Request) {
	log.Println("HANDLE")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		http.Error(w, "Failed to get summary", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./data/html/summary.html")
	if err != nil {
		http.Error(w, "server error1"+err.Error(), http.StatusInternalServerError)
		return
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	rank, err := h.rankRepo.Get(id)
	if errors.Is(err, model.ErrStatisticsNotFound) {
		channel := notification.GenerateChannel(id)
		data := struct {
			Text            string
			Rank            float64
			Similarity      int
			CentrifugoToken string
			CentrifugoURL   string
			Channel         string
			ProcessingID    string
			HasResult       bool
		}{
			Text:            "результаты",
			CentrifugoToken: generateCentrifugoToken(ip, channel),
			CentrifugoURL:   "ws://localhost:8000/connection/websocket",
			Channel:         channel,
			ProcessingID:    id.String(),
		}
		err = tmpl.Execute(w, data)
	} else {
		similarity := 0
		if rank.IsDuplicate {
			similarity = 1
		}
		data := struct {
			Text       string
			Rank       float64
			Similarity int
			HasResult  bool
		}{
			Text:       "результаты",
			Rank:       rank.Rank(),
			Similarity: similarity,
			HasResult:  true,
		}
		err = tmpl.Execute(w, data)
	}

	if err != nil {
		http.Error(w, "server error2 "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func generateCentrifugoToken(identifier string, channel string) string {
	claims := jwt.MapClaims{
		"sub":      identifier,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"channels": []string{channel},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("my_secret"))
	if err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		return ""
	}

	return signedToken
}
