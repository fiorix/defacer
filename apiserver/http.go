package apiserver

import (
	"image"
	"net/http"
	"os"
	"path"

	"github.com/prometheus/client_golang/prometheus"
)

// Handler provides the defacer HTTP API.
type Handler struct {
	Prefix    string
	ImageFile string
	NWorkers  uint
}

// Register registers the defacer API handlers to the given ServeMux.
func (h *Handler) Register(mux *http.ServeMux) error {
	df, err := h.newDefacer()
	if err != nil {
		return err
	}
	p := path.Clean(path.Join(h.Prefix, "v1"))
	mux.Handle(p+"/metrics", prometheus.Handler())
	mux.Handle(p+"/deface",
		prometheus.InstrumentHandler("deface", DefacerProxy(df, nil)))
	return nil
}

func (h *Handler) newDefacer() (Defacer, error) {
	nworkers := h.NWorkers
	if nworkers == 0 {
		nworkers = 100
	}
	if h.ImageFile == "" {
		return NewDefacerPool(nil, nworkers)
	}
	f, err := os.Open(h.ImageFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	overlay, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return NewDefacerPool(overlay, nworkers)
}
