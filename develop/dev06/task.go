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
	"strings"
)

type Flags struct {
	Delimiter string
	Fields    string
	Separated bool
}

func Cut(delimiter string, columns []int, separatedOnly bool, input io.Reader) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		if separatedOnly && !strings.Contains(line, delimiter) {
			continue
		}

		fields := strings.Split(line, delimiter)
		output := make([]string, len(columns))

		for i, col := range columns {
			if col >= 1 && col <= len(fields) {
				output[i] = fields[col-1]
			}
		}

		fmt.Println(strings.Join(output, delimiter))
	}
}

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

func parseInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func main() {
	flg := parseFlags()

	columns := []int{}
	columnsInput := strings.Split(flg.Fields, ",")
	for _, colStr := range columnsInput {
		col := strings.TrimSpace(colStr)
		if col != "" {
			columns = append(columns, parseInt(col))
		}
	}

	Cut(flg.Delimiter, columns, flg.Separated, os.Stdin)
}
