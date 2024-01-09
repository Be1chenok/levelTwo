package pattern

/*
	В данном файле паттерн фасад реализован на примере обработки изображения
	В системе есть несколько модулей: загрузка, изменение размера и применение фильтра на изображение
*/

/*
	Паттерн фасад может быть полезен,
	когда требуется упростить сложную систему
	и предоставить удобный интерфейс для её взаимодействия
*/

import "fmt"

// Структура загрузчика изображения
type imageLoader struct{}

func NewImageLoader() *imageLoader {
	return &imageLoader{}
}

func (imgl *imageLoader) LoadImage(filePath string) {
	fmt.Printf("\nimage loading: %s", filePath)
}

// Структура для изменения размера изображения
type imageResizer struct{}

func NewImageResizer() *imageResizer {
	return &imageResizer{}
}

func (imgr *imageResizer) ResizeImage(width, height int) {
	fmt.Printf("\nimage resolution change: width: %d, height: %d", width, height)
}

// Структура для фильтра изображения
type imageFilter struct{}

func NewImageFilter() *imageFilter {
	return &imageFilter{}
}

func (imgf *imageFilter) ApplyFilter(filterType string) {
	fmt.Printf("\nfilter application: %s", filterType)
}

// Структура фасада для обработки изображения
type imageProcessingFacade struct {
	loader  *imageLoader
	resizer *imageResizer
	filter  *imageFilter
}

func NewImageProcessingFacade() *imageProcessingFacade {
	return &imageProcessingFacade{}
}

func (facade *imageProcessingFacade) ProcessImage(filePath string, width, height int, filterType string) {
	facade.loader.LoadImage(filePath)
	facade.resizer.ResizeImage(width, height)
	facade.filter.ApplyFilter(filterType)
}

func main() {
	facade := NewImageProcessingFacade()
	facade.ProcessImage("image.png", 1920, 1080, "inversion")
}
