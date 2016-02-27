package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type FracShader struct {
	program uint32
	vao     uint32
	vbo     uint32

	size   mgl.Vec2
	scale  float32
	uniCam int32
}

const fracVertexSrc string = `
#version 150 core

in vec2 pos;
in vec2 fracPos;
out vec2 vPos;

uniform mat4 cam;

void main()
{
    vPos = vec2(cam * vec4(fracPos, 0, 1));
    gl_Position = vec4(pos, 0, 1);
}
` + "\x00"

const fracFragSrc string = `
#version 150 core

in vec2 vPos;

out vec4 outColor;

void main() {
    const float esc = 200;
    const float l2 = 0.6931471805599453;

    float i = 0;
    float x = vPos.x, y = vPos.y, xT, yT;
    for (; i < esc && (x*x+y*y) < 65536; i++) {
        xT = x*x-y*y+vPos.x;
        yT = 2*x*y+vPos.y;
        x = xT;
        y = yT;
    }

    float c = 0;
    if (i < esc) {
        float factor = log((log(x*x+y*y)/2)/l2) / l2;
        float c1 = (esc-i)/(esc+1);
        float c2 = (esc-i+1)/(esc+1);
        c = mix(c1, c2, factor);
    }

    outColor = vec4(c, c, c, 1);
}
` + "\x00"

func (s *FracShader) Init(size mgl.Vec2) {
	const (
		width  = 4
		height = 2
	)
	s.scale = Max(width/size.X(), height/size.Y())
	s.size = size

	s.program = CreateShaderProgram(fracVertexSrc, fracFragSrc)

	gl.UseProgram(s.program)
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)

	gl.GenBuffers(1, &s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)

	cam := mgl.Ident4()
	s.uniCam = gl.GetUniformLocation(s.program, gl.Str("cam\x00"))
	gl.UniformMatrix4fv(s.uniCam, 1, false, &cam[0])

	s.initArrayBuffer()
}

func (s *FracShader) Render() {
	gl.BindVertexArray(s.vao)
	gl.UseProgram(s.program)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(4))
}

func (s *FracShader) SetCamera(cam mgl.Mat4) {
	mat := mgl.Scale3D(s.scale, s.scale, 1).Mul4(cam)
	gl.UniformMatrix4fv(s.uniCam, 1, false, &mat[0])
}

func (s *FracShader) initArrayBuffer() {
	x, y := s.size.Mul(0.5).Elem()

	var data = []float32{ // position (x, y), fractal coords
		-1, 1, -x, y,
		-1, -1, -x, -y,
		1, 1, x, y,
		1, -1, x, -y,
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)

	posAttr := uint32(gl.GetAttribLocation(s.program, gl.Str("pos\x00")))
	gl.EnableVertexAttribArray(posAttr)
	gl.VertexAttribPointer(posAttr, 2, gl.FLOAT, false, 4*4, nil)

	texAttr := uint32(gl.GetAttribLocation(s.program, gl.Str("fracPos\x00")))
	gl.EnableVertexAttribArray(texAttr)
	gl.VertexAttribPointer(texAttr, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
}
