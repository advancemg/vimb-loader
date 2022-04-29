package models

import (
	goConvert "github.com/advancemg/go-convert"
)

type SwaggerChangeFilmsRequest struct {
	ChangeFilms []struct {
		FakeSpotIDs  string `json:"FakeSpotIDs"`
		CommInMplIDs string `json:"CommInMplIDs"`
	} `json:"ChangeFilms"`
}

type ChangeFilms struct {
	goConvert.UnsortedMap
}
