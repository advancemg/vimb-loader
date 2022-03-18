package utils

import (
	"fmt"
	"os"
)

type FieldValidateErrorType struct {
	Field   string `json:"field" bson:"field,omitempty"`
	Code    int    `json:"code" bson:"code,omitempty"`
	Message string `json:"message" bson:"message,omitempty"`
}

func CheckErr(err error) {
	if err != nil {
		fmt.Printf(`Api fatal error. Stop service. Err:[%v]`, err)
	}
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == `` {
		return fallback
	}
	return value
}
