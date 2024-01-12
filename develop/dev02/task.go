package main

import (
	"errors"
	"strconv"
	"strings"
)

/*
	Создать Go-функцию, осуществляющую примитивную распаковку строки,
	содержащую повторяющиеся символы/руны, например:
	"a4bc2d5e" => "aaaabccddddde"
	"abcd" => "abcd"
	"45" => "" (некорректная строка)
	"" => ""

	Дополнительно
	Реализовать поддержку escape-последовательностей.
	Например:
	qwe\4\5 => qwe45 (*)
	qwe\45 => qwe44444 (*)
	qwe\\5 => qwe\\\\\ (*)


	В случае если была передана некорректная строка,
	функция должна возвращать ошибку. Написать unit-тесты.

*/

// Глобальная переменная ошибки invalid string
var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	// Создание слайса рун из строки
	runes := []rune(str)

	// strings.Builder для формирования строки(поэлементного)
	builder := strings.Builder{}

	// Цикл по каждому индексу в слайсе
	for i := 0; i < len(runes); {
		currentRune := runes[i]
		// Если текущий символ буква, то добавляем его к результирующей строке
		if isLetter(currentRune) {
			builder.WriteString(string(currentRune))
			i++
		}
		// Если текущий символ цифра
		if isDigit(currentRune) {
			count := 0

			j := i
			prevIdx := i

			// Цикл для извлечения числа и определения повторений символа
			for j < len(runes) && isDigit(runes[j]) {
				// Преобразование руны в целое число
				temp, _ := strconv.Atoi(string(runes[j]))
				count = count*10 + temp
				j++
				i++
			}
			// Добавление в результирующую строку повторяющегося символа
			if prevIdx > 0 {
				writeString(runes[prevIdx-1], &builder, count-1)
				continue
			}

			// В случае, если перед цифрой нет символа для повторения
			return "", ErrInvalidString
		}
		// Если текущий символ слэш, то добавляем следующий символ к результирующей строке
		if string(currentRune) == `\` {
			if i < len(runes)-1 {
				builder.WriteString(string(runes[i+1]))
				i += 2
				continue
			}

			// В случае, если обратный слэш находится в конце строки или не имеет следующего символа
			return "", ErrInvalidString
		}
	}

	// Возвращаем полученную результирующую строку, и nil
	return builder.String(), nil
}

// Функция для проверки, является ли руна цифрой
func isDigit(r rune) bool {
	if r >= '0' && r <= '9' {
		return true
	}

	return false
}

// Функция для проверки, является ли руна буквой
func isLetter(r rune) bool {
	if !isDigit(r) && string(r) != `\` {
		return true
	}

	return false
}

// Функция для добавления символа r в результирующую строку builder count раз
func writeString(r rune, builder *strings.Builder, count int) {
	for i := 0; i < count; i++ {
		builder.WriteString(string(r))
	}
}
