package utils

import "encoding/base64"

func DecodeBase64(b64 string) (string, error) {
	result, err := base64.StdEncoding.DecodeString(b64)
	return string(result), err
}
