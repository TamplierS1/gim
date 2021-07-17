package main

import (
	"github.com/gdamore/tcell/v2"
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
func erase_prev_word(state *EditorState) {
	word_start := state.cursor.x
	// Find the start of the word
	for ; word_start >= 0; word_start-- {
		char, _, _, _ := state.screen.GetContent(word_start, state.cursor.y)

		if char != ' ' {
			break
		}

		move_back_cursor(state)
	}

	// Erase the word
	for i := word_start; i >= 0; i-- {
		char, _, _, _ := state.screen.GetContent(i, state.cursor.y)

		if char == ' ' {
			break
		}

		putchar_at_cursor(state, ' ')
		remove_symbol(state, state.cursor)
		move_back_cursor(state)
	}
}

// Adds a typed symbol into the structure that holds the text
// of the current file
func record_rune(state *EditorState, pos Vec2, r rune) {
	if pos.y >= len(state.text) {
		state.text = append(state.text, string(r))
		return
	}
	if pos.x >= len(state.text[pos.y]) {
		state.text[pos.y] = state.text[pos.y] + string(r)
		return
	}

	state.text[pos.y] = state.text[pos.y][:pos.x] + string(r) + state.text[pos.y][pos.x:]
}

// Remove the erased symbol from the `state.text`
func remove_symbol(state *EditorState, pos Vec2) (e error) {
	if pos.y >= len(state.text) ||
		pos.x >= len(state.text[pos.y]) {
		return Error{"Error: the symbol's index is out of bounds. Can not remove it" +
			"\n"}
	}

	state.text[pos.y] = state.text[pos.y][:pos.x] + state.text[pos.y][pos.x+1:]
	return nil
}

// Set a cell at the cursor's position to 'char'
func putchar_at_cursor(state *EditorState, char rune) {
	PutChar(state.screen, state.cursor.x, state.cursor.y, state.default_style, char)
}

// Moves all characters on line 'num' one cell to the right
// starting from 'index' on the screen
// func shift_line_to_right(state *EditorState) (e error) {
// 	prev_rune := ' '
// 	for i, r := range state.text[state.cursor.y][state.cursor.x:] {
// 		pos := Vec2{state.cursor.x + i, state.cursor.y}
// 		putchar_at(state, pos, prev_rune)
// 		prev_rune = r
// 	}
// 	return nil
// }

// Set a cell at the position 'pos' to 'char'
// func putchar_at(state *EditorState, pos Vec2, char rune) {
// 	PutChar(state.screen, pos.x, pos.y, state.default_style, char)
// }
