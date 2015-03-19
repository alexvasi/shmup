package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type PolyShader struct {
	program uint32
	vao     uint32
	vbo     uint32

	uniProj int32

	points []float32
	vbSize int
}

const polyVertexSrc string = `
#version 150 core

in vec2 pos;
in vec3 color;

out vec3 vColor;

uniform mat4 proj;

void main()
{
    vColor = color;
    gl_Position = proj * vec4(pos, 0, 1);
}
` + "\x00"

const polyFragSrc string = `
#version 150 core

in vec3 vColor;

out vec4 outColor;

void main() {
    outColor = vec4(vColor, 1);
}
` + "\x00"

func (s *PolyShader) Init(proj *mgl.Mat4) {
	s.program = CreateShaderProgram(polyVertexSrc, polyFragSrc)

	gl.UseProgram(s.program)
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)

	s.uniProj = gl.GetUniformLocation(s.program, gl.Str("proj\x00"))
	gl.UniformMatrix4fv(s.uniProj, 1, false, &proj[0])

	gl.GenBuffers(1, &s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	const stripe = 5 * 4

	posAttr := uint32(gl.GetAttribLocation(s.program, gl.Str("pos\x00")))
	gl.EnableVertexAttribArray(posAttr)
	gl.VertexAttribPointer(posAttr, 2, gl.FLOAT, false, stripe, nil)

	colorAttr := uint32(gl.GetAttribLocation(s.program, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colorAttr)
	gl.VertexAttribPointer(colorAttr, 3, gl.FLOAT, false, stripe,
		gl.PtrOffset(2*4))
}

func (s *PolyShader) Clear() {
	s.points = s.points[:0]
}

func (s *PolyShader) Render() {
	if len(s.points) == 0 {
		return
	}

	gl.BindVertexArray(s.vao)
	gl.UseProgram(s.program)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)

	if len(s.points) > s.vbSize {
		s.vbSize = len(s.points) * 2
		gl.BufferData(gl.ARRAY_BUFFER, s.vbSize*4, nil, gl.STREAM_DRAW)
	}
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(s.points)*4, gl.Ptr(s.points))

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.points)/5))
}

func (s *PolyShader) AddPoints(points []mgl.Vec2, color mgl.Vec3) {
	r, g, b := color.Elem()

	for _, p := range points {
		s.points = append(s.points, p.X(), p.Y(), r, g, b)
	}
}
