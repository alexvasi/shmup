package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
	polyProgram uint32
	polyVA      uint32
	polyVB      uint32
	scrProgram  uint32
	scrVA       uint32

	screen  uint32
	fb1     uint32
	fbTex1  uint32
	fb2     uint32
	fbTex2  uint32
	fb3     uint32
	fbTex3  uint32
	fb4     uint32
	fbTex4  uint32
	fbX1    uint32
	fbXTex1 uint32
	fbX2    uint32
	fbXTex2 uint32
	fbX3    uint32
	fbXTex3 uint32
	fbX4    uint32
	fbXTex4 uint32

	uniModel int32
	uniColor int32
	uniStep  int32
	size     mgl.Vec2
	fbSize   mgl.Vec2
}

const vertexShaderSrc string = `
#version 150 core

in vec2 position;

uniform mat4 ortho;
uniform mat4 model;

void main()
{
    gl_Position = ortho * model * vec4(position, 0, 1);
}
` + "\x00"

const fragShaderSrc string = `
#version 150 core

out vec4 outColor;

uniform vec3 color;

void main() {
    outColor = vec4(color, 1);
}
` + "\x00"

const scrVertexShaderSrc string = `
#version 150

in vec2 position;
in vec2 texpos;

out vec2 vTexpos;

void main()
{
    vTexpos = texpos;
    gl_Position = vec4(position, 0, 1.0);
}
` + "\x00"

const scrFragShaderSrc string = `
#version 150

in vec2 vTexpos;

out vec4 outColor;

uniform sampler2D tex;
uniform vec2 step;

void main() {
    vec4 color = 1./16 * texture(tex, vTexpos - 2*step);
    color += 4./16 * texture(tex, vTexpos - step);
    color += 6./16 * texture(tex, vTexpos);
    color += 4./16 * texture(tex, vTexpos + step);
    color += 1./16 * texture(tex, vTexpos + 2*step);

    // vec4 color = vec4(0.);
    // for (int x = -1; x <= 1; x++) {
    //   color += texture(tex, vTexpos + x*step) / 3;
    // }
    outColor = color;
}
` + "\x00"

func NewRenderer(sizeX, sizeY float32) *Renderer {
	r := &Renderer{
		size:   mgl.Vec2{sizeX, sizeY},
		fbSize: mgl.Vec2{sizeX, sizeY},
	}

	gl.Enable(gl.MULTISAMPLE)
	//gl.Enable(gl.BLEND)
	//gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	//gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	gl.BlendFunc(gl.ONE, gl.DST_ALPHA)

	r.fb1, r.fbTex1 = CreateFrameBuffer(r.fbSize.X(), r.fbSize.Y())
	r.fb2, r.fbTex2 = CreateFrameBuffer(r.fbSize.X()/2, r.fbSize.Y()/2)
	r.fb3, r.fbTex3 = CreateFrameBuffer(r.fbSize.X()/4, r.fbSize.Y()/4)
	r.fb4, r.fbTex4 = CreateFrameBuffer(r.fbSize.X()/8, r.fbSize.Y()/8)
	r.fbX1, r.fbXTex1 = CreateFrameBuffer(r.fbSize.X(), r.fbSize.Y())
	r.fbX2, r.fbXTex2 = CreateFrameBuffer(r.fbSize.X()/2, r.fbSize.Y()/2)
	r.fbX3, r.fbXTex3 = CreateFrameBuffer(r.fbSize.X()/4, r.fbSize.Y()/4)
	r.fbX4, r.fbXTex4 = CreateFrameBuffer(r.fbSize.X()/8, r.fbSize.Y()/8)
	r.polyProgram = CreateProgram(vertexShaderSrc, fragShaderSrc)
	r.scrProgram = CreateProgram(scrVertexShaderSrc, scrFragShaderSrc)

	// poly init
	gl.UseProgram(r.polyProgram)
	gl.GenVertexArrays(1, &r.polyVA)
	gl.BindVertexArray(r.polyVA)

	uniOrtho := gl.GetUniformLocation(r.polyProgram, gl.Str("ortho\x00"))
	ortho := mgl.Ortho2D(0, sizeX, 0, sizeY)
	gl.UniformMatrix4fv(uniOrtho, 1, false, &ortho[0])

	r.uniModel = gl.GetUniformLocation(r.polyProgram, gl.Str("model\x00"))
	r.uniColor = gl.GetUniformLocation(r.polyProgram, gl.Str("color\x00"))

	gl.GenBuffers(1, &r.polyVB)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.polyVB)
	gl.BufferData(gl.ARRAY_BUFFER, 100*4, nil, gl.DYNAMIC_DRAW)

	const stripeSize = 2 * 4 // position (x, y) float32
	posAttrib := uint32(gl.GetAttribLocation(r.polyProgram, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(posAttrib)
	gl.VertexAttribPointer(posAttrib, 2, gl.FLOAT, false, stripeSize, nil)

	// scr init
	gl.UseProgram(r.scrProgram)
	gl.GenVertexArrays(1, &r.scrVA)
	gl.BindVertexArray(r.scrVA)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	var data = []float32{
		// position (x, y), texcoord (u, v)
		-1, 1, 0, 1,
		-1, -1, 0, 0,
		1, 1, 1, 1,
		1, -1, 1, 0,
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)

	posAttrib = uint32(gl.GetAttribLocation(r.scrProgram, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(posAttrib)
	gl.VertexAttribPointer(posAttrib, 2, gl.FLOAT, false, 4*4, nil)

	texAttrib := uint32(gl.GetAttribLocation(r.scrProgram, gl.Str("texpos\x00")))
	gl.EnableVertexAttribArray(texAttrib)
	gl.VertexAttribPointer(texAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	gl.Uniform1i(gl.GetUniformLocation(r.scrProgram, gl.Str("tex\x00")), 0)
	gl.ActiveTexture(gl.TEXTURE0)

	r.uniStep = gl.GetUniformLocation(r.scrProgram, gl.Str("step\x00"))

	return r
}

func (r *Renderer) Clear() {
	gl.ClearColor(0, 0, 0, 0)

	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fb4)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fb3)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fb2)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fbX4)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fbX3)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fbX2)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	//gl.Viewport(0, 0, int32(r.size.X()), int32(r.size.Y()))

	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fbX1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fb1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	//gl.Viewport(0, 0, int32(r.fbSize.X()), int32(r.fbSize.Y()))

	gl.ClearColor(0.2, 0.4, 0.3, 1)
	gl.ClearColor(0, 0, 0, 1)
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.screen)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(0, 0, 0, 0)
	gl.Viewport(0, 0, int32(r.size.X()), int32(r.size.Y()))
}

