package actions

import (
	"encoding/json"
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
		switch err.(type) {
		case httpError:
			err := err.(httpError)
			res.Header().Add("Content-Type", "text/plain")
			res.WriteHeader(err.status)
			res.Write([]byte(err.message))
		default:
			panic(err)
		}
	}
}

func requestJson(req *http.Request, resource interface{}) {
	mime := req.Header.Get("Content-Type")
	if mime != "application/json" {
		fail(400, "bad content type")
	}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(resource)
	if err != nil {
		panic(err)
	}
}

func respondJson(res http.ResponseWriter, resource interface{}) {
	res.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	err := encoder.Encode(resource)
	if err != nil {
		panic(err)
	}
}
