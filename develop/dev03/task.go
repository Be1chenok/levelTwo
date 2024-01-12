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
	Input        string // -in входной файл
	Output       string // -out выходной файл
	ColumnSort   int    // -k
	NumericSort  bool   // -n
	ReverseSort  bool   // -r
	UniqueValues bool   // -u
}

/*
Переупорядочивает строки в массиве strs таким образом, чтобы строки,
начинающиеся с заглавной буквы, следовали после строк, начинающихся с
прописной буквы.
*/
func fixUpperCase(strs []string) []string {
	// Обход по всем строкам слайса strs
	for i := 0; i < len(strs)-1; i++ {
		// Текущая строка
		currentStr := []rune(strs[i])
		// Следующая строка
		nextStr := []rune(strs[i+1])

		/*
			Если текущая строка strs[i] начинается с заглавной буквы, а следующая строка strs[i+1]
			начинается с прописной буквы, то строки меняются местами.
		*/
		if unicode.ToLower(currentStr[0]) == unicode.ToLower(nextStr[0]) &&
			unicode.IsUpper(currentStr[0]) &&
			unicode.IsLower(nextStr[0]) {
			temp := strs[i]
			strs[i] = strs[i+1]
			strs[i+1] = temp
		}
	}

	return strs
}

/*
Сортировка слайса строк с использованием алгоритма быстрой сортировки
start - индекс начала сортировки
end - индекс конца сортировки
*/
func quickSort(strs []string, start, end int, numericSort bool) []string {
	if start < end {
		// Выбор опорного элемента
		pivot := strs[start]
		left := start
		right := end

		for left < right {
			// Если не числовая сортировка
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

			// Если не числовая сортировка
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

		// Помещение опорного элемента на правильное место
		strs[left] = pivot

		// Рекурсивное применение quickSort к двум частям массива
		strs = quickSort(strs, start, left-1, numericSort)
		strs = quickSort(strs, left+1, end, numericSort)
	}

	// Создание нового массива и копирование отсортированных строк
	sortedStrs := make([]string, len(strs))
	copy(sortedStrs, strs)
	return sortedStrs
}

// Удаляем дублиуаты из слайса строк
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

// Инвертирует порядок элементов в слайсе
func reversed(data []string) []string {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	return data
}

// Получение ключей из карты в виде массива строк.
func getKeysOfMap(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

// Сортирует массив строк strs по заданному столбцу
func sortByColumn(strs []string, flg Flags) []string {
	sourceMap := make(map[string]string)
	for _, str := range strs {
		columns := strings.Split(str, " ")
		if len(columns) > flg.ColumnSort {
			sourceMap[columns[flg.ColumnSort]] = str
			continue
		}
		sourceMap[columns[0]] = str
	}

	keysToBeSorted := getKeysOfMap(sourceMap)
	quickSort(keysToBeSorted, 0, len(keysToBeSorted)-1, flg.NumericSort)
	sortedStrs := make([]string, len(strs))

	for i, key := range keysToBeSorted {
		sortedStrs[i] = sourceMap[key]
	}

	return sortedStrs
}

func Sort(flg Flags) {
	// Считываем строки из входного файла
	strs := readFile(flg.Input)

	// Создаем пустой слайс result
	result := make([]string, len(strs))
	// копируем строки из strs в result
	copy(result, strs)

	if flg.ColumnSort >= 0 {
		// Сортируем массив по заданному столбцу
		result = sortByColumn(result, flg)
		// Решаем проблему с заглавными буквами
		result = fixUpperCase(result)
	} else {
		// Сортировка всего слайса
		result = quickSort(result, 0, len(result)-1, false)
		// Решаем проблему с заглавными буквами
		result = fixUpperCase(result)
	}

	if flg.UniqueValues {
		// Сортируем всего слайса
		result = quickSort(result, 0, len(result)-1, false)
		// Решаем проблему с заглавными буквами
		result = fixUpperCase(result)
		// Оставляем только уникальные строки
		result = onlyUnique(result)
	}

	if flg.ReverseSort {
		// Инвертируем порядок строк
		result = reversed(result)
	}

	// Записываем строки в выходной файл
	writeToFile(flg.Output, result)
}

// Чтение из файла
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

// Запись в файл
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

// Парсит аргументы командной строки
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

	// Если флаги in и out отсутствуют
	if *in == "" || *out == "" {
		log.Fatal("required: input and output files")
	}

	return flg
}

func main() {
	flg := parseFlags()
	Sort(flg)
}
