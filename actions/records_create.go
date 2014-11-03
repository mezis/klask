package actions

import (
	"github.com/gorilla/mux"
	"net/http"
)

func OnRecordsCreate(res http.ResponseWriter, req *http.Request) {
	defer failMessage(res)

	params := mux.Vars(req)
	idx := loadIndex(params["name"])

	record := idx.New()
	requestJson(req, &record)

	exists, err := idx.HasRecord(record.Id())
	abortOn(err)
	if exists {
		fail(http.StatusConflict, "already exists")
	}

	err = record.Persist()
	abortOn(err)

	respondJson(res, http.StatusCreated, record)
}
