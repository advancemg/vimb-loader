package utils

import (
	"encoding/base64"
)

func Base64DecodeString(cert string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(cert)
}
