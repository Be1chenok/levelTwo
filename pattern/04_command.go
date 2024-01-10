package pattern

import "fmt"

/*
	Паттерн команда используется для инкапсуляции запроса в виде объекта,
	позволяя параметризовать клиентов с различными запросами,
	организовывать их в очереди, а также поддерживать отмену операций
*/

/*
	Паттерн представлен на примере светодиодной ленты,
	которую можно включить и выключит
*/

// Интерфейс команды
type Command interface {
	Execute()
}

// Receiver - получатель команды
type ledStrip struct {
	isOn bool
	mode int
}

func (l *ledStrip) turnOn() {
	l.isOn = true
	fmt.Printf("LED strip is ON\n")
}

func (l *ledStrip) turnOff() {
	l.isOn = false
	fmt.Printf("LED strip is OFF\n")
}

// turnOnLEDCommand - команда включения световой ленты
type turnOnLEDCommand struct {
	ledStrip *ledStrip
}

// turnOffLEDCommand - команда выключения световой ленты
type turnOffLEDCommand struct {
	ledStrip *ledStrip
}

func (onCmd *turnOnLEDCommand) Execute() {
	onCmd.ledStrip.turnOn()
}

func (offCmd *turnOffLEDCommand) Execute() {
	offCmd.ledStrip.turnOff()
}

// Invoker - вызывающий объект
type remoteControl struct {
	pressOn  Command
	pressOff Command
}

func newRemoteControl(pressOn, pressOff Command) *remoteControl {
	return &remoteControl{
		pressOn:  pressOn,
		pressOff: pressOff,
	}
}

func (rc *remoteControl) pressOnExec() {
	rc.pressOn.Execute()
}

func (rc *remoteControl) pressOffExec() {
	rc.pressOff.Execute()
}

func main() {
	// создаем светодиодную ленту
	ledStrip := new(ledStrip)

	// Создаем команды
	onLED := turnOnLEDCommand{ledStrip: ledStrip}
	offLED := turnOffLEDCommand{ledStrip: ledStrip}

	// Создаем пульт управления
	rControl := newRemoteControl(&onLED, &offLED)

	// Нажимаем кнопки
	rControl.pressOnExec()
	rControl.pressOffExec()
}
