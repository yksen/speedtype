package app

import (
	"time"

	"github.com/integrii/flaggy"
	tb "github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"
)

var config struct {
	debug bool
}

var eventQueue chan tb.Event

func Init() {
	flaggy.SetName("speedtype")
	flaggy.SetDescription("Typing test")

	flaggy.Bool(&config.debug, "d", "debug", "Enable debug mode")
	flaggy.Parse()

	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.TimeOnly,
	})
	if config.debug {
		log.SetLevel(log.DebugLevel)
	}

	eventQueue = make(chan tb.Event)
	go func() {
		for {
			eventQueue <- tb.PollEvent()
		}
	}()
}

const (
	StateMenu = iota
	StateGame
	StateResults
)

var state = StateMenu
var cursor int

func Run() {
	tb.SetCursor(padding+borderSize, padding+borderSize)
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
		onResize()
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

func onResize() {
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

var buffer string

const (
	ActionAdd = iota
	ActionRemove
	ActionNone
)

func update(event tb.Event) {
	action := ActionNone
	if event.Ch != 0 {
		buffer += string(event.Ch)
		action = ActionAdd
	} else {
		switch event.Key {
		case tb.KeySpace:
			buffer += " "
			action = ActionAdd
		case tb.KeyBackspace, tb.KeyBackspace2:
			if len(buffer) > 0 {
				buffer = buffer[:len(buffer)-1]
				action = ActionRemove
			}
		}
	}

	switch action {
	case ActionAdd:
		cursorX, cursorY := getCursorPosition()
		tb.SetChar(cursorX, cursorY, rune(buffer[cursor]))
		cursor++
	case ActionRemove:
		cursor--
		cursorX, cursorY := getCursorPosition()
		tb.SetChar(cursorX, cursorY, ' ')
	}

	cursorX, cursorY := getCursorPosition()
	tb.SetCursor(cursorX, cursorY)
	tb.Flush()
}

func getCursorPosition() (int, int) {
	width, _ := getTerminalSize()
	cursorX := cursor % (width - 2*borderSize - 2*padding)
	cursorY := cursor / (width - 2*borderSize - 2*padding)
	return cursorX + padding + borderSize, cursorY + padding + borderSize
}

func render() {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	printBorder()
	printBuffer()
	tb.Flush()
}

const (
	ColorBgBorder = tb.ColorDefault
	ColorFgBorder = tb.ColorWhite
	ColorBgBuffer = tb.ColorDefault
	ColorFgBuffer = tb.ColorWhite
)

var borderSize = 1
var padding = 1

func printBorder() {
	width, height := getTerminalSize()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x < borderSize || x >= width-borderSize || y < borderSize || y >= height-borderSize {
				char := getBorderChar(x, y, width, height)
				tb.SetCell(x, y, char, ColorFgBorder, ColorBgBorder)
			}
		}
	}
}

func getTerminalSize() (int, int) {
	tb.Sync()
	return tb.Size()
}

const (
	BorderTopLeft     = '╭'
	BorderTopRight    = '╮'
	BorderBottomLeft  = '╰'
	BorderBottomRight = '╯'
	BorderHorizontal  = '─'
	BorderVertical    = '│'
)

func getBorderChar(x, y, width, height int) rune {
	if x < borderSize && y < borderSize {
		return BorderTopLeft
	}
	if x >= width-borderSize && y < borderSize {
		return BorderTopRight
	}
	if x < borderSize && y >= height-borderSize {
		return BorderBottomLeft
	}
	if x >= width-borderSize && y >= height-borderSize {
		return BorderBottomRight
	}
	if x < borderSize || x >= width-borderSize {
		return BorderVertical
	}
	if y < borderSize || y >= height-borderSize {
		return BorderHorizontal
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
				tb.SetCell(x, y, rune(buffer[index]), ColorFgBuffer, ColorBgBuffer)
			}
		}
	}
}
