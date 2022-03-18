package models

type StreamResponse struct {
	Body    []byte `json:"body"`
	Request string `json:"request"`
}
