package app

import (
	"github.com/integrii/flaggy"
	"github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"
)

var eventQueue chan termbox.Event

var config struct {
	options struct {
		debug bool
	}
}

func Init() {
	flaggy.SetName("speedtype")
	flaggy.SetDescription("Typing test")
	flaggy.SetVersion("0.1.0")

	flaggy.Bool(&config.options.debug, "d", "debug", "Enable debug mode")
	flaggy.Parse()

	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
	if config.options.debug {
		log.SetLevel(log.DebugLevel)
	}

	termbox.SetCursor(0, 0)
	eventQueue = make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
}

func Run() {
	for event := range eventQueue {
		handleEvent(event)
		update()
		render()
	}

}

func handleEvent(ev termbox.Event) {
	switch ev.Type {
	case termbox.EventKey:
		log.Debugf("Key: %v, Ch: %v", ev.Key, ev.Ch)
		if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
			close(eventQueue)
			return
		}
	case termbox.EventError:
		log.Fatalf("Error: %v", ev.Err)
	}
}

func update() {
}

func render() {
}

func Exit() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorBlack)
	termbox.Flush()
}
