package main

import "github.com/go-gl/glfw/v3.0/glfw"

type Timer struct {
	DT      float32
	updated float64
}

func NewTimer() Timer {
	return Timer{DT: 0, updated: glfw.GetTime()}
}

func (t *Timer) Tick() float32 {
	now := glfw.GetTime()

	t.DT = float32(now - t.updated)
	t.updated = now

	return t.DT
}
