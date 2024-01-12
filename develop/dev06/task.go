package main

/*
	Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные
	Поддержать флаги:
	-f - "fields" - выбрать поля (колонки)
	-d - "delimiter" - использовать другой разделитель
	-s - "separated" - только строки с разделителем
	Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Flags struct {
	Delimiter string // -d
	Fields    string // -f
	Separated bool   // -s
}

func Cut(delimiter string, columns []int, separatedOnly bool, input io.Reader) {
	// Сканер для построчного чтения
	scanner := bufio.NewScanner(input)
	// Цикл сканирования строк
	for scanner.Scan() {
		line := scanner.Text()

		// Проверка флага separatedOnly и пропуск строки, если она не содержит разделителя и флаг установлен
		if separatedOnly && !strings.Contains(line, delimiter) {
			continue
		}

		// Разбиение строки на поля с использованием указанного разделителя
		fields := strings.Split(line, delimiter)
		output := make([]string, len(columns))

		// Перебор указанных столбцов и добавление соответствующих значений в выходной массив
		for i, col := range columns {
			if col >= 1 && col <= len(fields) {
				output[i] = fields[col-1]
			}
		}

		// Вывод результата с введенным разделителем
		fmt.Println(strings.Join(output, delimiter))
	}
}

// Парсит аргументы коммандной строки
func parseFlags() Flags {
	delimiter := flag.String("d", "\t", "field delimiter")
	fields := flag.String("f", "", "selected columns")
	separated := flag.Bool("s", false, "only lines with delimiter")

	flag.Parse()

	flg := Flags{
		Delimiter: *delimiter,
		Fields:    *fields,
		Separated: *separated,
	}

	return flg
}

// Преобразует строку в целое число
func parseInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return num
}

// Создаем слайс для указанных столбцов
func calculateFields(flg *Flags) []int {
	fields := []int{}
	fieldsInput := strings.Split(flg.Fields, ",")
	for _, fieldStr := range fieldsInput {
		field := strings.TrimSpace(fieldStr)
		if field != "" {
			fields = append(fields, parseInt(field))
		}
	}

	return fields
}

func main() {
	flg := parseFlags()
	fields := calculateFields(&flg)

	Cut(flg.Delimiter, fields, flg.Separated, os.Stdin)
}
