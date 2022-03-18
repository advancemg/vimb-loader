package routes

import (
	"encoding/json"
	"net/http"
)

const (
	service = "vimb-loader"
)

type Status struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func setupResponse(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, token, refresh-token")
}

// Health godoc
// @Summary Проверка статуса сервиса
// @Description Проверка статуса сервиса
// @ID routes-health
// @Tags Проверка сервиса
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1 [get]
func Health(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Status{Status: "ok", Message: "Api server run!", Code: http.StatusOK})
}
