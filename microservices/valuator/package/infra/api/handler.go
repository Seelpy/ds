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
	textQueryService  query.TextQueryService
}

func NewHandler(textService service.TextService, statisticsQueryService query.StatisticsQueryService, textQueryService query.TextQueryService) *Handler {

	return &Handler{
		textService:       textService,
		statisticsService: statisticsQueryService,
		textQueryService:  textQueryService,
	}
}

func (h *Handler) CreateForm(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("./data/html/base.html", "./data/html/input.html")
	err = tmpl.Execute(w, map[string]interface{}{
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

	_, err := h.textService.Add(text)
	if err != nil {
		http.Error(w, "Failed to process text"+err.Error(), http.StatusInternalServerError)
		return
	}

	h.listImpl(w)
}

func (h *Handler) Statistics(w http.ResponseWriter, r *http.Request) {
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

	summary, err := h.statisticsService.GetSummary(id)
	if err != nil {
		http.Error(w, "Failed to get summary", http.StatusInternalServerError)
		return
	}

	rank := 1 - (float64(summary.SymbolStatistics.AlphabetCount) / float64(summary.SymbolStatistics.AllCount))
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
		TextID:     id,
		Rank:       rank,
		Similarity: similarity,
	}

	tmpl, err := template.ParseFiles("./data/html/base.html", "./data/html/summary.html")
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		http.Error(w, "Failed to get summary", http.StatusInternalServerError)
		return
	}

	err = h.textService.Remove(id)
	if err != nil {
		http.Error(w, "Failed to get summary", http.StatusInternalServerError)
		return
	}

	h.listImpl(w)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	h.listImpl(w)
}

func (h *Handler) listImpl(w http.ResponseWriter) {
	texts, err := h.textQueryService.List()
	if err != nil {
		http.Error(w, "Ошибка при получении списка текстов", http.StatusInternalServerError)
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
