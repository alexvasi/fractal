package main

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	TRS mgl.Mat3

	size          mgl.Vec2
	pos           mgl.Vec2
	origPos       mgl.Vec2
	moveSpeed     float32
	scale         float32
	scaleInSpeed  float32
	scaleOutSpeed float32
	angle         float32
	rotateSpeed   float32
	input         *Input
}

func NewCamera(input *Input, x, y, width, height float32) *Camera {
	c := &Camera{
		size:          mgl.Vec2{width, height},
		pos:           mgl.Vec2{x, y},
		origPos:       mgl.Vec2{x, y},
		moveSpeed:     Min(width, height) / 1,
		scale:         1,
		scaleInSpeed:  0.5,
		scaleOutSpeed: 1.5,
		rotateSpeed:   math.Pi,
		input:         input,
	}
	c.calcTRS()

	return c
}

func (c *Camera) Update(dt float32) {
	scaleSpeed := c.scaleInSpeed
	if c.input.Scale < 0 {
		scaleSpeed = c.scaleOutSpeed
	}
	c.scale -= dt * c.input.Scale * c.scale * scaleSpeed
	c.angle = NormalizeAngle(c.angle - dt*c.input.Rotate*c.rotateSpeed)

	p := c.input.Dir.Mul(dt * c.moveSpeed * c.scale)
	c.pos = c.pos.Add(mgl.Rotate2D(c.angle).Mul2x1(p))

	if c.input.ResetCam {
		c.pos = c.origPos
		c.scale = 1
		c.angle = 0
	}

	c.calcTRS()
}

func (c *Camera) calcTRS() {
	mat := mgl.Scale2D(c.scale, c.scale)
	mat = mgl.HomogRotate2D(c.angle).Mul3(mat)
	c.TRS = mgl.Translate2D(c.pos.X(), c.pos.Y()).Mul3(mat)
}