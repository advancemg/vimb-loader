package models

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestBudgetsBadgerQuery_FindJson(t *testing.T) {
	request := BudgetLoadBadgerRequest{Month: 201902}
	loadBudgets, err := request.LoadBudgets()
	if err != nil {
		fmt.Println(err.Error())
	}
	indent, _ := json.MarshalIndent(loadBudgets, "", " ")
	fmt.Println(string(indent))
}
