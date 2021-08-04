package middleware

import (
	"compress/flate"
	"compress/gzip"
	"net/http"
	"strings"
)


func NewGzipDeflateHandler(next http.Handler) http.Handler {
	return &GZipDeflateHandler{next}
}

type GZipDeflateHandler struct {
	next http.Handler
}

func (h *GZipDeflateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encodings := r.Header.Get("Accept-Encoding")

	if strings.Contains(encodings, "gzip") {
		h.serveGzipped(w, r)
	} else if strings.Contains(encodings, "deflate") {
		h.serveDeflate(w, r)
	} else {
		h.servePlain(w, r)
	}
}

func (h *GZipDeflateHandler) serveGzipped(w http.ResponseWriter, r *http.Request) {
	gzw := gzip.NewWriter(w)
	defer gzw.Close()
	w.Header().Set("Content-Encoding", "gzip")
	h.next.ServeHTTP(GzipResponseWriter{gzw, w}, r)
}

func (h *GZipDeflateHandler) serveDeflate(w http.ResponseWriter, r *http.Request) {
	dfw, _ := flate.NewWriter(w, flate.DefaultCompression)
	defer dfw.Close()
	w.Header().Set("Content-Encoding", "deflate")
	h.next.ServeHTTP(DeflateResponseWriter{dfw, w}, r)
}

func (h *GZipDeflateHandler) servePlain(w http.ResponseWriter, r *http.Request) {
	h.next.ServeHTTP(w, r)
}

type GzipResponseWriter struct {
	gw *gzip.Writer
	http.ResponseWriter
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	if _, ok := w.Header()["Content-Type"]; !ok {
		// If content type is not set, infer it from the uncompressed body.
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return w.gw.Write(b)
}

func (w GzipResponseWriter) Flush() {
	w.gw.Flush()
	if fw, ok := w.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}

type DeflateResponseWriter struct {
	dw *flate.Writer
	http.ResponseWriter
}

func (w DeflateResponseWriter) Write(b []byte) (int, error) {
	if _, ok := w.Header()["Content-Type"]; !ok {
		// If content type is not set, infer it from the uncompressed body.
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return w.dw.Write(b)
}

func (w DeflateResponseWriter) Flush() {
	w.dw.Flush()
	if fw, ok := w.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}
