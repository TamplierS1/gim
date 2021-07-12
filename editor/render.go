// Functions that render characters to the screen
package main

import "github.com/gdamore/tcell/v2"

// Draw 'text' in the box defined by 'x1, y1' and 'x2, y2'
// on the 'screen' with 'style'
func DrawText(screen tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	col, row := x1, y1
	for _, r := range text {
		if col > x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}

		screen.SetContent(col, row, r, nil, style)
		col++
	}
}

func PutChar(screen tcell.Screen, x, y int, style tcell.Style, character rune) {
	screen.SetContent(x, y, character, nil, style)
}
