package handlers

import (
	"encoding/json"
	"net/http"

	"kittens/data"
)

type nameRequest struct {
	Name string `json:"name"`
}

type idRespReq struct {
	Id int `json:"id"`
}

// Search is an http handler for our microservice
type SearchName struct {
	DataStore data.Store
}

func (s *SearchName) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	request := nameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || len(request.Name) < 1 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	kittens := s.DataStore.SearchName(request.Name)
	if len(kittens) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(rw).Encode(kittens)
}

type DeleteAll struct {
	DataStore data.Store
}

func (d *DeleteAll) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	d.DataStore.DeleteAll()
	rw.WriteHeader(http.StatusOK)
}

type DeleteName struct {
	DataStore data.Store
}

func (d *DeleteName) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	request := nameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || len(request.Name) < 1 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if d.DataStore.DeleteName(request.Name) {
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

type Insert struct {
	DataStore data.Store
}

func (d *Insert) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	request := data.Kitten{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Name == "" || request.Id != 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	resp := idRespReq{}
	if resp.Id = d.DataStore.Insert(request); resp.Id != 0 {
		json.NewEncoder(rw).Encode(resp)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

type List struct {
	DataStore data.Store
}

func (s *List) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	kittens := s.DataStore.List()
	if len(kittens) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(rw).Encode(kittens)
}

type DeleteId struct {
	DataStore data.Store
}

func (d *DeleteId) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	request := idRespReq{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Id == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if d.DataStore.DeleteId(request.Id) {
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

type SearchId struct {
	DataStore data.Store
}

func (d *SearchId) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	request := idRespReq{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Id == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	kitten := d.DataStore.SearchId(request.Id)
	if kitten == (data.Kitten{}) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(rw).Encode(&kitten)
}
