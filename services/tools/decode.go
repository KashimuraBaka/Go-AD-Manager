package tools

import (
	"encoding/base64"
	"strings"
)

func Base64Decode(str string) []byte {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return []byte(str)
	}
	return data
}

func Base64UrlDecode(input string) []byte {
	remainder := len(input) % 4
	if remainder != 0 {
		addlen := 4 - remainder
		input += strings.Repeat("=", addlen)
	}
	input = Strtr(input, "-_", "+/")
	return Base64Decode(input)
}
