package actions

import (
	"encoding/json"
	"github.com/juju/errgo"
	"net/http"
)

type httpError struct {
	status  int
	message string
}

func fail(status int, message string) {
	panic(httpError{status, message})
}

func abortOn(err error) {
	if err != nil {
		panic(err)
	}
}

func failMessage(res http.ResponseWriter) {
	if err := recover(); err != nil {
		switch e := err.(type) {
		case httpError:
			res.Header().Add("Content-Type", "text/plain")
			res.WriteHeader(e.status)
			res.Write([]byte(e.message))
		default:
			panic(e)
		}
	}
}

func requestJson(req *http.Request, resource interface{}) {
	mime := req.Header.Get("Content-Type")
	if mime != "application/json" {
		fail(http.StatusBadRequest, "bad content type")
	}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(resource)
	if err != nil {
		err = errgo.Mask(err)
		panic(err)
	}
}

func respondJson(res http.ResponseWriter, status int, resource interface{}) {
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(status)
	encoder := json.NewEncoder(res)
	err := encoder.Encode(resource)
	if err != nil {
		err = errgo.Mask(err)
		panic(err)
	}
}
