package app

import (
	"math/rand"
	"time"

	"github.com/integrii/flaggy"
	tb "github.com/nsf/termbox-go"
)

const (
	ActionNone = iota
	ActionAdd
	ActionRemove
)

const (
	StateNone = iota
	StateMenu
	StatePregame
	StateGame
	StateResult
)

var Config struct {
	Debug bool
}

var charsTyped int
var eventQueue chan tb.Event
var state = StateNone
var timeRemaining time.Duration
var words = []string{"the", "be", "of", "and", "a", "to", "in", "he", "have", "it", "that", "for", "they", "I", "with", "as", "not", "on", "she", "at", "by", "this", "we", "you", "do", "but", "from", "or", "which", "one", "would", "all", "will", "there", "say", "who", "make", "when", "can", "more", "if", "no", "man", "out", "other", "so", "what", "time", "up", "go", "about", "than", "into", "could", "state", "only", "new", "year", "some", "take", "come", "these", "know", "see", "use", "get", "like", "then", "first", "any", "work", "now", "may", "such", "give", "over", "think", "most", "even", "find", "day", "also", "after", "way", "many", "must", "look", "before", "great", "back", "through", "long", "where", "much", "should", "well", "people", "down", "own", "just", "because", "good", "each", "those", "feel", "seem", "how", "high", "too", "place", "little", "world", "very", "still", "nation", "hand", "old", "life", "tell", "write", "become", "here", "show", "house", "both", "between", "need", "mean", "call", "develop", "under", "last", "right", "move", "thing", "general", "school", "never", "same", "another", "begin", "while", "number", "part", "turn", "real", "leave", "might", "want", "point", "form", "off", "child", "few", "small", "since", "against", "ask", "late", "home", "interest", "large", "person", "end", "open", "public", "follow", "during", "present", "without", "again", "hold", "govern", "around", "possible", "head", "consider", "word", "program", "problem", "however", "lead", "system", "set", "order", "eye", "plan", "run", "keep", "face", "fact", "group", "play", "stand", "increase", "early", "course", "change", "help", "line"}

func Init() {
	flaggy.SetName("speedtype")
	flaggy.Bool(&Config.Debug, "d", "debug", "Enable debug mode")
	flaggy.Parse()

	UpdateTerminalSize()
	eventQueue = make(chan tb.Event)
	go func() {
		for {
			eventQueue <- tb.PollEvent()
		}
	}()
	changeState(StateMenu)
}

func Run() {
	for event := range eventQueue {
		if Config.Debug {
			printDebug()
		}
		handleEvent(event)
	}
}

func Exit() {
	close(eventQueue)
	tb.Clear(tb.ColorDefault, tb.ColorDefault|tb.AttrBold)
}

func changeState(newState int) {
	switch newState {
	case StatePregame:
		generateRandomWords()
		inputBuffer = make([]tb.Cell, 0)
		cursor = 0
		charsTyped = 0
	case StateGame:
		timeRemaining = 15 * time.Second
		ticker := time.NewTicker(1 * time.Second)
		go func() {
			for range ticker.C {
				timeRemaining -= 1 * time.Second
				if Config.Debug {
					printDebug()
				}
				if timeRemaining <= 0 {
					changeState(StateResult)
					break
				}
			}
		}()
	case StateResult:
	}

	state = newState
	Render()
	if Config.Debug {
		printDebug()
	}
}

func handleEvent(event tb.Event) {
	switch event.Type {
	case tb.EventKey:
		onKey(event)
	case tb.EventResize:
		onResize()
	}
}

func onKey(event tb.Event) {
	if shouldExit(event) {
		Exit()
	}
	update(event)
}

func onResize() {
	UpdateTerminalSize()
	Render()
}

func shouldExit(event tb.Event) bool {
	key := event.Key
	if key == tb.KeyEsc || key == tb.KeyCtrlC {
		return true
	}
	return false
}

func update(event tb.Event) {
	switch state {
	case StateMenu:
		if event.Key == tb.KeySpace {
			changeState(StatePregame)
		}
	case StatePregame:
		changeState(StateGame)
	case StateGame:
		updateGame(event)
	case StateResult:
		if event.Key == tb.KeySpace {
			changeState(StateMenu)
		}
	}
}

func updateGame(event tb.Event) {
	action := ActionNone
	if event.Ch != 0 {
		inputBuffer = append(inputBuffer, tb.Cell{Ch: event.Ch, Fg: tb.ColorDefault, Bg: tb.ColorDefault})
		action = ActionAdd
	} else {
		switch event.Key {
		case tb.KeySpace:
			inputBuffer = append(inputBuffer, tb.Cell{Ch: ' ', Fg: tb.ColorDefault, Bg: tb.ColorDefault})
			action = ActionAdd
		case tb.KeyBackspace, tb.KeyBackspace2:
			if len(inputBuffer) > 0 {
				inputBuffer = inputBuffer[:len(inputBuffer)-1]
				action = ActionRemove
			}
		}
	}

	switch action {
	case ActionAdd:
		cursorX, cursorY := GetCursorPosition(Area)
		tb.SetChar(cursorX, cursorY, inputBuffer[len(inputBuffer)-1].Ch)
		cursor++
	case ActionRemove:
		cursor--
		cursorX, cursorY := GetCursorPosition(Area)
		tb.SetChar(cursorX, cursorY, ' ')
	}

	cursorX, cursorY := GetCursorPosition(Area)
	tb.SetCursor(cursorX, cursorY)
	tb.Flush()
}

func generateRandomWords() {
	for i := 0; i < 1000; i++ {
		targetBuffer += words[rand.Intn(len(words))] + " "
	}
	targetBuffer = targetBuffer[:len(targetBuffer)-1]
}
