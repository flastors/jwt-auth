package utils

import (
	"encoding/base64"
)

func EncodeB64(str string) string {
	base64Str := base64.StdEncoding.EncodeToString([]byte(str))
	return string(base64Str)
}

func DecodeB64(base64Str string) (string, error) {
	str, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	return string(str), nil
}
