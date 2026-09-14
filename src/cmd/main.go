package main

import (
	"os"

	loadtest "github.com/RogerGCruz/pos-goexpert-desafio5/src"
)

func main() {
	code := loadtest.Run(os.Args[1:], os.Stdout)
	os.Exit(code)
}
