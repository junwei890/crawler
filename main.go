package main

import (
	"log"

	"github.com/junwei890/se-cli/src"
)

func main() {
	if err := src.Init(); err != nil {
		log.Fatal(err)
	}
}
