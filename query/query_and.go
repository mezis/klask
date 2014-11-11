package query

import (
	"github.com/juju/errgo"
)

// A query that intersects results, progressively transforming the set of IDs.
// Represented by a JSON array of generic queries, as the value of an $and key.
type query_and_t struct {
	query_sequence_t
}

func (self *query_and_t) Run(records string, ctx Context) (string, error) {
	conn := ctx.Conn()
	defer conn.Close()

	input := records
	for _, q := range self.queries {
		output, err := q.Run(input, ctx)
		if err != nil {
			return "", errgo.Mask(err)
		}

		if input != records {
			ctx.Keys().Release(input)
		}
		input = output
	}

	// after a run `input` contains either `records` or the last output
	return input, nil
}
