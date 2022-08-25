package pattern

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

/*
	Реализовать паттерн «строитель».
	Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern


	Порождающий паттерн «строитель» применяется, когда необходимо создавать сложные обьекты,
	конфигурация которых производится в несколько этапов.

	Чтобы избежать усложнения кода самого обьекта реализация его создания выносится в отдельные
	вспомогательные обьекты, называемые строителями.

	Каждый конкретный строитель знает все об инициализации обьекта с конкретной конфигурацией
	и реализует общий интерфейс создания обьекта, состоящий из набора известных шагов.

	Для удобства использования строителей часто создается вспомогательный класс, называемый «директором»,
	который принимает в себя конкретного строителя и берет на себя процесс вызова шагов строительства обьекта.

	Это помогает полностью скрыть процесс инициализации от клиентского кода.

	Данный паттерн особенно целесообразно применять когда для сложного объекта известны
	общие случаи использования. На основе этих требований создаются базовые темплейты.

	+++

	Поскольку строитель создает обьекты пошагово, у каждого строителя необходимо реализовать лишь те, которые нужны.
	Если бы для создания объектов был применен один общий метод / конструктор, то пришлось бы каждый раз передавать
	большое количество неиспользуемых параметров и изменять этот код каждый раз, при добавлении новых параметров,
	а это нарушает Open-Closed принцип.

	Паттерн позволяет переиспользовать один и тот же код для создания различных или похожих объектов.

	Процесс создания скрыт от клиентского кода.

	Исключается риск выдачи клиенту не до конца инициализированного обьекта.

	---

	Большое количество обьектов строителей (конфигураций) усложнит кодовую базу.

	Клиент может оказаться связан с конкретными реализациями строителей,
	если не удалось создать общий метод получения результата.


*/

// Product
type City struct {
	Name              string
	BuildingMaterials string
	BuildingCount     uint
	CitizenCount      uint
}

func (c *City) PrintInfo() {
	fmt.Printf("City: %s | Materials: %s | Building count: %d | Citizen count: %d\n",
		c.Name, c.BuildingMaterials, c.BuildingCount, c.CitizenCount)
}

// Director
type Director struct {
	builder ICityBuilder
}

func (d *Director) SetBuilder(b ICityBuilder) {
	d.builder = b
}

func (d *Director) BuildCity() {
	d.builder.SetName()
	d.builder.SetBuildingMaterials()
	d.builder.SetBuildingCount()
	d.builder.SetCitizenCount()
}

// Builder interface
type ICityBuilder interface {
	SetName()
	SetBuildingMaterials()
	SetBuildingCount()
	SetCitizenCount()
	GetCity() City
}

// Concrete builder
type RealCityBuilder struct {
	City
	MaxCitizenCount uint
}

func NewRealCityBuilder(MaxCitizenCount uint) *RealCityBuilder {
	return &RealCityBuilder{
		MaxCitizenCount: MaxCitizenCount,
	}
}
func (b *RealCityBuilder) SetName() {
	b.Name = "上海 (Shanghai)"
}
func (b *RealCityBuilder) SetBuildingMaterials() {
	b.BuildingMaterials = "concrete, steel, wood, glass"
}
func (b *RealCityBuilder) SetBuildingCount() {
	rand.Seed(time.Now().UnixNano())
	b.BuildingCount = uint(rand.Intn(100000) + 500)
}
func (b *RealCityBuilder) SetCitizenCount() {
	b.CitizenCount = uint(0.75 * float64(b.MaxCitizenCount))
}
func (b *RealCityBuilder) GetCity() City {
	return City{
		Name:              b.Name,
		BuildingMaterials: b.BuildingMaterials,
		BuildingCount:     b.BuildingCount,
		CitizenCount:      b.CitizenCount,
	}
}

// Concrete builder
type LegoCityBuilder struct {
	City
	MaxBuildingCount uint
}

func NewLegoCityBuilder(MaxBuildingCount uint) *LegoCityBuilder {
	return &LegoCityBuilder{
		MaxBuildingCount: MaxBuildingCount,
	}
}
func (b *LegoCityBuilder) SetName() {
	b.Name = "Lego city"
}
func (b *LegoCityBuilder) SetBuildingMaterials() {
	b.BuildingMaterials = "plastic"
}
func (b *LegoCityBuilder) SetBuildingCount() {
	b.BuildingCount = b.MaxBuildingCount / 2
}
func (b *LegoCityBuilder) SetCitizenCount() {
	b.CitizenCount = b.BuildingCount * 3
}
func (b *LegoCityBuilder) GetCity() City {
	return City{
		Name:              b.Name,
		BuildingMaterials: b.BuildingMaterials,
		BuildingCount:     b.BuildingCount,
		CitizenCount:      b.CitizenCount,
	}
}

//
func main2() {

	realCityBuilder := NewRealCityBuilder(uint(math.Pow(10, 7)))
	legoCityBuilder := NewLegoCityBuilder(500)
	director := Director{}

	director.SetBuilder(realCityBuilder)
	director.BuildCity()
	realCity := realCityBuilder.GetCity()
	// Output: City: 上海 (Shanghai) | Materials: concrete, steel, wood, glass | Building count: 36199 | Citizen count: 7500000
	realCity.PrintInfo()

	director.SetBuilder(legoCityBuilder)
	director.BuildCity()
	legoCity := legoCityBuilder.GetCity()
	// Output: City: Lego city | Materials: plastic | Building count: 250 | Citizen count: 750
	legoCity.PrintInfo()

}
