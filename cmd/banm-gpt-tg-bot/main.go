package main

import (
	"log"

	"github.com/beardedandnotmuch/gpt-tg-bot/internal/pkg/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
