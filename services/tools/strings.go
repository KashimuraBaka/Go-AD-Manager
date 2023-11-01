package tools

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
