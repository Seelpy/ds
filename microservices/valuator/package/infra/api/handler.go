package api

import (
	"github.com/gofrs/uuid"
	"html/template"
	"net/http"
	"valuator/package/app/query"
	"valuator/package/app/service"
)

type Handler struct {
	textService       service.TextService
	statisticsService query.StatisticsQueryService
	templates         *template.Template
}

func NewHandler(ts service.TextService, ss query.StatisticsQueryService) *Handler {
	templates := template.Must(template.ParseFiles("./data/html/index.html", "./data/html/input.html", "./data/html/summary.html", "./data/html/baseSummary.html"))
	return &Handler{
		textService:       ts,
		statisticsService: ss,
		templates:         templates,
	}
}

func (h *Handler) Index(w http.ResponseWriter, _ *http.Request) {
	err := h.templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Title": "Главная",
	})
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ProcessText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")

	textID, err := h.textService.Add(text)
	if err != nil {
		http.Error(w, "Failed to process text"+err.Error(), http.StatusInternalServerError)
		return
	}

	summary, err := h.statisticsService.GetSummary(textID)
	if err != nil {
		http.Error(w, "Failed to get summary", http.StatusInternalServerError)
		return
	}

	rank := float64(summary.SymbolStatistics.AlphabetCount) / float64(summary.SymbolStatistics.AllCount)
	similarity := 0
	if summary.UniqueStatistics.IsDuplicate {
		similarity = 1
	}

	data := struct {
		Title      string
		TextID     uuid.UUID
		Rank       float64
		Similarity int
	}{
		Title:      "Результаты",
		TextID:     textID,
		Rank:       rank,
		Similarity: similarity,
	}

	err = h.templates.ExecuteTemplate(w, "baseSummary.html", data)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}
