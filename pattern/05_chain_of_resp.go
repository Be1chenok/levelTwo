package pattern

import "fmt"

/*
	Паттерн Chain of Responsibility используется для организации обработки запросов или
	событий в виде цепочки объектов-обработчиков. Каждый объект-обработчик получает запрос
	и самостоятельно решает, может ли он обработать запрос. Если объект не может обработать запрос,
	он передает его следующему объекту в цепочке.
*/

/*
	В данном файле реализован паттерн цепочка вызовов на примере обработки покупок в интернет-магазине
*/

// Интерфейс обработчика
type PurchaseHandler interface {
	HandlePurchase(purchase *purchase)
	SetNextHandler(next PurchaseHandler)
}

// Структура представляющая покупку
type purchase struct {
	product   string
	quantity  int
	totalCost float32
}

func NewPurchase(product string, quantity int, totalCost float32) *purchase {
	return &purchase{
		product:   product,
		quantity:  quantity,
		totalCost: totalCost,
	}
}

// Базовая реализация обработчика
type basePurchaseHandler struct {
	nextHandler PurchaseHandler
}

func (bph *basePurchaseHandler) SetNextHandler(next PurchaseHandler) {
	bph.nextHandler = next
}

func (bph *basePurchaseHandler) HandlePurchase(purchase *purchase) {
	if bph.nextHandler != nil {
		bph.nextHandler.HandlePurchase(purchase)
	}
}

// Обработчик для скидок на товары
type discountHandler struct {
	basePurchaseHandler
}

// Обработчик для доставки
type deliveryHandler struct {
	basePurchaseHandler
}

func (dh *deliveryHandler) HandlePurchase(purchase *purchase) {
	if purchase.totalCost > 5000 {
		purchase.totalCost -= 5 // Применяем скидку на доставку
		fmt.Println("delivery discount applied")
	}

	// Передаем обработку следующему обработчику
	dh.basePurchaseHandler.HandlePurchase(purchase)
}

// Обработчик для подтверждения заказа
type confirmationHandler struct {
	basePurchaseHandler
}

func (ch *confirmationHandler) HandlePurchase(purchase *purchase) {
	fmt.Printf("order %s confirmed\n", purchase.product)

	// Передаем обработку следующему обработчику
	ch.basePurchaseHandler.HandlePurchase(purchase)
}

func main() {
	// Создаем цепочку обработчиков
	discountHandler := new(discountHandler)
	deliveryHandler := new(deliveryHandler)
	confirmationHandler := new(confirmationHandler)

	// Создаем и обрабатываем покупки
	discountHandler.SetNextHandler(deliveryHandler)
	deliveryHandler.SetNextHandler(confirmationHandler)

	firstPurchase := NewPurchase("laptop", 1, 90000.0)
	discountHandler.HandlePurchase(firstPurchase)

	secondPurchase := NewPurchase("phone", 2, 50000.0)
	discountHandler.HandlePurchase(secondPurchase)

	thirdPurchase := NewPurchase("lamp", 10, 3000.0)
	discountHandler.HandlePurchase(thirdPurchase)
}
