package actions

import (
	"net/http"
)

func OnIndicesIndex(res http.ResponseWriter, req *http.Request) {
	var data = []byte("{}")
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(500)
	res.Write(data)
}
