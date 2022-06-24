package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func FileToBase64(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Base64Encode:%w", err)
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Base64Encode:%w", err)
	}
	return []byte(base64.StdEncoding.EncodeToString(bytes)), nil
}
