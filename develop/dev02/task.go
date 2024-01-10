package main

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	runes := []rune(str)

	builder := strings.Builder{}

	idx := 0
	for idx < len(runes) {
		currentRune := runes[idx]
		if isLetter(currentRune) {
			builder.WriteString(string(currentRune))
			idx++
		}
		if isDigit(currentRune) {
			cnt := 0

			j := idx
			prevIdx := idx
			for j < len(runes) && isDigit(runes[j]) {
				buf, _ := strconv.Atoi(string(runes[j]))
				cnt = cnt*10 + buf
				j++
				idx++
			}
			if prevIdx > 0 {
				writeString(runes[prevIdx-1], &builder, cnt-1)
				continue
			}
			return "", ErrInvalidString
		}
		if string(currentRune) == `\` {
			if idx < len(runes)-1 {
				builder.WriteString(string(runes[idx+1]))
				idx += 2
				continue
			}
			return "", ErrInvalidString
		}
	}

	return builder.String(), nil
}

func isDigit(r rune) bool {
	if r >= '0' && r <= '9' {
		return true
	}
	return false
}

func isLetter(r rune) bool {
	if !isDigit(r) && string(r) != `\` {
		return true
	}
	return false
}

func writeString(r rune, builder *strings.Builder, count int) {
	for i := 0; i < count; i++ {
		builder.WriteString(string(r))
	}
}
