package pattern

/*
	Паттерн фабричного метода используется для создания объектов без явного указания их конкретных классов.
	Он позволяет делегировать процесс создания объектов подклассам,
	позволяя клиентскому коду работать с абстрактными типами данных, не зависящими от конкретной реализации.
*/

/*
	В данном файле реализован паттерн фабричный метод на примере способов доставки товара
*/

import (
	"fmt"
	"log"
)

const (
	car     = "car"
	scooter = "scooter"
)

// Интерфейс для способов доставки
type DeliveryMethod interface {
	Delivery() string
}

// Интерфейс для фабрик
type DeliveryFactory interface {
	Create() DeliveryMethod
}

// Фабричный метод
func CreateDeliveryFactoryByMethod(method string) (DeliveryFactory, error) {
	switch method {
	case car:
		return &carDeliveryFactory{}, nil
	case scooter:
		return &scooterDeliveryFactory{}, nil
	default:
		return nil, fmt.Errorf("%s is unknown", method)
	}
}

// Способ доставки автомобилем
type carDelivery struct{}

func (c *carDelivery) Delivery() string {
	return "car delivery"
}

// Способ доставки самокатом
type scooterDelivery struct{}

func (t *scooterDelivery) Delivery() string {
	return "scooter delivery"
}

// Фабрика для создания способа доставки автомобилем
type carDeliveryFactory struct{}

func (f *carDeliveryFactory) Create() DeliveryMethod {
	return &carDelivery{}
}

// Фабрика для создания способа доставки самокатом
type scooterDeliveryFactory struct{}

func (f *scooterDeliveryFactory) Create() DeliveryMethod {
	return &scooterDelivery{}
}

func main() {
	// создаем фабрику нужного способа доставки
	deliveryFactory, err := CreateDeliveryFactoryByMethod(scooter)
	if err != nil {
		log.Fatalf("failed to create delivery factory: %v", err)
	}

	//Создаем способ доставки, в данном случае доставка самокатом
	method := deliveryFactory.Create()

	//выводим каким способ будет осуществляться доставка
	fmt.Println(method.Delivery())
}
