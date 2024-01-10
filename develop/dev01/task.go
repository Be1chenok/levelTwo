package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

const timeHost = "0.beevik-ntp.pool.ntp.org"

func main() {
	time, err := ntp.Time(timeHost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(-1)
	}

	fmt.Fprintf(os.Stdout, "current time: %s\n", time.Format("15:04:05"))
	fmt.Fprintf(os.Stdout, "exact time: %s", time.Format("15:04:05.000000000"))
}
