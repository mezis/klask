package actions

import (
	"github.com/mezis/klask/config"
	"github.com/mezis/klask/index"
	"net/http"
)

func OnIndicesCreate(res http.ResponseWriter, req *http.Request) {
	defer failMessage(res)

	resource, err := index.New("unnamed", config.Pool())
	abortOn(err)

	requestJson(req, &resource)

	exists, err := resource.Exists()
	abortOn(err)
	if exists {
		fail(409, "already exists")
	}

	err = resource.Save()
	abortOn(err)

	res.WriteHeader(201)
	respondJson(res, resource)
}
