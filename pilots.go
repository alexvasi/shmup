package main

import mgl "github.com/go-gl/mathgl/mgl32"

type RoundShooter struct {
	IsDead bool

	ship     *Ship
	endPosX  float32
	angleMin float32
	angleMax float32
	dir      mgl.Vec2
	t        float32
}

func NewRoundShooter(ship *Ship, endPosX, angleMin, angleMax float32) *RoundShooter {
	rs := &RoundShooter{
		ship:     ship,
		endPosX:  endPosX,
		angleMin: angleMin,
		angleMax: angleMax,
		dir:      ship.Dir,
	}
	return rs
}

func (rs *RoundShooter) Update(dt float32, world *World) {
	const speed = 0.2
	const rSpeed = 0.5

	if rs.ship.IsDead {
		rs.IsDead = true
		return
	}

	if rs.ship.Pos.X() > rs.endPosX*world.Size.X() {
		rs.ship.Control(mgl.Vec2{-speed, 0}, false)
	} else {
		rs.ship.Control(mgl.Vec2{0, 0}, true)

		rs.t += rSpeed * dt
		for rs.t > 1 {
			rs.t -= 2
		}
		angle := rs.angleMin + mgl.Abs(rs.t)*(rs.angleMax-rs.angleMin)
		rs.ship.Dir = mgl.Rotate2D(angle).Mul2x1(rs.dir).Normalize()
	}
}

func (rs *RoundShooter) Die() {
	for rs.ship.Health() > 0 {
		rs.ship.Hit()
	}
}
