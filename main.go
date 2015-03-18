package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.0/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	const (
		title  = "Shmup"
		width  = 1366
		height = 768
	)

	defer HandlePanic()

	InitGLFW()
	defer glfw.Terminate()

	window := NewWindow(width, height, title)
	PanicOnError(gl.Init())

	input := NewInput(window)
	renderer := NewRenderer(width, height)

	color := mgl.Vec3{1, 1, 0}

	//glfw.SwapInterval(0)
	timer := NewTimer()
	for !window.ShouldClose() {
		timer.Tick()
		//fmt.Println(timer.DT)
		input.Process()

		renderer.Clear()

		renderer.DrawRect(mgl.Vec2{100, 100}, mgl.Vec2{100, 100}, color)
		renderer.DrawPoly(mgl.Vec2{500, 500}, mgl.Vec2{100, 100}, 5, color)
		renderer.DrawPoly(mgl.Vec2{800, 200}, mgl.Vec2{200, 100}, 30, color)

		renderer.DrawPoly(mgl.Vec2{1000, 500}, mgl.Vec2{100, 100}, 7, color)
		renderer.DrawPoly(mgl.Vec2{1000, 300}, mgl.Vec2{50, 50}, 7, color)
		renderer.DrawPoly(mgl.Vec2{1000, 100}, mgl.Vec2{20, 20}, 7, color)
		renderer.DrawPoly(mgl.Vec2{1000, 50}, mgl.Vec2{10, 10}, 7, color)
		renderer.DrawPoly(mgl.Vec2{1000, 50}, mgl.Vec2{5, 5}, 4, mgl.Vec3{1, 1, 1})

		renderer.DrawPoly(mgl.Vec2{1100, 50}, mgl.Vec2{5, 5}, 7, color)
		renderer.DrawPoly(mgl.Vec2{1100, 100}, mgl.Vec2{5, 5}, 3, color)
		renderer.DrawRect(mgl.Vec2{1100, 150}, mgl.Vec2{3, 3}, color)

		renderer.DrawPoly(mgl.Vec2{1200, 50}, mgl.Vec2{10, 5}, 7, color)
		renderer.DrawRect(mgl.Vec2{1300, 50}, mgl.Vec2{8, 4}, color)

		renderer.Neon()

		renderer.DrawRect(mgl.Vec2{100, 100}, mgl.Vec2{100, 100}, color)
		renderer.DrawPoly(mgl.Vec2{500, 500}, mgl.Vec2{100, 100}, 5, color)
		renderer.DrawPoly(mgl.Vec2{800, 200}, mgl.Vec2{200, 100}, 30, color)

		renderer.DrawPoly(mgl.Vec2{1000, 500}, mgl.Vec2{100, 100}, 7, color)
		renderer.DrawPoly(mgl.Vec2{1000, 300}, mgl.Vec2{50, 50}, 7, color)
		renderer.DrawPoly(mgl.Vec2{1000, 100}, mgl.Vec2{20, 20}, 7, color)
		renderer.DrawPoly(mgl.Vec2{1000, 50}, mgl.Vec2{10, 10}, 7, color)
		renderer.DrawPoly(mgl.Vec2{1000, 50}, mgl.Vec2{5, 5}, 4, mgl.Vec3{1, 1, 1})

		renderer.DrawPoly(mgl.Vec2{1100, 50}, mgl.Vec2{5, 5}, 7, color)
		renderer.DrawPoly(mgl.Vec2{1100, 100}, mgl.Vec2{5, 5}, 3, color)
		renderer.DrawRect(mgl.Vec2{1100, 150}, mgl.Vec2{3, 3}, color)

		renderer.DrawPoly(mgl.Vec2{1200, 50}, mgl.Vec2{10, 5}, 7, color)
		renderer.DrawRect(mgl.Vec2{1300, 50}, mgl.Vec2{8, 4}, color)

		renderer.DrawNeon()

		window.SwapBuffers()
	}

}

func HandlePanic() {
	if err := recover(); err != nil {
		fmt.Println(err)
		glfw.Terminate()
		os.Exit(1)
	}
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
