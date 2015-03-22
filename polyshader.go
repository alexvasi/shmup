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

	pointGroups  map[PolyGroup][]float32
	buffer       []float32
	bufferOffset map[PolyGroup]int
	vbSize       int
}

type PolyGroup int

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

	s.pointGroups = make(map[PolyGroup][]float32)
	s.bufferOffset = make(map[PolyGroup]int)
}

func (s *PolyShader) Clear() {
	for group, points := range s.pointGroups {
		s.pointGroups[group] = points[:0]
	}
	s.buffer = s.buffer[:0]
}

func (s *PolyShader) Render(groups ...PolyGroup) {
	gl.BindVertexArray(s.vao)
	gl.UseProgram(s.program)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)

	s.sendPoints()
	if len(groups) == 0 {
		s.draw(0, len(s.buffer))
	}
	for _, group := range groups {
		s.draw(s.bufferOffset[group], len(s.pointGroups[group]))
	}
}

func (s *PolyShader) AddPoints(points []mgl.Vec2, color mgl.Vec3, group PolyGroup) {
	r, g, b := color.Elem()

	for _, p := range points {
		s.pointGroups[group] = append(s.pointGroups[group],
			p.X(), p.Y(), r, g, b)
	}
}

func (s *PolyShader) sendPoints() {
	if len(s.buffer) > 0 {
		return // already sended
	}

	offset := 0
	for group, points := range s.pointGroups {
		s.buffer = append(s.buffer, points...)
		s.bufferOffset[group] = offset
		offset += len(points)
	}

	if len(s.buffer) > s.vbSize {
		s.vbSize = len(s.buffer) * 2
		gl.BufferData(gl.ARRAY_BUFFER, s.vbSize*4, nil, gl.STREAM_DRAW)
	}

	if len(s.buffer) > 0 {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(s.buffer)*4,
			gl.Ptr(s.buffer))
	}
}

func (s *PolyShader) draw(offset, count int) {
	if count > 0 {
		const recSize = 5
		gl.DrawArrays(gl.TRIANGLES, int32(offset/recSize),
			int32(count/recSize))
	}
}
