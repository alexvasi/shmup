package main

import (
	"fmt"

	"github.com/go-gl/glfw/v3.1/glfw"
)

type Timer struct {
	DT         float32
	TicksCount int

	updated    float64
	ticksTotal float32
}

func NewTimer() Timer {
	return Timer{
		updated: glfw.GetTime(),
	}
}

func (t *Timer) Tick() float32 {
	now := glfw.GetTime()
	t.DT = float32(now - t.updated)
	t.updated = now

	t.ticksTotal += t.DT
	t.TicksCount += 1

	return t.DT
}

func (t *Timer) Stat() string {
	avg := t.ticksTotal / float32(t.TicksCount)
	return fmt.Sprintf("%.2f %.f", avg*1000, 1/avg)
}

func (t *Timer) ResetCounter() {
	t.TicksCount = 0
	t.ticksTotal = 0
}
