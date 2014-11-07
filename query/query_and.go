package query

import (
	"github.com/juju/errgo"
	"github.com/mezis/klask/index"
)

// A query that intersects results, progressively transforming the set of IDs.
// Represented by a JSON array of generic queries, as the value of an $and key.
type query_and_t struct {
	query_sequence_t
}

func (self *query_and_t) Run(idx index.Index, targetKey string) error {
	err := self.query_sequence_t.Run(idx, "ZINTERSTORE", targetKey)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}
