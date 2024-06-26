package main

import (
	"github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"
	app "github.com/yksen/speedtype/app"
)

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	app.Init()
	app.Run()
}
