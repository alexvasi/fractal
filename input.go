package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Input struct {
	Dir         mgl.Vec2
	Scale       float32
	Rotate      float32
	ResetCam    bool
	NextFractal bool
	PrevFractal bool
	NextColor   bool
	PrevColor   bool
	Fullscreen  ToggleSwitch
	LockFrames  ToggleSwitch

	window   *glfw.Window
	joyState [19]bool
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
	i.Dir = mgl.Vec2{0, 0}
	i.Scale = 0
	i.Rotate = 0
	i.ResetCam = false
	i.NextFractal = false
	i.PrevFractal = false
	i.NextColor = false
	i.PrevColor = false
	i.Fullscreen.Toggled = false
	i.LockFrames.Toggled = false

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

	buttons := glfw.GetJoystickButtons(glfw.Joystick1)
	if i.isJoyPressed(10, buttons) || i.isJoyPressed(7, buttons) {
		i.PrevFractal = true
	}
	if i.isJoyPressed(11, buttons) || i.isJoyPressed(5, buttons) {
		i.NextFractal = true
	}
	if i.isJoyPressed(8, buttons) || i.isJoyPressed(6, buttons) {
		i.PrevColor = true
	}
	if i.isJoyPressed(9, buttons) || i.isJoyPressed(4, buttons) {
		i.NextColor = true
	}
	if i.isJoyPressed(12, buttons) {
		i.ResetCam = true
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

func (i *Input) isJoyPressed(btnIdx int, btnState []byte) bool {
	if btnIdx >= len(btnState) || btnIdx >= len(i.joyState) {
		return false
	}

	if btnState[btnIdx] > 0 {
		if i.joyState[btnIdx] {
			return false
		}
		i.joyState[btnIdx] = true
		return true
	}

	i.joyState[btnIdx] = false
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
	case glfw.KeyTab:
		i.NextFractal = true
	case glfw.KeyZ:
		i.PrevColor = true
	case glfw.KeyX:
		i.NextColor = true
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
