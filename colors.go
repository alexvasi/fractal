package main

import "github.com/lucasb-eyer/go-colorful"

type ColorInBand struct {
	hex string
	t   float64
}

var SOURCE_COLORS [][]ColorInBand = [][]ColorInBand{
	{
		{"#000000", 0.0},
		{"#ffff00", 0.3},
		{"#000000", 1.0},
	},
	{
		{"#000000", 0.0},
		{"#7c916f", 0.2},
		{"#e8ffd9", 0.3},
		{"#7c916f", 0.4},
		{"#000000", 1.0},
	},
	{
		{"#000000", 0.0},
		{"#df1b8b", 0.2},
		{"#faace6", 0.3},
		{"#df1b8b", 0.4},
		{"#000000", 1.0},
	},
	{
		{"#000100", 0.0},
		{"#44aa00", 0.2},
		{"#f6ffd5", 0.3},
		{"#44aa00", 0.6},
		{"#000100", 1.0},
	},
	{
		{"#000000", 0.0},
		{"#5fd3bc", 0.2},
		{"#aaffee", 0.3},
		{"#5fd3bc", 0.4},
		{"#000000", 1.0},
	},
	{
		{"#000000", 0.0},
		{"#b80000", 0.1},
		{"#f4ac00", 0.2},
		{"#fcfcf4", 0.3},
		{"#f4ac00", 0.4},
		{"#b80000", 0.5},
		{"#000000", 1.0},
	},
}

type Colors struct {
	bands [][]colorful.Color
	data  []uint8
	index int
	input *Input
}

func NewColors(input *Input) *Colors {
	const bandSize = 60

	c := &Colors{
		data:  make([]uint8, bandSize*3),
		input: input,
	}

	for _, source := range SOURCE_COLORS {
		c.addBand(source, bandSize)
	}

	c.calcData()
	return c
}

func (c *Colors) Data() []uint8 {
	return c.data
}

func (c *Colors) Update() {
	if c.input.NextColor {
		c.index += 1
		if c.index+1 > len(c.bands) {
			c.index = 0
		}
		c.calcData()
	} else if c.input.PrevColor {
		c.index -= 1
		if c.index < 0 {
			c.index = len(c.bands) - 1
		}
		c.calcData()
	}
}

func (c *Colors) calcData() {
	for i := range c.bands[c.index] {
		r, g, b := c.bands[c.index][i].RGB255()
		c.data[i*3+0] = r
		c.data[i*3+1] = g
		c.data[i*3+2] = b
	}
}

func (c *Colors) addBand(source []ColorInBand, size int) {
	var band []colorful.Color
	for i, s := 0, 0; i < size; i++ {
		pos := float64(i) / (float64(size) - 1)
		if pos > source[s+1].t {
			s += 1
		}

		from, _ := colorful.Hex(source[s].hex)
		to, _ := colorful.Hex(source[s+1].hex)
		factor := (pos - source[s].t) / (source[s+1].t - source[s].t)

		band = append(band, from.BlendHcl(to, factor).Clamped())
	}

	c.bands = append(c.bands, band)
}

var BANDS [][]struct {
	color string
	t     float32
} = [][]struct {
	color string
	t     float32
}{
	{
		{"#000000", 0.0},
		{"#ffff00", 0.5},
		{"#000000", 1.0},
	},
}
