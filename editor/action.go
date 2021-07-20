package main

import (
	"bufio"
	"log"
	"os"
)

type Action interface {
	execute(state *EditorState)
}

type CloseEditorAction struct{}

func (action CloseEditorAction) execute(state *EditorState) {
	exit(0, state.screen)
}

type EraseCharAction struct{}

func (action EraseCharAction) execute(state *EditorState) {
	move_back_cursor(state)
	// putchar_at_cursor(state, ' ')
	remove_symbol(state, state.cursor)
}

type EraseWordAction struct{}

func (action EraseWordAction) execute(state *EditorState) {
	erase_prev_word(state)
}

type NewLineAction struct{}

func (action NewLineAction) execute(state *EditorState) {
	record_rune(state, state.cursor, '\n')

	move_down_cursor(state)
	state.cursor.x = 0
}

type Direction int

const (
	MoveRight Direction = iota
	MoveLeft            = iota
	MoveDown            = iota
	MoveUp              = iota
)

type MoveCursorAction struct {
	move_where Direction
}

func (action MoveCursorAction) execute(state *EditorState) {
	switch action.move_where {
	case MoveRight:
		move_right_cursor(state)
	case MoveLeft:
		move_back_cursor(state)
	case MoveDown:
		move_down_cursor(state)
	case MoveUp:
		move_up_cursor(state)
	}
}

type EnterCharAction struct {
	r rune
}

func (action EnterCharAction) execute(state *EditorState) {
	record_rune(state, state.cursor, action.r)

	// putchar_at_cursor(state, action.r)
	move_right_cursor(state)
}

type SaveAndCloseEditorAction struct{}

func (action SaveAndCloseEditorAction) execute(state *EditorState) {
	file, err := os.Create(state.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		file.Close()
		exit(0, state.screen)
	}()

	writer := bufio.NewWriter(file)

	for i := range state.text {
		writer.WriteString(state.text[i])
	}

	writer.Flush()
}
