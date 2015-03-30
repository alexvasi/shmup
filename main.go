package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.0/glfw"
)

func init() {
	runtime.LockOSThread() // GLFW event handling must run on the main OS thread
	rand.Seed(time.Now().UnixNano())
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

	window, screenSize := NewFullScreenWindow(title)
	PanicOnError(gl.Init())

	InitSound()
	defer TerminateSound()
	LoadSoundFile("shoot_human", "shoot_human.wav")
	LoadSoundFile("shoot", "shoot.wav")
	LoadSoundFile("shoot_big", "shoot_big.wav")
	LoadSoundFile("boom", "boom.wav")
	LoadSoundFile("papa", "papa.wav")
	LoadSoundFile("intro", "intro.wav")
	LoadSoundFile("blip", "blip.wav")

	input := NewInput(window)
	renderer := NewRenderer(width, height, screenSize)
	world := NewWorld(width, height)
	game := NewGame(world, input)

	timer := NewTimer(window)
	for !window.ShouldClose() {
		timer.Tick()
		renderer.Clear()

		input.Process()
		timer.ShowTimings(input.debug, 60)

		game.Update(timer.DT)
		world.Draw(renderer)
		renderer.Render()

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
