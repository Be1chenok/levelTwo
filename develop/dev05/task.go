package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
	Реализовать утилиту фильтрации по аналогии с консольной утилитой
	(man grep — смотрим описание и основные параметры).

	Реализовать поддержку утилитой следующих ключей:
	-A - "after" печатать +N строк после совпадения
	-B - "before" печатать +N строк до совпадения
	-C - "context" (A+B) печатать ±N строк вокруг совпадения
	-c - "count" (количество строк)
	-i - "ignore-case" (игнорировать регистр)
	-v - "invert" (вместо совпадения, исключать)
	-F - "fixed", точное совпадение со строкой, не паттерн
	-n - "line num", напечатать номер строки
*/

type Flags struct {
	Input            string // in входной файл
	Pattern          string // -p паттерн
	AfterLines       int    // -A
	BeforeLines      int    // -B
	ContextLines     int    // -C
	CountLines       bool   // -c
	IgnoreCase       bool   // -i
	InvertMatch      bool   // -v
	FixedStringMatch bool   // -F
	PrintLineNumbers bool   // -n
}

func Grep(flg Flags) {
	// Открываем файл
	file, err := os.Open(flg.Input)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// Читаем файл построчно
	scanner := bufio.NewScanner(file)
	matchedLines := make([]string, 0)
	lineNum := 0
	matchFound := false

	// Цикл сканирования строк
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Проверка на соответствие условия поиска
		/*
			Если установлены флаги IgnoreCase и FixedStringMatch,
			выполняется сравнение без учета регистра символов
		*/
		if flg.IgnoreCase {
			if flg.FixedStringMatch {
				if strings.EqualFold(line, flg.Pattern) {
					matchFound = true
				}
				/*
					Если установлен только флаг IgnoreCase,
					выполняется поиск с игнорированием регистра символов
				*/
			} else {
				if strings.Contains(strings.ToLower(line), strings.ToLower(flg.Pattern)) {
					matchFound = true
				}
			}
		} else {
			/*
				Если установлен только флаг FixedStringMatch,
				выполняется поиск для точного совпадения строки
			*/
			if flg.FixedStringMatch {
				if line == flg.Pattern {
					matchFound = true
				}
				// Если ни один флаг не установлен, выполняется обычный поиск
			} else {
				if strings.Contains(line, flg.Pattern) {
					matchFound = true
				}
			}
		}

		// Обработка найденной строки
		if matchFound {
			/*
				Если установлен флаг CountLines,
				продолжается перебор следующих строк,
				без вывода самих строк,
				только подсчитывается количество совпадений
			*/
			if flg.CountLines {

				continue
			}

			// Если установлен флаг PrintLineNumbers, выводится номер строки
			if flg.PrintLineNumbers {

				fmt.Printf("%d:", lineNum)
			}

			// Вывод строки
			fmt.Println(line)
			matchedLines = append(matchedLines, line)
			matchFound = false

			// Обработка ситуации, когда соответствие не найдено
		} else {
			/*
				Если установлен флаг InvertMatch, выводятся строки,
				которые не соответствуют условию поиска,
				с учетом флага ContextLines, если он установлен
			*/
			if flg.InvertMatch {
				if flg.ContextLines > 0 {
					if len(matchedLines) > flg.ContextLines {
						matchedLines = matchedLines[1:]
					}
					matchedLines = append(matchedLines, line)
					fmt.Println(line)
				}
				/*
					Если флаг InvertMatch не установлен,выполняются действия
					в зависимости от флага BeforeLines
				*/
			} else {
				if flg.BeforeLines > 0 {
					/*
						Если найденных строк меньше или равно заданному значению BeforeLines,
						выводятся и сохраняются строки до текущей строки
					*/
					if len(matchedLines) <= flg.BeforeLines {
						fmt.Println(line)
						matchedLines = append(matchedLines, line)
						/*
							Если найденных строк больше заданного значения BeforeLines,
							удаляется самая ранняя найденная строка из слайса и
							добавляется текущая строка в конец слайса
						*/
					} else {
						matchedLines = matchedLines[1:]
						matchedLines = append(matchedLines, line)
						fmt.Println(line)
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
}

// Парсит аргументы командной строки
func parseFlags() Flags {
	afterLines := flag.Int("A", 0, "print N lines after each match")
	beforeLines := flag.Int("B", 0, "print N lines before each match")
	contextLines := flag.Int("C", 0, "print N lines around each match (before and after)")
	countLines := flag.Bool("c", false, "print count of matching lines")
	ignoreCase := flag.Bool("i", false, "perform case-insensitive matching")
	invertMatch := flag.Bool("v", false, "invert the match (exclude matching lines)")
	fixedStringMatch := flag.Bool("F", false, "search for fixed string instead of a pattern")
	printLineNumbers := flag.Bool("n", false, "print line numbers with output")

	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatalf("usage: go run task.go file pattern")
	}

	input := flag.Arg(0)
	pattern := flag.Arg(1)

	flg := Flags{
		Input:            input,
		Pattern:          pattern,
		AfterLines:       *afterLines,
		BeforeLines:      *beforeLines,
		ContextLines:     *contextLines,
		CountLines:       *countLines,
		IgnoreCase:       *ignoreCase,
		InvertMatch:      *invertMatch,
		FixedStringMatch: *fixedStringMatch,
		PrintLineNumbers: *printLineNumbers,
	}

	return flg
}

func main() {
	flg := parseFlags()
	Grep(flg)
}
