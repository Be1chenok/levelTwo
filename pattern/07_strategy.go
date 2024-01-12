package pattern

/*
	Паттерн стратегия позволяет выбирать один из нескольких вариантов поведения во время выполнения программы,
	путем инкапсуляции связанных операций в отдельные классы и предоставления им общего интерфейса.
*/

/*
	В данном файле реализован паттерн стратегия на примере оплаты заказа
*/

import "fmt"

const (
	qiwi     string = "qiwi"
	bankCard string = "bank card"
)

// интерфейс для оплаты
type Payment interface {
	Pay(sum float32) error
}

// структура продукта
type Product struct {
	Name  string
	Price float32
}

// Стратегия (позволяет варьировать способ оплаты не затрагивая логику заказа)
func processOrder(products []Product, payment Payment) error {
	var sum float32
	for _, product := range products {
		sum += product.Price
	}

	if err := payment.Pay(sum); err != nil {
		return fmt.Errorf("failed to make payment: %v", err)
	}

	return nil
}

// Структура оплаты по банковской карте
type cardPayment struct {
	cardNumber, date, cvv string
}

func NewCardPayment(cardNumber, date, cvv string) Payment {
	return &cardPayment{
		cardNumber: cardNumber,
		date:       date,
		cvv:        cvv,
	}
}

func (p *cardPayment) Pay(sum float32) error {
	fmt.Printf("payment via bank card for the amount: %0.2f", sum)
	return nil
}

// Структура оплаты по киви кошельку
type qiwiPayment struct {
	walletNumber, password string
}

func NewQIWIPayment(walletNumber, password string) Payment {
	return &qiwiPayment{
		walletNumber: walletNumber,
		password:     password,
	}
}

func (p *qiwiPayment) Pay(sum float32) error {
	fmt.Printf("payment via QIWI wallet for the amount: %0.2f", sum)
	return nil
}

func main() {
	products := []Product{
		{
			"Samsung Tv",
			45000.0,
		},
		{
			"iPhone 14",
			100000.0,
		},
		{
			"PlayStation 5",
			60000.0,
		},
	}

	var payment Payment
	payMethod := qiwi

	switch payMethod {
	case bankCard:
		payment = NewCardPayment("1234567887654321", "01/25", "123")
	case qiwi:
		payment = NewQIWIPayment("+79991112233", "qwerty123")
	}

	processOrder(products, payment)
}
