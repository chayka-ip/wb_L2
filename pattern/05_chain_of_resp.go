package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
	Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern

	Поведенческий паттерн «цепочка вызовов» позволяет последовательно передавать запрос
	по цепочке обработчиков до тех пор пока один из них не обработает его, либо
	каждый из них не выполнит свои функции по обработке запроса.

	Результат обработки может быть разным: прекращение обработки при невыполнении некоторых условий,
	частичная обработка запроса, передача управления последующему коду при успешной обработке и тд.

	Отдельная проверка / обработка помещается в изолированный обьект, реализующий общий интерфейс обработчка.
	Из этих обьектов можно выстраивать	цепочки неограниченной длины.

	Данный паттерн часто применяется при создании middleware компонентов, валидации данных,
	графических интерфейсах (для определения компонента при взаимодействии) и тд.

	+++

	Паттерн удобно применять, когда существует необходимость обрабатывать запросы разными спосообами,
	которых может быть много, или они часто меняются. Заранее может быть неизвестно, какие запросы придут
	и какие обработчики могут понадобиться.

	Цепочка может формироваться динамически исходя из нужд конкретного запроса.

	Обьединение обработчиков в цепочку гарантирует последовательный порядок выполнения проверок.

	Разбивка обработчкиков на небольшие обьекты, выполняющие одну функциию,
	повышает читабельность и поддерживаемость кода.

	Снижается связанность между запросом и обработчиками.

	---
	Запрос может остаться необработанным.

	Необходимо учитывать различные ситуации, возникающие при неполной или неуспешной обработке запроса.

	Есть риск возникновения большого количества типов обработчиков и разрастания кодовой базы.

*/

// Incoming data
type Request struct {
	Name     string
	Password string
	Data     string
	IsAdmin  bool
	IsValid  bool
}

// Handler interface
type Middleware interface {
	handle(*Request)
	setNext(Middleware)
}

// Concrete handler
type Authenticator struct {
	next Middleware
}

func (a *Authenticator) handle(r *Request) {
	if r.Name == "user" && r.Password == "pwd123" {
		fmt.Println("Authenticator: User is authenticated.")
		if a.next != nil {
			a.next.handle(r)
		}
		return
	}
	fmt.Println("Authenticator: User is not authenticated.")
}

func (a *Authenticator) setNext(m Middleware) {
	a.next = m
}

// Concrete handler
type Authorizer struct {
	next Middleware
}

func (a *Authorizer) handle(r *Request) {
	if r.IsAdmin {
		fmt.Println("Authorizer: Admin rights are granted.")
		if a.next != nil {
			a.next.handle(r)
		}
		return
	}
	fmt.Println("Authorizer: Unauthorized access.")
}

func (a *Authorizer) setNext(m Middleware) {
	a.next = m
}

// Concrete handler
type DataValidator struct {
	next Middleware
}

func (a *DataValidator) handle(r *Request) {
	if r.Data == "payload" {
		fmt.Println("DataValidator: Data is valid.")
		r.IsValid = true
		if a.next != nil {
			a.next.handle(r)
		}
		return
	}
}

func (a *DataValidator) setNext(m Middleware) {
	a.next = m
}

// Concrete handler
type DataProcessor struct {
}

func (a *DataProcessor) handle(r *Request) {
	if r.IsValid {
		fmt.Println("DataProcessor: Processing data...")
		return
	}
	fmt.Println("DataProcessor: Unable to handle request...")
}
func (a *DataProcessor) setNext(m Middleware) {
}

//
func main5() {
	request := Request{Name: "user", Password: "pwd123", Data: "payload", IsAdmin: true}

	m1 := &Authenticator{}
	m2 := &Authorizer{}
	m3 := &DataValidator{}
	m4 := &DataProcessor{}

	m1.setNext(m2)
	m2.setNext(m3)
	m3.setNext(m4)

	m1.handle(&request)

	/*
		Output:

		Authenticator: User is authenticated.
		Authorizer: Admin rights are granted.
		DataValidator: Data is valid.
		DataProcessor: Processing data...
	*/
}
