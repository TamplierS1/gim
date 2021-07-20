package main

// Find the position of the last character on a 'line'.
func find_last_char(state *EditorState, line int) int {
	for x := len(state.text[line]) - 1; x > 0; x-- {
		if state.text[line][x] != ' ' {
			return x
		}

	}
	return 0
}

// Find the position of the first character on a 'line'.
func find_first_char(state *EditorState, line int) int {
	for x := range state.text[line] {
		if state.text[line][x] != ' ' {
			return x
		}
	}

	return 0
}

// Find the position of the next character starting from `start_from`
func find_next_char(state *EditorState, start_from Vec2) (pos Vec2, e error) {
	// The search starts from the beginning of the line
	// except for the first line.
	is_first_line := true

	for y := range state.text[start_from.y:] {
		for x := range state.text[y+start_from.y] {
			if is_first_line {
				x = start_from.x + x
			}

			if x >= len(state.text[y+start_from.y]) {
				break
			}

			r := state.text[y+start_from.y][x]
			if r != ' ' && r != '\n' {
				return Vec2{x, y + start_from.y}, nil
			}
		}
		is_first_line = false
	}

	return Vec2{-1, -1}, Error{"Error: failed to find the next char."}
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

		// putchar_at_cursor(state, ' ')
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
// func putchar_at_cursor(state *EditorState, char rune) {
// 	PutChar(state.screen, state.cursor.x, state.cursor.y, state.default_style, char)
// }

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
