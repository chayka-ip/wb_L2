package main

import (
	"go-wget/wget"
	"os"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	os.Exit(wget.ExecuteCLI(os.Args[1:]))
}
