package main

import (
	"sort"
	"strings"
)

/*
	Напишите функцию поиска всех множеств анаграмм по словарю.
	Например:
	'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
	'листок', 'слиток' и 'столик' - другому.
	Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
	Выходные данные: Ссылка на мапу множеств анаграмм.
	Ключ - первое встретившееся в словаре слово из множества
	Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
	Множества из одного элемента не должны попасть в результат.
	Все слова должны быть приведены к нижнему регистру.
	В результате каждое слово должно встречаться только один раз.
	Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Сортировка массива строк с использование алгоритма быстрой сортировки
func quickSort(words []string, start, end int) *[]string {
	if start < end {
		// Выбор опорного элемента
		pivot := words[start]
		left := start + 1
		right := end

		for left <= right {
			for left <= right && words[left] <= pivot {
				left++
			}

			for left <= right && words[right] >= pivot {
				right--
			}

			if left < right {
				words[left], words[right] = words[right], words[left]
			}
		}

		words[start], words[right] = words[right], words[start]

		// Рекурсивное применение quickSort к двум частям массива
		words = *quickSort(words, start, right-1)
		words = *quickSort(words, right+1, end)
	}

	return &words
}

// Приводит строки к нижнему регистру
func toLower(words *[]string) *[]string {
	result := make([]string, len(*words))

	for i, str := range *words {
		result[i] = strings.ToLower(str)
	}

	return &result
}

// Сортирует символы в строке
func sortChars(word string) string {
	chars := strings.Split(word, "")
	sort.Strings(chars)

	return strings.Join(chars, "")
}

// Поиск анаграмм
func FindAnagrams(dict *[]string) *map[string][]string {
	// Приводим массив строк к нижнему регистру
	dict = toLower(dict)

	// Применяем быструю сортировку к масиву строк
	dict = quickSort(*dict, 0, len(*dict)-1)

	// Создаем временную карты для хранения анаграмм
	tempMap := make(map[string][]string)
	nameMap := make(map[string]string)

	// Поиск аннаграм в массиве строк
	for _, word := range *dict {
		if len(word) > 1 {
			sortedWord := sortChars(word)
			if _, ok := tempMap[sortedWord]; ok {
				tempMap[sortedWord] = append(tempMap[sortedWord], word)
				continue
			}
			tempMap[sortedWord] = make([]string, 0, 1)
			nameMap[sortedWord] = word
		}
	}

	// Формируем результирующую карту анаграмм
	resultMap := make(map[string][]string)
	for sortedWord, word := range tempMap {
		if len(word) != 0 {
			resultMap[nameMap[sortedWord]] = word
		}

	}

	return &resultMap
}
