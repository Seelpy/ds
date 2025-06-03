package api

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"html/template"
	"net/http"
	"valuator/package/app/authorization"
	"valuator/package/app/query"
	"valuator/package/app/service"
)

const (
	userIDKey      = "user_id_key"
	userCountryKey = "user_country_key"
)

type Handler struct {
	textService      service.TextService
	textQueryService query.TextQueryService
	secret           string
}

func NewHandler(textService service.TextService, textQueryService query.TextQueryService, secret string) *Handler {
	return &Handler{
		textService:      textService,
		textQueryService: textQueryService,
		secret:           secret,
	}
}

func (h *Handler) CreateForm(w http.ResponseWriter, r *http.Request) {
	_, ok := parseContext(w, r)
	if !ok {
		return
	}
	tmpl, err := template.ParseFiles("./data/html/base.html", "./data/html/input.html")
	err = tmpl.Execute(w, map[string]interface{}{
		"Title": "Главная",
	})
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("./data/html/base.html", "./data/html/login.html")
	err = tmpl.Execute(w, map[string]interface{}{})
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ProcessText(w http.ResponseWriter, r *http.Request) {
	actx, ok := parseContext(w, r)
	if !ok {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")

	_, err := h.textService.Add(actx, text)
	if err != nil {
		http.Error(w, "Failed to process text"+err.Error(), http.StatusInternalServerError)
		return
	}

	h.listImpl(actx, w)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	actx, ok := parseContext(w, r)
	if !ok {
		return
	}

	idStr := r.FormValue("id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		http.Error(w, "Failed to get summary", http.StatusInternalServerError)
		return
	}

	err = h.textService.Remove(actx, id)
	if err != nil {
		http.Error(w, "Failed to get summary", http.StatusInternalServerError)
		return
	}

	h.listImpl(actx, w)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx, ok := parseContext(w, r)
	if ok {
		h.listImpl(ctx, w)
	}
}

func (h *Handler) listImpl(ctx authorization.Context, w http.ResponseWriter) {
	texts, err := h.textQueryService.List(ctx)
	if err != nil {
		http.Error(w, "Ошибка при получении списка текстов:"+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./data/html/list.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct {
		Texts []query.TextData
	}{
		Texts: texts,
	})
	if err != nil {
		http.Error(w, "Ошибка при отображении шаблона", http.StatusInternalServerError)
		return
	}
}

func parseContext(w http.ResponseWriter, r *http.Request) (actx authorization.Context, ok bool) {
	ctx := r.Context()
	userID, err := uuid.FromString(ctx.Value(userIDKey).(string))
	fmt.Println("ASDASD4", err)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	country, ok := ctx.Value(userCountryKey).(string)
	fmt.Println("ASDASD5", ok)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	return authorization.NewContext(ctx, userID, country), ok
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// TODO: err
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
