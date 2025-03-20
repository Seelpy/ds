package api

import (
	"errors"
	"github.com/gofrs/uuid"
	"html/template"
	"log"
	"net/http"
	"rankcalculator/package/app/model"
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

	tmpl, err := template.ParseFiles("./data/html/base.html", "./data/html/summary.html")
	if err != nil {
		http.Error(w, "server error1"+err.Error(), http.StatusInternalServerError)
		return
	}

	statistics, err := h.rankRepo.Get(id)

	if errors.Is(err, model.ErrStatisticsNotFound) {
		data := struct {
			Title  string
			TextID uuid.UUID
		}{
			Title:  "Результаты",
			TextID: id,
		}
		err = tmpl.Execute(w, data)
	} else {
		similarity := 0
		if statistics.IsDuplicate {
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
			Rank:       float64(statistics.AllAlphabetCount) / float64(statistics.AllCount),
			Similarity: similarity,
		}
		err = tmpl.Execute(w, data)
	}

	if err != nil {
		http.Error(w, "server error2", http.StatusInternalServerError)
		return
	}
}
