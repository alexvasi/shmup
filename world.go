package main

import mgl "github.com/go-gl/mathgl/mgl32"

type World struct {
	Size      mgl.Vec2
	TimeSpeed float32

	ships    []*Ship
	missiles []*Missile
	objects  []WorldObject
}

type Race int

const (
	Human Race = iota
	Others
	Autopilot
)

type WorldObject interface {
	Update(dt float32)
	Draw(renderer *Renderer)
	IsDead() bool
}

func NewWorld(width, height float32) *World {
	w := World{
		Size:      mgl.Vec2{width, height},
		TimeSpeed: 1,
	}
	return &w
}

func (w *World) Update(dt float32) {
	livingShips := w.ships[:0]
	for _, s := range w.ships {
		s.Update(dt*w.TimeSpeed, w, w.ships)
		if !s.IsDead {
			livingShips = append(livingShips, s)
		}
	}
	w.ships = livingShips

	livingMissiles := w.missiles[:0]
	for _, m := range w.missiles {
		m.Update(dt*w.TimeSpeed, w, w.ships)
		w.killStrayedMissile(m)
		if !m.IsDead {
			livingMissiles = append(livingMissiles, m)
		}
	}
	w.missiles = livingMissiles

	livingObjects := w.objects[:0]
	for _, o := range w.objects {
		o.Update(dt * w.TimeSpeed)
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

func (w *World) ShipCount() int {
	return len(w.ships)
}

func (w *World) ResetMissilesAndShips() {
	w.ships = nil
	w.missiles = nil
}

func (w *World) killStrayedMissile(missile *Missile) {
	aabb := missile.AABB(mgl.Vec2{})
	if !CheckAABB(aabb, mgl.Vec4{0, 0, w.Size.X(), w.Size.Y()}) {
		missile.IsDead = true
	}
}
