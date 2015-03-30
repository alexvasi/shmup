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
	starShader NeonShader
}

const (
	PlainGroup PolyGroup = iota
	NeonGroup
	StarGroup
	EngineGroup
)

var WhiteColor = mgl.Vec4{1, 1, 1, 1}
var BlackColor = mgl.Vec4{0, 0, 0, 1}

func NewRenderer(width, height float32, screenSize mgl.Vec2) *Renderer {
	r := &Renderer{
		size:   mgl.Vec2{width, height},
		screen: screenSize,
		ortho:  mgl.Ortho2D(0, width, 0, height),
	}

	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.BLEND)

	r.polyShader.Init(&r.ortho)
	r.neonShader.Init(r.screen, 0.5, true)
	r.starShader.Init(r.screen, 1, false)

	return r
}

func (r *Renderer) Clear() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	r.polyShader.Clear()
	r.neonShader.Clear()
	r.starShader.Clear()
}

func (r *Renderer) Render() {
	r.starShader.BindFramebuffer()
	r.polyShader.Render(StarGroup, EngineGroup)
	r.starShader.Render()
	r.polyShader.Render(StarGroup, EngineGroup)

	r.polyShader.Render(PlainGroup, NeonGroup)

	r.neonShader.BindFramebuffer()
	r.polyShader.Render(NeonGroup)
	r.neonShader.Render()
}

func (r *Renderer) Draw(points []mgl.Vec2, color mgl.Vec4, group PolyGroup) {
	r.polyShader.AddPoints(points, color, group)
}

func (r *Renderer) DrawPoly(pos, size mgl.Vec2, sides int, color mgl.Vec4,
	group PolyGroup) {

	radius := size.Mul(0.5)
	points := mgl.Circle(radius.X(), radius.Y(), sides)
	for i := 0; i < len(points); i++ {
		points[i] = points[i].Add(pos)
	}

	r.polyShader.AddPoints(points, color, group)
}