func (r *Renderer) Neon() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fb2)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Viewport(0, 0, int32(r.fbSize.X()/2), int32(r.fbSize.Y())/2)
}

func (r *Renderer) DrawNeon() {
	// gl.BindFramebuffer(gl.READ_FRAMEBUFFER, r.screen)
	// gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, r.fb1)
	// gl.BlitFramebuffer(0, 0, int32(r.size.X()), int32(r.size.Y()),
	// 	0, 0, int32(r.fbSize.X()), int32(r.fbSize.Y()),
	// 	gl.COLOR_BUFFER_BIT, gl.NEAREST)

	//gl.Disable(gl.BLEND)
	//gl.Enable(gl.BLEND)
	//gl.ClearColor(0, 0, 0, 1)
	//gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindVertexArray(r.scrVA)
	gl.UseProgram(r.scrProgram)
	gl.Uniform2f(r.uniStep, 0, 0)
	gl.Disable(gl.BLEND)

	// gl.BindFramebuffer(gl.FRAMEBUFFER, r.fb2)
	// gl.Viewport(0, 0, int32(r.fbSize.X()/2), int32(r.fbSize.Y()/2))
	// gl.BindTexture(gl.TEXTURE_2D, r.fbTex1)
	// gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fb3)
	gl.Viewport(0, 0, int32(r.fbSize.X()/4), int32(r.fbSize.Y()/4))
	gl.BindTexture(gl.TEXTURE_2D, r.fbTex2)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fb4)
	gl.Viewport(0, 0, int32(r.fbSize.X()/8), int32(r.fbSize.Y()/8))
	gl.BindTexture(gl.TEXTURE_2D, r.fbTex3)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))

	gl.Viewport(0, 0, int32(r.fbSize.X()/8), int32(r.fbSize.Y()/8))
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fbX4)
	gl.Uniform2f(r.uniStep, 8/r.fbSize.X(), 0)
	gl.BindTexture(gl.TEXTURE_2D, r.fbTex4)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fb4)
	gl.Uniform2f(r.uniStep, 0, 8/r.fbSize.Y())
	gl.BindTexture(gl.TEXTURE_2D, r.fbXTex4)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))

	gl.Viewport(0, 0, int32(r.fbSize.X()/4), int32(r.fbSize.Y()/4))
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fbX3)
	gl.Uniform2f(r.uniStep, 4/r.fbSize.X(), 0)
	gl.BindTexture(gl.TEXTURE_2D, r.fbTex3)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fb3)
	gl.Uniform2f(r.uniStep, 0, 4/r.fbSize.Y())
	gl.BindTexture(gl.TEXTURE_2D, r.fbXTex3)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))

	gl.Viewport(0, 0, int32(r.fbSize.X()/2), int32(r.fbSize.Y()/2))
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fbX2)
	gl.Uniform2f(r.uniStep, 2/r.fbSize.X(), 0)
	gl.BindTexture(gl.TEXTURE_2D, r.fbTex2)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fb2)
	gl.Uniform2f(r.uniStep, 0, 2/r.fbSize.Y())
	gl.BindTexture(gl.TEXTURE_2D, r.fbXTex2)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))

	// gl.Viewport(0, 0, int32(r.fbSize.X()), int32(r.fbSize.Y()))
	// gl.BindFramebuffer(gl.FRAMEBUFFER, r.fbX1)
	// gl.Uniform2f(r.uniStep, 1/r.fbSize.X(), 0)
	// gl.BindTexture(gl.TEXTURE_2D, r.fbTex1)
	// gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
	// gl.BindFramebuffer(gl.FRAMEBUFFER, r.fb1)
	// gl.Uniform2f(r.uniStep, 0, 1/r.fbSize.Y())
	// gl.BindTexture(gl.TEXTURE_2D, r.fbXTex1)
	// gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))

	gl.Enable(gl.BLEND)
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.screen)
	gl.Viewport(0, 0, int32(r.size.X()), int32(r.size.Y()))
	gl.Uniform2f(r.uniStep, 0, 0)
	//gl.BindTexture(gl.TEXTURE_2D, r.fbTex1)
	//gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
	gl.BindTexture(gl.TEXTURE_2D, r.fbTex2)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
	gl.BindTexture(gl.TEXTURE_2D, r.fbTex3)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
	gl.BindTexture(gl.TEXTURE_2D, r.fbTex4)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
}

