// Simple console text editor
package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type Vec2 struct {
	x, y int
}

var (
	g_cursor_pos    = Vec2{0, 0}
	g_default_style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
)

// Find the position of the last character on a 'line'.
func find_last_char(screen tcell.Screen, line int) (x int) {
	screen_width, _ := screen.Size()
	for i := screen_width - 1; i > 0; i-- {
		char, _, _, _ := screen.GetContent(i, line)

		if char != ' ' {
			// I return i + 1 here because screen.GetContent
			// doesn't return '\n' as characters.
			return i + 1
		}
	}
	return 0
}

// Erases the word behind the cursor.
func erase_prev_word(screen tcell.Screen) {
	word_start := g_cursor_pos.x
	// Find the start of the word
	for ; word_start >= 0; word_start-- {
		char, _, _, _ := screen.GetContent(word_start, g_cursor_pos.y)

		if char != ' ' {
			break
		}

		move_back_cursor(screen)
	}

	for i := word_start; i >= 0; i-- {
		char, _, _, _ := screen.GetContent(i, g_cursor_pos.y)

		if char == ' ' {
			break
		}

		putchar_at_cursor(screen, ' ')
		move_back_cursor(screen)
	}
}

// Process events.
func handle_events(screen tcell.Screen) {
	event := screen.PollEvent()

	switch event := event.(type) {
	case *tcell.EventResize:
		screen.Sync()
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyESC: // Close the editor
			screen.Fini()
			os.Exit(0)
		case tcell.KeyBackspace2: // Erase characters
			if event.Modifiers() == tcell.ModAlt {
				erase_prev_word(screen)
			} else {
				move_back_cursor(screen)
				putchar_at_cursor(screen, ' ')
			}
		case tcell.KeyEnter: // Move the cursor to the next line
			move_down_cursor(screen)
			g_cursor_pos.x = 0
		case tcell.KeyRight: // Move the cursor
			move_right_cursor(screen)
		case tcell.KeyLeft:
			move_back_cursor(screen)
		case tcell.KeyUp:
			move_up_cursor(screen)
		case tcell.KeyDown:
			move_down_cursor(screen)
		default: // Enter text
			putchar_at_cursor(screen, event.Rune())
			move_right_cursor(screen)
		}
	}
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	screen.SetStyle(g_default_style)
	screen.EnableMouse()
	screen.EnablePaste()
	screen.ShowCursor(0, 0)

	screen.Clear()

	for {
		handle_events(screen)
		update_cursor(screen)

		screen.Show()
	}
}
