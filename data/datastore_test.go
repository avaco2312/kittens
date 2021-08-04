package data

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"kittens/datatest"

	"github.com/stretchr/testify/assert"
)

func TestDataStore(t *testing.T) {
	for _, tst := range []string{"MEM", "MGO", "CAS", "RED"} {
		store := *initStore(tst)
		var result bool
		for _, tt := range datatest.Testdata {
			if (tt.Expected == http.StatusBadRequest) || (tt.Expected == http.StatusMethodNotAllowed) {
				continue
			}
			switch tt.Url {
			case "/delete/all":
				store.DeleteAll()
				result = true
			case "/list":
				kittens := store.List()
				result = (len(kittens) != 0)
			case "/insert":
				kitten := Kitten{}
				json.NewDecoder(bytes.NewBuffer([]byte(tt.Body.(string)))).Decode(&kitten)
				result = (store.Insert(kitten) != 0)
			case "/search/name":
				kitten := Kitten{}
				json.NewDecoder(bytes.NewBuffer([]byte(tt.Body.(string)))).Decode(&kitten)
				kittens := store.SearchName(kitten.Name)
				result = (len(kittens) != 0)
			case "/search/id":
				kitten := Kitten{}
				json.NewDecoder(bytes.NewBuffer([]byte(tt.Body.(string)))).Decode(&kitten)
				result = (store.SearchId(kitten.Id) != nil)
			case "/delete/name":
				kitten := Kitten{}
				json.NewDecoder(bytes.NewBuffer([]byte(tt.Body.(string)))).Decode(&kitten)
				result = store.DeleteName(kitten.Name)
			case "/delete/id":
				kitten := Kitten{}
				json.NewDecoder(bytes.NewBuffer([]byte(tt.Body.(string)))).Decode(&kitten)
				result = store.DeleteId(kitten.Id)
			}
			assert.Equal(t, tt.Expected == http.StatusOK, result, tst+": "+tt.Name)
		}
	}
}

// Test para concurrencia -race, con nroclientes concurrentes
const nroclientes = 10

func TestDataStoreConcurrente(t *testing.T) {
	for _, tst := range []string{"MEM", "MGO", "CAS", "RED"} {
		store := *initStore(tst)
		var ch = make(chan struct{}, nroclientes)
		for i := 0; i < nroclientes; i++ {
			go func() {
				var result bool
				for _, tt := range datatest.Testdata {
					if (tt.Expected == http.StatusBadRequest) || (tt.Expected == http.StatusMethodNotAllowed) {
						continue
					}
					switch tt.Url {
					case "/delete/all":
						store.DeleteAll()
						result = true
					case "/list":
						kittens := store.List()
						result = (len(kittens) != 0)
					case "/insert":
						kitten := Kitten{}
						json.NewDecoder(bytes.NewBuffer([]byte(tt.Body.(string)))).Decode(&kitten)
						result = (store.Insert(kitten) != 0)
					case "/search/name":
						kitten := Kitten{}
						json.NewDecoder(bytes.NewBuffer([]byte(tt.Body.(string)))).Decode(&kitten)
						kittens := store.SearchName(kitten.Name)
						result = (len(kittens) != 0)
					case "/search/id":
						kitten := Kitten{}
						json.NewDecoder(bytes.NewBuffer([]byte(tt.Body.(string)))).Decode(&kitten)
						result = (store.SearchId(kitten.Id) != nil)
					case "/delete/name":
						kitten := Kitten{}
						json.NewDecoder(bytes.NewBuffer([]byte(tt.Body.(string)))).Decode(&kitten)
						result = store.DeleteName(kitten.Name)
					case "/delete/id":
						kitten := Kitten{}
						json.NewDecoder(bytes.NewBuffer([]byte(tt.Body.(string)))).Decode(&kitten)
						result = store.DeleteId(kitten.Id)
					}
					if nroclientes == 1 {
						assert.Equal(t, tt.Expected == http.StatusOK, result, tst+": "+tt.Name)
					}
				}
				ch <- struct{}{}
			}()
		}
		for i := 0; i < nroclientes; i++ {
			<-ch
		}
	}
}

func initStore(fl string) *Store {
	var store Store
	var err error
	switch fl {
	case "MGO":
		store, err = NewMongoStore("root:example@localhost")
		if err != nil {
			log.Fatal(err)
		}
	case "CAS":
		store, err = NewCassandraStore("localhost")
		if err != nil {
			log.Fatal(err)
		}
	case "RED":
		store, err = NewRediSearchStore("localhost")
		if err != nil {
			log.Fatal(err)
		}
	case "MEM":
		store, _ = NewMemoryStore()
	}
	return &store
}
