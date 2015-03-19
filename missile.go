package main

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Missile struct {
	pos      mgl.Vec2
	velocity mgl.Vec2
	size     float32
}

func NewMissile(pos mgl.Vec2, speed, size float32) *Missile {
	return &Missile{
		pos:      pos,
		velocity: mgl.Vec2{speed, 0},
		size:     size,
	}
}

func (m *Missile) Update(dt float32) {
	m.pos = m.pos.Add(m.velocity.Mul(dt))
}

func (m *Missile) Draw(renderer *Renderer) {
	renderer.DrawPoly(m.pos, mgl.Vec2{m.size * 2, m.size}, 10, mgl.Vec3{1, 1, 0})
}
