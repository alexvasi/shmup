package main

import mgl "github.com/go-gl/mathgl/mgl32"

type Particle struct {
	dead      bool
	pos       mgl.Vec2
	velocity  mgl.Vec2
	startSize float32
	endSize   float32
	color     mgl.Vec3

	lifetime float32
	ttl      float32
}

func NewParticle(pos mgl.Vec2, startSize, endSize, lifetime float32,
	color mgl.Vec3, velocity mgl.Vec2) *Particle {

	return &Particle{
		pos:       pos,
		startSize: startSize,
		endSize:   endSize,
		color:     color,
		lifetime:  lifetime,
		ttl:       lifetime,
		velocity:  velocity,
	}
}

func (p *Particle) Update(dt float32) {
	p.pos = p.pos.Add(p.velocity.Mul(dt))
	if p.ttl < 0 {
		p.dead = true
	}
	p.ttl -= dt
}

func (p *Particle) Draw(renderer *Renderer) {
	size := (p.startSize-p.endSize)*(p.ttl/p.lifetime) + p.endSize
	renderer.DrawPoly(p.pos, mgl.Vec2{size, size}, 10, p.color)
}

func (p *Particle) IsDead() bool {
	return p.dead
}
