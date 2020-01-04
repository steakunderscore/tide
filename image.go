package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

type Image struct {
	Path   string
	Pixels []Pixel
}

func init() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jepg", "jepg", jpeg.Decode, jpeg.DecodeConfig)

}

func NewImage(path string) *Image {
	return &Image{
		Path: path,
	}
}

func (i *Image) ReadFile() error {
	file, err := os.Open(i.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = i.getPixels(file)
	if err != nil {
		return err
	}
	return nil
}

// Get the bi-dimensional pixel array
func (i *Image) getPixels(file io.Reader) error {
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()

	i.Pixels = make([]Pixel, bounds.Max.X*bounds.Max.Y)

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			p := rgbaToPixel(img.At(x, y).RGBA())
			p.X = uint(x)
			p.Y = uint(y)
			i.Pixels = append(i.Pixels, p)
		}
	}

	return nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{
		X: 0, Y: 0,
		R: uint32(r / 257), G: uint32(g / 257), B: uint32(b / 257),
		A: uint32(a / 257)}
}

func (i *Image) Seralise() (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	for _, p := range i.Pixels {
		str := p.String()
		n, err := buf.WriteString(str)
		if n != len(str) || err != nil {
			return nil, err
		}
	}
	return buf, nil
}
