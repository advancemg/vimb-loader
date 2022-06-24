package routes

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/internal/models"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"net/http"
)

// PostChangeMPlanFilmPlannedInventory godoc
// @Summary Замена планового инвентаря роликов в медиаплане.
// @Description Метод допускает массовое редактирование, при этом все ролики должны относиться к одному медиаплану.
// @ID routes-change-mediaplan-film-planned-inventory
// @Tags Медиапланы
// @Param body body models.SwaggerChangeMPlanFilmPlannedInventoryRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.StreamResponse
// @Router /api/v1/mediaplan/change-film-planned-inventory [post]
func PostChangeMPlanFilmPlannedInventory(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request models.ChangeMPlanFilmPlannedInventory
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		(w).WriteHeader(http.StatusBadRequest)
		var response = utils.FieldValidateErrorType{
			Field:   "id",
			Message: fmt.Sprintf(`Ошибка %s`, err.Error()),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response, err := request.GetDataJson()
	if err != nil {
		(w).WriteHeader(http.StatusBadRequest)
		var response = utils.FieldValidateErrorType{
			Field:   "request",
			Message: fmt.Sprintf(`Ошибка %s`, err.Error()),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	(w).WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
