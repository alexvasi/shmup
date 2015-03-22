package main

import mgl "github.com/go-gl/mathgl/mgl32"

type Missile struct {
	IsDead   bool
	race     Race
	pos      mgl.Vec2
	velocity mgl.Vec2
	size     mgl.Vec2
}

func NewMissile(race Race, pos mgl.Vec2, speed, size float32) *Missile {
	return &Missile{
		race:     race,
		pos:      pos,
		velocity: mgl.Vec2{speed, 0},
		size:     mgl.Vec2{1.5 * size, size},
	}
}

func (m *Missile) Update(dt float32, world *World, ships []*Ship) {
	newPos := m.pos.Add(m.velocity.Mul(dt))

	halfSize := Max(m.size.X(), m.size.Y()) / 2
	aabbMin := MinVec2(newPos, m.pos)
	aabbMax := MaxVec2(newPos, m.pos)
	aabb := mgl.Vec4{aabbMin.X() - halfSize, aabbMin.Y() - halfSize,
		aabbMax.X() + halfSize, aabbMax.Y() + halfSize}

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
			ok, point := SegmentIntersection(newPos, m.pos,
				side[0], side[1])
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
		hitShip.Hit()
		m.IsDead = true
		world.AddObjects(NewParticle(hitPos, 15, 15, 0.05, mgl.Vec3{1, 1, 0}, mgl.Vec2{}))
	}
}

func (m *Missile) Draw(renderer *Renderer) {
	renderer.DrawPoly(m.pos, m.size, 10, mgl.Vec3{1, 1, 0})
}
