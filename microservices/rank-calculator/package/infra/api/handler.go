package api

import (
	"encoding/json"
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

const (
	userIDKey      = "user_id_key"
	userCountryKey = "user_country_key"
)

func NewHandler(rankRepo model.ReadOnlyTextStatisticsRepository, secret string) *Handler {
	return &Handler{
		rankRepo: rankRepo,
		secret:   secret,
	}
}

type Handler struct {
	rankRepo model.ReadOnlyTextStatisticsRepository
	secret   string
}

func (h *Handler) Statistics(w http.ResponseWriter, r *http.Request) {
	log.Println("HAaSDSADE1")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := parseContext(w, r)
	log.Println("HAaSDSADE2", ok)
	if !ok {
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

	rank, err := h.rankRepo.Get(userID, id)
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

func parseContext(w http.ResponseWriter, r *http.Request) (userID uuid.UUID, ok bool) {
	ctx := r.Context()
	userID, err := uuid.FromString(ctx.Value(userIDKey).(string))
	log.Println("HAaSDSADE3", err)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	return userID, true
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

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
