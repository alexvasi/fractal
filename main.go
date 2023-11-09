package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func init() {
	runtime.LockOSThread() // GLFW event handling must run on the main OS thread
	rand.Seed(time.Now().UnixNano())
}

func main() {
	const (
		title  = "Fractal"
		width  = 1366
		height = 768
	)

	fullscreen := flag.Bool("fs", false, "fullscreen mode")
	flag.Parse()

	defer HandlePanic()

	PanicOnError(InitGLFW())
	defer glfw.Terminate()

	window, screenSize := NewWindow(width, height, *fullscreen, title)
	PanicOnError(gl.Init())

	input := NewInput(window, *fullscreen)
	renderer := NewRenderer(width, height, screenSize)
	camera := NewCamera(input, 0, 0, width, height)
	fractal := NewFractal(input)
	colors := NewColors(input)
	timer := NewTimer()

	for !window.ShouldClose() {
		renderer.Clear()
		camera.Update(timer.DT)
		fractal.Update()
		colors.Update()
		renderer.SetFractal(camera.Proj(), fractal.Seed(), colors.Data())
		renderer.Render()

		window.SwapBuffers()

		timer.Tick()
		showTimerStat(timer, window, title)

		input.Process()
		switch {
		case input.LockFrames.ToggledTo(true):
			glfw.SwapInterval(1)
		case input.LockFrames.ToggledTo(false):
			glfw.SwapInterval(0)
		case input.Fullscreen.Toggled:
			window.Destroy()
			window, screenSize = NewWindow(width, height,
				input.Fullscreen.Value, title)
			if !input.LockFrames.Value {
				glfw.SwapInterval(0)
			}
			input.SetWindow(window)
			renderer = NewRenderer(width, height, screenSize)
		}
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

func showTimerStat(timer *Timer, window *glfw.Window, title string) {
	if timer.TicksCount == 30 {
		avg := timer.AvgTickTime()
		t := fmt.Sprintf("%s %.2fms %.f", title, avg*1000, 1/avg)
		window.SetTitle(t)
		timer.ResetCounter()
	}
}
