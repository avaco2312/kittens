package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"kittens/data"
	"kittens/handlers"
	"kittens/middleware"

	"github.com/gorilla/mux"
)

func main() {
	fl := flag.String("t", "MEM", "Datos a usar (MEM, CAS, MGO, RED) default (MEM)")
	flag.Parse()
	var store data.Store
	var err error
	mongohost := "localhost"
	cassandrahost := "localhost"
	redisearchhost := "localhost"
	switch *fl {
	case "MGO":
		store, err = data.NewMongoStore("root:example@" + mongohost)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Mongo")
	case "CAS":
		store, err = data.NewCassandraStore(cassandrahost)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Cassandra")
	case "RED":
		store, err = data.NewRediSearchStore(redisearchhost)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("RediSearch")
	case "MEM":
		store, _ = data.NewMemoryStore()
		log.Println("Memoria")
	}
	f, _ := os.Create("registro.log")
	defer f.Close()
	router := mux.NewRouter()
	middlewareRouter := middleware.NewLimitHandler(100, middleware.NewGzipDeflateHandler(middleware.NewLoginHandler(f, router)))
	s := router.PathPrefix("/kittens").Subrouter()
	s.Handle("/list", &handlers.List{DataStore: store}).Methods("GET")
	s.Handle("/search/id", &handlers.SearchId{DataStore: store}).Methods("GET")
	s.Handle("/search/name", &handlers.SearchName{DataStore: store}).Methods("GET")
	s.Handle("/insert", &handlers.Insert{DataStore: store}).Methods("POST")
	s.Handle("/delete/id", &handlers.DeleteId{DataStore: store}).Methods("DELETE")
	s.Handle("/delete/name", &handlers.DeleteName{DataStore: store}).Methods("DELETE")
	s.Handle("/delete/all", &handlers.DeleteAll{DataStore: store}).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8323", middlewareRouter))
}
