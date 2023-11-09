package main

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Fractal struct {
	index int
	seeds []mgl.Vec2
	input *Input
}

func NewFractal(input *Input) *Fractal {
	return &Fractal{
		seeds: []mgl.Vec2{
			{-0.4, 0.6},
			{0.285, 0.01},
			{-0.70176, -0.3842},
			{-0.8, 0.156},
			{-0.7269, 0.1889},
			{0, 0.8},
			{0, 0},
		},
		input: input,
	}
}

func (f *Fractal) Seed() mgl.Vec2 {
	return f.seeds[f.index]
}

func (f *Fractal) Update() {
	if f.input.NextFractal {
		f.index += 1
		if f.index+1 > len(f.seeds) {
			f.index = 0
		}
	} else if f.input.PrevFractal {
		f.index -= 1
		if f.index < 0 {
			f.index = len(f.seeds) - 1
		}
	}
}
