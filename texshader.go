package main

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type TexShader struct {
	program uint32
	vao     uint32
	vbo     uint32

	size    mgl.Vec2
	data    []uint8
	uniTex  int32
	texture uint32
}

const texVertexSrc string = `
#version 150 core

in vec2 pos;
in vec2 texpos;

out vec2 vTexpos;

void main()
{
    vTexpos = texpos;
    gl_Position = vec4(pos, 0, 1.0);
}
` + "\x00"

const texFragSrc string = `
#version 150 core

in vec2 vTexpos;

out vec4 outColor;

uniform sampler2D tex;

void main() {
    outColor = texture(tex, vTexpos);
}
` + "\x00"

func (s *TexShader) Init(size mgl.Vec2) {
	s.size = size
	s.data = make([]uint8, int(size.X()*size.Y())*4)

	s.program = CreateShaderProgram(texVertexSrc, texFragSrc)
	gl.UseProgram(s.program)
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)

	gl.GenBuffers(1, &s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)

	s.uniTex = gl.GetUniformLocation(s.program, gl.Str("tex\x00"))
	gl.Uniform1i(s.uniTex, 0)

	s.initArrayBuffer()
	s.initTexture()
}

func (s *TexShader) Render() {
	gl.BindVertexArray(s.vao)
	gl.UseProgram(s.program)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(s.size.X()),
		int32(s.size.Y()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(s.data))

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
}

func (s *TexShader) SetPixel(x, y float32, color mgl.Vec3) {
	if x >= s.size.X() || y >= s.size.Y() {
		panic(fmt.Sprintf("Pixel coords is out of range: %f, %f", x, y))
	}

	r, g, b := color.Elem()

	idx := int(y*s.size.X()*4) + int(x)*4
	s.data[idx] = uint8(r * 255)
	s.data[idx+1] = uint8(g * 255)
	s.data[idx+2] = uint8(b * 255)
}

func (s *TexShader) initArrayBuffer() {
	var data = []float32{ // position (x, y), texcoord (u, v)
		-1, 1, 0, 1,
		-1, -1, 0, 0,
		1, 1, 1, 1,
		1, -1, 1, 0,
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)

	posAttr := uint32(gl.GetAttribLocation(s.program, gl.Str("pos\x00")))
	gl.EnableVertexAttribArray(posAttr)
	gl.VertexAttribPointer(posAttr, 2, gl.FLOAT, false, 4*4, nil)

	texAttr := uint32(gl.GetAttribLocation(s.program, gl.Str("texpos\x00")))
	gl.EnableVertexAttribArray(texAttr)
	gl.VertexAttribPointer(texAttr, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
}

func (s *TexShader) initTexture() {
	gl.GenTextures(1, &s.texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
}
