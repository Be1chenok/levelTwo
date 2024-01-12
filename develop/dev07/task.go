package main

import (
	"fmt"
	"sync"
	"time"
)

/*
	Реализовать функцию, которая будет объединять один или более done-каналов в single-канал,
	если один из его составляющих каналов закроется.
	Очевидным вариантом решения могло бы стать выражение при использованием select,
	которое бы реализовывало эту связь, однако иногда неизвестно общее число done-каналов,
	с которыми вы работаете в рантайме. В этом случае удобнее использовать вызов единственной функции, которая,
	приняв на вход один или более or-каналов, реализовывала бы весь функционал.

	Определение функции:
	var or func(channels ...<- chan interface{}) <- chan interface{}
	Пример использования функции:
	sig := func(after time.Duration) <- chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or (
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf(“fone after %v”, time.Since(start))

*/

func or(channels ...<-chan interface{}) <-chan interface{} {
	// Создаем канал out который будет использоваться для объединения значений из done-каналов
	out := make(chan interface{})

	// WaitGroup для ожидания завершения всех done-каналов
	var wg sync.WaitGroup
	// Добавляем в счетчик число равное количеству done-каналов
	wg.Add(len(channels))

	// Для каждого done канала запускается отдельная горутина
	for _, ch := range channels {
		go func(ch <-chan interface{}) {
			for value := range ch {
				// запись значений из done-канала в out
				out <- value
			}
			// декрементируем счетчик
			wg.Done()
		}(ch)
	}

	go func() {
		// Ожидаем пока счетчик не станет равным 0
		wg.Wait()
		// Закрываем канал
		close(out)
	}()

	return out
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))

}
