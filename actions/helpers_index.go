package actions

import (
	"github.com/mezis/klask/config"
	"github.com/mezis/klask/index"
	"net/http"
)

func loadIndex(name string) index.Index {
	idx, err := index.New(name, config.Pool())
	abortOn(err)

	exist, err := idx.Exists()
	abortOn(err)
	if !exist {
		fail(http.StatusNotFound, "index does not exist")
	}

	err = idx.Load()
	abortOn(err)

	return idx
}
