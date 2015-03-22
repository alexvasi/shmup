package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
	size   mgl.Vec2
	screen mgl.Vec2
	ortho  mgl.Mat4

	polyShader PolyShader
	neonShader NeonShader
}

const (
	DefaultGroup PolyGroup = iota
	NeonGroup
)

func NewRenderer(width, height, screenWidth, screenHeight float32) *Renderer {
	r := &Renderer{
		size:   mgl.Vec2{width, height},
		screen: mgl.Vec2{screenWidth, screenHeight},
		ortho:  mgl.Ortho2D(0, width, 0, height),
	}

	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.BLEND)

	r.polyShader.Init(&r.ortho)
	r.neonShader.Init(r.screen)

	return r
}

func (r *Renderer) Clear() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	r.polyShader.Clear()
	r.neonShader.Clear()
}

func (r *Renderer) Render() {
	r.polyShader.Render()

	r.neonShader.BindFramebuffer()
	r.polyShader.Render(NeonGroup)
	r.neonShader.Render()
}

func (r *Renderer) Draw(points []mgl.Vec2, color mgl.Vec3) {
	r.polyShader.AddPoints(points, color, DefaultGroup)
}

func (r *Renderer) DrawNeon(points []mgl.Vec2, color mgl.Vec3) {
	r.polyShader.AddPoints(points, color, NeonGroup)
}

func (r *Renderer) DrawPoly(pos, size mgl.Vec2, sides int, color mgl.Vec3) {

	radius := size.Mul(0.5)
	points := mgl.Circle(radius.X(), radius.Y(), sides)
	for i := 0; i < len(points); i++ {
		points[i] = points[i].Add(pos)
	}

	r.polyShader.AddPoints(points, color, NeonGroup)
}
