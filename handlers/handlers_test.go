package handlers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"kittens/data"
	"kittens/datatest"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	for _, tst := range []string{"MEM", "MGO", "CAS", "RED"} {
		handlers := *initStore(tst)
		for _, tt := range datatest.Testdata {
			if tt.Expected == http.StatusMethodNotAllowed {
				continue
			}
			var r *http.Request
			rw := httptest.NewRecorder()
			if _, ok := tt.Body.(string); ok {
				r = httptest.NewRequest(tt.Method, tt.Url, bytes.NewBuffer([]byte(tt.Body.(string))))
			} else {
				r = httptest.NewRequest(tt.Method, tt.Url, nil)
			}
			handlers[tt.Url].ServeHTTP(rw, r)
			w := rw.Result()
			assert.Equal(t, tt.Expected, w.StatusCode, tst+": "+tt.Name)
		}
	}
}

// Test para concurrencia -race, con nroclientes concurrentes
const nroclientes = 10

func TestHandlerConcurrente(t *testing.T) {
	for _, tst := range []string{"MEM", "MGO", "CAS", "RED"} {
		handlers := *initStore(tst)
		ch := make(chan bool, nroclientes)
		for i := 0; i < nroclientes; i++ {
			go func() {
				for _, tt := range datatest.Testdata {
					if tt.Expected == http.StatusMethodNotAllowed {
						continue
					}
					var r *http.Request
					rw := httptest.NewRecorder()
					if _, ok := tt.Body.(string); ok {
						r = httptest.NewRequest(tt.Method, tt.Url, bytes.NewBuffer([]byte(tt.Body.(string))))
					} else {
						r = httptest.NewRequest(tt.Method, tt.Url, nil)
					}
					handlers[tt.Url].ServeHTTP(rw, r)
					w := rw.Result()
					if nroclientes == 1 {
						assert.Equal(t, tt.Expected, w.StatusCode, tst+": "+tt.Name)
					}
				}
				ch <- true
			}()
		}
		for i := 0; i < nroclientes; i++ {
			<-ch
		}
	}
}

func initStore(fl string) *map[string]http.Handler {
	var store data.Store
	var err error
	switch fl {
	case "MGO":
		store, err = data.NewMongoStore("root:example@localhost")
		if err != nil {
			log.Fatal(err)
		}
	case "CAS":
		store, err = data.NewCassandraStore("localhost")
		if err != nil {
			log.Fatal(err)
		}
	case "RED":
		store, err = data.NewRediSearchStore("localhost")
		if err != nil {
			log.Fatal(err)
		}
	case "MEM":
		store, _ = data.NewMemoryStore()
	}
	ahandlers := make(map[string]http.Handler)
	ahandlers["/list"] = &List{DataStore: store}
	ahandlers["/search/id"] = &SearchId{DataStore: store}
	ahandlers["/search/name"] = &SearchName{DataStore: store}
	ahandlers["/insert"] = &Insert{DataStore: store}
	ahandlers["/delete/id"] = &DeleteId{DataStore: store}
	ahandlers["/delete/name"] = &DeleteName{DataStore: store}
	ahandlers["/delete/all"] = &DeleteAll{DataStore: store}
	return &ahandlers
}
