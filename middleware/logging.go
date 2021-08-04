package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/VividCortex/ewma"
)

func NewLoginHandler(f io.Writer, next http.Handler) http.Handler {
	l := log.New(f, "kittens service: ", log.Ldate|log.Ltime)
	return &LoginHandler{l, next, ewma.NewMovingAverage()}
}

type LoginHandler struct {
	lo   *log.Logger
	Next http.Handler
	ma   ewma.MovingAverage
}

func (l LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lw := NewLoginResponseWriter(w)
	lw.request, _ = ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(lw.request))
	startTime := time.Now()
	l.Next.ServeHTTP(lw, r)
	l.ma.Add(float64(time.Since(startTime)))
	l.lo.Print(r.RemoteAddr, " ", r.Method, " ", r.RequestURI, " ", lw.statusCode, " ", http.StatusText(lw.statusCode), " ", string(lw.request), " ", string(lw.response), " ", l.ma.Value()/1000000)
}

func NewLoginResponseWriter(w http.ResponseWriter) *LoginResponseWriter {
	return &LoginResponseWriter{w, http.StatusOK, []byte{}, []byte{}}
}

type LoginResponseWriter struct {
	http.ResponseWriter
	statusCode int
	request    []byte
	response   []byte
}

func (l *LoginResponseWriter) Write(b []byte) (int, error) {
	l.response = b
	return l.ResponseWriter.Write(b)
}

func (l *LoginResponseWriter) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}
