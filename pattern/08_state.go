package pattern

import "fmt"

/*
	Реализовать паттерн «состояние».
	Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern

	Поведенческий паттерн «состояние» позволяет объектам менят поведение в зависимости от текущего состояния. Он может
	находиться в одном из кончечных состояний и разным образом реагировать на одни и те же события.

	Паттерн уместен при разработке графических редакторов, предоставляющих пользователю выбирать инструменты для работы.
	Пользовательский ввод будет одним и тем же (сигналы клавиатуры, мыши), а производимые действия описываются выбранным
	инструментом, то есть текущим состоянием редактоа.

	+++

	Паттерн применим, когда поведение обьекта значительно меняется в завиимости от состояния,
	при этом существует много состояний, которые часто изменяются.

	Паттерн может быть полезен, если в основных методах типа есть много условных операторов, опирающихся
	на значения полей данного типа.

	---

	Код может стать излишне усложнен, если состояний мало и они редко добавляются / изменяются.

	Могут возникнуть дополнительные накладные расходы при отладке, если состояний излишне много.

*/

// State interface
type IH2OState interface {
	Cool()
	Heat()
}

// Concrete state
type SolidH2OState struct {
	h2o *H2O
}

func (s *SolidH2OState) Cool() {
	fmt.Println("SolidH2OState: Making ice colder...")
}
func (s *SolidH2OState) Heat() {
	fmt.Println("SolidH2OState: Melting ice to water...")
	s.h2o.SetState(&s.h2o.water)
}

// Concrete state
type LiquidH2OState struct {
	h2o *H2O
}

func (s *LiquidH2OState) Cool() {
	fmt.Println("LiquidH2OState: Freezing water to ice...")
	s.h2o.SetState(&s.h2o.ice)
}
func (s *LiquidH2OState) Heat() {
	fmt.Println("LiquidH2OState: Evaporating water to steam...")
	s.h2o.SetState(&s.h2o.gas)
}

// Concrete state
type GaseousH2OState struct {
	h2o *H2O
}

func (s *GaseousH2OState) Cool() {
	fmt.Println("GaseousH2OState: Condensing steam to water...")
	s.h2o.SetState(&s.h2o.water)
}
func (s *GaseousH2OState) Heat() {
	fmt.Println("GaseousH2OState: Making steam hotter...")
}

// Context
type H2O struct {
	ice   SolidH2OState
	water LiquidH2OState
	gas   GaseousH2OState

	currentState IH2OState
}

func NewH2O() *H2O {
	w := &H2O{}
	ice := SolidH2OState{w}
	water := LiquidH2OState{w}
	gas := GaseousH2OState{w}

	w.ice = ice
	w.water = water
	w.gas = gas

	w.currentState = &water

	return w
}
func (h *H2O) SetState(state IH2OState) {
	h.currentState = state
}
func (h *H2O) Cool() {
	h.currentState.Cool()
}
func (h *H2O) Heat() {
	h.currentState.Heat()
}

func main8() {
	water := NewH2O()

	water.Heat()
	water.Heat()
	water.Cool()
	water.Cool()
	water.Cool()

	/*
		Output:

		LiquidH2OState: Evaporating water to steam...
		GaseousH2OState: Making steam hotter...
		GaseousH2OState: Condensing steam to water...
		LiquidH2OState: Freezing water to ice...
		SolidH2OState: Making ice colder...
	*/
}
