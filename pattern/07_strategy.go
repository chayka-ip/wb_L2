package pattern

import "fmt"

/*
	Реализовать паттерн «стратегия».
	Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern

	Поведенческий паттерн «стратегия» позволяет объединить семейство алгоритмов под общим интерфейсом
	и взаимозаменять их по необходимости во время выполнения.

	Паттерн может применяться при создании системы прокладки маршрута. Принцип построения маршрута остается неизменным,
	но реализация будет меняться в зависимости от типа транспора (можно построить пеший маршрут, маршрут на автомобиле,
	на общественном транспорте, на поезеде, на самолете, на корабле).

	+++

	Паттерн удобен, когда нужно использовать разные версии алгоритма одного вида для работы некоторого объекта.

	Паттерн позволяет скрыть реализации алгоритмов от внешнего мира.

	При помощи паттерна можно заменять используемый алгоритм в зависимости от текущего состояния программы и
	делегировать им требуемую работу.

	---

	Добавление новых типов в кодовую базу.

	Клиентский код должен понимать различия между конкретными реализациями алгоритма.

*/

// Strategy interface
type IDelivery interface {
	deliver(string)
}

// Concrete strategy
type PostDelivery struct{}

func NewPostDelivery() PostDelivery {
	return PostDelivery{}
}
func (d PostDelivery) deliver(address string) {
	fmt.Println("Delivering to post office...")
}

// Concrete strategy
type CourierDelivery struct{}

func NewCourierDelivery() CourierDelivery {
	return CourierDelivery{}
}

func (d CourierDelivery) deliver(address string) {
	fmt.Println("Delivering by address...")
}

// Concrete strategy
type PickupDelivery struct{}

func NewPickupDelivery() PickupDelivery {
	return PickupDelivery{}
}
func (d PickupDelivery) deliver(address string) {
	fmt.Println("Delivering to pick-up point...")
}

// Context
type OrderHandler struct {
	delivery  IDelivery
	orderInfo string
	address   string
}

func NewOrderHandler(delivery IDelivery) *OrderHandler {
	return &OrderHandler{
		delivery: delivery,
	}
}
func (o *OrderHandler) SetDeliveryStrategy(d IDelivery) {
	o.delivery = d
}

func (o *OrderHandler) SetOrderData(info, address string) {
	o.orderInfo = info
	o.address = address
}

func (o *OrderHandler) ProcessOrder() {
	fmt.Println("Checking order info...")
	fmt.Println("OK")
	fmt.Println("Preparing for delivery...")
	o.delivery.deliver(o.address)
}

func main7() {
	postDelivery := NewPostDelivery()
	courierDelivery := NewCourierDelivery()
	pickupDelivery := NewPickupDelivery()

	orderInfo1 := "Vladimir Vladimirov; laptop x2"
	orderAddr1 := "Moscow, Kremlin"

	orderInfo2 := "Dmitri Ivanov; computer chair x1"
	orderAddr2 := "St. Petersburg, Line 1, b. 16"

	orderInfo3 := "Otto Miranen; keyboard"
	orderAddr3 := "Helsinki, Skölvägen 20"

	op := NewOrderHandler(courierDelivery)
	op.SetOrderData(orderInfo1, orderAddr1)
	op.ProcessOrder()

	op.SetOrderData(orderInfo2, orderAddr2)
	op.SetDeliveryStrategy(pickupDelivery)
	op.ProcessOrder()

	op.SetOrderData(orderInfo3, orderAddr3)
	op.SetDeliveryStrategy(postDelivery)
	op.ProcessOrder()

	/*
		Output:

		Checking order info...
		OK
		Preparing for delivery...
		Delivering by address...

		Checking order info...
		OK
		Preparing for delivery...
		Delivering to pick-up point...

		Checking order info...
		OK
		Preparing for delivery...
		Delivering to post office...
	*/
}
