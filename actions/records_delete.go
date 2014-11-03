package actions

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mezis/klask/index"
	"net/http"
)

func OnRecordsDelete(res http.ResponseWriter, req *http.Request) {
	defer failMessage(res)

	params := mux.Vars(req)
	idx := loadIndex(params["name"])

	var id index.Id
	_, err := fmt.Sscanf(params["id"], "%d", &id)
	if err != nil {
		fail(http.StatusBadRequest, "bad id")
	}

	err = idx.Del(id)
	abortOn(err)

	res.WriteHeader(http.StatusNoContent)
}
