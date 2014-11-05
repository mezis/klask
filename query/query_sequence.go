package query

import (
	"github.com/juju/errgo"
	"github.com/mezis/klask/index"
)

// A query that contains a list of other queries. Used to parse $and and $or
// with the same code.
// Represented with a JSON array.
type query_sequence_t struct {
	queries []Query
}

func (self *query_sequence_t) parse(idx index.Index, parsed interface{}) error {
	switch node := parsed.(type) {
	case []interface{}:
		self.queries = make([]Query, 0, len(node))
		for _, n := range node {
			q := new(query_generic_t)
			if err := q.parse(idx, n); err != nil {
				return errgo.Mask(err)
			}
			self.queries = append(self.queries, q)
		}
		return nil
	default:
		return errgo.Newf("bad subquery of type '%T', expected an array (%v)", node, node)
	}
}
