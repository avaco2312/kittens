package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"kittens/datatest"
)

func BenchmarkSearchHandlerMEM(b *testing.B) {
	handlers := *initStore("MEM")
	for i := 0; i < b.N; i++ {
		for _, tt := range datatest.Testdata {
			var r *http.Request
			rw := httptest.NewRecorder()
			if _, ok := tt.Body.(string); ok {
				r = httptest.NewRequest(tt.Method, tt.Url, bytes.NewBuffer([]byte(tt.Body.(string))))
			} else {
				r = httptest.NewRequest(tt.Method, tt.Url, nil)
			}
			handlers[tt.Url].ServeHTTP(rw, r)
		}
	}
}

func BenchmarkSearchHandlerMGO(b *testing.B) {
	handlers := *initStore("MGO")
	for i := 0; i < b.N; i++ {
		for _, tt := range datatest.Testdata {
			var r *http.Request
			rw := httptest.NewRecorder()
			if _, ok := tt.Body.(string); ok {
				r = httptest.NewRequest(tt.Method, tt.Url, bytes.NewBuffer([]byte(tt.Body.(string))))
			} else {
				r = httptest.NewRequest(tt.Method, tt.Url, nil)
			}
			handlers[tt.Url].ServeHTTP(rw, r)
		}
	}
}

func BenchmarkSearchHandlerCAS(b *testing.B) {
	handlers := *initStore("CAS")
	for i := 0; i < b.N; i++ {
		for _, tt := range datatest.Testdata {
			var r *http.Request
			rw := httptest.NewRecorder()
			if _, ok := tt.Body.(string); ok {
				r = httptest.NewRequest(tt.Method, tt.Url, bytes.NewBuffer([]byte(tt.Body.(string))))
			} else {
				r = httptest.NewRequest(tt.Method, tt.Url, nil)
			}
			handlers[tt.Url].ServeHTTP(rw, r)
		}
	}
}

func BenchmarkSearchHandlerRED(b *testing.B) {
	handlers := *initStore("RED")
	for i := 0; i < b.N; i++ {
		for _, tt := range datatest.Testdata {
			var r *http.Request
			rw := httptest.NewRecorder()
			if _, ok := tt.Body.(string); ok {
				r = httptest.NewRequest(tt.Method, tt.Url, bytes.NewBuffer([]byte(tt.Body.(string))))
			} else {
				r = httptest.NewRequest(tt.Method, tt.Url, nil)
			}
			handlers[tt.Url].ServeHTTP(rw, r)
		}
	}
}
