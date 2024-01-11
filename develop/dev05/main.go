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
	FilePath         string
	Pattern          string
	AfterLines       int
	BeforeLines      int
	ContextLines     int
	CountLines       bool
	IgnoreCase       bool
	InvertMatch      bool
	FixedStringMatch bool
	PrintLineNumbers bool
}

func Grep(flg Flags) {
	file, err := os.Open(flg.FilePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	matchedLines := make([]string, 0)
	lineNum := 0
	matchFound := false

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if flg.IgnoreCase {
			if flg.FixedStringMatch {
				if strings.EqualFold(line, flg.Pattern) {
					matchFound = true
				}
			} else {
				if strings.Contains(strings.ToLower(line), strings.ToLower(flg.Pattern)) {
					matchFound = true
				}
			}
		} else {
			if flg.FixedStringMatch {
				if line == flg.Pattern {
					matchFound = true
				}
			} else {
				if strings.Contains(line, flg.Pattern) {
					matchFound = true
				}
			}
		}

		if matchFound {
			if flg.CountLines {
				continue
			}

			if flg.PrintLineNumbers {
				fmt.Printf("%d:", lineNum)
			}

			fmt.Println(line)
			matchedLines = append(matchedLines, line)
			matchFound = false
		} else {
			if flg.InvertMatch {
				if flg.ContextLines > 0 {
					if len(matchedLines) > flg.ContextLines {
						matchedLines = matchedLines[1:]
					}
					matchedLines = append(matchedLines, line)
					fmt.Println(line)
				}
			} else {
				if flg.BeforeLines > 0 {
					if len(matchedLines) <= flg.BeforeLines {
						fmt.Println(line)
						matchedLines = append(matchedLines, line)
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

func parseFlags() Flags {
	filePath := flag.String("in", "", "file to search in")
	pattern := flag.String("p", "", "pattern to search for")
	afterLines := flag.Int("A", 0, "print N lines after each match")
	beforeLines := flag.Int("B", 0, "print N lines before each match")
	contextLines := flag.Int("C", 0, "print N lines around each match (before and after)")
	countLines := flag.Bool("c", false, "print count of matching lines")
	ignoreCase := flag.Bool("i", false, "perform case-insensitive matching")
	invertMatch := flag.Bool("v", false, "invert the match (exclude matching lines)")
	fixedStringMatch := flag.Bool("F", false, "search for fixed string instead of a pattern")
	printLineNumbers := flag.Bool("n", false, "print line numbers with output")

	flag.Parse()

	flg := Flags{
		FilePath:         *filePath,
		Pattern:          *pattern,
		AfterLines:       *afterLines,
		BeforeLines:      *beforeLines,
		ContextLines:     *contextLines,
		CountLines:       *countLines,
		IgnoreCase:       *ignoreCase,
		InvertMatch:      *invertMatch,
		FixedStringMatch: *fixedStringMatch,
		PrintLineNumbers: *printLineNumbers,
	}

	if *pattern == "" || *filePath == "" {
		log.Fatal("required: pattern and file")
	}

	return flg
}

func main() {
	flg := parseFlags()
	Grep(flg)
}
