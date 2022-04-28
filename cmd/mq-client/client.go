package main

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/mq-broker"
)

func main() {
	qName := `test-q`
	config := mq_broker.InitConfig()
	err := config.DeclareSimpleQueue(qName)
	if err != nil {
		return
	}
	ch, err := config.Channel()
	if err != nil {
		return
	}
	defer ch.Close()
	err = ch.Qos(1, 0, false)
	messages, err := ch.Consume(qName, "",
		false,
		false,
		false,
		false,
		nil)
	for msg := range messages {
		var bodyJson = map[string]interface{}{}
		err := json.Unmarshal(msg.Body, &bodyJson)
		if err != nil {
			return
		}
		fmt.Printf("Message: %v\n", bodyJson)
		msg.Ack(false)
	}
}
