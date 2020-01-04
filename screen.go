package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	logger = log.New(os.Stdout, "pixilator: ", log.Lshortfile)
)

// GetScreenSize is a static method to probe a screen for its size
func GetScreenSize(addr string) (uint, uint) {
	screen, err := NewScreen(addr)
	if err != nil {
		logger.Panicf("couldn't open connection to server, %s", err)
	}
	defer screen.Close()
	err = screen.UpdateSize()
	if err != nil {
		logger.Panicf("couldn't get screen size, %s", err)
	}
	return screen.SizeX, screen.SizeY
}

// Screen represents a pixelflut screen
type Screen struct {
	Conn    net.Conn
	Addr    string
	OffsetX uint
	OffsetY uint
	SizeX   uint
	SizeY   uint
}

// NewScreen sets up a new screen ready to be written to
func NewScreen(address string) (*Screen, error) {
	s := &Screen{}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("failed to open connection %s, %s", address, err)
		return s, err
	}
	s.Conn = conn
	return s, nil
}

// SetOffset sets an offset on on the screen
func (s *Screen) SetOffset(x, y uint) error {
	s.OffsetX = x
	s.OffsetY = y
	o := []byte(fmt.Sprintf("OFFSET %d %d", x, y))
	return s.write(o)
}

// WriteData writes the data out to the screen
// Assumes IP fragmentation will do its job for any packets which are too large
func (s *Screen) WriteData(data *bytes.Buffer) error {
	var err error
	err = s.write(data.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (s *Screen) write(b []byte) error {
	c, err := s.Conn.Write(b)
	if err != nil || c != len(b) {
		return fmt.Errorf("error writing, wrote %d of %d, err: %v", c, len(b), err)
	}
	return nil
}

// UpdateSize for the screen we are talking to
func (s *Screen) UpdateSize() error {
	err := s.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return fmt.Errorf("SetReadDeadline failed: %s", err)
	}
	sizeCommand := []byte("SIZE\n")
	n, err := s.Conn.Write(sizeCommand)
	if err != nil || n != len(sizeCommand) {
		return fmt.Errorf("Writing to screen failed: %s", err)
	}
	// returns something like "2560 1920\n"
	logger.Println("Reading size")
	buf := make([]byte, 20)
	n, err = s.Conn.Read(buf)
	if err != nil {
		return err
	}
	logger.Println("Read size")
	if n < 3 {
		return fmt.Errorf("Invalid size given by server: %s", string(buf))
	}
	lines := strings.SplitN(string(buf), "\n", 2)
	sizes := strings.Split(lines[0], " ")
	x, err := strconv.ParseUint(sizes[1], 10, 32)
	if err != nil {
		return fmt.Errorf("Invalid size given by server: %s", err)
	}
	y, err := strconv.ParseUint(sizes[2], 10, 32)
	if err != nil {
		return fmt.Errorf("Invalid size given by server: %s", err)
	}
	s.SizeX = uint(x)
	s.SizeY = uint(y)
	return nil
}

// Close the screen connection
func (s *Screen) Close() error {
	return s.Conn.Close()
}
