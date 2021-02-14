package fractal

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/pprof"
	"sync"
)

// Configuration
const (
	// Position and height
	px = -0.5557506
	py = -0.55560
	ph = 0.000000001
	//px = -2
	//py = -1.2
	//ph = 2.5

	// Quality
	//	imgWidth     = 1024
	//	imgHeight    = 1024
	maxIter      = 1500
	samples      = 50
	linearMixing = true

	showProgress = true
	profileCpu   = false
)

var RemoteWorkers = true

//type RGB struct{ r, g, b int }

// RenderLine computes 1 line of the fractal and can be spun of to a separate engine
func RenderLine(y int, imgHeight, imgWidth int) []int64 {
	ratio := float64(imgWidth) / float64(imgHeight)

	rgb := make([]int64, imgWidth*3)
	for x := 0; x < imgWidth; x++ {
		var r, g, b int64
		for i := 0; i < samples; i++ {
			nx := ph*ratio*((float64(x)+RandFloat64())/float64(imgWidth)) + px
			ny := ph*((float64(y)+RandFloat64())/float64(imgHeight)) + py
			c := paint(mandelbrotIter(nx, ny, maxIter))
			if linearMixing {
				r += int64(RGBToLinear(c.R))
				g += int64(RGBToLinear(c.G))
				b += int64(RGBToLinear(c.B))
			} else {
				r += int64(c.R)
				g += int64(c.G)
				b += int64(c.B)
			}
		}
		// store rgb as 3 int64 values
		rgb[x*3] = r
		rgb[x*3+1] = g
		rgb[x*3+2] = b
	}
	return rgb
}

// addLineToImage takes the colors of 'col' and adds them to the image.
func addLineToImage(img *image.RGBA, y int, rgb []int64) {
	var cr, cg, cb uint8
	var numRgb = len(rgb) / 3
	for x := 0; x < numRgb; x++ {
		r := rgb[x*3]
		g := rgb[x*3+1]
		b := rgb[x*3+2]
		if linearMixing {
			cr = LinearToRGB(uint16(float64(r) / float64(samples)))
			cg = LinearToRGB(uint16(float64(g) / float64(samples)))
			cb = LinearToRGB(uint16(float64(b) / float64(samples)))
		} else {
			cr = uint8(float64(r) / float64(samples))
			cg = uint8(float64(g) / float64(samples))
			cb = uint8(float64(b) / float64(samples))
		}
		img.SetRGBA(x, y, color.RGBA{R: cr, G: cg, B: cb, A: 255})
	}
}

var HostUrl string

func Render(img *image.RGBA, imgWidth, imgHeight int, numWorker int) {
	if profileCpu {
		f, err := os.Create("profile.prof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if HostUrl == "" && RemoteWorkers {
		panic("HostUrl not set while remote workers are expected")
	}

	jobs := make(chan int)

	// create a series of workers
	var wg sync.WaitGroup
	wg.Add(numWorker)
	for i := 0; i < numWorker; i++ {
		go func() {
			defer wg.Done()
			for y := range jobs {
				var rgb []int64
				if RemoteWorkers {
					requestUrl := fmt.Sprintf("%s/GetFractalLine?y=%d&imgWidth=%d&imgHeight=%d", HostUrl, y, imgWidth, imgHeight)
					resp, err := http.Get(requestUrl)
					if err != nil {
						panic(err)
					}
					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						panic(err)
					}
					rgb = BytesToInt64(body)
				} else {
					rgb = RenderLine(y, imgHeight, imgWidth)
				}

				addLineToImage(img, y, rgb)
			}
		}()
	}

	// and send task to the set of workers
	for y := 0; y < imgHeight; y++ {
		jobs <- y
		if showProgress && y%50 == 0 {
			fmt.Printf("\r%d/%d (%d%%)", y, imgHeight, int(100*(float64(y)/float64(imgHeight))))
		}
	}

	// ensure all results have been processed before proceeding.
	close(jobs)
	wg.Wait()
	if showProgress {
		fmt.Printf("\r%d/%[1]d (100%%)\n", imgHeight)
	}
}

func paint(r float64, n int) color.RGBA {
	insideSet := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	if r > 4 {
		return hslToRGB(float64(n)/800*r, 1, 0.5)
	}

	return insideSet
}

func mandelbrotIter(px, py float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x*x, y*y, x*y
		if xx+yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + px
		y = 2*xy + py
	}

	return xx + yy, maxIter
}

// by u/Boraini
//func mandelbrotIterComplex(px, py float64, maxIter int) (float64, int) {
//	var current complex128
//	pxpy := complex(px, py)
//
//	for i := 0; i < maxIter; i++ {
//		magnitude := cmplx.Abs(current)
//		if magnitude > 2 {
//			return magnitude * magnitude, i
//		}
//		current = current * current + pxpy
//	}
//
//	magnitude := cmplx.Abs(current)
//	return magnitude * magnitude, maxIter
//}
