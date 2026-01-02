package stringex

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
	if str == "" {
		return false
	}
	matched, _ := regexp.MatchString(`^[0-9A-Za-z\x{4E00}-\x{9FA5}]+$`, str)
	return !matched
}
