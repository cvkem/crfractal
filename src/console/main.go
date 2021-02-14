package main

import (
	"os"
	"runtime"

	"github.com/cvkem/crfractal/fractal"
)

var numWorkers int = runtime.NumCPU()

// duplicates
const (
	imgWidth  = 1024
	imgHeight = 1024
)
const (
	textX = 100
	textY = 200
)

func main() {
	f, err := os.Create("result.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fractal.RemoteWorkers = false
	fractal.Mandelbrot(f, numWorkers, "no host-url needed")
}
