package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/mezis/klask/actions"
)

const (
	serveAddress = ":3001"
)

func main() {
	middleware := negroni.Classic()
	router := mux.NewRouter()

	router.HandleFunc("/", actions.OnRootGet).Methods("GET")
	router.HandleFunc("/indices", actions.OnIndicesIndex).Methods("GET")
	router.HandleFunc("/indices", actions.OnIndicesCreate).Methods("POST")
	router.HandleFunc("/indices/{name}", actions.OnIndicesShow).Methods("GET")

	middleware.UseHandler(router)
	middleware.Run(serveAddress)
}
