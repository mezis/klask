package actions

import (
	"github.com/gorilla/mux"
	"github.com/mezis/klask/config"
	"github.com/mezis/klask/index"
	"net/http"
)

func OnIndicesShow(res http.ResponseWriter, req *http.Request) {
	conn := config.Pool().Get()
	defer conn.Close()
	defer failMessage(res)

	params := mux.Vars(req)
	resource, err := index.New(params["name"], conn)
	abortOn(err)

	exists, err := resource.Exists()
	abortOn(err)
	if !exists {
		fail(404, "does not exist")
	}

	err = resource.Load()
	abortOn(err)

	respondJson(res, resource)
}
