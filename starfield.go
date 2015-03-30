package main

import (
	"math"
	"math/rand"

	mgl "github.com/go-gl/mathgl/mgl32"
)

const StarMinSize = 1
const StarMaxSize = 4

type Star struct {
	Pos  mgl.Vec2
	Size float32
}

type StarStratum struct {
	size        mgl.Vec2
	stars       []Star
	stripeCount int
	speed       float32
	speedFactor float32
	color       mgl.Vec4
	minStarSize float32
	maxStarSize float32

	lastX      float32
	lastStripe int
}

type Starfield []*StarStratum

func NewStarfield(width, height float32) Starfield {
	const count = 20
	const stars = 50
	const minSpeed = 5
	const maxSpeed = 100
	minColor := HexColor("#001f5d", 1)
	maxColor := HexColor("#b0caff", 1)

	size := mgl.Vec2{width, height}

	sf := make(Starfield, count)
	for i := range sf {
		speed := minSpeed + (maxSpeed-minSpeed)*float32(i)/count
		color := BlendColors(minColor, maxColor, float32(i)/count)
		sf[i] = NewStarStratum(size, stars, speed, color)
	}
	return sf
}

func (sf Starfield) Update(dt float32) {
	for i := range sf {
		sf[i].Update(dt)
	}
}

func (sf Starfield) Draw(renderer *Renderer) {
	for i := range sf {
		sf[i].Draw(renderer)
	}
}

func (Starfield) IsDead() bool {
	return false
}

func (sf Starfield) ChangeSpeed(factor float32) {
	for i := range sf {
		sf[i].speedFactor = factor
	}
}

func NewStarStratum(size mgl.Vec2, count int, speed float32,
	color mgl.Vec4) *StarStratum {

	const extraStripes = 2

	stripeCount := int(math.Floor(math.Sqrt(float64(count)) + 0.5))
	ss := StarStratum{
		size:        size,
		stars:       make([]Star, stripeCount*(stripeCount+extraStripes)),
		stripeCount: stripeCount,
		speed:       speed,
		speedFactor: 1,
		color:       color,
	}

	for i := 0; i < stripeCount+extraStripes; i++ {
		ss.generateNextStripe()
	}

	return &ss
}

func (ss *StarStratum) Update(dt float32) {
	shift := dt * ss.speed * ss.speedFactor

	for i := range ss.stars {
		ss.stars[i].Pos[0] -= shift
	}
	ss.lastX -= shift

	stripeWidth := ss.size.X() / float32(ss.stripeCount)
	if ss.lastX < ss.size.X()+stripeWidth/2 {
		ss.generateNextStripe()
	}
}

func (ss *StarStratum) Draw(renderer *Renderer) {
	model := make([]mgl.Vec2, len(StarModel))

	for _, star := range ss.stars {
		for i, point := range StarModel {
			model[i] = point.Mul(star.Size / 2).Add(star.Pos)
		}
		renderer.Draw(model, ss.color, StarGroup)
	}
}

func (ss *StarStratum) IsDead() bool {
	return false
}

func (ss *StarStratum) starStripe(stripeNum int) []Star {
	from := stripeNum * ss.stripeCount
	to := from + ss.stripeCount
	return ss.stars[from:to]
}

func (ss *StarStratum) generateNextStripe() {
	width := ss.size.X() / float32(ss.stripeCount)
	height := ss.size.Y() / float32(ss.stripeCount)

	starStripe := ss.starStripe(ss.lastStripe)
	for i := 0; i < len(starStripe); i++ {
		starStripe[i].Pos[0] = ss.lastX + width*rand.Float32()
		starStripe[i].Pos[1] = height*float32(i) + height*rand.Float32()

		size := StarMinSize + (StarMaxSize-StarMinSize)*rand.Float32()
		starStripe[i].Size = size
	}

	const extraStripes = 2
	ss.lastX += width
	ss.lastStripe = (ss.lastStripe + 1) % (ss.stripeCount + extraStripes)
}
