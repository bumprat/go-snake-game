package main

import (
	"log"
	"os"

	"github.com/bumprat/go-snakes-game/snakes"
)

func main() {
	logger := log.New(
		os.Stdout,
		"snakes: ",
		log.Lmicroseconds|log.Lshortfile,
	)

	err := snakes.NewDefault()
	if err != nil {
		logger.Fatalln(err.Error())
	}
	snakes.Start()
}
