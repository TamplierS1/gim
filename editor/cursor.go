package main

import "log"

// Scroll the screen down
func scroll_down(state *EditorState) {
	if state.upper_bound >= len(state.text) {
		return
	}

	state.upper_bound++
}

// Scroll the screen up
func scroll_up(state *EditorState) {
	if state.upper_bound == 0 {
		return
	}

	state.upper_bound--
}

// Move the cursor 1 cell to the right
func move_right_cursor(state *EditorState) {
	if state.cursor.x >= len(state.text[state.cursor.y]) {
		old_pos := state.cursor

		var e error
		state.cursor, e = find_next_char(state, state.cursor)

		if e != nil {
			log.Print(e)
			state.cursor = old_pos
		}
		return
	} else {
		state.cursor.x++
	}

	screen_x, screen_y := state.screen.Size()

	if state.cursor.x > screen_x {
		state.cursor.x = 1
		state.cursor.y++
	}
	if state.cursor.y >= screen_y {
		scroll_down(state)
		state.cursor.y--
	}
}

// Move the cursor one cell to the back.
func move_back_cursor(state *EditorState) {
	state.cursor.x--

	if state.cursor.x < 0 && state.cursor.y == 0 {
		state.cursor.x = 0
	} else if state.cursor.x < 0 {
		state.cursor.y--
		state.cursor.x = find_last_char(state, state.cursor.y)
	}

	if state.cursor.y < state.upper_bound {
		scroll_up(state)
		state.cursor.y = 0
	}
}

// Move the cursor one cell up.
func move_up_cursor(state *EditorState) {
	if state.cursor.y <= 0 {
		scroll_up(state)
		return
	}
	if state.cursor.y < 0 {
		state.cursor.y = 0
		return
	}

	if state.cursor.x >= len(state.text[state.cursor.y-1]) ||
		state.text[state.cursor.y-1][state.cursor.x] == ' ' {
		state.cursor.x = find_last_char(state, state.cursor.y-1)
	}
	state.cursor.y--
}

// Move the cursor one cell down.
func move_down_cursor(state *EditorState) {
	// Restrict cursor from going beyond existing characters.
	if state.upper_bound+state.cursor.y >= len(state.text) {
		return
	}

	// Restrict cursor from exiting the screen.
	_, screen_y := state.screen.Size()
	if state.cursor.y >= screen_y {
		scroll_down(state)
		state.cursor.y--
		return
	}

	if state.cursor.x >= len(state.text[state.cursor.y+1]) ||
		state.text[state.cursor.y+1][state.cursor.x] == ' ' {
		state.cursor.x = find_first_char(state, state.cursor.y+1)
	}

	state.cursor.y++
}

// Updates the cursor's coordinates.
func update_cursor(state *EditorState) {
	state.screen.ShowCursor(state.cursor.x, state.cursor.y)
}
