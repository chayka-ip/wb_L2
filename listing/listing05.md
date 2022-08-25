Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error
```
- Тип ```error``` в Go является интерфейсом. Переменной ```err``` был присвоен нулевой указатель на тип ```customError```. Это привело к тому, что интерфейс стал ненулевым. Интерфейс может быть равен ```nil``` только если и тип и значение этого интерфейса равны ```nil```.
- Чтобы избежать вышерассмотренной проблемы при возврате ошибок из функции необходимо использовать интерфейс ```error``` вместо указателя на тип ошибки.	
