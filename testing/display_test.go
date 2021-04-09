package main

type NoOpDisplay struct {}

func NewNoOpDisplay() *NoOpDisplay { return &NoOpDisplay{} }

func (*NoOpDisplay) Init() {}

func (*NoOpDisplay) Clear() {}

func (*NoOpDisplay) Write(string) {}

func (*NoOpDisplay) Close() {}

func (*NoOpDisplay) Flush() {}
