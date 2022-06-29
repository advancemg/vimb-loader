package routes

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/storage/mongodb-backup"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"net/http"
)

// PostMongoBackup godoc
// @Summary Создает backup MongoDB.
// @Description - Создаем резервную копию MongoDВ. При пустых полях берет значения их config файла.
// @ID routes-post-backup
// @Tags Backup
// @Param body body mongodb_backup.SwaggerBackupRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.StreamResponse
// @Router /api/v1/backup [post]
func PostMongoBackup(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request mongodb_backup.Config
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
	if request.Port == "" || request.Host == "" || request.DB == "" || request.Username == "" || request.Password == "" {
		request = *mongodb_backup.InitConfig()
	}
	response, err := request.RunBackup()
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
