package query

import (
	"github.com/juju/errgo"
	"github.com/mezis/klask/index"
	"github.com/mezis/klask/util/tempkey"
)

// A query that merges results.
// Represented by a JSON array of generic queries, as the value of an $or key.
type query_or_t struct {
	query_sequence_t
}

func (self *query_or_t) Run(idx index.Index, targetKey string) error {
	err := self.query_sequence_t.Run(idx, "ZUNIONSTORE", targetKey)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}
