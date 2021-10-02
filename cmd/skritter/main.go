package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/billglover/character-stats/skritter"
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

	var token string
	flag.StringVar(&token, "token", "", "Skritter API token")
	flag.Parse()

	if token == "" {
		return errors.New("must provide --token")
	}

	client := skritter.NewClient(token)
	err := client.Items()
	if err != nil {
		return err
	}

	return nil
}
