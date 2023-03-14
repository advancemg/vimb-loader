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
// @Security ApiKeyAuth
// @ID routes-post-backup
// @Tags Backup
// @Param body body mongodb_backup.SwaggerBackupRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.StreamResponse
// @Failure 401 "Error: Unauthorized"
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

// PostListBackups godoc
// @Summary Список бэкапов MongoDB.
// @Description - Возвращает список бэкапов, s3Key.
// @Security ApiKeyAuth
// @ID routes-post-backup-list
// @Tags Backup
// @Param body body mongodb_backup.SwaggerListBackupsRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.StreamResponse
// @Failure 401 "Error: Unauthorized"
// @Router /api/v1/backup-list [post]
func PostListBackups(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	response, err := mongodb_backup.ListBackups()
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

// PostMongoRestore godoc
// @Summary Восстановление backup MongoDB.
// @Description - Восстановление БД из резервной копии. Указать s3Key бэкапа.
// @Security ApiKeyAuth
// @ID routes-post-restore
// @Tags Backup
// @Param body body mongodb_backup.SwaggerRestoreRequest true  "Запрос"
// @Accept json
// @Produce json
// @Success 200 {object} models.StreamResponse
// @Failure 401 "Error: Unauthorized"
// @Router /api/v1/backup-restore [post]
func PostMongoRestore(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	var request mongodb_backup.SwaggerRestoreRequest
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
	dbConfig := *mongodb_backup.InitConfig()
	err = dbConfig.Restore(request.S3Key)
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
	json.NewEncoder(w).Encode(request)
}
