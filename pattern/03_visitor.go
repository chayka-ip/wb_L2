package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
	Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern


	Поведенческий паттерн «посетитель» позволяет разместить поведение некоторого типа за пределами этого типа.
	Для этого создается вспомогательный обьект - посетитель, который знает о существовании всех типов,
	которые он должен посещать, принимает обьекты данных типов в качестве аргуметов и реализует
	необходимое поведение.

	Посещаемые типы должны реализовывать метод принятия посетителя
	для обеспечения механизма двойной диспетчеризации.

	Каждый посетитель должен выполнять одну операцию.
	Это продиктовано принципом единственной ответственности.

	В качестве примера практического использования можно привести случай, когда существует несколько разнородных типов,
	представляющих некоторые данные, и возникает необходимость экспортировать эти данные в различные форматы.
	Логика экспорта инкапсулируется в посетителей и не загрязняет основной код.

	+++

	Паттерн хорошо подходит для случаев, когда необходимо совершать некоторую операцию
	над множеством обьектов различных типов.
	Операция может быть чужеродной для кода самих типов или применяться не для всех типов
	одной иерархии, а поэтому должна быть вынесена за пределы этих типов.

	Подход упрощает добавление новых операций, а код самих типов не изменяется (иногда нет даже возможности его изменить).

	---

	Кодовую базу, использующую данный паттерн, сложно поддерживать, если иерархия типов часто меняется.

	Иногда для совершения операции требуется доступ к приватным полям, 	что приведет
	к нарушению инкапсуляции или невозможность применить данный паттерн.

*/

type Visitor interface {
	VisitSchool(*School)
	VisitFactory(*Factory)
	VisitBusinessCenter(*BusinessCenter)
}

// Concrete visitor
type TaxCollectorVisitor struct {
}

func (v TaxCollectorVisitor) VisitSchool(s *School) {
	fmt.Printf("TaxCollectorVisitor: Collecting taxes from: %s\n", s.GetType())
}
func (v TaxCollectorVisitor) VisitFactory(f *Factory) {
	fmt.Printf("TaxCollectorVisitor: Collecting taxes from: %s\n", f.GetType())
}
func (v TaxCollectorVisitor) VisitBusinessCenter(b *BusinessCenter) {
	fmt.Printf("TaxCollectorVisitor: Collecting taxes from: %s\n", b.GetType())
}

// Concrete visitor
type EmployeeInspectorVisitor struct {
}

func (v EmployeeInspectorVisitor) VisitSchool(s *School) {
	fmt.Printf("EmployeeInspectorVisitor: Inspecting stuff in: %s\n", s.GetType())
}
func (v EmployeeInspectorVisitor) VisitFactory(f *Factory) {
	fmt.Printf("EmployeeInspectorVisitor: Inspecting stuff in: %s\n", f.GetType())

}
func (v EmployeeInspectorVisitor) VisitBusinessCenter(b *BusinessCenter) {
	fmt.Printf("EmployeeInspectorVisitor: Inspecting stuff in: %s\n", b.GetType())
}

//
type Institution interface {
	GetType() string
	Accept(Visitor)
}

// Concrete Institution
type School struct {
}

func (s *School) GetType() string {
	return "School"
}
func (s *School) Accept(v Visitor) {
	v.VisitSchool(s)
}

// Concrete Institution
type Factory struct {
}

func (f *Factory) GetType() string {
	return "Factory"
}
func (f *Factory) Accept(v Visitor) {
	v.VisitFactory(f)
}

// Concrete Institution
type BusinessCenter struct {
}

func (b *BusinessCenter) GetType() string {
	return "Business Center"
}
func (b *BusinessCenter) Accept(v Visitor) {
	v.VisitBusinessCenter(b)
}

func main3() {
	taxVisitor := TaxCollectorVisitor{}
	inspectorVisitor := EmployeeInspectorVisitor{}

	school := &School{}
	factory := &Factory{}
	bc := &BusinessCenter{}

	school.Accept(taxVisitor)
	factory.Accept(taxVisitor)
	bc.Accept(taxVisitor)

	school.Accept(inspectorVisitor)
	factory.Accept(inspectorVisitor)
	bc.Accept(inspectorVisitor)

	/*
		Output:

		TaxCollectorVisitor: Collecting taxes from: School
		TaxCollectorVisitor: Collecting taxes from: Factory
		TaxCollectorVisitor: Collecting taxes from: Business Center
		EmployeeInspectorVisitor: Inspecting stuff in: School
		EmployeeInspectorVisitor: Inspecting stuff in: Factory
		EmployeeInspectorVisitor: Inspecting stuff in: Business Center
	*/
}
