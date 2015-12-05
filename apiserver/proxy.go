package apiserver

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
)

// encoderFunc is an adapter function for image encoders.
type encoderFunc func(w io.Writer, m image.Image) error

// proxy is the defacer http proxy.
type proxy struct {
	Defacer  Defacer
	Client   *http.Client
	ErrorLog *log.Logger
}

// DefacerProxy does magic.
func DefacerProxy(df Defacer, cli *http.Client, logger *log.Logger) http.Handler {
	return &proxy{
		Defacer:  df,
		Client:   cli,
		ErrorLog: logger,
	}
}

// ServeHTTP implements the http.Handler interface.
func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	status, err := p.handler(w, r)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
}

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
//
// Copied from net/http/httputil.
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func (p *proxy) handler(w http.ResponseWriter, r *http.Request) (int, error) {
	url := r.FormValue("url")
	if url == "" {
		return http.StatusBadRequest, errors.New("Missing `url` param")
	}
	resp, err := p.req(url, r)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	defer resp.Body.Close()
	// clear response headers
	resp.Header.Del("Content-Length")
	for _, hdr := range hopHeaders {
		resp.Header.Del(hdr)
	}
	// copy response headers
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	var enc encoderFunc
	switch resp.Header.Get("Content-Type") {
	case "image/gif":
		enc = func(w io.Writer, m image.Image) error {
			return gif.Encode(w, m, nil)
		}
	case "image/jpeg":
		enc = func(w io.Writer, m image.Image) error {
			return jpeg.Encode(w, m, nil)
		}
	case "image/png":
		enc = png.Encode
	default:
		io.Copy(w, resp.Body)
		return 0, nil
	}
	img, err := p.Defacer.Deface(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	enc(w, img)
	return 0, nil
}

func (p *proxy) req(url string, r *http.Request) (*http.Response, error) {
	req, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		p.logf("failed to create request to %q: %v", url, err)
		return nil, err
	}
	req.Header.Set("User-Agent", r.Header.Get("User-Agent"))
	resp, err := p.Client.Do(req)
	if err != nil {
		p.logf("failed to exec request to %q: %v", url, err)
		return nil, err
	}
	return resp, nil
}

func (p *proxy) logf(format string, args ...interface{}) {
	if p.ErrorLog != nil {
		p.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
