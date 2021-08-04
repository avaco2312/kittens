package main

import (
	"bytes"
	"log"
	"net/http"
	"testing"

	"kittens/datatest"

	"github.com/stretchr/testify/assert"
)

func TestKittens(t *testing.T) {
	client := &http.Client{}
	var request *http.Request
	var err error
	for _, tt := range datatest.Testdata {
		if _, ok := tt.Body.(string); ok {
			request, err = http.NewRequest(tt.Method, "http://localhost:8323/kittens"+tt.Url, bytes.NewBuffer([]byte(tt.Body.(string))))
		} else {
			request, err = http.NewRequest(tt.Method, "http://localhost:8323/kittens"+tt.Url, nil)
		}
		if err != nil {
			log.Fatal(err)
		}
		resp, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, tt.Expected, resp.StatusCode, tt.Name)
		resp.Body.Close()
	}
}

// Test para concurrencia -race, con nroclientes concurrentes
const nroclientes = 10

func TestKittensConcurrente(t *testing.T) {
	ch := make(chan bool, nroclientes)
	for i := 0; i < nroclientes; i++ {
		go func() {
			client := &http.Client{}
			var request *http.Request
			var err error
			for _, tt := range datatest.Testdata {
				if _, ok := tt.Body.(string); ok {
					request, err = http.NewRequest(tt.Method, "http://localhost:8323/kittens"+tt.Url, bytes.NewBuffer([]byte(tt.Body.(string))))
				} else {
					request, err = http.NewRequest(tt.Method, "http://localhost:8323/kittens"+tt.Url, nil)
				}
				if err != nil {
					log.Fatal(err)
				}
				resp, err := client.Do(request)
				if err != nil {
					log.Fatal(err)
				}
				if nroclientes == 1 {
					assert.Equal(t, tt.Expected, resp.StatusCode, tt.Name)
				}
				resp.Body.Close()
			}
			ch <- true
		}()
	}
	for i := 0; i < nroclientes; i++ {
		<-ch
	}
}