func (r *Renderer) DrawRect(pos, size mgl.Vec2, color mgl.Vec3) {
	gl.BindVertexArray(r.polyVA)
	gl.UseProgram(r.polyProgram)

	model := mgl.Scale3D(size.X()/2, size.Y()/2, 1)
	model = mgl.Translate3D(pos.X(), pos.Y(), 0).Mul4(model)
	gl.UniformMatrix4fv(r.uniModel, 1, false, &model[0])
	gl.Uniform3f(r.uniColor, color[0], color[1], color[2])

	rect := []float32{
		-1, 1,
		-1, -1,
		1, 1,
		1, -1,
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, r.polyVB)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(rect)*4, gl.Ptr(rect))
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
}

func (r *Renderer) DrawPoly(pos, size mgl.Vec2, sides int, color mgl.Vec3) {
	gl.BindVertexArray(r.polyVA)
	gl.UseProgram(r.polyProgram)

	model := mgl.Scale3D(size.X()/2, size.Y()/2, 1)
	model = mgl.Translate3D(pos.X(), pos.Y(), 0).Mul4(model)
	gl.UniformMatrix4fv(r.uniModel, 1, false, &model[0])
	gl.Uniform3f(r.uniColor, color[0], color[1], color[2])

	points := make([]float32, 0, 2*(sides+1))
	points = append(points, 0, 0)
	for s := 0; s <= sides; s++ {
		theta := 2 * math.Pi / float64(sides) * float64(s)
		sin, cos := math.Sincos(theta)
		points = append(points, float32(cos), float32(sin))
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, r.polyVB)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(points)*4, gl.Ptr(points))
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, int32(len(points)/2))
}

func CreateProgram(vertexShader, fragShader string) uint32 {
	vertex := CompileShader(vertexShader, gl.VERTEX_SHADER)
	defer gl.DeleteShader(vertex)

	fragment := CompileShader(fragShader, gl.FRAGMENT_SHADER)
	defer gl.DeleteShader(fragment)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.BindFragDataLocation(program, 0, gl.Str("outColor\x00"))

	LinkProgram(program)
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

func LinkProgram(program uint32) {
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

func CreateFrameBuffer(width, height float32) (frameBuffer, texColorBuffer uint32) {
	w, h := int32(width), int32(height)

	gl.GenFramebuffers(1, &frameBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, frameBuffer)

	gl.GenTextures(1, &texColorBuffer)
	gl.BindTexture(gl.TEXTURE_2D, texColorBuffer)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA,
		gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0,
		gl.TEXTURE_2D, texColorBuffer, 0)

	res := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if res != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprint("error creating framebuffer:", res, gl.GetError()))
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	return frameBuffer, texColorBuffer
}
