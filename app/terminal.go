package app

import (
	"fmt"

	tb "github.com/nsf/termbox-go"
)

// Color constants
const (
	ColorBgMain  = tb.ColorDefault
	ColorFgMain  = tb.ColorWhite
	ColorFgWords = tb.ColorLightGray
)

// Border characters
const (
	BorderTopLeft     = '╭'
	BorderTopRight    = '╮'
	BorderBottomLeft  = '╰'
	BorderBottomRight = '╯'
	BorderHorizontal  = '─'
	BorderVertical    = '│'
)

// Layout constants
const (
	borderSize        = 1
	paddingHorizontal = 3
	paddingVertical   = 1
	paddingRatio      = 3
)

type Rect struct {
	X, Y, Width, Height int
}

/*
Terminal layout:
- Full: Full terminal size
- Body: Terminal size without borders
- Area: Terminal size without borders and padding

	╭Full────────╮
	│Body	     │
	│    Area    │
	│            │
	╰────────────╯
*/
var Full, Body, Area Rect
var Debug Rect

var cursor int
var inputBuffer []tb.Cell
var targetBuffer string

func UpdateTerminalSize() {
	tb.Sync()
	FullWidth, FullHeight := tb.Size()
	Full = Rect{0, 0, FullWidth, FullHeight}
	Body = Rect{borderSize, borderSize, FullWidth - 2*borderSize, FullHeight - 2*borderSize}
	Area = Rect{borderSize + paddingHorizontal, borderSize + paddingVertical,
		FullWidth - 2*borderSize - 2*paddingHorizontal,
		FullHeight - 2*borderSize - 2*paddingVertical,
	}
	Debug = Rect{0, 0, FullWidth, 2}
}

func GetCursorPosition(rect Rect) (int, int) {
	return rect.X + cursor%rect.Width, rect.Y + cursor/rect.Width
}

func GetBorderChar(x, y int, rect Rect) rune {
	switch {
	case x == rect.X && y == rect.Y:
		return BorderTopLeft
	case x == rect.X+rect.Width-1 && y == rect.Y:
		return BorderTopRight
	case x == rect.X && y == rect.Y+rect.Height-1:
		return BorderBottomLeft
	case x == rect.X+rect.Width-1 && y == rect.Y+rect.Height-1:
		return BorderBottomRight
	case x == rect.X || x == rect.X+rect.Width-1:
		return BorderVertical
	case y == rect.Y || y == rect.Y+rect.Height-1:
		return BorderHorizontal
	}
	return ' '
}

func Render() {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	printBorder(Full)

	switch state {
	case StateMenu:
		printMenu()
	case StateGame:
		printGame()
	case StateResult:
		printResult()
	}

	tb.Flush()
}

func printBorder(rect Rect) {
	width, height := rect.Width, rect.Height
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x < borderSize || x >= width-borderSize || y < borderSize || y >= height-borderSize {
				char := GetBorderChar(x, y, rect)
				tb.SetCell(x, y, char, ColorFgMain, ColorBgMain)
			}
		}
	}
}

func printMenu() {
	width, height := Area.Width, Area.Height
	text := "Press Space to start"
	textWidth := len(text)
	x := (width - textWidth) / 2
	y := height / 2
	for i, char := range text {
		tb.SetCell(x+i, y, char, ColorFgMain, ColorBgMain)
	}
}

func printGame() {
	printBuffer(Area, targetBuffer, ColorFgWords, ColorBgMain)
	printCellBuffer(Area, inputBuffer)
}

func printResult() {
	width, height := Area.Width, Area.Height
	text := "Press Space to restart"
	textWidth := len(text)
	x := (width - textWidth) / 2
	y := height / 2
	for i, char := range text {
		tb.SetCell(x+i, y, char, ColorFgMain, ColorBgMain)
	}
}

func printBuffer(rect Rect, buffer string, fg tb.Attribute, bg tb.Attribute) {
	cellBuffer := make([]tb.Cell, len(buffer))
	for i, char := range buffer {
		cellBuffer[i] = tb.Cell{Ch: char, Fg: fg, Bg: bg}
	}
	printCellBuffer(rect, cellBuffer)
}

func printCellBuffer(rect Rect, buffer []tb.Cell) {
	for x := rect.X; x < rect.X+rect.Width; x++ {
		for y := rect.Y; y < rect.Y+rect.Height; y++ {
			index := x - rect.X + (y-rect.Y)*rect.Width
			if index < len(buffer) {
				cell := buffer[index]
				tb.SetCell(x, y, cell.Ch, cell.Fg, cell.Bg)
			}
		}
	}
}

func printDebug() {
	buffer := fmt.Sprintf("state: %d, cursor: %d, timeRemaining: %02d", state, cursor, int(timeRemaining.Seconds()))
	printBuffer(Debug, buffer, tb.ColorCyan, ColorBgMain)
	tb.Flush()
}
