package actions

import (
	"github.com/gorilla/mux"
	"github.com/mezis/klask/config"
	"github.com/mezis/klask/index"
	"net/http"
)

func OnIndicesShow(res http.ResponseWriter, req *http.Request) {
	defer failMessage(res)

	params := mux.Vars(req)
	resource, err := index.New(params["name"], config.Pool())
	abortOn(err)

	exists, err := resource.Exists()
	abortOn(err)
	if !exists {
		fail(http.StatusNotFound, "does not exist")
	}

	err = resource.Load()
	abortOn(err)

	respondJson(res, http.StatusOK, resource)
}
