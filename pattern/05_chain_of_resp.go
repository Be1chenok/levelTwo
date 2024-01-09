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
type purchaseHandler interface {
	handlePurchase(purchase *purchase)
	setNextHandler(next purchaseHandler)
}

// Структура представляющая покупку
type purchase struct {
	product   string
	quantity  int
	totalCost float32
}

func newPurchase(product string, quantity int, totalCost float32) *purchase {
	return &purchase{
		product:   product,
		quantity:  quantity,
		totalCost: totalCost,
	}
}

// Базовая реализация обработчика
type basePurchaseHandler struct {
	nextHandler purchaseHandler
}

func (bph *basePurchaseHandler) setNextHandler(next purchaseHandler) {
	bph.nextHandler = next
}

func (bph *basePurchaseHandler) handlePurchase(purchase *purchase) {
	if bph.nextHandler != nil {
		bph.nextHandler.handlePurchase(purchase)
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

func (dh *deliveryHandler) handlePurchase(purchase *purchase) {
	if purchase.totalCost > 5000 {
		purchase.totalCost -= 5 // Применяем скидку на доставку
		fmt.Println("delivery discount applied")
	}

	// Передаем обработку следующему обработчику
	dh.basePurchaseHandler.handlePurchase(purchase)
}

// Обработчик для подтверждения заказа
type confirmationHandler struct {
	basePurchaseHandler
}

func (ch *confirmationHandler) handlePurchase(purchase *purchase) {
	fmt.Printf("order %s confirmed\n", purchase.product)

	// Передаем обработку следующему обработчику
	ch.basePurchaseHandler.handlePurchase(purchase)
}

func main() {
	// Создаем цепочку обработчиков
	discountHandler := new(discountHandler)
	deliveryHandler := new(deliveryHandler)
	confirmationHandler := new(confirmationHandler)

	// Создаем и обрабатываем покупки
	discountHandler.setNextHandler(deliveryHandler)
	deliveryHandler.setNextHandler(confirmationHandler)

	firstPurchase := newPurchase("laptop", 1, 90000.0)
	discountHandler.handlePurchase(firstPurchase)

	secondPurchase := newPurchase("phone", 2, 50000.0)
	discountHandler.handlePurchase(secondPurchase)

	thirdPurchase := newPurchase("lamp", 10, 3000.0)
	discountHandler.handlePurchase(thirdPurchase)
}
