package models

import (
	"encoding/json"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
)

const BudgetTable = "budgets"

type Budget struct {
	AgreementId int
}

type BudgetsUpdateRequest struct {
	Body []byte
}

func BudgetStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := BudgetsUpdateQueue
		amqpConfig := mq_broker.InitConfig()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			errorCh <- err
		}
		ch, err := amqpConfig.Channel()
		if err != nil {
			errorCh <- err
		}
		err = ch.Qos(1, 0, false)
		messages, err := ch.Consume(qName, "",
			false,
			false,
			false,
			false,
			nil)
		for msg := range messages {
			var bodyJson MqUpdateMessage
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			/*read from s3 by s3Key*/
			req := BudgetsUpdateRequest{
				Body: []byte{},
			}
			err = req.Update()
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}

		defer close(errorCh)
	}()
	return errorCh
}

func (request *BudgetsUpdateRequest) Update() error {

	return nil
}
