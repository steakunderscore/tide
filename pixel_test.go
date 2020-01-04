package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromHexRGB(t *testing.T) {
	pt := Pixel{}
	err := pt.FromHexRGB("F1F2F3")
	pe := Pixel{R: 0xF1, G: 0xF2, B: 0xF3}
	assert.Nil(t, err)
	assert.Equal(t, pe, pt)
}

func TestString(t *testing.T) {
	p := Pixel{
		X: 1,
		Y: 2,
		R: 3,
		G: 4,
		B: 0xde,
		A: 255,
	}
	assert.Equal(t, "PX 1 2 0304DE\n", p.String())
	p.A = 0x0A
	assert.Equal(t, "PX 1 2 0304DE0A\n", p.String())
}
