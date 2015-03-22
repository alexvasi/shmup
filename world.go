package main

import mgl "github.com/go-gl/mathgl/mgl32"

type World struct {
	size mgl.Vec2

	ships    []*Ship
	missiles []*Missile
	objects  []WorldObject
}

type Race int

const (
	Human Race = iota
	Others
)

type WorldObject interface {
	Update(dt float32)
	Draw(renderer *Renderer)
	IsDead() bool
}

func NewWorld(width, height float32) *World {
	w := World{
		size: mgl.Vec2{width, height},
	}
	return &w
}

func (w *World) Update(dt float32) {
	livingShips := w.ships[:0]
	for _, s := range w.ships {
		s.Update(dt, w)
		if !s.IsDead {
			livingShips = append(livingShips, s)
		}
	}
	w.ships = livingShips

	livingMissiles := w.missiles[:0]
	for _, m := range w.missiles {
		m.Update(dt, w, w.ships)
		if !m.IsDead {
			livingMissiles = append(livingMissiles, m)
		}
	}
	w.missiles = livingMissiles

	livingObjects := w.objects[:0]
	for _, o := range w.objects {
		o.Update(dt)
		if !o.IsDead() {
			livingObjects = append(livingObjects, o)
		}
	}
	w.objects = livingObjects
}

func (w *World) Draw(renderer *Renderer) {
	for _, o := range w.objects {
		o.Draw(renderer)
	}

	for _, s := range w.ships {
		s.Draw(renderer)
	}

	for _, m := range w.missiles {
		m.Draw(renderer)
	}
}

func (w *World) AddShips(ships ...*Ship) {
	w.ships = append(w.ships, ships...)
}

func (w *World) AddMissiles(missiles ...*Missile) {
	w.missiles = append(w.missiles, missiles...)
}

func (w *World) AddObjects(objects ...WorldObject) {
	w.objects = append(w.objects, objects...)
}
