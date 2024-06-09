package app

import (
	tb "github.com/nsf/termbox-go"
)

const (
	ColorBgMain  = tb.ColorDefault
	ColorFgMain  = tb.ColorWhite
	ColorFgWords = tb.ColorLightGray
)

const (
	BorderTopLeft     = '╭'
	BorderTopRight    = '╮'
	BorderBottomLeft  = '╰'
	BorderBottomRight = '╯'
	BorderHorizontal  = '─'
	BorderVertical    = '│'
)

const borderSize = 1
const padding = 2

var cursor int

func GetTerminalSize() (int, int) {
	tb.Sync()
	return tb.Size()
}

func GetCursorPosition() (int, int) {
	width, _ := GetTerminalSize()
	cursorX := cursor % (width - 2*borderSize - 2*padding)
	cursorY := cursor / (width - 2*borderSize - 2*padding)
	return cursorX + padding + borderSize, cursorY + padding + borderSize
}

func GetBorderChar(x, y, width, height int) rune {
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

func Render() {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	printBorder()

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

func printBorder() {
	width, height := GetTerminalSize()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x < borderSize || x >= width-borderSize || y < borderSize || y >= height-borderSize {
				char := GetBorderChar(x, y, width, height)
				tb.SetCell(x, y, char, ColorFgMain, ColorBgMain)
			}
		}
	}
}

func printMenu() {
	width, height := GetTerminalSize()
	text := "Press Space to start"
	textWidth := len(text)
	x := (width - textWidth) / 2
	y := height / 2
	for i, char := range text {
		tb.SetCell(x+i, y, char, ColorFgMain, ColorBgMain)
	}
}

func printGame() {
	printStringBuffer(targetBuffer, ColorFgWords)
	printCellBuffer(inputBuffer)
}

func printStringBuffer(buffer string, fg tb.Attribute) {
	terminalWidth, terminalHeight := GetTerminalSize()
	offset := borderSize + padding
	areaWidth, areaHeight := terminalWidth-2*offset, terminalHeight-2*offset
	bufferWidth := len(buffer) / areaHeight
	bufferStartWidth := offset + (areaWidth-bufferWidth)/2
	for x := offset; x < areaWidth; x++ {
		for y := offset; y < areaHeight; y++ {
			index := y - offset + (x-offset)*(areaHeight)
			if index < len(buffer) {
				char := buffer[index]
				tb.SetCell(x+bufferStartWidth, y, rune(char), fg, ColorBgMain)
			}
		}
	}
}

func printCellBuffer(buffer []tb.Cell) {
	terminalWidth, terminalHeight := GetTerminalSize()
	offset := borderSize + padding
	areaWidth, areaHeight := terminalWidth-2*offset, terminalHeight-2*offset
	bufferHeight := len(buffer) / areaWidth
	bufferStartHeight := offset + (areaHeight-bufferHeight)/2
	for x := offset; x < areaWidth; x++ {
		for y := offset; y < areaHeight; y++ {
			index := x - offset + (y-offset)*(areaWidth)
			if index < len(buffer) {
				cell := buffer[index]
				tb.SetCell(x, y+bufferStartHeight, cell.Ch, cell.Fg, cell.Bg)
			}
		}
	}
}

func printResult() {
	width, height := GetTerminalSize()
	text := "Press Space to restart"
	textWidth := len(text)
	x := (width - textWidth) / 2
	y := height / 2
	for i, char := range text {
		tb.SetCell(x+i, y, char, ColorFgMain, ColorBgMain)
	}
}
