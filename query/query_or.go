package query

import (
	"github.com/juju/errgo"
	"github.com/mezis/klask/index"
)

// A query that merges results.
// Represented by a JSON array of generic queries, as the value of an $or key.
type query_or_t struct {
	query_sequence_t
}

func (self *query_or_t) Run(idx index.Index, targetKey string) error {
	return errgo.New("not implemented")
}
