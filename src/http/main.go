package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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

var port = 8080

var defNumWorker = runtime.NumCPU() * 2

func main() {
	initParams()

	//Create the default mux
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)

	mux.HandleFunc("/GetFractal", MandelbrotHandler)

	mux.HandleFunc("/GetFractalLine", MandelbrotLineHandler)

	//Create the server.
	portStr := fmt.Sprintf(":%d", port)
	s := &http.Server{
		Addr:    portStr,
		Handler: mux,
	}
	log.Println("Starting fractal-http: go to localhost" + portStr)
	s.ListenAndServe()
}

var once sync.Once

func initParams() {
	var err error
	nw := os.Getenv("numWorker")
	if nw != "" {
		if defNumWorker, err = strconv.Atoi(nw); err != nil {
			panic(err)
		}
	}
	p := os.Getenv("port")
	if p != "" {
		if port, err = strconv.Atoi(p); err != nil {
			panic(err)
		}
	}
}

const homePage = `<html>
  <head/>
  <body>
  Usage: https://&lt;host&gt;/GetFractal[?numWorker=10]
  </body>
</html>`

func homeHandler(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(200)
	res.Header().Add("Content-Type", "html")
	res.Write([]byte(homePage))
}

func MandelbrotHandler(res http.ResponseWriter, req *http.Request) {
	var data bytes.Buffer

	// take the hostRequest of the first call to spin up worker tasks
	// only needed for this call
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
	numWorker := getParamOptInt(params, "numWorker", int64(defNumWorker))

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

func getParamOptInt(params url.Values, key string, defVal int64) int64 {
	strVal, present := params[key] //filters=["color", "price", "brand"]
	if present {
		if len(strVal) == 1 {
			i, err := strconv.ParseInt(strVal[0], 10, 64)
			if err != nil {
				panic(err)
			}
			return i
		} else {
			panic(fmt.Sprintf("Key '%s' is present but has %d values", key, len(strVal)))
		}
	}
	return defVal
}
