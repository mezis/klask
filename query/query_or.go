package query

import (
	"github.com/juju/errgo"
	// "github.com/mezis/klask/index"
)

// A query that merges results.
// Represented by a JSON array of generic queries, as the value of an $or key.
type query_or_t struct {
	query_sequence_t
}

func (self *query_or_t) Run(records string, ctx Context) (string, error) {
	// result, err := self.query_sequence_t.Run("ZUNIONSTORE", records, ctx)
	// if err != nil {
	// 	return "", errgo.Mask(err)
	// }
	// return result, nil
	return "", errgo.New("not implemented")
}
