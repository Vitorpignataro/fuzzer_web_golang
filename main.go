package main

import (
	"fuzzer/app"
	"log"
	"os"
)

func main() {
	application := app.Fuzzer()

	if erro := application.Run(os.Args); erro != nil{
		log.Fatal(erro)
	}
}