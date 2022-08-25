package pattern

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
	"unicode/utf8"
)

/*
	Реализовать паттерн «фасад».
	Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern


	Структурный паттерн «фасад» удобен, когда нужно предоставить доступ к сложной системе или фреймворку.
	Использование больших систем напрямую может сопровождаться такими трудностями, как необходимость разбираться в порядке
	инициализации служебных объектов, архитектуре и устройстве всех сущностей, их взаимодействии между собой.
	В любом случае это приведет к высокой связности клиентского кода с кодом используемого пакета / библиотеки / фреймворка.

	Фасад применяется в большинстве фреймворков. Так же фасадом можно считать UI, API, CLI.

	+++

	Фасад предоставляет простой и удобный интерфейс для работы с подобными системами.
	Клиенский код не зависит от подробностей реализации подсистемы и легко может заменить ее на другую
	в случае необходимости.

	---

	Фасад может слишком сильно разрастись и превратиться в "божественный объект", который выполняет все и сразу.
	Поддержка подобных обьектов сложна и нелегка. "Божественный объект" является антипаттерном.


*/

/////

type NASANegotiator struct {
}

func (n NASANegotiator) GetAproveFromNASA() error {
	fmt.Println("NASA Negotiator: OK, they would not mind.")
	return nil
}

/////

type FuelManager struct {
}

func (f FuelManager) PrepareFuel(travelDistance uint64) {
	gallons := 1024 * travelDistance
	fmt.Printf("Fuel manager: %d gallons of fuel prepared for the launch.\n", gallons)
}

/////

type RoutePoint struct {
	P uint64
}

func NewRoutePoint(i uint64) RoutePoint {
	return RoutePoint{
		P: i,
	}
}

func (p *RoutePoint) GetPointInfo() string {
	return fmt.Sprint(p.P)
}

/////

type SystemConfigurator struct {
}

func (s SystemConfigurator) ConfigureAllSystems() {
	fmt.Println("System Configurator: All systems are configured and ready for launch.")
}

/////

type Route struct {
	Points []RoutePoint
}

func NewRoute() *Route {
	return &Route{
		Points: make([]RoutePoint, 1),
	}
}

func (r *Route) AddPoint(i uint64) {
	r.Points = append(r.Points, NewRoutePoint(i))
}

func (r *Route) GetNextPointInfo() (string, error) {
	if r.IsNotEmpty() {
		return r.Points[0].GetPointInfo(), nil
	}
	return "", errors.New("Route is empty")
}

func (r *Route) AdvanceRoute() {
	if r.IsEmpty() {
		fmt.Println("Destination has been reached already.")
		return
	}

	r.Points = r.Points[1:]
	if r.IsEmpty() {
		fmt.Println("You've reached the destination. Great job!")
		return
	}
}

func (r *Route) IsEmpty() bool {
	return len(r.Points) == 0
}

func (r *Route) IsNotEmpty() bool {
	return !r.IsEmpty()
}

func (r *Route) GetTotalDistance() (x uint64) {
	for _, v := range r.Points {
		x += v.P
	}
	return
}

/////

type Navigator struct {
}

func NewNavigator() Navigator {
	return Navigator{}
}

func (n Navigator) BuildRoute(from, to string) (*Route, error) {
	if from == "" || to == "" {
		return nil, errors.New("route ends are incorrect")
	}
	s := 6 * (utf8.RuneCountInString(from) - utf8.RuneCountInString(to))
	if s == 0 {
		s = rand.Intn(10) + 10
	}
	if s < 0 {
		s *= -1
	}

	route := NewRoute()

	for i := 1; i <= s; i++ {
		route.AddPoint(uint64(80 * i))
	}

	return route, nil
}

/////

type RocketLauncherFacade struct {
	Route *Route
}

func NewRocketLauncherFacade() *RocketLauncherFacade {
	return &RocketLauncherFacade{
		Route: NewRoute(),
	}
}

func (f *RocketLauncherFacade) Delay() {
	time.Sleep(650 * time.Millisecond)
}

func (f *RocketLauncherFacade) Launch(destination string, timeBeforeLaunch time.Duration) {
	fmt.Println("Facade: Negotiating with NASA...")
	nasa := NASANegotiator{}
	if err := nasa.GetAproveFromNASA(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Facade: Creating navigator...")
	navigator := NewNavigator()

	fmt.Println("Facade: Building new route...")
	f.Delay()
	currentLocation := f.GetCurrentLocation()
	r, err := navigator.BuildRoute(currentLocation, destination)
	if err != nil {
		log.Fatal(err)
	}

	f.Route = r

	fmt.Println("Facade: Calculating travel distance...")
	d := r.GetTotalDistance()
	f.Delay()
	fmt.Printf("Facade: Travel distance: %d\n", d)

	fmt.Println("Facade: Replenishing fuel reserves...")
	f.Delay()
	fm := FuelManager{}
	fm.PrepareFuel(d)

	fmt.Println("Facade: Configuring all systems...")
	f.Delay()
	sc := SystemConfigurator{}
	sc.ConfigureAllSystems()

	fmt.Println("Facade: Starting countdown...")
	f.Countdown(timeBeforeLaunch)

	fmt.Println("Facade: GO!")
	f.StartTravel()
}

func (f *RocketLauncherFacade) StartTravel() {
	for {
		if f.Route.IsEmpty() {
			return
		}
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Facade: Another travel point is passed...")
		f.Route.AdvanceRoute()
	}
}

func (f *RocketLauncherFacade) GetCurrentLocation() string {
	return "Milkiway, Earth"
}

func (f *RocketLauncherFacade) Countdown(d time.Duration) {
	sec := d.Seconds()
	for i := sec; i >= 0; i-- {
		fmt.Printf("%d\n", int(i))
		time.Sleep(time.Second)
	}
}

/////

func main1() {
	desination := "Andromeda Nebula"
	timeBeforeLaunch := 5 * time.Second
	launcher := NewRocketLauncherFacade()
	launcher.Launch(desination, timeBeforeLaunch)

	/*
		Output:

		Facade: Negotiating with NASA...
		NASA Negotiator: OK, they would not mind.
		Facade: Creating navigator...
		Facade: Building new route...
		Facade: Calculating travel distance...
		Facade: Travel distance: 1680
		Facade: Replenishing fuel reserves...
		Fuel manager: 1720320 gallons of fuel prepared for the launch.
		Facade: Configuring all systems...
		System Configurator: All systems are configured and ready for launch.
		Facade: Starting countdown...
		5
		4
		3
		2
		1
		0
		Facade: GO!
		Facade: Another travel point is passed...
		Facade: Another travel point is passed...
		Facade: Another travel point is passed...
		Facade: Another travel point is passed...
		Facade: Another travel point is passed...
		Facade: Another travel point is passed...
		Facade: Another travel point is passed...
		You've reached the destination. Great job!
	*/
}
