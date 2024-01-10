package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Утилита sort ===
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
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Flags struct {
	ColumnSort   int
	NumericSort  bool
	ReverseSort  bool
	UniqueValues bool
}

func fixUpperCaseOrder(strs []string) {
	for i := 0; i < len(strs)-1; i++ {
		str1 := []rune(strs[i])
		str2 := []rune(strs[i+1])

		if unicode.ToLower(str1[0]) == unicode.ToLower(str2[0]) &&
			unicode.IsUpper(str1[0]) &&
			unicode.IsLower(str2[0]) {
			buf := strs[i]
			strs[i] = strs[i+1]
			strs[i+1] = buf
		}
	}
}

func stringQuickSort(strs []string, start, end int, numericSort bool) {
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
				r, err := strconv.Atoi(strs[right])
				if err != nil {
					log.Fatalf("not number: %s", strs[right])
				}
				b, err := strconv.Atoi(pivot)
				if err != nil {
					log.Fatalf("not number: %s", strs[right])
				}
				for left < right && r >= b {

					right--

					r, err = strconv.Atoi(strs[right])
					if err != nil {
						log.Fatalf("not number: %s", strs[right])
					}
					b, err = strconv.Atoi(pivot)
					if err != nil {
						log.Fatalf("not number: %s", strs[right])
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
				l, err := strconv.Atoi(strs[left])
				if err != nil {
					log.Fatalf("not number: %s", strs[right])
				}
				b, err := strconv.Atoi(pivot)
				if err != nil {
					log.Fatalf("not number: %s", strs[right])
				}
				for left < right && l <= b {
					left++
					l, err = strconv.Atoi(strs[left])
					if err != nil {
						log.Fatalf("not number: %s", strs[right])
					}
					b, err = strconv.Atoi(pivot)
					if err != nil {
						log.Fatalf("not number: %s", strs[right])
					}
				}
			}

			if left < right {
				strs[right] = strs[left]
				right--
			}
		}

		strs[left] = pivot

		stringQuickSort(strs, start, left-1, numericSort)
		stringQuickSort(strs, left+1, end, numericSort)
	}
}

func onlyUnique(data []string) []string {
	res := make([]string, 0, len(data))
	m := make(map[string]bool)
	for _, str := range data {
		if _, ok := m[str]; !ok {
			m[str] = true
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

func sortByColumn(strs []string, flg Flags) {
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
	stringQuickSort(keysToBeSorted, 0, len(keysToBeSorted)-1, flg.NumericSort)
	for i, key := range keysToBeSorted {
		strs[i] = srcMap[key]
	}
}

func Sort(inFilePath, outFilePath string, flg Flags) {
	strs := readFile(inFilePath)

	result := make([]string, len(strs))
	copy(result, strs)

	if flg.ColumnSort >= 0 {
		sortByColumn(result, flg)
		fixUpperCaseOrder(result)
	} else {
		stringQuickSort(result, 0, len(result)-1, false)
		fixUpperCaseOrder(result)
	}

	if flg.UniqueValues {
		stringQuickSort(result, 0, len(result)-1, false)
		fixUpperCaseOrder(result)
		result = onlyUnique(result)
	}

	if flg.ReverseSort {
		reversed(result)
	}

	writeToFile(outFilePath, result)
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

func main() {
	flg := Flags{
		ColumnSort:   2,
		NumericSort:  false,
		ReverseSort:  false,
		UniqueValues: false,
	}

	Sort("text.txt", "sorted.txt", flg)
}
