package main

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"
)

const ShipMaxSpeed = 600
const ShipFireRate = 10

type Ship struct {
	Race     Race
	IsDead   bool
	pos      mgl.Vec2
	dir      mgl.Vec2
	velocity mgl.Vec2
	model    ShipModel
	hp       int
	damaged  bool

	fire     bool
	cooldown float32
}

func NewShip(race Race, posX, posY float32) *Ship {
	s := &Ship{
		Race:  race,
		pos:   mgl.Vec2{posX, posY},
		dir:   mgl.Vec2{1, 0},
		model: PlayerModel,
		hp:    10,
	}
	return s
}

func (s *Ship) Control(thrust mgl.Vec2, fire bool) {
	s.fire = fire

	s.velocity = thrust.Mul(ShipMaxSpeed)
	speed := s.velocity.Len()
	if speed > ShipMaxSpeed {
		s.velocity = s.velocity.Mul(ShipMaxSpeed / speed)
	}
}

func (s *Ship) Update(dt float32, world *World) {
	if s.hp <= 0 {
		s.IsDead = true
		s.makeExplosion(world)
		return
	}

	s.pos = s.pos.Add(s.velocity.Mul(dt))

	s.cooldown -= dt
	if s.cooldown < 0 && s.fire {
		s.cooldown += 1. / ShipFireRate
		missile := NewMissile(s.Race, s.transform(s.model.gun)[0], 1000, 5)
		world.AddMissiles(missile)
	} else if s.cooldown < 0 {
		s.cooldown = 0
	}
}

func (s *Ship) Draw(renderer *Renderer) {
	color := s.model.color
	if s.damaged {
		s.damaged = false
		color = mgl.Vec3{1, 1, 1}
	}

	points := s.transform(s.model.hull...)
	renderer.Draw(points, color)
	x1, y1, x2, y2 := s.AABB().Elem()

	renderer.DrawPoly(mgl.Vec2{x1, y1}, mgl.Vec2{3, 3}, 3, mgl.Vec3{1})
	renderer.DrawPoly(mgl.Vec2{x2, y2}, mgl.Vec2{3, 3}, 3, mgl.Vec3{1})
	renderer.DrawPoly(mgl.Vec2{x1, y2}, mgl.Vec2{3, 3}, 3, mgl.Vec3{1})
	renderer.DrawPoly(mgl.Vec2{x2, y1}, mgl.Vec2{3, 3}, 3, mgl.Vec3{1})
}

func (s *Ship) Hit() {
	s.hp -= 1
	s.damaged = true
}

func (s *Ship) AABB() mgl.Vec4 {
	halfSize := Max(s.model.size.X(), s.model.size.Y()) / 2
	aabb := mgl.Vec4{
		s.pos.X() - halfSize,
		s.pos.Y() - halfSize,
		s.pos.X() + halfSize,
		s.pos.Y() + halfSize,
	}
	return aabb
}

func (s *Ship) Sides() [][2]mgl.Vec2 {
	points := s.transform(s.model.hull...)

	sides := make([][2]mgl.Vec2, 0, len(points))
	for i := 0; i < len(points)/3; i++ {
		p1 := points[i*3+0]
		p2 := points[i*3+1]
		p3 := points[i*3+2]
		sides = append(sides,
			[2]mgl.Vec2{p1, p2},
			[2]mgl.Vec2{p2, p3},
			[2]mgl.Vec2{p3, p1},
		)
	}

	return sides
}

func (s *Ship) transform(points ...mgl.Vec2) []mgl.Vec2 {
	dirAngle := -float32(math.Atan2(float64(s.dir.X()), float64(s.dir.Y())))

	mat := mgl.Scale2D(s.model.size.X()/2, s.model.size.Y()/2)
	mat = mgl.HomogRotate2D(dirAngle).Mul3(mat)
	mat = mgl.Translate2D(s.pos.X(), s.pos.Y()).Mul3(mat)

	result := make([]mgl.Vec2, len(points))
	for i, p := range points {
		result[i] = mat.Mul3x1(mgl.Vec3{p.X(), p.Y(), 1}).Vec2()
	}
	return result
}

func (s *Ship) makeExplosion(world *World) {
	bigBoom := NewParticle(
		s.pos,
		Max(s.model.size.Elem()),
		0,
		0.5,
		mgl.Vec3{1, 1, 1},
		mgl.Vec2{},
	)
	world.AddObjects(bigBoom)

	const count = 10
	const velocity = 100

	for i := 0; i < count; i++ {
		sin, cos := math.Sincos(float64(i) * 2 * math.Pi / count)
		v := mgl.Vec2{float32(sin) * velocity, float32(cos) * velocity}
		p := NewParticle(
			s.pos,
			40,
			0,
			2,
			s.model.color,
			v,
		)
		world.AddObjects(p)
	}
}
