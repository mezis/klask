package actions

import (
	"net/http"
)

func OnRootGet(res http.ResponseWriter, req *http.Request) {
	var text = []byte("<h1>Welcome to Klask.</h1>")
	res.Header().Add("Content-Type", "text/html")
	res.Write(text)
}
