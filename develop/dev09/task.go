package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Flags struct {
	Url        string // -url
	OutputPath string // -out
	Depth      int    // -d
}

// Wget с использование глубины поиска
func Wget(flg Flags) {
	inputLink, ok := normalizedInputUrl(flg.Url)
	if !ok {
		log.Fatalf("invalid url: correct url: https://example.com/ or http://example.com/")
	}

	links := []string{inputLink}

	if err := createOutputFolder(flg.OutputPath); err != nil {
		log.Fatalf("failed to create output folder: %v", err)
	}

	if err := os.Chdir(flg.OutputPath); err != nil {
		log.Fatalf("failed to go into folder: %v", err)
	}

	currentLinks := links
	var newLinks []string
	var err error

	currentDepth := 0
	for currentDepth < flg.Depth {
		for i := range currentLinks {
			newLinks, err = extractLinks(currentLinks[i])
			if err != nil {
				log.Print(err)
			}
			links = append(links, newLinks...)
		}
		currentDepth++
		currentLinks = newLinks
	}

	for _, link := range links {
		downloadLinks(link)
	}
}

// Скачивание сайтов
func downloadLinks(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("%s failed to make GET request: %v", url, err)
	}
	defer resp.Body.Close()

	fileName := getFileName(url)
	file := createFile(fileName)
	defer file.Close()

	if err := writeToFile(file, resp.Body); err != nil {
		return fmt.Errorf("failed to write to file %s: %v", fileName, err)
	}

	return nil
}

// Запись в файл
func writeToFile(file *os.File, body io.Reader) error {
	if _, err := io.Copy(file, body); err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

// Создание папки с загруженными сайтами
func createOutputFolder(outputPath string) error {
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if err := os.Mkdir(outputPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

// Создание файла в котором будет сайт
func createFile(fileName string) *os.File {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("%s failed to create file: %v", fileName, err)
	}

	return file
}

// Извлечение ссылок
func extractLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s failed to make GET request: %v", url, err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s failed to parse body: %v", err, url)
	}

	links := visitNode([]string{}, doc)
	links = normalizedLinks(url, links)

	return links, nil
}

// Обход ссылок в документе и добавление их ко всем ссылкам
func visitNode(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visitNode(links, c)
	}

	return links
}

// Получение имени файла из url сайта
func getFileName(url string) string {
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '/' && url[i-1] == '/' {
			url = url[i+1:]
		}
	}
	url, _ = strings.CutSuffix(url, "/")
	result := fmt.Sprintf("%s.html", strings.ReplaceAll(url, "/", "-"))

	return result
}

// Проверяет первую ссылку на правильность ввода и по надобности в конце добавляет /
func normalizedInputUrl(url string) (string, bool) {
	if !strings.Contains(url, "://") {
		return "", false
	}

	if !strings.HasSuffix(url, "/") {
		return fmt.Sprintf("%s/", url), true
	}

	return url, true
}

// Нормализация ссылок
func normalizedLinks(parentUrl string, links []string) []string {
	for i := range links {
		if !strings.Contains(links[i], "://") {
			links[i], _ = strings.CutPrefix(links[i], "/")
			links[i] = fmt.Sprintf("%s%s", parentUrl, links[i])
		}
	}

	return links
}

// Парсит аргументы командной строки
func parseFlags() Flags {
	depth := flag.Int("d", 0, "traversal depth")
	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatal("usage: go task.go -d <depth> url output_path")
	}

	url := flag.Arg(0)
	out := flag.Arg(1)

	flg := Flags{
		Url:        url,
		OutputPath: out,
		Depth:      *depth,
	}

	// Если глубина обхода меньше 0
	if *depth < 0 {
		log.Fatalf("invalid depth")
	}

	return flg
}

func main() {
	flg := parseFlags()
	Wget(flg)
}
