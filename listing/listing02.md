Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```

Ответ:
```
2
1

В функции test() выведет 2 так как x именован на уровне функции
defer изменит x после возврата из функции
Во функции anotherTest() выведется 1,
так как defer уже не может повлиять на возвращенное значение

defer не выполнится пока окружающая его функция не завершит свое выполнение.

```
