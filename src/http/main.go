package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"sync"

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

	//Create the default mux
	mux := http.NewServeMux()

	mux.HandleFunc("/GetFractal", MandelbrotHandler)

	mux.HandleFunc("/GetFractalLine", MandelbrotLineHandler)

	//Create the server.
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Starting fractal-http: go to localhost:8080 ")
	s.ListenAndServe()
}

var once sync.Once

func MandelbrotHandler(res http.ResponseWriter, req *http.Request) {
	var data bytes.Buffer

	// take the hostRequest of the first call to spin up worker tasks
	setHostUrl := func() {
		if req.TLS == nil {
			fractal.HostUrl = "http://" + req.Host
		} else {
			fractal.HostUrl = "https://" + req.Host
		}
		log.Println("Set host-url for remote workers to: ", fractal.HostUrl)
	}
	once.Do(setHostUrl)

	params := req.URL.Query()
	numWorker := getParam(params, "numWorker")

	fractal.Mandelbrot(&data, int(numWorker))

	res.WriteHeader(200)
	res.Header().Add("Content-Type", "image/png")
	res.Write(data.Bytes())
}

func MandelbrotLineHandler(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	y := getParam(params, "y")
	imgWidth := getParam(params, "imgWidth")
	imgHeight := getParam(params, "imgHeight")

	rgb := fractal.RenderLine(int(y), int(imgHeight), int(imgWidth))

	res.WriteHeader(200)
	res.Write(fractal.Int64ToBytes(rgb))

}

func getParam(params url.Values, key string) int64 {
	strVal, present := params[key] //filters=["color", "price", "brand"]
	if !present || len(strVal) == 0 {
		panic(fmt.Sprintf("Manditory key '%s' is missing (or has no values)", key))
	}
	if len(strVal) > 1 {
		panic(fmt.Sprintf("Manditory key '%s' has %d values (expected exactly 1)", key, len(strVal)))
	}
	i, err := strconv.ParseInt(strVal[0], 10, 64)
	if err != nil {
		panic(err)
	}

	return i

}
