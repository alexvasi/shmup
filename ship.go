package main

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

const ShipMaxSpeed = 600
const ShipFireRate = 10

type Ship struct {
	pos  mgl.Vec2
	size mgl.Vec2
	dir  mgl.Vec2

	velocity mgl.Vec2
	thrust   mgl.Vec2

	fire     bool
	cooldown float32
}

func NewShip() *Ship {
	s := &Ship{
		size: mgl.Vec2{100, 100},
	}
	return s
}

func (s *Ship) Thrust(force mgl.Vec2, fire bool) {
	s.thrust = force
	s.fire = fire
}

func (s *Ship) Update(dt float32) []*Missile {
	oldV := s.velocity
	s.velocity = s.thrust.Mul(ShipMaxSpeed)

	speed := s.velocity.Len()
	if speed > ShipMaxSpeed {
		s.velocity = s.velocity.Mul(ShipMaxSpeed / speed)
	}

	s.pos = s.pos.Add(oldV.Add(s.velocity).Mul(0.5 * dt))

	s.cooldown -= dt
	if s.cooldown < 0 && s.fire {
		s.cooldown += 1. / ShipFireRate
		const S = 5
		return []*Missile{
			NewMissile(mgl.Vec2{s.pos.X(), s.pos.Y() + 4*50}, 1000, S),
			NewMissile(mgl.Vec2{s.pos.X(), s.pos.Y() + 3*50}, 800, S),
			NewMissile(mgl.Vec2{s.pos.X(), s.pos.Y() + 2*50}, 600, S),
			NewMissile(mgl.Vec2{s.pos.X(), s.pos.Y() + 1*50}, 400, S),
			NewMissile(mgl.Vec2{s.pos.X(), s.pos.Y()}, 200, S),
			NewMissile(mgl.Vec2{s.pos.X(), s.pos.Y() - 1*50}, 1800, S),
			NewMissile(mgl.Vec2{s.pos.X(), s.pos.Y() - 2*50}, 1600, S),
			NewMissile(mgl.Vec2{s.pos.X(), s.pos.Y() - 3*50}, 1200, S),
			NewMissile(mgl.Vec2{s.pos.X(), s.pos.Y() - 4*50}, 2000, S),
		}
	} else if s.cooldown < 0 {
		s.cooldown = 0
	}

	return nil
}

func (s *Ship) Draw(renderer *Renderer) {
	renderer.DrawPoly(s.pos, s.size, 30, mgl.Vec3{1, 1, 1})
}
