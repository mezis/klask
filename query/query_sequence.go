package query

import (
	"github.com/juju/errgo"
	// "github.com/mezis/klask/index"
)

// A query that contains a list of other queries. Used to parse $and and $or
// with the same code.
// Represented with a JSON array.
type query_sequence_t struct {
	queries []Query
}

func (self *query_sequence_t) parse(parsed interface{}) error {
	switch node := parsed.(type) {
	case []interface{}:
		self.queries = make([]Query, 0, len(node))
		for _, n := range node {
			q := new(query_generic_t)
			if err := q.parse(n); err != nil {
				return errgo.Mask(err)
			}
			self.queries = append(self.queries, q)
		}
		return nil
	default:
		return errgo.Newf("bad subquery of type '%T', expected an array (%v)", node, node)
	}
}

// Operation should be ZINTERSTORE or ZUNIONSTORE
// func (self *query_sequence_t) Run(idx index.Index, operation string, targetKey string) error {
// 	conn := idx.Conn()
// 	defer conn.Close()

// 	N := len(self.queries)

// 	tempkeys, err := tempkey.NewSlice(conn, N)
// 	defer tempkeys.Clear()

// 	for k, query := range self.queries {
// 		err := query.Run(idx, tempkeys[k].Name())
// 		if err != nil {
// 			return errgo.Mask(err)
// 		}
// 	}

// 	keys := make([]interface{}, N+2)
// 	keys[0] = targetKey
// 	keys[1] = N
// 	for k, key := range tempkeys {
// 		keys[k+2] = key.Name()
// 	}

// 	_, err = conn.Do(operation, keys...)
// 	if err != nil {
// 		return errgo.Mask(err)
// 	}

// 	return nil
// }
