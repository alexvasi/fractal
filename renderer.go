package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
	size   mgl.Vec2
	screen mgl.Vec2
	ortho  mgl.Mat4

	fracShader FracShader
}

func NewRenderer(width, height float32, screenSize mgl.Vec2) *Renderer {
	r := &Renderer{
		size:   mgl.Vec2{width, height},
		screen: screenSize,
	}
	gl.Enable(gl.MULTISAMPLE)

	r.fracShader.Init(r.size)

	return r
}

func (r *Renderer) Clear() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (r *Renderer) Render() {
	r.fracShader.Render()
}

func (r *Renderer) SetCamera(mat mgl.Mat4) {
	r.fracShader.SetCamera(mat)
}
