package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bumprat/go-snakes-game/snakes"
)

func main() {
	logger := log.New(
		os.Stdout,
		"snakes: ",
		log.Lmicroseconds|log.Lshortfile,
	)

	var width, height string
	fmt.Print("Input game stage width (>=5):")
	fmt.Scanln(&width)
	w, err := strconv.Atoi(width)
	if err != nil {
		logger.Fatalln(err.Error())
	}

	fmt.Print("Input game stage height (>=5):")
	fmt.Scanln(&height)
	h, err := strconv.Atoi(height)
	if err != nil {
		logger.Fatalln(err.Error())
	}

	fmt.Println(w, h)
	er := snakes.Init(byte(w), byte(h))
	if er != nil {
		logger.Fatalln(er.Error())
	}

	snakes.Start()
}
