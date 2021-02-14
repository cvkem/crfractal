package fractal

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"time"
)

const (
	imgWidth  = 1024
	imgHeight = 1024
)
const (
	textX = 10
	textY = 20
)

func Mandelbrot(w io.Writer, numWorker int) {
	log.Println("Allocating image...")
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	log.Println("Rendering...")
	start := time.Now()
	Render(img, imgWidth, imgHeight, numWorker)
	end := time.Now()

	log.Println("Done rendering in", end.Sub(start))

	log.Println("Adding statistics to image")
	AddLabel(img, textX, textY, fmt.Sprintf("On %d runners rendered in %v", numWorker, end.Sub(start)))

	log.Println("Encoding image...")

	err := png.Encode(w, img)
	if err != nil {
		panic(err)
	}
	log.Println("Done!")
}
