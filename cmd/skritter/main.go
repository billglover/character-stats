package main

import (
	"log"
	"os"
)

func main() {

	log := log.New(os.Stdout, "api : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	err := run(log)
	if err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {
	log.Println("Hello, world!")
	return nil
}
