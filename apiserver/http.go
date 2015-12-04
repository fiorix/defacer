package apiserver

import (
	"image"
	"net/http"
	"os"
	"path"

	"github.com/fiorix/defacer/apiserver/internal"
	"github.com/prometheus/client_golang/prometheus"
)

// Handler provides the defacer HTTP API. The zero value of Handler
// is a valid Handler.
type Handler struct {
	Prefix    string // default: "/"
	ImageFile string // default: internal deface image
	Workers   uint   // default: 100
}

// Register registers the defacer API handlers to the given ServeMux.
func (h *Handler) Register(mux *http.ServeMux, cli *http.Client) error {
	if h.Prefix == "" {
		h.Prefix = "/"
	}
	if h.Workers == 0 {
		h.Workers = 100
	}
	df, err := h.newDefacer()
	if err != nil {
		return err
	}
	p := path.Clean(path.Join(h.Prefix, "v1"))
	mux.Handle(p+"/metrics", prometheus.Handler())
	proxy := DefacerProxy(df, cli, nil)
	mux.Handle(p+"/deface", prometheus.InstrumentHandler("deface", proxy))
	return nil
}

// newDefacer creates a defacer pool based on the handler's configuration.
// If ImageFile is empty, we load the default internal deface image.
func (h *Handler) newDefacer() (Defacer, error) {
	var err error
	var overlay image.Image
	switch h.ImageFile {
	case "":
		overlay, err = internal.DefaultDefaceImage()
		if err != nil {
			return nil, err
		}
	default:
		f, err := os.Open(h.ImageFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		overlay, _, err = image.Decode(f)
		if err != nil {
			return nil, err
		}
	}
	return NewDefacerPool(NewImageResizer(overlay), h.Workers)
}
