package routes

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/internal/models"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"net/http"
)

// PostGetAdvMessages godoc
// @Summary Возвращает список роликов, созданных в указанный период времени. Список отфильтрован по правам НВ.
// @Description Возвращает список роликов, созданных в указанный период времени. Список отфильтрован по правам НВ.
// @ID routes-get-adv-messages
// @Tags Справочники
// @Param body body models.SwaggerGetAdvMessagesRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.StreamResponse
// @Router /api/v1/adv-messages [post]
func PostGetAdvMessages(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request models.GetAdvMessages
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

// PostLoadAdvMessages godoc
// @Summary Создание задач, на загрузку роликов.
// @Description Создание задач, на загрузку роликов, за выбранный период.
// @ID routes-load-adv-messages
// @Tags Справочники
// @Param body body models.AdvMessagesLoadRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.CommonResponse
// @Router /api/v1/adv-messages/load [post]
func PostLoadAdvMessages(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request models.AdvMessagesLoadRequest
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

// PostAdvMessagesQuery godoc
// @Summary Загрузку сохраненных роликов.
// @Description Динамический запрос на загрузку сохраненных данных. Логические операторы: eq ==, ne !=, gt >, lt <, ge >=, le <=, in in, isnil is nil.
// @ID routes-query-adv-messages
// @Tags Справочники
// @Param body body models.AdvMessageQuery true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.CommonResponse
// @Router /api/v1/adv-messages/query [post]
func PostAdvMessagesQuery(w http.ResponseWriter, r *http.Request) {
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
	response, err := request.QueryAdvMessages()
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
