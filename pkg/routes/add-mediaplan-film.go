package routes

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/internal/models"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"net/http"
)

// PostAddMPlanFilm godoc
// @Summary “Втягивает” ролик в медиаплан. Период ролика в медиаплане устанавливается равным периоду медиаплана.
// @Description “Втягивает” ролик в медиаплан. Период ролика в медиаплане устанавливается равным периоду медиаплана.
// @ID routes-add-mediaplan-film
// @Tags Медиапланы
// @Param body body models.SwaggerAddMPlanFilmRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.StreamResponse
// @Router /api/v1/mediaplan/film [post]
func PostAddMPlanFilm(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request models.AddMPlanFilm
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
