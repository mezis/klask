package actions

import (
	// "github.com/mezis/klask/config"
	"github.com/gorilla/mux"
	// "github.com/mezis/klask/index"
	"github.com/mezis/klask/query"
	"net/http"
)

func OnRecordsIndex(res http.ResponseWriter, req *http.Request) {
	defer failMessage(res)

	params := mux.Vars(req)
	idx := loadIndex(params["name"])

	query := query.New(idx)
	requestJson(req, &query)

	// ids, err := query.Run(0, 3600)
	// abortOn(err)

	// respondJson(res, http.StatusOk, ids)
}
