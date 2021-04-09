package varchive

import (
	"goncurses"
)

type Display interface {
	Init() 
	Clear() 
	Write(message string) 
	Close() 
	Flush() 
}


type DisplayImpl struct {
	window *goncurses.Window
}

func NewDisplay() DisplayImpl {
	return DisplayImpl{&goncurses.Window{}}
}

func (d DisplayImpl) Init() {
	window, err := goncurses.Init()

	if err != nil {
		fatal(err.Error())
	}
	d.window = window
}

func (d DisplayImpl) Clear() {
	d.window.Erase()
}

func (d DisplayImpl) Write(message string) {
	d.window.Println(message)
}

func (d DisplayImpl) Close() {
	goncurses.End()
}

func (d DisplayImpl) Flush() {
	d.window.Refresh()
}
