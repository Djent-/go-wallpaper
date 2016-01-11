package main

import (
	wdb "github.com/djent-/go-walldatabase"
	ui "github.com/gizak/termui"
)

struct Pane type {
	List ui.List // termui struct
	TotalItems []string // entire list of items
	CurrentIndex int // index of selected list item
}

struct Screen type {
	Title string
	LeftPane Pane
	RightPane Pane
	StatusBar ui.Par
}

func main() {
	var CurrentPane Pane
	var LeftPane Pane
	ui.Handle("sys/kbd/esc", func(ui.Event) {
		// press esc to quit
		ui.StopLoop()
	})
}