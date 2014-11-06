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
	tempkeys := make([]index.Tempkey, len(self.queries))

	for k, _ := range tempkeys {
		key, err := idx.NewTempKey()
		if err != nil {
			return errgo.Mask(err)
		}
		tempkeys[k] = key
		defer key.Clear()
	}

	for k, query := range self.queries {
		err := query.Run(idx, tempkeys[k].Get())
		if err != nil {
			return errgo.Mask(err)
		}
	}
	return nil
}
