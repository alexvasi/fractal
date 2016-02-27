package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

func InitGLFW() error {
	err := glfw.Init()
	if err != nil {
		return err
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	//glfw.WindowHint(glfw.RefreshRate, 60)
	glfw.WindowHint(glfw.Samples, 4)

	return nil
}

func NewWindow(w, h int, fullscreen bool, title string) (*glfw.Window, mgl.Vec2) {
	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	if fullscreen {
		w, h = mode.Width, mode.Height
	} else {
		glfw.WindowHint(glfw.Visible, glfw.False)
		monitor = nil
	}

	window, err := glfw.CreateWindow(w, h, title, monitor, nil)
	PanicOnError(err)

	if !fullscreen {
		window.SetPos((mode.Width-w)/2, (mode.Height-h)/2)
		window.Show()
	}

	window.MakeContextCurrent()
	return window, mgl.Vec2{float32(w), float32(h)}
}
