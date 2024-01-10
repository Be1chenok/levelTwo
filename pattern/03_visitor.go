package pattern

/*
	В данном файле реализован паттерн посетитель на примере товаров в магазине
	Следует определить цену конкретного товара по скидке
*/

/*
	Паттерн Visitor используется для разделения алгоритмов и структур данных, на которых они оперируют.
	Он позволяет добавлять новые операции к существующим структурам данных, не изменяя самих структур.
*/

import "fmt"

// Интерфейс посетителя товаров
type ItemVisitor interface {
	VisitBook(book book)
	VisitNotepad(notepad notepad)
}

// Общий интерфейс для всех товаров
type Item interface {
	accept(visitor ItemVisitor)
}

// Структура книги
type book struct {
	title    string
	price    float32
	discount int
}

func newBook(title string, price float32, discount int) Item {
	return &book{
		title:    title,
		price:    price,
		discount: discount,
	}
}

func (b book) accept(visitor ItemVisitor) {
	visitor.VisitBook(b)
}

// Структура блокнота
type notepad struct {
	title    string
	price    float32
	discount int
}

func newNotepad(title string, price float32, discount int) Item {
	return &notepad{
		title:    title,
		price:    price,
		discount: discount,
	}
}

func (n notepad) accept(visitor ItemVisitor) {
	visitor.VisitNotepad(n)
}

// Структура посетителя
type saleVisitor struct{}

func NewSaleVisitor() ItemVisitor {
	return &saleVisitor{}
}

func (v saleVisitor) VisitBook(book book) {
	fmt.Printf("book: %s\n", book.title)
	fmt.Printf("Price: %.2f $\n", book.price)
	discountedPrice := calculateDiscountPrice(book.price, book.discount)
	fmt.Printf("Discounted price: %.2f $\n", discountedPrice)
}

func (v saleVisitor) VisitNotepad(notepad notepad) {
	fmt.Printf("notepad: %s\n", notepad.title)
	fmt.Printf("Price: %.2f $\n", notepad.price)
	discountedPrice := calculateDiscountPrice(notepad.price, notepad.discount)
	fmt.Printf("Discounted price: %.2f $\n", discountedPrice)
}

// Функция для расчета скидочной цены
func calculateDiscountPrice(price float32, discount int) float32 {
	if discount > 0 && discount <= 100 {
		return ((float32(100 - discount)) / 100) * price
	}
	return price
}

func main() {
	// Создаем посетителя
	visitor := NewSaleVisitor()

	// Создаем товары
	grokkingAlgorithms := newBook("Grokking Algorithms", 100.0, 50)
	simpleNotepad := newNotepad("Simple Notepad", 5.25, 15)

	// Применяем посетителя к товарам в результате чего выводится название, цена и цена по скидке
	grokkingAlgorithms.accept(visitor)
	simpleNotepad.accept(visitor)
}
