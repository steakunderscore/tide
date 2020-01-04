package main

import (
	"fmt"
	"strconv"
)

// Pixel struct example
type Pixel struct {
	X uint
	Y uint
	R uint32
	G uint32
	B uint32
	A uint32
}

func (p *Pixel) FromHexRGB(hexRGB string) error {
	if len(hexRGB) != 6 {
		return fmt.Errorf("Incorrect format given, should be like FFFFFF")
	}

	col, err := strconv.ParseInt(hexRGB, 16, 32)
	if err != nil {
		return err
	}

	p.R = uint32((col >> 16) & 0xFF)
	p.G = uint32((col >> 8) & 0xFF)
	p.B = uint32((col >> 0) & 0xFF)
	return nil
}

func (p Pixel) String() string {
	if p.A == 255 {
		return fmt.Sprintf("PX %d %d %02X%02X%02X\n", p.X, p.Y, p.R, p.G, p.B)
	}
	return fmt.Sprintf("PX %d %d %02X%02X%02X%02X\n", p.X, p.Y, p.R, p.G, p.B, p.A)

}

// Location sets the location of this pixel
func (p *Pixel) Location(x, y uint) {
	p.X = x
	p.Y = y
}
