package main

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Mandelbrot struct {
	size mgl.Vec2
	trs  mgl.Mat3
}

func NewMandelbrot(width, height float32) *Mandelbrot {
	const (
		x1 = -2
		x2 = 2
		y1 = -1
		y2 = 1

		centerX = x1 + (x2-x1)/2
		centerY = y1 + (y2-y1)/2
	)

	xRatio := (x2 - x1) / width
	yRatio := (y2 - y1) / height
	ratio := Max(xRatio, yRatio)

	mat := mgl.Scale2D(ratio, ratio)
	mat = mgl.Translate2D(centerX, centerY).Mul3(mat)

	m := &Mandelbrot{
		size: mgl.Vec2{width, height},
		trs:  mat,
	}

	return m
}

func (m *Mandelbrot) Draw(renderer *Renderer, cameraTRS mgl.Mat3) {
	mat := m.trs.Mul3(cameraTRS)

	halfX, halfY := m.size.Mul(0.5).Elem()

	for x := float32(0); x < m.size.X(); x++ {
		for y := float32(0); y < m.size.Y(); y++ {
			v := mgl.Vec3{x - halfX, y - halfY, 1}
			x2, y2, _ := mat.Mul3x1(v).Elem()
			c := m.calc(x2, y2)
			renderer.Plot(x, y, mgl.Vec3{c, c, c})
		}
	}
}

func (m *Mandelbrot) calc(x0, y0 float32) float32 {
	esc := 30
	i := 0
	x2, y2 := x0*x0, y0*y0
	for x, y := x0, y0; x2+y2 < 2*2 && i < esc; i++ {
		x, y = x2-y2+x0, 2*x*y+y0
		x2, y2 = x*x, y*y
	}

	return float32(esc-i) / float32(esc)
}

func (m *Mandelbrot) calcSmooth(x0, y0 float32) float32 {
	const (
		maxIter = 200
		limit   = 1 << 16
	)

	i := 0
	x2, y2 := x0*x0, y0*y0
	for x, y := x0, y0; x2+y2 < limit && i < maxIter; i++ {
		x, y = x2-y2+x0, 2*x*y+y0
		x2, y2 = x*x, y*y
	}

	var factor float32
	if i < maxIter {
		zn := math.Log(float64(x2+y2)) / 2
		factor = float32(math.Log(zn/math.Log(2)) / math.Log(2))
	}

	color1 := float32(maxIter-i) / float32(maxIter+1)
	color2 := float32(maxIter-i+1) / float32(maxIter+1)

	return (color2-color1)*factor + color1
}
