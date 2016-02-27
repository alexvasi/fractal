package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Input struct {
	Dir        mgl.Vec2
	Scale      float32
	Rotate     float32
	ResetCam   bool
	Fullscreen ToggleSwitch
	LockFrames ToggleSwitch

	window *glfw.Window
}

type ToggleSwitch struct {
	Value   bool
	Toggled bool
}

func NewInput(w *glfw.Window, fullscreen bool) *Input {
	i := &Input{
		Fullscreen: ToggleSwitch{Value: fullscreen},
		LockFrames: ToggleSwitch{Value: true},
	}
	i.SetWindow(w)
	return i
}

func (i *Input) Process() {
	i.Fullscreen.Toggled = false
	i.LockFrames.Toggled = false
	i.Dir = mgl.Vec2{0, 0}
	i.Scale = 0
	i.Rotate = 0
	i.ResetCam = false

	glfw.PollEvents()

	axes := glfw.GetJoystickAxes(glfw.Joystick1)
	if len(axes) > 1 {
		i.Dir[0] = axes[0]
		i.Dir[1] = axes[1] * -1
	}

	if len(axes) > 3 {
		i.Rotate = axes[2]
		i.Scale = axes[3] * -1
	}

	if len(axes) > 15 {
		i.Scale += axes[15]
		i.Scale -= axes[14]
	}

	buttons := glfw.GetJoystickButtons(glfw.Joystick1)
	if len(buttons) > 12 {
		if buttons[12] > 0 {
			i.ResetCam = true
		}
	}

	i.handleDirButtons(&i.Dir[0], glfw.KeyA, glfw.KeyD)
	i.handleDirButtons(&i.Dir[1], glfw.KeyS, glfw.KeyW)
	i.handleDirButtons(&i.Scale, glfw.KeyDown, glfw.KeyUp)
	i.handleDirButtons(&i.Scale, glfw.KeyLeftShift, glfw.KeySpace)
	i.handleDirButtons(&i.Rotate, glfw.KeyLeft, glfw.KeyRight)
	i.handleDirButtons(&i.Rotate, glfw.KeyQ, glfw.KeyE)

	i.normalizeDir(&i.Dir)
	i.Scale = mgl.Clamp(i.Scale, -1, 1)

	i.Rotate = mgl.Clamp(i.Rotate, -1, 1)
	if mgl.Abs(i.Rotate) < 0.3 {
		i.Rotate = 0
	}
}

func (i *Input) SetWindow(w *glfw.Window) {
	i.window = w
	w.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	w.SetInputMode(glfw.StickyKeysMode, glfw.False)
	w.SetKeyCallback(i.keyCallback)
}

func (i *Input) isPressed(keys ...glfw.Key) bool {
	for _, key := range keys {
		if i.window.GetKey(key) == glfw.Press {
			return true
		}
	}
	return false
}

func (i *Input) keyCallback(w *glfw.Window, key glfw.Key, scan int,
	action glfw.Action, m glfw.ModifierKey) {

	if action != glfw.Press {
		return
	}

	switch key {
	case glfw.KeyGraveAccent:
		i.LockFrames.Toggle()
	case glfw.KeyF:
		i.Fullscreen.Toggle()
	case glfw.KeyEscape:
		i.window.SetShouldClose(true)
	case glfw.KeyR:
		i.ResetCam = true
	}
}

func (i *Input) normalizeDir(dir *mgl.Vec2) {
	if !dir.ApproxEqual(mgl.Vec2{0, 0}) && dir.Len() > 1 {
		*dir = dir.Normalize()
	}
}

func (i *Input) handleDirButtons(value *float32, decBtn, incBtn glfw.Key) {
	if i.isPressed(incBtn) {
		*value = 1
	} else if i.isPressed(decBtn) {
		*value = -1
	}
}

func (ts *ToggleSwitch) Toggle() bool {
	ts.Value = !ts.Value
	ts.Toggled = true
	return ts.Value
}

func (ts *ToggleSwitch) ToggledTo(v bool) bool {
	return ts.Toggled && ts.Value == v
}
