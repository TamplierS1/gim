// Simple console text editor
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type Vec2 struct {
	x, y int
}

type EditorState struct {
	cursor        Vec2
	default_style tcell.Style
	screen        tcell.Screen
	filename      string
	actions       []Action
	text          []string
	// The line from which to start drawing text,
	// it's needed for scrolling.
	upper_bound int
}

func add_action(state *EditorState, action Action) {
	state.actions = append(state.actions, action)
}

func execute_actions(state *EditorState) {
	for _, action := range state.actions {
		action.execute(state)
		state.actions = state.actions[1:]
	}
}

// Draws the recorded text to the screen,
// starting from the line number `upper_bound`
func show_text(state *EditorState) {
	for y := range state.text[state.upper_bound:] {
		for x, r := range state.text[y+state.upper_bound] {
			PutChar(state.screen, x, y, state.default_style, r)
		}
	}
}

var g_last_key_pressed tcell.Key

func handle_events(state *EditorState) {
	event := state.screen.PollEvent()

	switch event := event.(type) {
	case *tcell.EventResize:
		state.screen.Sync()
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyESC:
			add_action(state, CloseEditorAction{})
		case tcell.KeyBackspace2:
			if event.Modifiers() == tcell.ModAlt {
				add_action(state, EraseWordAction{})
			} else {
				add_action(state, EraseCharAction{})
			}
		case tcell.KeyEnter:
			add_action(state, NewLineAction{})
		case tcell.KeyRight:
			add_action(state, MoveCursorAction{MoveRight})
		case tcell.KeyLeft:
			add_action(state, MoveCursorAction{MoveLeft})
		case tcell.KeyUp:
			add_action(state, MoveCursorAction{MoveUp})
		case tcell.KeyDown:
			add_action(state, MoveCursorAction{MoveDown})
		case tcell.KeyCtrlZ:
			if g_last_key_pressed == tcell.KeyCtrlZ {
				add_action(state, SaveAndCloseEditorAction{})
			}
		default:
			add_action(state, EnterCharAction{event.Rune()})
		}
		g_last_key_pressed = event.Key()
	}
}

func parse_args(state *EditorState) {
	switch len(os.Args) {
	case 1:
		fmt.Printf("usage: gim <filename>\n")
		exit(0, state.screen)
	case 2:
		state.filename = os.Args[1]
		state.text = read_file(state.filename, state.screen)
	default:
		state.filename = os.Args[1]
		state.text = read_file(state.filename, state.screen)
		fmt.Println("Warning: too many arguments.")
	}
}

func read_file(filename string, screen tcell.Screen) []string {
	file, err := os.Open(filename)
	if err != nil {
		file, err = os.Create(filename)
		if err != nil {
			fatal_exit(1, screen, err)
		}
	}

	reader := bufio.NewReader(file)

	text := make([]string, 0)
	for {
		line, err := reader.ReadString('\n')
		text = append(text, line)
		if err != nil { // EOF
			// log.Print(err)
			file.Close()
			return text
		}
	}
}

func fatal_exit(code int, screen tcell.Screen, err error) {
	screen.Fini()
	log.Fatal(err)
	os.Exit(code)
}

func exit(code int, screen tcell.Screen) {
	screen.Fini()
	os.Exit(code)
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	state := EditorState{
		default_style: tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset),
		screen:        screen,
		actions:       make([]Action, 0),
		text:          make([]string, 0),
		upper_bound:   0,
	}

	parse_args(&state)

	screen.SetStyle(state.default_style)
	screen.EnableMouse()
	screen.EnablePaste()
	screen.ShowCursor(0, 0)

	screen.Clear()

	for {
		screen.Clear()

		handle_events(&state)
		execute_actions(&state)
		update_cursor(&state)

		show_text(&state)
		screen.Show()

	}
}
