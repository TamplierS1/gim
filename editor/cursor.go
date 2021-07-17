package main

// Move the cursor 1 cell to the right
func move_right_cursor(state *EditorState) {
	state.cursor.x++
	screen_x, screen_y := state.screen.Size()

	if state.cursor.x > screen_x {
		state.cursor.x = 1
		state.cursor.y++
	}
	if state.cursor.y > screen_y {
		// TODO: the screen should scroll down here
	}
}

// Move the cursor one cell to the back.
func move_back_cursor(state *EditorState) {
	state.cursor.x--

	if state.cursor.x < 0 && state.cursor.y == 0 {
		state.cursor.x = 0
	} else if state.cursor.x < 0 {
		state.cursor.y--
		state.cursor.x = find_last_char(state.screen, state.cursor.y)
	}
}

// Move the cursor one cell up.
func move_up_cursor(cursor *Vec2) {
	cursor.y--

	if cursor.y < 0 {
		cursor.y = 0
	}
}

// Move the cursor one cell up.
func move_down_cursor(state *EditorState) {
	state.cursor.y++

	screen_y, _ := state.screen.Size()
	if state.cursor.y >= screen_y {
		state.cursor.y--
	}
}

// Updates the cursor's coordinates.
func update_cursor(state *EditorState) {
	state.screen.ShowCursor(state.cursor.x, state.cursor.y)
}
