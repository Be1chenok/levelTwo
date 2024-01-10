package pattern

import "fmt"

/*
	В данном файле реализован паттерн строитель на примере строительства дома
*/

/*
	Паттерн строитель используется для создания сложных объектов шаг за шагом.
	Он позволяет создавать различные объекты,используя один и тот же процесс конструирования
*/

// Интерфейс для строительства дома
type HouseBuilder interface {
	SetWindowType()
	SetDoorType()
	SetNumFloor()
	GetHouse() house
}

// Структура дома
type house struct {
	windowType string
	doorType   string
	floor      int
}

func NewHouse(windowType, doorType string, floor int) house {
	return house{
		windowType: windowType,
		doorType:   doorType,
		floor:      floor,
	}
}

// Структура для строительства деревянного дома
type woodenBuilder struct {
	h house
}

func newWoodenBuilder() *woodenBuilder {
	return &woodenBuilder{}
}

func (wb *woodenBuilder) GetHouse() house {
	return NewHouse(wb.h.windowType, wb.h.doorType, wb.h.floor)
}

func (wb *woodenBuilder) SetWindowType() {
	wb.h.windowType = "wooden window"
}

func (wb *woodenBuilder) SetDoorType() {
	wb.h.doorType = "wooden door"
}

func (wb *woodenBuilder) SetNumFloor() {
	wb.h.floor = 1
}

// Структура для строительства каменного дома
type stoneBuilder struct {
	h house
}

func newStoneBuilder() *stoneBuilder {
	return &stoneBuilder{}
}

func (sb *stoneBuilder) GetHouse() house {
	return NewHouse(sb.h.windowType, sb.h.doorType, sb.h.floor)
}

func (sb *stoneBuilder) SetWindowType() {
	sb.h.windowType = "stone window"
}

func (sb *stoneBuilder) SetDoorType() {
	sb.h.doorType = "stone door"
}

func (sb *stoneBuilder) SetNumFloor() {
	sb.h.floor = 2
}

// Директор - управляет процессом постройки дома
type director struct {
	builder HouseBuilder
}

func NewDirector(builder HouseBuilder) *director {
	return &director{
		builder: builder,
	}
}

func (d *director) setBuilder(builder HouseBuilder) {
	d.builder = builder
}

func (d *director) buildHouse() house {
	d.builder.SetWindowType()
	d.builder.SetDoorType()
	d.builder.SetNumFloor()

	return d.builder.GetHouse()
}

func getBuilder(builderType string) HouseBuilder {
	switch builderType {
	case "wooden":
		return newWoodenBuilder()
	case "stone":
		return newStoneBuilder()
	default:
		return nil
	}
}

func main() {
	woodenBuilder := getBuilder("wooden")
	stoneBuilder := getBuilder("stone")

	// Создаем директора и говорим, что будем строить деревянный дом
	director := NewDirector(woodenBuilder)

	// Строим и получаем дом
	woodenHouse := director.buildHouse()

	fmt.Printf("wooden house:\n window type: %s\n door type: %s\n num floor: %d\n",
		woodenHouse.windowType, woodenHouse.doorType, woodenHouse.floor)

	// Говорим что будем строить каменный дом
	director.setBuilder(stoneBuilder)

	// Строим и получаем дом
	stoneHouse := director.buildHouse()

	fmt.Printf("stone house:\n window type: %s\n door type: %s\n num floor: %d\n",
		stoneHouse.windowType, stoneHouse.doorType, stoneHouse.floor)
}
