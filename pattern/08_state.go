package pattern

/*
	Паттерн "Состояние" используется для управления объектом,
	чье поведение зависит от его текущего состояния.
	Он позволяет объекту динамически изменять свое поведение,
	исходя из внутреннего состояния, без привязки к
	конкретным условиям или переполнению кода условиями.
*/

/*
	В данном файле реализован шаблон State на примере режимов стиральной машины
*/

import "fmt"

// интерфейс состояния с методом оповещения
type WashingMachineState interface {
	Alert()
}

// Состояние стиральной машины (хлопок)
type cottonMode struct{}

func (m *cottonMode) Alert() {
	fmt.Println("washing machine set to cotton mode")
}

// Состояние стиральной машины (синтетика)
type syntheticsMode struct{}

func (m *syntheticsMode) Alert() {
	fmt.Println("washing machine set to synthetics mode")
}

// Состояние стиральной машины (шерсть)
type woolMode struct{}

func (m *woolMode) Alert() {
	fmt.Println("washing machine set to wool mode")
}

// Структура контекста с инкапсулированным состоянием
type WashingMachineCtx struct {
	state WashingMachineState
}

func (ctx *WashingMachineCtx) Alert() {
	ctx.state.Alert()
}

func main() {
	// Создаем контекст
	context := WashingMachineCtx{}

	// Меняем режим стирки, что приводит к изменению действия метода Alert
	context.state = &cottonMode{}
	context.Alert()

	context.state = &syntheticsMode{}
	context.Alert()

	context.state = &woolMode{}
	context.Alert()
}
