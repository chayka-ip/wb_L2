package pattern

import "fmt"

/*
	Реализовать паттерн «фабричный метод».
	Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern


	Порождающий паттерн «фабричный метод» позволяет создавать обьекты различных типов без явного указания этих типов.
	В Go невозможно реализовать классический вариант фабричного метода ввиду отсутствия классов
	и механизмов наследования. Наиболее приближенной реализацией является простая фабрика. Возвращаемый тип
	определяется динамически на основе переданных аргументов.

	Паттерн может быть использован при создании библиотек компонентов интерфейса. Основной код библиотеки
	будет работать с интерфейсами (Button, MenuItem, Dialog etc), но сами компоненты будут создаваться
	в стиле, определяемом, например, операционной системой, в которой запущено приложение.

	+++

	Паттерн уместен, когда заранее неизвестны типы или зависимости типов, с которыми планируется работать в клиентском коде.
	Код создания объектов отделен от кода их применения.

	Паттерн позволяет удобно добавлять новые типы, не затрагивая код основной программы.

	При использовании паттерна можно подменять создание новых обьектов уже существующими, сохраненными в кеше.

	---

	При очень большом количестве обьектов фабричный метод может сильно разрастись.

*/

const (
	devicePC     = "PC"
	deviceLaptop = "LAPTOP"
)

// Device interface
type IDevice interface {
	setType(string)

	setNumCPU(uint8)
	setRAM(uint64)
	setHDD(uint64)

	getTypeName() string
	getNumCPU() uint8
	getRAM() uint64
	getHDD() uint64
}

// Concrete product
type Device struct {
	typeName string
	ram      uint64
	hdd      uint64
	numCPU   uint8
}

func (d *Device) setType(t string) {
	d.typeName = t
}
func (d *Device) setNumCPU(n uint8) {
	d.numCPU = n
}
func (d *Device) setRAM(n uint64) {
	d.ram = n
}
func (d *Device) setHDD(n uint64) {
	d.hdd = n
}
func (d *Device) getTypeName() string {
	return d.typeName
}
func (d *Device) getNumCPU() uint8 {
	return d.numCPU
}
func (d *Device) getRAM() uint64 {
	return d.ram
}
func (d *Device) getHDD() uint64 {
	return d.hdd
}

// Concrete product
type PC struct {
	Device
}

func NewPC() *PC {
	return &PC{
		Device: Device{
			typeName: devicePC,
			ram:      64,
			hdd:      1024,
			numCPU:   12,
		},
	}
}

// Concrete product
type Laptop struct {
	Device
}

func NewLaptop() *Laptop {
	return &Laptop{
		Device: Device{
			typeName: deviceLaptop,
			ram:      32,
			hdd:      512,
			numCPU:   4,
		},
	}
}

// Factory method
func getDevice(deviceType string) (IDevice, error) {
	if deviceType == devicePC {
		return NewPC(), nil
	}
	if deviceType == deviceLaptop {
		return NewLaptop(), nil
	}
	return nil, fmt.Errorf("unknown device: %s", deviceType)
}

// Client code

func CompareDevices(a IDevice, b IDevice) {
	f := func(s, ap, bp string) {
		fmt.Printf("%s| a: %s | b: %s\n", s, ap, bp)
	}
	f("Type", a.getTypeName(), b.getTypeName())
	f("CPU", fmt.Sprint(a.getNumCPU()), fmt.Sprint(b.getNumCPU()))
	f("RAM", fmt.Sprint(a.getRAM()), fmt.Sprint(b.getRAM()))
	f("HDD", fmt.Sprint(a.getHDD()), fmt.Sprint(b.getHDD()))
}

func main6() {
	pc, _ := getDevice(devicePC)
	laptop, _ := getDevice(deviceLaptop)
	CompareDevices(pc, laptop)

	_, err := getDevice("fruitPhone")
	if err != nil {
		fmt.Println(err)
	}

	/*
		Output:

		Type| a: PC | b: LAPTOP
		CPU| a: 12 | b: 4
		RAM| a: 64 | b: 32
		HDD| a: 1024 | b: 512
		unknown device: fruitPhone
	*/

}
