package main

import (
	"math/rand"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type CargoStage struct {
	MinCount int
	MaxCount int
}

type FighterStage struct {
	ships []*Ship
}

type IntroStage struct {
	Stars Starfield
	Ship  *Ship

	time      float32
	totalTime float32
	sound     bool
}

type OutroStage struct {
	Stars Starfield
	Ship  *Ship

	time      float32
	totalTime float32
	started   bool
}

type FinalStage struct {
	papa      *Ship
	papaSound bool
	upPilot   *RoundShooter
	downPilot *RoundShooter
}

type DeathStage struct {
	Ship    *Ship
	time    float32
	started bool
}

func (s *CargoStage) Init(world *World) {
	const minSpeed = 0.2
	const maxSpeed = 0.5
	const posX = 1.05

	models := []*ShipModel{&CargoModel, &CargoModel2}
	count := s.MinCount
	if s.MaxCount > s.MinCount {
		rand.Intn(s.MaxCount - s.MinCount)
	}
	for i := 0; i < count; i++ {
		posY := float32(i+1) / (float32(count) + 1)
		speed := minSpeed + rand.Float32()*(maxSpeed-minSpeed)

		ship := NewShip(Others, models[rand.Intn(len(models))])
		ship.Pos = mgl.Vec2{world.Size.X() * posX, world.Size.Y() * posY}
		ship.Control(mgl.Vec2{-speed, 0}, false)
		world.AddShips(ship)
	}
}

func (s *CargoStage) Update(dt float32, world *World) bool {
	return world.ShipCount() > 1
}

func (s *FighterStage) Init(world *World) {
	const count = 8
	const minSpeed = 0.1
	const maxSpeed = 0.5
	const posX = 1.1

	s.ships = make([]*Ship, count)
	for i := 0; i < count; i++ {
		pos := float32(i+1) / (float32(count) + 1)
		speed := minSpeed + rand.Float32()*(maxSpeed-minSpeed)

		ship := NewShip(Others, &FighterModel)
		ship.Pos = mgl.Vec2{world.Size.X() * posX, world.Size.Y() * pos}
		ship.Control(mgl.Vec2{-speed, 0}, false)

		world.AddShips(ship)
		s.ships[i] = ship
	}
}

func (s *FighterStage) Update(dt float32, world *World) bool {
	const posX = 0.9

	if world.ShipCount() <= 1 {
		s.ships = []*Ship{}
		return false
	}

	for _, s := range s.ships {
		if s.Pos.X() <= world.Size.X()*posX {
			s.Control(mgl.Vec2{}, true)
		}
	}

	return true
}

func (s *FinalStage) Init(world *World) {
	const posX = 1.1
	const papaSpeed = 0.05

	s.papaSound = false
	s.papa = NewShip(Others, &PapaModel)
	s.papa.Pos = mgl.Vec2{world.Size.X() * posX, world.Size.Y() / 2}
	s.papa.Control(mgl.Vec2{-papaSpeed, 0}, false)
	world.AddShips(s.papa)

	s.upPilot = nil
	s.downPilot = nil
}
func (s *FinalStage) Update(dt float32, world *World) bool {
	const papaPosX = 0.9

	if !s.papaSound && s.papa.Pos.X() <= world.Size.X() {
		s.papaSound = true
		PlaySound("papa", 1, 0.8)
	}

	if s.papa.Pos.X() <= world.Size.X()*papaPosX {
		s.papa.Control(mgl.Vec2{}, true)
	}

	if s.upPilot == nil || s.upPilot.IsDead {
		s.upPilot = s.spawnShip(world, 0.9, 0, mgl.DegToRad(70))
	} else {
		s.upPilot.Update(dt, world)
	}

	if s.downPilot == nil || s.downPilot.IsDead {
		s.downPilot = s.spawnShip(world, 0.1, 0, mgl.DegToRad(-70))
	} else {
		s.downPilot.Update(dt, world)
	}

	if s.papa.IsDead {
		s.upPilot.Die()
		s.downPilot.Die()
		return false
	}

	return true
}

func (s *FinalStage) spawnShip(world *World, posY, angleMin,
	angleMax float32) *RoundShooter {

	const posX = 1.1
	const endPosX = 0.9

	ship := NewShip(Others, &ShooterModel)
	ship.Pos = mgl.Vec2{world.Size.X() * posX, world.Size.Y() * posY}
	pilot := NewRoundShooter(ship, endPosX, angleMin, angleMax)
	world.AddShips(ship)

	return pilot
}

func (s *IntroStage) Init(world *World) {
	const speed = 0.1
	const totalTime = 17

	s.Ship.Pos = mgl.Vec2{world.Size.X() * -0.1, world.Size.Y() / 2}
	s.Ship.Race = Autopilot
	s.Ship.Control(mgl.Vec2{speed, 0}, false)

	s.totalTime = totalTime
	s.time = 0
	s.sound = false
}

func (s *IntroStage) Update(dt float32, world *World) bool {
	const posX = 0.2
	const speed = 1000
	const starsTime = 1

	if s.Ship.Pos.X() > world.Size.X()*posX {
		s.Ship.Race = Human
	} else if s.Ship.Pos.X() > 0 && !s.sound {
		s.sound = true
		PlaySound("intro", 1, 1)
	}

	s.time += dt
	t := mgl.Clamp((starsTime-s.time)/s.totalTime, 0, 1)
	s.Stars.ChangeSpeed(1 + speed*t)

	if s.time > s.totalTime {
		s.Ship.Race = Human
		s.Stars.ChangeSpeed(1)
		return false
	}

	return true
}

func (s *OutroStage) Init(world *World) {
	const totalTime = 3

	s.Ship.Race = Autopilot
	s.Ship.Control(mgl.Vec2{}, false)

	s.totalTime = totalTime
	s.time = 0
	s.started = false

	world.ResetMissilesAndShips()
	world.AddShips(s.Ship)
}

func (s *OutroStage) Update(dt float32, world *World) bool {
	const starSpeed = 100
	const speed = 0.5

	pos := mgl.Vec2{world.Size.X() * 0.6, world.Size.Y() / 2}
	move := pos.Sub(s.Ship.Pos)

	if s.started {
		s.time += dt
		t := mgl.Clamp(s.time/s.totalTime, 0, 1)
		s.Stars.ChangeSpeed(1 + starSpeed*t)
		if s.time > s.totalTime {
			return false
		}
	} else {
		if move.Len() > world.Size.Y()*0.1 {
			s.Ship.Pos = s.Ship.Pos.Add(move.Mul(dt))
		} else {
			s.started = true
			s.Ship.Control(mgl.Vec2{speed, 0}, false)
			PlaySound("blip", 1, 1)
		}
	}

	return true
}

func (s *DeathStage) Init(world *World) {
	const totalTime = 5
	const timeSpeed = 0.2

	s.time = totalTime
	s.started = false

	s.Ship.Revive()
	world.TimeSpeed = timeSpeed
}

func (s *DeathStage) Update(dt float32, world *World) bool {
	const finalTime = 0.5
	const sizeBoost = 5

	s.time -= dt
	if !s.started && s.time < finalTime {
		s.started = true
		world.AddObjects(NewParticle(
			s.Ship.Pos,
			0,
			Max(world.Size.Elem())*sizeBoost,
			finalTime,
			BlackColor,
		))
		PlaySound("blip", 1, 1)
	}

	if s.time < 0 {
		world.ResetMissilesAndShips()
		world.TimeSpeed = 1

		s.Ship.Pos = mgl.Vec2{world.Size.X() * 0.1, world.Size.Y() / 2}
		world.AddShips(s.Ship)
		return false
	}

	return true
}
