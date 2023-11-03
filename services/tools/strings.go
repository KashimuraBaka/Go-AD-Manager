package tools

import (
	"regexp"
)

func Strtr(str string, from string, to string) string {
	stream := []rune(str)
	for _, char := range from {
		for i, r := range stream {
			if r == char {
				stream[i] = char
			}
		}
	}
	return string(stream)
}

func HasSysbol(str string) bool {
	matched, err := regexp.MatchString("^[0-9A-Za-z\u4e00-\u9fa5]+$", str)
	if err != nil {
		return true
	}
	return matched
}
