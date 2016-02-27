package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
	size   mgl.Vec2
	screen mgl.Vec2
	ortho  mgl.Mat4

	texShader TexShader
}

func NewRenderer(width, height float32, screenSize mgl.Vec2) *Renderer {
	r := &Renderer{
		size:   mgl.Vec2{width, height},
		screen: screenSize,
		ortho:  mgl.Ortho2D(0, width, 0, height),
	}
	gl.Enable(gl.MULTISAMPLE)

	r.texShader.Init(r.size)

	return r
}

func (r *Renderer) Clear() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (r *Renderer) Render() {
	r.texShader.Render()
}

func (r *Renderer) Plot(x, y float32, color mgl.Vec3) {
	r.texShader.SetPixel(x, y, color)
}
