package main

import (
	"github.com/advancemg/vimb-loader/pkg/mq-broker"
	"time"
)

func main() {
	qName := `test-q`
	config := mq_broker.InitConfig()
	err := config.DeclareSimpleQueue(qName)
	if err != nil {
		return
	}
	increment := 0
	for {
		m := map[string]interface{}{
			"increment": increment,
		}
		err := config.PublishJson(qName, m)
		if err != nil {
			break
		}
		time.Sleep(1 * time.Second)
		increment++
	}
}
