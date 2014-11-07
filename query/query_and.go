package query

import (
	"github.com/juju/errgo"
	"github.com/mezis/klask/index"
	"github.com/mezis/klask/util/tempkey"
)

// A query that intersects results, progressively transforming the set of IDs.
// Represented by a JSON array of generic queries, as the value of an $and key.
type query_and_t struct {
	query_sequence_t
}

func (self *query_and_t) Run(idx index.Index, targetKey string) error {
	conn := idx.Conn()
	defer conn.Close()

	tempkeys, err := tempkey.NewSlice(conn, len(self.queries))
	defer tempkeys.Clear()

	for k, query := range self.queries {
		err := query.Run(idx, tempkeys[k].Name())
		if err != nil {
			return errgo.Mask(err)
		}
	}

	keys := make([]interface{}, len(tempkeys)+1)
	keys[0] = targetKey
	for k, key := range tempkeys {
		keys[k+1] = key.Name()
	}

	_, err = conn.Do("SINTERSTORE", keys...)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
