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


// handy for 'dry run' mode, or when no console is available
type NoOpDisplay struct {}

func NewNoOpDisplay() *NoOpDisplay { return &NoOpDisplay{} }

func (*NoOpDisplay) Init() {}

func (*NoOpDisplay) Clear() {}

func (*NoOpDisplay) Write(string) {}

func (*NoOpDisplay) Close() {}

func (*NoOpDisplay) Flush() {}


type DisplayImpl struct {
	window *goncurses.Window
}

func NewDisplay() *DisplayImpl {
	return &DisplayImpl{}
}

func (d *DisplayImpl) Init() {
	window, err := goncurses.Init()

	if err != nil {
		fatal(err.Error())
	}
	d.window = window
}

func (d *DisplayImpl) Clear() {
	d.window.Erase()
}

func (d *DisplayImpl) Write(message string) {
	d.window.Println(message)
}

func (d *DisplayImpl) Close() {
	goncurses.End()
}

func (d *DisplayImpl) Flush() {
	d.window.Refresh()
}
