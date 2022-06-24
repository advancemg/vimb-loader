package routes

import (
	"encoding/json"
	"github.com/advancemg/vimb-loader/internal/models"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"net/http"
)

type MqInfo struct {
	Name      string `json:"name"`
	Messages  int    `json:"messages"`
	Consumers int    `json:"consumers"`
}

// GetQueuesMetrics godoc
// @Summary Состояние очередей.
// @Description Состояние очередей.
// @ID routes-get-mq-queue-metrics
// @Tags MQ
// @Accept json
// @Produce json
// @Success 200 {object} []routes.MqInfo
// @Router /api/v1/mq/queues [get]
func GetQueuesMetrics(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		(w).WriteHeader(http.StatusOK)
		return
	}
	config := mq_broker.InitConfig()
	var response []MqInfo
	for _, qName := range models.QueueNames {
		info, _ := config.GetQueueInfo(qName)
		response = append(response, MqInfo{
			Name:      qName,
			Messages:  info.Messages,
			Consumers: info.Consumers,
		})
	}
	(w).WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
