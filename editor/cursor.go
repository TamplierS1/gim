package main

import "github.com/gdamore/tcell/v2"

// Move the cursor 1 cell to the right
func move_right_cursor(screen tcell.Screen) {
	g_cursor_pos.x++
	screen_x, screen_y := screen.Size()

	if g_cursor_pos.x > screen_x {
		g_cursor_pos.x = 1
		g_cursor_pos.y++
	}
	if g_cursor_pos.y > screen_y {
		// TODO: the screen should scroll down here
	}
}

// Move the cursor one cell to the back.
func move_back_cursor(screen tcell.Screen) {
	g_cursor_pos.x--

	if g_cursor_pos.x < 0 && g_cursor_pos.y == 0 {
		g_cursor_pos.x = 0
	} else if g_cursor_pos.x < 0 {
		g_cursor_pos.y--
		g_cursor_pos.x = find_last_char(screen, g_cursor_pos.y)
	}
}

// Move the cursor one cell up.
func move_up_cursor(screen tcell.Screen) {
	g_cursor_pos.y--

	if g_cursor_pos.y < 0 {
		g_cursor_pos.y = 0
	}
}

// Move the cursor one cell up.
func move_down_cursor(screen tcell.Screen) {
	g_cursor_pos.y++

	screen_y, _ := screen.Size()
	if g_cursor_pos.y >= screen_y {
		g_cursor_pos.y--
	}
}

// Updates the cursor's coordinates.
func update_cursor(screen tcell.Screen) {
	screen.ShowCursor(g_cursor_pos.x, g_cursor_pos.y)
}

// Set a cell at the cursor's position to 'rune'
func putchar_at_cursor(screen tcell.Screen, char rune) {
	PutChar(screen, g_cursor_pos.x, g_cursor_pos.y, g_default_style, char)
}
