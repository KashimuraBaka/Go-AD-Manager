package tools

import (
	"encoding/base64"
	"strings"
)

func Base64UrlEncode(input []byte) string {
	base64 := base64.StdEncoding.EncodeToString(input)
	base64 = Strtr(base64, "+/", "-_")
	return strings.ReplaceAll(base64, "=", "")
}
