package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"time"
)

// SeraliseData generates a static data set to save time computating strings
func SeraliseData(maxX, maxY uint) *bytes.Buffer {
	data := bytes.NewBuffer([]byte{})
	p := Pixel{A: 255}
	p.FromHexRGB(opts.RGBColour)
	var x, y uint
	for x = 0; x < maxX; x++ {
		for y = 0; y < maxY; y++ {
			p.Location(x, y)
			_, err := data.WriteString(p.String())
			if err != nil {
				panic(err)
			}
		}
	}
	return data
}

var opts struct {
	Address   string
	RGBColour string
	ImgFile   string
	XOffset   uint
	YOffset   uint
	Threads   uint

	MaxX uint
	MaxY uint
}

func init() {
	flag.StringVar(&opts.Address, "address", "localhost:1337", "address of server to connect to")
	flag.StringVar(&opts.RGBColour, "rgb-colour", "FFFFFF", "To draw a color on the screen")
	flag.StringVar(&opts.ImgFile, "image", "", "The image you want drawn")
	flag.UintVar(&opts.XOffset, "x-offset", 0, "Offset in X direction")
	flag.UintVar(&opts.YOffset, "y-offset", 0, "Offset in Y direction")
	flag.UintVar(&opts.Threads, "threads", 20, "number of threads to run")
}

func main() {
	flag.Parse()
	opts.MaxX, opts.MaxY = GetScreenSize(opts.Address)
	log.Printf("Running against screen size %dx%dpx", opts.MaxX, opts.MaxY)

	var data *bytes.Buffer
	var err error
	if opts.ImgFile != "" {
		image := NewImage(opts.ImgFile)
		image.ReadFile()
		data, err = image.Seralise()
		if err != nil {
			logger.Fatalf("Failed to read image file: %s", err)
		}
	} else {
		data = SeraliseData(opts.MaxX, opts.MaxY)
	}

	fmt.Printf("Going to write %d pixels with %d bytes of data\n", opts.MaxX*opts.MaxY, data.Len())
	for t := uint(0); t < opts.Threads; t++ {
		fmt.Printf("Starting thread %d\n", t)
		go run(data)
		time.Sleep(5 * time.Millisecond)
	}
	for true {
		time.Sleep(5 * time.Second)
		fmt.Println("running")
	}

}

func run(data *bytes.Buffer) {
	screen, err := NewScreen(opts.Address)
	if err != nil {
		logger.Printf("couldn't open connection to server, %s", err)
		return
	}
	if opts.XOffset != 0 || opts.YOffset != 0 {
		err = screen.SetOffset(opts.XOffset, opts.YOffset)
		if err != nil {
			logger.Fatalf("failed to set offset: %s", err)
		}
	}
	for true {
		err = screen.WriteData(data)
		if err != nil {
			log.Print(err)
		}
	}
}
