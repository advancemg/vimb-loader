package routes

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/internal/models"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"net/http"
)

// PostGetRanks godoc
// @Summary Возвращает справочник рангов размещения.
// @Description Возвращает справочник рангов размещения.
// @Security ApiKeyAuth
// @ID routes-get-ranks
// @Tags Справочники
// @Param body body models.SwaggerGetRanksRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.StreamResponse
// @Failure 401 "Error: Unauthorized"
// @Router /api/v1/ranks [post]
func PostGetRanks(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request models.GetRanks
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

// PostLoadRanks godoc
// @Summary Создание задач, на загрузку рангов размещения.
// @Description Создание задач, на загрузку рангов размещения.
// @Security ApiKeyAuth
// @ID routes-load-ranks
// @Tags Справочники
// @Param body body models.RanksLoadRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.CommonResponse
// @Failure 401 "Error: Unauthorized"
// @Router /api/v1/ranks/load [post]
func PostLoadRanks(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request models.RanksLoadRequest
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
	response, err := request.InitTasks()
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

// PostRanksQuery godoc
// @Summary Загрузку сохраненных рангов размещения.
// @Description Динамический запрос на загрузку сохраненных данных. Логические операторы: eq ==, ne !=, gt >, lt <, ge >=, le <=, in in, isnil is nil.
// @Security ApiKeyAuth
// @ID routes-query-ranks
// @Tags Справочники
// @Param body body models.RankQuery true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.CommonResponse
// @Failure 401 "Error: Unauthorized"
// @Router /api/v1/ranks/query [post]
func PostRanksQuery(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request models.Any
	err := json.NewDecoder(r.Body).Decode(&request.Body)
	if err != nil {
		(w).WriteHeader(http.StatusBadRequest)
		var response = utils.FieldValidateErrorType{
			Field:   "id",
			Message: fmt.Sprintf(`Ошибка %s`, err.Error()),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response, err := request.QueryRanks()
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
