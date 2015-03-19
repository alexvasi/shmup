package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type NeonShader struct {
	program uint32
	vao     uint32
	vbo     uint32

	buffers [6]Framebuffer
	screen  mgl.Vec2
	result  Framebuffer

	uniStep int32
	uniBlur int32
	uniFade int32
	uniTex  int32
}

const neonVertexSrc string = `
#version 150 core

in vec2 pos;
in vec2 texpos;

out vec2 vTexpos;

void main()
{
    vTexpos = texpos;
    gl_Position = vec4(pos, 0, 1.0);
}
` + "\x00"

const neonFragSrc string = `
#version 150 core

in vec2 vTexpos;

out vec4 outColor;

uniform sampler2D tex;
uniform vec2 step;
uniform bool blur;
uniform bool fade;

void main() {
    vec4 color;

    if (blur) {
        color = 5./16 * texture(tex, vTexpos - step);
        color += 6./16 * texture(tex, vTexpos);
        color += 5./16 * texture(tex, vTexpos + step);
    }
    else if (fade) {
        color = vec4(0, 0, 0, 0.5);
    }
    else {
        color = texture(tex, vTexpos);
    }

    outColor = color;
}
` + "\x00"

func (s *NeonShader) Init(screen mgl.Vec2) {
	s.program = CreateShaderProgram(neonVertexSrc, neonFragSrc)
	s.screen = screen
	for i, size := 0, screen.Mul(0.5); i < len(s.buffers)/2; i++ {
		s.buffers[i*2].Init(size.X(), size.Y())
		s.buffers[i*2+1].Init(size.X(), size.Y())
		size = size.Mul(0.5)
	}
	s.result.Init(screen.X()/2, screen.Y()/2)

	gl.UseProgram(s.program)
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)

	gl.GenBuffers(1, &s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)

	var data = []float32{ // position (x, y), texcoord (u, v)
		-1, 1, 0, 1,
		-1, -1, 0, 0,
		1, 1, 1, 1,
		1, -1, 1, 0,
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)

	posAttr := uint32(gl.GetAttribLocation(s.program, gl.Str("pos\x00")))
	gl.EnableVertexAttribArray(posAttr)
	gl.VertexAttribPointer(posAttr, 2, gl.FLOAT, false, 4*4, nil)

	texAttr := uint32(gl.GetAttribLocation(s.program, gl.Str("texpos\x00")))
	gl.EnableVertexAttribArray(texAttr)
	gl.VertexAttribPointer(texAttr, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	s.uniTex = gl.GetUniformLocation(s.program, gl.Str("tex\x00"))
	s.uniStep = gl.GetUniformLocation(s.program, gl.Str("step\x00"))
	s.uniBlur = gl.GetUniformLocation(s.program, gl.Str("blur\x00"))
	s.uniFade = gl.GetUniformLocation(s.program, gl.Str("fade\x00"))
}

func (s *NeonShader) Clear() {
	s.buffers[0].Bind()
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, int32(s.screen.X()), int32(s.screen.Y()))
}

func (s *NeonShader) BindFramebuffer() {
	s.buffers[0].Bind()
}

func (s *NeonShader) Render() {
	gl.BindVertexArray(s.vao)
	gl.UseProgram(s.program)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	gl.Disable(gl.BLEND)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(s.uniTex, 0)
	gl.Uniform1i(s.uniFade, gl.FALSE)

	// downsample
	gl.Uniform1i(s.uniBlur, gl.FALSE)
	s.buffers[0].BindTexture()
	for i := 1; i < len(s.buffers)/2; i++ {
		s.buffers[i*2].Bind()
		s.draw()
		s.buffers[i*2].BindTexture()
	}

	// blur
	gl.Uniform1i(s.uniBlur, gl.TRUE)
	for i := 1; i < len(s.buffers)/2; i++ {
		bufA, bufB := s.buffers[i*2], s.buffers[i*2+1]
		bufB.Bind()
		bufA.BindTexture()
		gl.Uniform2f(s.uniStep, 1.5/bufA.size.X(), 0)
		s.draw()

		bufA.Bind()
		bufB.BindTexture()
		gl.Uniform2f(s.uniStep, 0, 1.5/bufA.size.Y())
		s.draw()
	}
	gl.Uniform1i(s.uniBlur, gl.FALSE)

	// combine
	s.buffers[0].Bind()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE)
	for i := 1; i < len(s.buffers)/2; i++ {
		s.buffers[i*2].BindTexture()
		s.draw()
	}

	// fade
	s.result.Bind()
	gl.Uniform1i(s.uniFade, gl.TRUE)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	s.draw()
	gl.Uniform1i(s.uniFade, gl.FALSE)
	gl.BlendFunc(gl.ONE, gl.ONE)
	s.buffers[0].BindTexture()
	s.draw()

	// blend to screen
	gl.BlendFunc(gl.ONE, gl.DST_ALPHA)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, int32(s.screen.X()), int32(s.screen.Y()))
	s.result.BindTexture()
	s.draw()
}

func (s *NeonShader) draw() {
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
}
