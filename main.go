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
	defaceImage := flag.String("overlay-image", "", "overlay image for the defacer")
	nworkers := flag.Uint("workers", 500, "number of defacer workers")
	flag.Parse()
	handler := &apiserver.Handler{
		Prefix:    *apiPrefix,
		ImageFile: *defaceImage,
		NWorkers:  *nworkers,
	}
	if err := handler.Register(http.DefaultServeMux); err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Addr:         *httpAddr,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
