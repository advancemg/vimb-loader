package models

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMediaplanLoadBadgerRequest_LoadMediaplan(t *testing.T) {
	request := MediaplanLoadBadgerRequest{Month: 201902}
	loadBudgets, err := request.LoadMediaplan()
	if err != nil {
		fmt.Println(err.Error())
	}
	indent, _ := json.MarshalIndent(loadBudgets, "", " ")
	fmt.Println(string(indent))
}
