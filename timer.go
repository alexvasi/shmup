package main

import (
	"fmt"

	"github.com/go-gl/glfw/v3.1/glfw"
)

type Timer struct {
	DT      float32
	updated float64
	window  *glfw.Window

	show         bool
	ticksAvg     float32
	ticksCounter int
	ticksTotal   int
}

func NewTimer(window *glfw.Window) Timer {
	return Timer{
		DT:      0,
		updated: glfw.GetTime(),
		window:  window,
	}
}

func (t *Timer) Tick() float32 {
	now := glfw.GetTime()

	t.DT = float32(now - t.updated)
	t.updated = now

	if t.show {
		t.ticksAvg += t.DT / float32(t.ticksTotal)
		t.ticksCounter += 1
		if t.ticksCounter == t.ticksTotal {
			text := fmt.Sprintf("%.2f %.f", t.ticksAvg*1000, 1/t.ticksAvg)
			fmt.Println(text)
			t.window.SetTitle(text)
			t.ticksAvg = 0
			t.ticksCounter = 0
		}
	}

	return t.DT
}

func (t *Timer) ShowTimings(show bool, every int) {
	if t.show == show {
		return
	}

	t.show = show
	if show {
		glfw.SwapInterval(0)
		t.ticksAvg = 0
		t.ticksCounter = 0
		t.ticksTotal = every
	} else {
		glfw.SwapInterval(1)
	}
}
