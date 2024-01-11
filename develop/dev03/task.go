package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/*
	Отсортировать строки (man sort)
	Основное
	Поддержать ключи
		-k — указание колонки для сортировки
		-n — сортировать по числовому значению
		-r — сортировать в обратном порядке
		-u — не выводить повторяющиеся строки
	Дополнительное
	Поддержать ключи
		-M — сортировать по названию месяца
		-b — игнорировать хвостовые пробелы
		-c — проверять отсортированы ли данные
		-h — сортировать по числовому значению с учётом суффиксов
	Программа должна проходить все тесты.
	Код должен проходить проверки go vet и golint.
*/

type Flags struct {
	Input        string
	Output       string
	ColumnSort   int
	NumericSort  bool
	ReverseSort  bool
	UniqueValues bool
}

func fixUpperCase(strs []string) []string {
	for i := 0; i < len(strs)-1; i++ {
		str1 := []rune(strs[i])
		str2 := []rune(strs[i+1])

		if unicode.ToLower(str1[0]) == unicode.ToLower(str2[0]) &&
			unicode.IsUpper(str1[0]) &&
			unicode.IsLower(str2[0]) {
			temp := strs[i]
			strs[i] = strs[i+1]
			strs[i+1] = temp
		}
	}

	return strs
}

func quickSort(strs []string, start, end int, numericSort bool) []string {
	if start < end {
		pivot := strs[start]
		left := start
		right := end

		for left < right {
			if !numericSort {
				for left < right && strings.ToLower(strs[right]) >= strings.ToLower(pivot) {
					right--
				}
			} else {
				p, err := strconv.Atoi(pivot)
				if err != nil {
					log.Fatalf("not a number: %s", pivot)
				}
				for left < right {
					r, err := strconv.Atoi(strs[right])
					if err != nil {
						log.Fatalf("not a number: %s", strs[right])
					}
					if r >= p {
						right--
					} else {
						break
					}
				}
			}

			if left < right {
				strs[left] = strs[right]
				left++
			}

			if !numericSort {
				for left < right && strings.ToLower(strs[left]) <= strings.ToLower(pivot) {
					left++
				}
			} else {
				p, err := strconv.Atoi(pivot)
				if err != nil {
					log.Fatalf("not a number: %s", pivot)
				}
				for left < right {
					l, err := strconv.Atoi(strs[left])
					if err != nil {
						log.Fatalf("not a number: %s", strs[left])
					}
					if l <= p {
						left++
					} else {
						break
					}
				}
			}

			if left < right {
				strs[right] = strs[left]
				right--
			}
		}

		strs[left] = pivot

		strs = quickSort(strs, start, left-1, numericSort)
		strs = quickSort(strs, left+1, end, numericSort)
	}

	sortedStrs := make([]string, len(strs))
	copy(sortedStrs, strs)
	return sortedStrs
}

func onlyUnique(data []string) []string {
	res := make([]string, 0, len(data))
	m := make(map[string]struct{})
	for _, str := range data {
		if _, ok := m[str]; !ok {
			m[str] = struct{}{}
			res = append(res, str)
		}
	}
	return res
}

func reversed(data []string) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func getKeysOfMap(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func sortByColumn(strs []string, flg Flags) []string {
	srcMap := make(map[string]string)
	for _, str := range strs {
		columns := strings.Split(str, " ")
		if len(columns) > flg.ColumnSort {
			srcMap[columns[flg.ColumnSort]] = str
			continue
		}
		srcMap[columns[0]] = str
	}
	keysToBeSorted := getKeysOfMap(srcMap)
	quickSort(keysToBeSorted, 0, len(keysToBeSorted)-1, flg.NumericSort)
	sortedStrs := make([]string, len(strs))
	for i, key := range keysToBeSorted {
		sortedStrs[i] = srcMap[key]
	}
	return sortedStrs
}

func Sort(flg Flags) {
	strs := readFile(flg.Input)

	result := make([]string, len(strs))
	copy(result, strs)

	if flg.ColumnSort >= 0 {
		result = sortByColumn(result, flg)
		result = fixUpperCase(result)
	} else {
		result = quickSort(result, 0, len(result)-1, false)
		result = fixUpperCase(result)
	}

	if flg.UniqueValues {
		result = quickSort(result, 0, len(result)-1, false)
		result = fixUpperCase(result)
		result = onlyUnique(result)
	}

	if flg.ReverseSort {
		reversed(result)
	}

	writeToFile(flg.Output, result)
}

func readFile(fileName string) []string {
	var rows []string

	file, err := os.Open(fileName)
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatalf("failed to close file: %v", err)
		}
	}(file)

	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		rows = append(rows, sc.Text())
	}

	return rows
}

func writeToFile(fileName string, strs []string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatalf("failed to close file: %v", err)
		}
	}(file)

	for _, str := range strs {
		if _, err := fmt.Fprintln(file, str); err != nil {
			log.Fatalf("failed to write to file: %v", err)
		}
	}
}

func parseFlags() Flags {
	in := flag.String("in", "", "input file path")
	out := flag.String("out", "", "output file path")
	columnSort := flag.Int("k", -1, "column number for sorting")
	numericSort := flag.Bool("n", false, "sort numerically")
	reverseSort := flag.Bool("r", false, "sort in reverse order")
	uniqueValues := flag.Bool("u", false, "show only unique values")

	flag.Parse()

	flg := Flags{
		Input:        *in,
		Output:       *out,
		ColumnSort:   *columnSort,
		NumericSort:  *numericSort,
		ReverseSort:  *reverseSort,
		UniqueValues: *uniqueValues,
	}

	if *in == "" || *out == "" {
		log.Fatal("required: input and output files")
	}

	return flg
}

func main() {
	flg := parseFlags()
	Sort(flg)
}
