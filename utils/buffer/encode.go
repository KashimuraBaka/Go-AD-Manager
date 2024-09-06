package buffer

import (
	"encoding/base64"
	"strings"

	"gitee.com/Kashimura/go-baka-control/utils/stringex"
)

func Base64UrlEncode(input []byte) string {
	base64 := base64.StdEncoding.EncodeToString(input)
	base64 = stringex.Strtr(base64, "+/", "-_")
	return strings.ReplaceAll(base64, "=", "")
}
