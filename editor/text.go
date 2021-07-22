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

	// Can't remove the beginning of the file
	if pos.x <= 0 && pos.y <= 0 {
		state.text[0] = state.text[0][1:]
		return
	}

	// If the user deletes a line.
	if state.text[pos.y][pos.x] == '\n' {
		// Copy the contents of the previous line (excluding the '\n')
		if state.text[pos.y+1][0] == '\n' {
			state.text[pos.y] += state.text[pos.y+1][1:]
		} else {
			state.text[pos.y] += state.text[pos.y+1]
		}

		// Remove the deleted line.
		state.text = append(state.text[:pos.y+1], state.text[pos.y+2:]...)
	} else {
		state.text[pos.y] = state.text[pos.y][:pos.x] + state.text[pos.y][pos.x+1:]
	}

	return nil
}
