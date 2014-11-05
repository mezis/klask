package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/mezis/klask/actions"
	"github.com/mezis/klask/middleware"
)

const (
	serveAddress = ":3001"
)

func main() {
	stack := negroni.New(middleware.NewRecovery(), negroni.NewLogger())
	router := mux.NewRouter()

	router.HandleFunc("/", actions.OnRootGet).Methods("GET")

	router.HandleFunc("/indices", actions.OnIndicesIndex).Methods("GET")
	router.HandleFunc("/indices", actions.OnIndicesCreate).Methods("POST")
	router.HandleFunc("/indices/{name}", actions.OnIndicesShow).Methods("GET")

	router.HandleFunc("/indices/{name}/records", actions.OnRecordsIndex).Methods("GET")
	router.HandleFunc("/indices/{name}/records", actions.OnRecordsCreate).Methods("POST")
	router.HandleFunc("/indices/{name}/records/{id}", actions.OnRecordsDelete).Methods("DELETE")

	stack.UseHandler(router)
	stack.Run(serveAddress)
}
