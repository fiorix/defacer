package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/fiorix/defacer/apiserver"
)

func main() {
	httpAddr := flag.String("http", ":8080", "[ip]:port to listen on for HTTP")
	apiPrefix := flag.String("api-prefix", "/api", "prefix for API handlers")
	timeout := flag.Duration("timeout", 60*time.Second, "timeout for downloading images")
	nworkers := flag.Uint("workers", 50, "number of defacer workers")
	defaceImage := flag.String("overlay-image", "", "overlay image for the defacer")
	flag.Parse()
	handler := &apiserver.Handler{
		Prefix:    *apiPrefix,
		Workers:   *nworkers,
		ImageFile: *defaceImage,
	}
	cli := &http.Client{Timeout: *timeout}
	log.Println("Starting workers, please wait...")
	if err := handler.Register(http.DefaultServeMux, cli); err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Addr:         *httpAddr,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	log.Println("Starting HTTP server on", *httpAddr)
	log.Fatal(srv.ListenAndServe())
}
