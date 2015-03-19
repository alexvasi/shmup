package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Framebuffer struct {
	buf  uint32
	tex  uint32
	size mgl.Vec2
}

func CreateShaderProgram(vertexShader, fragShader string) uint32 {
	vertex := CompileShader(vertexShader, gl.VERTEX_SHADER)
	defer gl.DeleteShader(vertex)

	fragment := CompileShader(fragShader, gl.FRAGMENT_SHADER)
	defer gl.DeleteShader(fragment)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.BindFragDataLocation(program, 0, gl.Str("outColor\x00"))

	LinkShaderProgram(program)
	return program
}

func CompileShader(source string, shaderType uint32) (shader uint32) {
	shader = gl.CreateShader(shaderType)

	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLen)

		logText := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(shader, logLen, nil, gl.Str(logText))

		panic(fmt.Sprintf("Shader compilation error:\n%v\n%v",
			logText, source))
	}

	return shader
}

func LinkShaderProgram(program uint32) {
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLen)

		logText := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(program, logLen, nil, gl.Str(logText))

		panic(fmt.Sprint("Shader program linking error:\n", logText))
	}
}

func CreateFrameBuffer(width, height float32) (frameBuffer, texColor uint32) {
	w, h := int32(width), int32(height)

	gl.GenFramebuffers(1, &frameBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, frameBuffer)

	gl.GenTextures(1, &texColor)
	gl.BindTexture(gl.TEXTURE_2D, texColor)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA,
		gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0,
		gl.TEXTURE_2D, texColor, 0)

	res := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if res != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprint("error creating framebuffer:", res, gl.GetError()))
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	return frameBuffer, texColor
}

func (fb *Framebuffer) Init(width, height float32) {
	fb.buf, fb.tex = CreateFrameBuffer(width, height)
	fb.size[0], fb.size[1] = width, height
}

func (fb *Framebuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb.buf)
	gl.Viewport(0, 0, int32(fb.size.X()), int32(fb.size.Y()))
}

func (fb *Framebuffer) BindTexture() {
	gl.BindTexture(gl.TEXTURE_2D, fb.tex)
}
