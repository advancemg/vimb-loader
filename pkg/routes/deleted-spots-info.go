package routes

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/internal/models"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"net/http"
)

// PostGetDeletedSpotInfo godoc
// @Summary Возвращает информацию об удаленных за период спотах
// @Description Возвращает информацию об удаленных за период спотах
// @ID routes-get-deleted-spots-info
// @Tags Споты
// @Param body body models.SwaggerGetDeletedSpotInfoRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.StreamResponse
// @Router /api/v1/spot/deleted-info [post]
func PostGetDeletedSpotInfo(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request models.GetDeletedSpotInfo
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
