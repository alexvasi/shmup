package main

import mgl "github.com/go-gl/mathgl/mgl32"

type Missile struct {
	IsDead   bool
	race     Race
	pos      mgl.Vec2
	velocity mgl.Vec2
	size     mgl.Vec2
	color    mgl.Vec4
}

func NewMissile(race Race, pos, velocity, size mgl.Vec2, color mgl.Vec4) *Missile {
	return &Missile{
		race:     race,
		pos:      pos,
		velocity: velocity,
		size:     size,
		color:    color,
	}
}

func (m *Missile) Update(dt float32, world *World, ships []*Ship) {
	newPos := m.pos.Add(m.velocity.Mul(dt))
	aabb := m.AABB(m.velocity.Mul(dt))

	var hitShip *Ship
	var hitPos mgl.Vec2
	var hitDistance = mgl.MaxValue

	for _, s := range ships {
		if m.race == s.Race {
			continue
		}
		if !CheckAABB(aabb, s.AABB()) {
			continue
		}
		for _, side := range s.Sides() {
			ok, point := m.intersection(newPos, side[0], side[1])
			if ok {
				hitShip = s
				distance := m.pos.Sub(point).Len()
				if hitDistance > distance {
					hitPos = point
					hitDistance = distance
				}
			}
		}
		if hitShip != nil {
			break
		}
	}

	if hitShip == nil {
		m.pos = newPos
	} else {
		const ttl = 0.05
		size := m.size.Y() * 3

		hitShip.Hit()
		m.IsDead = true
		explosion := NewParticle(hitPos, size, size, ttl, m.color)
		world.AddObjects(explosion)
	}
}

func (m *Missile) Draw(renderer *Renderer) {
	const huge = 25

	sides := 10
	if Max(m.size.Elem()) > huge {
		sides *= 2
	}
	renderer.DrawPoly(m.pos, m.size, sides, m.color, NeonGroup)
}

func (m *Missile) AABB(movement mgl.Vec2) mgl.Vec4 {
	newPos := m.pos.Add(movement)
	halfSize := Max(m.size.X(), m.size.Y()) / 2
	aabbMin := MinVec2(newPos, m.pos)
	aabbMax := MaxVec2(newPos, m.pos)
	aabb := mgl.Vec4{aabbMin.X() - halfSize, aabbMin.Y() - halfSize,
		aabbMax.X() + halfSize, aabbMax.Y() + halfSize}
	return aabb
}

func (m *Missile) intersection(newPos, p1, p2 mgl.Vec2) (bool, mgl.Vec2) {
	ok, point := SegmentIntersection(newPos, m.pos, p1, p2)
	if ok {
		return ok, point
	}

	for _, x := range []float32{-0.5, 0.5} {
		for _, y := range []float32{-0.5, 0.5} {
			a1 := m.pos.Add(mgl.Vec2{m.size.X() * x, m.size.Y() * y})
			a2 := newPos.Add(mgl.Vec2{m.size.X() * x, m.size.Y() * y})
			ok, point := SegmentIntersection(a1, a2, p1, p2)
			if ok {
				return ok, point
			}
		}
	}

	return false, mgl.Vec2{}
}
