package main

import (
	"log"

	"preselect/api"
)

func main() {
	app := api.New(api.Config{})
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
