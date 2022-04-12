package routes

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/models"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"net/http"
)

// PostGetChannels godoc
// @Summary Возвращает отфильтрованный по правам НВ список каналов (точек врезки) указанного направления продаж.
// @Description Результат включает активные в настоящий момент для размещения каналы, а также каналы, которые станут активными в течение ближайших трех месяцев.
// @ID routes-get-channels
// @Tags Справочники
// @Param body body models.SwaggerGetChannelsRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.StreamResponse
// @Router /api/v1/channels [post]
func PostGetChannels(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request models.GetChannels
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

// PostLoadChannels godoc
// @Summary Создание задач, на загрузку каналов.
// @Description Создание задач, на загрузку каналов, за выбранный период.
// @ID routes-load-channels
// @Tags Справочники
// @Param body body models.ChannelLoadRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.CommonResponse
// @Router /api/v1/channels/load [post]
func PostLoadChannels(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request models.ChannelLoadRequest
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

// PostLoadBadgerChannels godoc
// @Summary Загрузку сохраненных каналов.
// @Description Загрузку сохраненных каналов, по ID направлению продаж.
// @ID routes-load-badger-channels
// @Tags Справочники
// @Param body body models.ChannelLoadRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.CommonResponse
// @Router /api/v1/channels/badger/load [post]
func PostLoadBadgerChannels(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request models.ChannelLoadRequest
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
	response, err := request.LoadChannels()
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
