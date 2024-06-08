package app

import (
	"time"

	"github.com/integrii/flaggy"
	tb "github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"
)

var config struct {
	options struct {
		debug bool
	}
}
var eventQueue chan tb.Event
var buffer string

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
		onResize(event)
	case tb.EventError:
		log.Errorf("Error: %v", event.Err)
	}
}

func onKey(event tb.Event) {
	if shouldExit(event) {
		Exit()
	}
	update(event)
}

func onResize(_ tb.Event) {
	render()
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

func update(event tb.Event) {
	if event.Ch != 0 {
		buffer += string(event.Ch)
	}
	render()
}

func render() {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	border()
	printBuffer()
	tb.Flush()
}

const (
	borderFgColor = tb.ColorWhite
	borderBgColor = tb.ColorDefault
	bufferFgColor = tb.ColorWhite
	bufferBgColor = tb.ColorDefault
)

const (
	TopLeft     = '╭'
	TopRight    = '╮'
	BottomLeft  = '╰'
	BottomRight = '╯'
	Horizontal  = '─'
	Vertical    = '│'
)

var borderSize = 1
var padding = 1

func border() {
	width, height := getTerminalSize()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x < borderSize || x >= width-borderSize || y < borderSize || y >= height-borderSize {
				char := getBorderChar(x, y, width, height)
				tb.SetCell(x, y, char, borderFgColor, borderBgColor)
			}
		}
	}
}

func getTerminalSize() (int, int) {
	tb.Sync()
	return tb.Size()
}

func getBorderChar(x, y, width, height int) rune {
	if x < borderSize && y < borderSize {
		return TopLeft
	}
	if x >= width-borderSize && y < borderSize {
		return TopRight
	}
	if x < borderSize && y >= height-borderSize {
		return BottomLeft
	}
	if x >= width-borderSize && y >= height-borderSize {
		return BottomRight
	}
	if x < borderSize || x >= width-borderSize {
		return Vertical
	}
	if y < borderSize || y >= height-borderSize {
		return Horizontal
	}
	return ' '
}

func printBuffer() {
	width, height := getTerminalSize()
	offset := borderSize + padding
	for x := offset; x < width-offset; x++ {
		for y := offset; y < height-offset; y++ {
			index := x - offset + (y-offset)*(width-2*offset)
			if index < len(buffer) {
				tb.SetCell(x, y, rune(buffer[index]), bufferFgColor, bufferBgColor)
			}
		}
	}
}
