package routes

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/models"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"net/http"
)

// PostGetMPLans godoc
// @Summary Возвращает доступные медиапланы. Минимальный кванта периода начитки - 1 месяц (с первого по последнее число).
// @Description Возвращает доступные медиапланы. Минимальный кванта периода начитки - 1 месяц (с первого по последнее число).
// @ID routes-get-mediaplans
// @Tags Медиапланы
// @Param body body models.SwaggerGetMPLansRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.StreamResponse
// @Router /api/v1/mediaplan [post]
func PostGetMPLans(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request models.GetMPLans
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
	response, err := request.GetData()
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
