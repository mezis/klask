package actions

import (
	"github.com/mezis/klask/config"
	"github.com/mezis/klask/index"
	"net/http"
)

func OnIndicesIndex(res http.ResponseWriter, req *http.Request) {
	defer failMessage(res)

	indices := make([]index.Index, 0)
	for x := range index.Each(config.Pool()) {
		switch val := x.(type) {
		case index.Index:
			indices = append(indices, val)
		case error:
			abortOn(val)
		default:
			panic(val)
		}
	}

	respondJson(res, http.StatusOK, indices)
}
