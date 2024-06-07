package app

import (
	"time"

	"github.com/integrii/flaggy"
	tb "github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"
)

var eventQueue chan tb.Event

var config struct {
	options struct {
		debug bool
	}
}

func Init() {
	flaggy.SetName("speedtype")
	flaggy.SetDescription("Typing test")

	flaggy.Bool(&config.options.debug, "d", "debug", "Enable debug mode")
	flaggy.Parse()

	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.TimeOnly,
	})
	if config.options.debug {
		log.SetLevel(log.DebugLevel)
	}

	tb.SetInputMode(tb.InputEsc)
	eventQueue = make(chan tb.Event)
	go func() {
		for {
			eventQueue <- tb.PollEvent()
		}
	}()
}

func Run() {
	render()
	for event := range eventQueue {
		handleEvent(event)
	}
}

func handleEvent(event tb.Event) {
	switch event.Type {
	case tb.EventKey:
		onKey(event)
	case tb.EventResize:
		render()
	case tb.EventError:
		log.Errorf("Error: %v", event.Err)
	}
}

func onKey(event tb.Event) {
	if shouldExit(event) {
		Exit()
	}
	update()
}

func shouldExit(event tb.Event) bool {
	key := event.Key
	if key == tb.KeyEsc || key == tb.KeyCtrlC {
		return true
	}
	return false
}

func Exit() {
	close(eventQueue)
	tb.Clear(tb.ColorDefault, tb.ColorDefault|tb.AttrBold)
}

func update() {

}

func render() {
	border()
	tb.Flush()
}

const (
	TopLeft     = '╭'
	TopRight    = '╮'
	BottomLeft  = '╰'
	BottomRight = '╯'
	Horizontal  = '─'
	Vertical    = '│'
)

func border() {
	width, height := getTerminalSize()
	fgColor := tb.ColorWhite
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				char := getBorderChar(x, y, width, height)
				tb.SetCell(x, y, char, fgColor, tb.ColorDefault)
			}
		}
	}
}

func getTerminalSize() (int, int) {
	tb.Sync()
	return tb.Size()
}

func getBorderChar(x, y, width, height int) rune {
	if x == 0 && y == 0 {
		return TopLeft
	}
	if x == width-1 && y == 0 {
		return TopRight
	}
	if x == 0 && y == height-1 {
		return BottomLeft
	}
	if x == width-1 && y == height-1 {
		return BottomRight
	}
	if x == 0 || x == width-1 {
		return Vertical
	}
	if y == 0 || y == height-1 {
		return Horizontal
	}
	return ' '
}
