package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	go_convert "github.com/advancemg/go-convert"
	"io"
	"os"
)

type VimbResponse struct {
	FilePath string
}

func (req *VimbResponse) Convert(key string) (interface{}, error) {
	file, err := os.Open(req.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	zipBuffer := new(bytes.Buffer)
	_, err = io.Copy(zipBuffer, file)
	if err != nil {
		return nil, fmt.Errorf("not copy zip data, %w", err)
	}
	toJson, err := go_convert.ZipXmlToJson(zipBuffer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("not ZipXmlToJson, %w", err)
	}
	var jsonMap map[string]interface{}
	err = json.Unmarshal(toJson, &jsonMap)
	if err != nil {
		return nil, err
	}
	var recursiveGetField func(js map[string]interface{}) interface{}
	recursiveGetField = func(js map[string]interface{}) interface{} {
		for k, v := range js {
			if k == "attributes" {
				continue
			}
			if _, ok := v.([]interface{}); ok {
				continue
			}
			if k == key {
				return v
			}
			return recursiveGetField(v.(map[string]interface{}))
		}
		return nil
	}
	return recursiveGetField(jsonMap), nil
}
