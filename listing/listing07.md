Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
1
2
3
4
5
6
7
8
0
0
0
0
...

```

- Приведенная выше программа преобразует переданные числа в каналы и соединяет их в один общий канал. После записи всех чисел исходные каналы закрываются. В Go можно считывать значения из закрытого канала (при этом возвращается нулевое значение для типа данного канала), что и происходит в функции merge. Это приводит к тому, что программа будет бесконечно читать дефолтные значения из общего канала, который нигде не закрывается.
- Для корректной работы программы необходимо внести следующие изменения в функцию ```merge```:

```go
func merge(a, b <-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)

	writer := func(c <-chan int) {
		defer wg.Done()
		for v := range c {
			out <- v
		}
	}

	go writer(a)
	go writer(b)

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
```