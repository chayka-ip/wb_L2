package pattern

import "fmt"

/*
	Реализовать паттерн «команда».
	Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern


	Поведенческий паттерн «команда» позволяет превратить запрос или вызов метода в обьект,
	с которым можно производить различные полезные операции, такие как передача в качестве
	параметра, логирование, копирование, хранение, поддержка отмены операций.

	Команда повсеместно применяется при создании интерфейсов взаимодействия пользователя
	с программой (UI, CLI, API, Hot Keys), при создании макросов, инструментов автоматизации,
	удаленном вызове процедур (RPC), оценке времени выполнения в progress bar,
	создании транзакций и тд.

	+++

	Паттерн помогает создать дополнительный уровень абстракции между точкой вызова
	и принимающей стороной. Вся необходимая информация для выполнения команды
	инкапсулируется	внутри нее. Команда может вызываться из различных компонентов
	или участков кода, которые никак не связаны между собой.

	Паттерн позволяет создавать последовательности команд или комплексные команды из более простых.

	Работа с командами как с обьектами дает полный контроль над условиями их выполнения,
	позволяя откладывать его, отменять, сериализовывать, передавать по сети и т.д.

	---

	Расширяет кодовую базу ввиду добавления новых типов.


*/

// Command interface
type Command interface {
	Execute()
}

// Concrete command
type SaveCommand struct {
	state State
}

func (c *SaveCommand) Execute() {
	c.state.Save()
}

// Concrete command
type LoadCommand struct {
	state State
}

func (c *LoadCommand) Execute() {
	c.state.Load()
}

// Concrete command
type ResetCommand struct {
	state State
}

func (c *ResetCommand) Execute() {
	c.state.Reset()
}

// Sender
type Button struct {
	command Command
}

func (b *Button) Press() {
	fmt.Println("Button is pressed.")
	b.command.Execute()
}

// Sender
type Shortcut struct {
	command Command
}

func (s *Shortcut) Use() {
	fmt.Println("Shortcut is used")
	s.command.Execute()
}

// Business logic interface
type State interface {
	Save()
	Load()
	Reset()
}

// Receiver
type Progress struct {
}

func (p Progress) Save() {
	fmt.Println("Progress is saved.")
}

func (p Progress) Load() {
	fmt.Println("Progress is loaded.")
}

func (p Progress) Reset() {
	fmt.Println("Progress is reset.")
}

func main4() {
	p := Progress{}

	save := SaveCommand{state: p}
	load := LoadCommand{state: p}
	reset := ResetCommand{state: p}

	saveButton := &Button{command: &save}
	loadButton := &Button{command: &load}
	resetButton := &Button{command: &reset}

	saveShortcut := &Shortcut{command: &save}
	loadShortcut := &Shortcut{command: &load}
	resetShortcut := &Shortcut{command: &reset}

	saveButton.Press()
	resetButton.Press()
	loadButton.Press()

	saveShortcut.Use()
	resetShortcut.Use()
	loadShortcut.Use()

	/*

		Output:

		Button is pressed.
		Progress is saved.

		Button is pressed.
		Progress is reset.

		Button is pressed.
		Progress is loaded.


		Shortcut is used
		Progress is saved.

		Shortcut is used
		Progress is reset.

		Shortcut is used
		Progress is loaded.

	*/
}
