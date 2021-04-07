package varchive

import (
	"goncurses"
)

type Display struct {
	window *goncurses.Window
}

func NewDisplay() *Display {
	return &Display{&goncurses.Window{}}
}

func (d *Display) Init() {
	window, err := goncurses.Init()

	if err != nil {
		fatal(err.Error())
	}
	d.window = window
}

func (d *Display) Clear() {
	d.window.Erase()
}

func (d *Display) Write(message string) {
	d.window.Println(message)
}

func (d *Display) Close() {
	goncurses.End()
}

func (d *Display) Flush() {
	d.window.Refresh()
}
