package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"authentication/package/app/model"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	userRepository model.UserRepository
	secret         string
}

func NewHandler(userRepo model.UserRepository, secret string) *Handler {
	return &Handler{
		userRepository: userRepo,
		secret:         secret,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	type loginInput struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	var input loginInput
	if err := decodeJSONBody(w, r, &input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.userRepository.GetByLogin(input.Login)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	passwordHash := h.hash(input.Password)
	if user.PasswordHash != passwordHash {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	payload := jwt.MapClaims{
		"sub": user.UserID,
		"c":   user.Country,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(h.secret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": tokenString,
	})
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization header")
		return
	}

	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(h.secret), nil
	})

	if err != nil || !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	userID, err := uuid.FromString(claims["sub"].(string))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID in token")
		return
	}

	user, err := h.userRepository.GetByID(userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	newPayload := jwt.MapClaims{
		"sub": user.UserID,
		"c":   user.Country,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newPayload)
	newTokenString, err := newToken.SignedString([]byte(h.secret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    newTokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": newTokenString,
	})
}

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	type registrationRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Country  string `json:"country"`
	}

	var request registrationRequest
	if err := decodeJSONBody(w, r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	input := model.UserRegisterInput{
		Login:        request.Login,
		PasswordHash: h.hash(request.Password),
		Country:      model.Country(request.Country),
	}

	user, err := h.userRepository.Register(input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error registering user")
		return
	}

	newPayload := jwt.MapClaims{
		"sub": user.UserID,
		"c":   user.Country,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newPayload)
	newTokenString, err := newToken.SignedString([]byte(h.secret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    newTokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": newTokenString,
	})
}

func (h *Handler) hash(password string) string {
	return "hash_" + password
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
