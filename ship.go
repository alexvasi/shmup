package main

import (
	"math"
	"math/rand"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Ship struct {
	Race   Race
	IsDead bool
	Pos    mgl.Vec2
	Dir    mgl.Vec2

	velocity mgl.Vec2
	model    *ShipModel
	hp       int
	damaged  bool
	trs      mgl.Mat3

	fire     bool
	cooldown []float32
	engineCD []float32
}

func NewShip(race Race, model *ShipModel) *Ship {
	s := &Ship{
		Race:     race,
		model:    model,
		Dir:      mgl.Vec2{1, 0},
		hp:       model.Hp,
		cooldown: make([]float32, len(model.Guns)),
		engineCD: make([]float32, len(model.Engines)),
	}

	if race != Human {
		s.Dir[0] *= -1
	}

	return s
}

func (s *Ship) Control(thrust mgl.Vec2, fire bool) {
	s.fire = fire

	s.velocity = thrust.Mul(s.model.Speed)
	speed := s.velocity.Len()
	if speed > s.model.Speed {
		s.velocity = s.velocity.Mul(s.model.Speed / speed)
	}
}

func (s *Ship) Update(dt float32, world *World, ships []*Ship) {
	for _, other := range ships {
		if s.Race != other.Race && s.collides(other) {
			other.Hit()
			s.Hit()
		}
	}

	if s.hp <= 0 {
		s.IsDead = true
		s.makeExplosion(world)
		return
	}

	s.Pos = s.Pos.Add(s.velocity.Mul(dt))
	s.StayInWorld(world)
	s.trs = s.calcTRS(1)
	s.updateGuns(dt, world)
	s.updateEngines(dt, world)
}

func (s *Ship) updateGuns(dt float32, world *World) {
	for i, gun := range s.model.Guns {
		s.cooldown[i] -= dt
		for ; s.fire && s.cooldown[i] < 0; s.cooldown[i] += 1 / gun.Rate {
			pos := s.transformPoint(gun.Pos)
			v := s.Dir.Mul(gun.Speed)
			m := NewMissile(s.Race, pos, v, gun.Size, gun.Color)
			world.AddMissiles(m)
			if len(gun.Sound) > 0 {
				PlaySound(gun.Sound, gun.SoundGain, gun.SoundPitch)
			}
		}
		if !s.fire && s.cooldown[i] < 0 {
			s.cooldown[i] = 0
		}
	}
}

func (s *Ship) updateEngines(dt float32, world *World) {
	const speed = 100
	const MaxSideVelocity = 90

	for i, engine := range s.model.Engines {
		s.engineCD[i] -= dt

		proj := s.velocity.Dot(s.Dir)
		rej := s.velocity.Sub(s.Dir.Mul(proj))
		active := engine.MinVelocity <= proj && MaxSideVelocity > rej.Len()

		for ; active && s.engineCD[i] < 0; s.engineCD[i] += 1 / engine.Rate {
			shift := rand.Float32() - 0.5
			pos := s.transformPoint(engine.Pos)
			posDir := mgl.Vec2{-s.Dir.Y(), s.Dir.X()}
			pos = pos.Add(posDir.Mul(shift * engine.Size / 2))

			particle := NewParticle(
				pos,
				engine.ParticleSize[1],
				engine.ParticleSize[0],
				float32(engine.TTL)+engine.TTL*rand.Float32(),
				engine.Color,
			)
			particle.Velocity = s.Dir.Mul(-speed)
			particle.RenderGroup = EngineGroup
			world.AddObjects(particle)
		}
		if !active && s.engineCD[i] < 0 {
			s.engineCD[i] = 0
		}
	}
}

func (s *Ship) Draw(renderer *Renderer) {
	if s.damaged {
		s.damaged = false
		renderer.Draw(s.transform(s.model.Hull), WhiteColor, PlainGroup)
	} else {
		renderer.Draw(s.transform(s.model.Hull), s.model.Color1, PlainGroup)
		renderer.Draw(s.transformScale(s.model.Hull, 0.5), s.model.Color2, PlainGroup)
	}
}

func (s *Ship) Hit() {
	if s.Race != Autopilot {
		s.hp -= 1
		s.damaged = true
	}
}

func (s *Ship) AABB() mgl.Vec4 {
	halfSize := Max(s.model.Size.X(), s.model.Size.Y()) / 2
	aabb := mgl.Vec4{
		s.Pos.X() - halfSize,
		s.Pos.Y() - halfSize,
		s.Pos.X() + halfSize,
		s.Pos.Y() + halfSize,
	}
	return aabb
}

func (s *Ship) Sides() [][2]mgl.Vec2 {
	points := s.transform(s.model.Hull)

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

func (s *Ship) StayInWorld(world *World) {
	width, height := world.Size.Elem()
	minX, minY, maxX, maxY := s.AABB().Elem()

	if s.Race == Human {
		s.Pos[0] -= Min(minX, 0)
		s.Pos[1] -= Min(minY, 0)
		s.Pos[0] += Min(width-maxX, 0)
		s.Pos[1] += Min(height-maxY, 0)
	} else if s.Race == Others && maxX < 0 {
		s.IsDead = true // silently go away to another world
	}
}

func (s *Ship) Health() float32 {
	return float32(s.hp) / float32(s.model.Hp)
}

func (s *Ship) Revive() {
	s.hp = s.model.Hp
	s.IsDead = false
	s.damaged = false
}

func (s *Ship) calcTRS(scale float32) mgl.Mat3 {
	dirAngle := -float32(math.Atan2(float64(s.Dir.X()), float64(s.Dir.Y())))

	mat := mgl.Scale2D(scale*s.model.Size.X()/2, scale*s.model.Size.Y()/2)
	mat = mgl.HomogRotate2D(dirAngle).Mul3(mat)
	mat = mgl.Translate2D(s.Pos.X(), s.Pos.Y()).Mul3(mat)

	return mat
}

func (s *Ship) transform(points []mgl.Vec2) []mgl.Vec2 {
	result := make([]mgl.Vec2, len(points))
	for i, p := range points {
		result[i] = s.trs.Mul3x1(mgl.Vec3{p.X(), p.Y(), 1}).Vec2()
	}
	return result
}

func (s *Ship) transformPoint(p mgl.Vec2) mgl.Vec2 {
	return s.trs.Mul3x1(mgl.Vec3{p.X(), p.Y(), 1}).Vec2()
}

func (s *Ship) transformScale(points []mgl.Vec2, scale float32) []mgl.Vec2 {
	mat := s.calcTRS(scale)

	result := make([]mgl.Vec2, len(points))
	for i, p := range points {
		result[i] = mat.Mul3x1(mgl.Vec3{p.X(), p.Y(), 1}).Vec2()
	}
	return result
}

func (s *Ship) makeExplosion(world *World) {
	boomTTL := 0.5 * s.model.BlowupFactor
	size := 15 * s.model.BlowupFactor
	ttl := 2 * s.model.BlowupFactor
	count := 35 * s.model.BlowupFactor
	velocityMin := 50 * s.model.BlowupFactor
	velocityMax := 200 * s.model.BlowupFactor

	bigBoom := NewParticle(
		s.Pos,
		Max(s.model.Size.Elem())*s.model.BlowupFactor,
		0,
		boomTTL,
		WhiteColor,
	)
	world.AddObjects(bigBoom)

	colors := []mgl.Vec4{s.model.Color1, s.model.Color2}
	for i := 0; i < int(count); i++ {
		angleMin := float64(i) * 2 * math.Pi / float64(count)
		angleMax := float64(i+1) * 2 * math.Pi / float64(count)
		angle := angleMin + (angleMax-angleMin)*rand.Float64()
		sin, cos := math.Sincos(angle)

		v := velocityMin + (velocityMax-velocityMin)*rand.Float32()
		p := NewParticle(
			s.Pos,
			(size/2)*rand.Float32()+size/2,
			0,
			(ttl/2)*rand.Float32()+ttl/2,
			colors[rand.Intn(2)],
		)
		p.Velocity = mgl.Vec2{float32(sin) * v, float32(cos) * v}
		world.AddObjects(p)
	}

	if s.Race == Human {
		PlaySound("boom", 1, 0.3)
	} else {
		PlaySound("boom", 1, 1/s.model.BlowupFactor)
	}
}

func (s *Ship) collides(other *Ship) bool {
	if CheckAABB(s.AABB(), other.AABB()) {
		for _, side := range s.Sides() {
			for _, oSide := range other.Sides() {
				ok, _ := SegmentIntersection(side[0], side[1],
					oSide[0], oSide[1])
				if ok {
					return true
				}
			}
		}
	}
	return false
}
